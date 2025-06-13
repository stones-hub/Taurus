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
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CorsMiddleware adds CORS headers to the response
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//  w.Header().Set("Access-Control-Allow-Origin", "https://your-allowed-origin.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		setCorsToTrace(r, "*", "GET, POST, PUT, DELETE, OPTIONS", "Content-Type, Authorization", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setCorsToTrace(r *http.Request, origin string, methods string, headers string, maxAge string) {
	if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
		span.SetAttributes(attribute.String("Access-Control-Allow-Origin", fmt.Sprintf("origin: %v, methods: %v, headers: %v, maxAge: %v", origin, methods, headers, maxAge)))
	}
}
