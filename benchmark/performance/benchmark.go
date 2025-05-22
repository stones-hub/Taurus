package performance

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// BenchmarkResult 存储性能测试结果
type BenchmarkResult struct {
	TestName         string    `json:"test_name"`
	Timestamp        time.Time `json:"timestamp"`
	TotalRequests    int       `json:"total_requests"`
	Concurrency      int       `json:"concurrency"`
	Duration         float64   `json:"duration"`
	RequestsPerSec   float64   `json:"requests_per_sec"`
	TimePerRequest   float64   `json:"time_per_request"`
	FailedRequests   int       `json:"failed_requests"`
	MinResponseTime  float64   `json:"min_response_time"`
	MaxResponseTime  float64   `json:"max_response_time"`
	MeanResponseTime float64   `json:"mean_response_time"`
}

// RunHTTPBenchmark 运行 HTTP 性能测试
func RunHTTPBenchmark(url string, totalRequests, concurrency int) (*BenchmarkResult, error) {
	// 构建 ab 命令
	cmd := exec.Command("ab",
		"-n", strconv.Itoa(totalRequests),
		"-c", strconv.Itoa(concurrency),
		"-g", "benchmark_results.dat",
		url)

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("执行性能测试失败: %v\n输出: %s", err, output)
	}

	// 解析结果
	result := &BenchmarkResult{
		TestName:      "HTTP Benchmark",
		Timestamp:     time.Now(),
		TotalRequests: totalRequests,
		Concurrency:   concurrency,
	}

	// 解析 ab 输出
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		switch {
		case strings.Contains(line, "Time per request"):
			fmt.Sscanf(line, "Time per request: %f [ms]", &result.TimePerRequest)
		case strings.Contains(line, "Requests per second"):
			fmt.Sscanf(line, "Requests per second: %f", &result.RequestsPerSec)
		case strings.Contains(line, "Failed requests"):
			fmt.Sscanf(line, "Failed requests: %d", &result.FailedRequests)
		case strings.Contains(line, "Total transferred"):
			// 可以添加更多指标解析
		}
	}

	return result, nil
}

// RunGRPCBenchmark 运行 gRPC 性能测试
func RunGRPCBenchmark(protoFile, service, method string, totalRequests, concurrency int) (*BenchmarkResult, error) {
	// 构建 ghz 命令
	cmd := exec.Command("ghz",
		"--insecure",
		"--proto", protoFile,
		"--call", fmt.Sprintf("%s.%s", service, method),
		"-n", strconv.Itoa(totalRequests),
		"-c", strconv.Itoa(concurrency),
		"-d", `{"name":"test"}`,
		"127.0.0.1:50051")

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("执行 gRPC 性能测试失败: %v\n输出: %s", err, output)
	}

	// 解析结果
	result := &BenchmarkResult{
		TestName:      "gRPC Benchmark",
		Timestamp:     time.Now(),
		TotalRequests: totalRequests,
		Concurrency:   concurrency,
	}

	// 解析 ghz 输出
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		switch {
		case strings.Contains(line, "Requests/sec"):
			fmt.Sscanf(line, "Requests/sec: %f", &result.RequestsPerSec)
		case strings.Contains(line, "Average"):
			fmt.Sscanf(line, "Average: %f ms", &result.MeanResponseTime)
		case strings.Contains(line, "Fastest"):
			fmt.Sscanf(line, "Fastest: %f ms", &result.MinResponseTime)
		case strings.Contains(line, "Slowest"):
			fmt.Sscanf(line, "Slowest: %f ms", &result.MaxResponseTime)
		}
	}

	return result, nil
}

// SaveResults 保存测试结果
func SaveResults(result *BenchmarkResult, filename string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化测试结果失败: %v", err)
	}

	return os.WriteFile(filename, data, 0644)
}
