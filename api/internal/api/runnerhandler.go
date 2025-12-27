package api

import (
	"encoding/json"
	"net/http"
	"ohmycode_api/pkg/util"
	"strconv"
	"time"
)

func (s *Service) handleWsRunnerConnection(w http.ResponseWriter, r *http.Request) {
	s.HandleWs(w, r, s.runnerMessageHandler, s.runnerWork)
}

func (s *Service) runnerWork(client *wsClient) (ok bool) {
	updatesCh, unsubscribe := s.taskStore.Subscribe()
	defer unsubscribe()

	// Wait for init
	for client.getRunner() == nil {
		select {
		case <-client.done:
			return true
		case <-client.runnerSetCh:
		}
	}

	heartbeat := time.NewTicker(2 * time.Second)
	defer heartbeat.Stop()

	// Send initial tasks eagerly (in case tasks existed before runner connected)
	if r := client.getRunner(); r != nil {
		s.runnerStore.TouchRunner(r.ID)
		tasks := s.taskStore.GetTasksForRunner(r)
		if len(tasks) != 0 {
			if err := client.send(tasks); err != nil {
				util.Log("runnerWork: send tasks error: " + err.Error())
				return false
			}
		}
	}

	for {
		select {
		case <-client.done:
			return true
		case <-heartbeat.C:
			if r := client.getRunner(); r != nil {
				s.runnerStore.TouchRunner(r.ID)
			}
		case <-updatesCh:
			r := client.getRunner()
			if r == nil {
				continue
			}
			s.runnerStore.TouchRunner(r.ID)
			tasks := s.taskStore.GetTasksForRunner(r)
			if len(tasks) == 0 {
				continue
			}
			if err := client.send(tasks); err != nil {
				util.Log("runnerWork: send tasks error: " + err.Error())
				return false
			}
		}
	}
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
		client.setRunner(s.runnerStore.SetRunner(i.RunnerId, i.IsPublic))
		return true
	}

	if client.getRunner() == nil {
		util.Log("runnerMessageHandler: nil runner: " + i.RunnerId)
		return true
	}

	runner := client.getRunner()
	switch i.Action {
	case "set_result":
		task := s.taskStore.GetTask(runner.ID, i.Lang, i.Hash)
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
