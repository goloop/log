package log_test

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/goloop/log/v2"
	"github.com/goloop/log/v2/layout"
	"github.com/goloop/log/v2/level"
	"github.com/goloop/trit"
)

func callSiteLogger() (*log.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	lg := log.New()
	lg.SetOutputs(log.Output{
		Name:      "t",
		Writer:    buf,
		Levels:    level.Default,
		Layouts:   layout.ShortFilePath | layout.FuncName | layout.LineNumber,
		TextStyle: trit.True,
	})
	return lg, buf
}

// TestCallSiteMethod verifies that a layout-enabled logger reports the
// caller's file and exact line, not the logger's own internals.
func TestCallSiteMethod(t *testing.T) {
	lg, buf := callSiteLogger()

	_, _, line, _ := runtime.Caller(0)
	lg.Info("x") // must be reported as the call site (line+1)

	out := buf.String()
	want := fmt.Sprintf("callsite_test.go:%d", line+1)
	if !strings.Contains(out, want) {
		t.Errorf("expected %q in output, got: %s", want, out)
	}
	if strings.Contains(out, "logger.go") {
		t.Errorf("frame leaked logger internals: %s", out)
	}
}

// TestCallSitePackageLevel verifies the same for the package-level wrapper,
// which adds an extra internal frame.
func TestCallSitePackageLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := log.SetOutputs(log.Output{
		Name:    "t",
		Writer:  buf,
		Levels:  level.Default,
		Layouts: layout.ShortFilePath | layout.LineNumber,
	}); err != nil {
		t.Fatal(err)
	}
	defer log.SetOutputs(log.Stdout, log.Stderr) // restore global logger

	log.SetSkipStackFrames(0) // independent of any prior test's mutation

	_, _, line, _ := runtime.Caller(0)
	log.Info("x") // must be reported as the call site (line+1)

	out := buf.String()
	want := fmt.Sprintf("callsite_test.go:%d", line+1)
	if !strings.Contains(out, want) {
		t.Errorf("expected %q in output, got: %s", want, out)
	}
}

// TestSetSkipStackFramesAccepts pins that SetSkipStackFrames stores the given
// value instead of clamping it.
func TestSetSkipStackFramesAccepts(t *testing.T) {
	lg := log.New()
	if got := lg.SetSkipStackFrames(3); got != 3 {
		t.Errorf("SetSkipStackFrames(3) = %d, want 3", got)
	}
	if got := lg.SkipStackFrames(); got != 3 {
		t.Errorf("SkipStackFrames() = %d, want 3", got)
	}
}
