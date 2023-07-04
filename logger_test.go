package log

import (
	"os"
	"strings"
	"testing"

	"github.com/goloop/log/level"
	"github.com/goloop/trit"
)

// TestEcho tests the echo method of the Logger.
func TestEcho(t *testing.T) {
	// Create a new logger.
	logger := New("TEST-PREFIX:")

	// Classical test.
	r, w, _ := os.Pipe()
	err := logger.SetOutputs(Output{
		Name:   "test",
		Writer: w,
		Levels: level.Default,
	})
	if err != nil {
		t.Fatal(err)
	}

	logger.echo(nil, level.Debug, "test %s", "message")
	outC := make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out := <-outC

	if !strings.Contains(out, "test message") {
		t.Errorf("echo did not write the correct TEXT message: %s", out)
	}

	// As JSON.
	r, w, _ = os.Pipe()
	logger.SetOutputs(Output{
		Name:       "test",
		Writer:     w,
		Levels:     level.Default,
		WithPrefix: trit.False,
	})

	logger.echo(nil, level.Debug, "test %s", "message")
	outC = make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out = <-outC

	if strings.Contains(out, "TEST-PREFIX") {
		t.Errorf("the prefix should not appear in this test: %s", out)
	}

	// As JSON.
	r, w, _ = os.Pipe()
	logger.SetOutputs(Output{
		Name:      "test",
		Writer:    w,
		Levels:    level.Default,
		TextStyle: trit.False,
	})

	logger.echo(nil, level.Debug, "test %s", "message")
	outC = make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out = <-outC

	if !strings.Contains(out, "\"level\":\"DEBUG\"") {
		t.Errorf("echo did not write the correct JSON message: %s", out)
	}

	// Disabled.
	r, w, _ = os.Pipe()
	logger.SetOutputs(Output{
		Name:    "test",
		Writer:  w,
		Enabled: trit.False,
	})

	logger.echo(nil, level.Debug, "test %s", "message")
	outC = make(chan string)
	go ioCopy(r, outC)
	w.Close()
	out = <-outC

	if len(out) != 0 {
		t.Errorf("should not write anything: %s", out)
	}
}

// TestEchoWithTextFormatting tests the echo method with Text formatting.
func TestEchoWithTextFormatting(t *testing.T) {
	tests := []struct {
		name   string
		format string
		in     []interface{}
		want   string
	}{
		{
			name:   "Empty format",
			format: "",
			in:     []interface{}{"hello", "world"},
			want:   "helloworld", // used fmt.Print
		},
		{
			name:   "System formatStr",
			format: formatStr,
			in:     []interface{}{"hello", "world"},
			want:   "helloworld", // used fmt.Print
		},
		{
			name:   "System formatStrLn",
			format: formatStrLn,
			in:     []interface{}{"hello", "world"},
			want:   " hello world\n", // used fmt.Println
		},
		{
			name:   "Custom formats",
			format: "[%d]-%s is %v",
			in:     []interface{}{777, "message", true},
			want:   "[777]-message is true", // used fmt.Printf
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New()
			logger.SetSkipStackFrames(2)
			r, w, _ := os.Pipe()
			logger.SetOutputs(Output{
				Name:   "test",
				Writer: w,
				Levels: level.Default,
			})

			logger.echo(nil, level.Debug, tt.format, tt.in...)
			outC := make(chan string)
			go ioCopy(r, outC)
			w.Close()
			out := <-outC
			if !strings.Contains(out, tt.want) {
				t.Errorf("Expression `%v` does not contain `%v`", out, tt.want)
			}
		})
	}
}

// TestEchoWithJSONFormatting tests the echo method with JSON formatting.
func TestEchoWithJSONFormatting(t *testing.T) {
	tests := []struct {
		name   string
		format string
		in     []interface{}
		want   string
	}{
		{
			name:   "Empty format",
			format: "",
			in:     []interface{}{"hello", "world"},
			want:   "helloworld", // used fmt.Print
		},
		{
			name:   "System formatStr",
			format: formatStr,
			in:     []interface{}{"hello", "world"},
			want:   "helloworld", // used fmt.Print
		},
		{
			name:   "System formatStrLn",
			format: formatStrLn,
			in:     []interface{}{"hello", "world"},
			want:   "hello world", // used fmt.Println with Trim
		},
		{
			name:   "Custom formats",
			format: "[%d]-%s is %v",
			in:     []interface{}{777, "message", true},
			want:   "[777]-message is true", // used fmt.Printf
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New()
			logger.SetSkipStackFrames(2)
			r, w, _ := os.Pipe()
			logger.SetOutputs(Output{
				Name:      "test",
				Writer:    w,
				Levels:    level.Default,
				TextStyle: trit.False,
			})

			logger.echo(nil, level.Debug, tt.format, tt.in...)
			outC := make(chan string)
			go ioCopy(r, outC)
			w.Close()
			out := <-outC
			if !strings.Contains(out, tt.want) {
				t.Errorf("Expression `%v` does not contain `%v`", out, tt.want)
			}
		})
	}
}

//
// The others of the method is rolled through global function, see log_test.go
//
