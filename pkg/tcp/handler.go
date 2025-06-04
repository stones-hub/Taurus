package tcp

import "log"

// ---------------------------------------------------------------------------------------------------------------------
// Handler 定义了连接事件处理的接口, 注意如果handler中实现了子协程的逻辑切记需要监听ctx.Done()，否则子协程不会退出，造成协程泄漏
// ---------------------------------------------------------------------------------------------------------------------
// 实现者需要处理各种连接生命周期事件。
type Handler interface {
	OnConnect(conn *Connection)                      // 当新连接建立时调用
	OnMessage(conn *Connection, message interface{}) // 当收到消息时调用
	OnClose(conn *Connection)                        // 当连接关闭时调用
	OnError(conn *Connection, err error)             // 当发生错误时调用
}

var handlers = make(map[string]Handler)

func RegisterHandler(name string, handler Handler) {
	if _, ok := handlers[name]; ok {
		log.Printf("Handler %s already registered", name)
		return
	}
	handlers[name] = handler
}

func GetHandler(name string) Handler {
	if handler, ok := handlers[name]; ok {
		return handler
	}
	return &defaultHandler{}
}

type defaultHandler struct{}

func (h *defaultHandler) OnConnect(conn *Connection) {
	conn.SetAttr("handler", "default")
	conn.SetAttr("ip", conn.RemoteAddr())
	log.Printf("连接建立: %v", conn.RemoteAddr())
}

func (h *defaultHandler) OnMessage(conn *Connection, message interface{}) {
	log.Printf("收到消息: %v", message)
}

func (h *defaultHandler) OnClose(conn *Connection) {
	log.Printf("连接关闭: %v", conn.RemoteAddr())
}

func (h *defaultHandler) OnError(conn *Connection, err error) {
	log.Printf("连接错误: %v, 错误: %v", conn.RemoteAddr(), err)
}
