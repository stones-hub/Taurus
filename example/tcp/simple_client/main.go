// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

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
var wg sync.WaitGroup

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
	// 处理服务器的响应消息
	log.Printf("收到服务端消息: %+v", message)
}

func showPrompt() {
	fmt.Print("> ")
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

	wg.Add(1)
	// 启动一个goroutine来接收服务器消息
	go func() {
		defer wg.Done()
		for {
			msg, err := c.SimpleReceive()
			if err != nil {
				log.Printf("接收消息失败: %v", err)
				return
			}
			log.Printf("\n收到服务端消息: %+v", msg)
			showPrompt()
		}
	}()

	// 处理用户输入
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
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

			err = c.SimpleSend(&json.Message{
				Type:     JoinRoom,
				Sequence: sequence,
				Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": "进入房间"},
			})

			if err != nil {
				log.Printf("发送消息失败: %v", err)
				showPrompt()
				continue
			}

		case "leave":
			err = c.SimpleSend(&json.Message{
				Type:     LeaveRoom,
				Sequence: sequence,
				Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": "离开房间"},
			})
			if err != nil {
				log.Printf("发送消息失败: %v", err)
				showPrompt()
				continue
			}
			c.Close()
			return

		case "quit":
			c.Close()
			return

		case "help":
			showHelp()

		default:
			err = c.SimpleSend(&json.Message{
				Type:     Chat,
				Sequence: sequence,
				Data:     map[string]interface{}{"room_id": roomID, "user_id": userID, "message": text},
			})

			if err != nil {
				log.Printf("发送消息失败: %v", err)
				showPrompt()
				continue
			}
			log.Printf("发送消息成功: %s", text)
		}
		sequence++
		showPrompt()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("读取输入错误: %v", err)
	}

	wg.Wait()

}
