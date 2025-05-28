package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryClientInterceptor 重试拦截器
func RetryClientInterceptor(maxRetries int) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var err error
		for i := 0; i < maxRetries; i++ {
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err == nil {
				return nil
			}
			if status.Code(err) == codes.Unavailable {
				time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
				continue
			}
			return err
		}
		return err
	}
}
