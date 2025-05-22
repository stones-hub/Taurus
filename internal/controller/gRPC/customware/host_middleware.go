package customware

import (
	"Taurus/pkg/grpc/attributes"
	"Taurus/pkg/grpc/server"
	"Taurus/pkg/grpc/server/middleware"
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

func init() {
	server.RegisterMiddleware(HostMiddleware())

	// System default middleware
	// MetricsMiddleware
	server.RegisterMiddleware(middleware.MetricsMiddleware())
	// LoggingMiddleware
	server.RegisterMiddleware(middleware.LoggingMiddleware())
}
