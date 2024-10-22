package util

import (
	"context"
	"fmt"
	"log"
	"time"
)

const RequestStartTimeCtxKey string = "RequestStartTime"

func Log(ctx context.Context, msg string) {
	if ctx == nil {
		ctx = context.Background()
	}
	elapsedTimeStr := ""
	if startTime, ok := ctx.Value(RequestStartTimeCtxKey).(time.Time); ok {
		elapsedTime := time.Since(startTime)
		elapsedTimeStr = fmt.Sprintf(" (%0.3f)", elapsedTime.Seconds())
	}
	log.Printf("%s%s: %s\n", time.Now().Format("2006-01-02 15:04:05.000"), elapsedTimeStr, msg)
}
