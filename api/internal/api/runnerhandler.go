package api

import (
	"encoding/json"
	"net/http"
	"ohmycode_api/pkg/util"
	"strconv"
)

func (s *Service) handleWsRunnerConnection(w http.ResponseWriter, r *http.Request) {
	s.HandleWs(w, r, s.runnerMessageHandler, s.runnerWork)
}

func (s *Service) runnerWork(client *wsClient) (ok bool) {
	if client.runner == nil {
		return true
	}

	s.runnerStore.TouchRunner(client.runner.ID)
	tasks := s.taskStore.GetTasksForRunner(client.runner)
	if len(tasks) == 0 {
		return true
	}

	if err := client.send(tasks); err != nil {
		util.Log("runnerWork: send tasks error: " + err.Error())
		return false
	}
	return true
}

func (s *Service) runnerMessageHandler(client *wsClient, message []byte) (ok bool) {
	var i input
	err := json.Unmarshal(message, &i)
	if err != nil {
		util.Log("runnerMessageHandler: Cannot unmarshal: " + string(message))
		return true
	}

	if i.Action == "init" {
		if !util.IsUuid(i.RunnerId) {
			util.Log("runnerMessageHandler: Wrong runner_id: " + i.RunnerId)
			return false
		}
		client.runner = s.runnerStore.SetRunner(i.RunnerId, i.IsPublic)
		return true
	}

	if client.runner == nil {
		util.Log("runnerMessageHandler: nil runner: " + i.RunnerId)
		return true
	}

	switch i.Action {
	case "set_result":
		task := s.taskStore.GetTask(client.runner.ID, i.Lang, i.Hash)
		if task == nil {
			util.Log("runnerMessageHandler: task not found: " + i.Lang + ", " + strconv.Itoa(int(i.Hash)))
			return true
		}
		file, err := s.fileStore.GetFile(task.FileId)
		if err != nil {
			util.Log("runnerMessageHandler: file getting: " + i.FileId + ", err: " + err.Error())
			return true
		}
		if file == nil {
			util.Log("runnerMessageHandler: file not found: " + i.FileId)
			return true
		}
		if i.Result == "" {
			i.Result = "_"
		}
		if err := file.SetResult(i.Result); err != nil {
			util.Log("runnerMessageHandler: SetResult: " + i.FileId + ", err: " + err.Error())
			return true
		}
		s.taskStore.DeleteTask(file.ID)
	default:
		util.Log("runnerMessageHandler: Unknown message type: " + string(message))
	}
	return true
}
