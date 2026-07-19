package health

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/logger"
	servicehealth "github.com/asl-open/asl-core/services/api/internal/services/health"
)

var Module = fx.Provide(New)

type Handler interface {
	// Health reports whether the process is running. It never checks
	// dependencies.
	Health(c *gin.Context)
	// Ready reports whether the service can actually serve traffic (its
	// dependencies are reachable).
	Ready(c *gin.Context)
}

type Params struct {
	fx.In

	Checker servicehealth.Checker
	Logger  logger.Logger
}

type handler struct {
	checker servicehealth.Checker
	logger  logger.Logger
}

func New(p Params) Handler {
	return &handler{checker: p.Checker, logger: p.Logger}
}
