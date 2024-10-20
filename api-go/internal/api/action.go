package api

import (
	"encoding/json"
	"io"
	"net/http"
	"ohmycode_api/pkg/util"
)

func errorHandler(w http.ResponseWriter, str string, code int) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": str})
}

func handleAction(w http.ResponseWriter, r *http.Request) *input {
	if r.Method != http.MethodPost {
		errorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorHandler(w, "Unable to read input body", http.StatusInternalServerError)
		return nil
	}

	var i input
	err = json.Unmarshal(body, &i)
	if err != nil {
		errorHandler(w, "Invalid JSON input", http.StatusBadRequest)
		return nil
	}

	if !util.IsUuid(i.Session) {
		errorHandler(w, "Invalid: session", http.StatusBadRequest)
		return nil
	}

	if !util.IsUuid(i.User) {
		errorHandler(w, "Invalid: user", http.StatusBadRequest)
		return nil
	}

	return &i
}
