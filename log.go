package log

import (
	"io"
	"strings"
	"sync"

	"github.com/goloop/g"
)

var (
	self *Logger // is the default logger instance
	mu   sync.Mutex
)

// New returns a new Logger object with optional prefixes for log messages.
// If one or more prefixes are provided, they will be concatenated with
// hyphens and prepended to each log message. Leading and trailing whitespace
// characters are stripped from each prefix before concatenation.
//
// Example usage:
//
//	// Create a logger without any prefixes.
//	logger := log.New()
//	logger.Info("Hello, World!")
//
//	// Create a logger with a single prefix 'MYAPP'.
//	loggerWithPrefix := log.New("MYAPP")
//	loggerWithPrefix.Info("Hello, World!")
//
//	// Create a logger with multiple prefixes 'MYAPP' and 'PREFIX'.
//	// The prefixes will be joined as 'MYAPP-PREFIX'.
//	loggerWithMultiplePrefixes := log.New("MYAPP", "PREFIX")
//	loggerWithMultiplePrefixes.Info("Hello, World!")
//
//	// Create a logger with prefixes containing leading/trailing whitespace.
//	// The whitespace is removed and prefixes are joined as 'MYAPP-PREFIX'.
//	loggerWithWhitespace := log.New(" MYAPP ", " PREFIX ")
//	loggerWithWhitespace.Info("Hello, World!")
//
// Parameters:
//
//	prefixes: Optional string values to prepend to each log message.
//	          If multiple prefixes are provided, they will be concatenated
//	          with hyphens.
//
// Returns:
//
//	A new Logger instance configured with the provided prefixes.
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
	mu.Lock()
	defer mu.Unlock()

	if self == nil {
		self = New()
		skip := skipStackFrames + 1 // self works at the imported package level
		self.SetSkipStackFrames(skip)
	}
}

// Copy returns copy of the log object.
func Copy() *Logger {
	mu.Lock()
	defer mu.Unlock()
	return self.Copy()
}

// SetSkipStackFrames sets skip stack frames level.
func SetSkipStackFrames(skips int) {
	mu.Lock()
	defer mu.Unlock()
	self.SetSkipStackFrames(skips)
}

// SkipStackFrames returns skip stack frames level.
func SkipStackFrames() int {
	return self.SkipStackFrames()
}

// SetPrefix sets the name of the logger object.
func SetPrefix(prefix string) string {
	mu.Lock()
	defer mu.Unlock()
	return self.SetPrefix(prefix)
}

// Prefix returns the name of the log object.
func Prefix() string {
	mu.Lock()
	defer mu.Unlock()
	return self.Prefix()
}

// SetOutputs sets the outputs of the log object.
func SetOutputs(outputs ...Output) error {
	mu.Lock()
	defer mu.Unlock()
	return self.SetOutputs(outputs...)
}

// EditOutputs edits the outputs of the log object.
func EditOutputs(outputs ...Output) error {
	mu.Lock()
	defer mu.Unlock()
	return self.EditOutputs(outputs...)
}

// DeleteOutputs deletes the outputs of the log object.
func DeleteOutputs(names ...string) {
	mu.Lock()
	defer mu.Unlock()
	self.DeleteOutputs(names...)
}

// Outputs returns a list of outputs.
func Outputs(names ...string) []Output {
	mu.Lock()
	defer mu.Unlock()
	return self.Outputs(names...)
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Fpanic(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fpanic(w, a...)
}

// Fpanicf creates message with Panic level, according to a format
// specifier and writes to w.
func Fpanicf(w io.Writer, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fpanicf(w, format, a...)
}

// Fpanicln creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func Fpanicln(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fpanicln(w, a...)
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func Panic(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Panic(a...)
}

// Panicf creates message with Panic level, according to a format specifier
// and writes to log.Writer.
func Panicf(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Panicf(format, a...)
}

// Panicln creates message with Panic, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Panicln(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Panicln(a...)
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Ffatal(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ffatal(w, a...)
}

// Ffatalf creates message with Fatal level, according to a format
// specifier and writes to w.
func Ffatalf(w io.Writer, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ffatalf(w, format, a...)
}

// Ffatalln creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ffatalln(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ffatalln(w, a...)
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Fatal(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fatal(a...)
}

// Fatalf creates message with Fatal level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Fatalf(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fatalf(format, a...)
}

// Fatalln creates message with Fatal, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Fatalln(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fatalln(a...)
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Ferror(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ferror(w, a...)
}

// Ferrorf creates message with Error level, according to a format
// specifier and writes to w.
func Ferrorf(w io.Writer, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ferrorf(w, format, a...)
}

// Ferrorln creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ferrorln(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ferrorln(w, a...)
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Error(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Error(a...)
}

// Errorf creates message with Error level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Errorf(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Errorf(format, a...)
}

// Errorln creates message with Error, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Errorln(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Errorln(a...)
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Fwarn(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fwarn(w, a...)
}

// Fwarnf creates message with Warn level, according to a format
// specifier and writes to w.
func Fwarnf(w io.Writer, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fwarnf(w, format, a...)
}

// Fwarnln creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Fwarnln(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fwarnln(w, a...)
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Warn(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Warn(a...)
}

// Warnf creates message with Warn level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Warnf(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Warnf(format, a...)
}

// Warnln creates message with Warn, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Warnln(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Warnln(a...)
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Finfo(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Finfo(w, a...)
}

// Finfof creates message with Info level, according to a format
// specifier and writes to w.
func Finfof(w io.Writer, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Finfof(w, format, a...)
}

// Finfoln creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Finfoln(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Finfoln(w, a...)
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Info(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Info(a...)
}

// Infof creates message with Info level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Infof(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Infof(format, a...)
}

// Infoln creates message with Info, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Infoln(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Infoln(a...)
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Fdebug(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fdebug(w, a...)
}

// Fdebugf creates message with Debug level, according to a format
// specifier and writes to w.
func Fdebugf(w io.Writer, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fdebugf(w, format, a...)
}

// Fdebugln creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Fdebugln(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Fdebugln(w, a...)
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Debug(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Debug(a...)
}

// Debugf creates message with Debug level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Debugf(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Debugf(format, a...)
}

// Debugln creates message with Debug, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Debugln(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Debugln(a...)
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func Ftrace(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ftrace(w, a...)
}

// Ftracef creates message with Trace level, according to a format
// specifier and writes to w.
func Ftracef(w io.Writer, format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ftracef(w, format, a...)
}

// Ftraceln creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ftraceln(w io.Writer, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Ftraceln(w, a...)
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Trace(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Trace(a...)
}

// Tracef creates message with Trace level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Tracef(format string, a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Tracef(format, a...)
}

// Traceln creates message with Trace, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func Traceln(a ...any) {
	mu.Lock()
	defer mu.Unlock()
	self.Traceln(a...)
}
