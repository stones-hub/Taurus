package hooks

import (
	"Taurus/pkg/grpc/attributes"
	"log"

	"google.golang.org/grpc"
)

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

// gRPC host 中间件
func HostMiddleware() attributes.UnaryMiddleware {
	return func(handler grpc.UnaryHandler) grpc.UnaryHandler {
		log.Printf("gRPC -> 自定义 HostMiddleware")
		return handler
	}
}
