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

// TestWS_SetLang_UpdatePropagates verifies set_lang changes the lang in the file snapshot.
func TestWS_SetLang_UpdatePropagates(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsInit(t, conn, wsFileA, "File")
	wsRead(t, conn) // initial snapshot: lang=markdown

	wsSend(t, conn, map[string]interface{}{
		"action": "set_lang",
		"lang":   "go",
	})

	msg := wsRead(t, conn)
	if lang, _ := msg["lang"].(string); lang != "go" {
		t.Errorf("lang: got %q, want 'go'", lang)
	}
}

// TestWS_SetLang_Locked_NoUpdate verifies set_lang is silently ignored when the file is locked.
func TestWS_SetLang_Locked_NoUpdate(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsInit(t, conn, wsFileA, "File")
	wsRead(t, conn)

	wsSend(t, conn, map[string]interface{}{
		"action":    "set_locked",
		"is_locked": true,
	})
	wsRead(t, conn) // snapshot after lock (is_locked=true, lang=markdown)

	wsSend(t, conn, map[string]interface{}{
		"action": "set_lang",
		"lang":   "go",
	})

	// set_lang is ignored; send set_name to get the next snapshot and verify lang unchanged.
	wsSend(t, conn, map[string]interface{}{
		"action":    "set_name",
		"file_name": "Probe",
	})
	msg := wsRead(t, conn)
	if lang, _ := msg["lang"].(string); lang != "markdown" {
		t.Errorf("lang should stay 'markdown' while locked, got %q", lang)
	}
}

// TestWS_SetEncrypted_Enable verifies set_encrypted=true sets encrypted=true and generates ro_token.
func TestWS_SetEncrypted_Enable(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsInit(t, conn, wsFileA, "File")
	wsRead(t, conn)

	wsSend(t, conn, map[string]interface{}{
		"action":    "set_encrypted",
		"encrypted": true,
	})

	msg := wsRead(t, conn)
	if enc, _ := msg["encrypted"].(bool); !enc {
		t.Error("expected encrypted=true after set_encrypted")
	}
	if tok, _ := msg["ro_token"].(string); tok == "" {
		t.Error("expected non-empty ro_token after enabling encryption")
	}
}

// TestWS_ROToken_BlocksSetContent verifies that a client using the read-only token cannot modify content.
func TestWS_ROToken_BlocksSetContent(t *testing.T) {
	svc := newWsTestService()
	ts := wsTestServer(t, svc)

	// Owner: create file and enable encryption to get ro_token.
	connOwner := wsConnect(t, ts)
	wsInit(t, connOwner, wsFileA, "Encrypted File")
	wsRead(t, connOwner)

	wsSend(t, connOwner, map[string]interface{}{
		"action":    "set_encrypted",
		"encrypted": true,
	})
	snap := wsRead(t, connOwner)
	roToken, _ := snap["ro_token"].(string)
	if roToken == "" {
		t.Fatal("ro_token not returned after set_encrypted")
	}

	// RO client: connect with ro_token.
	connRO := wsConnect(t, ts)
	wsSend(t, connRO, map[string]interface{}{
		"action":    "init",
		"file_id":   wsFileA,
		"file_name": "Encrypted File",
		"app_id":    wsApp,
		"user_id":   wsApp,
		"lang":      "markdown",
		"content":   "",
		"ro_token":  roToken,
	})
	wsRead(t, connRO)

	// RO client tries to set content — should be silently ignored.
	wsSend(t, connRO, map[string]interface{}{
		"action":  "set_content",
		"content": "should not be saved",
	})

	// Verify via a fresh third client that content is unchanged.
	connVerify := wsConnect(t, ts)
	wsInit(t, connVerify, wsFileA, "Encrypted File")
	snapVerify := wsRead(t, connVerify)
	if content, ok := snapVerify["content"]; ok {
		if c, _ := content.(string); c == "should not be saved" {
			t.Error("RO client must not be able to set_content")
		}
	}
}

// TestWS_InvalidFileId_DisconnectsClient verifies init with a malformed file_id closes the connection.
func TestWS_InvalidFileId_DisconnectsClient(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsSend(t, conn, map[string]interface{}{
		"action":    "init",
		"file_id":   "not-valid!!",
		"file_name": "File",
		"app_id":    wsApp,
		"user_id":   wsApp,
		"lang":      "markdown",
		"content":   "",
	})

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Error("expected connection close after invalid file_id")
	}
}

// TestWS_UnknownAction_ConnRemainsOpen verifies that an unknown action does not close the connection.
func TestWS_UnknownAction_ConnRemainsOpen(t *testing.T) {
	ts := wsTestServer(t, newWsTestService())
	conn := wsConnect(t, ts)

	wsInit(t, conn, wsFileA, "File")
	wsRead(t, conn)

	wsSend(t, conn, map[string]interface{}{
		"action": "totally_unknown_action_xyz",
	})

	// Connection should survive; a valid set_name should still work.
	wsSend(t, conn, map[string]interface{}{
		"action":    "set_name",
		"file_name": "Still Works",
	})
	msg := wsRead(t, conn)
	if name, _ := msg["name"].(string); name != "Still Works" {
		t.Errorf("name after unknown action: got %q, want 'Still Works'", name)
	}
}

// TestWS_TwoClients_SetContent_BothReceiveUpdate verifies three-client fan-out: C reads A's content.
func TestWS_TwoClients_SetContent_BothReceiveUpdate(t *testing.T) {
	svc := newWsTestService()
	ts := wsTestServer(t, svc)

	const fileID = "CCCCCCCCCCCCCCCCCCCCCC"
	connA := wsConnect(t, ts)
	connC := wsConnect(t, ts)

	wsInit(t, connA, fileID, "File")
	wsRead(t, connA)

	wsInit(t, connC, fileID, "File")
	wsRead(t, connC)

	wsSend(t, connA, map[string]interface{}{
		"action":  "set_content",
		"content": "fan out test",
	})

	msgA := wsRead(t, connA)
	msgC := wsRead(t, connC)
	for label, msg := range map[string]map[string]interface{}{"A": msgA, "C": msgC} {
		if c, _ := msg["content"].(string); c != "fan out test" {
			t.Errorf("client %s content: got %q, want 'fan out test'", label, c)
		}
	}
}
