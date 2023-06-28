package layout

/*
import "testing"

// TestFormatFlagIsValid tests FormatFlag.IsValid method.
func TestFormatFlagIsValid(t *testing.T) {
	type test struct {
		value  FormatFlag
		result bool
	}

	tests := []test{
		{FullPathFormat, true},
		{FuncNameFormat, true},
		{LineNumberFormat, true},
		{FullPathFormat + FullPathFormat, true}, // FilePath + FilePath == FuncName
		{FullPathFormat + FuncNameFormat, false},
		{FormatFlag(maxFormatConfig + 1), false},
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

// TestFormatConfigIsValid tests FormatConfig.IsValid method.
func TestFormatConfigIsValid(t *testing.T) {
	type test struct {
		value  FormatConfig
		result bool
	}

	tests := []test{
		{FormatConfig(FullPathFormat), true},
		{FormatConfig(FuncNameFormat), true},
		{FormatConfig(LineNumberFormat), true},
		{FormatConfig(FullPathFormat + FullPathFormat), true},
		{FormatConfig(FullPathFormat + FullPathFormat + LineNumberFormat), true},
		{maxFormatConfig + 1, false},
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

// TestFormatConfigPrivateHas tests FormatConfig.Has method.
func TestFormatConfigPrivateHas(t *testing.T) {
	type test struct {
		value  FormatConfig
		target FormatFlag
		result bool
	}

	tests := []test{
		{FormatConfig(FullPathFormat), FullPathFormat, true},
		{FormatConfig(FullPathFormat), FuncNameFormat, false},
		{FormatConfig(FuncNameFormat), FuncNameFormat, true},
		{FormatConfig(FuncNameFormat), LineNumberFormat, false},
		{FormatConfig(LineNumberFormat), LineNumberFormat, true},
		{FormatConfig(LineNumberFormat), FullPathFormat, false},
		{FormatConfig(FullPathFormat + FuncNameFormat), FullPathFormat, true},
		{FormatConfig(FullPathFormat + FuncNameFormat), FuncNameFormat, true},
		{FormatConfig(FullPathFormat + FuncNameFormat), LineNumberFormat, false},
		{FormatConfig(FullPathFormat + FuncNameFormat + LineNumberFormat), FullPathFormat, true},
		{FormatConfig(FullPathFormat + FuncNameFormat + LineNumberFormat), FuncNameFormat, true},
		{FormatConfig(FullPathFormat + FuncNameFormat + LineNumberFormat), LineNumberFormat, true},
		{FormatConfig(FullPathFormat + FuncNameFormat + LineNumberFormat), 0, false},
		{FormatConfig(FullPathFormat + LineNumberFormat), FuncNameFormat, false},
		{FormatConfig(FullPathFormat + LineNumberFormat), LineNumberFormat, true},
		{FormatConfig(FullPathFormat + LineNumberFormat), FullPathFormat, true},
		{FormatConfig(FullPathFormat + LineNumberFormat), None, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Has(s.target); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatConfigSet tests FormatConfig.Set method.
func TestFormatConfigSet(t *testing.T) {
	type test struct {
		value  []FormatFlag
		result FormatConfig
	}

	tests := []test{
		{[]FormatFlag{FullPathFormat}, FormatConfig(FullPathFormat)},
		{[]FormatFlag{FullPathFormat, FuncNameFormat}, FormatConfig(FullPathFormat + FuncNameFormat)},
		{[]FormatFlag{FuncNameFormat, LineNumberFormat}, FormatConfig(FuncNameFormat + LineNumberFormat)},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat},
			FormatConfig(FuncNameFormat + LineNumberFormat),
		},
	}

	for i, s := range tests {
		var f FormatConfig
		f.Set(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestFormatConfigSetError tests FormatConfig.Set method with invalid flag values.
func TestFormatConfigSetError(t *testing.T) {
	type test struct {
		value  []FormatFlag
		result bool
	}

	tests := []test{
		{[]FormatFlag{FullPathFormat}, true},
		{[]FormatFlag{FullPathFormat, FuncNameFormat}, true},
		{[]FormatFlag{FuncNameFormat, LineNumberFormat}, true},
		{[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat}, true},
		{[]FormatFlag{FormatFlag(maxFormatConfig) + 1}, false},
		{[]FormatFlag{FuncNameFormat, None}, false},
	}

	for i, s := range tests {
		var f FormatConfig
		_, err := f.Set(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, err == nil)
		}
	}
}

// TestFormatConfigAdd tests Add method.
func TestFormatConfigAdd(t *testing.T) {
	type test struct {
		def    []FormatFlag
		value  []FormatFlag
		result FormatConfig
	}

	tests := []test{
		{
			[]FormatFlag{FullPathFormat},
			[]FormatFlag{FullPathFormat},
			FormatConfig(FullPathFormat),
		},
		{
			[]FormatFlag{FuncNameFormat},
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			FormatConfig(FullPathFormat + FuncNameFormat),
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{FuncNameFormat, LineNumberFormat},
			FormatConfig(FullPathFormat + FuncNameFormat + LineNumberFormat),
		},
		{
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat},
			FormatConfig(FullPathFormat + FuncNameFormat + LineNumberFormat),
		},
		{
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FullPathFormat},
			FormatConfig(FullPathFormat + FuncNameFormat + LineNumberFormat),
		},
	}

	for i, s := range tests {
		var f FormatConfig
		f.Set(s.def...)
		f.Add(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestFormatConfigAddError tests Add method with invalid flag values.
func TestFormatConfigAddError(t *testing.T) {
	type test struct {
		value  []FormatFlag
		result bool
	}

	tests := []test{
		{[]FormatFlag{FullPathFormat}, true},
		{[]FormatFlag{FullPathFormat, FuncNameFormat}, true},
		{[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat}, true},
		{[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat, FuncNameFormat}, true},
		{[]FormatFlag{FuncNameFormat, FormatFlag(maxFormatConfig) + 1, FuncNameFormat}, false},
		{[]FormatFlag{FuncNameFormat, None, FuncNameFormat}, false},
		{[]FormatFlag{None}, false},
	}

	for i, s := range tests {
		var f FormatConfig
		_, err := f.Add(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestFormatConfigDelete tests Delete method.
func TestFormatConfigDelete(t *testing.T) {
	type test struct {
		def    []FormatFlag
		value  []FormatFlag
		result FormatConfig
	}

	tests := []test{
		{
			[]FormatFlag{FullPathFormat},
			[]FormatFlag{FullPathFormat},
			FormatConfig(None),
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{FuncNameFormat},
			FormatConfig(FullPathFormat),
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{FuncNameFormat, LineNumberFormat},
			FormatConfig(FullPathFormat),
		},
		{
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat},
			FormatConfig(FullPathFormat),
		},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat},
			[]FormatFlag{LineNumberFormat, FullPathFormat, LineNumberFormat},
			FormatConfig(FuncNameFormat),
		},
		{
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FullPathFormat},
			FormatConfig(None),
		},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FullPathFormat},
			[]FormatFlag{},
			FormatConfig(FuncNameFormat + LineNumberFormat + FullPathFormat),
		},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FullPathFormat},
			[]FormatFlag{FuncNameFormat},
			FormatConfig(LineNumberFormat + FullPathFormat),
		},
	}

	for i, s := range tests {
		var f FormatConfig
		f.Set(s.def...)
		f.Delete(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestFormatConfigDeleteError tests Delete method with invalid flag values.
func TestFormatConfigDeleteError(t *testing.T) {
	type test struct {
		value  []FormatFlag
		result bool
	}

	tests := []test{
		{[]FormatFlag{FullPathFormat}, true},
		{[]FormatFlag{FullPathFormat, FuncNameFormat}, true},
		{[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat}, true},
		{[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat, FuncNameFormat}, true},
		{
			[]FormatFlag{
				FuncNameFormat,
				FormatFlag(maxFormatConfig + 1),
				FuncNameFormat,
			},
			false,
		},
		{[]FormatFlag{None}, false},
	}

	for i, s := range tests {
		var f FormatConfig
		_, err := f.Delete(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestFormatConfigAll tests All method.
func TestFormatConfigAll(t *testing.T) {
	type test struct {
		def    []FormatFlag
		value  []FormatFlag
		result bool
	}

	tests := []test{
		{
			[]FormatFlag{FullPathFormat},
			[]FormatFlag{FullPathFormat},
			true,
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{FuncNameFormat},
			true,
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{LineNumberFormat},
			false,
		},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat},
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			false,
		},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat, None},
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			false,
		},
	}

	for i, s := range tests {
		var f FormatConfig
		f.Set(s.def...)
		if ok, _ := f.All(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatConfigAny tests Any method.
func TestFormatConfigAny(t *testing.T) {
	type test struct {
		def    []FormatFlag
		value  []FormatFlag
		result bool
	}

	tests := []test{
		{
			[]FormatFlag{FullPathFormat},
			[]FormatFlag{FullPathFormat},
			true,
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{FuncNameFormat},
			true,
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{LineNumberFormat},
			false,
		},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat},
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			true,
		},
		{
			[]FormatFlag{FuncNameFormat, LineNumberFormat, FuncNameFormat},
			[]FormatFlag{FullPathFormat, FullPathFormat, LineNumberFormat, FullPathFormat},
			true,
		},
		{
			[]FormatFlag{FullPathFormat, FuncNameFormat},
			[]FormatFlag{LineNumberFormat, FullPathFormat},
			true,
		},
	}

	for i, s := range tests {
		var f FormatConfig
		f.Set(s.def...)
		if ok, _ := f.Any(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatConfigFilePath tests FormatConfig.FilePath method.
func TestFormatConfigFilePath(t *testing.T) {
	type test struct {
		value  FormatConfig
		result bool
	}

	tests := []test{
		{FormatConfig(FullPathFormat), true},
		{FormatConfig(FuncNameFormat), false},
		{FormatConfig(LineNumberFormat), false},
		{FormatConfig(FullPathFormat + FuncNameFormat), true},
		{FormatConfig(FuncNameFormat + LineNumberFormat), false},
		{maxFormatConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.FilePath(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatConfigFuncName tests FormatConfig.FuncName method.
func TestFormatConfigFuncName(t *testing.T) {
	type test struct {
		value  FormatConfig
		result bool
	}

	tests := []test{
		{FormatConfig(FullPathFormat), false},
		{FormatConfig(FuncNameFormat), true},
		{FormatConfig(LineNumberFormat), false},
		{FormatConfig(FullPathFormat + FuncNameFormat), true},
		{FormatConfig(FuncNameFormat + LineNumberFormat), true},
		{FormatConfig(FullPathFormat + LineNumberFormat), false},
		{maxFormatConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.FuncName(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatConfigLineNumber tests FormatConfig.LineNumber method.
func TestFormatConfigLineNumber(t *testing.T) {
	type test struct {
		value  FormatConfig
		result bool
	}

	tests := []test{
		{FormatConfig(FullPathFormat), false},
		{FormatConfig(FuncNameFormat), false},
		{FormatConfig(LineNumberFormat), true},
		{FormatConfig(FullPathFormat + FuncNameFormat), false},
		{FormatConfig(FuncNameFormat + LineNumberFormat), true},
		{FormatConfig(LineNumberFormat + FullPathFormat), true},
		{maxFormatConfig + 1, false},
		{None, false},
		{0, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.LineNumber(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}
*/
