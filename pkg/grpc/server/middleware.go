package server

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// 定义中间件类型
type UnaryMiddleware func(grpc.UnaryHandler) grpc.UnaryHandler

type StreamMiddleware func(grpc.StreamHandler) grpc.StreamHandler

// 监控中间件
func MetricsMiddleware() UnaryMiddleware {
	return func(next grpc.UnaryHandler) grpc.UnaryHandler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			start := time.Now()

			// 调用下一个处理函数
			resp, err := next(ctx, req)

			// 记录处理时间
			duration := time.Since(start)
			fmt.Printf("处理时间: %v\n", duration)

			return resp, err
		}
	}
}

// 日志中间件
func LoggingMiddleware() UnaryMiddleware {
	return func(next grpc.UnaryHandler) grpc.UnaryHandler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			fmt.Printf("请求开始: %v\n", req)

			resp, err := next(ctx, req)

			fmt.Printf("请求结束: %v, 错误: %v\n", resp, err)
			return resp, err
		}
	}
}

// Notice: This function is EXPERIMENTAL and may be changed or removed in the future.
// middlewares chain
func ChainUnaryMiddleware(middlewares ...UnaryMiddleware) UnaryMiddleware {
	return func(next grpc.UnaryHandler) grpc.UnaryHandler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Notice: This function is EXPERIMENTAL and may be changed or removed in the future.
// middlewares chain
func ChainStreamMiddleware(middlewares ...StreamMiddleware) StreamMiddleware {
	return func(next grpc.StreamHandler) grpc.StreamHandler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
