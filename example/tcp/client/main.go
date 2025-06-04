package main

import (
	"Taurus/pkg/tcp/client"
	"Taurus/pkg/tcp/protocol"
	"Taurus/pkg/tcp/protocol/json"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
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
}

func (h *ClientHandler) OnClose(ctx context.Context, conn net.Conn) {
	log.Printf("连接关闭: %s", conn.RemoteAddr())
}

func (h *ClientHandler) OnError(ctx context.Context, conn net.Conn, err error) {
	log.Printf("连接错误: %v", err)
}

func (h *ClientHandler) OnConnect(ctx context.Context, conn net.Conn) {
	log.Printf("已连接到服务器: %s", conn.RemoteAddr())
}

func (h *ClientHandler) OnMessage(ctx context.Context, conn net.Conn, message interface{}) {
	// 解析message
	// 解析message
	msg, ok := message.(*json.Message)
	if !ok {
		log.Printf("解析消息失败: %+v", message)
	} else {
		// 将纳秒时间戳转换为time.Time
		t := time.Unix(0, msg.Timestamp)
		timeStr := t.Format("15:04:05")
		fmt.Printf("\n[%s] 👤用户-%v: %v\n", timeStr, msg.Data["user_id"], msg.Data["message"])
	}
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
	handler := &ClientHandler{}

	// 创建客户端
	c, err = client.New(":8080",
		protocol.JSON, // 使用JSON协议
		handler,
		client.WithMaxMsgSize(1024*1024), // 1MB
		client.WithBufferSize(1024),      // 缓冲区大小
		client.WithConnectionTimeout(5*time.Second), // 连接超时
		client.WithIdleTimeout(5*time.Minute),       // 空闲超时
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

	// 初始化readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "quit",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	// 初始化sequence
	var sequence uint32 = 1

	// 启动goroutine处理输入
	go func() {
		for {
			line, err := rl.Readline()
			if err != nil { // io.EOF, readline.ErrInterrupt
				c.Close()
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				showPrompt()
				continue
			}

			input := strings.Split(line, " ")
			msgType := input[0]
			msgData := strings.Join(input[1:], " ")

			switch msgType {
			case "join":
				roomID, err = strconv.Atoi(msgData)
				if err != nil {
					log.Printf("房间ID格式错误: %v", err)
					continue
				}

				err = c.Send(&json.Message{
					Type:     JoinRoom,
					Sequence: sequence,
					Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": "进入房间"},
				})

				if err != nil {
					log.Printf("发送消息失败: %v", err)
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
					continue
				}
				return

			case "quit":
				c.Close()
				return

			case "help":
				showHelp()

			default:
				err = c.Send(&json.Message{
					Type:     Chat,
					Sequence: sequence,
					Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": line},
				})

				if err != nil {
					log.Printf("发送消息失败: %v", err)
					continue
				}
			}
			sequence++
		}
	}()

	c.Start()
	log.Println("客户端已关闭")
}

func showHelp() {
	fmt.Println("\n=== 聊天室命令帮助 ===")
	fmt.Println("join <room_id>  - 加入指定房间")
	fmt.Println("leave           - 离开当前房间")
	fmt.Println("<message>  	 - 发送聊天消息")
	fmt.Println("help            - 显示此帮助信息")
	fmt.Println("quit            - 退出程序")
	fmt.Println("====================")
}
