package store

import (
	"errors"
	"ohmycode_api/internal/model"
	"sync"
	"time"
)

type Store struct {
	mutex      *sync.RWMutex
	files      map[string]*model.File
	filesMutex map[string]*sync.Mutex
	runners    map[string]*model.Runner
	db         *Db
}

func NewStore(dbConfig DBConfig) *Store {
	return &Store{
		mutex:      &sync.RWMutex{},
		files:      make(map[string]*model.File),
		filesMutex: make(map[string]*sync.Mutex),
		runners:    make(map[string]*model.Runner),
		db:         newDb(dbConfig),
	}
}

func (s *Store) GetFileOrCreate(fileId, fileName, lang, content, userId, userName string) (*model.File, error) {
	file, err := s.GetFile(fileId)
	if err != nil {
		return nil, err
	}
	if file != nil {
		return file, nil
	}

	defer s.lockFileMutex(fileId).Unlock()

	runnerId := ""
	runner := s.GetPublicRunner()
	if runner != nil {
		runnerId = runner.ID
	}
	file = &model.File{
		ID:               fileId,
		Name:             fileName,
		Lang:             lang,
		Content:          content,
		Writer:           "",
		RunnerId:         runnerId,
		UpdatedAt:        time.Time{},
		ContentUpdatedAt: time.Time{},
		RunnerCheckedAt:  time.Time{},
		Users:            nil,
	}
	file.TouchByUser(userId, userName)

	s.mutex.Lock()
	s.files[fileId] = file
	s.mutex.Unlock()

	return file, nil
}

func (s *Store) GetAllFiles() map[string]*model.File {
	return s.files
}

func (s *Store) GetFile(fileId string) (*model.File, error) {
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

func (s *Store) GetPublicRunner() *model.Runner {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var result *model.Runner
	for _, runner := range s.runners {
		if runner.IsPublic {
			if result == nil || result.CheckedAt.Before(runner.CheckedAt) {
				result = runner
			}
		}
	}
	return result
}

func (s *Store) lockFileMutex(fileId string) *sync.Mutex {
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
