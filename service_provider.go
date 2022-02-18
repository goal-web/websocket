package websocket

import (
	"github.com/goal-web/contracts"
	"sync"
)

type ServiceProvider struct {
}

func (s ServiceProvider) Register(application contracts.Application) {
	application.Singleton("websocket", func(config contracts.Config) contracts.WebSocket {
		var wsConfig = config.Get("websocket").(Config)

		upgrader = wsConfig.Upgrader

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
