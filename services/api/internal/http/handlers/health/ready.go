package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const dependencyCheckTimeout = 2 * time.Second

func (h *handler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), dependencyCheckTimeout)
	defer cancel()

	if err := h.checker.Ready(ctx); err != nil {
		// The error is logged, not returned to the client - it could
		// contain internal details (host, port, driver internals).
		h.logger.Error(c.Request.Context(), "readiness check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
