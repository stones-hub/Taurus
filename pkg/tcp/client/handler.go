package client

import (
	"context"
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

// DefaultHandler 提供默认的处理器实现
type DefaultHandler struct{}

func (h *DefaultHandler) OnConnect(ctx context.Context, conn net.Conn)                  {}
func (h *DefaultHandler) OnMessage(ctx context.Context, conn net.Conn, msg interface{}) {}
func (h *DefaultHandler) OnClose(ctx context.Context, conn net.Conn)                    {}
func (h *DefaultHandler) OnError(ctx context.Context, conn net.Conn, err error)         {}
