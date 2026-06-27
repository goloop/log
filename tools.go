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

// logPackage is the import path of this package. Frames whose function
// belongs to it are part of the logger's own call chain and are skipped
// when locating the caller's stack frame.
const logPackage = "github.com/goloop/log/v2"

// The pcPool recycles the program-counter buffers used while locating the
// call site, keeping captureFrame off the allocation path.
var pcPool = sync.Pool{New: func() any {
	pcs := make([]uintptr, 64)
	return &pcs
}}

// The captureFrame returns the first stack frame outside this package — the
// code that called into the logger — skipping skip additional frames for
// user wrapper layers. The boolean result is false when no such frame can be
// found. Unlike a fixed skip count it is immune to changes in the depth of
// the logger's internal call chain.
func captureFrame(skip int) (*stackFrame, bool) {
	pcsp := pcPool.Get().(*[]uintptr)
	defer pcPool.Put(pcsp)
	pcs := *pcsp

	n := runtime.Callers(2, pcs) // skip runtime.Callers and captureFrame
	if n == 0 {
		return &stackFrame{}, false
	}

	frames := runtime.CallersFrames(pcs[:n])
	insideLog := true
	for {
		frame, more := frames.Next()
		if frame.Function == "" && frame.PC == 0 {
			break
		}

		// Skip the logger's own frames (echo, emit, the level methods and
		// the package-level wrappers) until the first caller frame.
		if insideLog {
			if inLogPackage(frame.Function) {
				if !more {
					break
				}
				continue
			}
			insideLog = false
		}

		// Then skip the requested number of user wrapper frames.
		if skip > 0 {
			skip--
			if !more {
				break
			}
			continue
		}

		return frameOf(frame), true
	}

	return &stackFrame{}, false
}

// The inLogPackage reports whether fn is a function declared in this package
// (e.g. "github.com/goloop/log/v2.(*Logger).Info"), as opposed to a caller,
// a subpackage (".../v2/level.…") or the external test package
// (".../v2_test.…").
func inLogPackage(fn string) bool {
	return len(fn) > len(logPackage) &&
		fn[:len(logPackage)] == logPackage &&
		fn[len(logPackage)] == '.'
}

// The frameOf converts a runtime.Frame into a stackFrame with a short
// function name.
func frameOf(f runtime.Frame) *stackFrame {
	sf := &stackFrame{
		FileLine:    f.Line,
		FuncName:    f.Function,
		FuncAddress: f.Entry,
		FilePath:    f.File,
	}
	if i := strings.LastIndexByte(sf.FuncName, '.'); i >= 0 {
		sf.FuncName = sf.FuncName[i+1:]
	}

	return sf
}

// The frameFromPC builds a stackFrame from a single program counter. The
// boolean result is false when pc does not map to a known function, in
// which case the returned frame is empty.
func frameFromPC(pc uintptr) (*stackFrame, bool) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return &stackFrame{}, false
	}

	sf := &stackFrame{FuncName: fn.Name(), FuncAddress: fn.Entry()}
	sf.FilePath, sf.FileLine = fn.FileLine(pc)
	if i := strings.LastIndexByte(sf.FuncName, '.'); i >= 0 {
		sf.FuncName = sf.FuncName[i+1:]
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
	fields []logField,
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

	// Read the layout bits once instead of decoding the mask repeatedly.
	hasFile := o.Layouts.FilePath()
	hasLine := o.Layouts.LineNumber()
	hasFunc := o.Layouts.FuncName()
	hasAddr := o.Layouts.FuncAddress()

	// File path.
	// The FullPath takes precedence over ShortPath.
	if hasFile {
		if o.Layouts.FullFilePath() {
			buf.WriteString(sf.FilePath)
		} else {
			buf.WriteString(cutFilePath(shortPathSections, sf.FilePath))
		}

		if hasLine {
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(sf.FileLine))
		}

		buf.WriteString(o.Space)
	}

	// Line number.
	if hasLine && !hasFile {
		buf.WriteString(strconv.Itoa(sf.FileLine))
		buf.WriteString(o.Space)
	}

	// Function name.
	if hasFunc {
		buf.WriteString(sf.FuncName)
		if hasAddr {
			buf.WriteByte(':')
			writeHex(buf, sf.FuncAddress)
		}
		buf.WriteString(o.Space)
	}

	// Function address.
	if hasAddr && !hasFunc {
		writeHex(buf, sf.FuncAddress)
		buf.WriteString(o.Space)
	}

	// Append the pre-rendered message body. The body already encodes the
	// operand separators and any trailing newline (for the println kind),
	// so the header is simply prefixed to it.
	buf.WriteString(body)

	// Structured fields (from the slog bridge) as space-separated key=value.
	for i := range fields {
		buf.WriteByte(' ')
		buf.WriteString(fields[i].key)
		buf.WriteByte('=')
		writeValue(buf, fields[i].val)
	}
}

// The writeValue writes v in its textual form (mirroring fmt's default
// formatting) for a key=value structured field in text output.
func writeValue(buf *bytes.Buffer, v any) {
	switch x := v.(type) {
	case string:
		buf.WriteString(x)
	case nil:
		buf.WriteString("<nil>")
	default:
		fmt.Fprint(buf, x)
	}
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
	fields []logField,
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
	} else if len(fields) == 0 {
		buf.Write(data)
	} else {
		writeObjectWithFields(buf, data, fields)
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

// The writeObjectWithFields writes the marshalled object with the structured
// fields spliced in as typed JSON keys before the closing brace.
func writeObjectWithFields(buf *bytes.Buffer, data []byte, fields []logField) {
	if len(data) < 2 || data[len(data)-1] != '}' {
		buf.Write(data) // unexpected shape; emit as-is
		return
	}

	buf.Write(data[:len(data)-1]) // everything but the closing brace
	comma := len(data) > 2        // object already carries at least one field
	for i := range fields {
		if comma {
			buf.WriteByte(',')
		}
		comma = true

		writeJSONString(buf, fields[i].key)
		buf.WriteByte(':')
		writeJSONValue(buf, fields[i].val)
	}
	buf.WriteByte('}')
}

// The writeJSONValue encodes v as a JSON value into buf. Common scalar types
// are written directly with strconv (no allocation); strings are escaped;
// any other type falls back to json.Marshal for correctness.
func writeJSONValue(buf *bytes.Buffer, v any) {
	switch x := v.(type) {
	case nil:
		buf.WriteString("null")
	case string:
		writeJSONString(buf, x)
	case bool:
		if x {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case int:
		buf.WriteString(strconv.FormatInt(int64(x), 10))
	case int64:
		buf.WriteString(strconv.FormatInt(x, 10))
	case time.Duration:
		buf.WriteString(strconv.FormatInt(int64(x), 10))
	case uint64:
		buf.WriteString(strconv.FormatUint(x, 10))
	case float64:
		buf.WriteString(strconv.FormatFloat(x, 'g', -1, 64))
	default:
		if vb, err := json.Marshal(v); err == nil {
			buf.Write(vb)
		} else {
			buf.WriteString("null")
		}
	}
}

// The writeJSONString writes s as a quoted, escaped JSON string into buf. The
// common no-escape case writes the whole string at once without allocating.
// Unlike encoding/json it does not HTML-escape '<', '>' and '&'; the result
// is still valid JSON.
func writeJSONString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' {
			continue
		}
		if start < i {
			buf.WriteString(s[start:i])
		}
		switch c {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			const hex = "0123456789abcdef"
			buf.WriteString(`\u00`)
			buf.WriteByte(hex[c>>4])
			buf.WriteByte(hex[c&0xF])
		}
		start = i + 1
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
}
