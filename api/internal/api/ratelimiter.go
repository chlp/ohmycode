package api

import (
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	runRateLimitPerWindow = 30
	runRateLimitWindow    = time.Minute
)

type runBucket struct {
	count   int
	resetAt time.Time
}

type ipRunLimiter struct {
	mu      sync.Mutex
	buckets map[string]*runBucket
}

func newIpRunLimiter() *ipRunLimiter {
	return &ipRunLimiter{buckets: make(map[string]*runBucket)}
}

func (l *ipRunLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	b, ok := l.buckets[ip]
	if !ok || now.After(b.resetAt) {
		l.buckets[ip] = &runBucket{count: 1, resetAt: now.Add(runRateLimitWindow)}
		return true
	}
	if b.count >= runRateLimitPerWindow {
		return false
	}
	b.count++
	return true
}

func (l *ipRunLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	for ip, b := range l.buckets {
		if now.After(b.resetAt) {
			delete(l.buckets, ip)
		}
	}
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
