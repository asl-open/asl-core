// Command migrate applies, rolls back, or inspects database migrations. It
// is a separate binary from cmd/api on purpose - the server never migrates
// the database on its own, migrations are a deliberate, explicit action.
//
// Run from the repository root:
//
//	go run ./services/api/cmd/migrate up
//	go run ./services/api/cmd/migrate down
//	go run ./services/api/cmd/migrate version
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/asl-open/asl-core/pkg/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|version>")
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Errors below deliberately never include the DSN - it may contain
	// credentials.
	m, err := migrate.New(cfg.GetString("migration.source"), cfg.GetString("database.dsn"))
	if err != nil {
		log.Fatalf("failed to initialize migrate: %v", err)
	}

	switch cmd := os.Args[1]; cmd {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "version":
		version, dirty, verr := m.Version()
		if verr != nil {
			log.Fatalf("failed to get migration version: %v", verr)
		}
		fmt.Printf("version %d, dirty=%v\n", version, dirty)
		return
	default:
		log.Fatalf("unknown command %q (expected up, down or version)", cmd)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migration %s failed: %v", os.Args[1], err)
	}

	fmt.Printf("migration %s completed\n", os.Args[1])
}
