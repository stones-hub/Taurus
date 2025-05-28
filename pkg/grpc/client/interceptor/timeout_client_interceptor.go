package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// 客户端拦截器

// TimeoutClientInterceptor 超时拦截器
func TimeoutClientInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
