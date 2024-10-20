package api

import (
	"fmt"
	"net/http"
	"ohmycode_api/pkg/util"
	"time"
)

type input struct {
	Session     string      `json:"session"`
	SessionName string      `json:"session_name"`
	User        string      `json:"user"`
	UserName    string      `json:"user_name"`
	Code        string      `json:"code"`
	Lang        string      `json:"lang"`
	IsKeepAlive bool        `json:"is_keep_alive"`
	LastUpdate  util.OhTime `json:"last_update"`
}

const keepAliveRequestMaxDuration = 30 * time.Second

func (s *Service) HandleFileGetUpdateRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	i := handleAction(w, r)
	if i == nil {
		return
	}

	file, err := s.store.GetFile(i.Session)
	if err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		// todo: return null stated file
		responseErr(r.Context(), w, "Wrong file", http.StatusNotFound)
		return
	}

	userFound := false
	for _, user := range file.Users {
		if user == i.User {
			userFound = true
			break
		}
	}
	if !userFound {
		println(1)
		// file.InsertUser(i.User)
	} else {
		println(2)
		// file.UpdateUserOnline(i.User)
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

			// file.UpdateUserOnline(i.User)
			// file.CleanupUsers() - send into file worker
			// file.CleanupWriter() - send into file worker

			select {
			case <-ctx.Done():
				fmt.Println("Client connection closed")
				responseOk(w, "Client connection closed")
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}

		}
	}

	// file.CleanupUsers() - send into file worker
	// file.CleanupWriter() - send into file worker

	//response, err := json.Marshal(file)
	//if err != nil {
	//	responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	responseOk(w, file)
}

func (s *Service) HandleFileSetCodeRequest(w http.ResponseWriter, r *http.Request) {
	i := handleAction(w, r)
	if i == nil {
		return
	}
	file, err := s.store.GetFile(i.Session)
	if err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		// todo: create new file
		responseErr(r.Context(), w, "Wrong file", http.StatusNotFound)
		return
	}

	userFound := false
	for _, user := range file.Users {
		if user == i.User {
			userFound = true
			break
		}
	}
	if !userFound {
		println(1)
		// file.InsertUser(i.User)
	} else {
		println(2)
		// file.UpdateUserOnline(i.User)
	}

	if err = file.SetCode(i.Code, i.User); err != nil {
		responseOk(w, nil)
	} else {
		responseErr(r.Context(), w, err.Error(), http.StatusBadRequest)
	}
}

//// getSession получает сессию или создает новую
//func getSession(sessionID, userID, userName, lang string) *Session {
//	session := &Session{ID: sessionID, Writer: userID, Lang: lang}
//	if userName != "" {
//		session.Users = append(session.Users, User{ID: userID, UserName: userName})
//	}
//	// Здесь можно добавить логику сохранения сессии
//	return session
//}
