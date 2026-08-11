package log

import (
	"context"
	"testing"
)

func TestContextRoundTrip(t *testing.T) {
	l := New("api")
	ctx := WithContext(context.Background(), l)

	if got := FromContext(ctx); got != l {
		t.Errorf("FromContext() = %v, want the logger that went in", got)
	}
	if got := FromContextOr(ctx, New("fallback")); got != l {
		t.Error("FromContextOr() preferred the fallback over the stored logger")
	}
}

// A call site must not have to branch on whether a logger was provided: the
// point of the fallback is that a missing logger degrades to logging
// somewhere, not to a panic and not to silence.
func TestFromContextWithoutOne(t *testing.T) {
	if got := FromContext(context.Background()); got != nil {
		t.Errorf("FromContext() = %v, want nil when none was stored", got)
	}

	fallback := New("fallback")
	if got := FromContextOr(context.Background(), fallback); got != fallback {
		t.Error("FromContextOr() did not fall back")
	}

	// A caller who stored nil deliberately gets the fallback too, rather
	// than a nil logger that panics on first use.
	ctx := WithContext(context.Background(), nil)
	if got := FromContextOr(ctx, fallback); got != fallback {
		t.Error("a nil logger in the context was returned instead of the fallback")
	}
}

// A nil context is a caller mistake that should not become a panic inside a
// logging helper, of all places.
func TestContextHelpersTolerateNilContext(t *testing.T) {
	//nolint:staticcheck // passing a nil context is exactly what is tested.
	if ctx := WithContext(nil, New("api")); FromContext(ctx) == nil {
		t.Error("WithContext(nil, l) lost the logger")
	}
	//nolint:staticcheck // passing a nil context is exactly what is tested.
	if got := FromContext(nil); got != nil {
		t.Errorf("FromContext(nil) = %v, want nil", got)
	}
}
