package util

import (
	"context"
	"log/slog"
	"os"
	"time"
)

const RequestStartTimeCtxKey string = "RequestStartTime"

var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

// Log is a backwards-compatible logger that accepts (ctx, msg) or (msg).
// For structured logging, use LogInfo / LogError.
func Log(args ...interface{}) {
	var ctx context.Context
	var msg string
	switch len(args) {
	case 1:
		msg, _ = args[0].(string)
	case 2:
		ctx, _ = args[0].(context.Context)
		msg, _ = args[1].(string)
	}
	attrs := []any{}
	if ctx != nil {
		if startTime, ok := ctx.Value(RequestStartTimeCtxKey).(time.Time); ok {
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
