package log

import "testing"

// TestLevelsSet tests Set method.
func TestLevelsSet(t *testing.T) {
	var (
		levels = make(Levels, 3)
		result = map[Level]bool{
			DEBUG: true,
			ERROR: false,
			FATAL: false,
			INFO:  true,
			TRACE: false,
			WARN:  true,
		}
	)

	l := levels.Set(DEBUG, WARN, INFO)

	count := 0
	for level, value := range result {
		if v := levels[level]; v != value {
			t.Errorf("test for %s level is failed, "+
				"expected %t but %t", level, value, v)
		}

		if value {
			count++
		}
	}

	if len(l) != count {
		t.Errorf("test for result len is failed, "+
			"expected %d but %d", count, len(l))
	}
}

// TestLevelsAdd tests Add method.
func TestLevelsAdd(t *testing.T) {
	var (
		levels = make(Levels, 3)
		result = map[Level]bool{
			DEBUG: true,
			ERROR: false,
			FATAL: false,
			INFO:  true,
			TRACE: true,
			WARN:  true,
		}
	)

	levels.Set(DEBUG, WARN)
	l := levels.Add(INFO, TRACE)

	count := 0
	for level, value := range result {
		if v := levels[level]; v != value {
			t.Errorf("test for %s level is failed, "+
				"expected %t but %t", level, value, v)
		}

		if value {
			count++
		}
	}

	if len(l) != count {
		t.Errorf("test for result len is failed, "+
			"expected %d but %d", count, len(l))
	}
}

// TestLevelsDelete tests Delete method.
func TestLevelsDelete(t *testing.T) {
	var (
		levels = make(Levels, 3)
		result = map[Level]bool{
			DEBUG: true,
			ERROR: false,
			FATAL: false,
			INFO:  false,
			TRACE: false,
			WARN:  false,
		}
	)

	levels.Set(DEBUG, WARN, INFO)
	l := levels.Delete(ERROR, FATAL, INFO, TRACE, WARN)

	count := 0
	for level, value := range result {
		if v := levels[level]; v != value {
			t.Errorf("test for %s level is failed, "+
				"expected %t but %t", level, value, v)
		}

		if value {
			count++
		}
	}

	if len(l) != count {
		t.Errorf("test for result len is failed, "+
			"expected %d but %d", count, len(l))
	}
}

// TestLevelsAll tests All method.
func TestLevelsAll(t *testing.T) {
	type test struct {
		result bool
		levels []Level
	}

	var (
		levels = Levels{
			DEBUG: true,
			ERROR: false,
			FATAL: false,
			INFO:  true,
			TRACE: false,
			WARN:  true,
		}
		tests = []test{
			{true, []Level{DEBUG, INFO, WARN}},
			{true, []Level{DEBUG, INFO}},
			{true, []Level{DEBUG}},
			{false, []Level{DEBUG, ERROR}},
			{false, []Level{ERROR}},
			{false, []Level{ERROR, TRACE}},
			{false, []Level{ERROR, TRACE, INFO}},
		}
	)

	for i, s := range tests {
		if r := levels.All(s.levels...); s.result != r {
			t.Errorf("test %d is failed, "+
				"expected %t but %t", i, s.result, r)
		}
	}
}

// TestLevelsAny tests Any method.
func TestLevelsAny(t *testing.T) {
	type test struct {
		result bool
		levels []Level
	}

	var (
		levels = Levels{
			DEBUG: true,
			ERROR: false,
			FATAL: false,
			INFO:  true,
			TRACE: false,
			WARN:  true,
		}
		tests = []test{
			{true, []Level{DEBUG, INFO, WARN}},
			{true, []Level{DEBUG, INFO}},
			{true, []Level{DEBUG}},
			{true, []Level{DEBUG, ERROR}},
			{false, []Level{ERROR}},
			{false, []Level{ERROR, TRACE}},
			{true, []Level{ERROR, TRACE, INFO}},
		}
	)

	for i, s := range tests {
		if r := levels.Any(s.levels...); s.result != r {
			t.Errorf("test %d is failed, "+
				"expected %t but %t", i, s.result, r)
		}
	}
}
