package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

// Room 代表一个聊天室
type Room struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
}

// WebSocketHub 管理多个聊天室
type WebSocketHub struct {
	rooms map[string]*Room
}

// NewWebSocketHub 创建一个新的 WebSocketHub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		rooms: make(map[string]*Room),
	}
}

// GetOrCreateRoom 获取或创建一个房间
func (hub *WebSocketHub) GetOrCreateRoom(roomName string) *Room {
	room, exists := hub.rooms[roomName]
	if !exists {
		room = &Room{
			clients:   make(map[*websocket.Conn]bool),
			broadcast: make(chan []byte),
		}
		hub.rooms[roomName] = room
		go room.start()
	}
	return room
}

// AdminBroadcast 向指定房间广播消息
func (hub *WebSocketHub) AdminBroadcast(roomName string, message []byte) {
	if room, exists := hub.rooms[roomName]; exists {
		room.BroadcastMessage(message)
	}
}

// start 启动房间的广播协程
func (room *Room) start() {
	for {
		message := <-room.broadcast
		for client := range room.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error broadcasting message to client: %v\n", err)
				client.Close()
				delete(room.clients, client)
			}
		}
	}
}

// AddClient 添加一个新的 WebSocket 客户端到房间
func (room *Room) AddClient(conn *websocket.Conn) {
	room.clients[conn] = true
}

// RemoveClient 移除一个 WebSocket 客户端从房间
func (room *Room) RemoveClient(conn *websocket.Conn) {
	delete(room.clients, conn)
}

// BroadcastMessage 向房间内的客户端广播消息
func (room *Room) BroadcastMessage(message []byte) {
	room.broadcast <- message
}
