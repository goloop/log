package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goloop/g"
	"github.com/goloop/is"
	"github.com/goloop/log/v2/layout"
	"github.com/goloop/log/v2/level"
	"github.com/goloop/trit"
)

const (
	// The skipStackFrames is the default number of additional caller frames
	// to skip when capturing the call site (0 means the direct caller). The
	// logger's own internal frames are skipped automatically.
	skipStackFrames = 0

	// The maxSkipStackFrames caps SetSkipStackFrames to a sane upper bound;
	// captureFrame walks a bounded number of frames anyway.
	maxSkipStackFrames = 32

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

	// The errorHandler, if set, is called when writing a log message to an
	// output fails. When nil the logger is best-effort (errors are ignored).
	errorHandler func(o Output, n int, err error)

	// The levelMask is the cached union of all enabled outputs' level masks.
	// It is read lock-free at the start of emit to skip, without taking the
	// lock, any level that no configured output is interested in.
	levelMask atomic.Uint32

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
		errorHandler:    logger.errorHandler,
		outputs:         map[string]*Output{},
	}

	instance.SetOutputs(outputs...)
	return instance
}

// The isNilWriter reports whether w is unusable: an untyped nil, or a typed
// nil pointer/interface/map/slice/channel/func. A zero-size writer such as
// io.Discard is a valid writer and is accepted.
func isNilWriter(w io.Writer) bool {
	if w == nil {
		return true
	}

	switch v := reflect.ValueOf(w); v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Map,
		reflect.Slice, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return false
	}
}

// SetErrorHandler sets a function that is called whenever writing a log
// message to one of the outputs fails. Passing nil restores the default
// best-effort behaviour (write errors are ignored). The handler is invoked
// outside the logger's lock, so it may safely call back into the logger.
func (logger *Logger) SetErrorHandler(handler func(o Output, n int, err error)) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.errorHandler = handler
}

// The recomputeLevelMask refreshes the cached union of all enabled outputs'
// level masks. Callers must hold the write lock.
func (logger *Logger) recomputeLevelMask() {
	var mask level.Level
	for _, o := range logger.outputs {
		if o.Enabled.IsTrue() {
			mask |= o.Levels
		}
	}
	logger.levelMask.Store(uint32(mask))
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

	// Store the value as given, capped at a sane upper bound. The call site
	// is located by skipping the logger's own frames automatically, so this
	// is purely the number of extra user wrapper frames to skip.
	if skip > maxSkipStackFrames {
		skip = maxSkipStackFrames
	}

	logger.skipStackFrames = skip
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
		if isNilWriter(o.Writer) {
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
	logger.recomputeLevelMask()
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

	// Copy-on-write: build a new map and new Output values for the edited
	// entries, so the previously published outputs are never mutated in
	// place. This lets emit read them without holding the lock. On any error
	// the new map is discarded and the current state is preserved.
	next := make(map[string]*Output, len(logger.outputs))
	for name, o := range logger.outputs {
		next[name] = o
	}

	for _, in := range outputs {
		old, ok := next[in.Name]
		if !ok {
			return fmt.Errorf("output not found '%s'", in.Name)
		}

		// Set the new value if it is specified, otherwise keep the old one.
		//
		// Note: g.Value returns the first non-empty value.
		edited := *old
		edited.Writer = g.Value(in.Writer, edited.Writer)
		edited.Layouts = g.Value(in.Layouts, edited.Layouts)
		edited.Levels = g.Value(in.Levels, edited.Levels)

		edited.Space = g.Value(in.Space, edited.Space)
		edited.WithPrefix = g.Value(in.WithPrefix, edited.WithPrefix)
		edited.WithColor = g.Value(in.WithColor, edited.WithColor)
		edited.Enabled = g.Value(in.Enabled, edited.Enabled)
		edited.TextStyle = g.Value(in.TextStyle, edited.TextStyle)
		edited.TimestampFormat = g.Value(in.TimestampFormat, edited.TimestampFormat)
		edited.LevelFormat = g.Value(in.LevelFormat, edited.LevelFormat)

		next[in.Name] = &edited
	}

	logger.outputs = next
	logger.recomputeLevelMask()
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

	if len(names) == 0 {
		return
	}

	// Copy-on-write so an unlocked read in emit is never affected by the
	// deletion (the published map is replaced, not mutated in place).
	next := make(map[string]*Output, len(logger.outputs))
	for n, o := range logger.outputs {
		next[n] = o
	}
	for _, name := range names {
		delete(next, name)
	}

	logger.outputs = next
	logger.recomputeLevelMask()
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

// emitKind selects how operands are rendered into the message body:
// print joins them without separators, println adds spaces and a trailing
// newline, printf applies a user-supplied format string. The body is
// rendered once by the calling method, so echo and the message builders
// never forward a user format into fmt.Sprintf (which keeps go vet quiet).
type emitKind uint8

const (
	kindPrint   emitKind = iota // operands joined without separators
	kindPrintln                 // spaces between operands + trailing newline
	kindPrintf                  // user-supplied format string
)

// logField is a single structured key/value pair carried alongside a log
// message (used by the slog bridge). The value keeps its original type so
// JSON outputs can render it as a typed field rather than text.
type logField struct {
	key string
	val any
}

// The echo writes a pre-rendered message body to every configured output
// whose level mask contains l and which is enabled. When w is not nil, the
// message is additionally written to an ad-hoc output backed by w using the
// Default settings; logger.outputs itself is never mutated.
func (logger *Logger) echo(w io.Writer, l level.Level, kind emitKind, body string) {
	logger.emit(w, l, kind, body, nil, nil)
}

// The emit is the core of echo. When frame is non-nil it is used as the
// stack frame instead of capturing one — the slog bridge passes the real
// call site carried in the record's program counter.
func (logger *Logger) emit(
	w io.Writer,
	l level.Level,
	kind emitKind,
	body string,
	frame *stackFrame,
	fields []logField,
) {
	// Fast path: when logging only to the configured outputs (w == nil), skip
	// the whole call — no lock, no snapshot — if no enabled output is
	// interested in this level. levelMask is a lock-free cached union.
	if w == nil && level.Level(logger.levelMask.Load())&l != l {
		return
	}

	// Snapshot the immutable outputs map and config under the read lock,
	// then release it and write outside the lock: a slow or re-entrant
	// writer must not block configuration or hold the lock. The map and its
	// entries are never mutated after publication (SetOutputs, EditOutputs
	// and DeleteOutputs replace them copy-on-write), so reading them without
	// the lock is safe.
	logger.mu.RLock()
	outputs := logger.outputs
	prefix := logger.prefix
	skip := logger.skipStackFrames
	handler := logger.errorHandler
	logger.mu.RUnlock()

	// The ad-hoc writer (Fxxx) is an extra local target with the default
	// settings; it never touches logger.outputs.
	var adhoc *Output
	if w != nil {
		a := Default
		a.Writer = w
		a.isSystem = true
		adhoc = &a
	}

	// First pass: is any output interested, and does any need a stack frame?
	// Time and the stack frame are skipped entirely when nobody is interested.
	interested, needFrame := false, false
	for _, o := range outputs {
		if o.Levels&l == l && o.Enabled.IsTrue() {
			interested = true
			needFrame = needFrame || o.Layouts != 0
		}
	}
	if adhoc != nil && adhoc.Levels&l == l && adhoc.Enabled.IsTrue() {
		interested = true
		needFrame = needFrame || adhoc.Layouts != 0
	}

	if !interested {
		return
	}

	// One timestamp (and at most one stack frame) shared by all outputs.
	now := time.Now()
	sf := emptyFrame
	if needFrame {
		if frame != nil {
			sf = frame
		} else {
			sf, _ = captureFrame(skip)
		}
	}

	// Second pass: render and write each message outside the lock.
	for _, o := range outputs {
		writeOutput(o, l, kind, prefix, now, sf, body, fields, handler)
	}
	if adhoc != nil {
		writeOutput(adhoc, l, kind, prefix, now, sf, body, fields, handler)
	}
}

// The writeOutput renders the message for one output into a pooled buffer
// and writes the raw bytes, reporting any write error to handler. It runs
// outside the logger lock.
func writeOutput(
	o *Output,
	l level.Level,
	kind emitKind,
	prefix string,
	now time.Time,
	sf *stackFrame,
	body string,
	fields []logField,
	handler func(o Output, n int, err error),
) {
	if o.Levels&l != l || !o.Enabled.IsTrue() {
		return
	}

	p := prefix
	if !o.WithPrefix.IsTrue() {
		p = ""
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	if o.TextStyle.IsTrue() {
		appendText(buf, p, l, now, o, sf, body, fields)
	} else {
		appendObject(buf, p, l, now, o, sf, kind, body, fields)
	}

	n, err := o.Writer.Write(buf.Bytes())
	bufPool.Put(buf)

	if err != nil && handler != nil {
		handler(*o, n, err)
	}
}

// Enabled reports whether at least one enabled output would emit a message
// at level l. Use it to guard the preparation of expensive log arguments
// that should be skipped when no output is interested in the level:
//
//	if logger.Enabled(level.Debug) {
//		logger.Debug(expensiveDump())
//	}
func (logger *Logger) Enabled(l level.Level) bool {
	logger.mu.RLock()
	defer logger.mu.RUnlock()

	for _, o := range logger.outputs {
		if o.Levels&l == l && o.Enabled.IsTrue() {
			return true
		}
	}

	return false
}

// Fpanic creates message with Panic level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are added between operands when neither is a string.
func (logger *Logger) Fpanic(w io.Writer, a ...any) {
	msg := fmt.Sprint(a...)
	logger.echo(w, level.Panic, kindPrint, msg)
	panic(msg)
}

// Fpanicf creates message with Panic level, according to a format
// specifier and writes to the configured outputs and additionally to w.
func (logger *Logger) Fpanicf(w io.Writer, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	logger.echo(w, level.Panic, kindPrintf, msg)
	panic(msg)
}

// Fpanicln creates message with Panic level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are always added between operands and a newline is appended.
func (logger *Logger) Fpanicln(w io.Writer, a ...any) {
	msg := fmt.Sprintln(a...)
	logger.echo(w, level.Panic, kindPrintln, msg)
	panic(msg)
}

// Panic creates message with Panic level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Panic(a ...any) {
	msg := fmt.Sprint(a...)
	logger.echo(nil, level.Panic, kindPrint, msg)
	panic(msg)
}

// Panicf creates message with Panic level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Panicf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	logger.echo(nil, level.Panic, kindPrintf, msg)
	panic(msg)
}

// Panicln creates message with Panic, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Panicln(a ...any) {
	msg := fmt.Sprintln(a...)
	logger.echo(nil, level.Panic, kindPrintln, msg)
	panic(msg)
}

// Ffatal creates message with Fatal level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are added between operands when neither is a string.
func (logger *Logger) Ffatal(w io.Writer, a ...any) {
	logger.echo(w, level.Fatal, kindPrint, fmt.Sprint(a...))
	exit(logger.fatalStatusCode)
}

// Ffatalf creates message with Fatal level, according to a format
// specifier and writes to the configured outputs and additionally to w.
func (logger *Logger) Ffatalf(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Fatal, kindPrintf, fmt.Sprintf(format, a...))
	exit(logger.fatalStatusCode)
}

// Ffatalln creates message with Fatal level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are always added between operands and a newline is appended.
func (logger *Logger) Ffatalln(w io.Writer, a ...any) {
	logger.echo(w, level.Fatal, kindPrintln, fmt.Sprintln(a...))
	exit(logger.fatalStatusCode)
}

// Fatal creates message with Fatal level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Fatal(a ...any) {
	logger.echo(nil, level.Fatal, kindPrint, fmt.Sprint(a...))
	exit(logger.fatalStatusCode)
}

// Fatalf creates message with Fatal level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Fatalf(format string, a ...any) {
	logger.echo(nil, level.Fatal, kindPrintf, fmt.Sprintf(format, a...))
	exit(logger.fatalStatusCode)
}

// Fatalln creates message with Fatal, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Fatalln(a ...any) {
	logger.echo(nil, level.Fatal, kindPrintln, fmt.Sprintln(a...))
	exit(logger.fatalStatusCode)
}

// Ferror creates message with Error level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are added between operands when neither is a string.
func (logger *Logger) Ferror(w io.Writer, a ...any) {
	logger.echo(w, level.Error, kindPrint, fmt.Sprint(a...))
}

// Ferrorf creates message with Error level, according to a format
// specifier and writes to the configured outputs and additionally to w.
func (logger *Logger) Ferrorf(w io.Writer, f string, a ...any) {
	logger.echo(w, level.Error, kindPrintf, fmt.Sprintf(f, a...))
}

// Ferrorln creates message with Error level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are always added between operands and a newline is appended.
func (logger *Logger) Ferrorln(w io.Writer, a ...any) {
	logger.echo(w, level.Error, kindPrintln, fmt.Sprintln(a...))
}

// Error creates message with Error level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Error(a ...any) {
	logger.echo(nil, level.Error, kindPrint, fmt.Sprint(a...))
}

// Errorf creates message with Error level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Errorf(f string, a ...any) {
	logger.echo(nil, level.Error, kindPrintf, fmt.Sprintf(f, a...))
}

// Errorln creates message with Error, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Errorln(a ...any) {
	logger.echo(nil, level.Error, kindPrintln, fmt.Sprintln(a...))
}

// Fwarn creates message with Warn level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are added between operands when neither is a string.
func (logger *Logger) Fwarn(w io.Writer, a ...any) {
	logger.echo(w, level.Warn, kindPrint, fmt.Sprint(a...))
}

// Fwarnf creates message with Warn level, according to a format
// specifier and writes to the configured outputs and additionally to w.
func (logger *Logger) Fwarnf(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Warn, kindPrintf, fmt.Sprintf(format, a...))
}

// Fwarnln creates message with Warn level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are always added between operands and a newline is appended.
func (logger *Logger) Fwarnln(w io.Writer, a ...any) {
	logger.echo(w, level.Warn, kindPrintln, fmt.Sprintln(a...))
}

// Warn creates message with Warn level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Warn(a ...any) {
	logger.echo(nil, level.Warn, kindPrint, fmt.Sprint(a...))
}

// Warnf creates message with Warn level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Warnf(format string, a ...any) {
	logger.echo(nil, level.Warn, kindPrintf, fmt.Sprintf(format, a...))
}

// Warnln creates message with Warn, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Warnln(a ...any) {
	logger.echo(nil, level.Warn, kindPrintln, fmt.Sprintln(a...))
}

// Finfo creates message with Info level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are added between operands when neither is a string.
func (logger *Logger) Finfo(w io.Writer, a ...any) {
	logger.echo(w, level.Info, kindPrint, fmt.Sprint(a...))
}

// Finfof creates message with Info level, according to a format
// specifier and writes to the configured outputs and additionally to w.
func (logger *Logger) Finfof(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Info, kindPrintf, fmt.Sprintf(format, a...))
}

// Finfoln creates message with Info level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are always added between operands and a newline is appended.
func (logger *Logger) Finfoln(w io.Writer, a ...any) {
	logger.echo(w, level.Info, kindPrintln, fmt.Sprintln(a...))
}

// Info creates message with Info level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Info(a ...any) {
	logger.echo(nil, level.Info, kindPrint, fmt.Sprint(a...))
}

// Infof creates message with Info level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Infof(format string, a ...any) {
	logger.echo(nil, level.Info, kindPrintf, fmt.Sprintf(format, a...))
}

// Infoln creates message with Info, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Infoln(a ...any) {
	logger.echo(nil, level.Info, kindPrintln, fmt.Sprintln(a...))
}

// Fdebug creates message with Debug level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are added between operands when neither is a string.
func (logger *Logger) Fdebug(w io.Writer, a ...any) {
	logger.echo(w, level.Debug, kindPrint, fmt.Sprint(a...))
}

// Fdebugf creates message with Debug level, according to a format
// specifier and writes to the configured outputs and additionally to w.
func (logger *Logger) Fdebugf(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Debug, kindPrintf, fmt.Sprintf(format, a...))
}

// Fdebugln creates message with Debug level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are always added between operands and a newline is appended.
func (logger *Logger) Fdebugln(w io.Writer, a ...any) {
	logger.echo(w, level.Debug, kindPrintln, fmt.Sprintln(a...))
}

// Debug creates message with Debug level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Debug(a ...any) {
	logger.echo(nil, level.Debug, kindPrint, fmt.Sprint(a...))
}

// Debugf creates message with Debug level, according to a format specifier
// and writes to log.Writer. It returns the number of bytes written and any
// write error encountered.
func (logger *Logger) Debugf(format string, a ...any) {
	logger.echo(nil, level.Debug, kindPrintf, fmt.Sprintf(format, a...))
}

// Debugln creates message with Debug, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Debugln(a ...any) {
	logger.echo(nil, level.Debug, kindPrintln, fmt.Sprintln(a...))
}

// Ftrace creates message with Trace level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are added between operands when neither is a string.
func (logger *Logger) Ftrace(w io.Writer, a ...any) {
	logger.echo(w, level.Trace, kindPrint, fmt.Sprint(a...))
}

// Ftracef creates message with Trace level, according to a format
// specifier and writes to the configured outputs and additionally to w.
func (logger *Logger) Ftracef(w io.Writer, format string, a ...any) {
	logger.echo(w, level.Trace, kindPrintf, fmt.Sprintf(format, a...))
}

// Ftraceln creates message with Trace level, using the default formats
// for its operands and writes them to the configured outputs and to w.
// Spaces are always added between operands and a newline is appended.
func (logger *Logger) Ftraceln(w io.Writer, a ...any) {
	logger.echo(w, level.Trace, kindPrintln, fmt.Sprintln(a...))
}

// Trace creates message with Trace level, using the default formats
// for its operands and writes to log.Writer. Spaces are added between
// operands when neither is a string.
func (logger *Logger) Trace(a ...any) {
	logger.echo(nil, level.Trace, kindPrint, fmt.Sprint(a...))
}

// Tracef creates message with Trace level, according to a format specifier
// and writes to log.Writer.
func (logger *Logger) Tracef(format string, a ...any) {
	logger.echo(nil, level.Trace, kindPrintf, fmt.Sprintf(format, a...))
}

// Traceln creates message with Trace, level using the default formats
// for its operands and writes to log.Writer. Spaces are always added
// between operands and a newline is appended.
func (logger *Logger) Traceln(a ...any) {
	logger.echo(nil, level.Trace, kindPrintln, fmt.Sprintln(a...))
}
