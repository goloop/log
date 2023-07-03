package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/goloop/g"
	"github.com/goloop/log/layout"
	"github.com/goloop/log/level"
	"github.com/goloop/trit"
)

const (
	// The skipStackFrames specifies the default number of stack frames to
	// skip before the program counter stack is collected.
	skipStackFrames = 4

	// The shortPathSections is the number of sections in the short path
	// of the file. I.e. if it is set to display the path to the file as
	// short, this value determines how many rightmost sections will be
	// displayed,  e.g. for full path as /very/long/path/to/project/main.go
	// with shortPathSections as 3 will be displayed .../to/project/main.go.
	shortPathSections = 3

	// The fatalStatusCode is the default status code for the fatal level.
	fatalStatusCode = 1

	// The outWithPrefix is default values for the output prefix parameters.
	outWithPrefix = trit.True

	// The outWithColor is default values for the output color parameters.
	outWithColor = trit.False

	// The outText is default values for the output text parameters.
	outText = trit.True

	// The outEnabled is default values for the output enabled parameters.
	outEnabled = trit.True

	// The outSpace is the space between the blocks of the
	// output prefix.
	outSpace = " "

	// The outTimestampFormat is the default time format for the output.
	outTimestampFormat = "2006/01/02 15:04:05"

	// The outLevelFormat is the default level format for the output.
	outLevelFormat = "%s"
)

var (
	// Stdout standard rules for displaying logger information's
	// messages in the console.
	Stdout = Output{
		Name:            "stdout",
		Writer:          os.Stdout,
		Space:           outSpace,
		Layouts:         layout.Default,
		Levels:          level.Info | level.Debug | level.Trace,
		WithPrefix:      outWithPrefix,
		WithColor:       outWithColor,
		Enabled:         outEnabled,
		Text:            outText,
		TimestampFormat: outTimestampFormat,
		LevelFormat:     outLevelFormat,
	}

	// Stderr standard rules for displaying logger errors
	// messages in the console.
	Stderr = Output{
		Name:            "stderr",
		Writer:          os.Stderr,
		Space:           outSpace,
		Layouts:         layout.Default,
		Levels:          level.Panic | level.Fatal | level.Error | level.Warn,
		WithPrefix:      outWithPrefix,
		WithColor:       outWithColor,
		Enabled:         outEnabled,
		Text:            outText,
		TimestampFormat: outTimestampFormat,
		LevelFormat:     outLevelFormat,
	}

	// Default ...
	Default = Output{
		Name:            "default",
		Writer:          os.Stdout,
		Space:           outSpace,
		Layouts:         layout.Default,
		Levels:          Stdout.Levels | Stderr.Levels,
		WithPrefix:      outWithPrefix,
		WithColor:       outWithColor,
		Enabled:         outEnabled,
		Text:            outText,
		TimestampFormat: outTimestampFormat,
		LevelFormat:     outLevelFormat,
	}

	// The exit is ...
	exit = os.Exit
)

// Output is the type of the logging output configuration.
// Specifies the destination where the log and display
// parameters will be output.
type Output struct {
	// Name is the name of the output. It is used to identify
	// the output in the list of outputs.
	Name string

	// Writer is the point where the login data will be output,
	// for example os.Stdout or text file descriptor.
	Writer io.Writer

	// Space is the space between the blocks of the
	// output prefix.
	Space string

	// Layouts is the flag-holder where flags responsible for
	// formatting the log message prefix.
	Layouts layout.Layout

	// Levels is the flag-holder where flags responsible for
	// levels of the logging: Panic, Fatal, Error, Warn, Info etc.
	Levels level.Level

	// WithPrefix is the flag that determines whether to show the prefix
	// in the log-message. I.e., if the prefix is set for the logger,
	// it will be convenient for os.Stdout to display it, since several
	// applications can send messages at the same time, but this is
	// not necessary for the log file because each application has
	// its own log file.
	//
	// This only works when the prefix is set to the logger.
	// By default, the prefix is enabled.
	//
	// Values are given by numerical marks, where:
	//  - values less than zero are considered false;
	//  - values greater than zero are considered true;
	//  - the value set to 0 is considered the default value
	//    (or don't change, for edit mode).
	WithPrefix trit.Trit

	// WithColor is the flag that determines whether to use color for the
	// log-message. Each message level has its own color. This is handy for
	// the console to visually see the problem or debug information.
	//
	// By default, the color is disabled.
	//
	// Values are given by numerical marks, where:
	//  - values less than zero are considered false;
	//  - values greater than zero are considered true;
	//  - the value set to 0 is considered the default value
	//    (or don't change, for edit mode).
	//
	// The color scheme works only for UNIX-like systems.
	// The color scheme works for flat format only (i.e. display of log
	// messages in the form of text, not as JSON).
	WithColor trit.Trit

	// Enabled is the flag that determines whether to enable the output.
	//
	// By default, the new output is enable.
	//
	// Values are given by numerical marks, where:
	//  - values less than zero are considered false;
	//  - values greater than zero are considered true;
	//  - the value set to 0 is considered the default value
	//    (or don't change, for edit mode).
	Enabled trit.Trit

	// Text is the flag that determines whether to use text format for the
	// log-message. Otherwise, the result will be displayed in JSON format.
	//
	// Values are given by numerical marks, where:
	//  - values less than zero are considered false;
	//  - values greater than zero are considered true;
	//  - the value set to 0 is considered the default value
	//    (or don't change, for edit mode).
	Text trit.Trit

	// TimestampFormat is the format of the timestamp in the log-message.
	// Must be specified in the format of the time.Format() function.
	TimestampFormat string

	// LevelFormat is the format of the level in the log-message.
	// Allows you to add additional information or a label around
	// the ID of the level, for example, to display the level in
	// square brackets: [LEVEL_NAME] - we need to specify the
	// format as "[%s]".
	LevelFormat string
}

// Logger is a structure that encapsulates logging functionality.
// It provides an interface to log messages with varying levels
// of severity, and can be configured to format and output these
// messages in different ways.
type Logger struct {
	// The skipStackFrames specifies the number of stack frames to skip before
	// the program counter stack is collected. For example, if skip is equal
	// to 4, then the top four frames of the stack will be ignored.
	skipStackFrames int

	// The fatalStatusCode is the status code that will be returned
	// to the operating system when the program is terminated
	// by the Fatal() method.
	fatalStatusCode int

	// The prefix (optional) it is a special value that is inserted before
	// the log-message. It can be used to identify the application that
	// generates the message, if several different applications output
	// to one output.
	prefix string

	// The outputs is the list of the logging outputs.
	// The logger can output the log message to several outputs
	// at once, for example, to the console and to the log file.
	outputs map[string]*Output

	// The mu is the mutex for the log object.
	mu sync.RWMutex
}

// Copy returns copy of the logger object.
func (logger *Logger) Copy() *Logger {
	// Lock the log object for change.
	logger.mu.RLock()
	defer logger.mu.RUnlock()

	// Get outputs.
	outputs := make([]Output, 0, len(logger.outputs))
	for _, o := range logger.outputs {
		outputs = append(outputs, *o)
	}

	instance := &Logger{
		skipStackFrames: logger.skipStackFrames,
		fatalStatusCode: logger.fatalStatusCode,
		prefix:          logger.prefix,
		outputs:         map[string]*Output{},
	}

	instance.SetOutputs(outputs...)
	return instance
}

// SetSkipStackFrames sets the number of stack frames to skip before
// the program counter stack is collected.
//
// If the specified value is less than zero, the value does not change.
// If too large a value is specified, the maximum allowable value
// will be set.
func (logger *Logger) SetSkipStackFrames(skip int) int {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	// Cannot be a negative value.
	if skip < 0 {
		return logger.skipStackFrames
	}

	// Too big a skip can cause panic in the getStackFrame.
	// Take the highest possible value.
	for {
		func() {
			// If you panic, reduce the skip value by one.
			defer func() {
				if r := recover(); r != nil {
					skip--
				}
			}()

			// If the skip is too large, it can cause panic.
			logger.skipStackFrames = skip
			getStackFrame(skip + 1) // plus 1 because we in inner func
		}()

		// If the subtraction is already nowhere or
		// the skip didn't cause panic (was not reduced).
		if skip <= 0 || logger.skipStackFrames == skip {
			break
		}
	}

	return logger.skipStackFrames
}

// SkipStackFrames returns the number of stack frames to skip before
// the program counter stack is collected.
func (logger *Logger) SkipStackFrames() int {
	logger.mu.RLock()
	defer logger.mu.RUnlock()
	return logger.skipStackFrames
}

// SetPrefix sets the prefix to the logger. The prefix is a special value
// that is inserted before the log-message. It can be used to identify
// the application that generates the message, if several different
// applications output to one output.
func (logger *Logger) SetPrefix(prefix string) string {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.prefix = prefix
	return logger.prefix
}

// Prefix returns logger prefix.
func (logger *Logger) Prefix() string {
	logger.mu.RLock()
	defer logger.mu.RUnlock()
	return logger.prefix
}

// SetOutputs clears all installed outputs and installs new ones from the list.
// If new outputs are not specified, the logger will not output information
// anywhere. We must explicitly specify at least one output, or use the
// defaults: log.Stdout or/and log.Stderr.
//
// Example usage:
//
//	// Show only informational levels.
//	logger.SetOutputs(log.Stdout)
//
//	// Show only error and fatal levels.
//	logger.SetOutputs(log.Stderr)
//
//	// Show all levels.
//	logger.SetOutputs(log.Stdout, log.Stderr)
//
//	// Use custom settings.
//	// Output to the console is different from output to a file.
//	f, err := os.OpenFile("e.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer f.Close()
//
//	logger.SetOutputs(
//	   log.Output{
//	       Name:        "errors",
//	       Writer:      f,
//	       Levels:      log.PanicLevel|log.FatalLevel|log.ErrorLevel,
//	       Formats:     log.ShortPathFormat|log.LineNumberFormat,
//	       WithPrefix:  log.Yes,
//	   },
//	   log.Stdout,
//	   log.Stderr,
//	)
//
// If the specified output has an empty name or writer as nil, the function
// returns an error. An error is returned if two or more outputs have the
// same name. If the function returns an error, it doesn't change or clear
// the previously set values.
//
// Example usage:
//
//	// Show only informational levels.
//	logger.SetOutputs(log.Stdout)
//
//	// Try to update with incorrect data.
//	logger.SetOutputs(log.Output{}) // error: the 0 object has empty name
func (logger *Logger) SetOutputs(outputs ...Output) error {
	// Lock the logger.
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if len(outputs) == 0 {
		return fmt.Errorf("the outputs list is empty")
	}

	// Check the correctness of the data and create a temporary map.
	// If the data is not correct, we cannot change the data already
	// set previously.
	result := make(map[string]*Output, len(outputs))
	for i := range outputs {
		o := &outputs[i]

		// The name must be specified.
		if g.IsEmpty(o.Name) {
			return fmt.Errorf("the %d output has empty name", i)
		}

		// The writer must be specified.
		if o.Writer == nil {
			return fmt.Errorf("the %d output has nil writer", i)
		}

		// If the output is already in the list, then return an error.
		if _, ok := result[o.Name]; ok {
			return fmt.Errorf("output duplicate name '%s'", o.Name)
		}

		// Set default value for space.
		if g.IsEmpty(o.Space) {
			o.Space = outSpace
		}

		// Set default value for prefix.
		if o.WithPrefix.IsUnknown() {
			o.WithPrefix = outWithPrefix
		}

		// Set default value for color.
		if o.WithColor.IsUnknown() {
			o.WithColor = outWithColor
		}

		// Set default value for disabled.
		if o.Enabled.IsUnknown() {
			o.Enabled = outEnabled
		}

		// Set default value for text.
		if o.Text.IsUnknown() {
			o.Text = outText
		}

		// Set default value for timestamp.
		if g.IsEmpty(o.TimestampFormat) {
			o.TimestampFormat = outTimestampFormat
		}

		// Level format.
		if g.IsEmpty(o.LevelFormat) {
			o.LevelFormat = outLevelFormat
		}

		result[o.Name] = o
	}

	logger.outputs = result
	return nil
}

// EditOutputs updates the list of outputs.
// If the output with the specified name is not set,
// the function returns an error.
//
// To edit the output, we must specify its name and only those fields
// that will be edited. Fields that are not specified will not be changed.
//
// Example usage:
//
//	// Set default settings.
//	logger.SetOutputs(log.Stdout, log.Stderr)
//
//	// Change the settings for the "stdout" output.
//	logger.EditOutputs(log.Output{
//	    Name:    log.Stdout.Name,
//	    Levels:  log.PanicLevel|log.FatalLevel|log.ErrorLevel,
//	})
func (logger *Logger) EditOutputs(outputs ...Output) error {
	// Lock the logger.
	logger.mu.Lock()
	defer logger.mu.Unlock()

	if len(outputs) == 0 {
		return fmt.Errorf("the outputs list is empty")
	}

	// Check the correctness of the data and create a temporary map.
	// If the data is not correct, we cannot change the data already
	// set previously.
	result := make(map[string]*Output, len(outputs))
	for _, o := range outputs {
		out, ok := logger.outputs[o.Name]
		if !ok {
			return fmt.Errorf("output not found '%s'", o.Name)
		}

		// Update the output fields.
		// Only those values that are clearly indicated.
		if !g.IsEmpty(o.Writer) {
			out.Writer = o.Writer
		}

		if !g.IsEmpty(o.Space) {
			out.Space = o.Space
		}

		if !g.IsEmpty(o.Layouts) {
			out.Layouts = o.Layouts
		}

		if !g.IsEmpty(o.Levels) {
			out.Levels = o.Levels
		}

		if !g.IsEmpty(o.WithPrefix) {
			out.WithPrefix = o.WithPrefix
		}

		if !g.IsEmpty(o.WithColor) {
			out.WithColor = o.WithColor
		}

		// Setvalue for disabled.
		if !o.Enabled.IsUnknown() {
			out.Enabled = o.Enabled
		}

		// Set value for text.
		if !o.Text.IsUnknown() {
			out.Text = o.Text
		}

		// Set value for timestamp.
		if !g.IsEmpty(o.TimestampFormat) {
			out.TimestampFormat = o.TimestampFormat
		}

		// Level format.
		if !g.IsEmpty(o.LevelFormat) {
			out.LevelFormat = o.LevelFormat
		}

		result[o.Name] = out
	}

	// Update outputs.
	for n, o := range result {
		logger.outputs[n] = o
	}

	return nil
}

// DeleteOutputs deletes outputs by name.
//
// Example usage:
//
//	// Delete the "stdout" output.
//	logger.DeleteOutputs(log.Stdout.Name)
func (logger *Logger) DeleteOutputs(names ...string) {
	// Lock the logger.
	logger.mu.Lock()
	defer logger.mu.Unlock()

	for _, name := range names {
		delete(logger.outputs, name)
	}
}

// Outputs returns a list of outputs.
//
// Example usage:
//
//	// Get a list of the current outputs.
//	outputs := logger.Outputs()
//
//	// Add new output.
//	outputs = append(outputs, log.Stdout)
//
//	// Set new outputs.
//	logger.SetOutputs(outputs...)
func (logger *Logger) Outputs() []Output {
	logger.mu.RLock()
	defer logger.mu.RUnlock()

	result := make([]Output, 0, len(logger.outputs))
	for _, o := range logger.outputs {
		result = append(result, *o)
	}

	return result
}

// The echo is universal method creates a message of the fmt.Fprint format.
func (logger *Logger) echo(w io.Writer, l level.Level, f string, a ...any) {
	// Lock the log object for change.
	logger.mu.RLock()
	defer logger.mu.RUnlock()

	// Get the stack frame.
	sf := getStackFrame(logger.skipStackFrames)

	// If an additional value is set for the output (writer),
	// use it with the default settings.
	outputs := logger.outputs
	if w != nil {
		output := Default
		output.Writer = w
		outputs["*"] = &output
	}

	// Output message.
	for _, o := range logger.outputs {
		var msg string
		has, err := o.Levels.Contains(l)
		if !has || err != nil || !o.Enabled.IsTrue() {
			continue
		}

		if o.Text > 0 {
			msg = textMessage(logger.prefix, l, time.Now(), o, sf, f, a...)
		} else {
			msg = objectMessage(logger.prefix, l, time.Now(), o, sf, f, a...)
		}

		fmt.Fprint(o.Writer, msg)
	}
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Fpanic(w io.Writer, a ...any) {
	logger.echo(w, level.Panic, "%s", a...)
	panic(fmt.Sprint(a...))
}

// Fpanicf creates message with Panic level, according to a format
// specifier and writes to w.
func (logger *Logger) Fpanicf(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Panic, format, a...)
	panic(fmt.Sprintf(format, a...))
}

// Fpanicln creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func (logger *Logger) Fpanicln(w io.Writer, a ...any) {
	logger.echo(w, level.Panic, "%s\n", a...)
	panic(fmt.Sprintln(a...))
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Panic(a ...any) {
	logger.echo(nil, level.Panic, "%s", a...)
	panic(fmt.Sprint(a...))
}

// Panicf creates message with Panic level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Panicf(format string, a ...any) {
	logger.echo(nil, level.Panic, format, a...)
	panic(fmt.Sprintf(format, a...))
}

// Panicln creates message with Panic, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Panicln(a ...any) (int, error) {
	logger.echo(nil, level.Panic, "%s\n", a...)
	panic(fmt.Sprintln(a...))
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Ffatal(w io.Writer, a ...any) {
	logger.echo(w, level.Fatal, "%s", a...)
	exit(logger.fatalStatusCode)
}

// Ffatalf creates message with Fatal level, according to a format
// specifier and writes to w.
func (logger *Logger) Ffatalf(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Fatal, format, a...)
	exit(logger.fatalStatusCode)
}

// Ffatalln creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func (logger *Logger) Ffatalln(w io.Writer, a ...any) {
	logger.echo(w, level.Fatal, "%s\n", a...)
	exit(logger.fatalStatusCode)
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Fatal(a ...any) {
	logger.echo(nil, level.Fatal, "%s", a...)
	exit(logger.fatalStatusCode)
}

// Fatalf creates message with Fatal level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Fatalf(format string, a ...any) {
	logger.echo(nil, level.Fatal, format, a...)
	exit(logger.fatalStatusCode)
}

// Fatalln creates message with Fatal, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Fatalln(a ...any) {
	logger.echo(nil, level.Fatal, "%s\n", a...)
	exit(logger.fatalStatusCode)
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Ferror(w io.Writer, a ...any) {
	logger.echo(w, level.Error, "%s", a...)
}

// Ferrorf creates message with Error level, according to a format
// specifier and writes to w.
func (logger *Logger) Ferrorf(w io.Writer, f string, a ...any) {
	logger.echo(w, level.Error, f, a...)
}

// Ferrorln creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func (logger *Logger) Ferrorln(w io.Writer, a ...any) {
	logger.echo(w, level.Error, "%s\n", a...)
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Error(a ...any) {
	logger.echo(nil, level.Error, "%s", a...)
}

// Errorf creates message with Error level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Errorf(f string, a ...any) {
	logger.echo(nil, level.Error, f, a...)
}

// Errorln creates message with Error, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Errorln(a ...any) {
	logger.echo(nil, level.Error, "%s\n", a...)
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Fwarn(w io.Writer, a ...any) {
	logger.echo(w, level.Warn, "%s", a...)
}

// Fwarnf creates message with Warn level, according to a format
// specifier and writes to w.
func (logger *Logger) Fwarnf(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Warn, format, a...)
}

// Fwarnln creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func (logger *Logger) Fwarnln(w io.Writer, a ...any) {
	logger.echo(w, level.Warn, "%s\n", a...)
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Warn(a ...any) {
	logger.echo(nil, level.Warn, "%s", a...)
}

// Warnf creates message with Warn level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Warnf(format string, a ...any) {
	logger.echo(nil, level.Warn, format, a...)
}

// Warnln creates message with Warn, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Warnln(a ...any) {
	logger.echo(nil, level.Warn, "%s\n", a...)
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Finfo(w io.Writer, a ...any) {
	logger.echo(w, level.Info, "%s", a...)
}

// Finfof creates message with Info level, according to a format
// specifier and writes to w.
func (logger *Logger) Finfof(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Info, format, a...)
}

// Finfoln creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func (logger *Logger) Finfoln(w io.Writer, a ...any) {
	logger.echo(w, level.Info, "%s\n", a...)
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Info(a ...any) {
	logger.echo(nil, level.Info, "%s", a...)
}

// Infof creates message with Info level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Infof(format string, a ...any) {
	logger.echo(nil, level.Info, format, a...)
}

// Infoln creates message with Info, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Infoln(a ...any) {
	logger.echo(nil, level.Info, "%s\n", a...)
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Fdebug(w io.Writer, a ...any) {
	logger.echo(w, level.Debug, "%s", a...)
}

// Fdebugf creates message with Debug level, according to a format
// specifier and writes to w.
func (logger *Logger) Fdebugf(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Debug, format, a...)
}

// Fdebugln creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func (logger *Logger) Fdebugln(w io.Writer, a ...any) {
	logger.echo(w, level.Debug, "%s\n", a...)
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Debug(a ...any) {
	logger.echo(nil, level.Debug, "%s", a...)
}

// Debugf creates message with Debug level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (logger *Logger) Debugf(format string, a ...any) {
	logger.echo(nil, level.Debug, format, a...)
}

// Debugln creates message with Debug, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Debugln(a ...any) {
	logger.echo(nil, level.Debug, "%s\n", a...)
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Ftrace(w io.Writer, a ...any) {
	logger.echo(w, level.Trace, "%s", a...)
}

// Ftracef creates message with Trace level, according to a format
// specifier and writes to w.
func (logger *Logger) Ftracef(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Trace, format, a...)
}

// Ftraceln creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are always added between
// operands and a newline is appended.
func (logger *Logger) Ftraceln(w io.Writer, a ...any) {
	logger.echo(w, level.Trace, "%s\n", a...)
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Trace(a ...any) {
	logger.echo(nil, level.Trace, "%s", a...)
}

// Tracef creates message with Trace level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Tracef(format string, a ...any) {
	logger.echo(nil, level.Trace, format, a...)
}

// Traceln creates message with Trace, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Traceln(a ...any) {
	logger.echo(nil, level.Trace, "%s\n", a...)
}
