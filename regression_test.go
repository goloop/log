package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/goloop/log/v2/layout"
	"github.com/goloop/log/v2/level"
	"github.com/goloop/trit"
)

// TestConcurrentFxxx guards against the data race once present in echo: many
// goroutines calling an F-method concurrently must not race or panic.
func TestConcurrentFxxx(t *testing.T) {
	logger := New()
	if err := logger.SetOutputs(Output{
		Name:   "t",
		Writer: io.Discard,
		Levels: level.Default,
	}); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				logger.Ferror(io.Discard, "msg")
			}
		}()
	}
	wg.Wait()
}

// TestFerrorIsolation checks that an F-method writes to both the ad-hoc
// writer and the configured outputs, and never leaks an output entry.
func TestFerrorIsolation(t *testing.T) {
	var cfg, adhoc bytes.Buffer
	logger := New()
	if err := logger.SetOutputs(Output{
		Name:   "cfg",
		Writer: &cfg,
		Levels: level.Error,
	}); err != nil {
		t.Fatal(err)
	}

	before := len(logger.Outputs())
	logger.Ferror(&adhoc, "HELLO")

	if !strings.Contains(adhoc.String(), "HELLO") {
		t.Errorf("ad-hoc writer missing message: %q", adhoc.String())
	}
	if !strings.Contains(cfg.String(), "HELLO") {
		t.Errorf("configured output missing message: %q", cfg.String())
	}
	if after := len(logger.Outputs()); after != before {
		t.Errorf("output leak: %d -> %d", before, after)
	}
}

// TestJSONKeys pins the JSON field names emitted by the logger.
func TestJSONKeys(t *testing.T) {
	var buf bytes.Buffer
	logger := New()
	if err := logger.SetOutputs(Output{
		Name:      "t",
		Writer:    &buf,
		Levels:    level.Default,
		TextStyle: trit.False,
		Layouts:   layout.FullFilePath | layout.LineNumber | layout.FuncName,
	}); err != nil {
		t.Fatal(err)
	}

	logger.Info("x")

	var m map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &m); err != nil {
		t.Fatalf("invalid JSON %q: %v", buf.String(), err)
	}

	for _, k := range []string{"filePath", "lineNumber", "funcName"} {
		if _, ok := m[k]; !ok {
			t.Errorf("missing expected key %q in %v", k, m)
		}
	}
	for _, k := range []string{"file", "line"} {
		if _, ok := m[k]; ok {
			t.Errorf("unexpected legacy key %q", k)
		}
	}
}

func FuzzCutFilePath(f *testing.F) {
	for _, s := range []string{"", "/", "main.go", "/a/b/c/main.go", "////"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, path string) {
		_ = cutFilePath(shortPathSections, path) // must not panic
	})
}

func FuzzNew(f *testing.F) {
	for _, s := range []string{"", "APP", " ", "ПРЕФІКС"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, prefix string) {
		_ = New(prefix) // must not panic
	})
}

// TestLevelMaskUpdatesOnEdit guards the lock-free level fast-gate: it must
// track configuration changes made through EditOutputs.
func TestLevelMaskUpdatesOnEdit(t *testing.T) {
	var buf bytes.Buffer
	logger := New()
	if err := logger.SetOutputs(Output{
		Name:   "t",
		Writer: &buf,
		Levels: level.Info,
	}); err != nil {
		t.Fatal(err)
	}

	// Debug is not in any output's mask → gated out.
	logger.Debug("x")
	if buf.Len() != 0 {
		t.Errorf("Debug should be gated, got: %q", buf.String())
	}

	// Enable Debug via EditOutputs → mask refreshed → now emitted.
	if err := logger.EditOutputs(Output{
		Name:   "t",
		Levels: level.Info | level.Debug,
	}); err != nil {
		t.Fatal(err)
	}
	logger.Debug("hello")
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("Debug should emit after EditOutputs, got: %q", buf.String())
	}

	// Disable the output → mask empties → gated out again.
	buf.Reset()
	if err := logger.EditOutputs(Output{
		Name:    "t",
		Enabled: trit.False,
	}); err != nil {
		t.Fatal(err)
	}
	logger.Info("nope")
	if buf.Len() != 0 {
		t.Errorf("disabled output should emit nothing, got: %q", buf.String())
	}
}

// TestSetDefault checks that the package-level default logger can be swapped.
func TestSetDefault(t *testing.T) {
	orig := Log()
	defer SetDefault(orig)

	custom := New("CUSTOM")
	SetDefault(custom)
	if Log() != custom {
		t.Error("SetDefault did not install the custom logger")
	}

	SetDefault(nil) // must be ignored
	if Log() != custom {
		t.Error("SetDefault(nil) should be ignored")
	}
}

// TestAcceptsIoDiscard guards BUG-06: io.Discard is a valid writer.
func TestAcceptsIoDiscard(t *testing.T) {
	logger := New()
	if err := logger.SetOutputs(Output{
		Name:   "t",
		Writer: io.Discard,
		Levels: level.Default,
	}); err != nil {
		t.Errorf("io.Discard should be accepted, got: %v", err)
	}
}

// TestRejectsTypedNilWriter ensures a typed-nil writer is still rejected.
func TestRejectsTypedNilWriter(t *testing.T) {
	logger := New()
	var w *bytes.Buffer // typed nil
	if err := logger.SetOutputs(Output{
		Name:   "t",
		Writer: w,
		Levels: level.Default,
	}); err == nil {
		t.Error("typed-nil writer should be rejected")
	}
}

type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("disk full")
}

// TestErrorHandler guards BUG-03: write errors reach the handler.
func TestErrorHandler(t *testing.T) {
	logger := New()
	var got error
	logger.SetErrorHandler(func(o Output, n int, err error) { got = err })
	if err := logger.SetOutputs(Output{
		Name:   "t",
		Writer: errWriter{},
		Levels: level.Default,
	}); err != nil {
		t.Fatal(err)
	}

	logger.Info("x")
	if got == nil {
		t.Error("error handler was not called on a failing write")
	}
}

type reentrantWriter struct{ lg *Logger }

func (w *reentrantWriter) Write(p []byte) (int, error) {
	w.lg.SetPrefix("x") // acquires the write lock; must not deadlock
	return len(p), nil
}

// TestReentrantWriterNoDeadlock guards BUG-04: the logger lock is released
// before user writes, so a writer may call back into the logger.
func TestReentrantWriterNoDeadlock(t *testing.T) {
	logger := New()
	w := &reentrantWriter{lg: logger}
	if err := logger.SetOutputs(Output{
		Name:   "t",
		Writer: w,
		Levels: level.Default,
	}); err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		logger.Info("x")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock: Info did not return (lock held during Write)")
	}
}
