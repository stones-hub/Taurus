// Package telemetry 提供了分布式追踪的核心接口定义。
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Tracer 定义追踪器接口，提供基本的追踪操作。
// 实现此接口可以创建自定义的追踪器，用于特定组件或场景的追踪。
type Tracer interface {
	// Start 开始一个新的span，返回带有span的上下文和span本身。
	// ctx: 父上下文
	// name: span名称
	// opts: span选项，如标签、事件等
	Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)

	// Extract 从载体中提取追踪上下文。
	// ctx: 父上下文
	// carrier: 携带追踪信息的载体，如HTTP头
	// 返回: 包含追踪信息的新上下文
	Extract(ctx context.Context, carrier interface{}) context.Context

	// Inject 将追踪上下文注入到载体中。
	// ctx: 包含追踪信息的上下文
	// carrier: 用于携带追踪信息的载体，如HTTP头
	Inject(ctx context.Context, carrier interface{})
}

// TracerProvider 定义追踪器提供者接口。
// 负责创建和管理追踪器实例，以及处理追踪器的生命周期。
type TracerProvider interface {
	// Tracer 返回一个命名的追踪器实例。
	// name: 追踪器名称，通常是组件或模块名
	// 返回: 追踪器实例
	Tracer(name string) Tracer

	// Shutdown 关闭追踪器提供者，确保所有数据都被导出。
	// ctx: 用于控制关闭操作的上下文
	// 返回: 关闭过程中的错误，如果有的话
	Shutdown(ctx context.Context) error
}

// Component 定义可追踪组件接口。
// 实现此接口的组件可以自动集成到追踪系统中。
type Component interface {
	// Name 返回组件名称。
	// 返回: 组件的唯一标识名称
	Name() string

	// Init 初始化组件的追踪功能。
	// provider: 追踪器提供者
	// 返回: 初始化过程中的错误，如果有的话
	Init(provider TracerProvider) error
}
