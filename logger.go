package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/goloop/g"
	"github.com/goloop/is"
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

	// The outTextStyle is default values for the output text parameters.
	outTextStyle = trit.True

	// The outEnabled is default values for the output enabled parameters.
	outEnabled = trit.True

	// The outSpace is the space between the blocks of the
	// output prefix.
	outSpace = " "

	// The outTimestampFormat is the default time format for the output.
	outTimestampFormat = "2006/01/02 15:04:05"

	// The outLevelFormat is the default level format for the output.
	outLevelFormat = "%s"

	// There is a difference between text-style message formatting
	// and JSON-style formatting.
	//
	// So, in text style: print simply prints all elements, pasting
	// them together without separators; println - adds spaces between
	// elements and adds \n at the end of the formed line; and printf -
	// outputs according to a special format template.
	//
	// At the same time, for JSON: print should insert all elements
	// without delimiters (as in text style), but all JSON blocks* don't
	// contain \n at the end; but println inserts all message elements
	// with a space added, and does not add \n to the end of the message,
	// instead, each new JSON block* is printed on a new line (\n is appended
	// to the JSON block); printf - \n is appended to the JSON block* and the
	// message format matches the specified user format.
	//
	// So, to determine if a custom format pattern is specified, we use
	// two constants formatPrint and formatPrintln, which help determine
	// the type of function (print, println, or printf) that renders
	// the message, and the formatting methods for text or JSON style.
	//
	//  * JSON block is the final text representation of the logger in
	//   JSON format written as text.

	// The formatPrint is the system format string
	// for print-type functions.
	formatPrint = "{$ string-without-a-newline-character $}"

	// The formatPrintln is the system format string
	// for println-type functions.
	formatPrintln = "{$ string-with-a-newline-character $}"
)

var (
	// Stdout standard rules for displaying logger information's
	// messages in the console.
	Stdout = Output{
		Name:            "stdout",
		Writer:          os.Stdout,
		Layouts:         layout.Default,
		Levels:          level.Info | level.Debug | level.Trace,
		Space:           outSpace,
		WithPrefix:      outWithPrefix,
		WithColor:       outWithColor,
		Enabled:         outEnabled,
		TextStyle:       outTextStyle,
		TimestampFormat: outTimestampFormat,
		LevelFormat:     outLevelFormat,
	}

	// Stderr standard rules for displaying logger errors
	// messages in the console.
	Stderr = Output{
		Name:            "stderr",
		Writer:          os.Stderr,
		Layouts:         layout.Default,
		Levels:          level.Panic | level.Fatal | level.Error | level.Warn,
		Space:           outSpace,
		WithPrefix:      outWithPrefix,
		WithColor:       outWithColor,
		Enabled:         outEnabled,
		TextStyle:       outTextStyle,
		TimestampFormat: outTimestampFormat,
		LevelFormat:     outLevelFormat,
	}

	// Default is output that processes all types of levels
	// and outputs the result to stdout.
	Default = Output{
		Name:            "default",
		Writer:          os.Stdout,
		Layouts:         layout.Default,
		Levels:          Stdout.Levels | Stderr.Levels,
		Space:           outSpace,
		WithPrefix:      outWithPrefix,
		WithColor:       outWithColor,
		Enabled:         outEnabled,
		TextStyle:       outTextStyle,
		TimestampFormat: outTimestampFormat,
		LevelFormat:     outLevelFormat,
	}

	// The exit causes the current program to exit with the given status code.
	// Redefined this function to be able to test *fatal* methods.
	exit = os.Exit
)

// Output is the type of the logging output configuration.
// Specifies the destination where the log and display
// parameters will be output.
type Output struct {
	// Name is the name of the output. It is used to identify
	// the output in the list of outputs.
	//
	// Mandatory parameter, cannot be empty.
	// The name must sound like a variable name in most programming languages:
	// it must have special characters, not start with a number, etc.
	// But the name can be a reserved word like return, for, if ... any.
	// The name can also be expressed as a selector name in CSS,
	// for example "some-log-file" (without the leading . or # symbols).
	Name string

	// Writer is the point where the login data will be output,
	// for example os.Stdout or text file descriptor.
	//
	// Mandatory parameter, cannot be empty.
	Writer io.Writer

	// Layouts is the flag-holder where flags responsible for
	// formatting the log message prefix.
	Layouts layout.Layout

	// Levels is the flag-holder where flags responsible for
	// levels of the logging: Panic, Fatal, Error, Warn, Info etc.
	Levels level.Level

	// Space is the space between the blocks of the
	// output prefix.
	Space string

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
	//
	// We can also use the github.com/goloop/trit package and
	// the trit.True or trit.False value.
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
	// We can also use the github.com/goloop/trit package and
	// the trit.True or trit.False value.
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
	//
	// We can also use the github.com/goloop/trit package and
	// the trit.True or trit.False value.
	Enabled trit.Trit

	// TextStyle is the flag that determines whether to use text style for
	// the log-message. Otherwise, the result will be displayed in JSON style.
	//
	// By default, the new output has text style.
	//
	// Values are given by numerical marks, where:
	//  - values less than zero are considered false;
	//  - values greater than zero are considered true;
	//  - the value set to 0 is considered the default value
	//    (or don't change, for edit mode).
	//
	// We can also use the github.com/goloop/trit package and
	// the trit.True or trit.False value.
	TextStyle trit.Trit

	// TimestampFormat is the format of the timestamp in the log-message.
	// Must be specified in the format of the time.Format() function.
	TimestampFormat string

	// LevelFormat is the format of the level in the log-message.
	// Allows us to add additional information or a label around
	// the ID of the level, for example, to display the level in
	// square brackets: [LEVEL_NAME] - we need to specify the
	// format as "[%s]".
	LevelFormat string

	// The isSystem is the flag that determines whether the output is system.
	// For example, this can be for all F* functions (Ferror, Finfo etc.) that
	// accept a target writer. Package generates a unique Output for them.
	isSystem bool
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
			// If panic, reduce the skip value by one.
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
		} else if !o.isSystem && !is.SelectorName(o.Name, true) {
			return fmt.Errorf("the %d output has incorrect name '%s'",
				i, o.Name)
		}

		// The writer must be specified.
		if g.IsEmpty(o.Writer) {
			return fmt.Errorf("the %d output has nil writer", i)
		}

		// If the output is already in the list, then return an error.
		if _, ok := result[o.Name]; ok {
			return fmt.Errorf("output duplicate name '%s'", o.Name)
		}

		// Set the new value if it is specified, otherwise set the default one.
		//
		// Note: g.Value returns the first non-empty value.
		o.Layouts = g.Value(o.Layouts, layout.Default)
		o.Levels = g.Value(o.Levels, level.Default)

		o.Space = g.Value(o.Space, outSpace)
		o.WithPrefix = g.Value(o.WithPrefix, outWithPrefix)
		o.WithColor = g.Value(o.WithColor, outWithColor)
		o.Enabled = g.Value(o.Enabled, outEnabled)
		o.TextStyle = g.Value(o.TextStyle, outTextStyle)
		o.TimestampFormat = g.Value(o.TimestampFormat, outTimestampFormat)
		o.LevelFormat = g.Value(o.LevelFormat, outLevelFormat)

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

		// Set the new value if it is specified, otherwise leave the old one.
		//
		// Note: g.Value returns the first non-empty value.
		out.Writer = g.Value(o.Writer, out.Writer)
		out.Layouts = g.Value(o.Layouts, out.Layouts)
		out.Levels = g.Value(o.Levels, out.Levels)

		out.Space = g.Value(o.Space, out.Space)
		out.WithPrefix = g.Value(o.WithPrefix, out.WithPrefix)
		out.WithColor = g.Value(o.WithColor, out.WithColor)
		out.Enabled = g.Value(o.Enabled, out.Enabled)
		out.TextStyle = g.Value(o.TextStyle, out.TextStyle)
		out.TimestampFormat = g.Value(o.TimestampFormat, out.TimestampFormat)
		out.LevelFormat = g.Value(o.LevelFormat, out.LevelFormat)

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
func (logger *Logger) Outputs(names ...string) []Output {
	logger.mu.RLock()
	defer logger.mu.RUnlock()

	// If the list of names is not empty, then we return only those outputs
	// that are specified in the list of names.
	if len(names) > 0 {
		result := make([]Output, 0, len(names))
		for _, name := range names {
			if o, ok := logger.outputs[name]; ok {
				result = append(result, *o)
			}
		}
		return result
	}

	// If the list of names is empty, then we return all outputs.
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
		output.isSystem = true
		outputs["*"] = &output // this name can be used for system names
	}

	// Output message.
	for _, o := range logger.outputs {
		var msg string
		has, err := o.Levels.Contains(l)
		if !has || err != nil || !o.Enabled.IsTrue() {
			continue
		}

		// Hide or show the prefix.
		prefix := logger.prefix
		if !o.WithPrefix.IsTrue() {
			prefix = ""
		}

		// Text or JSON representation of the message.
		if o.TextStyle.IsTrue() {
			msg = textMessage(prefix, l, time.Now(), o, sf, f, a...)
		} else {
			msg = objectMessage(prefix, l, time.Now(), o, sf, f, a...)
		}

		// Print message.
		fmt.Fprint(o.Writer, msg)
	}
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Fpanic(w io.Writer, a ...any) {
	logger.echo(w, level.Panic, formatPrint, a...)
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
	logger.echo(w, level.Panic, formatPrintln, a...)
	panic(fmt.Sprintln(a...))
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Panic(a ...any) {
	logger.echo(nil, level.Panic, formatPrint, a...)
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
	logger.echo(nil, level.Panic, formatPrintln, a...)
	panic(fmt.Sprintln(a...))
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Ffatal(w io.Writer, a ...any) {
	logger.echo(w, level.Fatal, formatPrint, a...)
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
	logger.echo(w, level.Fatal, formatPrintln, a...)
	exit(logger.fatalStatusCode)
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Fatal(a ...any) {
	logger.echo(nil, level.Fatal, formatPrint, a...)
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
	logger.echo(nil, level.Fatal, formatPrintln, a...)
	exit(logger.fatalStatusCode)
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Ferror(w io.Writer, a ...any) {
	logger.echo(w, level.Error, formatPrint, a...)
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
	logger.echo(w, level.Error, formatPrintln, a...)
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Error(a ...any) {
	logger.echo(nil, level.Error, formatPrint, a...)
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
	logger.echo(nil, level.Error, formatPrintln, a...)
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Fwarn(w io.Writer, a ...any) {
	logger.echo(w, level.Warn, formatPrint, a...)
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
	logger.echo(w, level.Warn, formatPrintln, a...)
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Warn(a ...any) {
	logger.echo(nil, level.Warn, formatPrint, a...)
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
	logger.echo(nil, level.Warn, formatPrintln, a...)
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Finfo(w io.Writer, a ...any) {
	logger.echo(w, level.Info, formatPrint, a...)
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
	logger.echo(w, level.Info, formatPrintln, a...)
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Info(a ...any) {
	logger.echo(nil, level.Info, formatPrint, a...)
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
	logger.echo(nil, level.Info, formatPrintln, a...)
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Fdebug(w io.Writer, a ...any) {
	logger.echo(w, level.Debug, formatPrint, a...)
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
	logger.echo(w, level.Debug, formatPrintln, a...)
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Debug(a ...any) {
	logger.echo(nil, level.Debug, formatPrint, a...)
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
	logger.echo(nil, level.Debug, formatPrintln, a...)
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes to w. Spaces are added between operands
// when neither is a string.
func (logger *Logger) Ftrace(w io.Writer, a ...any) {
	logger.echo(w, level.Trace, formatPrint, a...)
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
	logger.echo(w, level.Trace, formatPrintln, a...)
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Trace(a ...any) {
	logger.echo(nil, level.Trace, formatPrint, a...)
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
	logger.echo(nil, level.Trace, formatPrintln, a...)
}
