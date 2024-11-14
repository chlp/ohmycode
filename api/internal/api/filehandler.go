package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"ohmycode_api/pkg/util"
	"time"
)

const timeToSleepUntilNextFileUpdateSending = 500 * time.Millisecond

func (s *Service) HandleWsFile(w http.ResponseWriter, r *http.Request) {
	s.HandleWs(w, r, s.fileMessageHandler, s.fileWork)
}

func (s *Service) fileWork(client *wsClient) (ok bool) {
	if client.file == nil {
		return true
	}
	client.file.TouchByUser(client.userId, "")
	if client.file.UpdatedAt.After(client.lastUpdate) {
		client.lastUpdate = time.Now()
		if err := client.sendFile(); err != nil {
			util.Log("send file error: " + err.Error())
			return false
		}
		time.Sleep(timeToSleepUntilNextFileUpdateSending)
	}
	return true
}

func (s *Service) fileMessageHandler(client *wsClient, message []byte) (ok bool) {
	var i input
	err := json.Unmarshal(message, &i)
	if err != nil {
		util.Log("Cannot unmarshal: " + string(message))
		return true
	}

	if i.Action == "init" {
		client.userId = i.UserId
		client.file, err = s.fileStore.GetFileOrCreate(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
		if err != nil {
			util.Log("GetFile error: " + err.Error())
			return false
		}
		if client.file == nil {
			util.Log("GetFile not found")
			return false
		}
		return true
	}

	if client.file == nil {
		return true
	}

	switch i.Action {
	case "set_content":
		if err := client.file.SetContent(i.Content, i.UserId); err != nil {
			util.Log("set_content error: " + err.Error())
		}
	case "set_name":
		if !client.file.SetName(i.FileName) {
			util.Log("set_name error")
		}
	case "set_user_name":
		if !client.file.SetUserName(client.userId, i.FileName) {
			util.Log("set_user_name error")
		}
	case "set_lang":
		if !client.file.SetLang(i.Lang) {
			util.Log("set_lang error")
		}
	case "set_runner":
		if !client.file.SetRunnerId(i.RunnerId) {
			util.Log("set_runner error")
		}
	case "clean_result":
		s.taskStore.DeleteTask(client.file.ID)
		err = client.file.SetResult("")
		if err != nil {
			util.Log("set_runner error")
		}
	case "run_task":
		if !s.runnerStore.IsOnline(client.file.UsePublicRunner, client.file.RunnerId) {
			return true
		} else {
			client.file.SetWaitingForResult()
			s.taskStore.AddTask(client.file)
		}
	default:
		util.Log("Unknown message type: " + string(message))
	}
	return true
}

func (client *wsClient) sendFile() error {
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
