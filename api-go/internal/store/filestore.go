package store

import (
	"errors"
	"ohmycode_api/internal/model"
	"sync"
	"time"
)

type FileStore struct {
	mutex      *sync.RWMutex
	files      map[string]*model.File
	filesMutex map[string]*sync.Mutex
	db         *Db
}

func NewFileStore(dbConfig DBConfig) *FileStore {
	return &FileStore{
		mutex:      &sync.RWMutex{},
		files:      make(map[string]*model.File),
		filesMutex: make(map[string]*sync.Mutex),
		db:         newDb(dbConfig),
	}
}

func (s *FileStore) GetFileOrCreate(fileId, fileName, lang, content, userId, userName string) (*model.File, error) {
	file, err := s.GetFile(fileId)
	if err != nil {
		return nil, err
	}
	if file != nil {
		return file, nil
	}

	defer s.lockFileMutex(fileId).Unlock()

	file = &model.File{
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

	s.mutex.Lock()
	s.files[fileId] = file
	s.mutex.Unlock()

	return file, nil
}

func (s *FileStore) GetAllFiles() map[string]*model.File {
	return s.files
}

func (s *FileStore) GetFile(fileId string) (*model.File, error) {
	defer s.lockFileMutex(fileId).Unlock()

	if file, ok := s.files[fileId]; ok {
		return file, nil
	}

	filesRaw, err := s.db.Select("files", map[string]interface{}{"_id": fileId}, &model.File{})
	if err != nil {
		panic(err)
	}
	files, ok := filesRaw.([]model.File)
	if !ok {
		return nil, errors.New("problem with type conversion")
	}
	if len(files) == 0 {
		return nil, nil
	}
	if len(files) != 1 {
		return nil, errors.New("wrong count")
	}
	if files[0].ID != fileId {
		return nil, errors.New("problem with fileId")
	}

	s.mutex.Lock()
	s.files[fileId] = &files[0]
	s.mutex.Unlock()
	return &files[0], nil
}

func (s *FileStore) lockFileMutex(fileId string) *sync.Mutex {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if fileMutex, ok := s.filesMutex[fileId]; ok {
		return fileMutex
	}
	fileMutex := &sync.Mutex{}
	s.filesMutex[fileId] = fileMutex
	fileMutex.Lock()
	return fileMutex
}
