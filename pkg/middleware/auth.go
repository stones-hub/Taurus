package middleware

import (
	"Taurus/config"
	"Taurus/pkg/httpx"
	"net/http"
)

// ApiKeyAuthMiddleware validates the API key from the request headers
func ApiKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
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
