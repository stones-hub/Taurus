package config

import (
	"time"
)

// PerformanceConfig 性能测试配置
type PerformanceConfig struct {
	Protocol          string      `json:"protocol"`
	Concurrency       int         `json:"concurrency"`
	MaxConcurrency    int         `json:"max_concurrency"`
	RequestsPerSecond int         `json:"requests_per_second"`
	Duration          string      `json:"duration"`
	Timeout           string      `json:"timeout"`
	HTTP              *HTTPConfig `json:"http"`
	GRPC              *GRPCConfig `json:"grpc"`
}

// HTTPConfig HTTP 测试配置
type HTTPConfig struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// GRPCConfig gRPC 测试配置
type GRPCConfig struct {
	ServerAddr  string            `json:"server_addr"`
	Service     string            `json:"service"`
	Method      string            `json:"method"`
	RequestData string            `json:"request_data"`
	Metadata    map[string]string `json:"metadata"`
}

// GetDuration 获取测试持续时间
func (c *PerformanceConfig) GetDuration() (time.Duration, error) {
	return time.ParseDuration(c.Duration)
}

// GetTimeout 获取超时时间
func (c *PerformanceConfig) GetTimeout() (time.Duration, error) {
	return time.ParseDuration(c.Timeout)
}
