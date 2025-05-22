package main

import (
	"Taurus/pkg/wsocket"
	"log"

	"github.com/gorilla/websocket"
)

// 框架中已经集成了ws协议， 只需要按自己的业务需求，注册自己的handler即可
func main() {
	wsocket.RegisterHandler("demo", DemoHandler{})
}

type DemoHandler struct{}

func (h DemoHandler) Handle(conn *websocket.Conn, messageType int, message []byte) error {
	log.Printf("demo handler received message: %s", string(message))
	conn.WriteMessage(messageType, message)
	return nil
}
