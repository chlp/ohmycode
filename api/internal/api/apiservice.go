package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"os"
	"strconv"
	"time"
)

type Service struct {
	httpPort         int
	serveClientFiles bool
	fileStore        *store.FileStore
	runnerStore      *store.RunnerStore
	taskStore        *store.TaskStore
}

func NewService(httpPort int, serveClientFiles bool, fileStore *store.FileStore, runnerStore *store.RunnerStore, taskStore *store.TaskStore) *Service {
	return &Service{
		httpPort:         httpPort,
		serveClientFiles: serveClientFiles,
		fileStore:        fileStore,
		runnerStore:      runnerStore,
		taskStore:        taskStore,
	}
}

func (s *Service) Run() {
	util.Log("API Service started")
	mux := http.NewServeMux()

	mux.HandleFunc("/file/set_name", s.HandleFileSetNameRequest)
	mux.HandleFunc("/file/set_user_name", s.HandleFileSetUserNameRequest)
	mux.HandleFunc("/file/set_lang", s.HandleFileSetLangRequest)
	mux.HandleFunc("/file/set_runner", s.HandleFileSetRunnerRequest)

	mux.HandleFunc("/run/add_task", s.HandleRunAddTaskRequest)
	mux.HandleFunc("/run/get_tasks", s.HandleRunGetTasksRequest)

	mux.HandleFunc("/result/set", s.HandleResultSetRequest)
	mux.HandleFunc("/result/clean", s.HandleResultCleanRequest)

	mux.HandleFunc("/file", s.HandleWsFile)

	if s.serveClientFiles {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" && r.URL.Path != "/index.html" {
				file := "./static" + r.URL.Path
				if _, err := os.Stat(file); err == nil {
					http.ServeFile(w, r, file)
					return
				}
			}
			http.ServeFile(w, r, "./static/index.html")
			return
		})
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.httpPort), corsMiddleware(timerMiddleware(mux))))
}

func timerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), util.RequestStartTimeCtxKey, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func responseErr(ctx context.Context, w http.ResponseWriter, str string, code int) {
	util.Log(ctx, "action.responseErr: "+strconv.Itoa(code)+" ("+http.StatusText(code)+"): "+str)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": str})
}

func responseOk(w http.ResponseWriter, v interface{}) {
	if v == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonData, err := json.Marshal(v)
	if err != nil {
		responseErr(context.Background(), w, err.Error(), http.StatusInternalServerError)
		return
	}

	if string(jsonData) == "null" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}
