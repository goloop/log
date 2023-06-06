package log

import "testing"

// TestConfigFatalAllowed tests Config.FatalAllowed method.
func TestConfigFatalAllowed(t *testing.T) {
	type test struct {
		value  int
		result bool
	}

	tests := []test{
		{0, false},
		{1, true},
		{32, true},
	}

	for i, s := range tests {
		c := &Config{FatalStatusCode: s.value}
		if ok := c.FatalAllowed(); ok != s.result {
			t.Errorf("test for %d is failed, "+
				"expected %t but %t", i, s.result, ok)
		}
	}
}
