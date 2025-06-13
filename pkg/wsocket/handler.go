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

package wsocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Handler interface {
	Handle(conn *websocket.Conn, messageType int, message []byte) error
}

var handlers = make(map[string]Handler)

func RegisterHandler(name string, handler Handler) {
	if _, ok := handlers[name]; ok {
		log.Printf("handler %s already registered", name)
	}
	handlers[name] = handler
}

func GetHandler(name string) Handler {
	if handler, ok := handlers[name]; ok {
		return handler
	}
	return defaultHandler{}
}

type defaultHandler struct{}

func (h defaultHandler) Handle(conn *websocket.Conn, messageType int, message []byte) error {
	log.Printf("received message: %s", string(message))
	conn.WriteMessage(messageType, message)
	return nil
}
