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
	"Taurus/pkg/tcp"
	"Taurus/pkg/tcp/protocol"
	"Taurus/pkg/tcp/protocol/json"
	"fmt"
	"log"
	"sync"
)

// MessageHandler 实现了一个简单的消息处理器
type MessageHandler struct {
	rooms  map[string]*Room // 房间ID->房间信息
	roomMu sync.RWMutex
}

// Room 房间信息
type Room struct {
	ID    string                       // 房间ID
	Name  string                       // 房间名称
	Users map[uint64][]*tcp.Connection // 用户ID -> 该用户在本房间的所有连接
}

const (
	JoinRoom  = 1
	Chat      = 2
	LeaveRoom = 3
)

func (h *MessageHandler) OnConnect(conn *tcp.Connection) {
	log.Printf("新连接建立: %d", conn.ID())
}

func (h *MessageHandler) OnClose(conn *tcp.Connection) {

	// 在所有房间中查找并清理这个连接
	for _, room := range h.rooms {
		for userId, conns := range room.Users {
			if newConns := removeConn(conns, conn); len(newConns) < len(conns) {
				if len(newConns) == 0 {
					delete(room.Users, userId)
				} else {
					room.Users[userId] = newConns
				}
				break
			}
		}
	}
	log.Printf("连接关闭: %d", conn.ID())
}

func (h *MessageHandler) OnError(conn *tcp.Connection, err error) {
	if conn != nil {
		log.Printf("连接错误 [%d]: %v", conn.ID(), err)
	} else {
		log.Printf("服务器错误: %v", err)
	}
}

func (h *MessageHandler) OnMessage(conn *tcp.Connection, message interface{}) {
	msg, ok := message.(*json.Message)
	if !ok {
		log.Printf("消息格式错误")
		return
	}

	if msg.Data == nil {
		log.Printf("消息数据为空")
		return
	}

	switch msg.Type {
	case JoinRoom:
		h.handleJoinRoom(conn, msg)

	case Chat:
		// 处理聊天消息
		h.handleChatMessage(conn, msg)

	case LeaveRoom:
		// 处理离开房间请求
		h.handleLeaveRoom(conn, msg)
	default:
		log.Printf("未知消息类型: %s, %d, %d, %d", msg.Data, msg.Type, msg.Sequence, msg.Timestamp)
	}
}

// removeConn 从切片中移除指定连接
func removeConn(slice []*tcp.Connection, conn *tcp.Connection) []*tcp.Connection {
	for i, v := range slice {
		if v.ID() == conn.ID() {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (h *MessageHandler) handleJoinRoom(conn *tcp.Connection, msg *json.Message) {
	// JSON解析数字默认为float64，需要转换
	room_id := fmt.Sprintf("%.0f", msg.Data["room_id"].(float64))
	user_id := uint64(msg.Data["user_id"].(float64))

	h.roomMu.Lock()
	defer h.roomMu.Unlock()

	// 检查用户是否已经在其他房间
	for _, room := range h.rooms {
		if _, exists := room.Users[user_id]; exists {
			log.Printf("用户 %d 已经在房间 %s 中", user_id, room.ID)
			return
		}
	}

	// 如果房间不存在就创建
	if !h.isRoomExist(room_id) {
		h.rooms[room_id] = &Room{
			ID:    room_id,
			Name:  room_id,
			Users: make(map[uint64][]*tcp.Connection),
		}
	}

	// 将用户的连接添加到房间
	room := h.rooms[room_id]
	room.Users[user_id] = append(room.Users[user_id], conn)

	// 广播用户进入房间的消息
	h.broadcastToRoom(room_id, msg)
	log.Printf("用户 %d 的连接 %d 加入房间 %s", user_id, conn.ID(), room_id)
}

// 聊天信息
func (h *MessageHandler) handleChatMessage(conn *tcp.Connection, msg *json.Message) {
	room_id := fmt.Sprintf("%.0f", msg.Data["room_id"].(float64))
	log.Printf("服务端收到聊天消息, 房间(%s), 开始广播: %s\n", room_id, msg.Data)
	h.broadcastToRoom(room_id, msg)
}

func (h *MessageHandler) handleLeaveRoom(conn *tcp.Connection, msg *json.Message) {
	room_id := fmt.Sprintf("%.0f", msg.Data["room_id"].(float64))
	user_id := uint64(msg.Data["user_id"].(float64))

	h.roomMu.Lock()
	defer h.roomMu.Unlock()

	if h.isRoomExist(room_id) {
		room := h.rooms[room_id]
		// 从用户的连接列表中移除这个连接
		if conns, exists := room.Users[user_id]; exists {
			newConns := removeConn(conns, conn)
			if len(newConns) == 0 {
				// 如果用户没有任何连接了，从房间中删除该用户
				delete(room.Users, user_id)
			} else {
				room.Users[user_id] = newConns
			}
		}

		log.Printf("服务端收到离开房间消息, 房间(%s), 用户(%d), 开始广播: %s\n", room_id, user_id, msg.Data)
		// 广播用户离开消息
		h.broadcastToRoom(room_id, msg)
		// 关闭连接
		conn.Close()
	}
}

// broadcastToRoom 向房间内的所有用户的所有连接广播消息
func (h *MessageHandler) broadcastToRoom(room_id string, msg *json.Message) {
	if !h.isRoomExist(room_id) {
		return
	}

	room := h.rooms[room_id]
	for _, conns := range room.Users {
		for _, conn := range conns {
			if err := conn.Send(msg); err != nil {
				log.Printf("发送消息失败，连接ID: %d, 错误: %v", conn.ID(), err)
			}
		}
	}
}

// 判断房间在不在
func (h *MessageHandler) isRoomExist(room_id string) bool {
	_, ok := h.rooms[room_id]
	return ok
}

func main() {
	// 创建协议实例
	p, err := protocol.NewProtocol(
		protocol.WithType(protocol.JSON),
		protocol.WithMaxMessageSize(1024*1024), // 1MB
	)
	if err != nil {
		log.Fatalf("创建协议失败: %v", err)
	}

	// 创建handler
	handler := &MessageHandler{
		rooms: make(map[string]*Room),
	}

	// 创建服务器
	server, stop, err := tcp.NewServer(":8080", p, handler,
		tcp.WithMaxConnections(1000), // 最大1000个连接
	)

	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 启动服务器
	log.Println("服务器启动在 :8080")
	if err := server.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

	stop()
}
