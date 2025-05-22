package stability

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// StabilityResult 存储稳定性测试结果
type StabilityResult struct {
	TestName        string    `json:"test_name"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Duration        float64   `json:"duration_hours"`
	TotalRequests   int       `json:"total_requests"`
	SuccessCount    int       `json:"success_count"`
	FailureCount    int       `json:"failure_count"`
	ErrorRate       float64   `json:"error_rate"`
	AvgResponseTime float64   `json:"avg_response_time"`
	MaxResponseTime float64   `json:"max_response_time"`
	MinResponseTime float64   `json:"min_response_time"`
	Errors          []string  `json:"errors"`
}

// RunStabilityTest 运行稳定性测试
func RunStabilityTest(url string, duration time.Duration, concurrency int) (*StabilityResult, error) {
	result := &StabilityResult{
		TestName:  "Stability Test",
		StartTime: time.Now(),
		Errors:    make([]string, 0),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建通道来控制测试结束
	done := make(chan bool)
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	// 启动工作协程
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					start := time.Now()
					resp, err := client.Get(url)
					duration := time.Since(start).Seconds() * 1000 // 转换为毫秒

					mu.Lock()
					result.TotalRequests++
					if err != nil {
						result.FailureCount++
						result.Errors = append(result.Errors, fmt.Sprintf("请求失败: %v", err))
					} else {
						result.SuccessCount++
						if resp.StatusCode != http.StatusOK {
							result.FailureCount++
							result.Errors = append(result.Errors, fmt.Sprintf("HTTP状态码错误: %d", resp.StatusCode))
						}
						resp.Body.Close()
					}

					// 更新响应时间统计
					if result.MinResponseTime == 0 || duration < result.MinResponseTime {
						result.MinResponseTime = duration
					}
					if duration > result.MaxResponseTime {
						result.MaxResponseTime = duration
					}
					result.AvgResponseTime = (result.AvgResponseTime*float64(result.TotalRequests-1) + duration) / float64(result.TotalRequests)
					mu.Unlock()

					// 控制请求频率
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
	}

	// 等待测试时间结束
	<-ticker.C
	close(done)
	wg.Wait()

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime).Hours()
	if result.TotalRequests > 0 {
		result.ErrorRate = float64(result.FailureCount) / float64(result.TotalRequests) * 100
	}

	return result, nil
}

// SaveResults 保存测试结果
func SaveResults(result *StabilityResult, filename string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化测试结果失败: %v", err)
	}

	return os.WriteFile(filename, data, 0644)
}
