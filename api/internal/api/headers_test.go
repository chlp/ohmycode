package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func nopHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestHeadersMiddleware_CORS_AllowAll(t *testing.T) {
	h := headersMiddleware(nopHandler())
	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	if got := rw.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("ACAO: got %q, want *", got)
	}
}

func TestHeadersMiddleware_Options_Returns200AndSkipsInner(t *testing.T) {
	called := false
	h := headersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	req := httptest.NewRequest(http.MethodOptions, "/any", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("status: got %d, want 200", rw.Code)
	}
	if called {
		t.Error("inner handler should not be called for OPTIONS")
	}
}

func TestHeadersMiddleware_FileEndpoint_NoCacheStore(t *testing.T) {
	h := headersMiddleware(nopHandler())
	for _, path := range []string{"/file", "/runner"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rw := httptest.NewRecorder()
		h.ServeHTTP(rw, req)
		cc := rw.Header().Get("Cache-Control")
		if !strings.Contains(cc, "no-store") {
			t.Errorf("path %s: Cache-Control=%q, want no-store", path, cc)
		}
	}
}

func TestHeadersMiddleware_WsUpgradeRequest_NoCacheStore(t *testing.T) {
	h := headersMiddleware(nopHandler())
	req := httptest.NewRequest(http.MethodGet, "/any", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	cc := rw.Header().Get("Cache-Control")
	if !strings.Contains(cc, "no-store") {
		t.Errorf("WS upgrade: Cache-Control=%q, want no-store", cc)
	}
}

func TestHeadersMiddleware_OtherRoute_PublicCache(t *testing.T) {
	h := headersMiddleware(nopHandler())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	cc := rw.Header().Get("Cache-Control")
	if !strings.Contains(cc, "public") {
		t.Errorf("Cache-Control=%q, want 'public'", cc)
	}
	if !strings.Contains(cc, "max-age=86400") {
		t.Errorf("Cache-Control=%q, want 'max-age=86400'", cc)
	}
}

func TestHeadersMiddleware_AllowMethodsHeader(t *testing.T) {
	h := headersMiddleware(nopHandler())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)

	methods := rw.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("Access-Control-Allow-Methods should be set")
	}
}
