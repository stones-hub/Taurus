package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"Taurus/test/load"
	"Taurus/test/reports"
)

var (
	url                = flag.String("url", "http://127.0.0.1:9080/v1/api/", "测试目标URL")
	duration           = flag.Duration("duration", 5*time.Minute, "测试持续时间")
	initialConcurrency = flag.Int("initial-concurrency", 10, "初始并发数")
	peakConcurrency    = flag.Int("peak-concurrency", 100, "最大并发数")
	outputDir          = flag.String("output", "test/reports", "报告输出目录")
)

func main() {
	flag.Parse()

	// 创建输出目录
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 运行压力测试
	fmt.Println("开始压力测试...")
	loadResult, err := load.RunLoadTest(*url, *duration, *initialConcurrency, *peakConcurrency)
	if err != nil {
		log.Fatalf("运行压力测试失败: %v", err)
	}

	// 保存测试结果
	resultFile := filepath.Join(*outputDir, "load_test_result.json")
	if err := load.SaveResults(loadResult, resultFile); err != nil {
		log.Fatalf("保存测试结果失败: %v", err)
	}

	// 生成报告
	fmt.Println("生成测试报告...")
	report, err := reports.GenerateReport(*outputDir)
	if err != nil {
		log.Fatalf("生成报告失败: %v", err)
	}

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

	fmt.Printf("测试完成！报告已保存到: %s\n", *outputDir)
	fmt.Printf("HTML报告: %s\n", htmlFile)
	fmt.Printf("JSON报告: %s\n", jsonFile)
}
