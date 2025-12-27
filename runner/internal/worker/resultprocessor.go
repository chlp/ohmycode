package worker

import (
	"context"
	"fmt"
	"ohmycode_runner/internal/api"
	"ohmycode_runner/pkg/util"
	"os"
	"path/filepath"
	"strconv"
)

type ResultProcessor struct {
	apiClient *api.Client
	runnerId  string
	language  string
}

func NewResultProcessor(apiClient *api.Client, runnerId, lang string) *ResultProcessor {
	dir := getDirForResults(lang)
	_ = os.MkdirAll(dir, 0o755)
	// Ensure directories are listable/readable even if umask is restrictive.
	_ = os.Chmod(dir, 0o755)
	return &ResultProcessor{
		apiClient: apiClient,
		runnerId:  runnerId,
		language:  lang,
	}
}

func (rp *ResultProcessor) Process() (int, error) {
	resultsDir := getDirForResults(rp.language)
	resultFiles, err := os.ReadDir(resultsDir)
	if err != nil {
		util.Log(context.Background(), fmt.Sprintf("read results dir error: %v", err))
		return 0, err
	}

	processed := 0
	for _, entry := range resultFiles {
		if !isValidFile(entry) {
			continue
		}
		processed++
		err = rp.processOneFile(resultsDir, entry)
	}

	return processed, err
}

func isValidFile(entry os.DirEntry) bool {
	name := entry.Name()
	return !entry.IsDir() && name != "" && name[0] != '.'
}

func (rp *ResultProcessor) processOneFile(resultsDir string, entry os.DirEntry) error {
	filePath := filepath.Join(resultsDir, entry.Name())

	hash, err := strconv.ParseUint(entry.Name(), 10, 32)
	if err != nil {
		util.Log(context.Background(), fmt.Sprintf("wrong hash error: %v", err))
		_ = os.Remove(filePath)
		return err
	}

	result, err := os.ReadFile(filePath)
	if err != nil {
		util.Log(context.Background(), fmt.Sprintf("read result error: %v", err))
		_ = os.Remove(filePath)
		result = []byte("something wrong with result")
	}

	err = rp.apiClient.SetResult(&api.Task{
		RunnerId: rp.runnerId,
		Lang:     rp.language,
		Hash:     uint32(hash),
		Result:   string(result),
	})
	if err != nil {
		util.Log(context.Background(), fmt.Sprintf("set result error: %v", err))
	} else {
		_ = os.Remove(filePath)
	}
	return err
}

func getDirForResults(lang string) string {
	return fmt.Sprintf("%s/%s/results", dataFolderPath, lang)
}
