package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

// CPUStats 存储CPU统计信息
type CPUStats struct {
	Timestamp     time.Time `json:"timestamp"`
	ProcessCPU    float64   `json:"process_cpu_percent"`  // 进程CPU使用率
	SystemCPU     float64   `json:"system_cpu_percent"`   // 系统CPU使用率
	NumCPU        int       `json:"num_cpu"`              // CPU核心数
	NumGoroutine  int       `json:"num_goroutine"`        // Goroutine数量
	ProcessMemory uint64    `json:"process_memory_bytes"` // 进程内存使用量
}

// CPUMonitor CPU监控器
type CPUMonitor struct {
	interval time.Duration
	stop     chan struct{}
	stats    []CPUStats
}

// NewCPUMonitor 创建新的CPU监控器
func NewCPUMonitor(interval time.Duration) *CPUMonitor {
	return &CPUMonitor{
		interval: interval,
		stop:     make(chan struct{}),
		stats:    make([]CPUStats, 0),
	}
}

// Start 开始监控
func (m *CPUMonitor) Start() {
	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats, err := m.collectStats()
				if err != nil {
					fmt.Printf("收集CPU统计信息失败: %v\n", err)
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
func (m *CPUMonitor) Stop() {
	close(m.stop)
}

// GetStats 获取所有统计信息
func (m *CPUMonitor) GetStats() []CPUStats {
	return m.stats
}

// SaveStats 保存统计信息到文件
func (m *CPUMonitor) SaveStats(filename string) error {
	data, err := json.MarshalIndent(m.stats, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化CPU统计信息失败: %v", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// collectStats 收集CPU统计信息
func (m *CPUMonitor) collectStats() (CPUStats, error) {
	var stats CPUStats
	stats.Timestamp = time.Now()

	// 获取进程CPU使用率
	processCPU, err := cpu.Percent(0, false)
	if err != nil {
		return stats, fmt.Errorf("获取进程CPU使用率失败: %v", err)
	}
	if len(processCPU) > 0 {
		stats.ProcessCPU = processCPU[0]
	}

	// 获取系统CPU使用率
	systemCPU, err := cpu.Percent(0, true)
	if err != nil {
		return stats, fmt.Errorf("获取系统CPU使用率失败: %v", err)
	}
	if len(systemCPU) > 0 {
		stats.SystemCPU = systemCPU[0]
	}

	// 获取CPU核心数
	stats.NumCPU = runtime.NumCPU()

	// 获取Goroutine数量
	stats.NumGoroutine = runtime.NumGoroutine()

	// 获取进程内存使用量
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	stats.ProcessMemory = mem.Alloc

	return stats, nil
}

// GetAverageStats 获取平均统计信息
func (m *CPUMonitor) GetAverageStats() CPUStats {
	if len(m.stats) == 0 {
		return CPUStats{}
	}

	var avg CPUStats
	var totalProcessCPU, totalSystemCPU float64
	var totalProcessMemory uint64

	for _, stat := range m.stats {
		totalProcessCPU += stat.ProcessCPU
		totalSystemCPU += stat.SystemCPU
		totalProcessMemory += stat.ProcessMemory
	}

	count := float64(len(m.stats))
	avg.ProcessCPU = totalProcessCPU / count
	avg.SystemCPU = totalSystemCPU / count
	avg.ProcessMemory = totalProcessMemory / uint64(count)
	avg.NumCPU = m.stats[0].NumCPU
	avg.NumGoroutine = m.stats[len(m.stats)-1].NumGoroutine
	avg.Timestamp = time.Now()

	return avg
}
