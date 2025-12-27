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

type TaskDistributor struct {
	apiClient *api.Client
	runnerId  string
	languages map[string]interface{}
}

func NewTaskDistributor(apiClient *api.Client, runnerId string, languages []string) *TaskDistributor {
	td := &TaskDistributor{
		apiClient: apiClient,
		runnerId:  runnerId,
		languages: make(map[string]interface{}),
	}
	for _, lang := range languages {
		td.languages[lang] = nil
		// Avoid world-writable permissions by default; containers run as root anyway.
		_ = os.MkdirAll(getDirForRequests(lang), 0o755)
	}
	return td
}

func (td *TaskDistributor) Process() error {
	var err error
	tasks := td.apiClient.GetTasksRequest()
	for _, task := range tasks {
		err = td.moveTask(task)
		if err != nil {
			err = td.apiClient.SetResult(&api.Task{
				RunnerId: td.runnerId,
				Lang:     task.Lang,
				Hash:     task.Hash,
				Result:   err.Error(),
			})
			if err != nil {
				util.Log(context.Background(), fmt.Sprintf("set wrong result error: %v", err))
			}
		}
	}
	return err
}

func (td *TaskDistributor) moveTask(task *api.Task) error {
	_, ok := td.languages[task.Lang]
	if !ok {
		return fmt.Errorf("no runner for %s", task.Lang)
	}
	// Task files are data, not executables.
	fileName := strconv.FormatUint(uint64(task.Hash), 10)
	filePath := filepath.Join(getDirForRequests(task.Lang), fileName)
	tmpPath := filePath + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(task.Content), 0o644); err != nil {
		return fmt.Errorf("can not move task: %v", err)
	}
	if err := os.Rename(tmpPath, filePath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("can not finalize task: %v", err)
	}
	return nil
}

func getDirForRequests(lang string) string {
	return fmt.Sprintf("%s/%s/requests", dataFolderPath, lang)
}
