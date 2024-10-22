package api

import (
	"fmt"
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
	Lang        string      `json:"lang"`
	RunnerId    string      `json:"runner_id"`
	IsKeepAlive bool        `json:"is_keep_alive"`
	LastUpdate  util.OhTime `json:"last_update"`
}

const keepAliveRequestMaxDuration = 30 * time.Second

func (s *Service) HandleFileGetUpdateRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}

	if i.IsKeepAlive {
		startTime := time.Now()
		for {
			if file.UpdatedAt.After(i.LastUpdate.Time) {
				break
			}
			if time.Since(startTime) > keepAliveRequestMaxDuration {
				break
			}

			file.TouchByUser(i.UserId, "")

			select {
			case <-r.Context().Done():
				fmt.Println("Client connection closed")
				responseOk(w, "Client connection closed")
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}

		}
	}

	responseOk(w, file)
}

func (s *Service) HandleFileSetContentRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}

	if err := file.SetContent(i.Content, i.UserId); err != nil {
		responseOk(w, nil)
	} else {
		responseErr(r.Context(), w, err.Error(), http.StatusBadRequest)
	}
}

func (s *Service) HandleFileSetNameRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}
	if file.SetName(i.FileName) {
		responseOk(w, nil)
	} else {
		responseErr(r.Context(), w, "Wrong file name", http.StatusBadRequest)
	}
}

func (s *Service) HandleFileSetUserNameRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}
	if file.SetUserName(i.UserId, i.UserName) {
		responseOk(w, nil)
	} else {
		responseErr(r.Context(), w, "Wrong user name", http.StatusBadRequest)
	}
}

func (s *Service) HandleFileSetLangRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}
	if file.SetLang(i.Lang) {
		responseOk(w, nil)
	} else {
		responseErr(r.Context(), w, "Wrong user name", http.StatusBadRequest)
	}
}

func (s *Service) HandleFileSetRunnerRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}
	if file.SetRunnerId(i.RunnerId) {
		responseOk(w, nil)
	} else {
		responseErr(r.Context(), w, "Wrong user name", http.StatusBadRequest)
	}
}

func (s *Service) getFileOrCreateHandler(w http.ResponseWriter, r *http.Request) (*input, *model.File) {
	i := handleAction(w, r)
	if i == nil {
		return nil, nil
	}

	file, err := s.store.GetFileOrCreate(i.FileId, i.FileName, i.Lang, i.Content, i.UserId, i.UserName)
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
