package util

import (
	"context"
	"fmt"
	"time"
)

const RequestStartTimeCtxKey string = "RequestStartTime"

func Log(args ...interface{}) {
	var ctx context.Context = nil
	var msg string
	var ok bool

	switch len(args) {
	case 1:
		msg, ok = args[0].(string)
		if !ok {
			msg = fmt.Sprintf("Log error: first argument must be string, got %T", args[0])
		}
	case 2:
		ctx, _ = args[0].(context.Context)
		msg, ok = args[1].(string)
		if !ok {
			msg = fmt.Sprintf("Log error: second argument must be string, got %T", args[1])
		}
	default:
		msg = fmt.Sprintf("Log error: wrong usage, got %d arguments", len(args))
	}

	if ctx == nil {
		ctx = context.Background()
	}
	elapsedTimeStr := ""
	if startTime, ok := ctx.Value(RequestStartTimeCtxKey).(time.Time); ok {
		elapsedTime := time.Since(startTime)
		elapsedTimeStr = fmt.Sprintf(" (%0.3f)", elapsedTime.Seconds())
	}
	fmt.Printf("%s%s: %s\n", time.Now().Format("2006-01-02 15:04:05.000"), elapsedTimeStr, msg)
}
