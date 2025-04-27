package ws

import (
	"log"

	"github.com/google/wire"
	"github.com/gorilla/websocket"
)

type DemoWs struct {
}

var DemoWsSet = wire.NewSet(wire.Struct(new(DemoWs), "*"))

func (w *DemoWs) HandleMessage(conn *websocket.Conn, messageType int, message []byte) error {
	log.Println("Received message:", string(message))
	conn.WriteMessage(websocket.TextMessage, []byte("Hello, World!"))
	return nil
}
