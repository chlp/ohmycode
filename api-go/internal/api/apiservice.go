package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"strconv"
	"time"
)

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Run() {
	util.Log(nil, "API Service started")
	mux := http.NewServeMux()
	mux.HandleFunc("/file/get", s.HandleFileGetUpdateRequest)
	mux.HandleFunc("/file/set_content", s.HandleFileSetContentRequest)
	mux.HandleFunc("/file/set_name", s.HandleFileSetNameRequest)
	mux.HandleFunc("/file/set_user_name", s.HandleFileSetContentRequest)
	mux.HandleFunc("/file/set_lang", s.HandleFileSetLangRequest)
	mux.HandleFunc("/file/set_runner", s.HandleFileSetRunnerRequest)
	log.Fatal(http.ListenAndServe(":8081", requestTimerMiddleware(mux)))
}

func requestTimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), util.RequestStartTimeCtxKey, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func responseErr(ctx context.Context, w http.ResponseWriter, str string, code int) {
	util.Log(ctx, "action.responseErr: "+strconv.Itoa(code)+" ("+http.StatusText(code)+"): "+str)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": str})
}

func responseOk(w http.ResponseWriter, v interface{}) {
	if v != nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(v)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleAction(w http.ResponseWriter, r *http.Request) *input {
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

	if !util.IsUuid(i.FileId) {
		responseErr(r.Context(), w, "Invalid: file id is not uuid", http.StatusBadRequest)
		return nil
	}

	if !util.IsUuid(i.UserId) {
		responseErr(r.Context(), w, "Invalid: user id is not uuid", http.StatusBadRequest)
		return nil
	}

	return &i
}
