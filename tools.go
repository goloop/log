package log

import (
	"fmt"
)

// The getPrefix creates a log-message prefix without timestamp.
func getPrefix(level LevelFlag, config *Config, ss *StackSlice) string {
	var label, path, name, line string

	// Get level name.
	format, ok := config.Prefix.LevelFormat[level]
	if !ok {
		format = LevelFormat
	}
	if len(format) != 0 {
		label = fmt.Sprintf(
			format+"%s",
			LevelNames[level],
			config.Prefix.SpaceBetweenCells,
		)
	}

	// Configure prefix format.
	if ok, err := config.Formats.FilePath(); ok && err == nil {
		path = fmt.Sprintf(
			"%s%s",
			ss.FilePath,
			config.Prefix.SpaceBetweenCells,
		)
	}

	if ok, err := config.Formats.FuncName(); ok && err == nil {
		name = ss.FuncName
		if ok, err := config.Formats.LineNumber(); ok && err == nil {
			name += ":"
		} else {
			name += config.Prefix.SpaceBetweenCells
		}
	}

	if ok, err := config.Formats.LineNumber(); ok && err == nil {
		line = fmt.Sprintf(
			"%d%s",
			ss.FileLine,
			config.Prefix.SpaceBetweenCells,
		)
	}

	// Generate prefix.
	return config.Prefix.SpaceBetweenCells + label + path + name + line
}
