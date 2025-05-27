package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"Taurus/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// User 用户模型
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TraceMiddleware 实现追踪中间件
func TraceMiddleware(tracer trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从请求中获取父 span 的上下文
			ctx := r.Context()

			// 创建新的 span
			spanName := "http." + r.Method + "." + r.URL.Path
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.path", r.URL.Path),
				),
			)
			defer span.End()

			// 记录请求开始时间
			startTime := time.Now()

			// 包装 ResponseWriter 以捕获状态码
			wrapped := wrapResponseWriter(w)

			// 调用下一个处理器
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// 记录响应信息
			duration := time.Since(startTime)
			span.SetAttributes(
				attribute.Int("http.status_code", wrapped.statusCode),
				attribute.String("http.duration", duration.String()),
			)
		})
	}
}

// responseWriter 包装 http.ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// UserHandler 处理用户相关的请求
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取用户ID
	id := strings.TrimPrefix(r.URL.Path, "/user/")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	// 获取当前 span
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("user.id", id))

	// 模拟用户数据
	user := User{
		ID:   1,
		Name: "test user",
		Age:  20,
	}

	// 返回 JSON 响应
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// 1. 初始化追踪器提供者
	provider, err := telemetry.NewOTelProvider(
		telemetry.WithServiceName("http-demo"),
		telemetry.WithServiceVersion("v0.1.0"),
		telemetry.WithEnvironment("dev"),
	)
	if err != nil {
		log.Fatalf("init telemetry provider failed: %v", err)
	}
	defer provider.Shutdown(context.Background())

	// 2. 获取追踪器
	tracer := provider.Tracer("http-server")

	// 3. 创建 HTTP 处理器
	mux := http.NewServeMux()

	// 使用追踪中间件包装处理器
	handler := TraceMiddleware(tracer)(http.HandlerFunc(UserHandler))
	mux.Handle("/user/", handler)

	// 4. 启动 HTTP 服务器
	log.Printf("HTTP server listening at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("run http server failed: %v", err)
	}
}
