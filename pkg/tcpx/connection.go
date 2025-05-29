package tcpx

import (
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

// WithMaxMessageSize 设置最大消息大小
func WithMaxMessageSize(bytes int) ConnectionOption {
	return func(c *Connection) {
		c.maxMessageSize = bytes
	}
}

// WithBandwidthLimit 设置带宽限制
func WithBandwidthLimit(bytesPerSecond float64) ConnectionOption {
	return func(c *Connection) {
		c.bandwidth = rate.NewLimiter(rate.Limit(bytesPerSecond), 1)
	}
}

// Connection 表示单个 TCP 连接并管理其生命周期。
// 它处理消息读写、流量控制和资源管理。
type Connection struct {
	id       uint64             // 唯一连接标识符
	conn     net.Conn           // 底层 TCP 连接
	protocol protocol.Protocol  // 消息编码解码协议
	handler  Handler            // 连接事件处理器
	sendChan chan []byte        // 异步消息发送通道
	ctx      context.Context    // 生命周期管理上下文
	cancel   context.CancelFunc // 取消上下文的函数
	metrics  *Metrics           // 连接统计指标

	// 空闲超时管理
	idleTimeout time.Duration // 最大空闲时间

	// 流量控制
	maxMessageSize int           // 最大允许消息大小
	rateLimiter    *rate.Limiter // 消息频率限制器
	bandwidth      *rate.Limiter // 带宽使用限制器

	lastActiveTime atomic.Value // 连接最后活动时间戳

	waitGroup sync.WaitGroup // goroutine 同步等待组
	closed    int32          // 连接状态的原子标志
	attrs     sync.Map       // 线程安全的属性存储
	once      sync.Once      // 确保清理只执行一次
}

var globalConnectionID uint64 // 生成唯一连接 ID 的全局计数器

// NewConnection 创建一个新的连接实例。
// 它接受可选的配置选项，使用函数式选项模式进行配置。
func NewConnection(conn net.Conn, protocol protocol.Protocol, handler Handler, opts ...ConnectionOption) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Connection{
		id:       atomic.AddUint64(&globalConnectionID, 1), // 协程安全的，每次给globalConnectionID累加1（全局变量）, 并返回累加值, 作为连接id (唯一)
		conn:     conn,                                     // 底层tcp连接
		protocol: protocol,                                 // 消息编码解码协议
		handler:  handler,                                  // 连接处理器
		sendChan: make(chan []byte, 1024),                  // 异步消息发送通道
		ctx:      ctx,                                      // 生命周期管理上下文
		cancel:   cancel,                                   // 取消上下文的函数
		metrics:  NewMetrics(),                             // 连接层面的统计指标

		// 默认配置
		idleTimeout:    time.Minute * 5, // 默认 5 分钟超时, 具体指当连接消息的收发超过5分钟没有动静, 则认为连接已经死亡, 需要主动关闭
		maxMessageSize: 1024 * 1024,     // 默认最大消息大小 1MB, 如果消息大小超过这个值, 则认为消息太大, 需要丢弃

		// 这里用到了rate.NewLimiter(r Limit, b int) 来创建限流器, 限流器是golang.org/x/time/rate包中的一个限流器, 用于限制消息的速率, 限流器有三个主要方法, Allow, Reserve, Wait
		// 1. Allow: 判断是否能拿到令牌, 拿到了返回true，拿不到返回false
		// 2. Wait(ctx context.Context): 等待一个令牌(阻塞方法), 如果等待成功, 则返回nil, 否则返回一个error, 如果 context 被取消，返回 context 的错误
		// 3. Reserve: 希望能预约一个令牌，如果预约成功，则返回一个Reservation对象，否则返回一个error, Reservation中存储了令牌的可用时间, 适用于需要预判等待时间获取令牌的场景
		/*
			// 使用 Reserve
			r := limiter.Reserve() // 预约令牌
			if !r.OK() { // 如果预约失败，则返回一个error
			    return errors.New("rate limit exceeded")
			}
			// r.Delay() 返回的是令牌可用需要等待的时间
			if r.Delay() > 5*time.Second { // 判断等待时间是否太长，如果太长，则取消预约
			    r.Cancel() // 等待时间太长，取消预约
			    return errors.New("wait too long")
			}
			time.Sleep(r.Delay())
		*/
		// 对于参数解释：
		// r(rate.Limit(100)) 表示每秒允许通过100个令牌
		// b(1) 表示令牌桶的容量为1, 如果令牌桶满了, 则新的令牌会丢弃
		// 这样就意味着，即使每秒生成100个令牌，但是令牌桶的容量为1，所以每秒最多只能处理1个令牌, 因为能不能处理取决于令牌桶中能不能拿到令牌

		rateLimiter: rate.NewLimiter(rate.Limit(100), 100),             // 默认 每秒最多只能处理 100 条消息
		bandwidth:   rate.NewLimiter(rate.Limit(1024*1024), 1024*1024), // 默认 每秒最多只能处理 1MB 的数据
	}

	// 应用所有配置选项
	for _, opt := range opts {
		opt(c)
	}

	c.lastActiveTime.Store(time.Now()) // 初始化设置连接的最后活动时间
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
}

// readLoop 持续读取和处理传入消息。
// 它处理流量控制、消息验证和错误处理。
func (c *Connection) readLoop() {
	defer func() {
		c.waitGroup.Done()
		c.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 设置从连接中读取消息的最大等待时间（超过60s还未读取到数据，就返回读取超时）
			_ = c.conn.SetReadDeadline(time.Now().Add(time.Second * 60))

			// 读取消息前检查带宽限制
			if err := c.bandwidth.Wait(c.ctx); err != nil {
				c.handler.OnError(c, fmt.Errorf("bandwidth limit: %w", err))
				continue
			}

			// 记录实际消息处理的开始时间
			start := time.Now()
			// 读取连接中的数据，并解包(读取连接或解包都可能出错，错误返回)
			message, err := c.protocol.Unpack(c.conn)
			if err != nil {
				c.metrics.AddError()
				// 错误如果是连接本身的错误，则需要关闭连接, 其他的错误我们可以等待下一次读取,  并将当前的错误返回给上层
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					c.handler.OnError(c, ErrConnectionClosed)
					// 关闭连接
					return
				}
				continue
			}

			// 消息大小验证
			if data, ok := message.([]byte); ok {
				if len(data) > c.maxMessageSize { // 超出单条消息的大小限制
					c.handler.OnError(c, ErrMessageTooLarge)
					c.metrics.AddError()
					continue
				}

				// 更新接收字节数统计
				c.metrics.AddMessageReceived(int64(len(data)))

				// 速率限制（基于消息大小）
				if err := c.rateLimiter.WaitN(c.ctx, len(data)); err != nil {
					c.handler.OnError(c, fmt.Errorf("rate limit: %w", err))
					continue
				}
			}

			// 更新活跃时间和延迟指标
			c.updateActiveTime()
			c.metrics.SetMessageLatency(time.Since(start))

			// 处理消息
			c.handler.OnMessage(c, message)
		}
	}
}

// writeLoop 处理发送消息。
// 它应用流量控制并监控写操作。
func (c *Connection) writeLoop() {
	defer func() {
		c.waitGroup.Done()
		c.Close()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		case data := <-c.sendChan:
			start := time.Now()
			_ = c.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))

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
	defer c.waitGroup.Done()

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			lastActive := c.lastActiveTime.Load().(time.Time)
			if time.Since(lastActive) > c.idleTimeout {
				c.handler.OnError(c, ErrConnectionIdle)
				c.Close()
				return
			}
		}
	}
}

// Send 将消息排队等待发送。
// 它在排队前应用流量控制和大小验证。
func (c *Connection) Send(message interface{}) error {
	if atomic.LoadInt32(&c.closed) == 1 {
		return ErrConnectionClosed
	}

	data, err := c.protocol.Pack(message)
	if err != nil {
		return err
	}

	if len(data) > c.maxMessageSize {
		return ErrMessageTooLarge
	}

	if err := c.rateLimiter.Wait(c.ctx); err != nil {
		return err
	}

	if err := c.bandwidth.Wait(c.ctx); err != nil {
		return err
	}

	select {
	case c.sendChan <- data:
		return nil
	default:
		return ErrSendChannelFull
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

func (c *Connection) SetAttr(key, value interface{}) {
	c.attrs.Store(key, value)
}

func (c *Connection) GetAttr(key interface{}) (interface{}, bool) {
	return c.attrs.Load(key)
}

func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Connection) GetMetrics() map[string]interface{} {
	return c.metrics.GetStats()
}

func (c *Connection) updateActiveTime() {
	c.lastActiveTime.Store(time.Now())
}
