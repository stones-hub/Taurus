package gRPC

import (
	"Taurus/pkg/grpc/server"
	"context"

	"google.golang.org/grpc"
)

// gRPC host 中间件
func HostMiddleware() server.UnaryMiddleware {
	return func(handler grpc.UnaryHandler) grpc.UnaryHandler {
		return handler
	}
}

// gRPC auth 拦截器
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}
