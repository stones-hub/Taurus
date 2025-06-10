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
