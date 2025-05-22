package reports

import (
	"Taurus/test/load"
	"Taurus/test/monitor"
	"Taurus/test/performance"
	"Taurus/test/stability"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestReport 存储完整的测试报告
type TestReport struct {
	Title            string      `json:"title"`
	GeneratedAt      time.Time   `json:"generated_at"`
	TestDuration     string      `json:"test_duration"`
	Summary          Summary     `json:"summary"`
	HTTPResults      []Result    `json:"http_results"`
	GRPCResults      []Result    `json:"grpc_results"`
	LoadResults      []Result    `json:"load_results"`
	StabilityResults []Result    `json:"stability_results"`
	CPUStats         CPUStats    `json:"cpu_stats"`
	SystemStats      SystemStats `json:"system_stats"`
}

// Summary 存储测试摘要
type Summary struct {
	TotalTests      int     `json:"total_tests"`
	SuccessRate     float64 `json:"success_rate"`
	AvgResponseTime float64 `json:"avg_response_time"`
	MaxResponseTime float64 `json:"max_response_time"`
	MinResponseTime float64 `json:"min_response_time"`
	TotalRequests   int     `json:"total_requests"`
	ErrorRate       float64 `json:"error_rate"`
	AvgCPUUsage     float64 `json:"avg_cpu_usage"`
	MaxCPUUsage     float64 `json:"max_cpu_usage"`
	AvgMemoryUsage  uint64  `json:"avg_memory_usage"`
	MaxMemoryUsage  uint64  `json:"max_memory_usage"`
	AvgLoad1        float64 `json:"avg_load_1"`
	AvgLoad5        float64 `json:"avg_load_5"`
	AvgLoad15       float64 `json:"avg_load_15"`
	TotalDiskIO     uint64  `json:"total_disk_io"`
	TotalNetIO      uint64  `json:"total_net_io"`
	TotalGC         uint32  `json:"total_gc"`
	TotalGCPause    uint64  `json:"total_gc_pause"`
}

// SystemStats 存储系统统计信息
type SystemStats struct {
	LoadAvg1       float64 `json:"load_avg_1"`
	LoadAvg5       float64 `json:"load_avg_5"`
	LoadAvg15      float64 `json:"load_avg_15"`
	DiskReadBytes  uint64  `json:"disk_read_bytes"`
	DiskWriteBytes uint64  `json:"disk_write_bytes"`
	NetRecvBytes   uint64  `json:"net_recv_bytes"`
	NetSentBytes   uint64  `json:"net_sent_bytes"`
	NumGC          uint32  `json:"num_gc"`
	PauseTotalNs   uint64  `json:"pause_total_ns"`
	LastGC         uint64  `json:"last_gc"`
}

// CPUStats 存储CPU统计信息
type CPUStats struct {
	ProcessCPU    float64 `json:"process_cpu_percent"`
	SystemCPU     float64 `json:"system_cpu_percent"`
	NumCPU        int     `json:"num_cpu"`
	NumGoroutine  int     `json:"num_goroutine"`
	ProcessMemory uint64  `json:"process_memory_bytes"`
}

// ConvertFromMonitorCPUStats 将 monitor.CPUStats 转换为 reports.CPUStats
func ConvertFromMonitorCPUStats(m monitor.CPUStats) CPUStats {
	return CPUStats{
		ProcessCPU:    m.ProcessCPU,
		SystemCPU:     m.SystemCPU,
		NumCPU:        m.NumCPU,
		NumGoroutine:  m.NumGoroutine,
		ProcessMemory: m.ProcessMemory,
	}
}

// ConvertFromMonitorSystemStats 将 monitor.SystemStats 转换为 reports.SystemStats
func ConvertFromMonitorSystemStats(m monitor.SystemStats) SystemStats {
	return SystemStats{
		LoadAvg1:       m.LoadAvg1,
		LoadAvg5:       m.LoadAvg5,
		LoadAvg15:      m.LoadAvg15,
		DiskReadBytes:  m.DiskReadBytes,
		DiskWriteBytes: m.DiskWriteBytes,
		NetRecvBytes:   m.NetRecvBytes,
		NetSentBytes:   m.NetSentBytes,
		NumGC:          m.NumGC,
		PauseTotalNs:   m.PauseTotalNs,
		LastGC:         m.LastGC,
	}
}

// Result 存储单个测试结果
type Result struct {
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Duration        float64   `json:"duration_minutes"`
	TotalRequests   int       `json:"total_requests"`
	SuccessCount    int       `json:"success_count"`
	FailureCount    int       `json:"failure_count"`
	ErrorRate       float64   `json:"error_rate"`
	AvgResponseTime float64   `json:"avg_response_time"`
	MaxResponseTime float64   `json:"max_response_time"`
	MinResponseTime float64   `json:"min_response_time"`
}

// GenerateReport 生成测试报告
func GenerateReport(resultsDir string) (*TestReport, error) {
	report := &TestReport{
		Title:       "Taurus Framework 性能测试报告",
		GeneratedAt: time.Now(),
	}

	// 读取所有测试结果文件
	if err := filepath.Walk(resultsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// 只处理JSON文件
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取文件失败 %s: %v", path, err)
		}

		// 尝试解析不同类型的测试结果
		var result Result
		var loadResult load.LoadTestResult
		var stabilityResult stability.StabilityResult
		var performanceResult performance.BenchmarkResult

		// 根据文件路径判断结果类型
		switch {
		case strings.Contains(path, "load_test_result"):
			if err := json.Unmarshal(data, &loadResult); err != nil {
				return fmt.Errorf("解析负载测试结果失败 %s: %v", path, err)
			}
			result = convertLoadResult(loadResult)
		case strings.Contains(path, "stability_test_result"):
			if err := json.Unmarshal(data, &stabilityResult); err != nil {
				return fmt.Errorf("解析稳定性测试结果失败 %s: %v", path, err)
			}
			result = convertStabilityResult(stabilityResult)
		case strings.Contains(path, "performance_test_result"):
			if err := json.Unmarshal(data, &performanceResult); err != nil {
				return fmt.Errorf("解析性能测试结果失败 %s: %v", path, err)
			}
			result = convertPerformanceResult(performanceResult)
		default:
			// 尝试直接解析为通用结果
			if err := json.Unmarshal(data, &result); err != nil {
				return fmt.Errorf("解析结果失败 %s: %v", path, err)
			}
		}

		// 根据文件路径分类结果
		switch {
		case filepath.HasPrefix(path, filepath.Join(resultsDir, "http")):
			report.HTTPResults = append(report.HTTPResults, result)
		case filepath.HasPrefix(path, filepath.Join(resultsDir, "grpc")):
			report.GRPCResults = append(report.GRPCResults, result)
		case filepath.HasPrefix(path, filepath.Join(resultsDir, "load")):
			report.LoadResults = append(report.LoadResults, result)
		case filepath.HasPrefix(path, filepath.Join(resultsDir, "stability")):
			report.StabilityResults = append(report.StabilityResults, result)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// 计算摘要
	report.calculateSummary()

	return report, nil
}

// 转换负载测试结果为通用结果
func convertLoadResult(r load.LoadTestResult) Result {
	return Result{
		Name:            r.TestName,
		Type:            "load",
		StartTime:       r.StartTime,
		EndTime:         r.EndTime,
		Duration:        r.Duration,
		TotalRequests:   r.TotalRequests,
		SuccessCount:    r.SuccessCount,
		FailureCount:    r.FailureCount,
		ErrorRate:       r.ErrorRate,
		AvgResponseTime: r.AvgResponseTime,
		MaxResponseTime: r.MaxResponseTime,
		MinResponseTime: r.MinResponseTime,
	}
}

// 转换稳定性测试结果为通用结果
func convertStabilityResult(r stability.StabilityResult) Result {
	return Result{
		Name:            r.TestName,
		Type:            "stability",
		StartTime:       r.StartTime,
		EndTime:         r.EndTime,
		Duration:        r.Duration * 60, // 转换为分钟
		TotalRequests:   r.TotalRequests,
		SuccessCount:    r.SuccessCount,
		FailureCount:    r.FailureCount,
		ErrorRate:       r.ErrorRate,
		AvgResponseTime: r.AvgResponseTime,
		MaxResponseTime: r.MaxResponseTime,
		MinResponseTime: r.MinResponseTime,
	}
}

// 转换性能测试结果为通用结果
func convertPerformanceResult(r performance.BenchmarkResult) Result {
	return Result{
		Name:            r.TestName,
		Type:            "performance",
		StartTime:       r.Timestamp,
		EndTime:         r.Timestamp.Add(time.Duration(r.Duration) * time.Second),
		Duration:        r.Duration / 60, // 转换为分钟
		TotalRequests:   r.TotalRequests,
		SuccessCount:    r.TotalRequests - r.FailedRequests,
		FailureCount:    r.FailedRequests,
		ErrorRate:       float64(r.FailedRequests) / float64(r.TotalRequests) * 100,
		AvgResponseTime: r.MeanResponseTime,
		MaxResponseTime: r.MaxResponseTime,
		MinResponseTime: r.MinResponseTime,
	}
}

// calculateSummary 计算测试摘要
func (r *TestReport) calculateSummary() {
	var totalRequests, totalSuccess int
	var totalResponseTime float64
	var maxResponseTime, minResponseTime float64
	var firstTime = true

	allResults := append(append(append(r.HTTPResults, r.GRPCResults...), r.LoadResults...), r.StabilityResults...)
	r.Summary.TotalTests = len(allResults)

	for _, result := range allResults {
		totalRequests += result.TotalRequests
		totalSuccess += result.SuccessCount
		totalResponseTime += result.AvgResponseTime * float64(result.TotalRequests)

		if firstTime {
			maxResponseTime = result.MaxResponseTime
			minResponseTime = result.MinResponseTime
			firstTime = false
		} else {
			if result.MaxResponseTime > maxResponseTime {
				maxResponseTime = result.MaxResponseTime
			}
			if result.MinResponseTime < minResponseTime {
				minResponseTime = result.MinResponseTime
			}
		}
	}

	if totalRequests > 0 {
		r.Summary.SuccessRate = float64(totalSuccess) / float64(totalRequests) * 100
		r.Summary.AvgResponseTime = totalResponseTime / float64(totalRequests)
		r.Summary.ErrorRate = 100 - r.Summary.SuccessRate
	}
	r.Summary.MaxResponseTime = maxResponseTime
	r.Summary.MinResponseTime = minResponseTime
	r.Summary.TotalRequests = totalRequests

	// 添加CPU统计信息
	r.Summary.AvgCPUUsage = r.CPUStats.ProcessCPU
	r.Summary.MaxCPUUsage = r.CPUStats.SystemCPU
	r.Summary.AvgMemoryUsage = r.CPUStats.ProcessMemory
	r.Summary.MaxMemoryUsage = r.CPUStats.ProcessMemory

	// 添加系统统计信息
	r.Summary.AvgLoad1 = r.SystemStats.LoadAvg1
	r.Summary.AvgLoad5 = r.SystemStats.LoadAvg5
	r.Summary.AvgLoad15 = r.SystemStats.LoadAvg15
	r.Summary.TotalDiskIO = r.SystemStats.DiskReadBytes + r.SystemStats.DiskWriteBytes
	r.Summary.TotalNetIO = r.SystemStats.NetRecvBytes + r.SystemStats.NetSentBytes
	r.Summary.TotalGC = r.SystemStats.NumGC
	r.Summary.TotalGCPause = r.SystemStats.PauseTotalNs
}

// SaveHTML 保存HTML格式的报告
func (r *TestReport) SaveHTML(filename string) error {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .summary { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .result { margin: 20px 0; padding: 10px; border: 1px solid #ddd; }
        .chart { width: 100%; height: 300px; margin: 20px 0; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 8px; text-align: left; border: 1px solid #ddd; }
        th { background-color: #f5f5f5; }
    </style>
</head>
<body>
    <h1>{{.Title}}</h1>
    <p>生成时间: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</p>
    
    <div class="summary">
        <h2>测试摘要</h2>
        <p>总测试数: {{.Summary.TotalTests}}</p>
        <p>成功率: {{printf "%.2f" .Summary.SuccessRate}}%</p>
        <p>平均响应时间: {{printf "%.2f" .Summary.AvgResponseTime}}ms</p>
        <p>最大响应时间: {{printf "%.2f" .Summary.MaxResponseTime}}ms</p>
        <p>最小响应时间: {{printf "%.2f" .Summary.MinResponseTime}}ms</p>
        <p>总请求数: {{.Summary.TotalRequests}}</p>
        <p>错误率: {{printf "%.2f" .Summary.ErrorRate}}%</p>
        
        <h3>系统资源使用</h3>
        <p>CPU核心数: {{.CPUStats.NumCPU}}</p>
        <p>进程CPU使用率: {{printf "%.2f" .CPUStats.ProcessCPU}}%</p>
        <p>系统CPU使用率: {{printf "%.2f" .CPUStats.SystemCPU}}%</p>
        <p>Goroutine数量: {{.CPUStats.NumGoroutine}}</p>
        <p>进程内存使用: {{printf "%.2f" (div .CPUStats.ProcessMemory 1024.0 1024.0)}} MB</p>
        
        <h3>系统负载</h3>
        <p>1分钟负载: {{printf "%.2f" .SystemStats.LoadAvg1}}</p>
        <p>5分钟负载: {{printf "%.2f" .SystemStats.LoadAvg5}}</p>
        <p>15分钟负载: {{printf "%.2f" .SystemStats.LoadAvg15}}</p>
        
        <h3>I/O统计</h3>
        <p>磁盘读取: {{printf "%.2f" (div .SystemStats.DiskReadBytes 1024.0 1024.0)}} MB</p>
        <p>磁盘写入: {{printf "%.2f" (div .SystemStats.DiskWriteBytes 1024.0 1024.0)}} MB</p>
        <p>网络接收: {{printf "%.2f" (div .SystemStats.NetRecvBytes 1024.0 1024.0)}} MB</p>
        <p>网络发送: {{printf "%.2f" (div .SystemStats.NetSentBytes 1024.0 1024.0)}} MB</p>
        
        <h3>GC统计</h3>
        <p>GC次数: {{.SystemStats.NumGC}}</p>
        <p>GC暂停总时间: {{printf "%.2f" (div .SystemStats.PauseTotalNs 1000000.0)}} ms</p>
        <p>上次GC时间: {{printf "%.2f" (div .SystemStats.LastGC 1000000.0)}} ms</p>
    </div>

    <h2>HTTP测试结果</h2>
    {{range .HTTPResults}}
    <div class="result">
        <h3>{{.Name}}</h3>
        <table>
            <tr><th>指标</th><th>值</th></tr>
            <tr><td>总请求数</td><td>{{.TotalRequests}}</td></tr>
            <tr><td>成功数</td><td>{{.SuccessCount}}</td></tr>
            <tr><td>失败数</td><td>{{.FailureCount}}</td></tr>
            <tr><td>错误率</td><td>{{printf "%.2f" .ErrorRate}}%</td></tr>
            <tr><td>平均响应时间</td><td>{{printf "%.2f" .AvgResponseTime}}ms</td></tr>
        </table>
    </div>
    {{end}}

    <h2>gRPC测试结果</h2>
    {{range .GRPCResults}}
    <div class="result">
        <h3>{{.Name}}</h3>
        <table>
            <tr><th>指标</th><th>值</th></tr>
            <tr><td>总请求数</td><td>{{.TotalRequests}}</td></tr>
            <tr><td>成功数</td><td>{{.SuccessCount}}</td></tr>
            <tr><td>失败数</td><td>{{.FailureCount}}</td></tr>
            <tr><td>错误率</td><td>{{printf "%.2f" .ErrorRate}}%</td></tr>
            <tr><td>平均响应时间</td><td>{{printf "%.2f" .AvgResponseTime}}ms</td></tr>
        </table>
    </div>
    {{end}}

    <h2>压力测试结果</h2>
    {{range .LoadResults}}
    <div class="result">
        <h3>{{.Name}}</h3>
        <table>
            <tr><th>指标</th><th>值</th></tr>
            <tr><td>总请求数</td><td>{{.TotalRequests}}</td></tr>
            <tr><td>成功数</td><td>{{.SuccessCount}}</td></tr>
            <tr><td>失败数</td><td>{{.FailureCount}}</td></tr>
            <tr><td>错误率</td><td>{{printf "%.2f" .ErrorRate}}%</td></tr>
            <tr><td>平均响应时间</td><td>{{printf "%.2f" .AvgResponseTime}}ms</td></tr>
        </table>
    </div>
    {{end}}

    <h2>稳定性测试结果</h2>
    {{range .StabilityResults}}
    <div class="result">
        <h3>{{.Name}}</h3>
        <table>
            <tr><th>指标</th><th>值</th></tr>
            <tr><td>总请求数</td><td>{{.TotalRequests}}</td></tr>
            <tr><td>成功数</td><td>{{.SuccessCount}}</td></tr>
            <tr><td>失败数</td><td>{{.FailureCount}}</td></tr>
            <tr><td>错误率</td><td>{{printf "%.2f" .ErrorRate}}%</td></tr>
            <tr><td>平均响应时间</td><td>{{printf "%.2f" .AvgResponseTime}}ms</td></tr>
        </table>
    </div>
    {{end}}
</body>
</html>`

	funcMap := template.FuncMap{
		"div": func(a, b, c float64) float64 {
			return a / b / c
		},
	}

	t, err := template.New("report").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer f.Close()

	return t.Execute(f, r)
}

// SaveJSON 保存JSON格式的报告
func (r *TestReport) SaveJSON(filename string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化报告失败: %v", err)
	}

	return os.WriteFile(filename, data, 0644)
}
