package log

import (
	"fmt"
)

// The getPrefix creates a log-message prefix without timestamp.
func getPrefix(level LevelFlag, formats FormatConfig, ss *StackSlice) string {
	var label, path, name, line string

	// Get level name.
	label = fmt.Sprintf("[%s] ", LevelNames[level])

	// Configure prefix format.
	if ok, err := formats.FilePath(); ok && err == nil {
		path = fmt.Sprintf("%s ", ss.FilePath)
	}

	if ok, err := formats.FuncName(); ok && err == nil {
		name = ss.FuncName
		if ok, err := formats.LineNumber(); ok && err == nil {
			name += ":"
		} else {
			name += " "
		}
	}

	if ok, err := formats.LineNumber(); ok && err == nil {
		line = fmt.Sprintf("%d ", ss.FileLine)
	}

	// Generate prefix.
	return label + path + name + line
}
