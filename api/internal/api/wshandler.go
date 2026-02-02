package api

import (
	"net/http"
	"ohmycode_api/pkg/util"

	"github.com/gorilla/websocket"
)

type input struct {
	Action    string `json:"action"`
	FileId    string `json:"file_id"`
	FileName  string `json:"file_name"`
	AppId     string `json:"app_id"`
	UserId    string `json:"user_id"`
	UserName  string `json:"user_name"`
	Content   string `json:"content"`
	Hash      uint32 `json:"hash"`
	Lang      string `json:"lang"`
	RunnerId  string `json:"runner_id"`
	Result    string `json:"result"`
	IsPublic  bool   `json:"is_public"`
	VersionId string `json:"version_id"`
}

func (s *Service) HandleWs(w http.ResponseWriter, r *http.Request,
	messageHandler func(client *wsClient, message []byte) (ok bool),
	work func(client *wsClient) (ok bool)) {
	client := createWsClient(w, r, s.wsAllowedOrigins)
	if client == nil {
		return
	}
	defer client.close()

	go func() {
		for {
			select {
			case <-client.done:
				return
			default:
				wsMessageType, message, err := client.conn.ReadMessage()
				if err != nil {
					if !isIgnorableWsErr(err) {
						util.Log("websocket conn.ReadMessage err: " + err.Error())
					}
					client.close()
					return
				}
				if wsMessageType == websocket.CloseMessage {
					client.close()
					break
				}
				if wsMessageType != websocket.TextMessage {
					continue
				}
				if !messageHandler(client, message) {
					client.close()
					return
				}
			}
		}
	}()

	if !work(client) {
		client.close()
	}
}
