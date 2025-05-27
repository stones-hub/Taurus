package main

import (
	"context"
	"log"
	"net/http"

	"Taurus/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// tracingMiddleware 创建追踪中间件
func tracingMiddleware(provider telemetry.TracerProvider) func(http.Handler) http.Handler {
	tracer := provider.Tracer("http.server")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 从请求中提取 span context
			ctx := r.Context()

			// 创建新的 span
			opts := []trace.SpanStartOption{
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.target", r.URL.Path),
					attribute.String("http.scheme", r.URL.Scheme),
					attribute.String("http.host", r.Host),
					attribute.String("http.user_agent", r.UserAgent()),
				),
				trace.WithSpanKind(trace.SpanKindServer),
			}

			ctx, span := tracer.Start(ctx, r.URL.Path, opts...)
			defer span.End()

			// 将 span context 传递给下一个处理器
			r = r.WithContext(ctx)

			// 调用下一个处理器
			next.ServeHTTP(w, r)

			span.SetStatus(codes.Error, http.StatusText(http.StatusOK))
		})
	}
}

func main() {
	// 初始化 provider
	provider, err := telemetry.NewOTelProvider(
		telemetry.WithServiceName("http-demo"),
		telemetry.WithEnvironment("dev"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer provider.Shutdown(context.Background())

	// 创建路由
	mux := http.NewServeMux()

	// 注册路由处理器
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// 创建带追踪的处理器
	handler := tracingMiddleware(provider)(mux)

	// 启动服务器
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
