// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

// Package telemetry 提供了基于 OpenTelemetry 的分布式追踪功能。
// 该包支持配置化的链路追踪，可以灵活地设置服务信息、采样率、批处理等参数。
package telemetry

import (
	"time"
)

// ExportProtocol 定义导出协议类型
type ExportProtocol string

const (
	// ProtocolGRPC 使用 gRPC 协议导出
	ProtocolGRPC ExportProtocol = "grpc"
	// ProtocolHTTP 使用 HTTP 协议导出
	ProtocolHTTP ExportProtocol = "http"
	// ProtocolJSON 使用 HTTP/JSON 协议导出
	ProtocolJSON ExportProtocol = "json"
)

// Option 定义配置选项函数类型，用于设置 Provider 的各项参数。
// 使用函数式选项模式，可以灵活地配置所需的参数，无需传入完整的配置结构体。
type Option func(*options)

// options 内部配置结构体，存储所有可配置的参数。
type options struct {
	// 基础配置
	serviceName    string // 服务名称，用于标识追踪数据来源
	serviceVersion string // 服务版本号
	environment    string // 运行环境，如 development, staging, production

	// OTLP 导出器配置
	protocol ExportProtocol // 导出协议类型
	endpoint string         // OTLP 接收器地址，如 localhost:4317
	insecure bool           // 是否使用非安全连接
	timeout  time.Duration  // 导出超时时间

	// 采样配置
	samplingRatio float64 // 采样率，范围 0.0-1.0

	// 批处理配置
	batchTimeout       time.Duration // 批处理超时时间
	exportTimeout      time.Duration // 导出超时时间
	maxExportBatchSize int           // 最大导出批次大小
	maxQueueSize       int           // 最大队列大小
}

// defaultOptions 返回默认配置，提供合理的默认值
func defaultOptions() *options {
	return &options{
		serviceName:        "unknown-service",
		serviceVersion:     "1.0.0",
		environment:        "development",
		protocol:           ProtocolGRPC, // 默认使用 gRPC
		endpoint:           "localhost:4317",
		insecure:           true,
		timeout:            5 * time.Second,
		samplingRatio:      1.0,
		batchTimeout:       5 * time.Second,
		exportTimeout:      30 * time.Second,
		maxExportBatchSize: 512,
		maxQueueSize:       2048,
	}
}

// WithServiceName 设置服务名称
// name: 服务的唯一标识名称
func WithServiceName(name string) Option {
	return func(o *options) {
		o.serviceName = name
	}
}

// WithServiceVersion 设置服务版本
// version: 服务的版本号，如 1.0.0
func WithServiceVersion(version string) Option {
	return func(o *options) {
		o.serviceVersion = version
	}
}

// WithEnvironment 设置环境
// env: 运行环境，如 development, staging, production
func WithEnvironment(env string) Option {
	return func(o *options) {
		o.environment = env
	}
}

// WithEndpoint 设置OTLP端点
// endpoint: OTLP 接收器地址，如 localhost:4317
func WithEndpoint(endpoint string) Option {
	return func(o *options) {
		o.endpoint = endpoint
	}
}

// WithInsecure 设置是否使用非安全连接
// insecure: true表示使用非安全连接
func WithInsecure(insecure bool) Option {
	return func(o *options) {
		o.insecure = insecure
	}
}

// WithTimeout 设置超时时间
// timeout: 导出器操作的超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}

// WithSamplingRatio 设置采样率
// ratio: 采样率，范围 0.0-1.0，1.0 表示全采样
func WithSamplingRatio(ratio float64) Option {
	return func(o *options) {
		o.samplingRatio = ratio
	}
}

// WithBatchTimeout 设置批处理超时时间
// timeout: 批处理的超时时间
func WithBatchTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.batchTimeout = timeout
	}
}

// WithExportTimeout 设置导出超时时间
// timeout: 导出操作的超时时间
func WithExportTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.exportTimeout = timeout
	}
}

// WithMaxExportBatchSize 设置最大导出批次大小
// size: 单次导出的最大数据量
func WithMaxExportBatchSize(size int) Option {
	return func(o *options) {
		o.maxExportBatchSize = size
	}
}

// WithMaxQueueSize 设置最大队列大小
// size: 等待导出数据的最大队列长度
func WithMaxQueueSize(size int) Option {
	return func(o *options) {
		o.maxQueueSize = size
	}
}

// WithExportProtocol 设置导出协议类型
// protocol: 导出协议类型，支持 grpc、http、json
func WithExportProtocol(protocol ExportProtocol) Option {
	return func(o *options) {
		o.protocol = protocol
	}
}
