package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GormTracingHook GORM 的调用链监控钩子
type GormTracingHook struct {
	Tracer trace.Tracer
}

// Name 实现 gorm.Plugin 接口
func (h *GormTracingHook) Name() string {
	return "GormTracingHook"
}

// Initialize 实现 gorm.Plugin 接口， 给每个数据库操作添加 span， 利用 gorm 的回调机制， 在操作之前和之后添加 span
func (h *GormTracingHook) Initialize(db *gorm.DB) error {
	// 在查询之前开始 span
	_ = db.Callback().Create().Before("gorm:create").Register("tracing:before_create", h.before)
	_ = db.Callback().Query().Before("gorm:query").Register("tracing:before_query", h.before)
	_ = db.Callback().Delete().Before("gorm:delete").Register("tracing:before_delete", h.before)
	_ = db.Callback().Update().Before("gorm:update").Register("tracing:before_update", h.before)
	_ = db.Callback().Row().Before("gorm:row").Register("tracing:before_row", h.before)
	_ = db.Callback().Raw().Before("gorm:raw").Register("tracing:before_raw", h.before)

	// 在查询之后结束 span
	_ = db.Callback().Create().After("gorm:create").Register("tracing:after_create", h.after)
	_ = db.Callback().Query().After("gorm:query").Register("tracing:after_query", h.after)
	_ = db.Callback().Delete().After("gorm:delete").Register("tracing:after_delete", h.after)
	_ = db.Callback().Update().After("gorm:update").Register("tracing:after_update", h.after)
	_ = db.Callback().Row().After("gorm:row").Register("tracing:after_row", h.after)
	_ = db.Callback().Raw().After("gorm:raw").Register("tracing:after_raw", h.after)
	return nil
}

// before 在数据库操作之前创建 span
func (h *GormTracingHook) before(db *gorm.DB) {
	spanName := "gorm." + db.Statement.Schema.Table
	if db.Statement.Schema == nil {
		spanName = "gorm.raw"
	}

	ctx, span := h.Tracer.Start(db.Statement.Context, spanName,
		trace.WithAttributes(
			attribute.String("db.system", "mysql"),
			attribute.String("db.operation", db.Statement.Schema.Table),
			attribute.String("db.statement", db.Statement.SQL.String()),
		))
	db.Statement.Context = context.WithValue(ctx, "span", span)
}

// after 在数据库操作之后结束 span
func (h *GormTracingHook) after(db *gorm.DB) {
	if span, ok := db.Statement.Context.Value("span").(trace.Span); ok {
		defer span.End()
		if db.Error != nil {
			span.RecordError(db.Error)
			span.SetStatus(codes.Error, db.Error.Error())
		}
	}
}
