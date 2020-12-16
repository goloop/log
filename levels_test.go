package log

import "testing"

// TestLevelsPrivateHas tests has method.
func TestLevelsPrivateHas(t *testing.T) {
	type test struct {
		value  Levels
		target Levels
		result bool
	}

	var tests = []test{
		{Panic, Panic, true},
		{Panic, Info, false},
		{Info, Info, true},
		{Info, Debug, false},
		{Debug, Debug, true},
		{Debug, Panic, false},
		{Panic + Info, Panic, true},
		{Panic + Info, Info, true},
		{Panic + Info, Debug, false},
		{Panic + Info + Debug + Trace, Panic, true},
		{Panic + Error + Info + Debug, Info, true},
		{Panic + Info + Debug, Debug, true},
		{Panic + Info + Debug, Trace, false},
		{Panic + Debug, Info, false},
		{Panic + Debug, Debug, true},
		{Panic + Debug, Panic, true},
		{Panic + Debug, None, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.has(s.target); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelsIsValid tests IsValid method.
func TestLevelsIsValid(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Panic, true},
		{Info, true},
		{Debug, true},
		{Panic + Panic, true},
		{Panic + Panic + Debug, true},
		{maxLevelsValue + 1, false},
		{None, true},
	}

	for i, s := range tests {
		if ok := s.value.IsValid(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelsSet tests Set method.
func TestLevelsSet(t *testing.T) {
	type test struct {
		value  []Levels
		result Levels
	}

	var tests = []test{
		{[]Levels{Panic}, Panic},
		{[]Levels{Panic, Info}, Panic + Info},
		{[]Levels{Info, Debug}, Info + Debug},
		{[]Levels{Info, Debug, Info}, Info + Debug},
	}

	for i, s := range tests {
		var f Levels
		f.Set(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestLevelsSetError tests Set method with invalid flag values.
func TestLevelsSetError(t *testing.T) {
	type test struct {
		value  []Levels
		result bool
	}

	var tests = []test{
		{[]Levels{Panic}, true},
		{[]Levels{Panic, Info}, true},
		{[]Levels{Info, Debug}, true},
		{[]Levels{Info, Debug, Info}, true},
		{[]Levels{maxLevelsValue + 1}, false},
		{[]Levels{Info, maxLevelsValue + 1}, false},
	}

	for i, s := range tests {
		var f Levels
		_, err := f.Set(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, err == nil)
		}
	}
}

// TestLevelsAdd tests Add method.
func TestLevelsAdd(t *testing.T) {
	type test struct {
		def    []Levels
		value  []Levels
		result Levels
	}

	var tests = []test{
		{
			[]Levels{Panic},
			[]Levels{Panic},
			Panic,
		},
		{
			[]Levels{Info},
			[]Levels{Panic, Info},
			Panic + Info,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Info, Debug},
			Panic + Info + Debug,
		},
		{
			[]Levels{Debug, Panic},
			[]Levels{Info, Debug, Info},
			Panic + Info + Debug,
		},
		{
			[]Levels{Debug, Panic},
			[]Levels{Info, Debug, Info, None},
			Panic + Info + Debug,
		},
	}

	for i, s := range tests {
		var f Levels
		f.Set(s.def...)
		f.Add(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestLevelsAddError tests Add method with invalid flag values.
func TestLevelsAddError(t *testing.T) {
	type test struct {
		value  []Levels
		result bool
	}

	var tests = []test{
		{[]Levels{Panic}, true},
		{[]Levels{Panic, Info}, true},
		{[]Levels{Info, Debug, Info}, true},
		{[]Levels{Info, Debug, Info, Info}, true},
		{[]Levels{Info, maxLevelsValue + 1, Info}, false},
		{[]Levels{None}, true},
	}

	for i, s := range tests {
		var f Levels
		_, err := f.Add(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestLevelsDelete tests Delete method.
func TestLevelsDelete(t *testing.T) {
	type test struct {
		def    []Levels
		value  []Levels
		result Levels
	}

	var tests = []test{
		{
			[]Levels{Panic},
			[]Levels{Panic},
			None,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Info},
			Panic,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Info, Debug},
			Panic,
		},
		{
			[]Levels{Debug, Panic},
			[]Levels{Info, Debug, Info},
			Panic,
		},
		{
			[]Levels{Info, Debug},
			[]Levels{Debug, Panic, Debug},
			Info,
		},
		{
			[]Levels{Debug, Panic},
			[]Levels{Info, Debug, Panic},
			None,
		},
	}

	for i, s := range tests {
		var f Levels
		f.Set(s.def...)
		f.Delete(s.value...)
		if f != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %d but %d", i, s.result, f)
		}
	}
}

// TestLevelsDeleteError tests Delete method with invalid flag values.
func TestLevelsDeleteError(t *testing.T) {
	type test struct {
		value  []Levels
		result bool
	}

	var tests = []test{
		{[]Levels{Panic}, true},
		{[]Levels{Panic, Info}, true},
		{[]Levels{Info, Debug, Info}, true},
		{[]Levels{Info, Debug, Info, Info}, true},
		{[]Levels{Info, maxLevelsValue + 1, Info}, false},
		{[]Levels{None}, true},
	}

	for i, s := range tests {
		var f Levels
		_, err := f.Delete(s.value...)
		if (err == nil) != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, (err == nil))
		}
	}
}

// TestLevelsAll tests All method.
func TestLevelsAll(t *testing.T) {
	type test struct {
		def    []Levels
		value  []Levels
		result bool
	}

	var tests = []test{
		{
			[]Levels{Panic},
			[]Levels{Panic},
			true,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Info},
			true,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Debug},
			false,
		},
		{
			[]Levels{Info, Debug, Info},
			[]Levels{Debug, Panic},
			false,
		},
		{
			[]Levels{Info, Debug, Info, None},
			[]Levels{Debug, Panic},
			false,
		},
	}

	for i, s := range tests {
		var f Levels
		f.Set(s.def...)
		if ok, _ := f.All(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelsAny tests Any method.
func TestLevelsAny(t *testing.T) {
	type test struct {
		def    []Levels
		value  []Levels
		result bool
	}

	var tests = []test{
		{
			[]Levels{Panic},
			[]Levels{Panic},
			true,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Info},
			true,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Debug},
			false,
		},
		{
			[]Levels{Info, Debug, Info},
			[]Levels{Debug, Panic},
			true,
		},
		{
			[]Levels{Info, Debug, Info, None},
			[]Levels{Debug, Panic},
			true,
		},
		{
			[]Levels{Panic, Info},
			[]Levels{Debug, Panic},
			true,
		},
	}

	for i, s := range tests {
		var f Levels
		f.Set(s.def...)
		if ok, _ := f.Any(s.value...); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}
