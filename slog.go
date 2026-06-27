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
// The attrs field holds the flattened, group-qualified key/value pairs
// contributed by WithAttrs, while groups holds the active group names that
// qualify the keys of the record's own attributes.
type slogHandler struct {
	logger *Logger
	attrs  []logField
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

// Handle flattens the record's attributes (combined with the handler's own)
// and emits the message and fields through the logger, using the record's
// program counter as the call-site frame.
func (h *slogHandler) Handle(_ context.Context, r slog.Record) error {
	fields := h.attrs
	if r.NumAttrs() > 0 {
		combined := make([]logField, len(h.attrs), len(h.attrs)+r.NumAttrs())
		copy(combined, h.attrs)

		group := strings.Join(h.groups, ".")
		rec := make([]slog.Attr, 0, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			rec = append(rec, a)
			return true
		})
		fields = flattenAttrs(group, rec, combined)
	}

	var frame *stackFrame
	if r.PC != 0 {
		if f, ok := frameFromPC(r.PC); ok {
			frame = f
		}
	}

	h.logger.emit(nil, slogLevel(r.Level), kindPrint, r.Message, frame, fields)
	return nil
}

// WithAttrs returns a handler that remembers attrs (flattened and qualified
// with the current group) for every subsequent record.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	next := make([]logField, len(h.attrs), len(h.attrs)+len(attrs))
	copy(next, h.attrs)
	next = flattenAttrs(strings.Join(h.groups, "."), attrs, next)

	return &slogHandler{logger: h.logger, attrs: next, groups: h.groups}
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

	return &slogHandler{logger: h.logger, attrs: h.attrs, groups: groups}
}

// flattenAttrs appends each attribute to dst as a key/value field, resolving
// LogValuer values, qualifying keys with group and expanding group-valued
// attributes recursively (their key folded into the qualifier). Empty
// attributes are skipped.
func flattenAttrs(group string, attrs []slog.Attr, dst []logField) []logField {
	for _, a := range attrs {
		a.Value = a.Value.Resolve()
		if a.Equal(slog.Attr{}) {
			continue
		}

		key := a.Key
		if group != "" {
			if key != "" {
				key = group + "." + key
			} else {
				key = group
			}
		}

		if a.Value.Kind() == slog.KindGroup {
			dst = flattenAttrs(key, a.Value.Group(), dst)
			continue
		}

		dst = append(dst, logField{key: key, val: a.Value.Any()})
	}

	return dst
}
