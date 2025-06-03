package client

import (
	"context"
	"log"
	"net"
)

// Handler 定义了客户端事件处理接口
type Handler interface {
	// OnConnect 当连接建立时调用
	OnConnect(ctx context.Context, conn net.Conn)
	// OnMessage 当收到消息时调用
	OnMessage(ctx context.Context, conn net.Conn, msg interface{})
	// OnClose 当连接关闭时调用
	OnClose(ctx context.Context, conn net.Conn)
	// OnError 当发生错误时调用
	OnError(ctx context.Context, conn net.Conn, err error)
}

// TODO: 需要实现一个默认的处理器实现
type DefaultHandler struct{}

func (h *DefaultHandler) OnConnect(ctx context.Context, conn net.Conn) {
	log.Printf("连接建立: %s", conn.RemoteAddr())
}
func (h *DefaultHandler) OnMessage(ctx context.Context, conn net.Conn, msg interface{}) {
	log.Printf("收到消息: %v", msg)
}
func (h *DefaultHandler) OnClose(ctx context.Context, conn net.Conn) {
	log.Printf("连接关闭: %s", conn.RemoteAddr())
}
func (h *DefaultHandler) OnError(ctx context.Context, conn net.Conn, err error) {
	log.Printf("发生错误: %s, %v", conn.RemoteAddr(), err)
}
