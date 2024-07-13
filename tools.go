package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/goloop/g"
	"github.com/goloop/log/level"
)

// The stackFrame contains the top-level trace information
// where the logging method was called.
type stackFrame struct {
	FileLine    int     // line number
	FuncName    string  // function name
	FuncAddress uintptr // address of the function
	FilePath    string  // file path
}

// The ioCopy function is used to copy the output of a reader
// to a channel.
func ioCopy(r io.Reader, c chan string) {
	var buf bytes.Buffer
	io.Copy(&buf, r)
	c <- buf.String()
}

// The getStackFrame returns the stack slice. The skip argument
// is the number of stack frames to skip before taking a slice.
func getStackFrame(skip int) *stackFrame {
	sf := &stackFrame{}

	// Return program counters of function invocations on
	// the calling goroutine's stack and skipping function
	// call frames inside *Log.
	pc := make([]uintptr, skip+1) // program counters
	runtime.Callers(skip, pc)

	// Get a function at an address on the stack.
	fn := runtime.FuncForPC(pc[0])

	// Get name, path and line of the file.
	sf.FuncName = fn.Name()
	sf.FuncAddress = fn.Entry()
	sf.FilePath, sf.FileLine = fn.FileLine(pc[0])
	if r := strings.Split(sf.FuncName, "."); len(r) > 0 {
		sf.FuncName = r[len(r)-1]
	}

	return sf
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

// The textMessage creates a text message.
func textMessage(
	p string,
	l level.Level,
	t time.Time,
	o *Output,
	sf *stackFrame,
	f string,
	a ...any,
) string {
	// Generate log header.
	// The text before of the user's message, which includes the
	// prefix, the date and time of the event, the message level,
	// and additional format data (file, function, line etc.).
	sb := strings.Builder{}

	// Logger prefix.
	if p != "" {
		sb.WriteString(p)
		sb.WriteString(o.Space)
	}

	// Timestamp.
	sb.WriteString(t.Format(o.TimestampFormat))
	sb.WriteString(o.Space)

	// Level name.
	labels := level.Labels
	if o.WithColor.IsTrue() && runtime.GOOS != "windows" {
		labels = level.ColorLabels
	}

	if v, ok := labels[l]; ok {
		sb.WriteString(fmt.Sprintf(o.LevelFormat, v))
		sb.WriteString(o.Space)
	}

	// File path.
	// The FullPath takes precedence over ShortPath.
	if o.Layouts.FilePath() {
		if o.Layouts.FullFilePath() {
			sb.WriteString(sf.FilePath)
		} else {
			sb.WriteString(cutFilePath(shortPathSections, sf.FilePath))
		}

		if o.Layouts.LineNumber() {
			sb.WriteString(fmt.Sprintf(":%d", sf.FileLine))
		}

		sb.WriteString(o.Space)
	}

	// Line number.
	if o.Layouts.LineNumber() && !o.Layouts.FilePath() {
		sb.WriteString(fmt.Sprintf("%d%s", sf.FileLine, o.Space))
	}

	// Function name.
	if o.Layouts.FuncName() {
		sb.WriteString(sf.FuncName)
		if o.Layouts.FuncAddress() {
			sb.WriteString(fmt.Sprintf(":%#x", sf.FuncAddress))
		}
		sb.WriteString(o.Space)
	}

	// Function address.
	if o.Layouts.FuncAddress() && !o.Layouts.FuncName() {
		sb.WriteString(fmt.Sprintf("%#x%s", sf.FuncAddress, o.Space))
	}

	// Add message formatting.
	var msg string

	switch {
	case f == "":
		fallthrough
	case f == formatPrint:
		// For messages that are output on the same line, the task of
		// separating the messages falls on the user. We don't need to
		// add extra characters to user messages.
		// msg = fmt.Sprintf("%s%s%s", sb.String(), fmt.Sprint(a...), o.Space)
		msg = fmt.Sprintf("%s%s", sb.String(), fmt.Sprint(a...))
	case f == formatPrintln:
		msg = fmt.Sprintf("%s%s", sb.String(), fmt.Sprintln(a...))
	default:
		msg = fmt.Sprintf("%s%s", sb.String(), fmt.Sprintf(f, a...))
	}

	return msg
}

// The objectMessage creates a JSON message.
func objectMessage(
	p string,
	l level.Level,
	t time.Time,
	o *Output,
	sf *stackFrame,
	f string,
	a ...any,
) string {
	// Output object.
	// A general structure for outputting a log in JSON format.
	obj := struct {
		Prefix      string `json:"prefix,omitempty"`
		Level       string `json:"level,omitempty"`
		Timestamp   string `json:"timestamp,omitempty"`
		Message     string `json:"message,omitempty"`
		FilePath    string `json:"filePath,omitempty"`
		LineNumber  int    `json:"lineNumber,omitempty"`
		FuncName    string `json:"funcName,omitempty"`
		FuncAddress string `json:"funcAddress,omitempty"`
	}{}

	// Logger prefix.
	if p != "" {
		obj.Prefix = p
	}

	// Timestamp.
	obj.Timestamp = t.Format(o.TimestampFormat)

	// Level label.
	if v, ok := level.Labels[l]; ok {
		obj.Level = v
	}

	// File path, full path only.
	if o.Layouts.FilePath() {
		obj.FilePath = sf.FilePath
	}

	// Function name.
	if o.Layouts.FuncName() {
		obj.FuncName = sf.FuncName
	}

	// Function address.
	if o.Layouts.FuncAddress() {
		obj.FuncAddress = fmt.Sprintf("%#x", sf.FuncAddress)
	}

	// Line number.
	if o.Layouts.LineNumber() {
		obj.LineNumber = sf.FileLine
	}

	// Clean message for default -ln format.
	// Add message formatting.
	switch {
	case f == "":
		fallthrough
	case f == formatPrint:
		obj.Message = fmt.Sprint(a...)
	case f == formatPrintln:
		obj.Message = strings.TrimSuffix(fmt.Sprintln(a...), "\n")
	default:
		obj.Message = fmt.Sprintf(f, a...)
	}

	// Marshal object to JSON.
	data, err := json.Marshal(obj)
	data = g.If(err != nil, []byte{}, data)

	// Add JSON formatting.
	var msg string
	switch {
	case f == "":
		fallthrough
	case f == formatPrint:
		msg = string(data)
	default: // for formatStrLn and others
		msg = fmt.Sprintf("%s\n", data)
	}

	// Add space if necessary.
	if o.Space != "" {
		msg += o.Space
	}

	return msg
}

/*
// The getWriterID returns the unique ID of the object
// in the io.Writer interface.
//
// To identify duplicate objects of the os.Writer interface, the actual
// address of the object in memory is used. This does not guarantee 100%
// verification of uniqueness, but there is no need for it.
//
// Usually, these are 2-3 files for logging, which are added almost
// simultaneously, which guarantees a stable address of the object at
// the time of adding it to the list of outputs. And the issue of
// duplicates must be monitored by the developer who uses the logger.
//
// Therefore, checking for duplicates is a useful auxiliary function
// for detecting "stupid" logger creation errors.
func getWriterID(w io.Writer) uintptr {
	switch v := w.(type) {
	case *os.File:
		return reflect.ValueOf(v).Pointer()
	case *bytes.Buffer:
		return reflect.ValueOf(v).Pointer()
	case *strings.Builder:
		return reflect.ValueOf(v).Pointer()
	case *bufio.Writer:
		return reflect.ValueOf(v).Pointer()
	case *gzip.Writer:
		return reflect.ValueOf(v).Pointer()
	case *io.PipeWriter:
		return reflect.ValueOf(v).Pointer()
	}

	// Unknown type.
	// Get the address of the interface itself.
	return reflect.ValueOf(w).Pointer()
}
*/
