package api

import (
	"encoding/json"
	"io"
	"net/http"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"time"
)

type input struct {
	FileId      string      `json:"file_id"`
	FileName    string      `json:"file_name"`
	UserId      string      `json:"user_id"`
	UserName    string      `json:"user_name"`
	Content     string      `json:"content"`
	Hash        uint32      `json:"hash"`
	Lang        string      `json:"lang"`
	RunnerId    string      `json:"runner_id"`
	Result      string      `json:"result"`
	IsPublic    bool        `json:"is_public"`
	IsKeepAlive bool        `json:"is_keep_alive"`
	LastUpdate  util.OhTime `json:"last_update"`
}

const keepAliveRequestMaxDuration = 30 * time.Second

func getInput(w http.ResponseWriter, r *http.Request) *input {
	if r.Method != http.MethodPost {
		responseErr(r.Context(), w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responseErr(r.Context(), w, "Unable to read input body", http.StatusInternalServerError)
		return nil
	}

	var i input
	err = json.Unmarshal(body, &i)
	if err != nil {
		responseErr(r.Context(), w, "Invalid JSON input", http.StatusBadRequest)
		return nil
	}

	return &i
}

func getInputForFile(w http.ResponseWriter, r *http.Request) *input {
	i := getInput(w, r)
	if i == nil {
		return nil
	}

	if !util.IsUuid(i.FileId) {
		responseErr(r.Context(), w, "Invalid: file id is not uuid", http.StatusBadRequest)
		return nil
	}

	if !util.IsUuid(i.UserId) {
		responseErr(r.Context(), w, "Invalid: user id is not uuid", http.StatusBadRequest)
		return nil
	}

	return i
}

func (s *Service) getFileOrCreateHandler(w http.ResponseWriter, r *http.Request) (*input, *model.File) {
	i := getInputForFile(w, r)
	if i == nil {
		return nil, nil
	}

	file, err := s.fileStore.GetFileOrCreate(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
	if err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return nil, nil
	} else if file == nil {
		responseErr(r.Context(), w, "Wrong file", http.StatusNotFound)
		return nil, nil
	}

	file.TouchByUser(i.UserId, i.UserName)

	return i, file
}
