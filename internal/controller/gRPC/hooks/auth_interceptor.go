package hooks

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

// gRPC auth 拦截器
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("gRPC -> 自定义 AuthInterceptor: %v", req)
		return handler(ctx, req)
	}
}
