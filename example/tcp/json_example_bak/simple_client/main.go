package main

import (
	"Taurus/pkg/tcp/client"
	"Taurus/pkg/tcp/protocol"
	"Taurus/pkg/tcp/protocol/json"
	"context"
	"log"
	"net"
	"time"
)

// ClientHandler 实现了 client.Handler 接口
type ClientHandler struct{}

func (h *ClientHandler) OnConnect(ctx context.Context, conn net.Conn) {
	log.Printf("连接服务器成功")
}

func (h *ClientHandler) OnMessage(ctx context.Context, conn net.Conn, message interface{}) {
	msg := message.(*json.Message)
	log.Printf("收到服务器响应: type=%d, sequence=%d, data=%v", msg.Type, msg.Sequence, msg.Data)
}

func (h *ClientHandler) OnClose(ctx context.Context, conn net.Conn) {
	log.Printf("连接关闭")
}

func (h *ClientHandler) OnError(ctx context.Context, conn net.Conn, err error) {
	log.Printf("客户端错误: %v", err)
}

func main() {
	// 创建客户端实例
	c, err := client.New(":8080",
		protocol.JSON, // 使用JSON协议
		&ClientHandler{},
		client.WithMaxMsgSize(1024*1024),            // 1MB 最大消息大小
		client.WithBufferSize(1024),                 // 1KB 缓冲区大小
		client.WithConnectionTimeout(5*time.Second), // 连接超时5秒
		client.WithIdleTimeout(30*time.Second),      // 空闲超时30秒
		client.WithMaxRetries(3),                    // 最大重试3次
		client.WithBaseRetryDelay(time.Second),      // 初始重试延迟1秒
		client.WithMaxRetryDelay(5*time.Second),     // 最大重试延迟5秒
	)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 连接服务器
	if err := c.Connect(); err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}

	// 发送测试消息
	sequence := uint32(1)
	for {
		msg := &json.Message{
			Type:     1, // 假设 1 是 Echo 消息类型
			Sequence: sequence,
			Data: map[string]interface{}{
				"message": "Hello, Server!",
				"time":    time.Now().String(),
			},
		}

		log.Println(c.RemoteAddr(), c.LocalAddr())

		if err := c.Send(msg); err != nil {
			log.Printf("发送消息失败: %v", err)
			break
		} else {
			log.Printf("发送消息成功: %v", msg)
		}

		sequence++
		time.Sleep(time.Second * 3) // 每秒发送一条消息
	}

	// 关闭连接
	c.Close()
}
