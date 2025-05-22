package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"Taurus/benchmark/config"
	"Taurus/benchmark/vegeta"
)

func main() {
	// 解析命令行参数
	configFile := flag.String("config", "benchmark/config/vegeta.json", "配置文件路径")
	flag.Parse()

	// 读取配置文件
	configData, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析配置
	var perfConfig config.PerformanceConfig
	if err := json.Unmarshal(configData, &perfConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 创建测试器
	tester := vegeta.NewVegetaTester(&perfConfig)

	// 运行测试
	result, err := tester.RunTest()
	if err != nil {
		log.Fatalf("运行测试失败: %v", err)
	}

	// 打印结果
	fmt.Printf("测试完成！\n")
	fmt.Printf("总请求数: %d\n", result.Requests)
	fmt.Printf("请求速率: %.2f req/s\n", result.Rate)
	fmt.Printf("吞吐量: %.2f MB/s\n", result.Throughput)
	fmt.Printf("成功率: %.2f%%\n", result.Success*100)
	fmt.Printf("平均响应时间: %s\n", result.Latencies.Mean)
	fmt.Printf("P95 响应时间: %s\n", result.Latencies.P95)
	fmt.Printf("P99 响应时间: %s\n", result.Latencies.P99)
	fmt.Printf("最大响应时间: %s\n", result.Latencies.Max)
	fmt.Printf("详细报告已保存到 benchmark/reports/vegeta 目录\n")
}
