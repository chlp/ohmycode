package util

import (
	"context"
	"fmt"
	"time"
)

// ctxKey is a private context key type to avoid collisions with other packages.
type ctxKey int

const requestStartTimeCtxKey ctxKey = iota

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

func Log(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	var ctx context.Context
	var msg string

	// Optional leading context.Context; everything else becomes the message.
	if c, ok := args[0].(context.Context); ok {
		ctx = c
		msg = fmt.Sprint(args[1:]...)
	} else {
		msg = fmt.Sprint(args...)
	}

	if ctx == nil {
		ctx = context.Background()
	}
	elapsedTimeStr := ""
	if startTime, ok := RequestStartTime(ctx); ok {
		elapsedTime := time.Since(startTime)
		elapsedTimeStr = fmt.Sprintf(" (%0.3f)", elapsedTime.Seconds())
	}
	fmt.Printf("%s%s: %s\n", time.Now().Format("2006-01-02 15:04:05.000"), elapsedTimeStr, msg)
}
