package database_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx/fxtest"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/database"
)

// testDSN points at a real, reachable PostgreSQL instance on port 5439,
// database asl_core_test - the CI workflow provides this as a service
// container. Locally, run one yourself (e.g. a second docker-compose
// Postgres on that port/db) or override via TEST_DATABASE_DSN.
var testDSN = getTestDSN()

func getTestDSN() string {
	if v := os.Getenv("TEST_DATABASE_DSN"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5439/asl_core_test?sslmode=disable"
}

func newParams(t *testing.T, dsn string) database.Params {
	t.Helper()
	t.Setenv("DATABASE_DSN", dsn)

	cfg, err := config.New()
	require.NoError(t, err)

	return database.Params{
		Lifecycle: fxtest.NewLifecycle(t),
		Config:    cfg,
	}
}

func TestNew(t *testing.T) {
	t.Run("invalid dsn", func(t *testing.T) {
		conn, err := database.New(newParams(t, "not a valid dsn"))
		require.Error(t, err)
		require.Nil(t, conn)
	})

	t.Run("unreachable server", func(t *testing.T) {
		conn, err := database.New(newParams(t, "postgres://postgres:postgres@127.0.0.1:1/nonexistent?sslmode=disable&connect_timeout=1"))
		require.Error(t, err)
		require.Nil(t, conn)
	})

	t.Run("connects and pings successfully", func(t *testing.T) {
		conn, err := database.New(newParams(t, testDSN))
		require.NoError(t, err)
		require.NotNil(t, conn)
	})
}

func TestConn_Queries(t *testing.T) {
	conn, err := database.New(newParams(t, testDSN))
	require.NoError(t, err)

	ctx := t.Context()

	t.Run("Exec", func(t *testing.T) {
		tag, err := conn.Exec(ctx, "SELECT 1")
		require.NoError(t, err)
		require.NotNil(t, tag)
	})

	t.Run("Query", func(t *testing.T) {
		rows, err := conn.Query(ctx, "SELECT 1")
		require.NoError(t, err)
		defer rows.Close()

		require.True(t, rows.Next())
		var result int
		require.NoError(t, rows.Scan(&result))
		require.Equal(t, 1, result)
	})

	t.Run("QueryRow", func(t *testing.T) {
		var result int
		err := conn.QueryRow(ctx, "SELECT 1").Scan(&result)
		require.NoError(t, err)
		require.Equal(t, 1, result)
	})
}
