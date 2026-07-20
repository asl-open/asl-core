package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const dependencyCheckTimeout = 2 * time.Second

// Ready godoc
//
//	@Summary		Readiness check
//	@Description	Returns 200 if dependencies (currently just the database) are reachable, 503 otherwise. Dependency checks have a 2s timeout. The underlying error is logged server-side, never returned to the client.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	StatusResponse	"dependencies are reachable"
//	@Failure		503	{object}	StatusResponse	"a dependency is unreachable"
//	@Router			/ready [get]
//	@ID				readinessCheck
func (h *handler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), dependencyCheckTimeout)
	defer cancel()

	if err := h.checker.Ready(ctx); err != nil {
		// The error is logged, not returned to the client - it could
		// contain internal details (host, port, driver internals).
		h.logger.Error(c.Request.Context(), "readiness check failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, StatusResponse{Status: "unavailable"})
		return
	}

	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
