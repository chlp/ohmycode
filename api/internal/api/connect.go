package api

import (
	"bytes"
	"context"
	"encoding/json"
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
		util.Log(r.Context(), "HandleWsFile Upgrade: "+err.Error())
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

	conn.SetPongHandler(func(appData string) error {
		util.Log("pong handler")
		if err = conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			responseErr(r.Context(), w, err.Error(), http.StatusBadRequest)
		}
		return nil
	})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				util.Log("ping")
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					util.Log(context.Background(), "Ping error: "+err.Error())
					return
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	client := fileWsClient{
		isInitialized: false,
		conn:          conn,
	}

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				util.Log("r?")
				wsMessageType, message, err := conn.ReadMessage()
				if err != nil {
					util.Log("errrr" + err.Error())
					closeDone()
					return
				}
				if wsMessageType == websocket.PongMessage {
					util.Log("r:pong")
					break
				}
				if wsMessageType == websocket.PingMessage {
					util.Log("r:ping")
					break
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
					util.Log(r.Context(), "Cannot unmarshal: "+string(message))
					continue
				}
				switch i.Action {
				case "init":
					println("init")
					client.userId = i.UserId
					client.file, err = s.fileStore.GetFile(i.FileId)
					if err != nil {
						// todo: send error msg
						responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
						return
					}
					if client.file == nil {
						client.file = model.NewFile(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
						if client.file.UsePublicRunner {
							client.file.IsRunnerOnline = s.runnerStore.IsOnline(true, "")
						}
					}
					client.lastUpdate = time.Now()
					if err = client.sendFile(); err != nil {
						closeDone()
						return
					}
				case "set_content":
				case "set_name":
				case "set_user_name":
				case "set_lang":
				case "set_runner":
				default:
					util.Log(r.Context(), "Unknown message type: "+string(message))
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		default:
			if client.file != nil && client.file.UpdatedAt.After(client.lastUpdate) {
				if err = client.sendFile(); err != nil {
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (client *fileWsClient) sendFile() error {
	if client.file == nil {
		util.Log(context.Background(), "sendFile nil")
		return nil
	}

	jsonData, err := json.Marshal(client.file)
	if err != nil {
		return err
	}

	if bytes.Equal(jsonData, []byte("null")) {
		util.Log(context.Background(), "sendFile null")
		return nil
	}

	return client.conn.WriteMessage(websocket.TextMessage, jsonData)
}
