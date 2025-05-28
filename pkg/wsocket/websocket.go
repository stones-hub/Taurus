package wsocket

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

/*
HTTP 跨域：通过 CORS（跨域资源共享）头来控制，CorsMiddleware 已经处理了 HTTP 请求的跨域问题。
WebSocket 跨域：WebSocket 不依赖 CORS，而是通过 Origin 请求头来验证跨域。WebSocket 的跨域检查由服务器端的 CheckOrigin 方法控制。
*/

// Upgrader is used to upgrade HTTP connections to WebSocket connections
var upgrader websocket.Upgrader

// MessageHandler defines a function type for handling messages
type MessageHandler func(conn *websocket.Conn, messageType int, message []byte) error

// Initialize initializes the WebSocket upgrader
func Initialize() {
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// TODO Allow all origins for simplicity; customize as needed
			return true
		},
	}
	log.Println("WebSocket upgrader initialized")
}

// HandleWebSocket handles WebSocket connections with a custom message handler
func HandleWebSocket(w http.ResponseWriter, r *http.Request, handler MessageHandler) {
	defer func() { // websocket的特殊性，需要在处理函数中解决异常、错误问题， 不能用middleware来解决
		if err := recover(); err != nil {
			log.Printf("Recovered from panic in websocket: %v\n%s", err, debug.Stack())
		}
	}()

	// 对于websocket来说，每个请求是长连接，放在中间件来处理trace_id 不合适，所以需要手动生成
	hash := md5.Sum([]byte(uuid.New().String()))
	traceid := hex.EncodeToString(hash[:])

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection, traceid: %s, error: %v\n", traceid, err)
		http.Error(w, "Failed to establish websocket connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Printf("websocket connection established, traceid: %s\n", traceid)

	// Use the custom message handler.  handle sending message and receive message
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message, error: %v\n", err)
			break
		}

		log.Printf("Received message, traceid: %s, message: %s\n", traceid, message)

		// Call the custom message handler
		if err := handler(conn, messageType, message); err != nil {
			log.Printf("Error handling message, error: %v\n", err)
			break
		}
	}

	log.Printf("websocket connection closed, traceid: %s\n", traceid)
}

// HandleWebSocket handles WebSocket connections with a custom message handler
func HandleWebSocketRoom(w http.ResponseWriter, r *http.Request, handler MessageHandler, hub *WebSocketHub, roomName string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered from panic in websocket: %v\n", err)
		}
	}()

	// 验证用户身份
	userid, err := authenticateUser(r)
	if err != nil {
		log.Printf("Authentication failed: %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 检查用户是否有权进入房间
	if !checkRoomAccess(userid, roomName) {
		log.Printf("Access denied for user %s to room %s\n", userid, roomName)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection, error: %v\n", err)
		http.Error(w, "Failed to establish websocket connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	room := hub.GetOrCreateRoom(roomName)
	room.AddClient(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message, error: %v\n", err)
			break
		}

		log.Printf("Received message: %s\n", message)

		// 将消息发送到房间的广播通道
		room.BroadcastMessage(message)

		if err := handler(conn, messageType, message); err != nil {
			log.Printf("Error handling message, error: %v\n", err)
			break
		}
	}

	room.RemoveClient(conn)
}

// authenticateUser 验证用户身份
func authenticateUser(r *http.Request) (string, error) {
	// 在这里实现您的身份验证逻辑
	// 返回用户ID或错误
	log.Printf("authenticateUser, request: %v\n", r)
	return "userid", nil
}

// checkRoomAccess 检查用户是否有权进入房间
func checkRoomAccess(userid, roomName string) bool {
	// 在这里实现您的权限检查逻辑
	// 返回 true 表示有权进入，false 表示无权进入
	log.Printf("checkRoomAccess, userid: %s, roomName: %s\n", userid, roomName)
	return true
}
