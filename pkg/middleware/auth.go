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
	"Taurus/config"
	"Taurus/pkg/httpx"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ApiKeyAuthMiddleware validates the API key from the request headers
func ApiKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		setAuthorizationToTrace(r, apiKey)
		if apiKey == "" {
			httpx.SendResponse(w, http.StatusUnauthorized, "Authorization is empty", nil)
			return
		}

		// Validate the API key (this is a placeholder, replace with actual validation logic)
		if !isValidApiKey(apiKey) {
			httpx.SendResponse(w, http.StatusUnauthorized, "Authorization is invalid", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isValidAPIKey checks if the provided API key is valid
func isValidApiKey(apiKey string) bool {
	// Implement your API key validation logic here
	if apiKey == config.Core.Authorization {
		return true
	} else {
		return false
	}
}

// 将Authorization添加到trace中
func setAuthorizationToTrace(r *http.Request, authorization string) {
	if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
		span.SetAttributes(attribute.String("Authorization", authorization))
	}
}
