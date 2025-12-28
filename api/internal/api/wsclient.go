package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	stateMu     sync.RWMutex
	file        *model.File
	userId      string
	appId       string
	runner      *model.Runner
	lastUpdate  time.Time
	conn        *websocket.Conn
	done        chan struct{}
	close       func()
	writeMu     sync.Mutex
	fileSetCh   chan struct{}
	runnerSetCh chan struct{}
}

func isIgnorableWsErr(err error) bool {
	if err == nil {
		return false
	}
	// Normal/expected WS close errors.
	var closeErr *websocket.CloseError
	if errors.As(err, &closeErr) {
		return true
	}
	// Expected network timeouts (client gone / proxy closing).
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	// Common TCP-level disconnects.
	if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection reset by peer") ||
		strings.Contains(msg, "use of closed network connection") ||
		strings.Contains(msg, "i/o timeout") {
		return true
	}
	return false
}

func isWsOriginAllowed(r *http.Request, allowed []string) bool {
	// Non-browser clients may not send Origin.
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	// Default: allow all (backwards-compatible).
	if len(allowed) == 0 {
		return true
	}
	for _, a := range allowed {
		if a == "*" {
			return true
		}
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	originHostPort := u.Host
	originHost := u.Hostname()
	originFull := ""
	if u.Scheme != "" && u.Host != "" {
		originFull = u.Scheme + "://" + u.Host
	}

	for _, a := range allowed {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if a == "*" {
			return true
		}
		if strings.HasPrefix(a, "http://") || strings.HasPrefix(a, "https://") {
			if originFull != "" && a == originFull {
				return true
			}
			continue
		}
		// Allow host[:port] or bare hostname.
		if a == originHostPort || a == originHost {
			return true
		}
	}
	return false
}

const wsMessageLimit = 4 * (1 << 20) // 4 Mb
const (
	wsPongWait   = 10 * time.Second
	wsPingPeriod = 5 * time.Second
	wsWriteWait  = 5 * time.Second
)

func createWsClient(w http.ResponseWriter, r *http.Request, allowedOrigins []string) *wsClient {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return isWsOriginAllowed(r, allowedOrigins)
		},
	}
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		if !isIgnorableWsErr(err) {
			util.Log("WebSocket Upgrade error: " + err.Error())
		}
		return nil
	}
	conn.SetReadLimit(wsMessageLimit)

	done := make(chan struct{})
	var once sync.Once
	closeClient := func() {
		once.Do(func() {
			close(done)
			_ = conn.Close()
		})
	}

	_ = conn.SetReadDeadline(time.Now().Add(wsPongWait))

	client := wsClient{
		conn:        conn,
		done:        done,
		close:       closeClient,
		fileSetCh:   make(chan struct{}, 1),
		runnerSetCh: make(chan struct{}, 1),
	}
	go client.pingPongHandling()
	return &client
}

func (client *wsClient) pingPongHandling() {
	client.conn.SetPongHandler(func(appData string) error {
		if err := client.conn.SetReadDeadline(time.Now().Add(wsPongWait)); err != nil {
			util.Log("pong handler err: ", err.Error())
			client.close()
		}
		return nil
	})
	ticker := time.NewTicker(wsPingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-client.done:
			return
		case <-ticker.C:
			client.writeMu.Lock()
			_ = client.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
			err := client.conn.WriteMessage(websocket.PingMessage, nil)
			client.writeMu.Unlock()
			if err != nil {
				if !isIgnorableWsErr(err) {
					util.Log("Ping error: " + err.Error())
				}
				client.close()
				return
			}
		}
	}
}

func (client *wsClient) send(v interface{}) error {
	if v == nil {
		util.Log("wsClient.send nil")
		return nil
	}

	jsonData, err := json.Marshal(v)
	if err != nil {
		util.Log("wsClient.send json err: " + err.Error())
		return err
	}

	if bytes.Equal(jsonData, []byte("null")) {
		util.Log("wsClient.send null")
		return nil
	}

	client.writeMu.Lock()
	defer client.writeMu.Unlock()
	_ = client.conn.SetWriteDeadline(time.Now().Add(wsWriteWait))
	return client.conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (client *wsClient) setFile(file *model.File, appId, userId string) {
	client.stateMu.Lock()
	client.file = file
	client.appId = appId
	client.userId = userId
	client.lastUpdate = time.Time{}
	client.stateMu.Unlock()
	select {
	case client.fileSetCh <- struct{}{}:
	default:
	}
}

func (client *wsClient) getFile() *model.File {
	client.stateMu.RLock()
	f := client.file
	client.stateMu.RUnlock()
	return f
}

func (client *wsClient) getUserId() string {
	client.stateMu.RLock()
	u := client.userId
	client.stateMu.RUnlock()
	return u
}

func (client *wsClient) getAppId() string {
	client.stateMu.RLock()
	a := client.appId
	client.stateMu.RUnlock()
	return a
}

func (client *wsClient) getLastUpdate() time.Time {
	client.stateMu.RLock()
	t := client.lastUpdate
	client.stateMu.RUnlock()
	return t
}

func (client *wsClient) setLastUpdate(t time.Time) {
	client.stateMu.Lock()
	client.lastUpdate = t
	client.stateMu.Unlock()
}

func (client *wsClient) setRunner(runner *model.Runner) {
	client.stateMu.Lock()
	client.runner = runner
	client.stateMu.Unlock()
	select {
	case client.runnerSetCh <- struct{}{}:
	default:
	}
}

func (client *wsClient) getRunner() *model.Runner {
	client.stateMu.RLock()
	r := client.runner
	client.stateMu.RUnlock()
	return r
}
