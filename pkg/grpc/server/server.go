package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server gRPC服务器封装
type Server struct {
	server *grpc.Server   // gRPC服务器实例
	opts   *ServerOptions // 服务器配置
}

var GlobalgRPCServer *Server

// NewServer 创建新的gRPC服务器
func NewServer(opts ...ServerOption) (*Server, func(), error) {
	options := DefaultServerOptions()
	for _, opt := range opts {
		opt(options)
	}

	serverOpts := []grpc.ServerOption{}

	// 添加TLS配置
	if options.TLSConfig != nil {
		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(options.TLSConfig)))
	}

	// 添加KeepAlive配置
	if options.KeepAlive != nil {
		serverOpts = append(serverOpts, grpc.KeepaliveParams(*options.KeepAlive))
	}

	// 添加一元拦截器
	if len(options.UnaryMiddlewares) > 0 {
		serverOpts = append(serverOpts, grpc.UnaryInterceptor(chainUnaryServerWithMiddleware(options.UnaryMiddlewares, options.UnaryInterceptors)))
	} else if len(options.UnaryInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.UnaryInterceptor(chainUnaryServer(options.UnaryInterceptors...)))
	}

	// 添加流拦截器
	if len(options.StreamMiddlewares) > 0 {
		serverOpts = append(serverOpts, grpc.StreamInterceptor(chainStreamServerWithMiddleware(options.StreamMiddlewares, options.StreamInterceptors)))
	} else if len(options.StreamInterceptors) > 0 {
		serverOpts = append(serverOpts, grpc.StreamInterceptor(chainStreamServer(options.StreamInterceptors...)))
	}
	server := grpc.NewServer(serverOpts...)
	GlobalgRPCServer = &Server{
		server: server,
		opts:   options,
	}
	return GlobalgRPCServer, func() {
		GlobalgRPCServer.Stop()
		log.Println("GRPC server stopped successfully")
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	log.Println("Starting gRPC server on", s.opts.Address)
	lis, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	return s.server.Serve(lis)
}

// Stop 停止服务器
func (s *Server) Stop() {
	s.server.GracefulStop()
}

// Server 获取原始服务器实例
func (s *Server) Server() *grpc.Server {
	return s.server
}

// 同时处理中间件和拦截器
func chainUnaryServerWithMiddleware(mids []UnaryMiddleware, interceptors []grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
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

func chainStreamServerWithMiddleware(mids []StreamMiddleware, interceptors []grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
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

/*
# 1. 首先确保安装了必要的工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# 2. 生成gRPC代码
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/controller/gRPC/proto/user/user.proto
*/
