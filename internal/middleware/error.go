package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// ErrorHandlerMiddleware handles errors and recovers from panics in HTTP requests
func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				log.Printf("Recovered from panic: %v\n%s", err, debug.Stack())

				// Respond with a generic error message
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
