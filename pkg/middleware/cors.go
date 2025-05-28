package middleware

import (
	"fmt"
	"log"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CorsMiddleware adds CORS headers to the response
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("-------------------------------- CorsMiddleware --------------------------------")
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
