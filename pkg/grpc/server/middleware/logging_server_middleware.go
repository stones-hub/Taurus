package middleware

import (
	"Taurus/pkg/grpc/attributes"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

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
