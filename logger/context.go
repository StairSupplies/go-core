package logger

import (
	"context"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey int

// loggerKey is the key for logger values in contexts
const loggerKey contextKey = iota

// NewContext creates a new context with the provided logger
func NewContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithContext returns the logger associated with the context, or a default logger if none exists
func WithContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(loggerKey).(*Logger); ok {
		return logger
	}
	
	// Return a default logger
	logger, _ := New()
	return logger
}