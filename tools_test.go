package log

/*
import "testing"

// TestActiveLevels tests activeLevels function.
func TestActiveLevels(t *testing.T) {
	type test struct {
		levels map[Level]bool
		result []Level
	}

	var tests = []test{
		{
			map[Level]bool{
				"1": true,
				"2": false,
			},
			[]Level{"1"},
		},
		{
			map[Level]bool{
				"1": false,
				"2": true,
			},
			[]Level{"2"},
		},
		{
			map[Level]bool{
				"1": true,
				"2": false,
				"3": true,
			},
			[]Level{"1", "3"},
		},
	}

	for i, s := range tests {
		r := activeLevels(s.levels)
		if len(r) != len(s.result) {
			t.Errorf("test %d is failed, expected %v but %v", i, s.result, r)
		}
	}
}

// TestIn tests `in` function.
func TestIn(t *testing.T) {
	type test struct {
		level  Level
		levels []Level
		result bool
	}

	var tests = []test{
		{DEBUG, []Level{INFO, WARN, DEBUG, ERROR}, true},
		{WARN, []Level{WARN, INFO, DEBUG}, true},
		{ERROR, []Level{INFO, DEBUG}, false},
	}

	for i, s := range tests {
		if r := in(s.level, s.levels...); r != s.result {
			t.Errorf("test %d is failed, expected %v but %v", i, s.result, r)
		}
	}
}
*/
