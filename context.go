package log

import "context"

// loggerKey is the context key a request-scoped logger is stored under. It is
// an unexported type, so nothing outside this package can collide with it or
// replace what is there.
type loggerKey struct{}

// WithContext returns a copy of ctx carrying l.
//
// It exists so that a logger prepared once for a unit of work - a request, a
// job, a message - travels with that work instead of being threaded through
// every function that might want to log. The alternative in practice is not
// threading it: log lines end up written against a package-level logger, with
// nothing tying them to the request that produced them, and joining them
// afterwards means matching on timestamps.
//
// A nil logger is stored as given; [FromContext] reports it as absent, so a
// caller who cleared it deliberately gets the fallback rather than a panic.
func WithContext(ctx context.Context, l *Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext returns the logger stored in ctx, or nil when there is none.
//
// Use [FromContextOr] where a logger is required; this form is for code that
// wants to know whether one was provided at all.
func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return nil
	}
	l, _ := ctx.Value(loggerKey{}).(*Logger)
	return l
}

// FromContextOr returns the logger stored in ctx, or fallback when there is
// none:
//
//	log.FromContextOr(ctx, defaultLogger).Errorf("upstream refused: %v", err)
//
// This is the form to reach for at a call site. It never returns nil unless
// the fallback is nil, so a handler that may or may not run inside a prepared
// context does not need to branch, and a missing logger degrades to logging
// somewhere rather than to a panic or to silence.
func FromContextOr(ctx context.Context, fallback *Logger) *Logger {
	if l := FromContext(ctx); l != nil {
		return l
	}
	return fallback
}
