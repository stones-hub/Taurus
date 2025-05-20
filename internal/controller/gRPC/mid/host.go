package gRPC

import (
	"Taurus/pkg/grpc/server"
	"context"
	"log"

	"google.golang.org/grpc"
)

// gRPC host 中间件
func HostMiddleware() server.UnaryMiddleware {
	return func(handler grpc.UnaryHandler) grpc.UnaryHandler {
		log.Println("HostMiddleware")
		return handler
	}
}

// gRPC auth 拦截器
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("AuthInterceptor: %v", req)
		return handler(ctx, req)
	}
}

func init() {
	server.RegisterMiddleware(HostMiddleware())
	server.RegisterInterceptor(AuthInterceptor())
}
