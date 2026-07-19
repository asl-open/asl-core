// Package migration wraps golang-migrate for applying, rolling back, and
// inspecting database migrations. It is intentionally not an Fx module -
// migrations are a deliberate, explicit action (see services/api/cmd/migrate),
// not something that runs automatically as part of the app's Fx graph.
package migration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

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

var sequencePrefix = regexp.MustCompile(`^(\d{6})_`)

// Create adds a new empty up/down migration pair to dir (a filesystem
// path, not a "file://" source URL), named NNNNNN_name.{up,down}.sql,
// continuing dir's existing sequence. It returns the created files' paths.
func Create(dir, name string) (up, down string, err error) {
	next, err := nextSequence(dir)
	if err != nil {
		return "", "", err
	}

	base := fmt.Sprintf("%06d_%s", next, name)
	up = filepath.Join(dir, base+".up.sql")
	down = filepath.Join(dir, base+".down.sql")

	for _, path := range []string{up, down} {
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			return "", "", fmt.Errorf("failed to create %s: %w", path, err)
		}
	}

	return up, down, nil
}

func nextSequence(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("failed to read migrations directory %s: %w", dir, err)
	}

	highest := 0
	for _, entry := range entries {
		match := sequencePrefix.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}

		n, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}

		if n > highest {
			highest = n
		}
	}

	return highest + 1, nil
}
