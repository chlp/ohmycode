package api

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"ohmycode_runner/pkg/util"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	RunnerId string
	IsPublic bool
	ApiUrl   string

	mutex      *sync.Mutex
	tasksQueue []*Task

	socket *websocket.Conn
}

const timeToSleepBetweenMessagesHandling = 100 * time.Millisecond

func NewApiClient(runnerId string, isPublic bool, apiUrl string) *Client {
	apiClient := Client{
		RunnerId:   runnerId,
		IsPublic:   isPublic,
		ApiUrl:     apiUrl,
		mutex:      &sync.Mutex{},
		tasksQueue: make([]*Task, 0),
	}
	go apiClient.handleReconnection()
	go func() {
		for {
			apiClient.handleWebSocketMessages()
			time.Sleep(timeToSleepBetweenMessagesHandling)
		}
	}()
	return &apiClient
}

type InitMessage struct {
	Action   string `json:"action"`
	RunnerId string `json:"runner_id"`
	IsPublic bool   `json:"is_public"`
}

func (apiClient *Client) createWebSocket() {
	socket, _, err := websocket.DefaultDialer.Dial(apiClient.ApiUrl+"/runner", nil)
	if err != nil {
		util.Log("createWebSocket: can not dial, err: " + err.Error())
		return
	}
	initMessage := InitMessage{
		Action:   "init",
		RunnerId: apiClient.RunnerId,
		IsPublic: apiClient.IsPublic,
	}
	if err := socket.WriteJSON(initMessage); err != nil {
		util.Log("createWebSocket: json err: " + err.Error())
		return
	}
	apiClient.socket = socket
	util.Log("createWebSocket: connected!")
}

func (apiClient *Client) handleReconnection() {
	reconnectAttempts := 0
	for {
		if apiClient.socket == nil {
			reconnectAttempts++
			apiClient.createWebSocket()
		} else {
			reconnectAttempts = 0
		}
		delay := time.Duration(1000*math.Min(math.Pow(2, float64(reconnectAttempts)), 30) + rand.Float64()*3000)
		time.Sleep(delay * time.Millisecond)
	}
}

func (apiClient *Client) handleWebSocketMessages() {
	if apiClient.socket == nil {
		return
	}

	defer func() {
		_ = apiClient.socket.Close()
		apiClient.socket = nil
	}()

	for {
		wsMessageType, message, err := apiClient.socket.ReadMessage()
		if err != nil {
			var closeErr *websocket.CloseError
			if !errors.As(err, &closeErr) {
				util.Log("handleWebSocketMessages: read message err: " + err.Error())
			}
			return
		}

		if wsMessageType == websocket.CloseMessage {
			break
		}
		if wsMessageType != websocket.TextMessage {
			continue
		}

		var tasks []*Task
		if err := json.Unmarshal(message, &tasks); err != nil {
			util.Log("handleWebSocketMessages: json err: " + err.Error())
			continue
		}
		apiClient.mutex.Lock()
		apiClient.tasksQueue = append(apiClient.tasksQueue, tasks...)
		apiClient.mutex.Unlock()
	}
}

func (apiClient *Client) GetTasksRequest() []*Task {
	apiClient.mutex.Lock()
	tasks := apiClient.tasksQueue
	apiClient.tasksQueue = make([]*Task, 0)
	apiClient.mutex.Unlock()
	return tasks
}

func (apiClient *Client) SetResult(result *Task) error {
	if apiClient.socket == nil {
		util.Log("SetResult: can not set result now")
		time.Sleep(100 * time.Millisecond)
		return errors.New("can not set result now")
	}
	result.Action = "set_result"
	if err := apiClient.socket.WriteJSON(result); err != nil {
		util.Log("SetResult: json err: " + err.Error())
		return err
	}
	return nil
}
