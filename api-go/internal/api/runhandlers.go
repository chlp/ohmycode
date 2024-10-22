package api

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Service) HandleRunAddTaskRequest(w http.ResponseWriter, r *http.Request) {
	_, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}
	if !file.RunnerIsOnline() {
		responseErr(r.Context(), w, "Runner is not online", http.StatusBadRequest)
	}

	// todo: create task

	responseOk(w, nil)
}

func (s *Service) HandleRunGetTasksRequest(w http.ResponseWriter, r *http.Request) {
	i := getInput(w, r)
	if i == nil {
		return
	}

	tasks := make([]string, 0)
	startTime := time.Now()
	for {
		// receive all tasks from files and fill tasks. Mark not to give these tasks again for 3 seconds

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

func (s *Service) HandleRunSetTaskReceivedRequest(w http.ResponseWriter, r *http.Request) {
	i := getInput(w, r)
	if i == nil {
		return
	}

	// todo: find file with task, hash, lang, is_public||runner_id and mark as received
	// todo: if not found, return special status to remove task from runner

	//if err := file.SetContent(i.Content, i.UserId); err != nil {
	//	responseOk(w, nil)
	//} else {
	responseErr(r.Context(), w, "Not implemented", http.StatusNotImplemented)
	//}
}
