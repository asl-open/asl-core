// Package http contains the HTTP transport layer. Gin types must not leak
// outside this package tree.
package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/logger"
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
	Logger        logger.Logger
	Middleware    middleware.Middleware
	HealthHandler health.Handler
}

func New(p Params) error {
	engine := gin.New()
	engine.Use(p.Middleware.RequestID(), p.Middleware.Errors(), p.Middleware.Logging())

	registerRoutes(engine, p)

	server := &http.Server{
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	p.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := p.Config.GetString("http.addr")

			var lc net.ListenConfig
			listener, err := lc.Listen(ctx, "tcp", addr)
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
		OnStop: shutdownHook(p.Logger, server.Shutdown),
	})

	return nil
}

// shutdownHook wraps shutdown (typically (*http.Server).Shutdown) into an
// fx OnStop hook that logs and wraps any error, rather than letting it
// surface only through Fx's own event logger.
func shutdownHook(log logger.Logger, shutdown func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		if err := shutdown(ctx); err != nil {
			wrapped := fmt.Errorf("failed to shut down http server: %w", err)
			log.Error(ctx, "http server shutdown failed", "error", wrapped)
			return wrapped
		}
		return nil
	}
}
