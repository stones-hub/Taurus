package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"Taurus/test/load"
	"Taurus/test/monitor"
	"Taurus/test/reports"
)

var (
	url                = flag.String("url", "http://127.0.0.1:9080/v1/api/", "测试目标URL")
	duration           = flag.Duration("duration", 5*time.Minute, "测试持续时间")
	initialConcurrency = flag.Int("initial-concurrency", 10, "初始并发数")
	peakConcurrency    = flag.Int("peak-concurrency", 100, "最大并发数")
	outputDir          = flag.String("output", "test/reports", "报告输出目录")
	maxRetries         = flag.Int("max-retries", 3, "最大重试次数")
	retryInterval      = flag.Duration("retry-interval", 5*time.Second, "重试间隔")
	monitorInterval    = flag.Duration("monitor-interval", time.Second, "监控间隔")
)

// 检查服务器是否可访问
func checkServer(url string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("无法连接到服务器: %v", err)
	}
	defer resp.Body.Close()
	return nil
}

// 等待服务器启动
func waitForServer(url string, maxRetries int, retryInterval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		if err := checkServer(url); err == nil {
			return nil
		}
		log.Printf("等待服务器启动... (尝试 %d/%d)", i+1, maxRetries)
		time.Sleep(retryInterval)
	}
	return fmt.Errorf("服务器未响应，请确保服务器已启动")
}

func main() {
	flag.Parse()

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 等待服务器启动
	log.Println("检查服务器状态...")
	if err := waitForServer(*url, *maxRetries, *retryInterval); err != nil {
		log.Fatalf("服务器检查失败: %v", err)
	}

	// 启动CPU监控
	log.Println("启动CPU监控...")
	cpuMonitor := monitor.NewCPUMonitor(*monitorInterval)
	cpuMonitor.Start()
	defer cpuMonitor.Stop()

	// 启动系统监控
	log.Println("启动系统监控...")
	systemMonitor := monitor.NewSystemMonitor(*monitorInterval)
	systemMonitor.Start()
	defer systemMonitor.Stop()

	// 运行压力测试
	log.Println("开始压力测试...")
	loadResult, err := load.RunLoadTest(*url, *duration, *initialConcurrency, *peakConcurrency)
	if err != nil {
		log.Fatalf("运行压力测试失败: %v", err)
	}

	// 保存测试结果
	resultFile := filepath.Join(*outputDir, "load_test_result.json")
	if err := load.SaveResults(loadResult, resultFile); err != nil {
		log.Fatalf("保存测试结果失败: %v", err)
	}

	// 保存CPU监控结果
	cpuStatsFile := filepath.Join(*outputDir, "cpu_stats.json")
	if err := cpuMonitor.SaveStats(cpuStatsFile); err != nil {
		log.Printf("保存CPU统计信息失败: %v", err)
	}

	// 保存系统监控结果
	systemStatsFile := filepath.Join(*outputDir, "system_stats.json")
	if err := systemMonitor.SaveStats(systemStatsFile); err != nil {
		log.Printf("保存系统统计信息失败: %v", err)
	}

	// 生成报告
	log.Println("生成测试报告...")
	report, err := reports.GenerateReport(*outputDir)
	if err != nil {
		log.Fatalf("生成报告失败: %v", err)
	}

	// 添加监控统计信息到报告
	report.CPUStats = reports.ConvertFromMonitorCPUStats(cpuMonitor.GetAverageStats())
	report.SystemStats = reports.ConvertFromMonitorSystemStats(systemMonitor.GetAverageStats())

	// 保存HTML报告
	htmlFile := filepath.Join(*outputDir, "report.html")
	if err := report.SaveHTML(htmlFile); err != nil {
		log.Fatalf("保存HTML报告失败: %v", err)
	}

	// 保存JSON报告
	jsonFile := filepath.Join(*outputDir, "report.json")
	if err := report.SaveJSON(jsonFile); err != nil {
		log.Fatalf("保存JSON报告失败: %v", err)
	}

	log.Printf("测试完成！报告已保存到: %s\n", *outputDir)
	log.Printf("HTML报告: %s\n", htmlFile)
	log.Printf("JSON报告: %s\n", jsonFile)
	log.Printf("CPU统计信息: %s\n", cpuStatsFile)
	log.Printf("系统统计信息: %s\n", systemStatsFile)
}
