package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusResponse is the JSON body returned by /health and /ready.
//
//	@Description	Response body shared by /health and /ready.
type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

// Health godoc
//
//	@Summary		Liveness check
//	@Description	Always returns 200 while the process is running. Does not check dependencies.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	StatusResponse	"process is alive"
//	@Router			/health [get]
//	@ID				healthCheck
func (h *handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, StatusResponse{Status: "ok"})
}
