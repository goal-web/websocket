package websocket

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/http"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
	"github.com/gorilla/websocket"
)

var (
	upgrade = websocket.Upgrader{}
)

func New(controller contracts.WebSocketController) interface{} {
	return func(request *http.Request, serializer contracts.Serializer, socket contracts.WebSocket, handler contracts.ExceptionHandler) error {
		var ws, err = upgrade.Upgrade(request.Context.Response(), request.Request(), nil)

		if err != nil {
			logs.WithError(err).Error("websocket.New: Upgrade failed")
			return err
		}

		if err = controller.OnConnect(request); err != nil {
			logs.WithError(err).Error("websocket.New: OnConnect failed")
			return err
		}

		var conn = NewConnection(ws)
		socket.Add(conn)
		defer func() {
			if closeErr := socket.Close(conn.Fd()); closeErr != nil {
				logs.WithError(closeErr).Error("websocket.New: Connection close failed")
			}
		}()

		for {
			// Read
			msgType, msg, readErr := ws.ReadMessage()
			if readErr != nil {
				logs.WithError(readErr).Error("websocket.New: Failed to read message")
				return readErr
			}

			switch msgType {
			case websocket.TextMessage, websocket.BinaryMessage:
				go handleMessage(NewFrame(msg, conn, serializer), controller, handler)
			case websocket.CloseMessage:
				_ = socket.Close(conn.Fd())
			}
		}
	}
}

func handleMessage(frame contracts.WebSocketFrame, controller contracts.WebSocketController, handler contracts.ExceptionHandler) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			handler.Handle(Exception{
				Exception: exceptions.WithRecover(panicValue, contracts.Fields{
					"msg": frame.RawString(),
					"fd":  frame.Connection().Fd(),
				}),
			})
		}
	}()
	controller.OnMessage(frame)
}
