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
	file       *model.File
	userId     string
	appId      string
	runner     *model.Runner
	lastUpdate time.Time
	conn       *websocket.Conn
	done       chan struct{}
	close      func()
	writeMu    sync.Mutex
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
