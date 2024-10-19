package util

import (
	"log"
	"time"
)

var startTime time.Time

func Timer() float64 {
	if startTime.IsZero() {
		startTime = time.Now()
		return 0
	}
	return time.Since(startTime).Seconds()
}

func Log(str string) {
	log.Printf("%s (%0.3f): %s\n", time.Now().Format("2006-01-02 15:04:05.000"), Timer(), str)
}
