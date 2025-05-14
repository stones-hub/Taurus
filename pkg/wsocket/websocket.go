package wsocket

import (
	"Taurus/pkg/contextx"
	"log"
	"net/http"
	"runtime/debug"

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
			// Allow all origins for simplicity; customize as needed
			return true
		},
	}
	log.Println("WebSocket upgrader initialized")
}

// HandleWebSocket handles WebSocket connections with a custom message handler
func HandleWebSocket(w http.ResponseWriter, r *http.Request, handler MessageHandler) {
	defer func() { // websocket的特殊性，需要在处理函数中解决异常、错误问题， 不能用middleware来解决
		if err := recover(); err != nil {
			log.Printf("Recovered from panic in WebSocket: %v\n%s", err, debug.Stack())
		}
	}()

	// Retrieve traceID from context
	rc, ok := contextx.GetRequestContext(r.Context())
	if !ok {
		log.Println("Failed to retrieve traceID from context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	traceID := rc.TraceID

	log.Printf("New WebSocket connection, traceID: %s\n", traceID)

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection, traceID: %s, error: %v\n", traceID, err)
		http.Error(w, "Failed to establish WebSocket connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection established, traceID: %s\n", traceID)

	// Use the custom message handler
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message, traceID: %s, error: %v\n", traceID, err)
			break
		}

		log.Printf("Received message, traceID: %s, message: %s\n", traceID, message)

		// Call the custom message handler
		if err := handler(conn, messageType, message); err != nil {
			log.Printf("Error handling message, traceID: %s, error: %v\n", traceID, err)
			break
		}
	}

	log.Printf("WebSocket connection closed, traceID: %s\n", traceID)
}

// HandleWebSocket handles WebSocket connections with a custom message handler
func HandleWebSocketRoom(w http.ResponseWriter, r *http.Request, handler MessageHandler, hub *WebSocketHub, roomName string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered from panic in WebSocket: %v\n", err)
		}
	}()

	// 验证用户身份
	userID, err := authenticateUser(r)
	if err != nil {
		log.Printf("Authentication failed: %v\n", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 检查用户是否有权进入房间
	if !checkRoomAccess(userID, roomName) {
		log.Printf("Access denied for user %s to room %s\n", userID, roomName)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection, error: %v\n", err)
		http.Error(w, "Failed to establish WebSocket connection", http.StatusInternalServerError)
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
	return "userID", nil
}

// checkRoomAccess 检查用户是否有权进入房间
func checkRoomAccess(userID, roomName string) bool {
	// 在这里实现您的权限检查逻辑
	// 返回 true 表示有权进入，false 表示无权进入
	return true
}
