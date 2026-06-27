package log

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/goloop/log/v2/layout"
	"github.com/goloop/log/v2/level"
	"github.com/goloop/trit"
)

// newTextSlog returns a slog.Logger that writes text into a buffer through
// a goloop logger with the given layout.
func newTextSlog(t *testing.T, layouts layout.Layout) (*slog.Logger, *bytes.Buffer) {
	t.Helper()
	buf := &bytes.Buffer{}
	logger := New()
	err := logger.SetOutputs(Output{
		Name:      "test",
		Writer:    buf,
		Levels:    level.Default,
		Layouts:   layouts,
		TextStyle: trit.True,
	})
	if err != nil {
		t.Fatal(err)
	}
	return slog.New(logger.Handler()), buf
}

func TestSlogBasic(t *testing.T) {
	sl, buf := newTextSlog(t, 0)
	sl.Info("hello", "user", "bob", "n", 3)

	out := buf.String()
	for _, want := range []string{"INFO", "hello", "user=bob", "n=3"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q does not contain %q", out, want)
		}
	}
}

func TestSlogLevelMapping(t *testing.T) {
	cases := []struct {
		log   func(*slog.Logger)
		label string
	}{
		{func(l *slog.Logger) { l.Debug("m") }, "DEBUG"},
		{func(l *slog.Logger) { l.Info("m") }, "INFO"},
		{func(l *slog.Logger) { l.Warn("m") }, "WARNING"},
		{func(l *slog.Logger) { l.Error("m") }, "ERROR"},
	}

	for _, c := range cases {
		sl, buf := newTextSlog(t, 0)
		c.log(sl)
		if !strings.Contains(buf.String(), c.label) {
			t.Errorf("expected level %q in %q", c.label, buf.String())
		}
	}
}

func TestSlogWithAttrsAndGroup(t *testing.T) {
	sl, buf := newTextSlog(t, 0)
	sl = sl.With("svc", "api").WithGroup("req")
	sl.Info("done", "id", 7)

	out := buf.String()
	for _, want := range []string{"done", "svc=api", "req.id=7"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q does not contain %q", out, want)
		}
	}
}

func TestSlogEnabled(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New()
	if err := logger.SetOutputs(Output{
		Name:   "test",
		Writer: buf,
		Levels: level.Error, // only Error is emitted
	}); err != nil {
		t.Fatal(err)
	}
	sl := slog.New(logger.Handler())

	if sl.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("Debug should be disabled")
	}
	if !sl.Enabled(context.Background(), slog.LevelError) {
		t.Error("Error should be enabled")
	}
}

// TestSlogSourceFrame verifies that the record's program counter is used as
// the call-site frame, so file/line layouts point at the caller rather than
// at the bridge internals.
func TestSlogSourceFrame(t *testing.T) {
	sl, buf := newTextSlog(t, layout.ShortFilePath)
	sl.Info("x")

	if !strings.Contains(buf.String(), "slog_test.go") {
		t.Errorf("expected source file in %q", buf.String())
	}
}
