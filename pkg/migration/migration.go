// Package migration wraps golang-migrate for applying, rolling back, and
// inspecting database migrations. It is intentionally not an Fx module -
// migrations are a deliberate, explicit action (see services/api/cmd/migrate),
// not something that runs automatically as part of the app's Fx graph.
package migration

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// New returns a *migrate.Migrate configured to read migrations from source
// (e.g. "file://migrations/api") and apply them against dsn.
func New(source, dsn string) (*migrate.Migrate, error) {
	return migrate.New(source, dsn)
}

// Up applies all pending migrations. Having nothing to apply is not an
// error.
func Up(m *migrate.Migrate) error {
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// Down rolls back the most recently applied migration. Having nothing to
// roll back is not an error.
func Down(m *migrate.Migrate) error {
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
