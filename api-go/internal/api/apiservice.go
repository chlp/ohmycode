package api

import (
	"context"
	"encoding/json"
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

	mux.HandleFunc("/run/add_task", s.HandleRunAddTaskRequest)
	mux.HandleFunc("/run/get_tasks", s.HandleRunGetTasksRequest)
	mux.HandleFunc("/run/set_task_received", s.HandleRunSetTaskReceivedRequest)

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
