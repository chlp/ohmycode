package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ohmycode_api/internal/store"
)

func newHealthTestService() *Service {
	return &Service{
		fileStore:   store.NewFileStoreInMemory(),
		runnerStore: store.NewRunnerStore(),
	}
}

func healthTestServer(t *testing.T, svc *Service) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", svc.handleHealth)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func TestHealth_Returns200(t *testing.T) {
	ts := healthTestServer(t, newHealthTestService())
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want 200", resp.StatusCode)
	}
}

func TestHealth_ContentTypeJSON(t *testing.T) {
	ts := healthTestServer(t, newHealthTestService())
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type: got %q, want application/json", ct)
	}
}

func TestHealth_BodyWithInMemoryStore(t *testing.T) {
	ts := healthTestServer(t, newHealthTestService())
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body struct {
		Status        string `json:"status"`
		Mongo         string `json:"mongo"`
		PersistErrors int    `json:"persist_errors"`
		RunnersOnline int    `json:"runners_online"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("status: got %q, want 'ok'", body.Status)
	}
	if body.Mongo != "ok" {
		t.Errorf("mongo: got %q, want 'ok'", body.Mongo)
	}
	if body.PersistErrors != 0 {
		t.Errorf("persist_errors: got %d, want 0", body.PersistErrors)
	}
	if body.RunnersOnline != 0 {
		t.Errorf("runners_online: got %d, want 0", body.RunnersOnline)
	}
}

func TestHealth_RunnersOnlineCount(t *testing.T) {
	svc := newHealthTestService()
	const rID = "AAAAAAAAAAAAAAAAAAAAAA"
	svc.runnerStore.SetRunner(rID, false)

	ts := healthTestServer(t, svc)
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body struct {
		RunnersOnline int `json:"runners_online"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.RunnersOnline != 1 {
		t.Errorf("runners_online: got %d, want 1", body.RunnersOnline)
	}
}
