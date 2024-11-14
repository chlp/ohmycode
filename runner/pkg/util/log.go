package util

import (
	"context"
	"fmt"
	"log"
	"time"
)

const RequestStartTimeCtxKey string = "RequestStartTime"

func Log(args ...interface{}) {
	var ctx context.Context = nil
	var msg string

	switch len(args) {
	case 1:
		msg, _ = args[0].(string)
	case 2:
		ctx, _ = args[0].(context.Context)
		msg, _ = args[1].(string)
	default:
		log.Fatal("wrong Log usage")
		return
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
