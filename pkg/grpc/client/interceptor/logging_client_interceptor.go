package interceptor

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// LoggingClientInterceptor 日志拦截器
func LoggingClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.Printf("Client - Method:%s\tDuration:%s\tError:%v\n", method, time.Since(start), err)
		return err
	}
}
