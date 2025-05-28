package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor 创建一个认证拦截器
func AuthInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 将 token 添加到 metadata 中
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamAuthInterceptor 创建一个流式认证拦截器
func StreamAuthInterceptor(token string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		// 将 token 添加到 metadata 中
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", token)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
