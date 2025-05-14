package wsocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Handler interface {
	Handle(conn *websocket.Conn, messageType int, message []byte) error
}

var handlers = make(map[string]Handler)

func RegisterHandler(name string, handler Handler) {
	if _, ok := handlers[name]; ok {
		log.Printf("handler %s already registered", name)
	}
	handlers[name] = handler
}

func GetHandler(name string) Handler {
	if handler, ok := handlers[name]; ok {
		return handler
	}
	return defaultHandler{}
}

type defaultHandler struct{}

func (h defaultHandler) Handle(conn *websocket.Conn, messageType int, message []byte) error {
	log.Printf("received message: %s", string(message))
	conn.WriteMessage(messageType, message)
	return nil
}
