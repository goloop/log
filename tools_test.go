package log

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/goloop/log/layout"
	"github.com/goloop/log/level"
	"github.com/goloop/trit"
)

// TetsIoCopy tests ioCopy function.
func TestIoCopy(t *testing.T) {
	input := "Hello, World!"
	r := strings.NewReader(input)
	c := make(chan string)

	go ioCopy(r, c)
	result := <-c

	if result != input {
		t.Errorf("ioCopy failed, expected %v, got %v", input, result)
	}
}

// TestGetStackFrame tests getStackFrame function.
func TestGetStackFrame(t *testing.T) {
	frame := getStackFrame(2) // cuerrent function is TestGetStackFrame

	if frame == nil {
		t.Fatal("Expected frame to not be nil")
	}

	if frame.FuncName == "" {
		t.Errorf("Expected FuncName to not be empty")
	}

	if frame.FilePath == "" {
		t.Errorf("Expected FilePath to not be empty")
	}

	if frame.FileLine == 0 {
		t.Errorf("Expected FileLine to not be zero")
	}

	if frame.FuncAddress == 0 {
		t.Errorf("Expected FuncAddress to not be zero")
	}

	// Verify the function name
	expectedFuncName := "TestGetStackFrame"
	if frame.FuncName != expectedFuncName {
		t.Errorf("Expected function name to be '%s', got '%s'",
			expectedFuncName, frame.FuncName)
	}

	// Verify the file path
	_, fileName, _, _ := runtime.Caller(0)
	if !strings.Contains(frame.FilePath, fileName) {
		t.Errorf("Expected file path '%s' to contain '%s'",
			frame.FilePath, fileName)
	}
}

// TestGetStackFramePanicsOnNegativeSkip tests getStackFrame for panic.
func TestGetStackFramePanicsOnLargeSkip(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic")
		}
	}()

	getStackFrame(1024)
}

// TestCutFilePath tests cutFilePath function.
func TestCutFilePath(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		path     string
		expected string
	}{
		{
			name:     "Test case 1: Three sections, cut to two",
			n:        2,
			path:     "/path/to/file",
			expected: ".../to/file",
		},
		{
			name:     "Test case 2: Four sections, cut to two",
			n:        2,
			path:     "/path/to/another/file",
			expected: ".../another/file",
		},
		{
			name:     "Test case 3: One section, cut to two",
			n:        2,
			path:     "/file",
			expected: "/file",
		},
		{
			name:     "Test case 4: Three sections, cut to three",
			n:        3,
			path:     "/path/to/file",
			expected: "/path/to/file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := cutFilePath(test.n, test.path)

			if result != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

// TestTextMessage tests textMessage function.
func TestTextMessage(t *testing.T) {
	prefix := "test"
	level := level.Info
	timestamp := time.Now()
	output := &Stdout
	output.WithColor = trit.True
	output.Layouts = output.Layouts | layout.LineNumber | layout.FuncAddress
	stackframe := getStackFrame(2)

	tests := []struct {
		name string
		f    string
		a    []any
		e    string
	}{
		{
			name: "Text message with formatted string",
			f:    "formatted string %s",
			a:    []any{"value"},
			e:    "formatted string value",
		},
		{
			name: "Text message with multiple formatted values",
			f:    "formatted string with multiple values %s %d",
			a:    []any{"value", 1},
			e:    "formatted string with multiple values value 1",
		},
		{
			name: "Text message with no formatting",
			f:    "",
			a:    []any{"value"},
			e:    "value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := textMessage(
				prefix,
				level,
				timestamp,
				output,
				stackframe,
				test.f,
				test.a...,
			)

			if !strings.Contains(result, test.e) {
				t.Errorf("Message '%s' doesn't contains '%s'", result, test.e)
			}
		})
	}

	// Change layouts.
	output.Layouts = layout.FullFilePath
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := textMessage(
				prefix,
				level,
				timestamp,
				output,
				stackframe,
				test.f,
				test.a...,
			)

			if !strings.Contains(result, test.e) {
				t.Errorf("Message '%s' doesn't contains '%s'", result, test.e)
			}
		})
	}

	// Change layouts.
	output.Layouts = layout.LineNumber | layout.FuncAddress
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := textMessage(
				prefix,
				level,
				timestamp,
				output,
				stackframe,
				test.f,
				test.a...,
			)

			if !strings.Contains(result, test.e) {
				t.Errorf("Message '%s' doesn't contains '%s'", result, test.e)
			}
		})
	}
}

// TestObjectMessage tests objectMessage function.
func TestObjectMessage(t *testing.T) {
	prefix := "test"
	level := level.Info
	timestamp := time.Now()
	output := &Stdout
	stackframe := getStackFrame(2)

	tests := []struct {
		name string
		f    string
		a    []any
		e    map[string]interface{}
	}{
		{
			name: "Object message with formatted string",
			f:    "formatted string %s",
			a:    []any{"value"},
			e: map[string]interface{}{
				"prefix":      prefix,
				"level":       "INFO",
				"timestamp":   timestamp.Format(output.TimestampFormat),
				"message":     "formatted string value",
				"filePath":    stackframe.FilePath,
				"funcName":    stackframe.FuncName,
				"funcAddress": fmt.Sprintf("%#x", stackframe.FuncAddress),
				"lineNumber":  stackframe.FileLine,
			},
		},
		{
			name: "Object message with multiple formatted values",
			f:    "formatted string with multiple values %s %d",
			a:    []any{"value", 1},
			e: map[string]interface{}{
				"prefix":      prefix,
				"level":       "INFO",
				"timestamp":   timestamp.Format(output.TimestampFormat),
				"message":     "formatted string with multiple values value 1",
				"filePath":    stackframe.FilePath,
				"funcName":    stackframe.FuncName,
				"funcAddress": fmt.Sprintf("%#x", stackframe.FuncAddress),
				"lineNumber":  stackframe.FileLine,
			},
		},
		{
			name: "Object message with no formatting",
			f:    "",
			a:    []any{"value"},
			e: map[string]interface{}{
				"prefix":      prefix,
				"level":       "INFO",
				"timestamp":   timestamp.Format(output.TimestampFormat),
				"message":     "value",
				"filePath":    stackframe.FilePath,
				"funcName":    stackframe.FuncName,
				"funcAddress": fmt.Sprintf("%#x", stackframe.FuncAddress),
				"lineNumber":  stackframe.FileLine,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := objectMessage(
				prefix,
				level,
				timestamp,
				output,
				stackframe,
				test.f,
				test.a...,
			)

			// Unmarshal the JSON result into a map
			var resultObj map[string]interface{}
			err := json.Unmarshal([]byte(result), &resultObj)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if resultObj["level"] != test.e["level"] ||
				resultObj["prefix"] != test.e["prefix"] ||
				resultObj["message"] != test.e["message"] {
				t.Errorf("Expected '%v', got '%v'", test.e, resultObj)
			}
		})
	}
}

/*
// TestGetWriterID tests getWriterID function.
func TestGetWriterID(t *testing.T) {
	// Create several types that satisfy the io.Writer interface
	file, err := os.Create("test.txt")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer os.Remove("test.txt")
	defer file.Close()

	buffer := &bytes.Buffer{}
	builder := &strings.Builder{}
	writer := bufio.NewWriter(buffer)
	gzipWriter := gzip.NewWriter(buffer)
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()

	tests := []struct {
		name  string
		input io.Writer
	}{
		{"os.File", file},
		{"bytes.Buffer", buffer},
		{"strings.Builder", builder},
		{"bufio.Writer", writer},
		{"gzip.Writer", gzipWriter},
		{"io.PipeWriter", pipeWriter},
	}

	// Map to store the IDs of the writers
	writerIDs := make(map[uintptr]bool)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			id := getWriterID(test.input)
			if _, exists := writerIDs[id]; exists {
				t.Errorf("Non-unique writer ID returned for type: %s",
					test.name)
			}
			writerIDs[id] = true
		})
	}
}
*/
