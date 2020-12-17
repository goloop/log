package log

import "testing"

// TestLevelIsValid tests Level.IsValid method.
func TestLevelIsValid(t *testing.T) {
	type test struct {
		value  Level
		result bool
	}

	var tests = []test{
		{Panic, true},
		{Fatal, true},
		{Error, true},
		{Panic + Panic, true}, // Panic + Panic == Fatal
		{Panic + Fatal, false},
		{Level(maxLevelsValue + 1), false},
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

// TestLevelsIsValid tests Levels.IsValid method.
func TestLevelsIsValid(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Levels(Panic), true},
		{Levels(Fatal), true},
		{Levels(Error), true},
		{Levels(Panic + Panic), true},
		{Levels(Panic + Panic + Error), true},
		{Levels(maxLevelsValue + 1), false},
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

// TestLevelsPrivateHas tests Levels.Has method.
func TestLevelsPrivateHas(t *testing.T) {
	type test struct {
		value  Levels
		target Level
		result bool
	}

	var tests = []test{
		{Levels(Panic), Panic, true},
		{Levels(Panic), Fatal, false},
		{Levels(Fatal), Fatal, true},
		{Levels(Fatal), Error, false},
		{Levels(Error), Error, true},
		{Levels(Error), Panic, false},
		{Levels(Panic + Fatal), Panic, true},
		{Levels(Panic + Fatal), Fatal, true},
		{Levels(Panic + Fatal), Error, false},
		{Levels(Panic + Fatal + Error), Panic, true},
		{Levels(Panic + Fatal + Error), Fatal, true},
		{Levels(Panic + Fatal + Error), Error, true},
		{Levels(Panic + Fatal + Error), 0, false},
		{Levels(Panic + Error), Fatal, false},
		{Levels(Panic + Error), Error, true},
		{Levels(Panic + Error), Panic, true},
		{Levels(Panic + Error), None, false},
	}

	for i, s := range tests {
		if ok, _ := s.value.Has(s.target); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}

// TestLevelsSet tests Levels.Set method.
func TestLevelsSet(t *testing.T) {
	type test struct {
		value  []Level
		result Levels
	}

	var tests = []test{
		{[]Level{Panic}, Levels(Panic)},
		{[]Level{Panic, Fatal}, Levels(Panic + Fatal)},
		{[]Level{Fatal, Error}, Levels(Fatal + Error)},
		{
			[]Level{Fatal, Error, Fatal},
			Levels(Fatal + Error),
		},
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

// TestLevelsSetError tests Levels.Set method with invalid flag values.
func TestLevelsSetError(t *testing.T) {
	type test struct {
		value  []Level
		result bool
	}

	var tests = []test{
		{[]Level{Panic}, true},
		{[]Level{Panic, Fatal}, true},
		{[]Level{Fatal, Error}, true},
		{[]Level{Fatal, Error, Fatal}, true},
		{[]Level{Level(maxLevelsValue) + 1}, false},
		{[]Level{Fatal, None}, false},
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
		def    []Level
		value  []Level
		result Levels
	}

	var tests = []test{
		{
			[]Level{Panic},
			[]Level{Panic},
			Levels(Panic),
		},
		{
			[]Level{Fatal},
			[]Level{Panic, Fatal},
			Levels(Panic + Fatal),
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Fatal, Error},
			Levels(Panic + Fatal + Error),
		},
		{
			[]Level{Error, Panic},
			[]Level{Fatal, Error, Fatal},
			Levels(Panic + Fatal + Error),
		},
		{
			[]Level{Error, Panic},
			[]Level{Fatal, Error, Panic},
			Levels(Panic + Fatal + Error),
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
		value  []Level
		result bool
	}

	var tests = []test{
		{[]Level{Panic}, true},
		{[]Level{Panic, Fatal}, true},
		{[]Level{Fatal, Error, Fatal}, true},
		{[]Level{Fatal, Error, Fatal, Fatal}, true},
		{[]Level{Fatal, Level(maxLevelsValue) + 1, Fatal}, false},
		{[]Level{Fatal, None, Fatal}, false},
		{[]Level{None}, false},
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
		def    []Level
		value  []Level
		result Levels
	}

	var tests = []test{
		{
			[]Level{Panic},
			[]Level{Panic},
			Levels(None),
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Fatal},
			Levels(Panic),
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Fatal, Error},
			Levels(Panic),
		},
		{
			[]Level{Error, Panic},
			[]Level{Fatal, Error, Fatal},
			Levels(Panic),
		},
		{
			[]Level{Fatal, Error},
			[]Level{Error, Panic, Error},
			Levels(Fatal),
		},
		{
			[]Level{Error, Panic},
			[]Level{Fatal, Error, Panic},
			Levels(None),
		},
		{
			[]Level{Fatal, Error, Panic},
			[]Level{},
			Levels(Fatal + Error + Panic),
		},
		{
			[]Level{Fatal, Error, Panic},
			[]Level{Fatal},
			Levels(Error + Panic),
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
		value  []Level
		result bool
	}

	var tests = []test{
		{[]Level{Panic}, true},
		{[]Level{Panic, Fatal}, true},
		{[]Level{Fatal, Error, Fatal}, true},
		{[]Level{Fatal, Error, Fatal, Fatal}, true},
		{[]Level{Fatal, Level(maxLevelsValue + 1), Fatal}, false},
		{[]Level{None}, false},
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
		def    []Level
		value  []Level
		result bool
	}

	var tests = []test{
		{
			[]Level{Panic},
			[]Level{Panic},
			true,
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Fatal},
			true,
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Error},
			false,
		},
		{
			[]Level{Fatal, Error, Fatal},
			[]Level{Error, Panic},
			false,
		},
		{
			[]Level{Fatal, Error, Fatal, None},
			[]Level{Error, Panic},
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
		def    []Level
		value  []Level
		result bool
	}

	var tests = []test{
		{
			[]Level{Panic},
			[]Level{Panic},
			true,
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Fatal},
			true,
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Error},
			false,
		},
		{
			[]Level{Fatal, Error, Fatal},
			[]Level{Error, Panic},
			true,
		},
		{
			[]Level{Fatal, Error, Fatal},
			[]Level{Panic, Panic, Error, Panic},
			true,
		},
		{
			[]Level{Panic, Fatal},
			[]Level{Error, Panic},
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

// TestLevelsPanic tests Levels.Panic method.
func TestLevelsPanic(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Levels(Panic), true},
		{Levels(Fatal), false},
		{Levels(Error), false},
		{Levels(Panic + Fatal), true},
		{Levels(Fatal + Error), false},
		{Levels(maxLevelsValue + 1), false},
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

// TestLevelsFatal tests Levels.Fatal method.
func TestLevelsFatal(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Levels(Panic), false},
		{Levels(Fatal), true},
		{Levels(Error), false},
		{Levels(Panic + Fatal), true},
		{Levels(Fatal + Error), true},
		{Levels(Panic + Error), false},
		{Levels(maxLevelsValue + 1), false},
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

// TestLevelsError tests Levels.Error method.
func TestLevelsError(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Levels(Panic), false},
		{Levels(Fatal), false},
		{Levels(Error), true},
		{Levels(Panic + Fatal), false},
		{Levels(Fatal + Error), true},
		{Levels(Error + Panic), true},
		{Levels(maxLevelsValue + 1), false},
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

// TestLevelsInfo tests Levels.Info method.
func TestLevelsInfo(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Levels(Info), true},
		{Levels(Fatal), false},
		{Levels(Debug), false},
		{Levels(Panic + Fatal + Debug), false},
		{Levels(Fatal + Error + Info), true},
		{Levels(Error + Info + Panic), true},
		{Levels(maxLevelsValue + 1), false},
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

// TestLevelsDebug tests Levels.Debug method.
func TestLevelsDebug(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Levels(Info), false},
		{Levels(Fatal), false},
		{Levels(Debug), true},
		{Levels(Panic + Fatal + Debug), true},
		{Levels(Fatal + Error + Info), false},
		{Levels(Error + Info + Panic + Debug), true},
		{Levels(maxLevelsValue + 1), false},
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

// TestLevelsTrace tests Levels.Trace method.
func TestLevelsTrace(t *testing.T) {
	type test struct {
		value  Levels
		result bool
	}

	var tests = []test{
		{Levels(Info), false},
		{Levels(Fatal), false},
		{Levels(Trace), true},
		{Levels(Trace + Fatal + Debug), true},
		{Levels(Fatal + Error + Info), false},
		{Levels(Error + Trace + Panic + Debug), true},
		{Levels(maxLevelsValue + 1), false},
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
