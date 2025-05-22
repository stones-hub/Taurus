package vegeta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"Taurus/benchmark/config"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// VegetaTester Vegeta 测试器
type VegetaTester struct {
	config *config.PerformanceConfig
}

// NewVegetaTester 创建新的 Vegeta 测试器
func NewVegetaTester(config *config.PerformanceConfig) *VegetaTester {
	return &VegetaTester{
		config: config,
	}
}

// TestResult 测试结果
type TestResult struct {
	Latencies struct {
		Mean time.Duration `json:"mean"`
		P50  time.Duration `json:"p50"`
		P90  time.Duration `json:"p90"`
		P95  time.Duration `json:"p95"`
		P99  time.Duration `json:"p99"`
		Max  time.Duration `json:"max"`
	} `json:"latencies"`
	BytesIn struct {
		Total uint64  `json:"total"`
		Mean  float64 `json:"mean"`
	} `json:"bytes_in"`
	BytesOut struct {
		Total uint64  `json:"total"`
		Mean  float64 `json:"mean"`
	} `json:"bytes_out"`
	Earliest    time.Time      `json:"earliest"`
	Latest      time.Time      `json:"latest"`
	End         time.Time      `json:"end"`
	Duration    time.Duration  `json:"duration"`
	Wait        time.Duration  `json:"wait"`
	Requests    uint64         `json:"requests"`
	Rate        float64        `json:"rate"`
	Throughput  float64        `json:"throughput"`
	Success     float64        `json:"success"`
	StatusCodes map[string]int `json:"status_codes"`
	Errors      []string       `json:"errors"`
}

// RunTest 运行测试
func (t *VegetaTester) RunTest() (*TestResult, error) {
	// 创建目标
	targeter := t.createTargeter()
	if targeter == nil {
		return nil, fmt.Errorf("创建目标失败")
	}

	// 获取超时时间
	timeout, err := t.config.GetTimeout()
	if err != nil {
		return nil, fmt.Errorf("获取超时时间失败: %v", err)
	}

	// 创建攻击器
	attacker := vegeta.NewAttacker(
		vegeta.Timeout(timeout),
		vegeta.Workers(uint64(t.config.Concurrency)),
		vegeta.MaxWorkers(uint64(t.config.MaxConcurrency)),
	)

	// 设置测试参数
	duration, err := t.config.GetDuration()
	if err != nil {
		return nil, fmt.Errorf("获取测试持续时间失败: %v", err)
	}

	rate := vegeta.Rate{
		Freq: t.config.RequestsPerSecond,
		Per:  time.Second,
	}

	// 运行测试
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Benchmark Test") {
		metrics.Add(res)
	}
	metrics.Close()

	// 生成报告
	result := t.generateReport(&metrics)

	// 保存报告
	if err := t.saveReport(result); err != nil {
		return nil, fmt.Errorf("保存报告失败: %v", err)
	}

	return result, nil
}

// createTargeter 创建目标
func (t *VegetaTester) createTargeter() vegeta.Targeter {
	switch t.config.Protocol {
	case "http":
		return t.createHTTPTargeter()
	case "grpc":
		return t.createGRPCTargeter()
	default:
		return nil
	}
}

// createHTTPTargeter 创建 HTTP 目标
func (t *VegetaTester) createHTTPTargeter() vegeta.Targeter {
	// 转换 headers
	headers := make(http.Header)
	for key, value := range t.config.HTTP.Headers {
		headers.Set(key, value)
	}

	return vegeta.NewStaticTargeter(vegeta.Target{
		Method: t.config.HTTP.Method,
		URL:    t.config.HTTP.URL,
		Header: headers,
		Body:   []byte(t.config.HTTP.Body),
	})
}

// createGRPCTargeter 创建 gRPC 目标
func (t *VegetaTester) createGRPCTargeter() vegeta.Targeter {
	// TODO: 实现 gRPC 目标创建
	return nil
}

// generateReport 生成报告
func (t *VegetaTester) generateReport(metrics *vegeta.Metrics) *TestResult {
	return &TestResult{
		Latencies: struct {
			Mean time.Duration `json:"mean"`
			P50  time.Duration `json:"p50"`
			P90  time.Duration `json:"p90"`
			P95  time.Duration `json:"p95"`
			P99  time.Duration `json:"p99"`
			Max  time.Duration `json:"max"`
		}{
			Mean: metrics.Latencies.Mean,
			P50:  metrics.Latencies.P50,
			P90:  metrics.Latencies.P90,
			P95:  metrics.Latencies.P95,
			P99:  metrics.Latencies.P99,
			Max:  metrics.Latencies.Max,
		},
		BytesIn: struct {
			Total uint64  `json:"total"`
			Mean  float64 `json:"mean"`
		}{
			Total: metrics.BytesIn.Total,
			Mean:  metrics.BytesIn.Mean,
		},
		BytesOut: struct {
			Total uint64  `json:"total"`
			Mean  float64 `json:"mean"`
		}{
			Total: metrics.BytesOut.Total,
			Mean:  metrics.BytesOut.Mean,
		},
		Earliest:    metrics.Earliest,
		Latest:      metrics.Latest,
		End:         metrics.End,
		Duration:    metrics.Duration,
		Wait:        metrics.Wait,
		Requests:    metrics.Requests,
		Rate:        metrics.Rate,
		Throughput:  metrics.Throughput,
		Success:     metrics.Success,
		StatusCodes: metrics.StatusCodes,
		Errors:      metrics.Errors,
	}
}

// saveReport 保存报告
func (t *VegetaTester) saveReport(result *TestResult) error {
	// 创建报告目录
	reportDir := filepath.Join("benchmark", "reports", "vegeta")
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return fmt.Errorf("创建报告目录失败: %v", err)
	}

	// 生成报告文件名
	timestamp := time.Now().Format("20060102_150405")
	reportFile := filepath.Join(reportDir, fmt.Sprintf("report_%s.json", timestamp))

	// 保存 JSON 报告
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %v", err)
	}

	if err := os.WriteFile(reportFile, jsonData, 0644); err != nil {
		return fmt.Errorf("写入报告文件失败: %v", err)
	}

	// 生成 HTML 报告
	htmlReport := t.generateHTMLReport(result)
	htmlFile := filepath.Join(reportDir, fmt.Sprintf("report_%s.html", timestamp))
	if err := os.WriteFile(htmlFile, []byte(htmlReport), 0644); err != nil {
		return fmt.Errorf("写入 HTML 报告失败: %v", err)
	}

	return nil
}

// generateHTMLReport 生成 HTML 报告
func (t *VegetaTester) generateHTMLReport(result *TestResult) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Vegeta 测试报告</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .section { margin-bottom: 20px; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .metric { margin: 10px 0; }
        .metric-label { font-weight: bold; }
        .chart { width: 100%%; height: 400px; margin: 20px 0; }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>
    <div class="container">
        <h1>Vegeta 测试报告</h1>
        
        <div class="section">
            <h2>测试配置</h2>
            <div class="metric">
                <span class="metric-label">协议:</span> %s
            </div>
            <div class="metric">
                <span class="metric-label">并发数:</span> %d
            </div>
            <div class="metric">
                <span class="metric-label">持续时间:</span> %s
            </div>
        </div>

        <div class="section">
            <h2>性能指标</h2>
            <div class="metric">
                <span class="metric-label">总请求数:</span> %d
            </div>
            <div class="metric">
                <span class="metric-label">请求速率:</span> %.2f req/s
            </div>
            <div class="metric">
                <span class="metric-label">吞吐量:</span> %.2f MB/s
            </div>
            <div class="metric">
                <span class="metric-label">成功率:</span> %.2f%%
            </div>
        </div>

        <div class="section">
            <h2>响应时间</h2>
            <div class="metric">
                <span class="metric-label">平均响应时间:</span> %s
            </div>
            <div class="metric">
                <span class="metric-label">P50:</span> %s
            </div>
            <div class="metric">
                <span class="metric-label">P90:</span> %s
            </div>
            <div class="metric">
                <span class="metric-label">P95:</span> %s
            </div>
            <div class="metric">
                <span class="metric-label">P99:</span> %s
            </div>
            <div class="metric">
                <span class="metric-label">最大响应时间:</span> %s
            </div>
        </div>

        <div class="section">
            <h2>状态码分布</h2>
            <div id="statusCodesChart" class="chart"></div>
        </div>

        <div class="section">
            <h2>错误信息</h2>
            <ul>
                %s
            </ul>
        </div>
    </div>

    <script>
        // 状态码分布图表
        const statusCodesData = %s;
        new Chart(document.getElementById('statusCodesChart'), {
            type: 'pie',
            data: {
                labels: Object.keys(statusCodesData),
                datasets: [{
                    data: Object.values(statusCodesData),
                    backgroundColor: [
                        '#4CAF50',
                        '#2196F3',
                        '#FFC107',
                        '#F44336'
                    ]
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'top',
                    },
                    title: {
                        display: true,
                        text: 'HTTP 状态码分布'
                    }
                }
            }
        });
    </script>
</body>
</html>
`,
		t.config.Protocol,
		t.config.Concurrency,
		t.config.Duration,
		result.Requests,
		result.Rate,
		result.Throughput,
		result.Success*100,
		result.Latencies.Mean,
		result.Latencies.P50,
		result.Latencies.P90,
		result.Latencies.P95,
		result.Latencies.P99,
		result.Latencies.Max,
		formatErrors(result.Errors),
		formatStatusCodes(result.StatusCodes),
	)
}

// formatErrors 格式化错误信息
func formatErrors(errors []string) string {
	if len(errors) == 0 {
		return "<li>无错误</li>"
	}

	var result string
	for _, err := range errors {
		result += fmt.Sprintf("<li>%s</li>", err)
	}
	return result
}

// formatStatusCodes 格式化状态码数据
func formatStatusCodes(codes map[string]int) string {
	data, _ := json.Marshal(codes)
	return string(data)
}
