package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	// None means nothing.
	None = 0

	// The skip default stack offset values.
	skip = 4
)

// New returns new Log object. Accepts zero or more log-level flags
// as arguments. If logging levels are not specified, all possible
// log-levels will be activated.
func New(flags ...LevelFlag) (log *Log, err error) {
	log = &Log{
		skip:   skip,
		Writer: os.Stdout,
		Config: &Config{
			Formats:         FormatConfig(None),
			Levels:          LevelConfig(None),
			FatalStatusCode: 1,
			Prefix: &PrefixConfig{
				TimestampFormat:   TimestampFormat,
				SpaceBetweenCells: SpaceBetweenCells,
				LevelFormat:       LevelFormatConfig{},
			},
		},
	}

	if len(flags) > 0 {
		_, err = log.Config.Levels.Set(flags...)
	} else {
		_, err = log.Config.Levels.Set(Panic, Fatal, Error,
			Warn, Info, Debug, Trace)
	}

	return log, err
}

// Log is the logger object.
type Log struct {
	// The skip is default stack offset.
	skip int

	// Writer is the message receiver object (os.Stdout by default).
	Writer io.Writer

	Config *Config
}

// The echo method creates a message of the fmt.Fprint format.
// It is used as a base for all levels of logging in the fmt.Fprint format.
func (l *Log) echo(skip int, level LevelFlag, w io.Writer,
	a ...interface{},
) (n int, err error) {
	ss := getStackSlice(skip)

	// If the level is not supported.
	if v, err := l.Config.Levels.Has(level); !v || err != nil {
		return 0, err
	}

	// Generate prefix.
	timestamp := time.Now().Format(l.Config.Prefix.TimestampFormat)
	prefix := timestamp + getPrefix(level, l.Config, ss)
	a = append([]interface{}{prefix}, a...)

	return fmt.Fprint(w, a...)
}

// The echof method creates a message of the fmt.Fprintf format.
// It is used as a base for all levels of logging in the fmt.Fprintf format.
func (l *Log) echof(skip int, level LevelFlag, w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	ss := getStackSlice(skip)

	// If the level is not supported.
	if v, err := l.Config.Levels.Has(level); !v || err != nil {
		return 0, err
	}

	// Generate log prefix.
	timestamp := time.Now().Format(l.Config.Prefix.TimestampFormat)
	prefix := timestamp + getPrefix(level, l.Config, ss) + format

	return fmt.Fprintf(w, prefix, a...)
}

// The echoln method creates a message of the fmt.Fprintln format.
// It is used as a base for all levels of logging in the fmt.Fprintln format.
func (l *Log) echoln(skip int, level LevelFlag, w io.Writer,
	a ...interface{},
) (n int, err error) {
	ss := getStackSlice(skip)

	// If the level is not supported.
	if v, err := l.Config.Levels.Has(level); !v || err != nil {
		return 0, err
	}

	// Generate log prefix.
	timestamp := time.Now().Format(l.Config.Prefix.TimestampFormat)
	prefix := timestamp + getPrefix(level, l.Config, ss)

	return fmt.Fprint(w, prefix+fmt.Sprintln(a...))
}

// Copy returns copy of the log object.
func (l *Log) Copy() *Log {
	return &Log{
		Writer: l.Writer,
		skip:   l.skip,
		Config: &Config{
			Levels:          l.Config.Levels,
			Formats:         l.Config.Formats,
			FatalStatusCode: l.Config.FatalStatusCode,
			Prefix: &PrefixConfig{
				TimestampFormat:   l.Config.Prefix.TimestampFormat,
				SpaceBetweenCells: l.Config.Prefix.SpaceBetweenCells,
				LevelFormat:       l.Config.Prefix.LevelFormat,
			},
		},
	}
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fpanic(w io.Writer, a ...interface{}) (n int, err error) {
	n, err = l.echo(l.skip, Panic, w, a...)
	if ok, _ := l.Config.Levels.Has(Panic); ok {
		panic(fmt.Sprint(a...))
	}

	return
}

// Fpanicf creates message with Panic level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fpanicf(w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	n, err = l.echof(l.skip, Panic, w, format, a...)
	if ok, _ := l.Config.Levels.Has(Panic); ok {
		panic(fmt.Sprintf(format, a...))
	}

	return
}

// Fpanicln creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Fpanicln(w io.Writer, a ...interface{}) (n int, err error) {
	n, err = l.echoln(l.skip, Panic, w, a...)
	if ok, _ := l.Config.Levels.Has(Panic); ok {
		panic(fmt.Sprintln(a...))
	}

	return
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Panic(a ...interface{}) (n int, err error) {
	n, err = l.echo(l.skip, Panic, l.Writer, a...)
	if ok, _ := l.Config.Levels.Has(Panic); ok {
		panic(fmt.Sprint(a...))
	}

	return
}

// Panicf creates message with Panic level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Panicf(format string, a ...interface{}) (n int, err error) {
	n, err = l.echof(l.skip, Panic, l.Writer, format, a...)
	if ok, _ := l.Config.Levels.Has(Panic); ok {
		panic(fmt.Sprintf(format, a...))
	}

	return
}

// Panicln creates message with Panic, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Panicln(a ...interface{}) (n int, err error) {
	n, err = l.echoln(l.skip, Panic, l.Writer, a...)
	if ok, _ := l.Config.Levels.Has(Panic); ok {
		panic(fmt.Sprintln(a...))
	}

	return
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ffatal(w io.Writer, a ...interface{}) (n int, err error) {
	n, err = l.echo(l.skip, Fatal, w, a...)
	if ok, _ := l.Config.Levels.Has(Fatal); ok && l.Config.FatalAllowed() {
		os.Exit(l.Config.FatalStatusCode)
	}

	return
}

// Ffatalf creates message with Fatal level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ffatalf(w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	n, err = l.echof(l.skip, Fatal, w, format, a...)
	if ok, _ := l.Config.Levels.Has(Fatal); ok && l.Config.FatalAllowed() {
		os.Exit(l.Config.FatalStatusCode)
	}

	return
}

// Ffatalln creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Ffatalln(w io.Writer, a ...interface{}) (n int, err error) {
	n, err = l.echoln(l.skip, Fatal, w, a...)
	if ok, _ := l.Config.Levels.Has(Fatal); ok && l.Config.FatalAllowed() {
		os.Exit(l.Config.FatalStatusCode)
	}

	return
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Fatal(a ...interface{}) (n int, err error) {
	n, err = l.echo(l.skip, Fatal, l.Writer, a...)
	if ok, _ := l.Config.Levels.Has(Fatal); ok && l.Config.FatalAllowed() {
		os.Exit(l.Config.FatalStatusCode)
	}

	return
}

// Fatalf creates message with Fatal level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Fatalf(format string, a ...interface{}) (n int, err error) {
	n, err = l.echof(l.skip, Fatal, l.Writer, format, a...)
	if ok, _ := l.Config.Levels.Has(Fatal); ok && l.Config.FatalAllowed() {
		os.Exit(l.Config.FatalStatusCode)
	}

	return
}

// Fatalln creates message with Fatal, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Fatalln(a ...interface{}) (n int, err error) {
	n, err = l.echoln(l.skip, Fatal, l.Writer, a...)
	if ok, _ := l.Config.Levels.Has(Fatal); ok && l.Config.FatalAllowed() {
		os.Exit(l.Config.FatalStatusCode)
	}

	return
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ferror(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Error, w, a...)
}

// Ferrorf creates message with Error level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ferrorf(w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	return l.echof(l.skip, Error, w, format, a...)
}

// Ferrorln creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Ferrorln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Error, w, a...)
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Error(a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Error, l.Writer, a...)
}

// Errorf creates message with Error level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Errorf(format string, a ...interface{}) (n int, err error) {
	return l.echof(l.skip, Error, l.Writer, format, a...)
}

// Errorln creates message with Error, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Errorln(a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Error, l.Writer, a...)
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fwarn(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Warn, w, a...)
}

// Fwarnf creates message with Warn level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fwarnf(w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	return l.echof(l.skip, Warn, w, format, a...)
}

// Fwarnln creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Fwarnln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Warn, w, a...)
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Warn(a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Warn, l.Writer, a...)
}

// Warnf creates message with Warn level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Warnf(format string, a ...interface{}) (n int, err error) {
	return l.echof(l.skip, Warn, l.Writer, format, a...)
}

// Warnln creates message with Warn, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Warnln(a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Warn, l.Writer, a...)
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Finfo(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Info, w, a...)
}

// Finfof creates message with Info level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Finfof(w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	return l.echof(l.skip, Info, w, format, a...)
}

// Finfoln creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Finfoln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Info, w, a...)
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Info(a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Info, l.Writer, a...)
}

// Infof creates message with Info level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Infof(format string, a ...interface{}) (n int, err error) {
	return l.echof(l.skip, Info, l.Writer, format, a...)
}

// Infoln creates message with Info, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Infoln(a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Info, l.Writer, a...)
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fdebug(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Debug, w, a...)
}

// Fdebugf creates message with Debug level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Fdebugf(w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	return l.echof(l.skip, Debug, w, format, a...)
}

// Fdebugln creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Fdebugln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Debug, w, a...)
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Debug(a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Debug, l.Writer, a...)
}

// Debugf creates message with Debug level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Debugf(format string, a ...interface{}) (n int, err error) {
	return l.echof(l.skip, Debug, l.Writer, format, a...)
}

// Debugln creates message with Debug, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Debugln(a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Debug, l.Writer, a...)
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ftrace(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Trace, w, a...)
}

// Ftracef creates message with Trace level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Log) Ftracef(w io.Writer, format string,
	a ...interface{},
) (n int, err error) {
	return l.echof(l.skip, Trace, w, format, a...)
}

// Ftraceln creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Ftraceln(w io.Writer, a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Trace, w, a...)
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Log) Trace(a ...interface{}) (n int, err error) {
	return l.echo(l.skip, Trace, l.Writer, a...)
}

// Tracef creates message with Trace level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Log) Tracef(format string, a ...interface{}) (n int, err error) {
	return l.echof(l.skip, Trace, l.Writer, format, a...)
}

// Traceln creates message with Trace, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Log) Traceln(a ...interface{}) (n int, err error) {
	return l.echoln(l.skip, Trace, l.Writer, a...)
}
