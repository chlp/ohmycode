package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ohmycode_api/pkg/util"
	"time"
)

type Session struct {
	ID     string `json:"id"`
	Users  []User `json:"users"`
	Writer string `json:"writer"`
	Lang   string `json:"lang"`
}

type User struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
}

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

//
//func HandleSessionRequest(w http.ResponseWriter, r *http.Request) {
//	i := handleAction(w, r)
//	if i == nil {
//		return
//	}
//	switch i.Action {
//	case "getUpdate":
//		handleGetUpdate(w, i.Session, i.User, i.IsKeepAlive)
//	case "setSessionName":
//		handleSetSessionName(w, sessionID, userID, userName, lang, i)
//	// Добавьте другие действия...
//	default:
//		errorHandler(w, "Wrong action", http.StatusNotFound)
//	}
//}

const keepAliveRequestMaxDuration = 30 * time.Second

func (s *Service) HandleFileGetUpdateRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	i := handleAction(w, r)
	if i == nil {
		return
	}
	file, err := s.GetFile(i.Session)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		// todo: return null stated file
		errorHandler(w, "Wrong file", http.StatusNotFound)
		return
	}

	userFound := false

	if i.IsKeepAlive {
		for {
			if file.UpdatedAt.After(i.LastUpdate.Time) {
				break
			}

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
			// file.CleanupUsers()
			// file.CleanupWriter()

			select {
			case <-ctx.Done():
				fmt.Println("Client connection closed")
				errorHandler(w, "Client connection closed", http.StatusOK)
				return
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}

	if !userFound {
		for _, user := range file.Users {
			if user == i.User {
				userFound = true
				break
			}
		}
	}
	if !userFound {
		println(1)
		// file.InsertUser(i.User)
	} else {
		println(2)
		// file.UpdateUserOnline(i.User)
	}
	// file.CleanupUsers()
	// file.CleanupWriter()

	response, err := json.Marshal(file)
	if err != nil {
		errorHandler(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(response)
}

// getSession получает сессию или создает новую
func getSession(sessionID, userID, userName, lang string) *Session {
	session := &Session{ID: sessionID, Writer: userID, Lang: lang}
	if userName != "" {
		session.Users = append(session.Users, User{ID: userID, UserName: userName})
	}
	// Здесь можно добавить логику сохранения сессии
	return session
}
