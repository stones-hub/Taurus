package reports

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

// TestReport 存储完整的测试报告
type TestReport struct {
	Title            string    `json:"title"`
	GeneratedAt      time.Time `json:"generated_at"`
	TestDuration     string    `json:"test_duration"`
	Summary          Summary   `json:"summary"`
	HTTPResults      []Result  `json:"http_results"`
	GRPCResults      []Result  `json:"grpc_results"`
	LoadResults      []Result  `json:"load_results"`
	StabilityResults []Result  `json:"stability_results"`
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

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取文件失败 %s: %v", path, err)
		}

		var result Result
		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf("解析结果失败 %s: %v", path, err)
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

	t, err := template.New("report").Parse(tmpl)
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
