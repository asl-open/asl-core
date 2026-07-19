package ping

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping is a temporary test endpoint used to verify the HTTP server is wired
// up correctly. It will be replaced by the real health endpoint in a later
// issue.
func (h *handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
