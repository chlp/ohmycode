package api

import (
	"fmt"
	"net/http"
	"ohmycode_api/pkg/util"
	"time"
)

type input struct {
	SessionId   string      `json:"session_id"`
	SessionName string      `json:"session_name"`
	UserId      string      `json:"user_id"`
	UserName    string      `json:"user_name"`
	Content     string      `json:"content"`
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

	file, err := s.store.GetFile(i.SessionId)
	if err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		// todo: return null stated file
		responseErr(r.Context(), w, "Wrong file", http.StatusNotFound)
		return
	}

	file.TouchByUser(i.UserId, i.UserName)

	if i.IsKeepAlive {
		startTime := time.Now()
		for {
			if file.UpdatedAt.After(i.LastUpdate.Time) {
				break
			}
			if time.Since(startTime) > keepAliveRequestMaxDuration {
				break
			}

			// file.UpdateUserOnline(i.UserId)
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
	file, err := s.store.GetFile(i.SessionId)
	if err != nil {
		responseErr(r.Context(), w, err.Error(), http.StatusInternalServerError)
		return
	}
	if file == nil {
		// todo: create new file
		responseErr(r.Context(), w, "Wrong file", http.StatusNotFound)
		return
	}

	file.TouchByUser(i.UserId, i.UserName)

	if err = file.SetCode(i.Content, i.UserId); err != nil {
		responseOk(w, nil)
	} else {
		responseErr(r.Context(), w, err.Error(), http.StatusBadRequest)
	}
}

//// getSession получает сессию или создает новую
//func getSession(sessionID, userID, userName, lang string) *SessionId {
//	session := &SessionId{ID: sessionID, Writer: userID, Lang: lang}
//	if userName != "" {
//		session.Users = append(session.Users, UserId{ID: userID, UserName: userName})
//	}
//	// Здесь можно добавить логику сохранения сессии
//	return session
//}
