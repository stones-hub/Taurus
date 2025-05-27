package main

import (
	"context"
	"log"
	"net"

	"Taurus/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// unaryServerInterceptor 创建一元 RPC 拦截器
func unaryServerInterceptor(provider telemetry.TracerProvider) grpc.UnaryServerInterceptor {
	tracer := provider.Tracer("grpc.server")

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 创建新的 span
		opts := []trace.SpanStartOption{
			trace.WithAttributes(
				// 添加 gRPC 特定的属性
				attribute.String("rpc.service", info.FullMethod),
				attribute.String("rpc.system", "grpc"),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		}

		ctx, span := tracer.Start(ctx, info.FullMethod, opts...)
		defer span.End()

		// 调用处理器
		resp, err := handler(ctx, req)

		// 处理错误
		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.RecordError(err)
		}

		return resp, err
	}
}

// streamServerInterceptor 创建流式 RPC 拦截器
func streamServerInterceptor(provider telemetry.TracerProvider) grpc.StreamServerInterceptor {
	tracer := provider.Tracer("grpc.server")

	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// 创建新的 span
		ctx := ss.Context()
		opts := []trace.SpanStartOption{
			trace.WithAttributes(
				// 添加 gRPC 特定的属性
				attribute.String("rpc.service", info.FullMethod),
				attribute.String("rpc.system", "grpc"),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		}

		ctx, span := tracer.Start(ctx, info.FullMethod, opts...)
		defer span.End()

		// 包装 ServerStream 以传递追踪上下文
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// 调用处理器
		err := handler(srv, wrappedStream)

		// 处理错误
		if err != nil {
			s, _ := status.FromError(err)
			span.SetStatus(codes.Error, s.Message())
			span.RecordError(err)
		}

		return err
	}
}

// wrappedServerStream 包装 grpc.ServerStream 以传递追踪上下文
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func main() {
	// 初始化 provider
	provider, err := telemetry.NewOTelProvider(
		telemetry.WithServiceName("grpc-demo"),
		telemetry.WithEnvironment("dev"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Shutdown(context.Background())

	// 创建带追踪的 gRPC 服务器
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryServerInterceptor(provider)),
		grpc.StreamInterceptor(streamServerInterceptor(provider)),
	)

	// 启动服务器
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server starting on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
