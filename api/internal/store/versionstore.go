package store

import (
	"context"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"time"
)

const (
	versionCollection = "file_versions"
	maxVersions       = 20
)

type VersionStore struct {
	db *Db
}

func NewVersionStore(dbConfig DBConfig) *VersionStore {
	return &VersionStore{
		db: newDb(dbConfig),
	}
}

func (vs *VersionStore) Close(ctx context.Context) error {
	if vs.db == nil {
		return nil
	}
	return vs.db.Close(ctx)
}

// SaveVersion creates a new version if enough time has passed since lastVersionedAt.
// Returns the creation time of the new version, or zero time if no version was created.
func (vs *VersionStore) SaveVersion(fileID, content, name, lang string, lastVersionedAt time.Time) (time.Time, error) {
	now := time.Now().UTC()
	currentDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if !lastVersionedAt.IsZero() {
		lastDay := time.Date(lastVersionedAt.Year(), lastVersionedAt.Month(), lastVersionedAt.Day(), 0, 0, 0, 0, time.UTC)
		if lastDay.Equal(currentDay) {
			return time.Time{}, nil
		}
	}

	version := model.FileVersion{
		ID:        util.GenUuid(),
		FileID:    fileID,
		Content:   content,
		Name:      name,
		Lang:      lang,
		CreatedAt: now,
	}

	if err := vs.db.InsertOne(versionCollection, &version); err != nil {
		return time.Time{}, err
	}

	return now, nil
}

func (vs *VersionStore) GetVersions(fileID string) ([]model.FileVersion, error) {
	cursor, err := vs.db.Find(
		versionCollection,
		map[string]interface{}{"file_id": fileID},
		map[string]interface{}{"created_at": -1},
		maxVersions+1,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var versions []model.FileVersion
	if err := cursor.All(context.Background(), &versions); err != nil {
		return nil, err
	}

	// Clean up old versions beyond maxVersions
	if len(versions) > maxVersions {
		for _, v := range versions[maxVersions:] {
			if err := vs.db.DeleteOne(versionCollection, map[string]interface{}{"_id": v.ID}); err != nil {
				util.Log("VersionStore.GetVersions: failed to delete old version: " + err.Error())
			}
		}
		versions = versions[:maxVersions]
	}

	return versions, nil
}

func (vs *VersionStore) GetVersion(versionID string) (*model.FileVersion, error) {
	var version model.FileVersion
	found, err := vs.db.FindOne(versionCollection, map[string]interface{}{"_id": versionID}, &version)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &version, nil
}
