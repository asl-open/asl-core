package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/asl-open/asl-core/pkg/errorsx"
	"github.com/asl-open/asl-core/pkg/response"
	"github.com/asl-open/asl-core/services/api/internal/http/reply"
)

// requestIDKey is the gin.Context key a request-ID middleware should set
// so error responses can include it. Nothing sets it yet - until a
// request-ID middleware lands, requestID falls back to echoing the
// X-Request-Id header if the client sent one.
const requestIDKey = "request_id"

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

func requestID(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return c.GetHeader("X-Request-Id")
}
