package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"sync"
	"time"
)

type wsClient struct {
	stateMu sync.RWMutex
	file       *model.File
	userId     string
	appId      string
	runner     *model.Runner
	lastUpdate time.Time
	conn       *websocket.Conn
	done       chan struct{}
	close      func()
	writeMu    sync.Mutex
	fileSetCh  chan struct{}
	runnerSetCh chan struct{}
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true // could be return r.Host == "yourdomain.com"
	},
}

const wsMessageLimit = 4 * (1 << 20) // 4 Mb
const (
	wsPongWait   = 10 * time.Second
	wsPingPeriod = 5 * time.Second
	wsWriteWait  = 5 * time.Second
)

func createWsClient(w http.ResponseWriter, r *http.Request) *wsClient {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		util.Log("WebSocket Upgrade error: " + err.Error())
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
		conn:      conn,
		done:      done,
		close:     closeClient,
		fileSetCh: make(chan struct{}, 1),
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
				util.Log("Ping error: " + err.Error())
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
