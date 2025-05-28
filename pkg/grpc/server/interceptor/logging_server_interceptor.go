package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// 服务端拦截器工厂方法

// LoggingServerInterceptor 日志拦截器
func LoggingServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		log.Printf("Method: %s, Duration: %s, Error: %v", info.FullMethod, time.Since(start), err)
		return resp, err
	}
}
