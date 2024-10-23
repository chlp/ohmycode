package api

import (
	"fmt"
	"net/http"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"time"
)

func (s *Service) HandleRunAddTaskRequest(w http.ResponseWriter, r *http.Request) {
	_, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}
	if s.runnerStore.IsOnline(file.UsePublicRunner, file.RunnerId) {
		responseErr(r.Context(), w, "Runner is not online", http.StatusBadRequest)
	}

	s.taskStore.SetTask(file)

	responseOk(w, nil)
}

func (s *Service) HandleRunGetTasksRequest(w http.ResponseWriter, r *http.Request) {
	i := getInput(w, r)
	if i == nil {
		return
	}

	if !util.IsUuid(i.RunnerId) {
		responseErr(r.Context(), w, "Invalid: runner id is not uuid", http.StatusBadRequest)
		return
	}

	tasks := make([]*model.Task, 0)
	startTime := time.Now()
	for {
		s.runnerStore.SetRunner(i.RunnerId, i.IsPublic)

		tasks = s.taskStore.GetTasksForRunner(i.RunnerId, i.IsPublic)

		if len(tasks) > 0 || !i.IsKeepAlive || time.Since(startTime) > keepAliveRequestMaxDuration {
			break
		}

		select {
		case <-r.Context().Done():
			fmt.Println("Client connection closed")
			responseOk(w, "Client connection closed")
			return
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}

	responseOk(w, tasks)
}

func (s *Service) HandleRunAckTaskRequest(w http.ResponseWriter, r *http.Request) {
	i := getInput(w, r)
	if i == nil {
		return
	}

	task := s.taskStore.GetTask(i.RunnerId, i.Lang, i.Hash)
	if task == nil {
		responseErr(r.Context(), w, "Task not found", http.StatusNotFound)
		return
	}
	task.AcknowledgedByRunnerAt = time.Now()

	responseOk(w, nil)
}
