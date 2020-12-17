package log

import "testing"

// TestFormatFlagIsValid tests FormatFlag.IsValid method.
func TestFormatFlagIsValid(t *testing.T) {
	type test struct {
		value  FormatFlag
		result bool
	}

	var tests = []test{
		{FilePath, true},
		{FuncName, true},
		{LineNumber, true},
		{FilePath + FilePath, true}, // FilePath + FilePath == FuncName
		{FilePath + FuncName, false},
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

	var tests = []test{
		{FormatConfig(FilePath), true},
		{FormatConfig(FuncName), true},
		{FormatConfig(LineNumber), true},
		{FormatConfig(FilePath + FilePath), true},
		{FormatConfig(FilePath + FilePath + LineNumber), true},
		{FormatConfig(maxFormatConfig + 1), false},
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

	var tests = []test{
		{FormatConfig(FilePath), FilePath, true},
		{FormatConfig(FilePath), FuncName, false},
		{FormatConfig(FuncName), FuncName, true},
		{FormatConfig(FuncName), LineNumber, false},
		{FormatConfig(LineNumber), LineNumber, true},
		{FormatConfig(LineNumber), FilePath, false},
		{FormatConfig(FilePath + FuncName), FilePath, true},
		{FormatConfig(FilePath + FuncName), FuncName, true},
		{FormatConfig(FilePath + FuncName), LineNumber, false},
		{FormatConfig(FilePath + FuncName + LineNumber), FilePath, true},
		{FormatConfig(FilePath + FuncName + LineNumber), FuncName, true},
		{FormatConfig(FilePath + FuncName + LineNumber), LineNumber, true},
		{FormatConfig(FilePath + FuncName + LineNumber), 0, false},
		{FormatConfig(FilePath + LineNumber), FuncName, false},
		{FormatConfig(FilePath + LineNumber), LineNumber, true},
		{FormatConfig(FilePath + LineNumber), FilePath, true},
		{FormatConfig(FilePath + LineNumber), None, false},
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

	var tests = []test{
		{[]FormatFlag{FilePath}, FormatConfig(FilePath)},
		{[]FormatFlag{FilePath, FuncName}, FormatConfig(FilePath + FuncName)},
		{[]FormatFlag{FuncName, LineNumber}, FormatConfig(FuncName + LineNumber)},
		{
			[]FormatFlag{FuncName, LineNumber, FuncName},
			FormatConfig(FuncName + LineNumber),
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

	var tests = []test{
		{[]FormatFlag{FilePath}, true},
		{[]FormatFlag{FilePath, FuncName}, true},
		{[]FormatFlag{FuncName, LineNumber}, true},
		{[]FormatFlag{FuncName, LineNumber, FuncName}, true},
		{[]FormatFlag{FormatFlag(maxFormatConfig) + 1}, false},
		{[]FormatFlag{FuncName, None}, false},
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

	var tests = []test{
		{
			[]FormatFlag{FilePath},
			[]FormatFlag{FilePath},
			FormatConfig(FilePath),
		},
		{
			[]FormatFlag{FuncName},
			[]FormatFlag{FilePath, FuncName},
			FormatConfig(FilePath + FuncName),
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{FuncName, LineNumber},
			FormatConfig(FilePath + FuncName + LineNumber),
		},
		{
			[]FormatFlag{LineNumber, FilePath},
			[]FormatFlag{FuncName, LineNumber, FuncName},
			FormatConfig(FilePath + FuncName + LineNumber),
		},
		{
			[]FormatFlag{LineNumber, FilePath},
			[]FormatFlag{FuncName, LineNumber, FilePath},
			FormatConfig(FilePath + FuncName + LineNumber),
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

	var tests = []test{
		{[]FormatFlag{FilePath}, true},
		{[]FormatFlag{FilePath, FuncName}, true},
		{[]FormatFlag{FuncName, LineNumber, FuncName}, true},
		{[]FormatFlag{FuncName, LineNumber, FuncName, FuncName}, true},
		{[]FormatFlag{FuncName, FormatFlag(maxFormatConfig) + 1, FuncName}, false},
		{[]FormatFlag{FuncName, None, FuncName}, false},
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

	var tests = []test{
		{
			[]FormatFlag{FilePath},
			[]FormatFlag{FilePath},
			FormatConfig(None),
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{FuncName},
			FormatConfig(FilePath),
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{FuncName, LineNumber},
			FormatConfig(FilePath),
		},
		{
			[]FormatFlag{LineNumber, FilePath},
			[]FormatFlag{FuncName, LineNumber, FuncName},
			FormatConfig(FilePath),
		},
		{
			[]FormatFlag{FuncName, LineNumber},
			[]FormatFlag{LineNumber, FilePath, LineNumber},
			FormatConfig(FuncName),
		},
		{
			[]FormatFlag{LineNumber, FilePath},
			[]FormatFlag{FuncName, LineNumber, FilePath},
			FormatConfig(None),
		},
		{
			[]FormatFlag{FuncName, LineNumber, FilePath},
			[]FormatFlag{},
			FormatConfig(FuncName + LineNumber + FilePath),
		},
		{
			[]FormatFlag{FuncName, LineNumber, FilePath},
			[]FormatFlag{FuncName},
			FormatConfig(LineNumber + FilePath),
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

	var tests = []test{
		{[]FormatFlag{FilePath}, true},
		{[]FormatFlag{FilePath, FuncName}, true},
		{[]FormatFlag{FuncName, LineNumber, FuncName}, true},
		{[]FormatFlag{FuncName, LineNumber, FuncName, FuncName}, true},
		{
			[]FormatFlag{
				FuncName,
				FormatFlag(maxFormatConfig + 1),
				FuncName,
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

	var tests = []test{
		{
			[]FormatFlag{FilePath},
			[]FormatFlag{FilePath},
			true,
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{FuncName},
			true,
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{LineNumber},
			false,
		},
		{
			[]FormatFlag{FuncName, LineNumber, FuncName},
			[]FormatFlag{LineNumber, FilePath},
			false,
		},
		{
			[]FormatFlag{FuncName, LineNumber, FuncName, None},
			[]FormatFlag{LineNumber, FilePath},
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

	var tests = []test{
		{
			[]FormatFlag{FilePath},
			[]FormatFlag{FilePath},
			true,
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{FuncName},
			true,
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{LineNumber},
			false,
		},
		{
			[]FormatFlag{FuncName, LineNumber, FuncName},
			[]FormatFlag{LineNumber, FilePath},
			true,
		},
		{
			[]FormatFlag{FuncName, LineNumber, FuncName},
			[]FormatFlag{FilePath, FilePath, LineNumber, FilePath},
			true,
		},
		{
			[]FormatFlag{FilePath, FuncName},
			[]FormatFlag{LineNumber, FilePath},
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

	var tests = []test{
		{FormatConfig(FilePath), true},
		{FormatConfig(FuncName), false},
		{FormatConfig(LineNumber), false},
		{FormatConfig(FilePath + FuncName), true},
		{FormatConfig(FuncName + LineNumber), false},
		{FormatConfig(maxFormatConfig + 1), false},
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

	var tests = []test{
		{FormatConfig(FilePath), false},
		{FormatConfig(FuncName), true},
		{FormatConfig(LineNumber), false},
		{FormatConfig(FilePath + FuncName), true},
		{FormatConfig(FuncName + LineNumber), true},
		{FormatConfig(FilePath + LineNumber), false},
		{FormatConfig(maxFormatConfig + 1), false},
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

	var tests = []test{
		{FormatConfig(FilePath), false},
		{FormatConfig(FuncName), false},
		{FormatConfig(LineNumber), true},
		{FormatConfig(FilePath + FuncName), false},
		{FormatConfig(FuncName + LineNumber), true},
		{FormatConfig(LineNumber + FilePath), true},
		{FormatConfig(maxFormatConfig + 1), false},
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
