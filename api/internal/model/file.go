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
	Content          *string   `json:"content,omitempty" bson:"content"`
	ContentUpdatedAt time.Time `json:"content_updated_at" bson:"content_updated_at"`
	Result           string    `json:"result" bson:"result"`
	Writer           string    `json:"writer_id" bson:"writer_id"`
	UsePublicRunner  bool      `json:"use_public_runner" bson:"use_public_runner"`
	RunnerId         string    `json:"runner_id" bson:"runner_id"`
	Users            []User    `json:"users" bson:"users"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	Persisted        bool      `json:"persisted"`

	IsWaitingForResult bool `json:"is_waiting_for_result"`
	IsRunnerOnline     bool `json:"is_runner_online"`

	PersistedAt time.Time `json:"-"`
	mutex       sync.Mutex
	updateCh    chan struct{}
}

type User struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	TouchedAt time.Time `json:"touched_at" bson:"touched_at"`
}

const (
	contentMaxLength               = 512 * (1 << 10) // 512 Kb
	durationIsActiveFromLastUpdate = 5 * time.Second
	durationIsWriterStillWriting   = 2 * time.Second
	durationForWaitingForResultMax = 20 * time.Second
	durationIsUnused               = 10 * time.Minute
)

func NewFile(fileId, fileName, lang, content, userId, userName string) *File {
	if fileName == "" {
		fileName = "File " + time.Now().Format("2006-01-02 15:04:05")
	}
	file := &File{
		ID:                 fileId,
		Name:               fileName,
		Lang:               lang,
		Content:            &content,
		Writer:             "",
		UsePublicRunner:    true,
		RunnerId:           "",
		UpdatedAt:          time.Now(),
		ContentUpdatedAt:   time.Now(),
		Users:              nil,
		IsWaitingForResult: false,
		Persisted:          false,
		IsRunnerOnline:     true, // todo: set correct
		updateCh:           make(chan struct{}, 1),
	}
	file.TouchByUser(userId, userName)
	return file
}

func (f *File) Updates() <-chan struct{} {
	f.lock()
	defer f.unlock()
	if f.updateCh == nil {
		f.updateCh = make(chan struct{}, 1)
	}
	return f.updateCh
}

func (f *File) signalUpdatedLocked() {
	if f.updateCh == nil {
		f.updateCh = make(chan struct{}, 1)
	}
	select {
	case f.updateCh <- struct{}{}:
	default:
	}
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
		f.signalUpdatedLocked()
	}
}

func (f *File) SetName(name string) bool {
	if !util.IsValidName(name) {
		return false
	}
	f.lock()
	defer f.unlock()
	if f.Name == name {
		return true
	}
	f.Persisted = true
	f.Name = name
	f.UpdatedAt = time.Now()
	f.signalUpdatedLocked()
	return true
}

func (f *File) SetLang(lang string) bool {
	if _, ok := Langs[lang]; !ok {
		return false
	}
	f.lock()
	defer f.unlock()
	if f.Lang == lang {
		return true
	}
	f.Lang = lang
	f.UpdatedAt = time.Now()
	f.signalUpdatedLocked()
	return true
}

func (f *File) SetContent(content, appId string) error {
	if len(content) > contentMaxLength {
		return errors.New("content is too long")
	}
	f.lock()
	defer f.unlock()
	if f.Writer != "" && f.Writer != appId {
		return errors.New("file is locked by another user")
	}

	f.Persisted = true
	f.Content = &content
	f.Writer = appId
	f.ContentUpdatedAt = time.Now()
	f.UpdatedAt = time.Now()
	f.signalUpdatedLocked()
	return nil
}

func (f *File) SetWaitingForResult() {
	if f.IsWaitingForResult {
		return
	}

	f.lock()
	defer f.unlock()

	f.Persisted = true
	f.IsWaitingForResult = true
	f.Result = "Started execution at " + time.Now().UTC().Format("15:04:05") + " UTC"
	f.UpdatedAt = time.Now()
	f.signalUpdatedLocked()
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

	f.Persisted = true
	f.IsWaitingForResult = false
	f.Result = result
	f.UpdatedAt = time.Now()
	f.signalUpdatedLocked()
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
				f.signalUpdatedLocked()
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
	f.lock()
	defer f.unlock()
	if f.RunnerId == runnerId {
		return true
	}
	f.Persisted = true
	f.RunnerId = runnerId
	f.UpdatedAt = time.Now()
	f.signalUpdatedLocked()
	return true
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
		f.signalUpdatedLocked()
	}
}

func (f *File) CleanupWaitingForResult() {
	f.lock()
	defer f.unlock()

	if f.IsWaitingForResult && time.Since(f.UpdatedAt) > durationForWaitingForResultMax {
		f.IsWaitingForResult = false
		f.UpdatedAt = time.Now()
		f.signalUpdatedLocked()
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
	f.lock()
	defer f.unlock()
	if f.Writer == "" {
		return
	}
	if time.Since(f.ContentUpdatedAt) > durationIsWriterStillWriting {
		f.Writer = ""
		f.UpdatedAt = time.Now()
		f.signalUpdatedLocked()
	}
}

func (f *File) lock() {
	f.mutex.Lock()
}

func (f *File) unlock() {
	f.mutex.Unlock()
}
