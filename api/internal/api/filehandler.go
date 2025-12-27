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
	// Wait for init
	for client.getFile() == nil {
		select {
		case <-client.done:
			return true
		case <-client.fileSetCh:
		}
	}

	touchTicker := time.NewTicker(1 * time.Second)
	defer touchTicker.Stop()

	var nextSendAllowed time.Time
	var updatesCh <-chan struct{}
	var subFile interface {
		SubscribeUpdates() chan struct{}
		UnsubscribeUpdates(ch chan struct{})
	}
	var subCh chan struct{}
	defer func() {
		if subFile != nil && subCh != nil {
			subFile.UnsubscribeUpdates(subCh)
		}
	}()

	// Send initial snapshot eagerly
	if f := client.getFile(); f != nil {
		subFile = f
		subCh = f.SubscribeUpdates()
		updatesCh = subCh
		dto := toFileDTO(f, true)
		if err := client.send(dto); err != nil {
			util.Log("fileWork: send file error: " + err.Error())
			return false
		}
		// lastUpdate must track the file's UpdatedAt we have sent, not wall-clock send time.
		// Otherwise we can drop updates that occurred while we were sending the previous message.
		client.setLastUpdate(dto.UpdatedAt)
		nextSendAllowed = time.Now().Add(timeToSleepUntilNextFileUpdateSending)
	}

	for {
		select {
		case <-client.done:
			return true
		case <-client.fileSetCh:
			// File was changed (SPA navigation), send new snapshot immediately.
			if f := client.getFile(); f != nil {
				// re-subscribe to the new file updates
				if subFile != nil && subCh != nil {
					subFile.UnsubscribeUpdates(subCh)
				}
				subFile = f
				subCh = f.SubscribeUpdates()
				updatesCh = subCh
				dto := toFileDTO(f, true)
				if err := client.send(dto); err != nil {
					util.Log("fileWork: send file error: " + err.Error())
					return false
				}
				client.setLastUpdate(dto.UpdatedAt)
				nextSendAllowed = time.Now().Add(timeToSleepUntilNextFileUpdateSending)
			}
		case <-touchTicker.C:
			if f := client.getFile(); f != nil {
				f.TouchByUser(client.getUserId(), "")
			}
		case <-updatesCh:
			// throttle sends
			if time.Now().Before(nextSendAllowed) {
				timer := time.NewTimer(time.Until(nextSendAllowed))
				select {
				case <-client.done:
					timer.Stop()
					return true
				case <-timer.C:
				}
			}

			f := client.getFile()
			if f == nil {
				continue
			}
			lastUpdate := client.getLastUpdate()
			snapMeta := f.Snapshot(false)
			if !snapMeta.UpdatedAt.After(lastUpdate) {
				continue
			}

			includeContent := true
			if snapMeta.ContentUpdatedAt.Before(lastUpdate) {
				includeContent = false
			}
			dto := toFileDTO(f, includeContent)
			if snapMeta.ID != dto.ID {
				util.Log("fileWork: new file, will not send now")
				continue
			}
			if err := client.send(dto); err != nil {
				util.Log("fileWork: send file error: " + err.Error())
				return false
			}

			client.setLastUpdate(dto.UpdatedAt)
			nextSendAllowed = time.Now().Add(timeToSleepUntilNextFileUpdateSending)
		}
	}
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
		file, err := s.fileStore.GetFileOrCreate(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
		if err != nil {
			util.Log("fileMessageHandler: GetFile error: " + err.Error())
			return false
		}
		if file == nil {
			util.Log("fileMessageHandler: GetFile not found")
			return false
		}
		client.setFile(file, i.AppId, i.UserId)
		return true
	}

	if client.getFile() == nil {
		util.Log("fileMessageHandler: nil file: " + i.RunnerId)
		return true
	}

	file := client.getFile()
	switch i.Action {
	case "set_content":
		if err := file.SetContent(i.Content, client.getAppId()); err != nil {
			util.Log("fileMessageHandler: set_content error: " + err.Error())
		}
	case "set_name":
		if !file.SetName(i.FileName) {
			util.Log("fileMessageHandler: set_name error")
		}
	case "set_user_name":
		if !file.SetUserName(client.getUserId(), i.UserName) {
			util.Log("fileMessageHandler: set_user_name error")
		}
	case "set_lang":
		if !file.SetLang(i.Lang) {
			util.Log("fileMessageHandler: set_lang error")
		}
	case "set_runner":
		if !file.SetRunnerId(i.RunnerId) {
			util.Log("fileMessageHandler: set_runner error")
		}
	case "clean_result":
		s.taskStore.DeleteTask(file.ID)
		err = file.SetResult("")
		if err != nil {
			util.Log("fileMessageHandler: set_runner error")
		}
	case "run_task":
		if !s.runnerStore.IsOnline(file.UsePublicRunner, file.RunnerId) {
			return true
		} else {
			file.SetWaitingForResult()
			s.taskStore.AddTask(file)
		}
	default:
		util.Log("fileMessageHandler: Unknown message type: " + string(message))
	}
	return true
}
