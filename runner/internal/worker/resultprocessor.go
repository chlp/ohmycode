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
	_ = os.MkdirAll(getDirForResults(lang), os.ModePerm)
	return &ResultProcessor{
		apiClient: apiClient,
		runnerId:  runnerId,
		language:  lang,
	}
}

func (rp *ResultProcessor) Process() error {
	resultsDir := getDirForResults(rp.language)
	resultFiles, err := os.ReadDir(resultsDir)
	if err != nil {
		util.Log(context.Background(), fmt.Sprintf("read results dir error: %v", err))
		return err
	}

	for _, entry := range resultFiles {
		if !isValidFile(entry) {
			continue
		}
		err = rp.processOneFile(resultsDir, entry)
	}

	return err
}

func isValidFile(entry os.DirEntry) bool {
	return !entry.IsDir() && entry.Name()[0] != '.'
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
