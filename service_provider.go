package websocket

import (
	"github.com/goal-web/contracts"
	"sync"
)

type ServiceProvider struct {
}

func (s ServiceProvider) Register(application contracts.Application) {
	application.Singleton("websocket", func() contracts.WebSocket {
		return &WebSocket{
			mutex:       sync.RWMutex{},
			connections: map[uint64]contracts.WebSocketConnection{},
			count:       0,
		}
	})
}

func (s ServiceProvider) Start() error {
	return nil
}

func (s ServiceProvider) Stop() {
}
