package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"

	"github.com/asl-open/asl-core/pkg/config"
)

func TestNew(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)

	log, err := New(Params{
		Config:    cfg,
		Lifecycle: fxtest.NewLifecycle(t),
	})
	require.NoError(t, err)
	require.NotNil(t, log)

	ctx := context.Background()

	log.Info(ctx, "test logger info message")
	log.Debug(ctx, "test logger debug message")
	log.Warn(ctx, "test logger warn message")
	log.Error(ctx, "test logger error message")
}

func Test_getLevel(t *testing.T) {
	t.Run("default level", func(t *testing.T) {
		cfg, err := config.New()
		require.NoError(t, err)

		require.Equal(t, "info", getLevel(cfg).String())
	})

	cases := map[string]string{
		"debug":   "debug",
		"info":    "info",
		"warning": "warn",
		"error":   "error",
	}

	for raw, want := range cases {
		t.Run(raw, func(t *testing.T) {
			t.Setenv("LOGGER_LEVEL", raw)

			cfg, err := config.New()
			require.NoError(t, err)

			require.Equal(t, want, getLevel(cfg).String())
		})
	}
}
