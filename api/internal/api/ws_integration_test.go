package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ohmycode_api/internal/store"

	"github.com/gorilla/websocket"
)

// Valid 22-char base62 IDs that pass util.IsValidId.
const (
	wsFileA = "AAAAAAAAAAAAAAAAAAAAAA"
	wsFileB = "BBBBBBBBBBBBBBBBBBBBBB"
	wsApp   = "CCCCCCCCCCCCCCCCCCCCCC"
)

func newWsTestService() *Service {
	return &Service{
		wsAllowedOrigins: []string{},
		fileStore:        store.NewFileStoreInMemory(),
	}
}

func wsTestServer(t *testing.T, svc *Service) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/file", svc.handleWsFileConnection)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func wsConnect(t *testing.T, ts *httptest.Server) *websocket.Conn {
	t.Helper()
	url := "ws" + strings.TrimPrefix(ts.URL, "http") + "/file"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("ws dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func wsSend(t *testing.T, conn *websocket.Conn, v interface{}) {
	t.Helper()
	data, _ := json.Marshal(v)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("ws write: %v", err)
	}
}

func wsRead(t *testing.T, conn *websocket.Conn) map[string]interface{} {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("ws read: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("ws read unmarshal: %v (raw: %s)", err, data)
	}
	return m
}

func wsInit(t *testing.T, conn *websocket.Conn, fileID, name string) {
	t.Helper()
	wsSend(t, conn, map[string]interface{}{
		"action":    "init",
		"file_id":   fileID,
		"file_name": name,
		"app_id":    wsApp,
		"user_id":   wsApp,
		"lang":      "markdown",
		"content":   "",
	})
}

// TestWS_Init_ReceivesFileSnapshot verifies the server sends the file snapshot after init.
func TestWS_Init_ReceivesFileSnapshot(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsInit(t, conn, wsFileA, "My Test File")
	msg := wsRead(t, conn)

	if id, _ := msg["id"].(string); id != wsFileA {
		t.Errorf("id: got %q, want %q", id, wsFileA)
	}
	if name, _ := msg["name"].(string); name != "My Test File" {
		t.Errorf("name: got %q, want 'My Test File'", name)
	}
	if lang, _ := msg["lang"].(string); lang != "markdown" {
		t.Errorf("lang: got %q, want 'markdown'", lang)
	}
}

// TestWS_SetContent_ClientReceivesUpdate verifies that set_content triggers a snapshot back to the sender.
func TestWS_SetContent_ClientReceivesUpdate(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsInit(t, conn, wsFileA, "File")
	wsRead(t, conn) // consume initial snapshot

	wsSend(t, conn, map[string]interface{}{
		"action":  "set_content",
		"content": "hello integration",
	})

	msg := wsRead(t, conn) // throttle ≤500 ms, read deadline 3 s
	if c, _ := msg["content"].(string); c != "hello integration" {
		t.Errorf("content: got %q, want 'hello integration'", c)
	}
}

// TestWS_TwoClients_BroadcastSync verifies that a content change from one client is broadcast to the other.
func TestWS_TwoClients_BroadcastSync(t *testing.T) {
	svc := newWsTestService()
	ts := wsTestServer(t, svc)

	connA := wsConnect(t, ts)
	connB := wsConnect(t, ts)

	wsInit(t, connA, wsFileB, "Shared File")
	wsRead(t, connA) // A's initial snapshot — A is now subscribed

	wsInit(t, connB, wsFileB, "Shared File")
	wsRead(t, connB) // B's initial snapshot — B is now subscribed

	wsSend(t, connA, map[string]interface{}{
		"action":  "set_content",
		"content": "broadcast me",
	})

	// Both A and B receive the update (throttle ≤500 ms, deadline 3 s each)
	msgA := wsRead(t, connA)
	msgB := wsRead(t, connB)

	if c, _ := msgA["content"].(string); c != "broadcast me" {
		t.Errorf("A content: got %q, want 'broadcast me'", c)
	}
	if c, _ := msgB["content"].(string); c != "broadcast me" {
		t.Errorf("B content: got %q, want 'broadcast me'", c)
	}
}

// TestWS_SetName_UpdatePropagates verifies the set_name action reflects in subsequent snapshots.
func TestWS_SetName_UpdatePropagates(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsInit(t, conn, wsFileA, "Original Name")
	wsRead(t, conn)

	wsSend(t, conn, map[string]interface{}{
		"action":    "set_name",
		"file_name": "Renamed",
	})

	msg := wsRead(t, conn)
	if name, _ := msg["name"].(string); name != "Renamed" {
		t.Errorf("name: got %q, want 'Renamed'", name)
	}
}
