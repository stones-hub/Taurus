package util

import (
	"log"
	"sync"
	"time"
)

// RateLimiter is a token bucket rate limiter.
type RateLimiter struct {
	capacity      int           // Maximum number of tokens
	tokens        int           // Current number of tokens
	fillInterval  time.Duration // Interval to add a token
	lastTokenTime time.Time     // Last time a token was added
	mutex         sync.Mutex    // Mutex to protect shared state
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(capacity int, fillInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity:      capacity,
		tokens:        capacity,
		fillInterval:  fillInterval,
		lastTokenTime: time.Now(),
	}
}

// Allow checks if a request can proceed. It returns true if allowed, false otherwise.
func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastTokenTime)

	// Add tokens based on elapsed time
	tokensToAdd := int(elapsed / rl.fillInterval)
	if tokensToAdd > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+tokensToAdd)
		rl.lastTokenTime = now
	}

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type CompositeRateLimiter struct {
	ipLimiters     map[string]*RateLimiter // ip限流器
	globalLimiter  *RateLimiter            // 全局限流器
	queue          []chan bool             // 没有拿到令牌的请求队列
	queueSignal    chan struct{}           // 通道用于通知队列有新请求
	ipCapacity     int                     // ip限流器容量
	globalCapacity int                     // 全局限流器容量
	mutex          sync.Mutex
}

func NewCompositeRateLimiter(ipCapacity, globalCapacity int, fillInterval time.Duration) *CompositeRateLimiter {
	compositeRateLimiter := &CompositeRateLimiter{
		ipLimiters:     make(map[string]*RateLimiter),
		globalLimiter:  NewRateLimiter(globalCapacity, fillInterval),
		ipCapacity:     ipCapacity,
		globalCapacity: globalCapacity,
		queue:          make([]chan bool, 0),
		queueSignal:    make(chan struct{}, 1), // 缓冲区为1，防止阻塞
	}

	// 启动goroutine来处理队列
	go compositeRateLimiter.processQueue()
	return compositeRateLimiter
}

func (compositeRateLimiter *CompositeRateLimiter) Allow(ip string) (bool, string) {
	compositeRateLimiter.mutex.Lock()

	// Check IP-specific limiter
	ipLimiter, exists := compositeRateLimiter.ipLimiters[ip]
	if !exists {
		ipLimiter = NewRateLimiter(compositeRateLimiter.ipCapacity, compositeRateLimiter.globalLimiter.fillInterval)
		compositeRateLimiter.ipLimiters[ip] = ipLimiter
	}

	// Check if both global and IP-specific limiters allow the request
	if compositeRateLimiter.globalLimiter.Allow() && ipLimiter.Allow() {
		compositeRateLimiter.mutex.Unlock()
		return true, ""
	}

	// If not allowed, add to queue
	log.Printf("Request from IP %s is denied and queued", ip)

	wait := make(chan bool)
	compositeRateLimiter.queue = append(compositeRateLimiter.queue, wait)

	// 通知队列有新请求
	select {
	case compositeRateLimiter.queueSignal <- struct{}{}:
	default: // 如果信号已经在通道中，则不阻塞
	}

	// Unlock before waiting to avoid holding the lock while blocked
	compositeRateLimiter.mutex.Unlock()

	// Wait for token with timeout
	select {
	case allowed := <-wait:
		return allowed, ""
	case <-time.After(5 * time.Second): // 5秒超时
		return false, "请求超时，请稍后再试！"
	}
}

func (compositeRateLimiter *CompositeRateLimiter) processQueue() {
	for range compositeRateLimiter.queueSignal { // 使用for range来监听信号
		compositeRateLimiter.mutex.Lock()
		for len(compositeRateLimiter.queue) > 0 && compositeRateLimiter.globalLimiter.Allow() {
			wait := compositeRateLimiter.queue[0]
			compositeRateLimiter.queue = compositeRateLimiter.queue[1:]
			wait <- true
			close(wait)
		}
		compositeRateLimiter.mutex.Unlock()
	}
}

// ------------------  例子 ------------------
/*

func main() {
	// 创建一个 CompositeRateLimiter 实例
	limiter := util.NewCompositeRateLimiter(1, 5, time.Minute) // 每分钟允许5个请求, 每个IP允许1个请求

	// 创建一个新的 HTTP 服务器
	mux := http.NewServeMux()

	// 使用 RateLimitMiddleware 包装处理器
	mux.Handle("/", middleware.RateLimitMiddleware(http.HandlerFunc(helloHandler), limiter))

	// 启动服务器
	http.ListenAndServe(":8080", mux)
}

*/
