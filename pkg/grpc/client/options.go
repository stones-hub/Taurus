package client

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// ClientOption 定义客户端配置选项
type ClientOption func(*ClientOptions)

// ClientOptions 包含所有客户端配置
type ClientOptions struct {
	// 基础配置
	Address   string
	Timeout   time.Duration
	TLSConfig *tls.Config
	Token     string

	// 高级配置
	KeepAlive          *keepalive.ClientParameters
	UnaryInterceptors  []grpc.UnaryClientInterceptor
	StreamInterceptors []grpc.StreamClientInterceptor
}

// DefaultClientOptions 返回默认配置
func DefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		Timeout: 5 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		KeepAlive: &keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		},
	}
}

// WithAddress 设置服务器地址
func WithAddress(addr string) ClientOption {
	return func(o *ClientOptions) {
		o.Address = addr
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.Timeout = timeout
	}
}

// WithInsecure set no use tls
func WithInsecure() ClientOption {
	return func(o *ClientOptions) {
		o.TLSConfig = nil // 设置为 nil 表示使用非安全连接
	}
}

// WithTLS 设置TLS配置
func WithTLS(config *tls.Config) ClientOption {
	return func(o *ClientOptions) {
		o.TLSConfig = config
	}
}

// WithToken 设置认证Token
func WithToken(token string) ClientOption {
	return func(o *ClientOptions) {
		o.Token = token
	}
}

// WithKeepAlive 设置KeepAlive配置
func WithKeepAlive(config *keepalive.ClientParameters) ClientOption {
	return func(o *ClientOptions) {
		o.KeepAlive = config
	}
}

// WithUnaryInterceptor 添加一元拦截器
func WithUnaryInterceptor(interceptor grpc.UnaryClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.UnaryInterceptors = append(o.UnaryInterceptors, interceptor)
	}
}

// WithStreamInterceptor 添加流拦截器
func WithStreamInterceptor(interceptor grpc.StreamClientInterceptor) ClientOption {
	return func(o *ClientOptions) {
		o.StreamInterceptors = append(o.StreamInterceptors, interceptor)
	}
}
