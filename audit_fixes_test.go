package log

import (
	"bytes"
	"encoding/json"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/goloop/log/v2/layout"
	"github.com/goloop/log/v2/level"
	"github.com/goloop/trit/v2"
)

// syncBuffer is a bytes.Buffer guarded by a mutex, used only to give the race
// detector a stateful shared writer target.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

// newLogger builds a logger with a single text output to w and no layout
// fields (so records are deterministic).
func newTextLogger(t *testing.T, w interface{ Write([]byte) (int, error) }) *Logger {
	t.Helper()
	l := New()
	if err := l.SetOutputs(Output{
		Name:    "test",
		Writer:  w,
		Levels:  level.Default,
		Layouts: layout.None,
	}); err != nil {
		t.Fatalf("SetOutputs: %v", err)
	}
	return l
}

// TestConcurrentWritesSerialized guards BUG-01: concurrent logging to one
// output with a non-concurrent-safe writer must not race.
func TestConcurrentWritesSerialized(t *testing.T) {
	var buf bytes.Buffer // deliberately unsynchronized
	l := New()
	if err := l.SetOutputs(Output{
		Name:   "test",
		Writer: &buf,
		Levels: level.Default,
	}); err != nil {
		t.Fatalf("SetOutputs: %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				l.Info("concurrent message")
			}
		}()
	}
	wg.Wait()
}

// TestTextRecordsAreDelimited guards BUG-03: successive non-ln text records are
// separated by a newline.
func TestTextRecordsAreDelimited(t *testing.T) {
	var buf syncBuffer
	l := newTextLogger(t, &buf)
	l.Info("first")
	l.Info("second")

	out := buf.buf.String()
	if strings.Count(out, "\n") != 2 {
		t.Errorf("expected two newline-delimited records, got %q", out)
	}
	if strings.Contains(out, "firstsecond") {
		t.Errorf("records are glued: %q", out)
	}
}

// TestLayoutNoneDisablesCallerInfo guards BUG-04: layout.None produces no
// file/func/line, unlike a zero Layouts (which means Default).
func TestLayoutNoneDisablesCallerInfo(t *testing.T) {
	var none, zero syncBuffer

	ln := newTextLogger(t, &none) // uses layout.None
	ln.Info("x")
	if strings.Contains(none.buf.String(), ".go:") {
		t.Errorf("layout.None still printed caller info: %q", none.buf.String())
	}

	lz := New()
	lz.SetOutputs(Output{Name: "z", Writer: &zero, Levels: level.Default, Layouts: 0})
	lz.Info("x")
	if !strings.Contains(zero.buf.String(), ".go:") {
		t.Errorf("Layouts:0 should default to caller info, got %q", zero.buf.String())
	}
}

// TestJSONIsValidJSONL guards BUG-05: JSON records are one valid object per
// line, with no leading/trailing space.
func TestJSONIsValidJSONL(t *testing.T) {
	var buf syncBuffer
	l := New()
	if err := l.SetOutputs(Output{
		Name:      "json",
		Writer:    &buf,
		Levels:    level.Default,
		Layouts:   layout.None,
		TextStyle: trit.False,
	}); err != nil {
		t.Fatalf("SetOutputs: %v", err)
	}
	l.Info("one")
	l.Info("two")

	lines := strings.Split(strings.TrimRight(buf.buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSONL lines, got %d: %q", len(lines), buf.buf.String())
	}
	for _, line := range lines {
		if line != strings.TrimSpace(line) {
			t.Errorf("line has surrounding space: %q", line)
		}
		if !json.Valid([]byte(line)) {
			t.Errorf("line is not valid JSON: %q", line)
		}
	}
}

// TestNewSinglePrefixTrimmed guards BUG-07: a single prefix has surrounding
// whitespace stripped, as documented.
func TestNewSinglePrefixTrimmed(t *testing.T) {
	if got := New(" MYAPP ").prefix; got != "MYAPP" {
		t.Errorf("single-prefix = %q, want %q", got, "MYAPP")
	}
}

// TestLevelWarnAccessor guards BUG-06: Level has a Warn accessor.
func TestLevelWarnAccessor(t *testing.T) {
	l := level.Warn
	if !l.Warn() {
		t.Error("Warn level .Warn() = false, want true")
	}
	other := level.Info
	if other.Warn() {
		t.Error("Info level .Warn() = true, want false")
	}
}

// TestLevelIsSingleUpperBound guards BUG-08: an out-of-range bit is not a valid
// single flag.
func TestLevelIsSingleUpperBound(t *testing.T) {
	valid := level.Warn
	if !valid.IsSingle() {
		t.Error("Warn.IsSingle() = false, want true")
	}
	invalid := level.Level(128)
	if invalid.IsSingle() {
		t.Error("Level(128).IsSingle() = true, want false (out of range)")
	}
}

// overlapWriter counts the calls to Write that are in flight at once and
// remembers whether two ever overlapped. It is itself synchronized, so the
// verdict does not depend on the outcome of a data race.
type overlapWriter struct {
	active  atomic.Int32
	overlap atomic.Bool
}

func (w *overlapWriter) Write(p []byte) (int, error) {
	if w.active.Add(1) > 1 {
		w.overlap.Store(true)
	}
	time.Sleep(time.Millisecond)
	w.active.Add(-1)
	return len(p), nil
}

// TestCopySharesWriteLock guards the copy of a logger: the copy and the
// original write to the same Writer under the same lock, so their writes
// never overlap. The same holds for a logger reconfigured with its own
// outputs and for outputs carried over with Outputs().
func TestCopySharesWriteLock(t *testing.T) {
	run := func(t *testing.T, a, b *Logger, w *overlapWriter) {
		t.Helper()
		var wg sync.WaitGroup
		for _, l := range []*Logger{a, b} {
			wg.Add(1)
			go func(l *Logger) {
				defer wg.Done()
				for i := 0; i < 20; i++ {
					l.Info("shared writer")
				}
			}(l)
		}
		wg.Wait()
		if w.overlap.Load() {
			t.Errorf("two loggers wrote to the same writer at once")
		}
	}

	t.Run("copy", func(t *testing.T) {
		w := &overlapWriter{}
		l := New()
		if err := l.SetOutputs(Output{Name: "test", Writer: w, Layouts: layout.None}); err != nil {
			t.Fatalf("SetOutputs: %v", err)
		}
		run(t, l, l.Copy(), w)
	})

	t.Run("outputs carried over", func(t *testing.T) {
		w := &overlapWriter{}
		l := New()
		if err := l.SetOutputs(Output{Name: "test", Writer: w, Layouts: layout.None}); err != nil {
			t.Fatalf("SetOutputs: %v", err)
		}
		other := New()
		if err := other.SetOutputs(l.Outputs()...); err != nil {
			t.Fatalf("SetOutputs: %v", err)
		}
		run(t, l, other, w)
	})

	t.Run("reconfigured with the same name", func(t *testing.T) {
		w := &overlapWriter{}
		l := New()
		if err := l.SetOutputs(Output{Name: "test", Writer: w, Layouts: layout.None}); err != nil {
			t.Fatalf("SetOutputs: %v", err)
		}
		c := l.Copy()
		// Reconfigure the original from a fresh literal; the copy keeps the
		// previous configuration and must still share the lock.
		if err := l.SetOutputs(Output{Name: "test", Writer: w, Layouts: layout.None}); err != nil {
			t.Fatalf("SetOutputs: %v", err)
		}
		run(t, l, c, w)
	})
}
