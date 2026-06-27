package log

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"testing"

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
		Writer: &nopWriter{}, // io.Discard is rejected by g.IsEmpty (empty struct)
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
