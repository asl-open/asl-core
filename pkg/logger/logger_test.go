package logger

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/requestid"
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

func Test_withRequestID(t *testing.T) {
	t.Run("appends request_id when ctx carries one", func(t *testing.T) {
		ctx := requestid.NewContext(context.Background(), "req-123")

		got := withRequestID(ctx, []any{"key", "value"})

		require.Equal(t, []any{"key", "value", "request_id", "req-123"}, got)
	})

	t.Run("leaves fields unchanged when ctx carries none", func(t *testing.T) {
		got := withRequestID(context.Background(), []any{"key", "value"})

		require.Equal(t, []any{"key", "value"}, got)
	})

	t.Run("does not mutate the caller's slice", func(t *testing.T) {
		fields := make([]any, 2, 4) // spare capacity, so a naive append could alias it
		fields[0], fields[1] = "key", "value"
		ctx := requestid.NewContext(context.Background(), "req-123")

		_ = withRequestID(ctx, fields)

		require.Len(t, fields, 2, "must not grow the caller's slice in place")
	})
}

func Test_isIgnorableSyncError(t *testing.T) {
	t.Run("EINVAL wrapped like os.Sync on a pipe/character device", func(t *testing.T) {
		err := &os.PathError{Op: "sync", Path: "/dev/stdout", Err: syscall.EINVAL}
		require.True(t, isIgnorableSyncError(err))
	})

	t.Run("ENOTTY", func(t *testing.T) {
		err := &os.PathError{Op: "sync", Path: "/dev/stdout", Err: syscall.ENOTTY}
		require.True(t, isIgnorableSyncError(err))
	})

	t.Run("other errors are not ignored", func(t *testing.T) {
		require.False(t, isIgnorableSyncError(fmt.Errorf("disk full")))
	})
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
