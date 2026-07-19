// Package http contains the HTTP transport layer. Gin types must not leak
// outside this package tree.
package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/services/api/internal/http/handlers"
	"github.com/asl-open/asl-core/services/api/internal/http/handlers/health"
	"github.com/asl-open/asl-core/services/api/internal/http/middleware"
)

var Module = fx.Options(
	handlers.Module,
	middleware.Module,
	fx.Invoke(New),
)

type Params struct {
	fx.In
	fx.Lifecycle

	Config        config.Config
	Middleware    middleware.Middleware
	HealthHandler health.Handler
}

func New(p Params) error {
	engine := gin.New()
	engine.Use(gin.Recovery(), p.Middleware.Logging())

	registerRoutes(engine, p)

	server := &http.Server{
		Handler: engine,
	}

	p.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := p.Config.GetString("http.addr")

			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("failed to listen on %s: %w", addr, err)
			}

			go func() {
				if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := server.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shut down http server: %w", err)
			}
			return nil
		},
	})

	return nil
}
