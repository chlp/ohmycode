package model

import "time"

type Task struct {
	FileId   string `json:"-"`
	Content  string `json:"content"`
	Lang     string `json:"lang"`
	Hash     uint32 `json:"hash"`
	RunnerId string `json:"-"`
	IsPublic bool   `json:"-"`

	GivenToRunnerAt time.Time
}
