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
	"Taurus/example/http_validator/controller"
	"Taurus/pkg/middleware"
	"log"
	"net/http"
)

func main() {
	// 创建控制器实例
	userCtrl := &controller.UserController{}

	// 创建路由
	mux := http.NewServeMux()

	// 注册路由，使用验证中间件
	mux.Handle("/api/user/create", middleware.ValidationMiddleware(&controller.UserRequest{})(http.HandlerFunc(userCtrl.CreateUser)))
	mux.Handle("/api/user/update", middleware.ValidationMiddleware(&controller.UpdateUserRequest{})(http.HandlerFunc(userCtrl.UpdateUser)))

	// 启动服务器
	log.Println("HTTP服务器启动在 :8080 端口")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
