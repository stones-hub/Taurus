// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

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
