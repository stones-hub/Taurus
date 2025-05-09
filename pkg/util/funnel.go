package util

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// 漏斗模型, 用于多协程处理数据

const (
	SF_PROCESS_NUM = 30
)

type SpecialFunnel struct {
	id              string
	closeChan       chan struct{}
	dataChan        chan interface{}
	wg              *sync.WaitGroup
	handler         func(data interface{})
	tickerCloseChan chan struct{}        // 定时器关闭通道
	processedCount  int64                // 已处理的数据条数
	closed          atomic.Bool          // 标记漏斗是否关闭
	heartbeat       func(*SpecialFunnel) // 心跳函数
}

type FunnelConfig struct {
	Cap       int
	Interval  int
	Handler   func(data interface{})
	Heartbeat func(*SpecialFunnel)
}

// 创建漏斗
func NewSpecialFunnel(config *FunnelConfig) (*SpecialFunnel, func(), error) {
	f := &SpecialFunnel{
		// 生成唯一ID
		id:              uuid.NewString(),
		closeChan:       make(chan struct{}),
		dataChan:        make(chan interface{}, config.Cap),
		wg:              &sync.WaitGroup{},
		handler:         config.Handler,
		tickerCloseChan: make(chan struct{}),
		processedCount:  0,
		heartbeat:       config.Heartbeat,
	}
	f.startWorkers()
	f.startTimer(config.Interval)
	return f, f.Close, nil
}

// 启动协程
func (f *SpecialFunnel) startWorkers() {
	for i := 0; i < SF_PROCESS_NUM; i++ {
		f.wg.Add(1)
		go f.worker()
	}
}

// 每个协程的worker
func (f *SpecialFunnel) worker() {
	defer f.wg.Done()
	for {
		select {
		case data, ok := <-f.dataChan:
			if !ok {
				log.Printf("SpecialFunnel[%s] worker 管道被关闭，协程将退出。\n", f.id)
				return
			}
			f.do(data)
		case <-f.closeChan:
			// 检查下dataChan是否还有数据, 如果有数据，就继续处理，知道处理完
			for {
				// 当你调用 close(f.dataChan) 关闭通道时，即便通道里还有数据，这些数据也不会丢失，协程仍然能够把通道里剩余的数据遍历完。
				// 1. 通道一旦关闭，就不能再向其发送数据，不过可以继续从通道接收数据，直到通道里的数据被全部接收完。
				// 2. 从已关闭的通道接收数据时，若通道里还有数据，接收操作会正常返回数据和 ok 为 true；若通道里的数据已全部接收完，接收操作会返回通道元素类型的零值和 ok 为 false。
				data, ok := <-f.dataChan
				if !ok {
					return
				}
				f.do(data)
			}
		}
	}
}

// 处理数据
func (f *SpecialFunnel) do(data interface{}) {
	if f.handler == nil {
		log.Printf("SpeialFunnel[%s] handler 为空，数据未被处理: %v", f.id, data)
		return
	}
	//  sync/atomic 包提供了 AddInt64、LoadInt64 等函数。
	// 原子计数, 协程安全
	atomic.AddInt64(&f.processedCount, 1)
	// 有可能handler是阻塞的，但是不可以用协程，避免无休止的开协程
	f.handler(data)
}

// 启动定时器, 定时器每隔interval秒检查一次已处理的数据条数， 按需启用即可
func (f *SpecialFunnel) startTimer(interval int) {
	go func() {
		// 创建定时器，每10秒检查一次已处理的数据条数
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if f.heartbeat != nil {
					f.heartbeat(f)
					// log.Printf("SpecialFunnel[%s] 已处理的数据条数: %d\n", f.id, atomic.LoadInt64(&f.processedCount))
				}
			case <-f.tickerCloseChan:
				log.Printf("SpecialFunnel[%s] 定时器被关闭，协程退出。\n", f.id)
				return
			}
		}
	}()
}

// 关闭漏斗
func (f *SpecialFunnel) Close() {

	if !f.closed.CompareAndSwap(false, true) {
		log.Printf("SpecialFunnel[%s] 漏斗已关闭, 重复调用。\n", f.id)
		return
	}

	// close 是发送通知，并不会阻塞
	close(f.closeChan)          // 通知所有协程开始处理剩余数据
	close(f.dataChan)           // 关闭数据通道阻止新数据
	f.wg.Wait()                 // 等待所有协程完成, 阻塞
	close(f.tickerCloseChan)    // 关闭定时器, 不阻塞
	time.Sleep(time.Second * 1) // 等待2秒，确保所有协程都退出
	log.Printf("所有协程已退出。\n")
}

// 添加数据
func (f *SpecialFunnel) AddData(data interface{}) {
	if f.closed.Load() {
		log.Printf("SpecialFunnel[%s]  漏斗已关闭, 无法添加数据。\n", f.id)
		return
	}

	select {
	case f.dataChan <- data:
	case <-time.After(time.Second * 60): // 如果通道满了，等待5后重试
		log.Printf("SpecialFunnel[%s]数据通道已满，丢弃数据 ： %v\n", f.id, data)
	}
}

// 获取已处理的数据条数
func (f *SpecialFunnel) GetProcessedCount() int64 {
	return atomic.LoadInt64(&f.processedCount)
}
