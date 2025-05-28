package middleware

import (
	"Taurus/pkg/contextx"
	"crypto/md5"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// 设置一个包装来记录Response的最终状态

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// 重写WriteHeader方法，记录响应状态码, 意味着后续凡是自定义的responseWriter，都需要调用这个方法
func (rw *responseWriter) WriteHeader(code int) {
	log.Printf("trace middleware write header: %v", code)
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// TraceMiddleware 实现追踪中间件
func TraceMiddleware(tracer trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("-------------------------------- TraceMiddleware --------------------------------")
			// 用x-request-id作为traceID
			requestid := r.Header.Get("X-Request-ID")
			if requestid == "" {
				requestid = uuid.New().String()
			}

			// 使用 MD5 生成 16 字节的 TraceID, 因为调用链监控只支持16进制
			hash := md5.Sum([]byte(requestid))
			var traceID trace.TraceID
			copy(traceID[:], hash[:])

			rc := &contextx.RequestContext{ // 生成一个唯一的RequestContext，里面记录了traceID和请求开始时间
				TraceID: traceID.String(), // 生成一个唯一的 traceID
				AtTime:  time.Now(),       // 记录请求开始时间
			}
			// 将自定义的上下文添加到请求中
			ctx := contextx.WithRequestContext(r.Context(), rc)

			// 创建新的spanContext
			spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: traceID,
			})

			// 将SpanContext注入到上下文中
			ctx = trace.ContextWithSpanContext(ctx, spanCtx)

			// 创建新的 span，并记录请求的详细信息
			spanName := "http." + r.Method + "." + r.URL.Path
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.path", r.URL.Path),
					attribute.String("http.trace_id", rc.TraceID),
					attribute.String("http.at_time", rc.AtTime.Format(time.RFC3339)),
				),
			)
			defer span.End()

			// 包装ResponseWriter，记录响应状态码
			wrapped := wrapResponseWriter(w)

			// 调用下一个处理器
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// 记录响应信息
			duration := time.Since(rc.AtTime)
			log.Printf("duration: %v, statusCode: %v", duration, wrapped.statusCode)
			span.SetAttributes(
				attribute.String("http.duration", duration.String()),  // 记录请求的持续时间
				attribute.Int("http.status_code", wrapped.statusCode), // 记录响应状态码
			)
		})
	}
}
