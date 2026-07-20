// Package requestid carries a per-request identifier through a plain
// context.Context, so anything holding the context - a handler, a
// service, pkg/logger - can read it without depending on the HTTP
// layer. The HTTP-facing middleware that populates it lives in
// services/api/internal/http/middleware.
package requestid

import "context"

type ctxKey struct{}

// NewContext returns a copy of ctx carrying id.
func NewContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// FromContext returns the request ID carried by ctx, if any.
func FromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKey{}).(string)
	return id, ok
}
