// Package telemetry 提供分布式追踪功能
package telemetry

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Protocol 定义了支持的协议类型
type Protocol = ExportProtocol

// OTelProvider OpenTelemetry 追踪提供者实现
type OTelProvider struct {
	opts           *options
	tracerProvider *sdktrace.TracerProvider // 追踪提供者实例
	once           sync.Once
}

// NewOTelProvider 创建新的 OpenTelemetry 追踪提供者
// 该函数会初始化所有必要的组件，包括导出器、资源属性和采样器
func NewOTelProvider(opts ...Option) (*OTelProvider, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	p := &OTelProvider{
		opts: options,
	}

	var err error
	p.once.Do(func() {
		fmt.Println("-------------------------------- create tracer provider --------------------------------")
		// 1. 创建 OTLP 导出器（决定数据导出/发送到哪里）
		var exporter sdktrace.SpanExporter
		exporter, err = p.createExporter()
		if err != nil {
			err = fmt.Errorf("create exporter failed: %w", err)
			return
		}

		// 2.创建资源 (采集的数据的标记，让后续的分析可以识别)
		var res *resource.Resource
		res, err = p.createResource()
		if err != nil {
			err = fmt.Errorf("create resource failed: %w", err)
			return
		}

		// 3. 创建追踪器的提供者, 将导出器和资源属性传递给追踪器提供者
		p.tracerProvider = p.createTracerProvider(exporter, res)

		// 将我们创建的追踪提供者设置为全局默认提供者
		// 这样其他包就可以直接使用 otel.GetTracerProvider() 获取到这个提供者
		// 就像一个总控制台，告诉系统："这就是我们要用的追踪系统"
		otel.SetTracerProvider(p.tracerProvider)

		// 设置如何在服务之间传递追踪信息
		// 约定一个通用的交流方式，让不同的服务能互相理解追踪信息
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C 标准的追踪上下文传播
			propagation.Baggage{},      // 用于传递自定义属性
		))
	})

	if err != nil {
		return nil, fmt.Errorf("setup provider failed: %w", err)
	}
	return p, nil
}

// createExporter 创建 OTLP 导出器
func (p *OTelProvider) createExporter() (sdktrace.SpanExporter, error) {
	endpoint := p.opts.endpoint
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	switch p.opts.protocol {
	case ProtocolGRPC:
		if endpoint == "" {
			endpoint = "localhost:4317"
		}
		var exporterOpts []otlptracegrpc.Option
		if p.opts.insecure {
			exporterOpts = append(exporterOpts,
				otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithEndpoint(endpoint))
		} else {
			exporterOpts = append(exporterOpts,
				otlptracegrpc.WithEndpoint(endpoint))
		}

		if p.opts.timeout > 0 {
			exporterOpts = append(exporterOpts, otlptracegrpc.WithTimeout(p.opts.timeout))
		}

		return otlptracegrpc.New(context.Background(), exporterOpts...)

	case ProtocolHTTP, ProtocolJSON:
		if endpoint == "" {
			endpoint = "localhost:4318"
		}
		var exporterOpts []otlptracehttp.Option
		if p.opts.insecure {
			exporterOpts = append(exporterOpts,
				otlptracehttp.WithInsecure(),
				otlptracehttp.WithEndpoint(endpoint))
		} else {
			exporterOpts = append(exporterOpts,
				otlptracehttp.WithEndpoint(endpoint))
		}

		if p.opts.timeout > 0 {
			exporterOpts = append(exporterOpts, otlptracehttp.WithTimeout(p.opts.timeout))
		}

		if p.opts.protocol == ProtocolJSON {
			exporterOpts = append(exporterOpts, otlptracehttp.WithHeaders(map[string]string{
				"Content-Type": "application/json",
			}))
		}

		return otlptracehttp.New(context.Background(), exporterOpts...)

	default:
		return nil, fmt.Errorf("unsupported protocol: %v", p.opts.protocol)
	}
}

// createResource 创建资源属性
func (p *OTelProvider) createResource() (*resource.Resource, error) {
	return resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(p.opts.serviceName),
			semconv.ServiceVersion(p.opts.serviceVersion),
			attribute.String("environment", p.opts.environment),
		),
	)
}

// createTracerProvider 创建追踪提供者实例
// 配置如何收集、处理和导出追踪数据
func (p *OTelProvider) createTracerProvider(exp sdktrace.SpanExporter, res *resource.Resource) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		// 1. 配置批处理导出器
		sdktrace.WithBatcher(exp,
			// 一次最多导出多少条数据
			// 比如设置为 512，就是凑够 512 条数据就会导出一次
			sdktrace.WithMaxExportBatchSize(p.opts.maxExportBatchSize),

			// 多久导出一次数据
			// 比如设置为 5s，即使数据没凑够 512 条，5s 后也会导出
			sdktrace.WithBatchTimeout(p.opts.batchTimeout),

			// 最多能缓存多少条待导出的数据
			// 比如设置为 2048，超过后新的数据就会被丢弃
			sdktrace.WithMaxQueueSize(p.opts.maxQueueSize),

			// 导出超时时间
			// 比如设置为 10s，如果导出超时，数据就会被丢弃
			sdktrace.WithExportTimeout(p.opts.exportTimeout),
		),

		// 2. 设置资源属性
		// 就像前面说的快递单，每条数据都会带上这些标签
		sdktrace.WithResource(res),

		// 3. 配置采样策略
		sdktrace.WithSampler(
			// 使用基于父 Span 的采样策略
			sdktrace.ParentBased(
				// 设置采样比例
				// samplingRatio 范围是 0-1
				// 0 表示不采样，1 表示全采样
				// 0.1 表示采样 10% 的数据
				sdktrace.TraceIDRatioBased(p.opts.samplingRatio),
			),
		),
	)
}

// Tracer 返回指定名称的追踪器
func (p *OTelProvider) Tracer(name string) trace.Tracer {
	return p.tracerProvider.Tracer(name)
}

// Shutdown 关闭追踪提供者 p.tracerProvider 会被关闭
func (p *OTelProvider) Shutdown(ctx context.Context) error {
	if p.tracerProvider == nil {
		return nil
	}
	return p.tracerProvider.Shutdown(ctx)
}
