package log

import (
	"io"
	"strings"

	"github.com/goloop/g"
)

// The self is the default logger instance.
var self *Logger

// New returns a new Logger object. We can optionally provide one or more
// prefixes that will be prepended to each log message. If multiple prefixes
// are provided, they will be joined with hyphens. Additionally, prefixes
// are stripped of leading and trailing whitespace characters.
//
// Example usage:
//
//	// The simplest use-case is to create a new logger without any prefixes.
//	logger := log.New()
//	logger.Info("Hello, World!")
//
//	// We can also add a prefix for our logger. Here, 'MYAPP' will be
//	// prepended to each log message.
//	loggerWithPrefix := log.New("MYAPP")
//	loggerWithPrefix.Info("Hello, World!")
//
//	// If multiple prefixes are provided, they will be joined with hyphens.
//	// Here, 'MYAPP-PREFIX' will be prepended to each log message.
//	loggerWithMultiplePrefixes := log.New("MYAPP", "PREFIX")
//	loggerWithMultiplePrefixes.Info("Hello, World!")
//
//	// Any leading and trailing whitespace characters are removed from
//	// prefixes. 'MYAPP-PREFIX' will be prepended to each log message.
//	loggerWithWhitespace := log.New(" MYAPP ", " PREFIX ")
//	loggerWithWhitespace.Info("Hello, World!")
func New(prefixes ...string) *Logger {
	// Generate prefix.
	prefix := ""
	if len(prefixes) != 0 {
		// Concatenate prefixes.
		if l := len(prefixes); l == 1 {
			// If there is only one prefix, use it as is.
			// In this case, no changes are made to the prefix.
			prefix = prefixes[0]
		} else if l > 1 {
			// Several words that characterize the prefix are given.
			// In this case, they must be combined as ONE-TWO-THREE in
			// upper case, removing all special characters such as spaces,
			// colons and \t, \i, \n.
			i := 0
			sb := strings.Builder{}
			for _, p := range prefixes {
				v := g.Trim(p, " \t\n\r")
				if v == "" {
					continue
				}

				// If one character is installed, it is added without
				// the separator '-' and is considered a marker of the
				// end of the prefix.
				// {"MYAPP:"} => "MYAPP:"
				// {"MYAPP", ":"} => "MYAPP:"
				// {"MY", "APP", ":"} => "MY-APP:"
				if l := len(v); i != 0 && l > 1 {
					sb.WriteString("-")
				}
				i++

				sb.WriteString(v)
			}

			// If the prefix is not empty, add a marker at the end.
			if sb.Len() > 0 {
				prefix = sb.String()
			}
		}
	}

	logger := &Logger{
		skipStackFrames: skipStackFrames,
		fatalStatusCode: fatalStatusCode,
		prefix:          prefix,
		outputs:         map[string]*Output{},
	}

	logger.SetOutputs(Stdout, Stderr)
	return logger
}

// Initializes the logger.
func init() {
	self = New()
	sik := skipStackFrames + 1 // self works at the imported package level
	self.SetSkipStackFrames(sik)
}

// Copy returns copy of the log object.
func Copy() *Logger {
	return self.Copy()
}

// SetSkipStackFrames sets skip stack frames level.
func SetSkipStackFrames(skips int) {
	self.SetSkipStackFrames(skips)
}

// SkipStackFrames returns skip stack frames level.
func SkipStackFrames() int {
	return self.SkipStackFrames()
}

// SetPrefix sets the name of the logger object.
func SetPrefix(prefix string) string {
	return self.SetPrefix(prefix)
}

// Prefix returns the name of the log object.
func Prefix() string {
	return self.Prefix()
}

// SetOutputs sets the outputs of the log object.
func SetOutputs(outputs ...Output) error {
	return self.SetOutputs(outputs...)
}

// EditOutputs edits the outputs of the log object.
func EditOutputs(outputs ...Output) error {
	return self.EditOutputs(outputs...)
}

// DeleteOutputs deletes the outputs of the log object.
func DeleteOutputs(names ...string) {
	self.DeleteOutputs(names...)
}

// Outputs returns a list of outputs.
func Outputs() []Output {
	return self.Outputs()
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Fpanic(w io.Writer, a ...any) {
	self.Fpanic(w, a...)
}

// Fpanicf creates message with Panic level, according to a format
// specifier and writes to w.
func Fpanicf(w io.Writer, format string, a ...any) {
	self.Fpanicf(w, format, a...)
}

// Fpanicln creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func Fpanicln(w io.Writer, a ...any) {
	self.Fpanicln(w, a...)
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func Panic(a ...any) {
	self.Panic(a...)
}

// Panicf creates message with Panic level, according to a format specifier
// and writes to log.Writer.
func Panicf(format string, a ...any) {
	self.Panicf(format, a...)
}

// Panicln creates message with Panic, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Panicln(a ...any) {
	self.Panicln(a...)
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Ffatal(w io.Writer, a ...any) {
	self.Ffatal(w, a...)
}

// Ffatalf creates message with Fatal level, according to a format
// specifier and writes to w.
func Ffatalf(w io.Writer, format string, a ...any) {
	self.Ffatalf(w, format, a...)
}

// Ffatalln creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ffatalln(w io.Writer, a ...any) {
	self.Ffatalln(w, a...)
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Fatal(a ...any) {
	self.Fatal(a...)
}

// Fatalf creates message with Fatal level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Fatalf(format string, a ...any) {
	self.Fatalf(format, a...)
}

// Fatalln creates message with Fatal, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Fatalln(a ...any) {
	self.Fatalln(a...)
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Ferror(w io.Writer, a ...any) {
	self.Ferror(w, a...)
}

// Ferrorf creates message with Error level, according to a format
// specifier and writes to w.
func Ferrorf(w io.Writer, format string, a ...any) {
	self.Ferrorf(w, format, a...)
}

// Ferrorln creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ferrorln(w io.Writer, a ...any) {
	self.Ferrorln(w, a...)
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Error(a ...any) {
	self.Error(a...)
}

// Errorf creates message with Error level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Errorf(format string, a ...any) {
	self.Errorf(format, a...)
}

// Errorln creates message with Error, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Errorln(a ...any) {
	self.Errorln(a...)
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Fwarn(w io.Writer, a ...any) {
	self.Fwarn(w, a...)
}

// Fwarnf creates message with Warn level, according to a format
// specifier and writes to w.
func Fwarnf(w io.Writer, format string, a ...any) {
	self.Fwarnf(w, format, a...)
}

// Fwarnln creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Fwarnln(w io.Writer, a ...any) {
	self.Fwarnln(w, a...)
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Warn(a ...any) {
	self.Warn(a...)
}

// Warnf creates message with Warn level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Warnf(format string, a ...any) {
	self.Warnf(format, a...)
}

// Warnln creates message with Warn, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Warnln(a ...any) {
	self.Warnln(a...)
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Finfo(w io.Writer, a ...any) {
	self.Finfo(w, a...)
}

// Finfof creates message with Info level, according to a format
// specifier and writes to w.
func Finfof(w io.Writer, format string, a ...any) {
	self.Finfof(w, format, a...)
}

// Finfoln creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Finfoln(w io.Writer, a ...any) {
	self.Finfoln(w, a...)
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Info(a ...any) {
	self.Info(a...)
}

// Infof creates message with Info level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Infof(format string, a ...any) {
	self.Infof(format, a...)
}

// Infoln creates message with Info, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Infoln(a ...any) {
	self.Infoln(a...)
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Fdebug(w io.Writer, a ...any) {
	self.Fdebug(w, a...)
}

// Fdebugf creates message with Debug level, according to a format
// specifier and writes to w.
func Fdebugf(w io.Writer, format string, a ...any) {
	self.Fdebugf(w, format, a...)
}

// Fdebugln creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Fdebugln(w io.Writer, a ...any) {
	self.Fdebugln(w, a...)
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Debug(a ...any) {
	self.Debug(a...)
}

// Debugf creates message with Debug level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Debugf(format string, a ...any) {
	self.Debugf(format, a...)
}

// Debugln creates message with Debug, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Debugln(a ...any) {
	self.Debugln(a...)
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Ftrace(w io.Writer, a ...any) {
	self.Ftrace(w, a...)
}

// Ftracef creates message with Trace level, according to a format
// specifier and writes to w.
func Ftracef(w io.Writer, format string, a ...any) {
	self.Ftracef(w, format, a...)
}

// Ftraceln creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ftraceln(w io.Writer, a ...any) {
	self.Ftraceln(w, a...)
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Trace(a ...any) {
	self.Trace(a...)
}

// Tracef creates message with Trace level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Tracef(format string, a ...any) {
	self.Tracef(format, a...)
}

// Traceln creates message with Trace, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Traceln(a ...any) {
	self.Traceln(a...)
}
