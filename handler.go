package websocket

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/http"
	"github.com/goal-web/supports/exceptions"
	"github.com/goal-web/supports/logs"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

func New(controller contracts.WebSocketController) any {
	return func(request *http.Request, serializer contracts.Serializer, socket contracts.WebSocket, handler contracts.ExceptionHandler) error {
		var ws, err = upgrader.Upgrade(request.Context.Response(), request.Request(), nil)

		if err != nil {
			logs.WithError(err).Error("websocket.New: Upgrade failed")
			return err
		}

		var fd = socket.GetFd()

		if err = controller.OnConnect(request, fd); err != nil {
			logs.WithError(err).Error("websocket.New: OnConnect failed")
			return err
		}

		var conn = NewConnection(ws, fd)
		socket.Add(conn)

		defer func() {
			controller.OnClose(fd)
			if closeErr := socket.Close(conn.Fd()); closeErr != nil {
				logs.WithError(closeErr).Error("websocket.New: Connection close failed")
			}
		}()

		for {
			// Read
			var msgType, msg, readErr = ws.ReadMessage()
			if readErr != nil {
				logs.WithError(readErr).Error("websocket.New: Failed to read message")
				return readErr
			}

			switch msgType {
			case websocket.TextMessage, websocket.BinaryMessage:
				go handleMessage(NewFrame(msg, conn, serializer), controller, handler)
			case websocket.CloseMessage:
				return nil
			}
		}
	}
}

func handleMessage(frame contracts.WebSocketFrame, controller contracts.WebSocketController, handler contracts.ExceptionHandler) {
	defer func() {
		if panicValue := recover(); panicValue != nil {
			handler.Handle(Exception{
				Exception: exceptions.WithRecover(panicValue),
			})
		}
	}()
	controller.OnMessage(frame)
}
