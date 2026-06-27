package log

import (
	"context"
	"log/slog"
	"strings"

	"github.com/goloop/log/v2/level"
)

// Handler returns a slog.Handler that routes slog records through this
// logger's configured outputs. The slog level is mapped onto the logger
// levels (Debug, Info, Warn, Error) and the record attributes are appended
// to the message as space-separated key=value pairs.
//
// In JSON-style outputs the attributes therefore appear inside the message
// field rather than as separate JSON keys.
func (logger *Logger) Handler() slog.Handler {
	return &slogHandler{logger: logger}
}

// NewSlog returns a *slog.Logger backed by a fresh logger whose outputs are
// the usual Stdout and Stderr. It is a convenience wrapper around
// New(prefixes...).Handler().
func NewSlog(prefixes ...string) *slog.Logger {
	return slog.New(New(prefixes...).Handler())
}

// The slogHandler adapts a Logger to the slog.Handler interface.
//
// The prefix field holds the already-rendered " key=value" fragments
// contributed by WithAttrs, while groups holds the active group names that
// qualify the keys of the record's own attributes.
type slogHandler struct {
	logger *Logger
	prefix string
	groups []string
}

// slogLevel maps a slog.Level onto a logger level.
func slogLevel(l slog.Level) level.Level {
	switch {
	case l < slog.LevelInfo:
		return level.Debug
	case l < slog.LevelWarn:
		return level.Info
	case l < slog.LevelError:
		return level.Warn
	default:
		return level.Error
	}
}

// Enabled reports whether a record at the given slog level would be emitted.
func (h *slogHandler) Enabled(_ context.Context, l slog.Level) bool {
	return h.logger.Enabled(slogLevel(l))
}

// Handle renders the record (message plus attributes) and emits it through
// the logger, using the record's program counter as the call-site frame.
func (h *slogHandler) Handle(_ context.Context, r slog.Record) error {
	var sb strings.Builder
	sb.WriteString(r.Message)
	sb.WriteString(h.prefix)

	group := strings.Join(h.groups, ".")
	r.Attrs(func(a slog.Attr) bool {
		appendAttr(&sb, group, a)
		return true
	})

	var frame *stackFrame
	if r.PC != 0 {
		if f, ok := frameFromPC(r.PC); ok {
			frame = f
		}
	}

	h.logger.emit(nil, slogLevel(r.Level), kindPrint, sb.String(), frame)
	return nil
}

// WithAttrs returns a handler that pre-renders attrs with the current group
// qualification and remembers them for every subsequent record.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	var sb strings.Builder
	sb.WriteString(h.prefix)

	group := strings.Join(h.groups, ".")
	for _, a := range attrs {
		appendAttr(&sb, group, a)
	}

	return &slogHandler{logger: h.logger, prefix: sb.String(), groups: h.groups}
}

// WithGroup returns a handler whose subsequent attribute keys are qualified
// with name.
func (h *slogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	groups := make([]string, len(h.groups)+1)
	copy(groups, h.groups)
	groups[len(h.groups)] = name

	return &slogHandler{logger: h.logger, prefix: h.prefix, groups: groups}
}

// appendAttr writes a single attribute as " group.key=value" into sb. Empty
// attributes are skipped and group-valued attributes are expanded
// recursively with their key folded into the qualifier.
func appendAttr(sb *strings.Builder, group string, a slog.Attr) {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return
	}

	if a.Value.Kind() == slog.KindGroup {
		sub := a.Value.Group()
		if len(sub) == 0 {
			return
		}

		inner := a.Key
		if inner != "" && group != "" {
			inner = group + "." + inner
		} else if inner == "" {
			inner = group
		}

		for _, ga := range sub {
			appendAttr(sb, inner, ga)
		}
		return
	}

	sb.WriteByte(' ')
	if group != "" {
		sb.WriteString(group)
		sb.WriteByte('.')
	}
	sb.WriteString(a.Key)
	sb.WriteByte('=')
	sb.WriteString(a.Value.String())
}
