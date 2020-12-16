package log

import "testing"

// TestFormatsPrivateHas tests has method.
func TestFormatsPrivateHas(t *testing.T) {
	type test struct {
		value  Formats
		target Formats
		result bool
	}

	var tests = []test{
		{FilePath, FilePath, true},
		{FilePath, FuncName, false},
		{FuncName, FuncName, true},
		{FuncName, LineNumber, false},
		{LineNumber, LineNumber, true},
		{LineNumber, FilePath, false},
		{FilePath + FuncName, FilePath, true},
		{FilePath + FuncName, FuncName, true},
		{FilePath + FuncName, LineNumber, false},
		{FilePath + FuncName + LineNumber, FilePath, true},
		{FilePath + FuncName + LineNumber, FuncName, true},
		{FilePath + FuncName + LineNumber, LineNumber, true},
		{FilePath + FuncName + LineNumber, maxFormatsValue, false},
		{FilePath + LineNumber, FuncName, false},
		{FilePath + LineNumber, LineNumber, true},
		{FilePath + LineNumber, FilePath, true},
		{FilePath + LineNumber, None, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.has(s.target); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatsIsValid tests IsValid method.
func TestFormatsIsValid(t *testing.T) {
	type test struct {
		value  Formats
		result bool
	}

	var tests = []test{
		{FilePath, true},
		{FuncName, true},
		{LineNumber, true},
		{FilePath + FilePath, true},
		{FilePath + FilePath + LineNumber, true},
		{maxFormatsValue + 1, false},
		{None, true},
	}

	for i, s := range tests {
		if ok := s.value.IsValid(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestFormatsSet tests Set method.
func TestFormatsSet(t *testing.T) {
	type test struct {
		value  []Formats
		result Formats
	}

	var tests = []test{
		{[]Formats{FilePath}, FilePath},
		{[]Formats{FilePath, FuncName}, FilePath + FuncName},
		{[]Formats{FuncName, LineNumber}, FuncName + LineNumber},
		{[]Formats{FuncName, LineNumber, FuncName}, FuncName + LineNumber},
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

// TestFormatsSetError tests Set method with invalid flag values.
func TestFormatsSetError(t *testing.T) {
	type test struct {
		value  []Formats
		result bool
	}

	var tests = []test{
		{[]Formats{FilePath}, true},
		{[]Formats{FilePath, FuncName}, true},
		{[]Formats{FuncName, LineNumber}, true},
		{[]Formats{FuncName, LineNumber, FuncName}, true},
		{[]Formats{maxFormatsValue + 1}, false},
		{[]Formats{FuncName, maxFormatsValue + 1}, false},
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
		def    []Formats
		value  []Formats
		result Formats
	}

	var tests = []test{
		{
			[]Formats{FilePath},
			[]Formats{FilePath},
			FilePath,
		},
		{
			[]Formats{FuncName},
			[]Formats{FilePath, FuncName},
			FilePath + FuncName,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{FuncName, LineNumber},
			FilePath + FuncName + LineNumber,
		},
		{
			[]Formats{LineNumber, FilePath},
			[]Formats{FuncName, LineNumber, FuncName},
			FilePath + FuncName + LineNumber,
		},
		{
			[]Formats{LineNumber, FilePath},
			[]Formats{FuncName, LineNumber, FuncName, None},
			FilePath + FuncName + LineNumber,
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
		value  []Formats
		result bool
	}

	var tests = []test{
		{[]Formats{FilePath}, true},
		{[]Formats{FilePath, FuncName}, true},
		{[]Formats{FuncName, LineNumber, FuncName}, true},
		{[]Formats{FuncName, LineNumber, FuncName, FuncName}, true},
		{[]Formats{FuncName, maxFormatsValue + 1, FuncName}, false},
		{[]Formats{None}, true},
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
		def    []Formats
		value  []Formats
		result Formats
	}

	var tests = []test{
		{
			[]Formats{FilePath},
			[]Formats{FilePath},
			None,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{FuncName},
			FilePath,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{FuncName, LineNumber},
			FilePath,
		},
		{
			[]Formats{LineNumber, FilePath},
			[]Formats{FuncName, LineNumber, FuncName},
			FilePath,
		},
		{
			[]Formats{FuncName, LineNumber},
			[]Formats{LineNumber, FilePath, LineNumber},
			FuncName,
		},
		{
			[]Formats{LineNumber, FilePath},
			[]Formats{FuncName, LineNumber, FilePath},
			None,
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
		value  []Formats
		result bool
	}

	var tests = []test{
		{[]Formats{FilePath}, true},
		{[]Formats{FilePath, FuncName}, true},
		{[]Formats{FuncName, LineNumber, FuncName}, true},
		{[]Formats{FuncName, LineNumber, FuncName, FuncName}, true},
		{[]Formats{FuncName, maxFormatsValue + 1, FuncName}, false},
		{[]Formats{None}, true},
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
		def    []Formats
		value  []Formats
		result bool
	}

	var tests = []test{
		{
			[]Formats{FilePath},
			[]Formats{FilePath},
			true,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{FuncName},
			true,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{LineNumber},
			false,
		},
		{
			[]Formats{FuncName, LineNumber, FuncName},
			[]Formats{LineNumber, FilePath},
			false,
		},
		{
			[]Formats{FuncName, LineNumber, FuncName, None},
			[]Formats{LineNumber, FilePath},
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
		def    []Formats
		value  []Formats
		result bool
	}

	var tests = []test{
		{
			[]Formats{FilePath},
			[]Formats{FilePath},
			true,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{FuncName},
			true,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{LineNumber},
			false,
		},
		{
			[]Formats{FuncName, LineNumber, FuncName},
			[]Formats{LineNumber, FilePath},
			true,
		},
		{
			[]Formats{FuncName, LineNumber, FuncName, None},
			[]Formats{LineNumber, FilePath},
			true,
		},
		{
			[]Formats{FilePath, FuncName},
			[]Formats{LineNumber, FilePath},
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
