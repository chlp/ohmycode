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
	ContentUpdatedAt time.Time `json:"content_updated_at" bson:"content_updated_at"`
	Result           string    `json:"result" bson:"result"`
	Writer           string    `json:"writer_id" bson:"writer_id"`
	UsePublicRunner  bool      `json:"use_public_runner" bson:"use_public_runner"`
	RunnerId         string    `json:"runner_id" bson:"runner_id"`
	Users            []User    `json:"users" bson:"users"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`

	IsWaitingForResult bool `json:"is_waiting_for_result"`
	IsRunnerOnline     bool `json:"is_runner_online"`

	PersistedAt time.Time
	mutex       *sync.Mutex
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
	durationForWaitingForResultMax = 20 * time.Second
	durationIsUnused               = 10 * time.Minute
)

func NewFile(fileId, fileName, lang, content, userId, userName string) *File {
	if fileName == "" {
		fileName = "File " + time.Now().Format("2006-01-02")
	}
	file := &File{
		ID:               fileId,
		Name:             fileName,
		Lang:             lang,
		Content:          content,
		Writer:           "",
		UsePublicRunner:  true,
		RunnerId:         "",
		UpdatedAt:        time.Now(),
		ContentUpdatedAt: time.Now(),
		Users:            nil,
	}
	file.TouchByUser(userId, userName)
	return file
}

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
		if userName == "" || !util.IsValidName(userName) {
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
	if f.Name == name {
		return true
	}
	f.Name = name
	f.UpdatedAt = time.Now()
	return true
}

func (f *File) SetLang(lang string) bool {
	if _, ok := Langs[lang]; !ok {
		return false
	}
	if f.Lang == lang {
		return true
	}
	f.Lang = lang
	f.UpdatedAt = time.Now()
	return true
}

func (f *File) SetContent(content, userId string) error {
	if len(content) > contentMaxLength {
		return errors.New("content is too long")
	}
	if f.Writer != "" && f.Writer != userId {
		return errors.New("file is locked by another user")
	}
	if f.Content == content {
		return nil
	}

	f.lock()
	defer f.unlock()

	f.Content = content
	f.Writer = userId
	f.ContentUpdatedAt = time.Now()
	f.UpdatedAt = time.Now()
	return nil
}

func (f *File) SetWaitingForResult() {
	if f.IsWaitingForResult {
		return
	}

	f.lock()
	defer f.unlock()

	f.IsWaitingForResult = true
	f.UpdatedAt = time.Now()
}

func (f *File) SetResult(result string) error {
	if len(result) > contentMaxLength {
		return errors.New("result is too long")
	}
	if f.IsWaitingForResult == false && f.Result == result {
		return nil
	}

	f.lock()
	defer f.unlock()

	f.IsWaitingForResult = false
	f.Result = result
	f.UpdatedAt = time.Now()
	return nil
}

func (f *File) SetUserName(userId, userName string) bool {
	if !util.IsValidName(userName) || !util.IsUuid(userId) {
		return false
	}

	f.lock()
	defer f.unlock()

	for i, user := range f.Users {
		if user.ID == userId {
			if f.Users[i].Name != userName {
				f.Users[i].Name = userName
				f.UpdatedAt = time.Now()
			}
			return true
		}
	}

	return false
}

func (f *File) SetRunnerId(runnerId string) bool {
	if !util.IsUuid(runnerId) {
		return false
	}
	if f.RunnerId == runnerId {
		return true
	}
	f.RunnerId = runnerId
	f.UpdatedAt = time.Now()
	return false
}

func (f *File) CleanupUsers() {
	f.lock()
	defer f.unlock()

	changed := false
	for i := len(f.Users) - 1; i >= 0; i-- {
		if time.Since(f.Users[i].TouchedAt) > durationIsActiveFromLastUpdate {
			f.Users = append(f.Users[:i], f.Users[i+1:]...)
			changed = true
		}
	}
	if changed {
		f.UpdatedAt = time.Now()
	}
}

func (f *File) CleanupWaitingForResult() {
	f.lock()
	defer f.unlock()

	if f.IsWaitingForResult && time.Since(f.UpdatedAt) > durationForWaitingForResultMax {
		f.IsWaitingForResult = false
		f.UpdatedAt = time.Now()
	}
}

func (f *File) IsUnused() bool {
	f.lock()
	defer f.unlock()

	if time.Since(f.UpdatedAt) > durationIsUnused {
		for _, user := range f.Users {
			if time.Since(user.TouchedAt) < durationIsUnused {
				return false
			}
		}
		return true
	}
	return false
}

func (f *File) CleanupWriter() {
	if f.Writer == "" {
		return
	}
	if time.Since(f.ContentUpdatedAt) > durationIsWriterStillWriting {
		f.Writer = ""
		f.UpdatedAt = time.Now()
	}
}

func (f *File) lock() {
	if f.mutex == nil {
		f.mutex = &sync.Mutex{}
	}
	f.mutex.Lock()
}

func (f *File) unlock() {
	f.mutex.Unlock()
}
