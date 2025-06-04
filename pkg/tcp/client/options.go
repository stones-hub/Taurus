package client

import "time"

// TCPClientOption 定义客户端选项函数类型
type TCPClientOption func(*Client)

// WithMaxMsgSize 设置最大消息大小
func WithMaxMsgSize(size uint32) TCPClientOption {
	return func(c *Client) {
		c.maxMsgSize = size
	}
}

// WithBufferSize 设置缓冲区大小
func WithBufferSize(size int) TCPClientOption {
	return func(c *Client) {
		c.bufferSize = size
		c.sendChan = make(chan interface{}, size)
	}
}

// WithConnectionTimeout 设置连接超时时间
func WithConnectionTimeout(timeout time.Duration) TCPClientOption {
	return func(c *Client) {
		c.connectionTimeout = timeout
	}
}

// WithIdleTimeout 设置空闲超时时间
func WithIdleTimeout(timeout time.Duration) TCPClientOption {
	return func(c *Client) {
		c.idleTimeout = timeout
	}
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) TCPClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithBaseRetryDelay 设置初始重试等待时间
func WithBaseRetryDelay(baseDelay time.Duration) TCPClientOption {
	return func(c *Client) {
		c.baseDelay = baseDelay
	}
}

// WithMaxRetryDelay 设置最大重试等待时间
func WithMaxRetryDelay(maxDelay time.Duration) TCPClientOption {
	return func(c *Client) {
		c.maxDelay = maxDelay
	}
}
