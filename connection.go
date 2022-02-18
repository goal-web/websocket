package websocket

import (
	"github.com/goal-web/contracts"
	"github.com/gorilla/websocket"
	"sync"
)

type Connection struct {
	fd    uint64
	ws    *websocket.Conn
	mutex sync.Mutex
}

func NewConnection(ws *websocket.Conn) contracts.WebSocketConnection {
	return &Connection{
		fd:    0,
		ws:    ws,
		mutex: sync.Mutex{},
	}
}

func (conn *Connection) SetFd(fd uint64) {
	conn.fd = fd
}

func (conn *Connection) Fd() uint64 {
	return conn.fd
}

func (conn *Connection) Close() error {
	return conn.ws.Close()
}

func (conn *Connection) Send(message interface{}) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn.ws.WriteJSON(message)
}

func (conn *Connection) SendBinary(bytes []byte) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn.ws.WriteMessage(websocket.BinaryMessage, bytes)
}

func (conn *Connection) SendBytes(bytes []byte) error {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn.ws.WriteMessage(websocket.TextMessage, bytes)
}
