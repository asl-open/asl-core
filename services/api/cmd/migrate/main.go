// Command migrate creates, applies, rolls back, or inspects database
// migrations. It is a separate binary from cmd/api on purpose - the server
// never migrates the database on its own, migrations are a deliberate,
// explicit action.
//
// Run from the repository root:
//
//	go run ./services/api/cmd/migrate create <name>
//	go run ./services/api/cmd/migrate up
//	go run ./services/api/cmd/migrate down
//	go run ./services/api/cmd/migrate version
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/migration"
)

// defaultMigrationSource lives here, not in pkg/config - it names this
// specific service's migration directory, so it isn't something the
// shared config package (used by every future service) should know about.
const defaultMigrationSource = "file://migrations/api"

func migrationSource(cfg config.Config) string {
	if v := cfg.GetString("migration.source"); v != "" {
		return v
	}
	return defaultMigrationSource
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <create|up|down|version>")
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cmd := os.Args[1]

	// create only touches the filesystem - it doesn't need a database
	// connection, unlike everything else below.
	if cmd == "create" {
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate create <name>")
		}

		dir := strings.TrimPrefix(migrationSource(cfg), "file://")
		up, down, err := migration.Create(dir, os.Args[2])
		if err != nil {
			log.Fatalf("failed to create migration: %v", err)
		}

		fmt.Printf("created %s\ncreated %s\n", up, down)
		return
	}

	// Errors below deliberately never include the DSN - it may contain
	// credentials.
	m, err := migration.New(migrationSource(cfg), cfg.GetString("database.dsn"))
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v", err)
	}

	switch cmd {
	case "up":
		err = migration.Up(m)
	case "down":
		err = migration.Down(m)
	case "version":
		version, dirty, verr := m.Version()
		if verr != nil {
			log.Fatalf("failed to get migration version: %v", verr)
		}
		fmt.Printf("version %d, dirty=%v\n", version, dirty)
		return
	default:
		log.Fatalf("unknown command %q (expected create, up, down or version)", cmd)
	}

	if err != nil {
		log.Fatalf("migration %s failed: %v", cmd, err)
	}

	fmt.Printf("migration %s completed\n", cmd)
}
