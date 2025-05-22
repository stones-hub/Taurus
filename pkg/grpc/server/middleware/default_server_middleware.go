package middleware

import (
	"Taurus/pkg/grpc/attributes"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

// 监控中间件
func MetricsMiddleware() attributes.UnaryMiddleware {
	return func(next grpc.UnaryHandler) grpc.UnaryHandler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			start := time.Now()

			// 调用下一个处理函数
			resp, err := next(ctx, req)

			// 记录处理时间
			duration := time.Since(start)

			log.Printf("Duration: %s, Error: %v", duration, err)

			return resp, err
		}
	}
}

// 日志中间件
func LoggingMiddleware() attributes.UnaryMiddleware {
	return func(next grpc.UnaryHandler) grpc.UnaryHandler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			fmt.Printf("Request: %v\n", req)

			resp, err := next(ctx, req)

			log.Printf("Response: %v, Error: %v", resp, err)
			return resp, err
		}
	}
}
