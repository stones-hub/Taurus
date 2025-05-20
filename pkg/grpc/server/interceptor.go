package server

import (
	"context"
	"log"
	"time"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 服务端拦截器工厂方法

// LoggingServerInterceptor 日志拦截器
func LoggingServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		log.Printf("Server - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err)
		return resp, err
	}
}

// RecoveryServerInterceptor 恢复拦截器
func RecoveryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				err = status.Error(codes.Internal, "Internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

// AuthServerInterceptor 认证拦截器
func AuthServerInterceptor(token string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 从metadata中获取token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		if tokens[0] != token {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}

// RateLimitServerInterceptor 限流拦截器
func RateLimitServerInterceptor(limit int) grpc.UnaryServerInterceptor {
	limiter := rate.NewLimiter(rate.Limit(limit), limit)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !limiter.Allow() {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(ctx, req)
	}
}

// 拦截器链 处理一元请求
func chainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	// interceptors = [] func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// handler =  func(ctx context.Context, req any) (any, error)
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(next grpc.UnaryHandler, interceptor grpc.UnaryServerInterceptor) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					return interceptor(ctx, req, info, next)
				}
			}(chain, interceptors[i])
		}

		return chain(ctx, req)
	}
}

// 拦截器链 处理流请求
func chainStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(next grpc.StreamHandler, interceptor grpc.StreamServerInterceptor) grpc.StreamHandler {
				return func(srv interface{}, ss grpc.ServerStream) error {
					return interceptor(srv, ss, info, next)
				}
			}(chain, interceptors[i])
		}
		return chain(srv, ss)
	}
}

// 同时处理中间件和拦截器
func chainUnaryServerWithMiddleware(mids []UnaryMiddleware, interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		midChain := handler

		for i := len(mids) - 1; i >= 0; i-- {
			midChain = mids[i](midChain)
		}

		chain := midChain

		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(next grpc.UnaryHandler, interceptor grpc.UnaryServerInterceptor) grpc.UnaryHandler {
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					return interceptor(ctx, req, info, next)
				}
			}(chain, interceptors[i])
		}

		return chain(ctx, req)
	}
}

func chainStreamServerWithMiddleware(mids []StreamMiddleware, interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		midChain := handler
		for i := len(mids) - 1; i >= 0; i-- {
			midChain = mids[i](midChain)
		}

		chain := midChain
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(next grpc.StreamHandler, interceptor grpc.StreamServerInterceptor) grpc.StreamHandler {
				return func(srv interface{}, ss grpc.ServerStream) error {
					return interceptor(srv, ss, info, next)
				}
			}(chain, interceptors[i])
		}
		return chain(srv, ss)
	}
}
