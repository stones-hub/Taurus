package customceptor

import (
	"Taurus/pkg/grpc/server"
	"Taurus/pkg/grpc/server/interceptor"
	"context"
	"log"

	"google.golang.org/grpc"
)

// gRPC auth 拦截器
func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("gRPC -> 自定义 AuthInterceptor: %v", req)
		return handler(ctx, req)
	}
}

// 初始化 gRPC 中间件 or 拦截器
func init() {
	server.RegisterInterceptor(AuthInterceptor())

	// System default interceptor
	// LoggingServerInterceptor
	server.RegisterInterceptor(interceptor.LoggingServerInterceptor())
	// RecoveryServerInterceptor
	server.RegisterInterceptor(interceptor.RecoveryServerInterceptor())
	// AuthServerInterceptor
	server.RegisterInterceptor(interceptor.AuthServerInterceptor("123456"))
	// RateLimitServerInterceptor
	server.RegisterInterceptor(interceptor.RateLimitServerInterceptor(10))

}
