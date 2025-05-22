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
