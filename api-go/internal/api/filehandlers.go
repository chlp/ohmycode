package api

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Service) HandleFileGetUpdateRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r) // todo: check if it will create empty file
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
