package store

import (
	"context"
	"ohmycode_api/internal/model"
	"ohmycode_api/pkg/util"
	"strings"
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

// SaveVersion saves a snapshot of the content as it was at contentUpdatedAt.
// A version is only created if no version has already been saved for the same calendar day
// as contentUpdatedAt. Returns contentUpdatedAt on success, or zero time if no version was created.
// After inserting it trims old versions beyond maxVersions so that cleanup is
// co-located with the write path rather than the read path.
func (vs *VersionStore) SaveVersion(fileID, content, name, lang string, contentUpdatedAt, lastVersionedAt time.Time) (time.Time, error) {
	if contentUpdatedAt.IsZero() {
		return time.Time{}, nil
	}
	contentUpdatedAt = contentUpdatedAt.UTC()
	contentDay := time.Date(contentUpdatedAt.Year(), contentUpdatedAt.Month(), contentUpdatedAt.Day(), 0, 0, 0, 0, time.UTC)

	if !lastVersionedAt.IsZero() {
		lastDay := time.Date(lastVersionedAt.Year(), lastVersionedAt.Month(), lastVersionedAt.Day(), 0, 0, 0, 0, time.UTC)
		if lastDay.Equal(contentDay) {
			return time.Time{}, nil
		}
	}

	preview := diffPreview(vs.latestContent(fileID), content)

	version := model.FileVersion{
		ID:        util.GenId(),
		FileID:    fileID,
		Content:   content,
		Name:      name,
		Lang:      lang,
		Preview:   preview,
		CreatedAt: contentUpdatedAt,
	}

	if err := vs.db.InsertOne(versionCollection, &version); err != nil {
		return time.Time{}, err
	}

	vs.trimOldVersions(fileID)

	return contentUpdatedAt, nil
}

// trimOldVersions deletes versions beyond maxVersions, keeping the most recent ones.
func (vs *VersionStore) trimOldVersions(fileID string) {
	cursor, err := vs.db.Find(
		versionCollection,
		map[string]interface{}{"file_id": fileID},
		map[string]interface{}{"created_at": -1},
		maxVersions+1,
	)
	if err != nil {
		util.Log("VersionStore.trimOldVersions: find error: " + err.Error())
		return
	}
	defer cursor.Close(context.Background())

	var versions []model.FileVersion
	if err := cursor.All(context.Background(), &versions); err != nil {
		util.Log("VersionStore.trimOldVersions: cursor error: " + err.Error())
		return
	}

	for _, v := range versions[min(maxVersions, len(versions)):] {
		if err := vs.db.DeleteOne(versionCollection, map[string]interface{}{"_id": v.ID}); err != nil {
			util.Log("VersionStore.trimOldVersions: delete error: " + err.Error())
		}
	}
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

	if len(versions) > maxVersions {
		versions = versions[:maxVersions]
	}

	return versions, nil
}

// latestContent returns the content of the most recent saved version for fileID, or "" if none.
func (vs *VersionStore) latestContent(fileID string) string {
	cursor, err := vs.db.Find(
		versionCollection,
		map[string]interface{}{"file_id": fileID},
		map[string]interface{}{"created_at": -1},
		1,
	)
	if err != nil {
		return ""
	}
	defer cursor.Close(context.Background())
	var versions []model.FileVersion
	if err := cursor.All(context.Background(), &versions); err != nil || len(versions) == 0 {
		return ""
	}
	return versions[0].Content
}

// diffPreview returns a one-line summary of what changed between oldContent and newContent.
// It finds the first line in newContent absent from oldContent ("+"), or the first line in
// oldContent absent from newContent ("-"). Falls back to first positionally changed line ("~").
func diffPreview(oldContent, newContent string) string {
	if oldContent == newContent {
		return ""
	}

	oldSet := toLineSet(oldContent)
	newSet := toLineSet(newContent)

	for _, l := range splitTrimmedLines(newContent) {
		if !oldSet[l] {
			return truncLine("+ " + l)
		}
	}
	for _, l := range splitTrimmedLines(oldContent) {
		if !newSet[l] {
			return truncLine("- " + l)
		}
	}

	// All lines exist in both — find first positionally different line
	oldLines := splitTrimmedLines(oldContent)
	newLines := splitTrimmedLines(newContent)
	for i := 0; i < min(len(oldLines), len(newLines)); i++ {
		if oldLines[i] != newLines[i] {
			return truncLine("~ " + newLines[i])
		}
	}
	return ""
}

func toLineSet(content string) map[string]bool {
	m := make(map[string]bool)
	for _, l := range splitTrimmedLines(content) {
		m[l] = true
	}
	return m
}

func splitTrimmedLines(content string) []string {
	raw := strings.Split(content, "\n")
	out := make([]string, 0, len(raw))
	for _, l := range raw {
		if t := strings.TrimSpace(l); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func truncLine(s string) string {
	r := []rune(s)
	if len(r) > 62 {
		return string(r[:62]) + "…"
	}
	return s
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
