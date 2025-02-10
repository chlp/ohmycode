package api

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"ohmycode_api/pkg/util"
	"time"
)

type input struct {
	Action   string `json:"action"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	AppId    string `json:"app_id"`
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Content  string `json:"content"`
	Hash     uint32 `json:"hash"`
	Lang     string `json:"lang"`
	RunnerId string `json:"runner_id"`
	Result   string `json:"result"`
	IsPublic bool   `json:"is_public"`
}

const timeToSleepBetweenWork = 100 * time.Millisecond

func (s *Service) HandleWs(w http.ResponseWriter, r *http.Request,
	messageHandler func(client *wsClient, message []byte) (ok bool),
	work func(client *wsClient) (ok bool)) {
	client := createWsClient(w, r)
	defer client.closeDone()

	go func() {
		for {
			select {
			case <-client.done:
				return
			default:
				wsMessageType, message, err := client.conn.ReadMessage()
				if err != nil {
					var closeErr *websocket.CloseError
					if !errors.As(err, &closeErr) {
						util.Log("websocket conn.ReadMessage err: " + err.Error())
					}
					client.closeDone()
					return
				}
				if wsMessageType == websocket.CloseMessage {
					client.closeDone()
					break
				}
				if wsMessageType != websocket.TextMessage {
					continue
				}
				if !messageHandler(client, message) {
					client.closeDone()
					return
				}
			}
		}
	}()

	for {
		select {
		case <-client.done:
			return
		default:
			if !work(client) {
				return
			}
			time.Sleep(timeToSleepBetweenWork)
		}
	}
}
