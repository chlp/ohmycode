package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"sync"
	"time"
)

type fileWsClient struct {
	isInitialized bool
	file          *model.File
	userId        string
	lastUpdate    time.Time
	conn          *websocket.Conn
	done          chan struct{}
	closeDone     func()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true // could be return r.Host == "yourdomain.com"
	},
}

const WsMessageLimit = 4 * (1 << 20) // 4 Mb

func (s *Service) HandleWsFile(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		util.Log("HandleWsFile Upgrade error: " + err.Error())
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	var once sync.Once
	closeDone := func() {
		once.Do(func() {
			close(done)
		})
	}
	defer closeDone()

	conn.SetReadLimit(WsMessageLimit)

	client := fileWsClient{
		isInitialized: false,
		conn:          conn,
		done:          done,
		closeDone:     closeDone,
	}

	go client.pingPongHandling()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				wsMessageType, message, err := conn.ReadMessage()
				if err != nil {
					var closeErr *websocket.CloseError
					if !errors.As(err, &closeErr) {
						util.Log("websocket conn.ReadMessage err: " + err.Error())
					}
					closeDone()
					return
				}
				if wsMessageType == websocket.CloseMessage {
					closeDone()
					break
				}
				if wsMessageType != websocket.TextMessage {
					continue
				}

				var i input
				err = json.Unmarshal(message, &i)
				if err != nil {
					util.Log("Cannot unmarshal: " + string(message))
					continue
				}

				switch i.Action {
				case "init":
					client.userId = i.UserId
					client.file, err = s.fileStore.GetFile(i.FileId)
					if err != nil {
						util.Log("GetFile error: " + err.Error())
						closeDone()
						return
					}
					if client.file == nil {
						client.file = model.NewFile(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
					}
				case "set_content":
				case "set_name":
				case "set_user_name":
				case "set_lang":
				case "set_runner":
				default:
					util.Log("Unknown message type: " + string(message))
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
			if client.file != nil {
				client.file.TouchByUser(client.userId, "")
				if client.file.UpdatedAt.After(client.lastUpdate) {
					client.lastUpdate = time.Now()
					if err = client.sendFile(); err != nil {
						return
					}
					time.Sleep(400 * time.Millisecond)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (client *fileWsClient) sendFile() error {
	if client.file == nil {
		util.Log("sendFile nil")
		return nil
	}

	jsonData, err := json.Marshal(client.file)
	if err != nil {
		util.Log("sendFile json err: " + err.Error())
		return err
	}

	if bytes.Equal(jsonData, []byte("null")) {
		util.Log("sendFile null")
		return nil
	}

	return client.conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (client *fileWsClient) pingPongHandling() {
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
