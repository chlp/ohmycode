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
	runner     *model.Runner
	lastUpdate time.Time
	conn       *websocket.Conn
	done       chan struct{}
	closeDone  func()
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

func createWsClient(w http.ResponseWriter, r *http.Request) *wsClient {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		util.Log("WebSocket Upgrade error: " + err.Error())
		return nil
	}
	conn.SetReadLimit(wsMessageLimit)

	done := make(chan struct{})
	var once sync.Once
	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}

	conn.SetReadLimit(wsMessageLimit)

	client := wsClient{
		conn:      conn,
		done:      done,
		closeDone: closeDone,
	}
	go client.pingPongHandling()
	return &client
}

func (client *wsClient) pingPongHandling() {
	client.conn.SetPongHandler(func(appData string) error {
		if err := client.conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			util.Log("pong handler err: ", err.Error())
			client.closeDone()
		}
		return nil
	})
	for {
		select {
		case <-client.done:
			return
		default:
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				util.Log("Ping error: " + err.Error())
				return
			}
			time.Sleep(5 * time.Second)
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

	return client.conn.WriteMessage(websocket.TextMessage, jsonData)
}
