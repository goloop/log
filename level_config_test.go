package log

import "testing"

// TestLevelFlagIsValid tests LevelFlag.IsValid method.
func TestLevelFlagIsValid(t *testing.T) {
	type test struct {
		value  LevelFlag
		result bool
	}

	var tests = []test{
		{Panic, true},
		{Fatal, true},
		{Error, true},
		{Panic + Panic, true}, // Panic + Panic == Fatal
		{Panic + Fatal, false},
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

	var tests = []test{
		{LevelConfig(Panic), true},
		{LevelConfig(Fatal), true},
		{LevelConfig(Error), true},
		{LevelConfig(Panic + Panic), true},
		{LevelConfig(Panic + Panic + Error), true},
		{LevelConfig(maxLevelConfig + 1), false},
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

	var tests = []test{
		{LevelConfig(Panic), Panic, true},
		{LevelConfig(Panic), Fatal, false},
		{LevelConfig(Fatal), Fatal, true},
		{LevelConfig(Fatal), Error, false},
		{LevelConfig(Error), Error, true},
		{LevelConfig(Error), Panic, false},
		{LevelConfig(Panic + Fatal), Panic, true},
		{LevelConfig(Panic + Fatal), Fatal, true},
		{LevelConfig(Panic + Fatal), Error, false},
		{LevelConfig(Panic + Fatal + Error), Panic, true},
		{LevelConfig(Panic + Fatal + Error), Fatal, true},
		{LevelConfig(Panic + Fatal + Error), Error, true},
		{LevelConfig(Panic + Fatal + Error), 0, false},
		{LevelConfig(Panic + Error), Fatal, false},
		{LevelConfig(Panic + Error), Error, true},
		{LevelConfig(Panic + Error), Panic, true},
		{LevelConfig(Panic + Error), None, false},
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

	var tests = []test{
		{[]LevelFlag{Panic}, LevelConfig(Panic)},
		{[]LevelFlag{Panic, Fatal}, LevelConfig(Panic + Fatal)},
		{[]LevelFlag{Fatal, Error}, LevelConfig(Fatal + Error)},
		{
			[]LevelFlag{Fatal, Error, Fatal},
			LevelConfig(Fatal + Error),
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

	var tests = []test{
		{[]LevelFlag{Panic}, true},
		{[]LevelFlag{Panic, Fatal}, true},
		{[]LevelFlag{Fatal, Error}, true},
		{[]LevelFlag{Fatal, Error, Fatal}, true},
		{[]LevelFlag{LevelFlag(maxLevelConfig) + 1}, false},
		{[]LevelFlag{Fatal, None}, false},
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

	var tests = []test{
		{
			[]LevelFlag{Panic},
			[]LevelFlag{Panic},
			LevelConfig(Panic),
		},
		{
			[]LevelFlag{Fatal},
			[]LevelFlag{Panic, Fatal},
			LevelConfig(Panic + Fatal),
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Fatal, Error},
			LevelConfig(Panic + Fatal + Error),
		},
		{
			[]LevelFlag{Error, Panic},
			[]LevelFlag{Fatal, Error, Fatal},
			LevelConfig(Panic + Fatal + Error),
		},
		{
			[]LevelFlag{Error, Panic},
			[]LevelFlag{Fatal, Error, Panic},
			LevelConfig(Panic + Fatal + Error),
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

	var tests = []test{
		{[]LevelFlag{Panic}, true},
		{[]LevelFlag{Panic, Fatal}, true},
		{[]LevelFlag{Fatal, Error, Fatal}, true},
		{[]LevelFlag{Fatal, Error, Fatal, Fatal}, true},
		{[]LevelFlag{Fatal, LevelFlag(maxLevelConfig) + 1, Fatal}, false},
		{[]LevelFlag{Fatal, None, Fatal}, false},
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

	var tests = []test{
		{
			[]LevelFlag{Panic},
			[]LevelFlag{Panic},
			LevelConfig(None),
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Fatal},
			LevelConfig(Panic),
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Fatal, Error},
			LevelConfig(Panic),
		},
		{
			[]LevelFlag{Error, Panic},
			[]LevelFlag{Fatal, Error, Fatal},
			LevelConfig(Panic),
		},
		{
			[]LevelFlag{Fatal, Error},
			[]LevelFlag{Error, Panic, Error},
			LevelConfig(Fatal),
		},
		{
			[]LevelFlag{Error, Panic},
			[]LevelFlag{Fatal, Error, Panic},
			LevelConfig(None),
		},
		{
			[]LevelFlag{Fatal, Error, Panic},
			[]LevelFlag{},
			LevelConfig(Fatal + Error + Panic),
		},
		{
			[]LevelFlag{Fatal, Error, Panic},
			[]LevelFlag{Fatal},
			LevelConfig(Error + Panic),
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

	var tests = []test{
		{[]LevelFlag{Panic}, true},
		{[]LevelFlag{Panic, Fatal}, true},
		{[]LevelFlag{Fatal, Error, Fatal}, true},
		{[]LevelFlag{Fatal, Error, Fatal, Fatal}, true},
		{[]LevelFlag{Fatal, LevelFlag(maxLevelConfig + 1), Fatal}, false},
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

	var tests = []test{
		{
			[]LevelFlag{Panic},
			[]LevelFlag{Panic},
			true,
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Fatal},
			true,
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Error},
			false,
		},
		{
			[]LevelFlag{Fatal, Error, Fatal},
			[]LevelFlag{Error, Panic},
			false,
		},
		{
			[]LevelFlag{Fatal, Error, Fatal, None},
			[]LevelFlag{Error, Panic},
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

	var tests = []test{
		{
			[]LevelFlag{Panic},
			[]LevelFlag{Panic},
			true,
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Fatal},
			true,
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Error},
			false,
		},
		{
			[]LevelFlag{Fatal, Error, Fatal},
			[]LevelFlag{Error, Panic},
			true,
		},
		{
			[]LevelFlag{Fatal, Error, Fatal},
			[]LevelFlag{Panic, Panic, Error, Panic},
			true,
		},
		{
			[]LevelFlag{Panic, Fatal},
			[]LevelFlag{Error, Panic},
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

	var tests = []test{
		{LevelConfig(Panic), true},
		{LevelConfig(Fatal), false},
		{LevelConfig(Error), false},
		{LevelConfig(Panic + Fatal), true},
		{LevelConfig(Fatal + Error), false},
		{LevelConfig(maxLevelConfig + 1), false},
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

	var tests = []test{
		{LevelConfig(Panic), false},
		{LevelConfig(Fatal), true},
		{LevelConfig(Error), false},
		{LevelConfig(Panic + Fatal), true},
		{LevelConfig(Fatal + Error), true},
		{LevelConfig(Panic + Error), false},
		{LevelConfig(maxLevelConfig + 1), false},
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

	var tests = []test{
		{LevelConfig(Panic), false},
		{LevelConfig(Fatal), false},
		{LevelConfig(Error), true},
		{LevelConfig(Panic + Fatal), false},
		{LevelConfig(Fatal + Error), true},
		{LevelConfig(Error + Panic), true},
		{LevelConfig(maxLevelConfig + 1), false},
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

	var tests = []test{
		{LevelConfig(Info), true},
		{LevelConfig(Fatal), false},
		{LevelConfig(Debug), false},
		{LevelConfig(Panic + Fatal + Debug), false},
		{LevelConfig(Fatal + Error + Info), true},
		{LevelConfig(Error + Info + Panic), true},
		{LevelConfig(maxLevelConfig + 1), false},
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

	var tests = []test{
		{LevelConfig(Info), false},
		{LevelConfig(Fatal), false},
		{LevelConfig(Debug), true},
		{LevelConfig(Panic + Fatal + Debug), true},
		{LevelConfig(Fatal + Error + Info), false},
		{LevelConfig(Error + Info + Panic + Debug), true},
		{LevelConfig(maxLevelConfig + 1), false},
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

	var tests = []test{
		{LevelConfig(Info), false},
		{LevelConfig(Fatal), false},
		{LevelConfig(Trace), true},
		{LevelConfig(Trace + Fatal + Debug), true},
		{LevelConfig(Fatal + Error + Info), false},
		{LevelConfig(Error + Trace + Panic + Debug), true},
		{LevelConfig(maxLevelConfig + 1), false},
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
