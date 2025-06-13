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
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ErrorHandlerMiddleware handles errors and recovers from panics in HTTP requests
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				log.Printf("Recovered from panic: %v\n%s", err, debug.Stack())
				setErrorToTrace(r, err)
				httpx.SendResponse(w, http.StatusInternalServerError, "Internal Server Error", nil)
			}
		}()
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func setErrorToTrace(r *http.Request, err interface{}) {
	if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
		span.SetAttributes(attribute.String("Error", fmt.Sprintf("%v", err)))
	}
}
