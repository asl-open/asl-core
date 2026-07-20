package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/asl-open/asl-core/pkg/errorsx"
	"github.com/asl-open/asl-core/pkg/requestid"
	"github.com/asl-open/asl-core/pkg/response"
	"github.com/asl-open/asl-core/services/api/internal/http/reply"
)

// Errors recovers panics and converts any error attached via c.Error
// during the request into the standard response.Response envelope.
// Register it in place of gin.Recovery() - it subsumes that behavior (a
// panic still can't crash the process) and additionally produces the
// same JSON body a handler-reported error would, instead of an empty
// one.
func (m *middleware) Errors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				m.logger.Error(c.Request.Context(), "panic recovered", "panic", r)
				reply.JSON(c, response.Internal().WithRequestID(requestID(c)))
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 && !c.Writer.Written() {
			err := c.Errors.Last().Err

			if appErr, ok := errorsx.As(err); !ok || appErr.Internal() {
				m.logger.Error(c.Request.Context(), "request error", "error", err)
			}

			reply.JSON(c, response.Err(err).WithRequestID(requestID(c)))
		}
	}
}

// requestID returns the ID RequestID() attached to the request context.
// Falls back to the raw header in case Errors() ever runs without
// RequestID() ahead of it in the chain.
func requestID(c *gin.Context) string {
	if id, ok := requestid.FromContext(c.Request.Context()); ok {
		return id
	}
	return c.GetHeader(headerRequestID)
}
