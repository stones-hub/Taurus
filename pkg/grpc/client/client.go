package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Client gRPC客户端封装
type Client struct {
	conn   *grpc.ClientConn
	opts   *ClientOptions
	ctx    context.Context
	cancel context.CancelFunc
}

// NewClient 创建新的gRPC客户端
func NewClient(opts ...ClientOption) (*Client, error) {
	options := DefaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)

	dialOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	}

	// 添加TLS配置
	if options.TLSConfig != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.TLSConfig)))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 添加拦截器
	if len(options.UnaryInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(chainUnaryClient(options.UnaryInterceptors...)))
	}
	if len(options.StreamInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithStreamInterceptor(chainStreamClient(options.StreamInterceptors...)))
	}

	// 添加KeepAlive配置
	if options.KeepAlive != nil {
		dialOpts = append(dialOpts, grpc.WithKeepaliveParams(*options.KeepAlive))
	}

	conn, err := grpc.Dial(options.Address, dialOpts...)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	return &Client{
		conn:   conn,
		opts:   options,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Conn 获取原始连接
func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

// Close 关闭连接
func (c *Client) Close() error {
	c.cancel()
	return c.conn.Close()
}
