package tcpx

import (
	"Taurus/pkg/tcpx/errors"
	"Taurus/pkg/tcpx/protocol"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// Handler 定义了连接事件处理的接口。
// 实现者需要处理各种连接生命周期事件。
type Handler interface {
	OnConnect(conn *Connection)                      // 当新连接建立时调用
	OnMessage(conn *Connection, message interface{}) // 当收到消息时调用
	OnClose(conn *Connection)                        // 当连接关闭时调用
	OnError(conn *Connection, err error)             // 当发生错误时调用
}

// ConnectionOption 定义了配置连接的函数类型
type ConnectionOption func(*Connection)

// WithSendChanSize 设置消息发送通道的大小
func WithSendChanSize(size int) ConnectionOption {
	return func(c *Connection) {
		c.sendChan = make(chan []byte, size)
	}
}

// WithIdleTimeout 设置连接的空闲超时时间
func WithIdleTimeout(timeout time.Duration) ConnectionOption {
	return func(c *Connection) {
		c.idleTimeout = timeout
	}
}

// WithRateLimit 设置消息速率限制
func WithRateLimit(messagesPerSecond float64) ConnectionOption {
	return func(c *Connection) {
		c.rateLimiter = rate.NewLimiter(rate.Limit(messagesPerSecond), 1)
	}
}

// WithMaxMessageSize 设置连接允许传输的最大消息大小
func WithMaxMessageSize(bytes int) ConnectionOption {
	return func(c *Connection) {
		c.maxMessageSize = bytes
	}
}

// Connection 表示单个 TCP 连接并管理其生命周期。
// 它处理消息读写、流量控制和资源管理。
type Connection struct {
	id             uint64             // 唯一连接标识符
	conn           net.Conn           // 底层 TCP 连接
	protocol       protocol.Protocol  // 消息编码解码协议
	handler        Handler            // 连接事件处理器
	ctx            context.Context    // 生命周期管理上下文
	cancel         context.CancelFunc // 取消上下文的函数
	metrics        *Metrics           // 连接统计指标
	lastActiveTime atomic.Value       // 连接最后活动时间戳
	closed         int32              // 连接状态的原子标志, 0: 连接正常, 1: 连接已关闭，
	attrs          sync.Map           // 线程安全的属性存储, 用于存储连接的属性
	once           sync.Once          // 确保清理只执行一次
	waitGroup      *sync.WaitGroup    // goroutine 同步等待组
	// 重试相关配置，用于解决从连接获取数据时，可能出现的临时错误
	maxRetryCount  int           // 最大重试次数, 默认3次
	baseRetryDelay time.Duration // 基础重试延迟, 默认1秒
	maxRetryDelay  time.Duration // 最大重试延迟, 默认10秒

	// 默认配置
	sendChan       chan []byte   // 异步消息发送通道
	idleTimeout    time.Duration // 连接最大空闲超时时间
	maxMessageSize int           // 连接允许传输的最大消息大小
	rateLimiter    *rate.Limiter // 消息频率限制器
}

var globalConnectionID uint64 // 生成唯一连接 ID 的全局计数器

// NewConnection 创建一个新的连接实例。
// 它接受可选的配置选项，使用函数式选项模式进行配置。
func NewConnection(conn net.Conn, protocol protocol.Protocol, handler Handler, opts ...ConnectionOption) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Connection{
		id:        atomic.AddUint64(&globalConnectionID, 1),
		conn:      conn,
		protocol:  protocol,
		handler:   handler,
		ctx:       ctx,
		cancel:    cancel,
		metrics:   NewMetrics(),
		closed:    0,
		attrs:     sync.Map{},
		once:      sync.Once{},
		waitGroup: &sync.WaitGroup{},

		// 重试相关配置, 默认3次, 1秒, 10秒
		maxRetryCount:  3,
		baseRetryDelay: 1 * time.Second,
		maxRetryDelay:  10 * time.Second,

		// 默认配置, 可以被配置选项覆盖
		sendChan:       make(chan []byte, 1024),               // 消息发送通道
		idleTimeout:    time.Minute * 5,                       // 默认5分钟空闲超时
		rateLimiter:    rate.NewLimiter(rate.Limit(100), 100), // 默认每秒100条消息
		maxMessageSize: 10 * 1024 * 1024,                      // 默认最大消息大小10MB
	}

	// 应用所有配置选项
	for _, opt := range opts {
		opt(c)
	}

	c.lastActiveTime.Store(time.Now())
	return c
}

// Start 启动连接的读写循环。
// 它启动独立的 goroutine 用于读取、写入和空闲检查。
func (c *Connection) Start() {
	c.waitGroup.Add(3)
	go c.readLoop()
	go c.writeLoop()
	go c.checkIdleLoop()
	c.handler.OnConnect(c)
	c.waitGroup.Wait()
}

// readLoop 持续读取和处理传入消息。
// 它处理流量控制、消息验证和错误处理。
func (c *Connection) readLoop() {
	defer func() {
		c.waitGroup.Done()
	}()

	// 预分配读取缓冲区, 16kB
	readBuf := make([]byte, 16*1024)
	// 消息缓冲区, 最大消息大小
	msgBuf := make([]byte, 0, c.maxMessageSize)
	// 重试相关变量, 用于处理临时错误
	retryCount := 0
	retryDelay := c.baseRetryDelay

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 判断是否达到发送速率限制
			if !c.rateLimiter.Allow() {
				// 如果达到发送速率限制，等待100ms
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// 1. 设置读取超时
			if err := c.conn.SetReadDeadline(time.Now().Add(c.idleTimeout)); err != nil {
				c.handler.OnError(c, fmt.Errorf("set read deadline failed: %w", err))
				return
			}

			// 检查缓冲区大小，如果过大，说明可能有大量无效数据，直接清空
			if len(msgBuf) > c.maxMessageSize {
				// 缓冲区过大，说明可能有大量无效数据，直接清空
				msgBuf = msgBuf[:0]
				c.metrics.AddError()
				c.handler.OnError(c, errors.ErrBufferOverflow)
				continue
			}

			// 2. 检查连接状态
			if atomic.LoadInt32(&c.closed) == 1 {
				return
			}

			// 3. 从连接读取数据
			n, err := c.conn.Read(readBuf)
			if err != nil {
				if err == io.EOF {
					c.handler.OnError(c, errors.ErrConnectionClosed)
					return
				}
				if errors.IsTemporaryError(err) {
					retryCount++
					c.handler.OnError(c, fmt.Errorf("temporary read error (attempt %d/%d): %w", retryCount, c.maxRetryCount, err))
					if retryCount > c.maxRetryCount {
						// 重试次数超过最大值，关闭连接
						c.handler.OnError(c, fmt.Errorf("max retry count exceeded: %w", err))
						return
					}
					// 使用指数退避策略计算下一次重试延迟
					retryDelay *= 2
					if retryDelay > c.maxRetryDelay {
						retryDelay = c.maxRetryDelay
					}
					time.Sleep(retryDelay)
					continue
				}
				// 非临时错误，直接关闭连接
				c.handler.OnError(c, fmt.Errorf("read error: %w", err))
				return
			}

			// 读取成功，重置重试相关变量
			retryCount = 0
			retryDelay = c.baseRetryDelay

			// 4. 追加到消息缓冲区
			msgBuf = append(msgBuf, readBuf[:n]...)
			// 清空已读取的数据, 以防万一
			readBuf = readBuf[:0]

			// 5. 尝试解析一个完整的消息
			start := time.Now()
			message, consumed, err := c.protocol.Unpack(msgBuf)

			// 6. 处理不同的错误情况
			switch err {
			case nil:
				// 成功解析一个完整的消息
				// 更新接收的消息数量
				c.metrics.AddMessageReceived(int64(consumed))
				// 设置消息最后处理时间
				c.metrics.SetMessageLatency(time.Since(start))
				// 更新连接最后活动时间
				c.updateActiveTime()

				// 处理消息
				c.handler.OnMessage(c, message)

				// 移除已处理的数据
				msgBuf = msgBuf[consumed:]

			case errors.ErrShortRead:
				// 数据不足，保留所有数据等待更多数据
				continue

			case errors.ErrMessageTooLarge:
				// 消息过大，丢弃指定长度
				c.metrics.AddError()
				c.handler.OnError(c, err)
				// 如果返回的consumed大于当前数据，说明需要丢弃所有数据
				if consumed > len(msgBuf) {
					msgBuf = msgBuf[:0]
				} else {
					msgBuf = msgBuf[consumed:]
				}

			case errors.ErrInvalidFormat:
				// 格式错误（比如魔数不在开头），丢弃指定长度的数据
				c.metrics.AddError()
				c.handler.OnError(c, err)
				// consumed表示魔数之前的数据长度，直接丢弃
				msgBuf = msgBuf[consumed:]

			case errors.ErrChecksum:
				// 校验错误，丢弃整个包
				c.metrics.AddError()
				c.handler.OnError(c, err)
				msgBuf = msgBuf[consumed:]

			default:
				// 其他错误（比如JSON解析错误），丢弃整个包
				c.metrics.AddError()
				c.handler.OnError(c, err)
				if consumed > 0 {
					msgBuf = msgBuf[consumed:]
				} else {
					// 无法恢复的错误，且没有指定丢弃长度，丢弃所有数据
					msgBuf = msgBuf[:0]
				}
			}
		}
	}
}

// writeLoop 处理发送消息。
func (c *Connection) writeLoop() {
	defer func() {
		c.waitGroup.Done()
	}()

	for {
		// 判断是否达到发送速率限制
		if !c.rateLimiter.Allow() {
			// 如果达到发送速率限制，等待100ms
			time.Sleep(100 * time.Millisecond)
			continue
		}

		select {
		case <-c.ctx.Done():
			return
		case data := <-c.sendChan:

			start := time.Now()
			err := c.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			if err != nil {
				c.handler.OnError(c, fmt.Errorf("set write deadline failed: %w", err))
				return
			}

			// 检查连接状态，如果已关闭，直接返回
			if atomic.LoadInt32(&c.closed) == 1 {
				return
			}

			n, err := c.conn.Write(data)
			if err != nil {
				c.handler.OnError(c, err)
				c.metrics.AddError()
				return
			}

			c.updateActiveTime()
			c.metrics.AddMessageSent(int64(n))
			c.metrics.SetMessageLatency(time.Since(start))
		}
	}
}

// checkIdleLoop 监控连接活动并关闭空闲连接。
func (c *Connection) checkIdleLoop() {
	defer func() {
		c.waitGroup.Done()
	}()

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			lastActive := c.lastActiveTime.Load().(time.Time)
			if time.Since(lastActive) > c.idleTimeout {
				c.handler.OnError(c, errors.ErrConnectionIdle)
				return
			}
		}
	}
}

// Send 将消息写入发送队列，等待发送。
func (c *Connection) Send(message interface{}) error {
	if atomic.LoadInt32(&c.closed) == 1 {
		return errors.ErrConnectionClosed
	}

	data, err := c.protocol.Pack(message)
	if err != nil {
		return err
	}

	select {
	case c.sendChan <- data:
		return nil
	default:
		return errors.ErrSendChannelFull
	}
}

// Close 优雅地关闭连接。
// 它确保清理只发生一次并等待所有 goroutine 完成。
func (c *Connection) Close() {
	c.once.Do(func() {
		if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
			return
		}
		c.cancel()
		_ = c.conn.Close()
		c.handler.OnClose(c)
		c.waitGroup.Wait()
	})
}

// 工具方法
func (c *Connection) ID() uint64 {
	return c.id
}

// SetAttr 设连接属性
func (c *Connection) SetAttr(key, value interface{}) {
	c.attrs.Store(key, value)
}

// GetAttr 获取连接属性
func (c *Connection) GetAttr(key interface{}) (interface{}, bool) {
	return c.attrs.Load(key)
}

// RemoteAddr 获取远程地址
func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

// LocalAddr 获取本地地址
func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// GetMetrics 获取连接的统计指标
func (c *Connection) GetMetrics() map[string]interface{} {
	return c.metrics.GetStats()
}

// updateActiveTime 更新连接最后活动时间
func (c *Connection) updateActiveTime() {
	c.lastActiveTime.Store(time.Now())
}
