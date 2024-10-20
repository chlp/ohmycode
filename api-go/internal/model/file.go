package model

import (
	"ohmycode_api/pkg/util"
	"time"
)

type File struct {
	ID              string    `json:"_id,omitempty"`
	Name            string    `json:"name"`
	Lang            string    `json:"lang"`
	Code            string    `json:"code"`
	Writer          string    `json:"writer"`
	Runner          string    `json:"runner"`
	UpdatedAt       time.Time `json:"updated_at"`
	CodeUpdatedAt   time.Time `json:"code_updated_at"`
	RunnerCheckedAt time.Time `json:"runner_checked_at"`
	Users           []string  `json:"users"` // todo: user
}

const (
	DefaultLang                    = "markdown"
	codeMaxLength                  = 32768
	durationIsActiveFromLastUpdate = 5 * time.Second
	durationIsWriterStillWriting   = 2 * time.Second
)

type Lang struct {
	Name        string // Название языка
	Highlighter string // Соответствующий хайлайтер
}

type Langs map[string]Lang

var LANGS = Langs{
	"go": {
		Name:        "GoLang",
		Highlighter: "go",
	},
	"java": {
		Name:        "Java",
		Highlighter: "text/x-java",
	},
	"json": {
		Name:        "JSON",
		Highlighter: "application/json",
	},
	"markdown": {
		Name:        "Markdown",
		Highlighter: "text/x-markdown",
	},
	"mysql8": {
		Name:        "MySQL 8",
		Highlighter: "sql",
	},
	"php82": {
		Name:        "PHP 8.2",
		Highlighter: "php",
	},
	"postgres13": {
		Name:        "PostgreSQL 13",
		Highlighter: "sql",
	},
}

func (s *File) SetName(name string) bool {
	if !util.IsValidString(name) {
		return false
	}
	s.Name = name
	// todo: to update
	return true
}

func (s *File) SetLang(lang string) bool {
	if _, ok := LANGS[lang]; !ok {
		return false
	}
	s.Lang = lang
	// todo: to update
	return true
}

func (s *File) SetCode(code, userId string) bool {
	if len(code) > codeMaxLength {
		return false
	}
	// validate code
	s.Code = code
	// todo: to update
	return true
}

func (s *File) UpdateTime() {
	s.UpdatedAt = time.Now()
}

func (s *File) SetUserName(userId, name string) bool {
	if !util.IsValidString(name) || !util.IsUuid(userId) {
		return false
	}
	s.Name = name
	// todo: to update
	return false
}

func (s *File) SetRunner(runner string) bool {
	if !util.IsUuid(runner) {
		return false
	}
	s.Runner = runner
	// todo: to update
	return false
}

func (s *File) UpdateRunnerCheckedAt(runner string, isPublic bool) {
	if !util.IsUuid(runner) {
		return
	}
	s.RunnerCheckedAt = time.Now()
	// todo: to update
}

func (s *File) CleanupWriter() {
	if s.Writer == "" {
		return
	}
	if time.Since(s.RunnerCheckedAt) > durationIsWriterStillWriting {
		s.Writer = ""
		// todo: to update
	}
}

func (s *File) RunnerIsOnline() bool {
	if s.RunnerCheckedAt.IsZero() {
		return false
	}
	return time.Since(s.RunnerCheckedAt) < durationIsActiveFromLastUpdate
}
