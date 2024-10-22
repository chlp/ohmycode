package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"ohmycode_api/internal/store"
	"ohmycode_api/pkg/util"
	"strconv"
)

type Service struct {
	store store.Store
}

func NewService(store store.Store) Service {
	return Service{
		store: store,
	}
}

//
//func (s *Service) GetNewestPublicRunnerCheckedAt() *time.Time {
//	var t *time.Time
//	for _, runner := range s.runners {
//		if runner.IsPublic {
//			if t == nil || t.Before(runner.CheckedAt) {
//				t = &runner.CheckedAt
//			}
//		}
//	}
//	return t
//}

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

	if !util.IsUuid(i.SessionId) {
		responseErr(r.Context(), w, "Invalid: session", http.StatusBadRequest)
		return nil
	}

	if !util.IsUuid(i.UserId) {
		responseErr(r.Context(), w, "Invalid: user", http.StatusBadRequest)
		return nil
	}

	return &i
}
