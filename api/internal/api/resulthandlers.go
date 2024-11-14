package api

import (
	"net/http"
)

func (s *Service) HandleResultSetRequest(w http.ResponseWriter, r *http.Request) {
	i := getInput(w, r)
	if i == nil {
		return
	}

	task := s.taskStore.GetTask(i.RunnerId, i.Lang, i.Hash)
	if task == nil {
		responseErr(r.Context(), w, "Task not found (result)", http.StatusNotFound)
		return
	}
	file, err := s.fileStore.GetFile(task.FileId)
	if err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		responseErr(r.Context(), w, "File not found (result)", http.StatusNotFound)
		return
	}

	if i.Result == "" {
		i.Result = "_"
	}
	if err := file.SetResult(i.Result); err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusBadRequest)
		return
	}
	s.taskStore.DeleteTask(file.ID)

	responseOk(w, nil)
}
