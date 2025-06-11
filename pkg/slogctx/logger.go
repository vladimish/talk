package slogctx

import (
	"context"
	"log/slog"
	"sync"
)

type contextKey struct{}

var logFieldsKey = contextKey{}

type logFields struct {
	mu     sync.RWMutex
	fields map[string]any
}

type ContextHandler struct {
	handler slog.Handler
}

func NewContextHandler(h slog.Handler) *ContextHandler {
	return &ContextHandler{handler: h}
}

// Handle processes a log record, adding any fields from the context.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if fields := getFields(ctx); fields != nil {
		fields.mu.RLock()
		for k, v := range fields.fields {
			r.AddAttrs(slog.Any(k, v))
		}
		fields.mu.RUnlock()
	}

	return h.handler.Handle(ctx, r)
}

// WithAttrs returns a new Handler with additional attributes.
func (h *ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{handler: h.handler.WithAttrs(attrs)}
}

// WithGroup returns a new Handler with a group name.
func (h *ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{handler: h.handler.WithGroup(name)}
}

// Enabled reports whether the handler handles records at the given level.
func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithField adds a single key-value pair to the context.
func WithField(ctx context.Context, key string, value any) context.Context {
	fields := getOrCreateFields(ctx)
	fields.mu.Lock()
	fields.fields[key] = value
	fields.mu.Unlock()
	return context.WithValue(ctx, logFieldsKey, fields)
}

// WithFields adds multiple key-value pairs to the context.
func WithFields(ctx context.Context, kvPairs map[string]any) context.Context {
	fields := getOrCreateFields(ctx)
	fields.mu.Lock()
	for k, v := range kvPairs {
		fields.fields[k] = v
	}
	fields.mu.Unlock()
	return context.WithValue(ctx, logFieldsKey, fields)
}

// WithFieldsFromContext copies fields from one context to another.
func WithFieldsFromContext(dst, src context.Context) context.Context {
	if srcFields := getFields(src); srcFields != nil {
		dstFields := getOrCreateFields(dst)

		srcFields.mu.RLock()
		dstFields.mu.Lock()
		for k, v := range srcFields.fields {
			dstFields.fields[k] = v
		}
		dstFields.mu.Unlock()
		srcFields.mu.RUnlock()

		return context.WithValue(dst, logFieldsKey, dstFields)
	}
	return dst
}

// ClearFields removes all fields from the context.
func ClearFields(ctx context.Context) context.Context {
	return context.WithValue(ctx, logFieldsKey, nil)
}

// GetFields returns a copy of all fields in the context.
func GetFields(ctx context.Context) map[string]any {
	fields := getFields(ctx)
	if fields == nil {
		return nil
	}

	fields.mu.RLock()
	defer fields.mu.RUnlock()

	result := make(map[string]any, len(fields.fields))
	for k, v := range fields.fields {
		result[k] = v
	}
	return result
}

func getFields(ctx context.Context) *logFields {
	if ctx == nil {
		return nil
	}
	fields, _ := ctx.Value(logFieldsKey).(*logFields)
	return fields
}

func getOrCreateFields(ctx context.Context) *logFields {
	fields := getFields(ctx)
	if fields == nil {
		fields = &logFields{
			fields: make(map[string]any),
		}
	} else {
		// Create a copy to avoid modifying parent context
		fields.mu.RLock()
		newFields := make(map[string]any, len(fields.fields))
		for k, v := range fields.fields {
			newFields[k] = v
		}
		fields.mu.RUnlock()
		fields = &logFields{
			fields: newFields,
		}
	}
	return fields
}
