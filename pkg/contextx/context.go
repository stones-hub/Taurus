package contextx

import (
	"context"
	"time"
)

type RequestContext struct {
	TraceID string
	AtTime  time.Time
	// Future fields for statistics or other metadata
	// UserID string
}

// ContextKey is a custom type to avoid context key collisions
type ContextKey string

const requestContextKey ContextKey = "requestContext"

// WithRequestContext adds a RequestContext to the context
func WithRequestContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, requestContextKey, rc)
}

// GetRequestContext retrieves the RequestContext from the context
func GetRequestContext(ctx context.Context) (*RequestContext, bool) {
	rc, ok := ctx.Value(requestContextKey).(*RequestContext)
	return rc, ok
}
