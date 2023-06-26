package log

import "testing"

// TestLevelFlagIsValid tests LevelFlag.IsValid method.
func TestLevelFlagIsValid(t *testing.T) {
	type test struct {
		value  LevelFlag
		result bool
	}

	tests := []test{
		{PanicLevel, true},
		{FatalLevel, true},
		{ErrorLevel, true},
		{PanicLevel + PanicLevel, true}, // Panic + Panic == Fatal
		{PanicLevel + FatalLevel, false},
		{LevelFlag(maxLevelConfig + 1), false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok := s.value.IsValid(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigIsValid tests LevelConfig.IsValid method.
func TestLevelConfigIsValid(t *testing.T) {
	type test struct {
		value  LevelConfig
		result bool
	}

	tests := []test{
		{LevelConfig(PanicLevel), true},
		{LevelConfig(FatalLevel), true},
		{LevelConfig(ErrorLevel), true},
		{LevelConfig(PanicLevel + PanicLevel), true},
		{LevelConfig(PanicLevel + PanicLevel + ErrorLevel), true},
		{maxLevelConfig + 1, false},
		{None, true},
		{0, true},
	}

	for i, s := range tests {
		if ok := s.value.IsValid(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigPrivateHas tests LevelConfig.Has method.
func TestLevelConfigPrivateHas(t *testing.T) {
	type test struct {
		value  LevelConfig
		target LevelFlag
		result bool
	}

	tests := []test{
		{LevelConfig(PanicLevel), PanicLevel, true},
		{LevelConfig(PanicLevel), FatalLevel, false},
		{LevelConfig(FatalLevel), FatalLevel, true},
		{LevelConfig(FatalLevel), ErrorLevel, false},
		{LevelConfig(ErrorLevel), ErrorLevel, true},
		{LevelConfig(ErrorLevel), PanicLevel, false},
		{LevelConfig(PanicLevel + FatalLevel), PanicLevel, true},
		{LevelConfig(PanicLevel + FatalLevel), FatalLevel, true},
		{LevelConfig(PanicLevel + FatalLevel), ErrorLevel, false},
		{LevelConfig(PanicLevel + FatalLevel + ErrorLevel), PanicLevel, true},
		{LevelConfig(PanicLevel + FatalLevel + ErrorLevel), FatalLevel, true},
		{LevelConfig(PanicLevel + FatalLevel + ErrorLevel), ErrorLevel, true},
		{LevelConfig(PanicLevel + FatalLevel + ErrorLevel), 0, false},
		{LevelConfig(PanicLevel + ErrorLevel), FatalLevel, false},
		{LevelConfig(PanicLevel + ErrorLevel), ErrorLevel, true},
		{LevelConfig(PanicLevel + ErrorLevel), PanicLevel, true},
		{LevelConfig(PanicLevel + ErrorLevel), None, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Has(s.target); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigSet tests LevelConfig.Set method.
func TestLevelConfigSet(t *testing.T) {
	type test struct {
		value  []LevelFlag
		result LevelConfig
	}

	tests := []test{
		{[]LevelFlag{PanicLevel}, LevelConfig(PanicLevel)},
		{[]LevelFlag{PanicLevel, FatalLevel}, LevelConfig(PanicLevel + FatalLevel)},
		{[]LevelFlag{FatalLevel, ErrorLevel}, LevelConfig(FatalLevel + ErrorLevel)},
		{
			[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel},
			LevelConfig(FatalLevel + ErrorLevel),
		},
	}

	for i, s := range tests {
		var f LevelConfig
		f.Set(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestLevelConfigSetError tests LevelConfig.Set method
// with invalid flag values.
func TestLevelConfigSetError(t *testing.T) {
	type test struct {
		value  []LevelFlag
		result bool
	}

	tests := []test{
		{[]LevelFlag{PanicLevel}, true},
		{[]LevelFlag{PanicLevel, FatalLevel}, true},
		{[]LevelFlag{FatalLevel, ErrorLevel}, true},
		{[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel}, true},
		{[]LevelFlag{LevelFlag(maxLevelConfig) + 1}, false},
		{[]LevelFlag{FatalLevel, None}, false},
	}

	for i, s := range tests {
		var f LevelConfig
		_, err := f.Set(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, err == nil)
		}
	}
}

// TestLevelConfigAdd tests Add method.
func TestLevelConfigAdd(t *testing.T) {
	type test struct {
		def    []LevelFlag
		value  []LevelFlag
		result LevelConfig
	}

	tests := []test{
		{
			[]LevelFlag{PanicLevel},
			[]LevelFlag{PanicLevel},
			LevelConfig(PanicLevel),
		},
		{
			[]LevelFlag{FatalLevel},
			[]LevelFlag{PanicLevel, FatalLevel},
			LevelConfig(PanicLevel + FatalLevel),
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{FatalLevel, ErrorLevel},
			LevelConfig(PanicLevel + FatalLevel + ErrorLevel),
		},
		{
			[]LevelFlag{ErrorLevel, PanicLevel},
			[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel},
			LevelConfig(PanicLevel + FatalLevel + ErrorLevel),
		},
		{
			[]LevelFlag{ErrorLevel, PanicLevel},
			[]LevelFlag{FatalLevel, ErrorLevel, PanicLevel},
			LevelConfig(PanicLevel + FatalLevel + ErrorLevel),
		},
	}

	for i, s := range tests {
		var f LevelConfig
		f.Set(s.def...)
		f.Add(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestLevelConfigAddError tests Add method with invalid flag values.
func TestLevelConfigAddError(t *testing.T) {
	type test struct {
		value  []LevelFlag
		result bool
	}

	tests := []test{
		{[]LevelFlag{PanicLevel}, true},
		{[]LevelFlag{PanicLevel, FatalLevel}, true},
		{[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel}, true},
		{[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel, FatalLevel}, true},
		{[]LevelFlag{FatalLevel, LevelFlag(maxLevelConfig) + 1, FatalLevel}, false},
		{[]LevelFlag{FatalLevel, None, FatalLevel}, false},
		{[]LevelFlag{None}, false},
	}

	for i, s := range tests {
		var f LevelConfig
		_, err := f.Add(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestLevelConfigDelete tests Delete method.
func TestLevelConfigDelete(t *testing.T) {
	type test struct {
		def    []LevelFlag
		value  []LevelFlag
		result LevelConfig
	}

	tests := []test{
		{
			[]LevelFlag{PanicLevel},
			[]LevelFlag{PanicLevel},
			LevelConfig(None),
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{FatalLevel},
			LevelConfig(PanicLevel),
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{FatalLevel, ErrorLevel},
			LevelConfig(PanicLevel),
		},
		{
			[]LevelFlag{ErrorLevel, PanicLevel},
			[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel},
			LevelConfig(PanicLevel),
		},
		{
			[]LevelFlag{FatalLevel, ErrorLevel},
			[]LevelFlag{ErrorLevel, PanicLevel, ErrorLevel},
			LevelConfig(FatalLevel),
		},
		{
			[]LevelFlag{ErrorLevel, PanicLevel},
			[]LevelFlag{FatalLevel, ErrorLevel, PanicLevel},
			LevelConfig(None),
		},
		{
			[]LevelFlag{FatalLevel, ErrorLevel, PanicLevel},
			[]LevelFlag{},
			LevelConfig(FatalLevel + ErrorLevel + PanicLevel),
		},
		{
			[]LevelFlag{FatalLevel, ErrorLevel, PanicLevel},
			[]LevelFlag{FatalLevel},
			LevelConfig(ErrorLevel + PanicLevel),
		},
	}

	for i, s := range tests {
		var f LevelConfig
		f.Set(s.def...)
		f.Delete(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestLevelConfigDeleteError tests Delete method with invalid flag values.
func TestLevelConfigDeleteError(t *testing.T) {
	type test struct {
		value  []LevelFlag
		result bool
	}

	tests := []test{
		{[]LevelFlag{PanicLevel}, true},
		{[]LevelFlag{PanicLevel, FatalLevel}, true},
		{[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel}, true},
		{[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel, FatalLevel}, true},
		{[]LevelFlag{FatalLevel, LevelFlag(maxLevelConfig + 1), FatalLevel}, false},
		{[]LevelFlag{None}, false},
	}

	for i, s := range tests {
		var f LevelConfig
		_, err := f.Delete(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestLevelConfigAll tests All method.
func TestLevelConfigAll(t *testing.T) {
	type test struct {
		def    []LevelFlag
		value  []LevelFlag
		result bool
	}

	tests := []test{
		{
			[]LevelFlag{PanicLevel},
			[]LevelFlag{PanicLevel},
			true,
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{FatalLevel},
			true,
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{ErrorLevel},
			false,
		},
		{
			[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel},
			[]LevelFlag{ErrorLevel, PanicLevel},
			false,
		},
		{
			[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel, None},
			[]LevelFlag{ErrorLevel, PanicLevel},
			false,
		},
	}

	for i, s := range tests {
		var f LevelConfig
		f.Set(s.def...)
		if ok, _ := f.All(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigAny tests Any method.
func TestLevelConfigAny(t *testing.T) {
	type test struct {
		def    []LevelFlag
		value  []LevelFlag
		result bool
	}

	tests := []test{
		{
			[]LevelFlag{PanicLevel},
			[]LevelFlag{PanicLevel},
			true,
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{FatalLevel},
			true,
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{ErrorLevel},
			false,
		},
		{
			[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel},
			[]LevelFlag{ErrorLevel, PanicLevel},
			true,
		},
		{
			[]LevelFlag{FatalLevel, ErrorLevel, FatalLevel},
			[]LevelFlag{PanicLevel, PanicLevel, ErrorLevel, PanicLevel},
			true,
		},
		{
			[]LevelFlag{PanicLevel, FatalLevel},
			[]LevelFlag{ErrorLevel, PanicLevel},
			true,
		},
	}

	for i, s := range tests {
		var f LevelConfig
		f.Set(s.def...)
		if ok, _ := f.Any(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigPanic tests LevelConfig.Panic method.
func TestLevelConfigPanic(t *testing.T) {
	type test struct {
		value  LevelConfig
		result bool
	}

	tests := []test{
		{LevelConfig(PanicLevel), true},
		{LevelConfig(FatalLevel), false},
		{LevelConfig(ErrorLevel), false},
		{LevelConfig(PanicLevel + FatalLevel), true},
		{LevelConfig(FatalLevel + ErrorLevel), false},
		{maxLevelConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Panic(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigFatal tests LevelConfig.Fatal method.
func TestLevelConfigFatal(t *testing.T) {
	type test struct {
		value  LevelConfig
		result bool
	}

	tests := []test{
		{LevelConfig(PanicLevel), false},
		{LevelConfig(FatalLevel), true},
		{LevelConfig(ErrorLevel), false},
		{LevelConfig(PanicLevel + FatalLevel), true},
		{LevelConfig(FatalLevel + ErrorLevel), true},
		{LevelConfig(PanicLevel + ErrorLevel), false},
		{maxLevelConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Fatal(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigError tests LevelConfig.Error method.
func TestLevelConfigError(t *testing.T) {
	type test struct {
		value  LevelConfig
		result bool
	}

	tests := []test{
		{LevelConfig(PanicLevel), false},
		{LevelConfig(FatalLevel), false},
		{LevelConfig(ErrorLevel), true},
		{LevelConfig(PanicLevel + FatalLevel), false},
		{LevelConfig(FatalLevel + ErrorLevel), true},
		{LevelConfig(ErrorLevel + PanicLevel), true},
		{maxLevelConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Error(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigInfo tests LevelConfig.Info method.
func TestLevelConfigInfo(t *testing.T) {
	type test struct {
		value  LevelConfig
		result bool
	}

	tests := []test{
		{LevelConfig(InfoLevel), true},
		{LevelConfig(FatalLevel), false},
		{LevelConfig(DebugLevel), false},
		{LevelConfig(PanicLevel + FatalLevel + DebugLevel), false},
		{LevelConfig(FatalLevel + ErrorLevel + InfoLevel), true},
		{LevelConfig(ErrorLevel + InfoLevel + PanicLevel), true},
		{maxLevelConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Info(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigDebug tests LevelConfig.Debug method.
func TestLevelConfigDebug(t *testing.T) {
	type test struct {
		value  LevelConfig
		result bool
	}

	tests := []test{
		{LevelConfig(InfoLevel), false},
		{LevelConfig(FatalLevel), false},
		{LevelConfig(DebugLevel), true},
		{LevelConfig(PanicLevel + FatalLevel + DebugLevel), true},
		{LevelConfig(FatalLevel + ErrorLevel + InfoLevel), false},
		{LevelConfig(ErrorLevel + InfoLevel + PanicLevel + DebugLevel), true},
		{maxLevelConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Debug(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelConfigTrace tests LevelConfig.Trace method.
func TestLevelConfigTrace(t *testing.T) {
	type test struct {
		value  LevelConfig
		result bool
	}

	tests := []test{
		{LevelConfig(InfoLevel), false},
		{LevelConfig(FatalLevel), false},
		{LevelConfig(TraceLevel), true},
		{LevelConfig(TraceLevel + FatalLevel + DebugLevel), true},
		{LevelConfig(FatalLevel + ErrorLevel + InfoLevel), false},
		{LevelConfig(ErrorLevel + TraceLevel + PanicLevel + DebugLevel), true},
		{maxLevelConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Trace(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}
