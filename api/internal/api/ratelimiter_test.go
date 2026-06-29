package api

import (
	"net/http/httptest"
	"testing"
)

// --- ipRunLimiter ---

func TestIpRunLimiter_AllowsFirst30(t *testing.T) {
	l := newIpRunLimiter()
	for i := 0; i < 30; i++ {
		if !l.Allow("1.2.3.4") {
			t.Fatalf("Allow returned false on call %d, expected true", i+1)
		}
	}
}

func TestIpRunLimiter_Blocks31st(t *testing.T) {
	l := newIpRunLimiter()
	for i := 0; i < 30; i++ {
		l.Allow("1.2.3.4")
	}
	if l.Allow("1.2.3.4") {
		t.Error("31st call should return false (rate limit hit)")
	}
}

func TestIpRunLimiter_DifferentIPs_Independent(t *testing.T) {
	l := newIpRunLimiter()
	for i := 0; i < 30; i++ {
		l.Allow("1.1.1.1")
	}
	if !l.Allow("2.2.2.2") {
		t.Error("different IP should start with a fresh bucket")
	}
}

func TestIpRunLimiter_BucketCountsCorrectly(t *testing.T) {
	l := newIpRunLimiter()
	const ip = "10.0.0.1"
	for i := 0; i < 29; i++ {
		if !l.Allow(ip) {
			t.Fatalf("call %d returned false, expected true", i+1)
		}
	}
	// 30th call: still within limit
	if !l.Allow(ip) {
		t.Error("30th call should be allowed")
	}
	// 31st call: over limit
	if l.Allow(ip) {
		t.Error("31st call should be blocked")
	}
}

// --- clientIP ---

func TestClientIP_XRealIPTakesPriority(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "192.168.1.1:9999"
	if got := clientIP(req); got != "10.0.0.1" {
		t.Errorf("clientIP: got %q, want '10.0.0.1'", got)
	}
}

func TestClientIP_FallsBackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.5:8080"
	if got := clientIP(req); got != "192.168.1.5" {
		t.Errorf("clientIP: got %q, want '192.168.1.5'", got)
	}
}

// --- isWsOriginAllowed ---

func TestIsWsOriginAllowed_NoOriginHeader_Allow(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	if !isWsOriginAllowed(req, []string{"example.com"}) {
		t.Error("missing Origin header should be allowed (non-browser clients)")
	}
}

func TestIsWsOriginAllowed_EmptyAllowedList_AllowAll(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://anyone.com")
	if !isWsOriginAllowed(req, []string{}) {
		t.Error("empty allowed list should allow all origins (backwards-compatible)")
	}
}

func TestIsWsOriginAllowed_WildcardAllowed_AllowAny(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://any.com")
	if !isWsOriginAllowed(req, []string{"*"}) {
		t.Error("wildcard entry should allow any origin")
	}
}

func TestIsWsOriginAllowed_MatchByHostname(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://example.com")
	if !isWsOriginAllowed(req, []string{"example.com"}) {
		t.Error("bare hostname match should be allowed")
	}
}

func TestIsWsOriginAllowed_MatchBySchemeHost(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://example.com")
	if !isWsOriginAllowed(req, []string{"https://example.com"}) {
		t.Error("full scheme+host match should be allowed")
	}
}

func TestIsWsOriginAllowed_NoMatch_Deny(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://evil.com")
	if isWsOriginAllowed(req, []string{"example.com", "trusted.com"}) {
		t.Error("non-matching origin should be denied")
	}
}

func TestIsWsOriginAllowed_SchemeHostMismatch_Deny(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://example.com") // http, not https
	if isWsOriginAllowed(req, []string{"https://example.com"}) {
		t.Error("http origin should not match https allowed entry")
	}
}
