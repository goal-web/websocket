package websocket

import (
	"github.com/goal-web/contracts"
)

func Default(handler func(frame contracts.WebSocketFrame)) interface{} {
	return New(&DefaultController{messageHandler: handler})
}

type DefaultController struct {
	messageHandler func(frame contracts.WebSocketFrame)
}

func (d *DefaultController) OnConnect(request contracts.HttpRequest) error {
	return nil
}

func (d *DefaultController) OnMessage(frame contracts.WebSocketFrame) {
	d.messageHandler(frame)
}
