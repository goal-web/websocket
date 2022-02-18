package websocket

import (
	"github.com/goal-web/contracts"
	"github.com/gorilla/websocket"
)

type Connection struct {
	fd uint64
	ws *websocket.Conn
}

func NewConnection(ws *websocket.Conn) contracts.WebSocketConnection {
	return &Connection{
		fd: 4,
		ws: ws,
	}
}

func (conn *Connection) SendBinary(bytes []byte) error {
	return conn.ws.WriteMessage(websocket.BinaryMessage, bytes)
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
	return conn.ws.WriteJSON(message)
}

func (conn *Connection) SendBytes(bytes []byte) error {
	return conn.ws.WriteMessage(websocket.TextMessage, bytes)
}
