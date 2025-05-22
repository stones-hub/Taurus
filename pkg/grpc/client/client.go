package client

import (
	"Taurus/pkg/grpc/attributes"
	"context"
	"crypto/tls"
	"fmt"
	"time"

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

	// 这个是client的context， 并不是 每次请求各种服务的context， 每次请求各种服务的context需要自己创建
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
		dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(attributes.ChainUnaryClient(options.UnaryInterceptors...)))
	}
	if len(options.StreamInterceptors) > 0 {
		dialOpts = append(dialOpts, grpc.WithStreamInterceptor(attributes.ChainStreamClient(options.StreamInterceptors...)))
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

func (c *Client) Options() *ClientOptions {
	return c.opts
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

func (c *Client) Token() string {
	return c.opts.Token
}

func (c *Client) Address() string {
	return c.opts.Address
}

func (c *Client) Timeout() time.Duration {
	return c.opts.Timeout
}

func (c *Client) TLSConfig() *tls.Config {
	return c.opts.TLSConfig
}
