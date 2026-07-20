package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/logger"
)

var Module = fx.Provide(New)

type Params struct {
	fx.In

	Logger logger.Logger
}

type Middleware interface {
	Logging() gin.HandlerFunc
	Errors() gin.HandlerFunc
	RequestID() gin.HandlerFunc
}

type middleware struct {
	logger logger.Logger
}

func New(p Params) Middleware {
	return &middleware{logger: p.Logger}
}
