package log

import (
	"fmt"
	"strings"
)

// The getPrefix creates a log-message prefix without timestamp.
func getPrefix(level LevelFlag, config *Config, ss *StackSlice) string {
	var label, path, name, line string

	// Get level name.
	label = fmt.Sprintf("[%s]%s", LevelNames[level], config.SpaceBetweenCells)

	// Configure prefix format.
	if ok, err := config.Formats.FilePath(); ok && err == nil {
		path = fmt.Sprintf("%s%s", ss.FilePath, config.SpaceBetweenCells)
	}

	if ok, err := config.Formats.FuncName(); ok && err == nil {
		name = ss.FuncName
		if ok, err := config.Formats.LineNumber(); ok && err == nil {
			name += ":"
		} else {
			name += config.SpaceBetweenCells
		}
	}

	if ok, err := config.Formats.LineNumber(); ok && err == nil {
		line = fmt.Sprintf("%d%s", ss.FileLine, config.SpaceBetweenCells)
	}

	// Generate prefix.
	r := config.SpaceBetweenCells + label + path + name + line
	return strings.TrimSuffix(r, config.SpaceBetweenCells) + " "
}
