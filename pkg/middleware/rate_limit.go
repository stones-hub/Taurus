package middleware

import (
	"Taurus/pkg/httpx"
	"Taurus/pkg/util"
	"net/http"
)

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(next http.Handler, limiter *util.CompositeRateLimiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr // 获取请求的IP地址

		allowed, message := limiter.Allow(ip)
		if !allowed {
			httpx.SendErrorResponse(w, http.StatusTooManyRequests, message)
			return
		}

		// 如果请求被允许，继续处理下一个中间件或处理器
		next.ServeHTTP(w, r)

	})
}
