package main

import (
	"Taurus/pkg/tcp"
	"Taurus/pkg/tcp/protocol"
	"Taurus/pkg/tcp/protocol/json"
	"log"
)

// EchoHandler 实现了 tcp.Handler 接口
type EchoHandler struct{}

func (h *EchoHandler) OnConnect(conn *tcp.Connection) {
	log.Printf("新连接建立: %d", conn.ID())
}

func (h *EchoHandler) OnMessage(conn *tcp.Connection, message interface{}) {
	msg := message.(*json.Message)
	log.Printf("收到消息: type=%d, sequence=%d, data=%v", msg.Type, msg.Sequence, msg.Data)

	// 回显消息
	response := &json.Message{
		Type:     msg.Type,
		Sequence: msg.Sequence,
		Data:     msg.Data,
	}
	if err := conn.Send(response); err != nil {
		log.Printf("发送消息失败: %v", err)
	} else {
		log.Printf("发送消息成功: %v", response)
	}
}

func (h *EchoHandler) OnClose(conn *tcp.Connection) {
	log.Printf("连接关闭: %d", conn.ID())
}

func (h *EchoHandler) OnError(conn *tcp.Connection, err error) {
	if conn != nil {
		log.Printf("连接 %d 发生错误: %v", conn.ID(), err)
	} else {
		log.Printf("服务器错误: %v", err)
	}
}

func main() {
	// 创建协议实例
	p, err := protocol.NewProtocol(
		protocol.WithType(protocol.JSON),
		protocol.WithMaxMessageSize(1024*1024), // 1MB 最大消息大小
	)
	if err != nil {
		log.Fatalf("创建协议失败: %v", err)
	}

	// 创建服务器实例
	server := tcp.NewServer(":8080",
		tcp.WithProtocol(p),
		tcp.WithHandler(&EchoHandler{}),
		tcp.WithMaxConnections(1000),
	)

	// 启动服务器
	log.Println("Echo 服务器启动在 :8080...")
	if err := server.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
