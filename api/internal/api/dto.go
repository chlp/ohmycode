package api

import (
	"ohmycode_api/internal/model"
	"time"
)

// fileDTO is a JSON snapshot of model.File without sync primitives.
// It intentionally keeps the JSON shape expected by the frontend.
type fileDTO struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Lang             string      `json:"lang"`
	Content          *string     `json:"content,omitempty"`
	ContentUpdatedAt time.Time   `json:"content_updated_at"`
	Result           string      `json:"result"`
	WriterID         string      `json:"writer_id"`
	Runner           string      `json:"runner"`
	Users            []model.User `json:"users"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Persisted        bool        `json:"persisted"`
	IsWaitingForResult bool      `json:"is_waiting_for_result"`
	IsRunnerOnline     bool      `json:"is_runner_online"`
}

func toFileDTO(f *model.File, includeContent bool) fileDTO {
	s := f.Snapshot(includeContent)
	dto := fileDTO{
		ID:                 s.ID,
		Name:               s.Name,
		Lang:               s.Lang,
		ContentUpdatedAt:   s.ContentUpdatedAt,
		Result:             s.Result,
		WriterID:           s.Writer,
		Runner:             s.RunnerId,
		Users:              s.Users,
		UpdatedAt:          s.UpdatedAt,
		Persisted:          s.Persisted,
		IsWaitingForResult: s.IsWaitingForResult,
		IsRunnerOnline:     s.IsRunnerOnline,
	}
	if includeContent {
		dto.Content = s.Content
	}
	return dto
}

type versionDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Lang      string    `json:"lang"`
	CreatedAt time.Time `json:"created_at"`
}

type versionsResponseDTO struct {
	Action   string       `json:"action"`
	Versions []versionDTO `json:"versions"`
}

type openFileResponseDTO struct {
	Action string `json:"action"`
	FileId string `json:"file_id"`
}


