package main

import (
	"Taurus/pkg/tcp/client"
	"Taurus/pkg/tcp/protocol"
	"Taurus/pkg/tcp/protocol/json"
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	JoinRoom  = 1 // 加入房间
	Chat      = 2 // 聊天
	LeaveRoom = 3 // 离开房间
)

var userID int = rand.Intn(1000000)
var roomID int

// 用户信息
type UserMessage struct {
	RoomID  int    `json:"room_id"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

// ClientHandler 实现了客户端的消息处理
type ClientHandler struct {
	stopCh   chan struct{}
	stopOnce sync.Once
}

func (h *ClientHandler) OnClose(ctx context.Context, conn net.Conn) {
	log.Printf("连接关闭: %s", conn.RemoteAddr())
	h.stopOnce.Do(func() {
		close(h.stopCh)
	})
}

func (h *ClientHandler) OnError(ctx context.Context, conn net.Conn, err error) {
	log.Printf("连接错误: %v", err)
}

func (h *ClientHandler) OnConnect(ctx context.Context, conn net.Conn) {
	log.Printf("已连接到服务器: %s", conn.RemoteAddr())
}

func (h *ClientHandler) OnMessage(ctx context.Context, conn net.Conn, message interface{}) {
	// 处理服务器的响应消息
	log.Printf("收到服务端消息: %+v", message)
	// 重新显示提示符
	showPrompt()
}

func showPrompt() {
	fmt.Print("> ")
}

func main() {
	var (
		err error
		c   *client.Client
	)

	// 创建handler
	handler := &ClientHandler{
		stopCh: make(chan struct{}),
	}

	// 创建客户端
	c, err = client.New(":8080",
		protocol.JSON, // 使用JSON协议
		handler,
		client.WithMaxMsgSize(1024*1024), // 1MB
		client.WithBufferSize(1024),      // 缓冲区大小
		client.WithConnectionTimeout(5*time.Second), // 连接超时
		client.WithIdleTimeout(30*time.Second),      // 空闲超时
		client.WithMaxRetries(3),                    // 最大重试次数
	)

	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	if err := c.Connect(); err != nil {
		log.Printf("连接服务器失败: %v", err)
		return
	}

	// 显示帮助信息
	showHelp()
	showPrompt()

	// 初始化sequence
	var sequence uint32 = 1

	// 创建一个channel用于通知标准输入处理goroutine退出
	inputDone := make(chan struct{})

	// 启动goroutine处理标准输入
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case <-inputDone:
				return
			default:
			}

			text := scanner.Text()
			if text == "" {
				showPrompt()
				continue
			}

			input := strings.Split(text, " ")
			msgType := input[0]
			msgData := strings.Join(input[1:], " ")

			switch msgType {
			case "join":
				roomID, err = strconv.Atoi(msgData)
				if err != nil {
					log.Printf("房间ID格式错误: %v", err)
					showPrompt()
					continue
				}

				err = c.Send(&json.Message{
					Type:     JoinRoom,
					Sequence: sequence,
					Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": "进入房间"},
				})

				if err != nil {
					log.Printf("发送消息失败: %v", err)
					showPrompt()
					continue
				}

			case "chat":
				if msgData == "" {
					log.Println("请输入聊天内容")
					showPrompt()
					continue
				}

				err = c.Send(&json.Message{
					Type:     Chat,
					Sequence: sequence,
					Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": msgData},
				})

				if err != nil {
					log.Printf("发送消息失败: %v", err)
					showPrompt()
					continue
				}

			case "leave":
				err = c.Send(&json.Message{
					Type:     LeaveRoom,
					Sequence: sequence,
					Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": "离开房间"},
				})
				if err != nil {
					log.Printf("发送消息失败: %v", err)
					showPrompt()
					continue
				}
				handler.stopOnce.Do(func() {
					close(handler.stopCh)
				})
				return

			case "quit":
				handler.stopOnce.Do(func() {
					close(handler.stopCh)
				})
				return

			case "help":
				showHelp()

			default:
				log.Printf("未知命令: %s", text)
				showHelp()
			}
			sequence++
			showPrompt()
		}
	}()

	// 等待服务端断开连接或用户主动退出
	<-handler.stopCh
	// 通知输入处理goroutine退出
	close(inputDone)
	c.GracefulClose()
	log.Println("客户端已关闭")
}

func showHelp() {
	fmt.Println("\n=== 聊天室命令帮助 ===")
	fmt.Println("join <room_id>  - 加入指定房间")
	fmt.Println("leave           - 离开当前房间")
	fmt.Println("chat <message>  - 发送聊天消息")
	fmt.Println("help            - 显示此帮助信息")
	fmt.Println("quit            - 退出程序")
	fmt.Println("====================")
}
