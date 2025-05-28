package hooks

import (
	"Taurus/pkg/grpc/attributes"
	"log"

	"google.golang.org/grpc"
)

// gRPC host 中间件
func HostMiddleware() attributes.UnaryMiddleware {
	return func(handler grpc.UnaryHandler) grpc.UnaryHandler {
		log.Printf("gRPC -> 自定义 HostMiddleware")
		return handler
	}
}
