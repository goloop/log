package log

import "testing"

// TestLevelFormatConfigDefault tests LevelFormatConfig.Default method.
func TestLevelFormatConfigDefault(t *testing.T) {
	var format = LevelFormatConfig{}

	format[Debug] = "%s"
	format.Default()
	if len(format) != 0 {
		t.Errorf("test is failed, "+
			"expected %d but %d", 0, len(format))
	}
}

// TestLevelFormatConfigSet tests LevelFormatConfig.Set method.
func TestLevelFormatConfigSet(t *testing.T) {
	var format = LevelFormatConfig{}

	format.Set("[%s]")
	for l, name := range LevelNames {
		if v, ok := format[l]; v != "[%s]" || !ok {
			t.Errorf("test for %s level is failed, "+
				"expected %s but %s", name, "[%s]", v)
		}
	}

}

// TestLevelFormatConfigColor tests LevelFormatConfig.Color method.
func TestLevelFormatConfigColor(t *testing.T) {
	var format = LevelFormatConfig{}

	format.Color()
	for l, name := range LevelNames {
		if _, ok := format[l]; !ok {
			t.Errorf("test for %s level is failed, expected value", name)
		}
	}
}

// TestLevelFormatConfigColorf tests LevelFormatConfig.Colorf method.
func TestLevelFormatConfigColorf(t *testing.T) {
	var format = LevelFormatConfig{}

	format.Colorf("[%s]")
	for l, name := range LevelNames {
		if _, ok := format[l]; !ok {
			t.Errorf("test for %s level is failed, expected value", name)
		}
	}
}
