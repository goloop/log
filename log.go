package log

import (
	"fmt"
	"io"
	"os"
)

// Timestamp is the format for displaying the time stamp in the log message.
const Timestamp = "01.02.2006 15:04:05"

// Log this is the logging object.
type Log struct {
	// Writer is the message receiver object (os.Stdout by default).
	Writer io.Writer

	// Timestamp is the format for displaying the
	// time stamp in the log message.
	Timestamp string

	// Levels map of available levels.
	Levels Levels

	// ShowFilePath if true appends the full path to the go-file,
	// the logging method was called.
	ShowFilePath bool

	// ShowFuncName if true, appends the function name where the
	// logging method was called.
	ShowFuncName bool

	// ShowFileLine if true appends the line number in the go-file
	// where the logging method was called.
	ShowFileLine bool

	// FatalStatusCode the exit code when calling the Fatal method.
	// Default - 1. If the code is <= 0, the forced exit will not occur.
	FatalStatusCode int

	// Skip default stack offset.
	Skip int
}

// New returns new Log object.
// Takes zero or more log levels as arguments. If logging levels are not
// specified, all possible logging levels will be activated, otherwise
// only the specified logging levels will be activated.
func New(levels ...Level) (*Log, error) {
	var log = Log{
		Writer:          os.Stdout,
		Timestamp:       Timestamp,
		Levels:          Levels{},
		ShowFilePath:    false,
		ShowFuncName:    true,
		ShowFileLine:    true,
		FatalStatusCode: 1,
		Skip:            SKIP,
	}

	if len(levels) > 0 {
		log.Levels.Set(levels...)
	} else {
		log.Levels.Set(FATAL, ERROR, WARN, INFO, DEBUG, TRACE)
	}

	return &log, nil
}

// The echo method creates a message of the fmt.Fprint format.
// It is used as a base for all levels of logging in the fmt.Fprint format.
func (l *Log) echo(skip int, w io.Writer, level Level,
	a ...interface{}) (n int, err error) {
	var trace = getTrace(skip)

	// If the level is not supported.
	if v := l.Levels[level]; !v {
		return 0, nil
	}

	// Generate log prefix.
	prefix := getPrefix(trace, "", l.Timestamp, level,
		l.ShowFilePath, l.ShowFuncName, l.ShowFileLine)
	a = append([]interface{}{prefix}, a...)

	return fmt.Fprint(w, a...)
}

// The echof method creates a message of the fmt.Fprintf format.
// It is used as a base for all levels of logging in the fmt.Fprintf format.
func (l *Log) echof(skip int, w io.Writer, level Level, format string,
	a ...interface{}) (n int, err error) {
	var trace = getTrace(skip)

	// If the level is not supported.
	if v := l.Levels[level]; !v {
		return 0, nil
	}

	// Generate log prefix.
	prefix := getPrefix(trace, format, l.Timestamp, level,
		l.ShowFilePath, l.ShowFuncName, l.ShowFileLine)

	return fmt.Fprintf(w, prefix, a...)
}

// The echoln method creates a message of the fmt.Fprintln format.
// It is used as a base for all levels of logging in the fmt.Fprintln format.
func (l *Log) echoln(skip int, w io.Writer, level Level,
	a ...interface{}) (n int, err error) {
	var trace = getTrace(skip)

	// If the level is not supported.
	if v := l.Levels[level]; !v {
		return 0, nil
	}

	// Generate log prefix.
	prefix := getPrefix(trace, "", l.Timestamp, level,
		l.ShowFilePath, l.ShowFuncName, l.ShowFileLine)
	prefix = prefix[:len(prefix)-1] // remove trailing space
	a = append([]interface{}{prefix}, a...)

	return fmt.Fprintln(w, a...)
}

// Copy returns copy of the log object.
func (l *Log) Copy() *Log {
	var log = Log{
		Writer:       l.Writer,
		Timestamp:    l.Timestamp,
		Levels:       l.Levels,
		ShowFilePath: l.ShowFilePath,
		ShowFuncName: l.ShowFuncName,
		ShowFileLine: l.ShowFileLine,
	}

	return &log
}

// Format sets the message prefix display configuration flags for display:
// file path, function name and file line.
func (l *Log) Format(showFilePath, showFuncName, showFileLine bool) {
	l.ShowFilePath = showFilePath
	l.ShowFuncName = showFuncName
	l.ShowFileLine = showFileLine
}

// Ffatal creates message with FATAL level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. Performs forced exit from the program
// with status - 1.
func (l *Log) Ffatal(w io.Writer, a ...interface{}) {
	l.echo(SKIP, w, FATAL, a...)
	if l.Levels.All(FATAL) && l.FatalStatusCode > 0 {
		os.Exit(l.FatalStatusCode)
	}
}

// Ffatalf creates message with FATAL level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered. Performs forced exit from the
// program with status - 1.
func (l *Log) Ffatalf(w io.Writer, format string, a ...interface{}) {
	l.echof(SKIP, w, FATAL, format, a...)
	if l.Levels.All(FATAL) && l.FatalStatusCode > 0 {
		os.Exit(l.FatalStatusCode)
	}
}

// Ffatalln creates message with FATAL level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. Performs forced exit from the
// program with status - 1.
func (l *Log) Ffatalln(w io.Writer, a ...interface{}) {
	l.echoln(SKIP, w, FATAL, a...)
	if l.Levels.All(FATAL) && l.FatalStatusCode > 0 {
		os.Exit(l.FatalStatusCode)
	}
}

// Fatal creates message with FATAL level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. Performs forced exit from the
// program with status - 1.
func (l *Log) Fatal(a ...interface{}) {
	l.echo(SKIP, l.Writer, FATAL, a...)
	if l.Levels.All(FATAL) && l.FatalStatusCode > 0 {
		os.Exit(l.FatalStatusCode)
	}
}

// Fatalf creates message with FATAL level, according to a format
// specifier and writes to log.Writer. Performs forced exit from
// the program with status - 1.
func (l *Log) Fatalf(format string, a ...interface{}) {
	l.echof(SKIP, l.Writer, FATAL, format, a...)
	if l.Levels.All(FATAL) && l.FatalStatusCode > 0 {
		os.Exit(l.FatalStatusCode)
	}
}

// Fatalln creates message with FATAL, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. Performs forced exit from
// the program with status - 1.
func (l *Log) Fatalln(a ...interface{}) {
	l.echoln(SKIP, l.Writer, FATAL, a...)
	if l.Levels.All(FATAL) && l.FatalStatusCode > 0 {
		os.Exit(l.FatalStatusCode)
	}
}

// Ferror creates message with ERROR level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ferror(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(SKIP, w, ERROR, a...)
}

// Ferrorf creates message with ERROR level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ferrorf(w io.Writer, format string,
	a ...interface{}) (n int, err error) {
	return l.echof(SKIP, w, ERROR, format, a...)
}

// Ferrorln creates message with ERROR level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Ferrorln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, w, ERROR, a...)
}

// Error creates message with ERROR level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Error(a ...interface{}) (n int, err error) {
	return l.echo(SKIP, l.Writer, ERROR, a...)
}

// Errorf creates message with ERROR level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Errorf(format string, a ...interface{}) (n int, err error) {
	return l.echof(SKIP, l.Writer, ERROR, format, a...)
}

// Errorln creates message with ERROR, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Errorln(a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, l.Writer, ERROR, a...)
}

// Fwarn creates message with WARN level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fwarn(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(SKIP, w, WARN, a...)
}

// Fwarnf creates message with WARN level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fwarnf(w io.Writer, format string,
	a ...interface{}) (n int, err error) {
	return l.echof(SKIP, w, WARN, format, a...)
}

// Fwarnln creates message with WARN level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Fwarnln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, w, WARN, a...)
}

// Warn creates message with WARN level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Warn(a ...interface{}) (n int, err error) {
	return l.echo(SKIP, l.Writer, WARN, a...)
}

// Warnf creates message with WARN level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Warnf(format string, a ...interface{}) (n int, err error) {
	return l.echof(SKIP, l.Writer, WARN, format, a...)
}

// Warnln creates message with WARN, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Warnln(a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, l.Writer, WARN, a...)
}

// Finfo creates message with INFO level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Finfo(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(SKIP, w, INFO, a...)
}

// Finfof creates message with INFO level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Finfof(w io.Writer, format string,
	a ...interface{}) (n int, err error) {
	return l.echof(SKIP, w, INFO, format, a...)
}

// Finfoln creates message with INFO level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Finfoln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, w, INFO, a...)
}

// Info creates message with INFO level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Info(a ...interface{}) (n int, err error) {
	return l.echo(SKIP, l.Writer, INFO, a...)
}

// Infof creates message with INFO level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Infof(format string, a ...interface{}) (n int, err error) {
	return l.echof(SKIP, l.Writer, INFO, format, a...)
}

// Infoln creates message with INFO, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Infoln(a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, l.Writer, INFO, a...)
}

// Fdebug creates message with DEBUG level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fdebug(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(SKIP, w, DEBUG, a...)
}

// Fdebugf creates message with DEBUG level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fdebugf(w io.Writer, format string,
	a ...interface{}) (n int, err error) {
	return l.echof(SKIP, w, DEBUG, format, a...)
}

// Fdebugln creates message with DEBUG level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Fdebugln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, w, DEBUG, a...)
}

// Debug creates message with DEBUG level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Debug(a ...interface{}) (n int, err error) {
	return l.echo(SKIP, l.Writer, DEBUG, a...)
}

// Debugf creates message with DEBUG level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Debugf(format string, a ...interface{}) (n int, err error) {
	return l.echof(SKIP, l.Writer, DEBUG, format, a...)
}

// Debugln creates message with DEBUG, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Debugln(a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, l.Writer, DEBUG, a...)
}

// Ftrace creates message with TRACE level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ftrace(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(SKIP, w, TRACE, a...)
}

// Ftracef creates message with TRACE level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ftracef(w io.Writer, format string,
	a ...interface{}) (n int, err error) {
	return l.echof(SKIP, w, TRACE, format, a...)
}

// Ftraceln creates message with TRACE level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Ftraceln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, w, TRACE, a...)
}

// Trace creates message with TRACE level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Trace(a ...interface{}) (n int, err error) {
	return l.echo(SKIP, l.Writer, TRACE, a...)
}

// Tracef creates message with TRACE level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Tracef(format string, a ...interface{}) (n int, err error) {
	return l.echof(SKIP, l.Writer, TRACE, format, a...)
}

// Traceln creates message with TRACE, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Traceln(a ...interface{}) (n int, err error) {
	return l.echoln(SKIP, l.Writer, TRACE, a...)
}
