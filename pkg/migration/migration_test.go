package migration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/migration"
)

// testDSN points at a real, reachable PostgreSQL instance - same
// requirement as pkg/database's tests. Overridable via TEST_DATABASE_DSN.
func testDSN() string {
	if v := os.Getenv("TEST_DATABASE_DSN"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5439/asl_core_test?sslmode=disable"
}

func TestNew(t *testing.T) {
	t.Run("invalid source", func(t *testing.T) {
		m, err := migration.New("not-a-valid-source", testDSN())
		require.Error(t, err)
		require.Nil(t, m)
	})

	t.Run("valid source and dsn", func(t *testing.T) {
		m, err := migration.New("file://../../migrations/api", testDSN())
		require.NoError(t, err)
		require.NotNil(t, m)
	})
}

func TestUpDown(t *testing.T) {
	m, err := migration.New("file://../../migrations/api", testDSN())
	require.NoError(t, err)

	require.NoError(t, migration.Up(m))
	require.NoError(t, migration.Up(m), "applying again should be a no-op, not an error")

	version, dirty, err := m.Version()
	require.NoError(t, err)
	require.False(t, dirty)
	require.Equal(t, uint(1), version)

	require.NoError(t, migration.Down(m))
	require.NoError(t, migration.Down(m), "rolling back again should be a no-op, not an error")
}
