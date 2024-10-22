package model

import (
	"errors"
	"ohmycode_api/pkg/util"
	"time"
)

type File struct {
	ID               string    `json:"id" bson:"_id,omitempty"`
	Name             string    `json:"name" bson:"name"`
	Lang             string    `json:"lang" bson:"lang"`
	Content          string    `json:"content" bson:"content"`
	Writer           string    `json:"writer_id" bson:"writer_id"`
	Runner           string    `json:"runner_id" bson:"runner_id"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	ContentUpdatedAt time.Time `json:"content_updated_at" bson:"content_updated_at"`
	RunnerCheckedAt  time.Time `json:"runner_checked_at" bson:"runner_checked_at"`
	Users            []User    `json:"users" bson:"users"`
}

type User struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	TouchedAt time.Time `json:"touched_at" bson:"touched_at"`
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

func (f *File) TouchByUser(userId, userName string) {
	if !util.IsUuid(userId) {
		return
	}
	isAlreadyExist := false
	for i, user := range f.Users {
		if user.ID == userId {
			isAlreadyExist = true
			f.Users[i].TouchedAt = time.Now()
			break
		}
	}
	if !isAlreadyExist {
		if !util.IsValidName(userName) {
			userName = util.RandomName()
		}
		f.Users = append(f.Users, User{
			ID:        userId,
			Name:      userName,
			TouchedAt: time.Now(),
		})
		f.UpdatedAt = time.Now()
	}
}

func (f *File) SetName(name string) bool {
	if !util.IsValidName(name) {
		return false
	}
	f.Name = name
	// todo: to update
	return true
}

func (f *File) SetLang(lang string) bool {
	if _, ok := LANGS[lang]; !ok {
		return false
	}
	f.Lang = lang
	// todo: to update
	return true
}

func (f *File) SetCode(code, userId string) error {
	if len(code) > codeMaxLength {
		return errors.New("code is too long")
	}
	// validate code
	f.Content = code
	// todo: to update
	// setWriter or err
	return nil
}

func (f *File) UpdateTime() {
	f.UpdatedAt = time.Now()
}

func (f *File) SetUserName(userId, name string) bool {
	if !util.IsValidName(name) || !util.IsUuid(userId) {
		return false
	}
	f.Name = name
	// todo: to update
	return false
}

func (f *File) SetRunner(runner string) bool {
	if !util.IsUuid(runner) {
		return false
	}
	f.Runner = runner
	// todo: to update
	return false
}

func (f *File) UpdateRunnerCheckedAt(runner string, isPublic bool) {
	if !util.IsUuid(runner) {
		return
	}
	f.RunnerCheckedAt = time.Now()
	// todo: to update
}

func (f *File) CleanupWriter() {
	if f.Writer == "" {
		return
	}
	if time.Since(f.RunnerCheckedAt) > durationIsWriterStillWriting {
		f.Writer = ""
		// todo: to update
	}
}

func (f *File) RunnerIsOnline() bool {
	if f.RunnerCheckedAt.IsZero() {
		return false
	}
	return time.Since(f.RunnerCheckedAt) < durationIsActiveFromLastUpdate
}
