package ping

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Handler interface {
	Ping(c *gin.Context)
}

type handler struct{}

func New() Handler {
	return &handler{}
}
