package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/goloop/log/v2/layout"
	"github.com/goloop/log/v2/level"
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
	stackframe := &stackFrame{
		FilePath:    "/tmp/example/main.go",
		FileLine:    42,
		FuncName:    "main",
		FuncAddress: 0x1234,
	}

	tests := []struct {
		name string
		body string
		e    string
	}{
		{
			name: "Text message with formatted string",
			body: "formatted string value",
			e:    "formatted string value",
		},
		{
			name: "Text message with multiple formatted values",
			body: "formatted string with multiple values value 1",
			e:    "formatted string with multiple values value 1",
		},
		{
			name: "Text message with no formatting",
			body: "value",
			e:    "value",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			appendText(&buf, prefix, level, timestamp, output, stackframe,
				test.body)
			result := buf.String()

			if !strings.Contains(result, test.e) {
				t.Errorf("Message '%s' doesn't contains '%s'", result, test.e)
			}
		})
	}

	// Change layouts.
	output.Layouts = layout.FullFilePath
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			appendText(&buf, prefix, level, timestamp, output, stackframe,
				test.body)
			result := buf.String()

			if !strings.Contains(result, test.e) {
				t.Errorf("Message '%s' doesn't contains '%s'", result, test.e)
			}
		})
	}

	// Change layouts.
	output.Layouts = layout.LineNumber | layout.FuncAddress
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			appendText(&buf, prefix, level, timestamp, output, stackframe,
				test.body)
			result := buf.String()

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
	stackframe := &stackFrame{
		FilePath:    "/tmp/example/main.go",
		FileLine:    42,
		FuncName:    "main",
		FuncAddress: 0x1234,
	}

	tests := []struct {
		name string
		kind emitKind
		body string
		e    map[string]interface{}
	}{
		{
			name: "Object message with formatted string",
			kind: kindPrintf,
			body: "formatted string value",
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
			kind: kindPrintf,
			body: "formatted string with multiple values value 1",
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
			kind: kindPrint,
			body: "value",
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
			var buf bytes.Buffer
			appendObject(&buf, prefix, level, timestamp, output, stackframe,
				test.kind, test.body)
			result := buf.String()

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

// The ioCopy reads everything from r and sends it as a single string on c.
// It is a helper used by the package tests to drain os.Pipe / io.Pipe ends.
func ioCopy(r io.Reader, c chan string) {
	var buf bytes.Buffer
	io.Copy(&buf, r)
	c <- buf.String()
}
