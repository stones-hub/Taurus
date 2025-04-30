package middleware

import (
	"Taurus/pkg/contextx"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// TraceMiddleware generates a traceID for each request and adds it to the context
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique traceID
		traceID := uuid.New().String()

		// Create a RequestContext and add it to the context
		rc := &contextx.RequestContext{
			TraceID: traceID,
			AtTime:  time.Now(),
		}
		ctx := contextx.WithRequestContext(r.Context(), rc)
		// Log the traceID for debugging
		log.Printf("Incoming request, traceID: %s, path: %s\n", traceID, r.URL.Path)
		// Pass the context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))

		// Log the duration of the request
		duration := time.Since(rc.AtTime)
		log.Printf("Request completed, traceID: %s, duration: %s\n", traceID, duration)
	})
}
