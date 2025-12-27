package api

import (
	"context"
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
	ctx context.Context

	RunnerId string
	IsPublic bool
	ApiUrl   string

	mutex      *sync.Mutex
	tasksQueue []*Task

	socketMu *sync.RWMutex
	writeMu  *sync.Mutex
	socket   *websocket.Conn
}

const timeToSleepBetweenMessagesHandling = 100 * time.Millisecond

func NewApiClient(ctx context.Context, runnerId string, isPublic bool, apiUrl string) *Client {
	apiClient := Client{
		ctx:        ctx,
		RunnerId:   runnerId,
		IsPublic:   isPublic,
		ApiUrl:     apiUrl,
		mutex:      &sync.Mutex{},
		socketMu:   &sync.RWMutex{},
		writeMu:    &sync.Mutex{},
		tasksQueue: make([]*Task, 0),
	}
	go apiClient.handleReconnection()
	go func() {
		for {
			select {
			case <-apiClient.ctx.Done():
				return
			default:
			}
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

func (apiClient *Client) getSocket() *websocket.Conn {
	apiClient.socketMu.RLock()
	socket := apiClient.socket
	apiClient.socketMu.RUnlock()
	return socket
}

func (apiClient *Client) setSocket(socket *websocket.Conn) {
	apiClient.socketMu.Lock()
	apiClient.socket = socket
	apiClient.socketMu.Unlock()
}

func (apiClient *Client) clearSocketIfCurrent(socket *websocket.Conn) {
	apiClient.socketMu.Lock()
	if apiClient.socket == socket {
		apiClient.socket = nil
	}
	apiClient.socketMu.Unlock()
}

func (apiClient *Client) createWebSocket() {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
	socket, _, err := dialer.DialContext(apiClient.ctx, apiClient.ApiUrl+"/runner", nil)
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
		_ = socket.Close()
		return
	}
	apiClient.setSocket(socket)
	util.Log("createWebSocket: connected!")
}

func (apiClient *Client) handleReconnection() {
	reconnectAttempts := 0
	for {
		select {
		case <-apiClient.ctx.Done():
			return
		default:
		}

		if apiClient.getSocket() == nil {
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
	socket := apiClient.getSocket()
	if socket == nil {
		return
	}

	defer func() {
		// Close should not race with concurrent writes (SetResult).
		apiClient.writeMu.Lock()
		_ = socket.Close()
		apiClient.writeMu.Unlock()
		apiClient.clearSocketIfCurrent(socket)
	}()

	for {
		wsMessageType, message, err := socket.ReadMessage()
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
	apiClient.writeMu.Lock()
	defer apiClient.writeMu.Unlock()

	socket := apiClient.getSocket()
	if socket == nil {
		util.Log("SetResult: can not set result now")
		time.Sleep(100 * time.Millisecond)
		return errors.New("can not set result now")
	}
	result.Action = "set_result"
	if err := socket.WriteJSON(result); err != nil {
		util.Log("SetResult: json err: " + err.Error())
		_ = socket.Close()
		apiClient.clearSocketIfCurrent(socket)
		return err
	}
	return nil
}
