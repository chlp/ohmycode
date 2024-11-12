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
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					responseErr(r.Context(), w, "Ping error: "+err.Error(), http.StatusConflict)
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
				wsMessageType, message, err := conn.ReadMessage()
				if err != nil {
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
					util.Log(r.Context(), "Cannot unmarshal: "+string(message))
					continue
				}
				switch i.Action {
				case "init":
					println("init")
					client.userId = i.UserId
					file, err := s.fileStore.GetFile(i.FileId)
					if err != nil {
						// todo: send error msg
						responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
						return
					}
					if file == nil {
						file = model.NewFile(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
						if file.UsePublicRunner {
							file.IsRunnerOnline = s.runnerStore.IsOnline(true, "")
						}
					}
					client.file = file
					client.lastUpdate = time.Now()
					if err = client.write(file); err != nil {
						responseErr(r.Context(), w, err.Error(), http.StatusBadRequest)
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
			// todo: send file updates
			println("send file updates")
			time.Sleep(1 * time.Second)

			// todo: if not initialized for 10 seconds -> close
		}
	}
}

func (client *fileWsClient) write(v interface{}) error {
	if v == nil {
		util.Log(context.Background(), "write nil")
		return nil
	}

	jsonData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if bytes.Equal(jsonData, []byte("null")) {
		util.Log(context.Background(), "write null")
		return nil
	}

	_ = client.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return client.conn.WriteMessage(websocket.TextMessage, jsonData)
}
