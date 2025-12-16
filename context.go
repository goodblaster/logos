package logos

import (
	"context"
)

type contextKey string

// CtxKeyLogger is the key used to store a Logger in context.Context.
const CtxKeyLogger contextKey = "logos.logger"

// FromContext retrieves a Logger from the context. If a logger is found in the context,
// it returns that logger. Otherwise, it returns the DefaultLogger.
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return DefaultLogger
	}

	if ctxLog, ok := ctx.Value(CtxKeyLogger).(Logger); ok && ctxLog.level != nil {
		return ctxLog
	}

	return DefaultLogger
}

// WithLogger returns a new context with the provided logger stored in it.
// The logger can later be retrieved using FromContext.
func WithLogger(ctx context.Context, logger Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, CtxKeyLogger, logger)
}
