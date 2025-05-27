package telemetry

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const redisSpanKey = contextKey("redis_span")

// RedisHook Redis 的调用链监控钩子
type RedisHook struct {
	Tracer trace.Tracer
}

// BeforeProcess 在命令执行前创建 span
func (h *RedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	spanName := "redis." + cmd.Name()
	ctx, span := h.Tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
			attribute.String("db.operation", cmd.Name()),
			attribute.String("db.statement", cmd.String()),
		))
	return context.WithValue(ctx, redisSpanKey, span), nil
}

// AfterProcess 在命令执行后结束 span
func (h *RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if span, ok := ctx.Value(redisSpanKey).(trace.Span); ok {
		defer span.End()
		if err := cmd.Err(); err != nil && err != redis.Nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}
	return nil
}

// BeforeProcessPipeline 在管道命令执行前创建 span
func (h *RedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	ctx, span := h.Tracer.Start(ctx, "redis.pipeline",
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
			attribute.String("db.operation", "pipeline"),
			attribute.Int("db.redis.num_cmd", len(cmds)),
		))
	return context.WithValue(ctx, redisSpanKey, span), nil
}

// AfterProcessPipeline 在管道命令执行后结束 span
func (h *RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if span, ok := ctx.Value(redisSpanKey).(trace.Span); ok {
		defer span.End()
		for _, cmd := range cmds {
			if err := cmd.Err(); err != nil && err != redis.Nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				break
			}
		}
	}
	return nil
}
