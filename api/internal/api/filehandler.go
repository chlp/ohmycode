package api

import (
	"encoding/json"
	"net/http"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"time"
)

const timeToSleepUntilNextFileUpdateSending = 500 * time.Millisecond

// saveVersionBeforeChange snapshots the current file content and saves it as a history
// version (using its ContentUpdatedAt as the version date) before the content is overwritten.
func (s *Service) saveVersionBeforeChange(file interface {
	Snapshot(includeContent bool) model.FileSnapshot
	SetVersionedAt(t time.Time)
}) {
	if s.versionStore == nil {
		return
	}
	snap := file.Snapshot(true)
	if snap.Content == nil || snap.ContentUpdatedAt.IsZero() {
		return
	}
	newVersionedAt, err := s.versionStore.SaveVersion(snap.ID, *snap.Content, snap.Name, snap.Lang, snap.ContentUpdatedAt, snap.VersionedAt)
	if err != nil {
		util.LogError("save version failed", "file_id", snap.ID, "error", err)
		return
	}
	if !newVersionedAt.IsZero() {
		file.SetVersionedAt(newVersionedAt)
	}
}

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
			if !isIgnorableWsErr(err) {
				util.Log("fileWork: send file error: " + err.Error())
			}
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
					if !isIgnorableWsErr(err) {
						util.Log("fileWork: send file error: " + err.Error())
					}
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
				if !isIgnorableWsErr(err) {
					util.Log("fileWork: send file error: " + err.Error())
				}
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
		if !util.IsValidId(i.FileId) || !util.IsValidId(i.UserId) {
			util.LogError("init: invalid file_id or user_id", "file_id", i.FileId, "user_id", i.UserId)
			return false
		}
		file, err := s.fileStore.GetFileOrCreate(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
		if err != nil {
			util.LogError("init: get file failed", "file_id", i.FileId, "error", err)
			return false
		}
		if file == nil {
			util.LogError("init: file not found", "file_id", i.FileId)
			return false
		}
		client.setFile(file, i.AppId, i.UserId)
		return true
	}

	if client.getFile() == nil {
		util.LogDebug("message before init", "action", i.Action, "file_id", i.FileId)
		return true
	}

	file := client.getFile()
	switch i.Action {
	case "set_content":
		s.saveVersionBeforeChange(file)
		if err := file.SetContent(i.Content, client.getAppId()); err != nil {
			util.LogError("set_content failed", "file_id", file.ID, "error", err)
		}
	case "set_name":
		if !file.SetName(i.FileName) {
			util.LogError("set_name failed", "file_id", file.ID)
		}
	case "set_user_name":
		if !file.SetUserName(client.getUserId(), i.UserName) {
			util.LogError("set_user_name failed", "file_id", file.ID, "user_id", client.getUserId())
		}
	case "set_lang":
		if !file.SetLang(i.Lang) {
			util.LogError("set_lang failed", "file_id", file.ID, "lang", i.Lang)
		}
	case "set_runner":
		if !file.SetRunnerId(i.RunnerId) {
			util.LogError("set_runner failed", "file_id", file.ID, "runner_id", i.RunnerId)
		}
	case "clean_result":
		s.taskStore.DeleteTask(file.ID)
		err = file.SetResult("")
		if err != nil {
			util.LogError("clean_result failed", "file_id", file.ID, "error", err)
		}
	case "run_task":
		snap := file.Snapshot(false)
		if snap.IsWaitingForResult {
			return true
		}
		if !s.runnerStore.IsOnline(snap.UsePublicRunner, snap.RunnerId) {
			return true
		}
		if !s.runLimiter.Allow(client.ip) {
			return true
		}
		file.SetWaitingForResult()
		s.taskStore.AddTask(file)
	case "run_task_with_content":
		snap := file.Snapshot(false)
		if snap.IsWaitingForResult {
			return true
		}
		if !s.runnerStore.IsOnline(snap.UsePublicRunner, snap.RunnerId) {
			return true
		}
		if !s.runLimiter.Allow(client.ip) {
			return true
		}
		s.saveVersionBeforeChange(file)
		if err := file.StartRunWithContent(i.Content, client.getAppId()); err != nil {
			util.LogError("start run failed", "file_id", file.ID, "error", err)
			return true
		}
		s.taskStore.AddTask(file)
	case "get_versions":
		versions, err := s.versionStore.GetVersions(file.ID)
		if err != nil {
			util.LogError("get_versions failed", "file_id", file.ID, "error", err)
			return true
		}
		dtos := make([]versionDTO, 0, len(versions))
		for _, v := range versions {
			dtos = append(dtos, versionDTO{
				ID:        v.ID,
				Name:      v.Name,
				Lang:      v.Lang,
				Preview:   v.Preview,
				CreatedAt: v.CreatedAt,
			})
		}
		resp := versionsResponseDTO{
			Action:   "versions",
			Versions: dtos,
		}
		if err := client.send(resp); err != nil {
			if !isIgnorableWsErr(err) {
				util.LogError("send versions failed", "file_id", file.ID, "error", err)
			}
			return false
		}
	case "restore_version":
		if i.VersionId == "" {
			util.LogError("restore_version: empty version_id", "file_id", file.ID)
			return true
		}
		version, err := s.versionStore.GetVersion(i.VersionId)
		if err != nil {
			util.LogError("restore_version failed", "file_id", file.ID, "version_id", i.VersionId, "error", err)
			return true
		}
		if version == nil {
			util.LogError("restore_version: not found", "file_id", file.ID, "version_id", i.VersionId)
			return true
		}
		if version.FileID != file.ID {
			util.LogError("restore_version: belongs to another file", "file_id", file.ID, "version_id", i.VersionId)
			return true
		}
		// Create new file with content from version
		newFileId := util.GenId()
		newFile, err := s.fileStore.GetFileOrCreate(newFileId, version.Name, version.Lang, version.Content, client.getUserId(), "")
		if err != nil {
			util.LogError("restore_version: create file failed", "file_id", file.ID, "error", err)
			return true
		}
		// Send redirect response to client
		resp := openFileResponseDTO{
			Action: "open_file",
			FileId: newFile.ID,
		}
		if err := client.send(resp); err != nil {
			if !isIgnorableWsErr(err) {
				util.LogError("send open_file failed", "file_id", file.ID, "error", err)
			}
			return false
		}
	default:
		util.LogDebug("unknown action", "action", i.Action, "file_id", file.ID)
	}
	return true
}
