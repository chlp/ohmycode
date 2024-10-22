package model

import (
	"errors"
	"ohmycode_api/pkg/util"
	"sync"
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
	mutex            *sync.Mutex
}

type User struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	TouchedAt time.Time `json:"touched_at" bson:"touched_at"`
}

const (
	contentMaxLength               = 32768
	durationIsActiveFromLastUpdate = 5 * time.Second
	durationIsWriterStillWriting   = 2 * time.Second
)

func (f *File) TouchByUser(userId, userName string) {
	if !util.IsUuid(userId) {
		return
	}

	f.lock()
	defer f.unlock()

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
	if _, ok := Langs[lang]; !ok {
		return false
	}
	f.Lang = lang
	// todo: to update
	return true
}

func (f *File) SetCode(code, userId string) error {
	if len(code) > contentMaxLength {
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

func (f *File) CleanupUsers() {
	f.lock()
	defer f.unlock()

	for i, user := range f.Users {
		if time.Since(user.TouchedAt) > durationIsActiveFromLastUpdate {
			f.Users = append(f.Users[:i], f.Users[i+1:]...)
		}
	}
}

func (f *File) CleanupWriter() {
	if f.Writer == "" {
		return
	}
	if time.Since(f.ContentUpdatedAt) > durationIsWriterStillWriting {
		f.Writer = ""
	}
}

func (f *File) RunnerIsOnline() bool {
	if f.RunnerCheckedAt.IsZero() {
		return false
	}
	return time.Since(f.RunnerCheckedAt) < durationIsActiveFromLastUpdate
}

func (f *File) lock() {
	if f.mutex == nil {
		f.mutex = &sync.Mutex{}
	}
	f.mutex.Lock()
}

func (f *File) unlock() {
	f.mutex.Lock()
}
