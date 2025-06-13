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
	"Taurus/pkg/wsocket"
	"log"

	"github.com/gorilla/websocket"
)

// 框架中已经集成了ws协议， 只需要按自己的业务需求，注册自己的handler即可
func main() {
	wsocket.RegisterHandler("demo", DemoHandler{})
}

type DemoHandler struct{}

func (h DemoHandler) Handle(conn *websocket.Conn, messageType int, message []byte) error {
	log.Printf("demo handler received message: %s", string(message))
	conn.WriteMessage(messageType, message)
	return nil
}
