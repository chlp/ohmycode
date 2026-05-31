package util

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// ctxKey is a private context key type to avoid collisions with other packages.
type ctxKey int

const requestStartTimeCtxKey ctxKey = iota

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

func WithRequestStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestStartTimeCtxKey, t)
}

func RequestStartTime(ctx context.Context) (time.Time, bool) {
	if ctx == nil {
		return time.Time{}, false
	}
	startTime, ok := ctx.Value(requestStartTimeCtxKey).(time.Time)
	return startTime, ok
}

// Log is a backwards-compatible unstructured logger.
// For structured logging, use LogInfo / LogError / LogDebug.
func Log(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	var ctx context.Context
	var msg string
	if c, ok := args[0].(context.Context); ok {
		ctx = c
		msg = fmt.Sprint(args[1:]...)
	} else {
		msg = fmt.Sprint(args...)
	}
	attrs := []any{}
	if ctx != nil {
		if startTime, ok := RequestStartTime(ctx); ok {
			attrs = append(attrs, "elapsed_s", time.Since(startTime).Seconds())
		}
	}
	logger.Info(msg, attrs...)
}

func LogInfo(msg string, keysAndValues ...any) {
	logger.Info(msg, keysAndValues...)
}

func LogError(msg string, keysAndValues ...any) {
	logger.Error(msg, keysAndValues...)
}

func LogDebug(msg string, keysAndValues ...any) {
	logger.Debug(msg, keysAndValues...)
}
