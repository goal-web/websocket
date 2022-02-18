package websocket

import (
	"errors"
	"github.com/goal-web/contracts"
	"sync"
)

var (
	ConnectionDontExistsErr = errors.New("connection does not exist")
)

type WebSocket struct {
	mutex       sync.RWMutex
	connections map[uint64]contracts.WebSocketConnection
	count       uint64
}

func (ws *WebSocket) Add(connect contracts.WebSocketConnection) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	ws.count++
	var fd = ws.count
	ws.connections[fd] = connect
	connect.SetFd(fd)
}

func (ws *WebSocket) Close(fd uint64) error {
	var conn, exists = ws.connections[fd]
	if exists {
		ws.mutex.Lock()
		defer ws.mutex.Unlock()
		delete(ws.connections, fd)
		return conn.Close()
	}

	return ConnectionDontExistsErr
}

func (ws *WebSocket) Send(fd uint64, message interface{}) error {
	var conn, exists = ws.connections[fd]
	if exists {
		return conn.Send(message)
	}

	return ConnectionDontExistsErr
}
