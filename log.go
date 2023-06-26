package log

import (
	"io"
	"os"
)

const (
	// None means nothing.
	None = 0

	// The skipStackFrames specifies the number of stack frames to skip before
	// the program counter stack is collected.
	skipStackFrames = 4

	// The shortPathSections is the number of sections
	// in the short path of the file.
	shortPathSections = 3
)

var self *Logger

// New returns new Logger object.
func New(prefixes ...string) *Logger {
	logger := &Logger{
		skipStackFrames: skipStackFrames,
		name:            "",
		writer:          os.Stdout, // the os.Stdout is default writer.
		config: &Config{
			Formats:         FormatConfig(DefaultFormat),
			Levels:          LevelConfig(DefaultLevel),
			FatalStatusCode: 1,
			Prefix: &PrefixConfig{
				TimestampFormat:   TimestampFormat,
				SpaceBetweenCells: SpaceBetweenCells,
				LevelFormat:       LevelFormatConfig{},
			},
		},
	}
	logger.SetName(prefixes...)
	return logger
}

// Initializes the logger.
func init() {
	self = New()
	self.SetSkip(skipStackFrames + 1) // because self works at the imported package level
}

// Copy returns copy of the log object.
func Copy() *Logger {
	return self.Copy()
}

// SetSkip returns the skip value of the log object.
func SetSkip(skips ...int) int {
	return self.SetSkip(skips...)
}

// SetWriter sets the writer of the log object.
func SetWriter(writers ...io.Writer) io.Writer {
	return self.SetWriter(writers...)
}

// SetName returns the name of the log object.
func SetName(prefixes ...string) string {
	return self.SetName(prefixes...)
}

// Configure sets the configuration of the log object.
func Configure(config *Config) {
	self.Configure(config)
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func Fpanic(w io.Writer, a ...interface{}) (int, error) {
	return self.Fpanic(w, a...)
}

// Fpanicf creates message with Panic level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func Fpanicf(w io.Writer, format string, a ...interface{}) (int, error) {
	return self.Fpanicf(w, format, a...)
}

// Fpanicln creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Fpanicln(w io.Writer, a ...interface{}) (int, error) {
	return self.Fpanicln(w, a...)
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Panic(a ...interface{}) (int, error) {
	return self.Panic(a...)
}

// Panicf creates message with Panic level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Panicf(format string, a ...interface{}) (int, error) {
	return self.Panicf(format, a...)
}

// Panicln creates message with Panic, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func Panicln(a ...interface{}) (int, error) {
	return self.Panicln(a...)
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func Ffatal(w io.Writer, a ...interface{}) (int, error) {
	return self.Ffatal(w, a...)
}

// Ffatalf creates message with Fatal level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func Ffatalf(w io.Writer, format string, a ...interface{}) (int, error) {
	return self.Ffatalf(w, format, a...)
}

// Ffatalln creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ffatalln(w io.Writer, a ...interface{}) (int, error) {
	return self.Ffatalln(w, a...)
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Fatal(a ...interface{}) (int, error) {
	return self.Fatal(a...)
}

// Fatalf creates message with Fatal level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Fatalf(format string, a ...interface{}) (int, error) {
	return self.Fatalf(format, a...)
}

// Fatalln creates message with Fatal, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func Fatalln(a ...interface{}) (int, error) {
	return self.Fatalln(a...)
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func Ferror(w io.Writer, a ...interface{}) (int, error) {
	return self.Ferror(w, a...)
}

// Ferrorf creates message with Error level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func Ferrorf(w io.Writer, format string, a ...interface{}) (int, error) {
	return self.Ferrorf(w, format, a...)
}

// Ferrorln creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ferrorln(w io.Writer, a ...interface{}) (int, error) {
	return self.Ferrorln(w, a...)
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Error(a ...interface{}) (int, error) {
	return self.Error(a...)
}

// Errorf creates message with Error level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Errorf(format string, a ...interface{}) (int, error) {
	return self.Errorf(format, a...)
}

// Errorln creates message with Error, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func Errorln(a ...interface{}) (int, error) {
	return self.Errorln(a...)
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func Fwarn(w io.Writer, a ...interface{}) (int, error) {
	return self.Fwarn(w, a...)
}

// Fwarnf creates message with Warn level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func Fwarnf(w io.Writer, format string, a ...interface{}) (int, error) {
	return self.Fwarnf(w, format, a...)
}

// Fwarnln creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Fwarnln(w io.Writer, a ...interface{}) (int, error) {
	return self.Fwarnln(w, a...)
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Warn(a ...interface{}) (int, error) {
	return self.Warn(a...)
}

// Warnf creates message with Warn level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Warnf(format string, a ...interface{}) (int, error) {
	return self.Warnf(format, a...)
}

// Warnln creates message with Warn, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func Warnln(a ...interface{}) (int, error) {
	return self.Warnln(a...)
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func Finfo(w io.Writer, a ...interface{}) (int, error) {
	return self.Finfo(w, a...)
}

// Finfof creates message with Info level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func Finfof(w io.Writer, format string, a ...interface{}) (int, error) {
	return self.Finfof(w, format, a...)
}

// Finfoln creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Finfoln(w io.Writer, a ...interface{}) (int, error) {
	return self.Finfoln(w, a...)
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Info(a ...interface{}) (int, error) {
	return self.Info(a...)
}

// Infof creates message with Info level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Infof(format string, a ...interface{}) (int, error) {
	return self.Infof(format, a...)
}

// Infoln creates message with Info, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func Infoln(a ...interface{}) (int, error) {
	return self.Infoln(a...)
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func Fdebug(w io.Writer, a ...interface{}) (int, error) {
	return self.Fdebug(w, a...)
}

// Fdebugf creates message with Debug level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func Fdebugf(w io.Writer, format string, a ...interface{}) (int, error) {
	return self.Fdebugf(w, format, a...)
}

// Fdebugln creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Fdebugln(w io.Writer, a ...interface{}) (int, error) {
	return self.Fdebugln(w, a...)
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Debug(a ...interface{}) (int, error) {
	return self.Debug(a...)
}

// Debugf creates message with Debug level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Debugf(format string, a ...interface{}) (int, error) {
	return self.Debugf(format, a...)
}

// Debugln creates message with Debug, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func Debugln(a ...interface{}) (int, error) {
	return self.Debugln(a...)
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string. It returns the number of bytes written
// and any write error encountered.
func Ftrace(w io.Writer, a ...interface{}) (int, error) {
	return self.Ftrace(w, a...)
}

// Ftracef creates message with Trace level, according to a format
// specifier and writes to w. It returns the number of bytes written
// and any write error encountered.
func Ftracef(w io.Writer, format string, a ...interface{}) (int, error) {
	return self.Ftracef(w, format, a...)
}

// Ftraceln creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended. It returns the number of bytes
// written and any write error encountered.
func Ftraceln(w io.Writer, a ...interface{}) (int, error) {
	return self.Ftraceln(w, a...)
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string. It returns the number of bytes
// written and any write error encountered.
func Trace(a ...interface{}) (int, error) {
	return self.Trace(a...)
}

// Tracef creates message with Trace level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func Tracef(format string, a ...interface{}) (int, error) {
	return self.Tracef(format, a...)
}

// Traceln creates message with Trace, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended. It returns the number
// of bytes written and any write error encountered.
func Traceln(a ...interface{}) (int, error) {
	return self.Traceln(a...)
}
