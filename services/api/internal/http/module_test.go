package http_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/logger"
	apihttp "github.com/asl-open/asl-core/services/api/internal/http"
	"github.com/asl-open/asl-core/services/api/internal/http/handlers/ping"
	"github.com/asl-open/asl-core/services/api/internal/http/middleware"
)

func newTestMiddleware() middleware.Middleware {
	return middleware.New(middleware.Params{Logger: &logger.MockLogger{}})
}

func TestNew(t *testing.T) {
	t.Run("err on listen", func(t *testing.T) {
		listener, err := net.Listen("tcp", ":0")
		require.NoError(t, err)
		defer listener.Close()

		t.Setenv("HTTP_ADDR", listener.Addr().String())

		cfg, err := config.New()
		require.NoError(t, err)

		lifecycle := fxtest.NewLifecycle(t)

		err = apihttp.New(apihttp.Params{
			Lifecycle:   lifecycle,
			Config:      cfg,
			Middleware:  newTestMiddleware(),
			PingHandler: ping.NewMockHandler(),
		})
		require.NoError(t, err)

		err = lifecycle.Start(t.Context())
		require.Error(t, err)
	})

	t.Run("lifecycle start and stop", func(t *testing.T) {
		t.Setenv("HTTP_ADDR", ":0")

		cfg, err := config.New()
		require.NoError(t, err)

		pingHandler := ping.NewMockHandler()
		lifecycle := fxtest.NewLifecycle(t)

		err = apihttp.New(apihttp.Params{
			Lifecycle:   lifecycle,
			Config:      cfg,
			Middleware:  newTestMiddleware(),
			PingHandler: pingHandler,
		})
		require.NoError(t, err)

		err = lifecycle.Start(t.Context())
		require.NoError(t, err)

		err = lifecycle.Stop(t.Context())
		require.NoError(t, err)
	})
}
