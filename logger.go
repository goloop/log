package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/goloop/g"
)

// TimestampFormt is default date and time format for a timestamp.
const (
	// FatalStatusCode is default value of the an exit code when
	// calling the Fatal method.
	FatalStatusCode = 1

	// TimestampFormat default value of the time and date format
	// for the timestamp in the log message.
	TimestampFormat = "2006/01/02 15:04:05"

	// SpaceBetweenCells is default value of the string that is set
	// between elements of information blocks in the logging prefix.
	SpaceBetweenCells = " "

	// LevelFormat is default value of the format string
	// of the log level substring.
	LevelFormat = "%s"
)

// Config is the type of logging configurations: message display
// parameters, log levels, etc.
type Config struct {
	// Formats is the flag-holder where flags responsible for
	// formatting the log message prefix.
	Formats FormatConfig

	// Levels is the flag-holder where flags responsible for
	// levels of the logging: Panic, Fatal, Error, Warn, Info etc.
	Levels LevelConfig

	// FatalStatusCode is an exit code when calling the Fatal method.
	// Default - 1. If the code is <= 0, the forced exit will not occur.
	FatalStatusCode int

	// Prefix is prefix config of the log-message.
	Prefix *PrefixConfig
}

// FatalAllowed reutrns ture if the exit code for Fatal method not equal zero.
func (c Config) FatalAllowed() bool {
	return c.FatalStatusCode > 0
}

// Logger is the logger object.
type Logger struct {
	// The skipStackFrames specifies the number of stack frames to skip before
	// the program counter stack is collected. For example, if skip is equal
	// to 4, then the top four frames of the stack will be ignored.
	skipStackFrames int

	// The name is the logger name.
	name string

	// The writer is the message receiver object (os.Stdout by default).
	writer io.Writer

	// The config is the logging configuration object.
	config *Config

	// The mu is the mutex for the log object.
	mu sync.RWMutex
}

// The echo method creates a message of the fmt.Fprint format.
// It is used as a base for all levels of logging in the fmt.Fprint format.
func (l *Logger) echo(
	skip int,
	level LevelFlag,
	w io.Writer,
	a ...interface{},
) (int, error) {
	// Lock the log object for change.
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Get the stack slice.
	ss := getStackSlice(skip)

	// If the level is not supported, this is not a logger error,
	// it just means that we should ignore the message of this type.
	if v, err := l.config.Levels.Has(level); !v || err != nil {
		return 0, nil // ... so, return zero bytes and nil.
	}

	// Generate log prefix.
	prefix := getPrefix(level, time.Now(), l, ss)
	a = append([]interface{}{prefix}, a...)

	return fmt.Fprint(w, a...)
}

// The echof method creates a message of the fmt.Fprintf format.
// It is used as a base for all levels of logging in the fmt.Fprintf format.
func (l *Logger) echof(
	skip int,
	level LevelFlag,
	w io.Writer,
	format string,
	a ...interface{},
) (int, error) {
	ss := getStackSlice(skip)

	// If the level is not supported.
	if v, err := l.config.Levels.Has(level); !v || err != nil {
		return 0, err
	}

	// Generate log prefix.
	prefix := getPrefix(level, time.Now(), l, ss) + format
	return fmt.Fprintf(w, prefix, a...)
}

// The echoln method creates a message of the fmt.Fprintln format.
// It is used as a base for all levels of logging in the fmt.Fprintln format.
func (l *Logger) echoln(
	skip int,
	level LevelFlag,
	w io.Writer,
	a ...interface{},
) (int, error) {
	ss := getStackSlice(skip)

	// If the level is not supported.
	if v, err := l.config.Levels.Has(level); !v || err != nil {
		return 0, err
	}

	// Generate log prefix.
	prefix := getPrefix(level, time.Now(), l, ss)

	return fmt.Fprint(w, prefix+fmt.Sprintln(a...))
}

// Copy returns copy of the log object.
func (l *Logger) Copy() *Logger {
	return &Logger{
		writer:          l.writer,
		name:            l.name,
		skipStackFrames: l.skipStackFrames,
		config: &Config{
			Levels:          l.config.Levels,
			Formats:         l.config.Formats,
			FatalStatusCode: l.config.FatalStatusCode,
			Prefix: &PrefixConfig{
				TimestampFormat:   l.config.Prefix.TimestampFormat,
				SpaceBetweenCells: l.config.Prefix.SpaceBetweenCells,
				LevelFormat:       l.config.Prefix.LevelFormat,
			},
		},
	}
}

// SetSkip returns the skip value of the log object.
func (l *Logger) SetSkip(skips ...int) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.skipStackFrames = skipStackFrames
	if len(skips) > 0 {
		l.skipStackFrames = g.Sum(skips...)
	}

	return l.skipStackFrames
}

// SetWriter sets the writer of the log object.
func (l *Logger) SetWriter(writers ...io.Writer) io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.writer = os.Stdout
	if len(writers) > 0 {
		if mw, ok := io.MultiWriter(writers...).(*os.File); ok {
			l.writer = mw
		}
	}

	return l.writer
}

// SetName returns the name of the log object.
func (l *Logger) SetName(names ...string) string {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.name = ""
	if len(names) != 0 {
		var name string

		// Concatenate prefixes.
		if l := len(names); l == 1 {
			// If there is only one prefix, use it as is.
			// In this case, no changes are made to the prefix.
			name = names[0]
		} else if l > 1 {
			// Several words that characterize the prefix are given.
			// In this case, they must be combined as ONE-TWO-THREE in
			// upper case, removing all special characters such as spaces,
			// colons and \t, \i, \n.
			i := 0
			sb := strings.Builder{}
			for _, p := range names {
				v := g.Trim(p, ": \t\n\r")
				if v == "" {
					continue
				}

				if i != 0 {
					sb.WriteString("-")
				}
				i++

				sb.WriteString(v)
			}

			// If the prefix is not empty, add a marker at the end.
			if sb.Len() > 0 {
				sb.WriteString(":") // add marker at the end
				name = sb.String()
			}
		}

		l.name = name
	}

	return l.name
}

// Configure sets the configuration of the log object.
func (l *Logger) Configure(config *Config) {
	l.mu.Lock()
	defer l.mu.Unlock()

	fsc := g.If(config.FatalStatusCode == 0, -1, config.FatalStatusCode)
	l.config = &Config{
		Formats:         FormatConfig(config.Formats),
		Levels:          LevelConfig(config.Levels),
		FatalStatusCode: fsc,
		Prefix:          l.config.Prefix,
	}

	if config.Prefix != nil {
		l.config.Prefix = &PrefixConfig{
			TimestampFormat:   config.Prefix.TimestampFormat,
			SpaceBetweenCells: config.Prefix.SpaceBetweenCells,
			LevelFormat:       config.Prefix.LevelFormat,
		}
	}
}

// Config returns the configuration of the log object.
func (l *Logger) Config() *Config {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.config
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Fpanic(w io.Writer, a ...interface{}) (int, error) {
	n, err := l.echo(l.skipStackFrames, PanicLevel, w, a...)
	if ok, _ := l.config.Levels.Has(PanicLevel); ok {
		panic(fmt.Sprint(a...))
	}

	return n, err
}

// Fpanicf creates message with Panic level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Fpanicf(w io.Writer, format string,
	a ...interface{},
) (int, error) {
	n, err := l.echof(l.skipStackFrames, PanicLevel, w, format, a...)
	if ok, _ := l.config.Levels.Has(PanicLevel); ok {
		panic(fmt.Sprintf(format, a...))
	}

	return n, err
}

// Fpanicln creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Fpanicln(w io.Writer, a ...interface{}) (int, error) {
	n, err := l.echoln(l.skipStackFrames, PanicLevel, w, a...)
	if ok, _ := l.config.Levels.Has(PanicLevel); ok {
		panic(fmt.Sprintln(a...))
	}

	return n, err
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Panic(a ...interface{}) (int, error) {
	n, err := l.echo(l.skipStackFrames, PanicLevel, l.writer, a...)
	if ok, _ := l.config.Levels.Has(PanicLevel); ok {
		panic(fmt.Sprint(a...))
	}

	return n, err
}

// Panicf creates message with Panic level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Logger) Panicf(format string, a ...interface{}) (int, error) {
	n, err := l.echof(l.skipStackFrames, PanicLevel, l.writer, format, a...)
	if ok, _ := l.config.Levels.Has(PanicLevel); ok {
		panic(fmt.Sprintf(format, a...))
	}

	return n, err
}

// Panicln creates message with Panic, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Logger) Panicln(a ...interface{}) (int, error) {
	n, err := l.echoln(l.skipStackFrames, PanicLevel, l.writer, a...)
	if ok, _ := l.config.Levels.Has(PanicLevel); ok {
		panic(fmt.Sprintln(a...))
	}

	return n, err
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Ffatal(w io.Writer, a ...interface{}) (int, error) {
	n, err := l.echo(l.skipStackFrames, FatalLevel, w, a...)
	if ok, _ := l.config.Levels.Has(FatalLevel); ok &&
		l.config.FatalAllowed() {
		os.Exit(l.config.FatalStatusCode)
	}

	return n, err
}

// Ffatalf creates message with Fatal level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Ffatalf(w io.Writer, format string,
	a ...interface{},
) (int, error) {
	n, err := l.echof(l.skipStackFrames, FatalLevel, w, format, a...)
	if ok, _ := l.config.Levels.Has(FatalLevel); ok &&
		l.config.FatalAllowed() {
		os.Exit(l.config.FatalStatusCode)
	}

	return n, err
}

// Ffatalln creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Ffatalln(w io.Writer, a ...interface{}) (int, error) {
	n, err := l.echoln(l.skipStackFrames, FatalLevel, w, a...)
	if ok, _ := l.config.Levels.Has(FatalLevel); ok &&
		l.config.FatalAllowed() {
		os.Exit(l.config.FatalStatusCode)
	}

	return n, err
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Fatal(a ...interface{}) (int, error) {
	n, err := l.echo(l.skipStackFrames, FatalLevel, l.writer, a...)
	if ok, _ := l.config.Levels.Has(FatalLevel); ok &&
		l.config.FatalAllowed() {
		os.Exit(l.config.FatalStatusCode)
	}

	return n, err
}

// Fatalf creates message with Fatal level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Logger) Fatalf(format string, a ...interface{}) (int, error) {
	n, err := l.echof(l.skipStackFrames, FatalLevel, l.writer, format, a...)
	if ok, _ := l.config.Levels.Has(FatalLevel); ok &&
		l.config.FatalAllowed() {
		os.Exit(l.config.FatalStatusCode)
	}

	return n, err
}

// Fatalln creates message with Fatal, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Logger) Fatalln(a ...interface{}) (int, error) {
	n, err := l.echoln(l.skipStackFrames, FatalLevel, l.writer, a...)
	if ok, _ := l.config.Levels.Has(FatalLevel); ok &&
		l.config.FatalAllowed() {
		os.Exit(l.config.FatalStatusCode)
	}

	return n, err
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Ferror(w io.Writer, a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, ErrorLevel, w, a...)
}

// Ferrorf creates message with Error level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Ferrorf(w io.Writer, format string,
	a ...interface{},
) (int, error) {
	return l.echof(l.skipStackFrames, ErrorLevel, w, format, a...)
}

// Ferrorln creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Ferrorln(w io.Writer, a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, ErrorLevel, w, a...)
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Error(a ...interface{}) (int, error) {
	fmt.Println(l.skipStackFrames, ErrorLevel, l.writer, a)
	return l.echo(l.skipStackFrames, ErrorLevel, l.writer, a...)
}

// Errorf creates message with Error level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Logger) Errorf(format string, a ...interface{}) (int, error) {
	return l.echof(l.skipStackFrames, ErrorLevel, l.writer, format, a...)
}

// Errorln creates message with Error, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Logger) Errorln(a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, ErrorLevel, l.writer, a...)
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Fwarn(w io.Writer, a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, WarnLevel, w, a...)
}

// Fwarnf creates message with Warn level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Fwarnf(w io.Writer, format string,
	a ...interface{},
) (int, error) {
	return l.echof(l.skipStackFrames, WarnLevel, w, format, a...)
}

// Fwarnln creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Fwarnln(w io.Writer, a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, WarnLevel, w, a...)
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Warn(a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, WarnLevel, l.writer, a...)
}

// Warnf creates message with Warn level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Logger) Warnf(format string, a ...interface{}) (int, error) {
	return l.echof(l.skipStackFrames, WarnLevel, l.writer, format, a...)
}

// Warnln creates message with Warn, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Logger) Warnln(a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, WarnLevel, l.writer, a...)
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Finfo(w io.Writer, a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, InfoLevel, w, a...)
}

// Finfof creates message with Info level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Finfof(w io.Writer, format string,
	a ...interface{},
) (int, error) {
	return l.echof(l.skipStackFrames, InfoLevel, w, format, a...)
}

// Finfoln creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Finfoln(w io.Writer, a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, InfoLevel, w, a...)
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Info(a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, InfoLevel, l.writer, a...)
}

// Infof creates message with Info level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Logger) Infof(format string, a ...interface{}) (int, error) {
	return l.echof(l.skipStackFrames, InfoLevel, l.writer, format, a...)
}

// Infoln creates message with Info, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Logger) Infoln(a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, InfoLevel, l.writer, a...)
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Fdebug(w io.Writer, a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, DebugLevel, w, a...)
}

// Fdebugf creates message with Debug level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Fdebugf(w io.Writer, format string,
	a ...interface{},
) (int, error) {
	return l.echof(l.skipStackFrames, DebugLevel, w, format, a...)
}

// Fdebugln creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Fdebugln(w io.Writer, a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, DebugLevel, w, a...)
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Debug(a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, DebugLevel, l.writer, a...)
}

// Debugf creates message with Debug level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Logger) Debugf(format string, a ...interface{}) (int, error) {
	return l.echof(l.skipStackFrames, DebugLevel, l.writer, format, a...)
}

// Debugln creates message with Debug, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Logger) Debugln(a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, DebugLevel, l.writer, a...)
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Ftrace(w io.Writer, a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, TraceLevel, w, a...)
}

// Ftracef creates message with Trace level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func (l *Logger) Ftracef(w io.Writer, format string,
	a ...interface{},
) (int, error) {
	return l.echof(l.skipStackFrames, TraceLevel, w, format, a...)
}

// Ftraceln creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Ftraceln(w io.Writer, a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, TraceLevel, w, a...)
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func (l *Logger) Trace(a ...interface{}) (int, error) {
	return l.echo(l.skipStackFrames, TraceLevel, l.writer, a...)
}

// Tracef creates message with Trace level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (l *Logger) Tracef(format string, a ...interface{}) (int, error) {
	return l.echof(l.skipStackFrames, TraceLevel, l.writer, format, a...)
}

// Traceln creates message with Trace, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func (l *Logger) Traceln(a ...interface{}) (int, error) {
	return l.echoln(l.skipStackFrames, TraceLevel, l.writer, a...)
}
