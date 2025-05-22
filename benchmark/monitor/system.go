package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/net"
)

// SystemStats 存储系统统计信息
type SystemStats struct {
	Timestamp      time.Time `json:"timestamp"`
	LoadAvg1       float64   `json:"load_avg_1"`       // 1分钟负载
	LoadAvg5       float64   `json:"load_avg_5"`       // 5分钟负载
	LoadAvg15      float64   `json:"load_avg_15"`      // 15分钟负载
	DiskReadBytes  uint64    `json:"disk_read_bytes"`  // 磁盘读取字节数
	DiskWriteBytes uint64    `json:"disk_write_bytes"` // 磁盘写入字节数
	NetRecvBytes   uint64    `json:"net_recv_bytes"`   // 网络接收字节数
	NetSentBytes   uint64    `json:"net_sent_bytes"`   // 网络发送字节数
	NumGC          uint32    `json:"num_gc"`           // GC次数
	PauseTotalNs   uint64    `json:"pause_total_ns"`   // GC暂停总时间
	LastGC         uint64    `json:"last_gc"`          // 上次GC时间
}

// SystemMonitor 系统监控器
type SystemMonitor struct {
	interval time.Duration
	stop     chan struct{}
	stats    []SystemStats
}

// NewSystemMonitor 创建新的系统监控器
func NewSystemMonitor(interval time.Duration) *SystemMonitor {
	return &SystemMonitor{
		interval: interval,
		stop:     make(chan struct{}),
		stats:    make([]SystemStats, 0),
	}
}

// Start 开始监控
func (m *SystemMonitor) Start() {
	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats, err := m.collectStats()
				if err != nil {
					fmt.Printf("收集系统统计信息失败: %v\n", err)
					continue
				}
				m.stats = append(m.stats, stats)
			case <-m.stop:
				return
			}
		}
	}()
}

// Stop 停止监控
func (m *SystemMonitor) Stop() {
	close(m.stop)
}

// GetStats 获取所有统计信息
func (m *SystemMonitor) GetStats() []SystemStats {
	return m.stats
}

// SaveStats 保存统计信息到文件
func (m *SystemMonitor) SaveStats(filename string) error {
	data, err := json.MarshalIndent(m.stats, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化系统统计信息失败: %v", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// collectStats 收集系统统计信息
func (m *SystemMonitor) collectStats() (SystemStats, error) {
	var stats SystemStats
	stats.Timestamp = time.Now()

	// 获取系统负载
	loadAvg, err := load.Avg()
	if err != nil {
		return stats, fmt.Errorf("获取系统负载失败: %v", err)
	}
	stats.LoadAvg1 = loadAvg.Load1
	stats.LoadAvg5 = loadAvg.Load5
	stats.LoadAvg15 = loadAvg.Load15

	// 获取磁盘I/O
	diskIO, err := disk.IOCounters()
	if err != nil {
		return stats, fmt.Errorf("获取磁盘I/O失败: %v", err)
	}
	for _, io := range diskIO {
		stats.DiskReadBytes += io.ReadBytes
		stats.DiskWriteBytes += io.WriteBytes
	}

	// 获取网络I/O
	netIO, err := net.IOCounters(true)
	if err != nil {
		return stats, fmt.Errorf("获取网络I/O失败: %v", err)
	}
	for _, io := range netIO {
		stats.NetRecvBytes += io.BytesRecv
		stats.NetSentBytes += io.BytesSent
	}

	// 获取GC信息
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	stats.NumGC = mem.NumGC
	stats.PauseTotalNs = mem.PauseTotalNs
	stats.LastGC = mem.LastGC

	return stats, nil
}

// GetAverageStats 获取平均统计信息
func (m *SystemMonitor) GetAverageStats() SystemStats {
	if len(m.stats) == 0 {
		return SystemStats{}
	}

	var avg SystemStats
	var totalLoad1, totalLoad5, totalLoad15 float64
	var totalDiskRead, totalDiskWrite uint64
	var totalNetRecv, totalNetSent uint64
	var totalNumGC uint32
	var totalPauseTotalNs uint64

	for _, stat := range m.stats {
		totalLoad1 += stat.LoadAvg1
		totalLoad5 += stat.LoadAvg5
		totalLoad15 += stat.LoadAvg15
		totalDiskRead += stat.DiskReadBytes
		totalDiskWrite += stat.DiskWriteBytes
		totalNetRecv += stat.NetRecvBytes
		totalNetSent += stat.NetSentBytes
		totalNumGC += stat.NumGC
		totalPauseTotalNs += stat.PauseTotalNs
	}

	count := float64(len(m.stats))
	avg.LoadAvg1 = totalLoad1 / count
	avg.LoadAvg5 = totalLoad5 / count
	avg.LoadAvg15 = totalLoad15 / count
	avg.DiskReadBytes = totalDiskRead / uint64(count)
	avg.DiskWriteBytes = totalDiskWrite / uint64(count)
	avg.NetRecvBytes = totalNetRecv / uint64(count)
	avg.NetSentBytes = totalNetSent / uint64(count)
	avg.NumGC = totalNumGC / uint32(count)
	avg.PauseTotalNs = totalPauseTotalNs / uint64(count)
	avg.LastGC = m.stats[len(m.stats)-1].LastGC
	avg.Timestamp = time.Now()

	return avg
}
