# Only migration targets exist for now - run/build/test/lint/docker
# targets land with the full Makefile (#13).

MIGRATE_BIN := $(HOME)/go/bin/migrate
MIGRATIONS_DIR := migrations/api

.PHONY: setup-migrate migrate-create migrate-up migrate-down migrate-version

setup-migrate:
	@if [ ! -x "$(MIGRATE_BIN)" ]; then \
		echo "Installing golang-migrate CLI..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi

## migrate-create: scaffold a new migration pair. Usage: make migrate-create name=add_users_table
migrate-create: setup-migrate
	@if [ -z "$(name)" ]; then echo "usage: make migrate-create name=<description>"; exit 1; fi
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq -digits 6 $(name)

## migrate-up: apply all pending migrations. Requires DATABASE_DSN.
migrate-up: setup-migrate
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" up

## migrate-down: roll back the most recently applied migration. Requires DATABASE_DSN.
migrate-down: setup-migrate
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" down 1

## migrate-version: show the current migration version. Requires DATABASE_DSN.
migrate-version: setup-migrate
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DATABASE_DSN)" version
