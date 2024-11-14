package api

import (
	"encoding/json"
	"net/http"
	"ohmycode_api/pkg/util"
	"time"
)

const timeToSleepUntilNextFileUpdateSending = 500 * time.Millisecond

func (s *Service) handleWsFileConnection(w http.ResponseWriter, r *http.Request) {
	s.HandleWs(w, r, s.fileMessageHandler, s.fileWork)
}

func (s *Service) fileWork(client *wsClient) (ok bool) {
	if client.file == nil {
		return true
	}
	client.file.TouchByUser(client.userId, "")
	if client.file.UpdatedAt.After(client.lastUpdate) {
		client.lastUpdate = time.Now()
		if err := client.send(client.file); err != nil {
			util.Log("fileWork: send file error: " + err.Error())
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
		util.Log("fileMessageHandler: Cannot unmarshal: " + string(message))
		return true
	}

	if i.Action == "init" {
		if !util.IsUuid(i.FileId) || !util.IsUuid(i.UserId) {
			util.Log("fileMessageHandler: Wrong file_id or user_id: " + i.FileId + ", " + i.UserId)
			return false
		}
		client.userId = i.UserId
		client.file, err = s.fileStore.GetFileOrCreate(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
		if err != nil {
			util.Log("fileMessageHandler: GetFile error: " + err.Error())
			return false
		}
		if client.file == nil {
			util.Log("fileMessageHandler: GetFile not found")
			return false
		}
		return true
	}

	if client.file == nil {
		util.Log("fileMessageHandler: nil file: " + i.RunnerId)
		return true
	}

	switch i.Action {
	case "set_content":
		if err := client.file.SetContent(i.Content, i.UserId); err != nil {
			util.Log("fileMessageHandler: set_content error: " + err.Error())
		}
	case "set_name":
		if !client.file.SetName(i.FileName) {
			util.Log("fileMessageHandler: set_name error")
		}
	case "set_user_name":
		if !client.file.SetUserName(client.userId, i.FileName) {
			util.Log("fileMessageHandler: set_user_name error")
		}
	case "set_lang":
		if !client.file.SetLang(i.Lang) {
			util.Log("fileMessageHandler: set_lang error")
		}
	case "set_runner":
		if !client.file.SetRunnerId(i.RunnerId) {
			util.Log("fileMessageHandler: set_runner error")
		}
	case "clean_result":
		s.taskStore.DeleteTask(client.file.ID)
		err = client.file.SetResult("")
		if err != nil {
			util.Log("fileMessageHandler: set_runner error")
		}
	case "run_task":
		if !s.runnerStore.IsOnline(client.file.UsePublicRunner, client.file.RunnerId) {
			return true
		} else {
			client.file.SetWaitingForResult()
			s.taskStore.AddTask(client.file)
		}
	default:
		util.Log("fileMessageHandler: Unknown message type: " + string(message))
	}
	return true
}
