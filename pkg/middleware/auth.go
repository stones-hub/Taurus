package middleware

import (
	"Taurus/config"
	"Taurus/pkg/httpx"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ApiKeyAuthMiddleware validates the API key from the request headers
func ApiKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("-------------------------------- ApiKeyAuthMiddleware --------------------------------")
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
