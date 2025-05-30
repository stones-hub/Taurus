package tcp

import (
	"Taurus/pkg/tcp/errors"
	"Taurus/pkg/tcp/protocol"
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Server 表示一个处理多个客户端连接的 TCP 服务器。
// 它管理连接生命周期、执行资源限制并提供监控功能。
type Server struct {
	started    int32              // 防止多次启动的原子标志
	addr       string             // 网络监听地址
	ctx        context.Context    // 生命周期管理的上下文
	cancel     context.CancelFunc // 取消上下文的函数
	baseDelay  time.Duration      // 初始重试延迟时间
	maxDelay   time.Duration      // 最大重试延迟时间
	maxRetries int                // 最大重试次数
	conns      sync.Map           // 线程安全的连接存储
	wg         *sync.WaitGroup    // 优雅关闭的等待组
	metrics    *Metrics           // 服务器指标收集器
	listener   net.Listener       // TCP 监听器

	// 默认配置, 可以被配置选项覆盖
	protocol protocol.Protocol // 消息处理的协议实现
	handler  Handler           // 业务逻辑处理器
	maxConns int32             // 最大并发连接数
	connChan chan struct{}     // 连接限制信号量
}

// ServerOption 定义了配置服务器的函数类型。
// 这遵循函数式选项模式以实现灵活配置。
type ServerOption func(*Server)

// WithProtocol 设置服务器的协议实现。
// 协议定义了消息如何编码和解码。
func WithProtocol(protocol protocol.Protocol) ServerOption {
	return func(s *Server) {
		s.protocol = protocol
	}
}

// WithHandler 设置服务器的消息处理器。
// 处理器实现了消息处理的业务逻辑。
func WithHandler(handler Handler) ServerOption {
	return func(s *Server) {
		s.handler = handler
	}
}

// WithMaxConnections 设置最大并发连接数。
// 同时初始化连接信号量通道。
func WithMaxConnections(maxConns int32) ServerOption {
	return func(s *Server) {
		s.maxConns = maxConns
		s.connChan = make(chan struct{}, maxConns)
	}
}

// NewServer 创建一个新的 TCP 服务器实例。
// 使用默认值初始化服务器并应用提供的选项。
func NewServer(addr string, opts ...ServerOption) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{
		addr:       addr,
		ctx:        ctx,
		cancel:     cancel,
		baseDelay:  1 * time.Second,   // 初始重试延迟时间
		maxDelay:   10 * time.Second,  // 最大重试延迟时间
		maxRetries: 3,                 // 最大重试次数
		conns:      sync.Map{},        // 线程安全的连接存储
		wg:         &sync.WaitGroup{}, // 优雅关闭的等待组
		metrics:    NewMetrics(),      // 服务器层面的统计指标

		maxConns: 1000, // 默认最大连接数
	}

	// 应用所有配置选项
	for _, opt := range opts {
		opt(s)
	}

	// 如果未设置则初始化连接信号量
	if s.connChan == nil {
		s.connChan = make(chan struct{}, s.maxConns)
	}

	return s
}

// Start 开始接受客户端连接。
// 确保服务器只启动一次并具备所需组件。
func (s *Server) Start() error {
	// 使用原子操作确保单次启动
	if !atomic.CompareAndSwapInt32(&s.started, 0, 1) {
		return errors.ErrServerAlreadyStarted
	}

	// 验证必需组件
	if s.protocol == nil {
		return errors.ErrProtocolNotSet
	}
	if s.handler == nil {
		return errors.ErrHandlerNotSet
	}

	// 创建 TCP 监听器
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeSystem, err, "listen failed")
	}
	s.listener = listener

	s.wg.Add(1)
	// 开协程在后台开始接受连接, 协程的退出，不会清理服务器资源，只是退出协程
	go s.acceptLoop()

	// 阻塞等待所有协程退出
	s.wg.Wait()
	log.Println("server stopped")
	return nil
}

// acceptLoop 在独立的 goroutine 中运行并处理传入连接。
// 它实现了连接限制、错误处理和重试机制。
func (s *Server) acceptLoop() {
	defer func() {
		// 如果acceptLoop退出，说明服务器已经关闭，需要清理服务器资源
		s.Stop()
		s.wg.Done()
		log.Println("acceptLoop exited")
	}()

	retries := 0         // 当前重试次数
	delay := s.baseDelay // 当前重试延迟时间

	for {
		// 先检查服务器状态
		if atomic.LoadInt32(&s.started) == 0 {
			return
		}

		// 先检查是否有可用槽位和服务状态
		select {
		case s.connChan <- struct{}{}: // 获取到槽位
			// 获取到槽位后接受连接, 有可能阻塞，不过没事，因为你已经获取到槽位了
			conn, err := s.listener.Accept()
			if err != nil {
				<-s.connChan // 释放槽位
				// 如果服务器已关闭，直接返回
				if atomic.LoadInt32(&s.started) == 0 {
					return
				}

				// 判断错误是否为临时性的
				if errors.IsTemporaryError(err) {
					if retries < s.maxRetries {
						retries++
						s.metrics.AddError()
						s.handler.OnError(nil, errors.ErrSystemOverload)
						time.Sleep(delay)
						delay *= 2
						if delay > s.maxDelay {
							delay = s.maxDelay
						}
						continue
					}
					// 重试次数达到最大值，退出协程
					s.metrics.AddError()
					s.handler.OnError(nil, errors.ErrSystemFatal)
					return
				}
				// 非临时性错误，退出协程
				s.metrics.AddError()
				s.handler.OnError(nil, errors.ErrSystemFatal)
				return
			}

			// 重置重试相关计数
			retries = 0
			delay = s.baseDelay

			// 创建并存储新连接
			c := NewConnection(conn, s.protocol, s.handler)
			s.conns.Store(c.ID(), c)
			s.metrics.AddConnection()

			s.wg.Add(1)
			// 启动协程处理单个连接
			go func() {
				// 当处理连接的协程退出，只需要清理当前连接的资源
				defer func() {
					s.conns.Delete(c.ID())
					s.metrics.RemoveConnection()
					<-s.connChan // 释放连接槽
					s.wg.Done()
				}()
				// Connection.Start()内部会处理连接的关闭
				c.Start()
			}()

		case <-s.ctx.Done(): // 服务器正在关闭
			return

		default: // 没有可用槽位，等待一会再试
			s.metrics.AddConnectionRefused()
			s.handler.OnError(nil, errors.ErrTooManyConnections)
			time.Sleep(time.Millisecond * 100) // 避免空转
		}
	}
}

// Stop 优雅地关闭服务器。
// 停止接受新连接并关闭现有连接。
func (s *Server) Stop() {
	// 1. 先将服务器标记为已关闭
	if !atomic.CompareAndSwapInt32(&s.started, 1, 0) {
		log.Println("server is already stopping")
		return
	}

	// 2. 关闭监听器，停止接受新连接
	if s.listener != nil {
		s.listener.Close()
	}

	// 3. 关闭所有现有连接, 确保所有连接资源能关闭
	s.conns.Range(func(key, value interface{}) bool {
		if conn, ok := value.(*Connection); ok {
			// 关闭连接会同时关闭socket和触发清理流程
			conn.Close()
		}
		return true
	})

	// 4. 取消上下文，确保所有协程都收到退出信号
	s.cancel()

	log.Println("server stopped completely")
}

// GetConnection 根据 ID 获取连接
func (s *Server) GetConnection(id uint64) (*Connection, bool) {
	if value, ok := s.conns.Load(id); ok {
		return value.(*Connection), true
	}
	return nil, false
}

// Broadcast 向所有连接广播消息
func (s *Server) Broadcast(message interface{}) {
	s.conns.Range(func(key, value interface{}) bool {
		conn := value.(*Connection)
		_ = conn.Send(message)
		return true
	})
}

// ConnectionCount 获取当前连接数
func (s *Server) ConnectionCount() int32 {
	return s.maxConns - int32(len(s.connChan))
}

// GetMetrics 返回当前服务器指标。
func (s *Server) GetMetrics() map[string]interface{} {
	return s.metrics.GetStats()
}
