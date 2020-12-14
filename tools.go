package log

import (
	"fmt"
	"time"
)

// The activeLevels returns only active level list.
func activeLevels(m map[Level]bool) []Level {
	var result = make([]Level, 0, len(m))
	for key, value := range m {
		if value {
			result = append(result, key)
		}
	}

	return result
}

// The getPrefix creates a log message prefix based on the collected call
// stack data, time and log level.
func getPrefix(trace *Trace, format, timestampFormat string, level Level,
	showFilePath, showFuncName, showFileLine bool) string {
	var path, name, line string
	timestamp := time.Now().Format(timestampFormat)

	if showFilePath || level == TRACE {
		path = fmt.Sprintf("%s ", trace.FilePath)
	}

	if showFuncName || level == TRACE {
		name = trace.FuncName
		if showFileLine || level == TRACE {
			name += ":"
		} else {
			name += " "
		}
	}

	if showFileLine || level == TRACE {
		line = fmt.Sprintf("%d ", trace.FileLine)
	}

	r := fmt.Sprintf("%s [%s] %s%s%s", timestamp, level, path, name, line)
	if len(format) > 0 {
		r += format
	}

	return r
}

// The in function returns true if levels contains specified log level.
func in(level Level, levels ...Level) bool {
	for _, item := range levels {
		if level == item {
			return true
		}
	}

	return false
}
