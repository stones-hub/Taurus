package middleware

import (
	"Taurus/pkg/httpx"
	"Taurus/pkg/util"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter *util.CompositeRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr // 获取请求的IP地址
			allowed, message := limiter.Allow(ip)
			setRateLimitToTrace(r, allowed, message)
			if !allowed {
				httpx.SendResponse(w, http.StatusTooManyRequests, message, nil)
				return
			}
			// 如果请求被允许，继续处理下一个中间件或处理器
			next.ServeHTTP(w, r)
		})
	}
}

func setRateLimitToTrace(r *http.Request, allowed bool, message string) {
	if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
		span.SetAttributes(attribute.String("RateLimit", fmt.Sprintf("allowed: %v, message: %v", allowed, message)))
	}
}
