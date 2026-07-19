package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Logging logs each completed HTTP request.
func (m *middleware) Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		m.logger.Info(c.Request.Context(), "http request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
		)
	}
}
