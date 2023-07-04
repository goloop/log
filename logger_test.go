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

//
// The others of the method is rolled through global function, see log_test.go
//
