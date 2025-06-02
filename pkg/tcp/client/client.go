package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	tcperr "Taurus/pkg/tcp/errors"
	"Taurus/pkg/tcp/protocol"
)

// TCPClientOption 定义客户端选项函数类型
type TCPClientOption func(*Client)

// WithMaxMsgSize 设置最大消息大小
func WithMaxMsgSize(size int) TCPClientOption {
	return func(c *Client) {
		c.maxMsgSize = size
	}
}

// WithBufferSize 设置缓冲区大小
func WithBufferSize(size int) TCPClientOption {
	return func(c *Client) {
		c.bufferSize = size
		c.sendChan = make(chan interface{}, size)
		c.recvBuf = make([]byte, size)
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

// Stats 统计信息
type Stats struct {
	// 消息统计
	MessagesSent     atomic.Int64
	MessagesReceived atomic.Int64
	BytesRead        atomic.Int64
	BytesWritten     atomic.Int64
	Errors           atomic.Int64
}

// NewStats 创建并初始化统计信息
func NewStats() Stats {
	return Stats{}
}

// AddMessageSent 增加发送消息计数
func (s *Stats) AddMessageSent(n int64) {
	s.MessagesSent.Add(n)
}

// AddMessageReceived 增加接收消息计数
func (s *Stats) AddMessageReceived(n int64) {
	s.MessagesReceived.Add(n)
}

// AddBytesRead 增加读取字节计数
func (s *Stats) AddBytesRead(n int64) {
	s.BytesRead.Add(n)
}

// AddBytesWritten 增加写入字节计数
func (s *Stats) AddBytesWritten(n int64) {
	s.BytesWritten.Add(n)
}

// AddError 增加错误计数
func (s *Stats) AddError(n int64) {
	s.Errors.Add(n)
}

// Client TCP客户端
type Client struct {
	// 基础配置
	address           string
	maxMsgSize        int
	bufferSize        int
	connectionTimeout time.Duration // 连接超时时间
	idleTimeout       time.Duration // 空闲超时时间
	maxRetries        int           // 最大重试次数
	baseDelay         time.Duration // 初始重试等待时间
	maxDelay          time.Duration // 最大重试等待时间

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	// 连接相关
	conn     net.Conn
	protocol protocol.Protocol
	handler  Handler

	// 消息通道
	sendChan chan interface{}
	recvBuf  []byte

	// 状态控制
	connected atomic.Bool
	closeOnce sync.Once

	// 统计信息
	stats Stats
}

// New 创建新的客户端实例
func New(address string, protocolType protocol.ProtocolType, handler Handler, opts ...TCPClientOption) (*Client, error) {
	if address == "" {
		return nil, fmt.Errorf("address cannot be empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		// 基础配置
		address:           address,
		maxMsgSize:        4096,             // 默认4KB
		bufferSize:        1024,             // 默认1KB
		connectionTimeout: 5 * time.Second,  // 默认连接超时5秒
		idleTimeout:       5 * time.Minute,  // 默认空闲超时5分钟
		maxRetries:        3,                // 默认最多重试3次
		baseDelay:         time.Second,      // 默认初始等待1秒
		maxDelay:          10 * time.Second, // 默认最大等待10秒

		// 上下文控制
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},

		// 统计信息
		stats: NewStats(),

		// 设置处理器
		handler: handler,
	}

	// 应用选项
	for _, opt := range opts {
		opt(client)
	}

	// 连接相关
	if client.handler == nil {
		client.handler = &DefaultHandler{}
	}
	p, err := protocol.NewProtocol(protocol.WithType(protocolType), protocol.WithMaxMessageSize(uint32(client.maxMsgSize)))
	if err != nil {
		return nil, fmt.Errorf("failed to create protocol: %v", err)
	}
	client.protocol = p

	// 消息通道
	client.sendChan = make(chan interface{}, client.bufferSize)
	client.recvBuf = make([]byte, client.bufferSize)

	// 状态控制
	client.connected.Store(false)

	return client, nil
}

// Connect 连接到服务器
func (c *Client) Connect() error {
	if c.connected.Load() {
		return fmt.Errorf("client already connected")
	}

	// 建立TCP连接
	dialer := &net.Dialer{Timeout: c.connectionTimeout} // 使用连接超时时间
	conn, err := dialer.DialContext(c.ctx, "tcp", c.address)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	c.conn = conn
	c.connected.Store(true)

	// 启动收发协程
	c.wg.Add(2)
	go c.readLoop()
	go c.writeLoop()

	// 通知连接建立
	c.handler.OnConnect(c.ctx, c.conn)

	return nil
}

// readLoop 读取循环
func (c *Client) readLoop() {
	defer func() {
		log.Println("readLoop closed")
		c.wg.Done()
		c.Close() // 直接调用 Close 进行清理
	}()

	// 预分配读取缓冲区, 16kB
	readBuf := make([]byte, 16*1024)
	// 消息缓冲区, 最大消息大小
	msgBuf := make([]byte, 0, c.maxMsgSize)
	// 重试相关变量
	retryCount := 0
	retryDelay := c.baseDelay

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 检查连接状态
			if !c.connected.Load() {
				return
			}

			// 检查缓冲区大小，如果过大，说明可能有大量无效数据，直接清空
			if len(msgBuf) > c.maxMsgSize {
				msgBuf = msgBuf[:0]
				c.stats.AddError(1)
				c.handler.OnError(c.ctx, c.conn, fmt.Errorf("buffer overflow"))
				continue
			}

			// 设置读取超时
			if err := c.conn.SetReadDeadline(time.Now().Add(c.idleTimeout)); err != nil {
				c.handler.OnError(c.ctx, c.conn, fmt.Errorf("set read deadline failed: %v", err))
				return
			}

			// 从连接读取数据
			n, err := c.conn.Read(readBuf)
			if err != nil {
				if err == io.EOF {
					c.handler.OnError(c.ctx, c.conn, fmt.Errorf("connection closed"))
					return
				}
				if tcperr.IsTemporaryError(err) {
					retryCount++
					c.handler.OnError(c.ctx, c.conn, fmt.Errorf("temporary read error (attempt %d/%d): %v", retryCount, c.maxRetries, err))
					if retryCount > c.maxRetries {
						// 重试次数超过最大值，关闭连接
						c.handler.OnError(c.ctx, c.conn, fmt.Errorf("max retry count exceeded"))
						return
					}

					retryDelay *= 2
					if retryDelay > c.maxDelay {
						retryDelay = c.maxDelay
					}
					// 使用指数退避策略
					time.Sleep(retryDelay)
					continue
				} else {
					c.handler.OnError(c.ctx, c.conn, fmt.Errorf("read error: %v", err))
					return
				}
			}

			// 读取成功，重置重试相关变量
			retryCount = 0
			retryDelay = c.baseDelay

			// 追加到消息缓冲区
			msgBuf = append(msgBuf, readBuf[:n]...)
			// 清空已读取的数据
			// readBuf = readBuf[:0]

			// 尝试解析消息
			message, consumed, err := c.protocol.Unpack(msgBuf)
			switch err {
			case nil:
				// 成功解析一个完整的消息
				c.stats.AddMessageReceived(1)
				c.stats.AddBytesRead(int64(consumed))
				c.handler.OnMessage(c.ctx, c.conn, message)
				// 移除已处理的数据
				msgBuf = msgBuf[consumed:]

			case tcperr.ErrShortRead:
				// 数据不足，等待更多数据
				continue

			case tcperr.ErrMessageTooLarge:
				// 消息过大，丢弃指定长度
				c.stats.AddError(1)
				c.handler.OnError(c.ctx, c.conn, err)
				if consumed > len(msgBuf) {
					msgBuf = msgBuf[:0]
				} else {
					msgBuf = msgBuf[consumed:]
				}

			case tcperr.ErrInvalidFormat:
				// 格式错误，丢弃指定长度的数据
				c.stats.AddError(1)
				c.handler.OnError(c.ctx, c.conn, err)
				msgBuf = msgBuf[consumed:]

			case tcperr.ErrChecksum:
				// 校验错误，丢弃整个包
				c.stats.AddError(1)
				c.handler.OnError(c.ctx, c.conn, err)
				msgBuf = msgBuf[consumed:]

			default:
				// 其他错误，丢弃整个包或指定长度
				c.stats.AddError(1)
				c.handler.OnError(c.ctx, c.conn, err)
				if consumed > 0 {
					msgBuf = msgBuf[consumed:]
				} else {
					msgBuf = msgBuf[:0]
				}
			}
		}
	}
}

// writeLoop 写入循环
func (c *Client) writeLoop() {
	defer func() {
		log.Println("writeLoop closed")
		c.wg.Done()
		c.Close() // 直接调用 Close 进行清理
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case data, ok := <-c.sendChan:
			if !ok {
				log.Println("sendChan closed")
				return
			}

			// 检查连接状态
			if !c.connected.Load() {
				return
			}

			// 设置写入超时
			if err := c.conn.SetWriteDeadline(time.Now().Add(c.idleTimeout)); err != nil {
				c.handler.OnError(c.ctx, c.conn, fmt.Errorf("set write deadline failed: %v", err))
				return
			}

			// 重试相关变量
			retryCount := 0
			retryDelay := c.baseDelay

			// 写入重试逻辑
			for {
				n, err := c.conn.Write(data.([]byte))
				if err != nil {
					if tcperr.IsTemporaryError(err) {
						retryCount++
						c.handler.OnError(c.ctx, c.conn, fmt.Errorf("temporary write error (attempt %d/%d): %v",
							retryCount, c.maxRetries, err))
						if retryCount > c.maxRetries {
							c.stats.AddError(1)
							return
						}
						// 使用指数退避策略
						retryDelay *= 2
						if retryDelay > c.maxDelay {
							retryDelay = c.maxDelay
						}
						time.Sleep(retryDelay)
						continue
					} else {
						c.handler.OnError(c.ctx, c.conn, fmt.Errorf("write error: %v", err))
						c.stats.AddError(1)
						return
					}
				}

				// 更新统计信息
				c.stats.AddMessageSent(1)
				c.stats.AddBytesWritten(int64(n))
				break
			}
		}
	}
}

// Send 发送消息
func (c *Client) Send(msg interface{}) error {
	// 使用 defer-recover 来处理 channel 关闭导致的 panic
	defer func() {
		if r := recover(); r != nil {
			c.handler.OnError(c.ctx, c.conn, fmt.Errorf("send message failed: %v", r))
		}
	}()

	if !c.connected.Load() {
		return fmt.Errorf("client not connected")
	}

	// 2. 打包消息
	data, err := c.protocol.Pack(msg)
	if err != nil {
		return err
	}

	select {
	case c.sendChan <- data:
		return nil
	case <-c.ctx.Done():
		return fmt.Errorf("client closed")
	default:
		return fmt.Errorf("send buffer full")
	}
}

// Close 关闭客户端连接, 回收资源
func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		// 发出关闭信号
		c.cancel()

		// 先将连接标记为关闭
		c.connected.Store(false)

		// 关闭发送通道
		close(c.sendChan)

		// 清理连接资源
		if c.conn != nil {
			_ = c.conn.Close()
			c.handler.OnClose(c.ctx, c.conn)
			c.conn = nil
		}
	})
	return nil
}

// String 返回客户端状态的字符串表示
func (c *Client) String() string {
	return fmt.Sprintf(
		"TCPClient{address: %s, connected: %v, stats: {"+
			"msgs_sent: %d, msgs_recv: %d, "+
			"bytes_read: %d, bytes_written: %d, "+
			"errors: %d}}",
		c.address,
		c.connected.Load(),
		c.stats.MessagesSent.Load(),
		c.stats.MessagesReceived.Load(),
		c.stats.BytesRead.Load(),
		c.stats.BytesWritten.Load(),
		c.stats.Errors.Load(),
	)
}

// IsConnected 检查是否已连接
func (c *Client) IsConnected() bool {
	return c.connected.Load()
}

// RemoteAddr 获取远程地址
func (c *Client) RemoteAddr() net.Addr {
	if c.conn != nil {
		return c.conn.RemoteAddr()
	}
	return nil
}

// LocalAddr 获取本地地址
func (c *Client) LocalAddr() net.Addr {
	if c.conn != nil {
		return c.conn.LocalAddr()
	}
	return nil
}
