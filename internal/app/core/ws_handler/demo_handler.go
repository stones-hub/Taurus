package ws_handler

import (
	"Taurus/pkg/wsocket"
	"log"

	"github.com/gorilla/websocket"
)

type DemoHandler struct{}

func (h DemoHandler) Handle(conn *websocket.Conn, messageType int, message []byte) error {
	log.Printf("demo handler received message: %s", string(message))
	conn.WriteMessage(messageType, message)
	return nil
}

func init() {
	wsocket.RegisterHandler("demo", DemoHandler{})
}
