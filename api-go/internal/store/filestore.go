package store

import (
	"errors"
	"ohmycode_api/internal/model"
	"sync"
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

func (fs *FileStore) GetFileOrCreate(fileId, fileName, lang, content, userId, userName string) (*model.File, error) {
	file, err := fs.GetFile(fileId)
	if err != nil {
		return nil, err
	}
	if file != nil {
		return file, nil
	}

	defer fs.lockFileMutex(fileId).Unlock()

	file = model.NewFile(fileId, fileName, lang, content, userId, userName)

	fs.mutex.Lock()
	fs.files[fileId] = file
	fs.mutex.Unlock()

	return file, nil
}

func (fs *FileStore) GetAllFiles() map[string]*model.File {
	return fs.files
}

func (fs *FileStore) GetFile(fileId string) (*model.File, error) {
	defer fs.lockFileMutex(fileId).Unlock()

	if file, ok := fs.files[fileId]; ok {
		return file, nil
	}

	filesRaw, err := fs.db.Select("files", map[string]interface{}{"_id": fileId}, &model.File{})
	if err != nil {
		return nil, err
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

	fs.mutex.Lock()
	fs.files[fileId] = &files[0]
	fs.mutex.Unlock()
	return &files[0], nil
}

func (fs *FileStore) lockFileMutex(fileId string) *sync.Mutex {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	var fileMutex *sync.Mutex
	var ok bool
	if fileMutex, ok = fs.filesMutex[fileId]; !ok {
		fileMutex = &sync.Mutex{}
		fs.filesMutex[fileId] = fileMutex
	}

	fileMutex.Lock()
	return fileMutex
}
