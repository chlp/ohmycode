package api

import (
	"net/http"
)

func (s *Service) HandleFileSetContentRequest(w http.ResponseWriter, r *http.Request) {
	i, file := s.getFileOrCreateHandler(w, r)
	if file == nil {
		return
	}

	if err := file.SetContent(i.Content, i.UserId); err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusBadRequest)
	} else {
		responseOk(w, nil)
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
