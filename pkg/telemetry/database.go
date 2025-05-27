package telemetry

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// WrapMySQL 包装 MySQL 驱动以支持追踪
func WrapMySQL(dsn string) (*sql.DB, error) {
	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// WrapRedis 包装 Redis 客户端以支持追踪
func WrapRedis(opts *redis.Options) *TracedRedis {
	client := redis.NewClient(opts)
	return NewTracedRedis(client)
}

// TracedDB 带追踪的数据库操作包装器
type TracedDB struct {
	*sql.DB
	tracer trace.Tracer
}

// NewTracedDB 创建带追踪的数据库操作包装器
func NewTracedDB(db *sql.DB) *TracedDB {
	return &TracedDB{
		DB:     db,
		tracer: otel.Tracer("mysql"),
	}
}

// QueryContext 执行查询并追踪
func (db *TracedDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx, span := db.tracer.Start(ctx, "mysql.query",
		trace.WithAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.statement", query),
		))
	defer span.End()

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return rows, err
}

// ExecContext 执行命令并追踪
func (db *TracedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, span := db.tracer.Start(ctx, "mysql.exec",
		trace.WithAttributes(
			semconv.DBSystemMySQL,
			attribute.String("db.statement", query),
		))
	defer span.End()

	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return result, err
}

// TracedRedis Redis 客户端的追踪包装器
type TracedRedis struct {
	*redis.Client
	tracer trace.Tracer
}

// NewTracedRedis 创建带追踪的 Redis 客户端
func NewTracedRedis(client *redis.Client) *TracedRedis {
	return &TracedRedis{
		Client: client,
		tracer: otel.Tracer("redis"),
	}
}

// Process 执行单个 Redis 命令并追踪
func (c *TracedRedis) Process(ctx context.Context, cmd redis.Cmder) error {
	ctx, span := c.tracer.Start(ctx, "redis."+cmd.Name(),
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
			attribute.String("db.operation", cmd.Name()),
		))
	defer span.End()

	err := c.Client.Process(ctx, cmd)
	if err != nil && err != redis.Nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

// Pipeline 执行 Redis 管道命令并追踪
func (c *TracedRedis) Pipeline() redis.Pipeliner {
	return &TracedPipeliner{
		Pipeliner: c.Client.Pipeline(),
		tracer:    c.tracer,
	}
}

// TracedPipeliner Redis 管道命令的追踪包装器
type TracedPipeliner struct {
	redis.Pipeliner
	tracer trace.Tracer
}

// Exec 执行管道命令并追踪
func (p *TracedPipeliner) Exec(ctx context.Context) ([]redis.Cmder, error) {
	ctx, span := p.tracer.Start(ctx, "redis.pipeline",
		trace.WithAttributes(
			attribute.String("db.system", "redis"),
		))
	defer span.End()

	cmds, err := p.Pipeliner.Exec(ctx)
	if err != nil && err != redis.Nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return cmds, err
}
