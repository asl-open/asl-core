package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/asl-open/asl-core/pkg/requestid"
)

const headerRequestID = "X-Request-Id"

const maxRequestIDLength = 128

// RequestID ensures every request has one: an existing, valid
// X-Request-Id request header is preserved as-is (so an upstream
// system's own correlation ID survives, whatever format it uses),
// otherwise a new UUID is generated. Either way, the ID is attached to
// the request context - so pkg/logger picks it up automatically in
// every log line from here on, with no further wiring needed - and
// echoed back as a response header.
func (m *middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(headerRequestID)
		if !isValidRequestID(id) {
			id = uuid.NewString()
		}

		c.Request = c.Request.WithContext(requestid.NewContext(c.Request.Context(), id))
		c.Writer.Header().Set(headerRequestID, id)

		c.Next()
	}
}

// isValidRequestID reports whether id is safe and reasonable to accept
// from a client as-is: non-empty, not absurdly long (avoids response
// header/log abuse), and free of control characters (avoids log
// injection via an embedded newline etc). It does not require id to be
// a UUID - accepting whatever correlation ID format an upstream system
// already uses is more useful than forcing everything to be a UUID.
func isValidRequestID(id string) bool {
	if id == "" || len(id) > maxRequestIDLength {
		return false
	}

	for _, r := range id {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}

	return true
}
