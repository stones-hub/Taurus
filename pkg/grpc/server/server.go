package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server gRPC服务器封装
type Server struct {
	server *grpc.Server   // gRPC服务器实例
	opts   *ServerOptions // 服务器配置
}

// NewServer 创建新的gRPC服务器
func NewServer(opts ...ServerOption) *Server {
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
	if len(options.UnaryInterceptors) > 0 || len(options.UnaryMiddlewares) > 0 {
		serverOpts = append(serverOpts, grpc.UnaryInterceptor(chainUnaryServerWithMiddleware(options.UnaryMiddlewares, options.UnaryInterceptors...)))
	}

	// 添加流拦截器
	if len(options.StreamInterceptors) > 0 || len(options.StreamMiddlewares) > 0 {
		serverOpts = append(serverOpts, grpc.StreamInterceptor(chainStreamServerWithMiddleware(options.StreamMiddlewares, options.StreamInterceptors...)))
	}

	server := grpc.NewServer(serverOpts...)

	return &Server{
		server: server,
		opts:   options,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
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
