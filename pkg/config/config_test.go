package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/config"
)

func TestNew_Defaults(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)
	require.Equal(t, ":8080", cfg.GetString("http.addr"))
}

func TestNew_EnvOverride(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")

	cfg, err := config.New()
	require.NoError(t, err)
	require.Equal(t, ":9090", cfg.GetString("http.addr"))
}
