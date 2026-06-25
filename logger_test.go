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

	logger.echo(nil, level.Debug, kindPrint, "test message")
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

	logger.echo(nil, level.Debug, kindPrint, "test message")
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

	logger.echo(nil, level.Debug, kindPrint, "test message")
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

	logger.echo(nil, level.Debug, kindPrint, "test message")
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
		name string
		kind emitKind
		body string
		want string
	}{
		{
			name: "Print kind",
			kind: kindPrint,
			body: "helloworld",
			want: "helloworld", // fmt.Sprint body
		},
		{
			name: "Println kind",
			kind: kindPrintln,
			body: "hello world\n",
			want: " hello world\n", // fmt.Sprintln body (header adds leading space)
		},
		{
			name: "Printf kind",
			kind: kindPrintf,
			body: "[777]-message is true",
			want: "[777]-message is true", // fmt.Sprintf body
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

			logger.echo(nil, level.Debug, tt.kind, tt.body)
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
		name string
		kind emitKind
		body string
		want string
	}{
		{
			name: "Print kind",
			kind: kindPrint,
			body: "helloworld",
			want: "helloworld", // fmt.Sprint body
		},
		{
			name: "Println kind",
			kind: kindPrintln,
			body: "hello world\n",
			want: "hello world", // fmt.Sprintln body, trailing newline trimmed
		},
		{
			name: "Printf kind",
			kind: kindPrintf,
			body: "[777]-message is true",
			want: "[777]-message is true", // fmt.Sprintf body
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

			logger.echo(nil, level.Debug, tt.kind, tt.body)
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
