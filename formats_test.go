package log

import "testing"

// TestFormatIsValid tests Format.IsValid method.
func TestFormatIsValid(t *testing.T) {
	type test struct {
		value  Format
		result bool
	}

	var tests = []test{
		{FilePath, true},
		{FuncName, true},
		{LineNumber, true},
		{FilePath + FilePath, true}, // FilePath + FilePath == FuncName
		{FilePath + FuncName, false},
		{Format(maxFormatsValue + 1), false},
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

// TestFormatsIsValid tests Formats.IsValid method.
func TestFormatsIsValid(t *testing.T) {
	type test struct {
		value  Formats
		result bool
	}

	var tests = []test{
		{Formats(FilePath), true},
		{Formats(FuncName), true},
		{Formats(LineNumber), true},
		{Formats(FilePath + FilePath), true},
		{Formats(FilePath + FilePath + LineNumber), true},
		{Formats(maxFormatsValue + 1), false},
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

// TestFormatsPrivateHas tests Formats.Has method.
func TestFormatsPrivateHas(t *testing.T) {
	type test struct {
		value  Formats
		target Format
		result bool
	}

	var tests = []test{
		{Formats(FilePath), FilePath, true},
		{Formats(FilePath), FuncName, false},
		{Formats(FuncName), FuncName, true},
		{Formats(FuncName), LineNumber, false},
		{Formats(LineNumber), LineNumber, true},
		{Formats(LineNumber), FilePath, false},
		{Formats(FilePath + FuncName), FilePath, true},
		{Formats(FilePath + FuncName), FuncName, true},
		{Formats(FilePath + FuncName), LineNumber, false},
		{Formats(FilePath + FuncName + LineNumber), FilePath, true},
		{Formats(FilePath + FuncName + LineNumber), FuncName, true},
		{Formats(FilePath + FuncName + LineNumber), LineNumber, true},
		{Formats(FilePath + FuncName + LineNumber), 0, false},
		{Formats(FilePath + LineNumber), FuncName, false},
		{Formats(FilePath + LineNumber), LineNumber, true},
		{Formats(FilePath + LineNumber), FilePath, true},
		{Formats(FilePath + LineNumber), None, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Has(s.target); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatsSet tests Formats.Set method.
func TestFormatsSet(t *testing.T) {
	type test struct {
		value  []Format
		result Formats
	}

	var tests = []test{
		{[]Format{FilePath}, Formats(FilePath)},
		{[]Format{FilePath, FuncName}, Formats(FilePath + FuncName)},
		{[]Format{FuncName, LineNumber}, Formats(FuncName + LineNumber)},
		{
			[]Format{FuncName, LineNumber, FuncName},
			Formats(FuncName + LineNumber),
		},
	}

	for i, s := range tests {
		var f Formats
		f.Set(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestFormatsSetError tests Formats.Set method with invalid flag values.
func TestFormatsSetError(t *testing.T) {
	type test struct {
		value  []Format
		result bool
	}

	var tests = []test{
		{[]Format{FilePath}, true},
		{[]Format{FilePath, FuncName}, true},
		{[]Format{FuncName, LineNumber}, true},
		{[]Format{FuncName, LineNumber, FuncName}, true},
		{[]Format{Format(maxFormatsValue) + 1}, false},
		{[]Format{FuncName, None}, false},
	}

	for i, s := range tests {
		var f Formats
		_, err := f.Set(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, err == nil)
		}
	}
}

// TestFormatsAdd tests Add method.
func TestFormatsAdd(t *testing.T) {
	type test struct {
		def    []Format
		value  []Format
		result Formats
	}

	var tests = []test{
		{
			[]Format{FilePath},
			[]Format{FilePath},
			Formats(FilePath),
		},
		{
			[]Format{FuncName},
			[]Format{FilePath, FuncName},
			Formats(FilePath + FuncName),
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{FuncName, LineNumber},
			Formats(FilePath + FuncName + LineNumber),
		},
		{
			[]Format{LineNumber, FilePath},
			[]Format{FuncName, LineNumber, FuncName},
			Formats(FilePath + FuncName + LineNumber),
		},
		{
			[]Format{LineNumber, FilePath},
			[]Format{FuncName, LineNumber, FilePath},
			Formats(FilePath + FuncName + LineNumber),
		},
	}

	for i, s := range tests {
		var f Formats
		f.Set(s.def...)
		f.Add(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestFormatsAddError tests Add method with invalid flag values.
func TestFormatsAddError(t *testing.T) {
	type test struct {
		value  []Format
		result bool
	}

	var tests = []test{
		{[]Format{FilePath}, true},
		{[]Format{FilePath, FuncName}, true},
		{[]Format{FuncName, LineNumber, FuncName}, true},
		{[]Format{FuncName, LineNumber, FuncName, FuncName}, true},
		{[]Format{FuncName, Format(maxFormatsValue) + 1, FuncName}, false},
		{[]Format{FuncName, None, FuncName}, false},
		{[]Format{None}, false},
	}

	for i, s := range tests {
		var f Formats
		_, err := f.Add(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestFormatsDelete tests Delete method.
func TestFormatsDelete(t *testing.T) {
	type test struct {
		def    []Format
		value  []Format
		result Formats
	}

	var tests = []test{
		{
			[]Format{FilePath},
			[]Format{FilePath},
			Formats(None),
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{FuncName},
			Formats(FilePath),
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{FuncName, LineNumber},
			Formats(FilePath),
		},
		{
			[]Format{LineNumber, FilePath},
			[]Format{FuncName, LineNumber, FuncName},
			Formats(FilePath),
		},
		{
			[]Format{FuncName, LineNumber},
			[]Format{LineNumber, FilePath, LineNumber},
			Formats(FuncName),
		},
		{
			[]Format{LineNumber, FilePath},
			[]Format{FuncName, LineNumber, FilePath},
			Formats(None),
		},
		{
			[]Format{FuncName, LineNumber, FilePath},
			[]Format{},
			Formats(FuncName + LineNumber + FilePath),
		},
		{
			[]Format{FuncName, LineNumber, FilePath},
			[]Format{FuncName},
			Formats(LineNumber + FilePath),
		},
	}

	for i, s := range tests {
		var f Formats
		f.Set(s.def...)
		f.Delete(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestFormatsDeleteError tests Delete method with invalid flag values.
func TestFormatsDeleteError(t *testing.T) {
	type test struct {
		value  []Format
		result bool
	}

	var tests = []test{
		{[]Format{FilePath}, true},
		{[]Format{FilePath, FuncName}, true},
		{[]Format{FuncName, LineNumber, FuncName}, true},
		{[]Format{FuncName, LineNumber, FuncName, FuncName}, true},
		{[]Format{FuncName, Format(maxFormatsValue + 1), FuncName}, false},
		{[]Format{None}, false},
	}

	for i, s := range tests {
		var f Formats
		_, err := f.Delete(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestFormatsAll tests All method.
func TestFormatsAll(t *testing.T) {
	type test struct {
		def    []Format
		value  []Format
		result bool
	}

	var tests = []test{
		{
			[]Format{FilePath},
			[]Format{FilePath},
			true,
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{FuncName},
			true,
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{LineNumber},
			false,
		},
		{
			[]Format{FuncName, LineNumber, FuncName},
			[]Format{LineNumber, FilePath},
			false,
		},
		{
			[]Format{FuncName, LineNumber, FuncName, None},
			[]Format{LineNumber, FilePath},
			false,
		},
	}

	for i, s := range tests {
		var f Formats
		f.Set(s.def...)
		if ok, _ := f.All(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatsAny tests Any method.
func TestFormatsAny(t *testing.T) {
	type test struct {
		def    []Format
		value  []Format
		result bool
	}

	var tests = []test{
		{
			[]Format{FilePath},
			[]Format{FilePath},
			true,
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{FuncName},
			true,
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{LineNumber},
			false,
		},
		{
			[]Format{FuncName, LineNumber, FuncName},
			[]Format{LineNumber, FilePath},
			true,
		},
		{
			[]Format{FuncName, LineNumber, FuncName},
			[]Format{FilePath, FilePath, LineNumber, FilePath},
			true,
		},
		{
			[]Format{FilePath, FuncName},
			[]Format{LineNumber, FilePath},
			true,
		},
	}

	for i, s := range tests {
		var f Formats
		f.Set(s.def...)
		if ok, _ := f.Any(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatsFilePath tests Formats.FilePath method.
func TestFormatsFilePath(t *testing.T) {
	type test struct {
		value  Formats
		result bool
	}

	var tests = []test{
		{Formats(FilePath), true},
		{Formats(FuncName), false},
		{Formats(LineNumber), false},
		{Formats(FilePath + FuncName), true},
		{Formats(FuncName + LineNumber), false},
		{Formats(maxFormatsValue + 1), false},
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

// TestFormatsFuncName tests Formats.FuncName method.
func TestFormatsFuncName(t *testing.T) {
	type test struct {
		value  Formats
		result bool
	}

	var tests = []test{
		{Formats(FilePath), false},
		{Formats(FuncName), true},
		{Formats(LineNumber), false},
		{Formats(FilePath + FuncName), true},
		{Formats(FuncName + LineNumber), true},
		{Formats(FilePath + LineNumber), false},
		{Formats(maxFormatsValue + 1), false},
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

// TestFormatsLineNumber tests Formats.LineNumber method.
func TestFormatsLineNumber(t *testing.T) {
	type test struct {
		value  Formats
		result bool
	}

	var tests = []test{
		{Formats(FilePath), false},
		{Formats(FuncName), false},
		{Formats(LineNumber), true},
		{Formats(FilePath + FuncName), false},
		{Formats(FuncName + LineNumber), true},
		{Formats(LineNumber + FilePath), true},
		{Formats(maxFormatsValue + 1), false},
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
