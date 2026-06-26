package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/goloop/log/v2/level"
)

// The bufPool recycles the byte buffers used to assemble a single log
// message, keeping the hot path allocation-free across calls.
var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

// The emptyFrame is a shared, read-only zero frame used when no output
// requests stack-frame information, so the builders never see a nil sf.
var emptyFrame = &stackFrame{}

// The stackFrame contains the top-level trace information
// where the logging method was called.
type stackFrame struct {
	FileLine    int     // line number
	FuncName    string  // function name
	FuncAddress uintptr // address of the function
	FilePath    string  // file path
}

// The getStackFrame returns the top stack frame after skipping skip frames.
// The boolean result is false when no valid frame can be captured (for
// example, when skip is larger than the call stack); in that case the
// returned frame is empty rather than causing a panic.
func getStackFrame(skip int) (*stackFrame, bool) {
	sf := &stackFrame{}

	// Return the single program counter at the requested depth. Only pc[0]
	// is ever read, so a one-element buffer is enough regardless of skip.
	pc := make([]uintptr, 1)
	if n := runtime.Callers(skip, pc); n == 0 {
		return sf, false
	}

	// Get a function at an address on the stack. A nil result means the
	// program counter does not map to a known function (skip too large).
	fn := runtime.FuncForPC(pc[0])
	if fn == nil {
		return sf, false
	}

	// Get name, path and line of the file.
	sf.FuncName = fn.Name()
	sf.FuncAddress = fn.Entry()
	sf.FilePath, sf.FileLine = fn.FileLine(pc[0])
	if r := strings.Split(sf.FuncName, "."); len(r) > 0 {
		sf.FuncName = r[len(r)-1]
	}

	return sf, true
}

// The cutFilePath keeps the last n sections of the path. It walks the
// string from the end counting separators, so no intermediate slice is
// allocated. The path is returned unchanged when it has n or fewer
// sections (i.e. n or fewer separators).
func cutFilePath(n int, path string) string {
	count, cut := 0, -1
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] != '/' {
			continue
		}

		count++
		if count == n {
			cut = i // start of the last n sections
		} else if count > n {
			return ".../" + path[cut+1:]
		}
	}

	return path
}

// The appendText writes a text-style message into buf: the header
// (prefix, timestamp, level, and the requested file/function/line layout)
// followed by the pre-rendered body.
func appendText(
	buf *bytes.Buffer,
	p string,
	l level.Level,
	t time.Time,
	o *Output,
	sf *stackFrame,
	body string,
) {
	// Logger prefix.
	if p != "" {
		buf.WriteString(p)
		buf.WriteString(o.Space)
	}

	// Timestamp.
	buf.WriteString(t.Format(o.TimestampFormat))
	buf.WriteString(o.Space)

	// Level name.
	labels := level.Labels
	if o.WithColor.IsTrue() && runtime.GOOS != "windows" {
		labels = level.ColorLabels
	}

	if v, ok := labels[l]; ok {
		if o.LevelFormat == "%s" { // common case: no formatting needed
			buf.WriteString(v)
		} else {
			fmt.Fprintf(buf, o.LevelFormat, v)
		}
		buf.WriteString(o.Space)
	}

	// File path.
	// The FullPath takes precedence over ShortPath.
	if o.Layouts.FilePath() {
		if o.Layouts.FullFilePath() {
			buf.WriteString(sf.FilePath)
		} else {
			buf.WriteString(cutFilePath(shortPathSections, sf.FilePath))
		}

		if o.Layouts.LineNumber() {
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(sf.FileLine))
		}

		buf.WriteString(o.Space)
	}

	// Line number.
	if o.Layouts.LineNumber() && !o.Layouts.FilePath() {
		buf.WriteString(strconv.Itoa(sf.FileLine))
		buf.WriteString(o.Space)
	}

	// Function name.
	if o.Layouts.FuncName() {
		buf.WriteString(sf.FuncName)
		if o.Layouts.FuncAddress() {
			buf.WriteByte(':')
			writeHex(buf, sf.FuncAddress)
		}
		buf.WriteString(o.Space)
	}

	// Function address.
	if o.Layouts.FuncAddress() && !o.Layouts.FuncName() {
		writeHex(buf, sf.FuncAddress)
		buf.WriteString(o.Space)
	}

	// Append the pre-rendered message body. The body already encodes the
	// operand separators and any trailing newline (for the println kind),
	// so the header is simply prefixed to it.
	buf.WriteString(body)
}

// The writeHex writes v as a "0x"-prefixed lowercase hexadecimal number,
// matching the %#x verb without allocating an intermediate string.
func writeHex(buf *bytes.Buffer, v uintptr) {
	buf.WriteString("0x")
	buf.WriteString(strconv.FormatUint(uint64(v), 16))
}

// The appendObject writes a JSON-style message into buf. On a marshal
// failure it falls back to a minimal text line so the message is never
// silently lost.
func appendObject(
	buf *bytes.Buffer,
	p string,
	l level.Level,
	t time.Time,
	o *Output,
	sf *stackFrame,
	kind emitKind,
	body string,
) {
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

	// The pre-rendered body already carries the operand separators. For the
	// println kind it ends with a newline that the JSON message field does
	// not need, so it is trimmed.
	if kind == kindPrintln {
		obj.Message = strings.TrimSuffix(body, "\n")
	} else {
		obj.Message = body
	}

	// Marshal object to JSON straight into the buffer. If marshalling ever
	// fails, fall back to a minimal text line (timestamp, level, message)
	// rather than dropping the log entry.
	data, err := json.Marshal(obj)
	if err != nil {
		buf.WriteString(obj.Timestamp)
		buf.WriteByte(' ')
		if obj.Level != "" {
			buf.WriteString(obj.Level)
			buf.WriteByte(' ')
		}
		buf.WriteString(obj.Message)
	} else {
		buf.Write(data)
	}

	// Add JSON formatting. The print kind keeps JSON blocks on a single
	// line; println and printf terminate each block with a newline.
	if kind != kindPrint {
		buf.WriteByte('\n')
	}

	// Add space if necessary.
	if o.Space != "" {
		buf.WriteString(o.Space)
	}
}
