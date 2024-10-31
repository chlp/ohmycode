package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ohmycode_runner/pkg/util"
	"time"
)

type Client struct {
	RunnerId string
	IsPublic bool
	ApiUrl   string
}

const getTasksEndpoint = "/run/get_tasks"
const setResultEndpoint = "/result/set"
const keepAliveRequestTimeout = 35 * time.Second
const requestTimeout = 3 * time.Second

type Response struct {
	IsOk bool
	Code int
	Data []Task
}

func NewApiClient(runnerId string, isPublic bool, apiUrl string) *Client {
	return &Client{RunnerId: runnerId, IsPublic: isPublic, ApiUrl: apiUrl}
}

type getTasksReq struct {
	RunnerId    string `json:"runner_id"`
	IsPublic    bool   `json:"is_public"`
	IsKeepAlive bool   `json:"is_keep_alive"`
}

func (apiClient *Client) GetTasksRequest() ([]*Task, error) {
	url := fmt.Sprintf("%s%s", apiClient.ApiUrl, getTasksEndpoint)

	jsonParams, err := json.Marshal(getTasksReq{
		RunnerId:    apiClient.RunnerId,
		IsPublic:    apiClient.IsPublic,
		IsKeepAlive: true,
	})
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
	if err != nil {
		return nil, fmt.Errorf("http request creating error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{
		Timeout: keepAliveRequestTimeout,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpClient do error: %v", err)
	}
	defer resp.Body.Close()

	var tasks []*Task
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("httpClient wrong code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %v", err)
	}

	if err = json.Unmarshal(body, &tasks); err != nil {
		return nil, fmt.Errorf("response unmarshal error: %v", err)
	}

	return tasks, nil
}

func (apiClient *Client) SetResult(result *Task) error {
	url := fmt.Sprintf("%s%s", apiClient.ApiUrl, setResultEndpoint)

	jsonParams, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
	if err != nil {
		return fmt.Errorf("http request creating error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{
		Timeout: requestTimeout,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("httpClient do error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		util.Log(context.Background(), fmt.Sprintf("task not found for lang: %v", result.Lang))
		return nil
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("httpClient wrong code: %v", resp.StatusCode)
	}
	return nil
}
