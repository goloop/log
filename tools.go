package log

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// stackSlice contains the top-level trace information
// where the logging method was called.
type stackSlice struct {
	FileLine int
	FuncName string
	FilePath string
}

// The getStackSlice returns the stack slice. The skip argument
// is the number of stack frames to skip before taking a slice.
func getStackSlice(skip int) *stackSlice {
	ss := &stackSlice{}

	// Return program counters of function invocations on
	// the calling goroutine's stack and skipping function
	// call frames inside *Log.
	pc := make([]uintptr, skip+1) // program counters
	runtime.Callers(skip, pc)

	// Get a function at an address on the stack.
	fn := runtime.FuncForPC(pc[0])

	// Get name, path and line of the file.
	ss.FuncName = fn.Name()
	ss.FilePath, ss.FileLine = fn.FileLine(pc[0])
	if r := strings.Split(ss.FuncName, "."); len(r) > 0 {
		ss.FuncName = r[len(r)-1]
	}

	return ss
}

// The getPrefix creates a log-message prefix without timestamp.
func getPrefix(level LevelFlag, t time.Time, l *Logger, ss *stackSlice) string {
	sb := strings.Builder{}

	// Logger name.
	if l.name != "" {
		sb.WriteString(l.name)
		sb.WriteString(l.config.Prefix.SpaceBetweenCells)
	}

	// Timestamp.
	sb.WriteString(t.Format(l.config.Prefix.TimestampFormat))
	sb.WriteString(l.config.Prefix.SpaceBetweenCells)

	// Level name.
	format, ok := l.config.Prefix.LevelFormat[level]
	if !ok {
		format = LevelFormat
	}

	if len(format) != 0 {
		sb.WriteString(fmt.Sprintf(
			format+"%s",
			LevelNames[level],
			l.config.Prefix.SpaceBetweenCells,
		))
	}

	// File path.
	// The FullPath takes precedence over ShortPath.
	if ok, err := l.config.Formats.FilePath(); ok && err == nil {
		if ok, err := l.config.Formats.FullPath(); ok && err == nil {
			sb.WriteString(ss.FilePath)
			/*
				sb.WriteString(fmt.Sprintf(
					"%s%s",
					ss.FilePath,
					l.config.Prefix.SpaceBetweenCells,
				))*/
		} else if ok, err := l.config.Formats.ShortPath(); ok && err == nil {
			sb.WriteString(cutFilePath(shortPathSections, ss.FilePath))
			/*
				sb.WriteString(fmt.Sprintf(
					"%s%s",
					cutFilePath(shortPathSections, ss.FilePath),
					l.config.Prefix.SpaceBetweenCells,
				))*/
		}

		if ok, err := l.config.Formats.LineNumber(); ok && err == nil {
			if ok, err := l.config.Formats.FuncName(); !ok && err == nil {
				sb.WriteString(":")
			} else {
				sb.WriteString(l.config.Prefix.SpaceBetweenCells)
			}
		} else {
			sb.WriteString(l.config.Prefix.SpaceBetweenCells)
		}
	}

	// Function name.
	if ok, err := l.config.Formats.FuncName(); ok && err == nil {
		sb.WriteString(ss.FuncName)
		if ok, err := l.config.Formats.LineNumber(); ok && err == nil {
			sb.WriteString(":")
		} else {
			sb.WriteString(l.config.Prefix.SpaceBetweenCells)
		}
	}

	// Line number.
	if ok, err := l.config.Formats.LineNumber(); ok && err == nil {
		sb.WriteString(fmt.Sprintf(
			"%d%s",
			ss.FileLine,
			l.config.Prefix.SpaceBetweenCells,
		))
	}

	return sb.String()
}

// The cutFilePath cuts the path to the file to the
// specified number of sections.
func cutFilePath(n int, path string) string {
	sections := strings.Split(path, "/")

	// If there are fewer or equal sections than n,
	// return the path unmodified.
	if len(sections) <= n+1 {
		return path
	}

	return ".../" + strings.Join(sections[len(sections)-n:], "/")
}
