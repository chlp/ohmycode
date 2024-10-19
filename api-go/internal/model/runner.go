package model

import (
	"time"
)

type Runner struct {
	ID        string    `json:"_id,omitempty"`
	IsPublic  bool      `json:"is_public"`
	CheckedAt time.Time `json:"checked_at"`
}
