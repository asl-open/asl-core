# Development commands for the ASL Core monorepo. Run `make help` for a
# list of available targets.

.DEFAULT_GOAL := help

MIGRATE_BIN := $(HOME)/go/bin/migrate
MIGRATIONS_DIR := migrations/api
GOLANGCI_LINT_BIN := $(HOME)/go/bin/golangci-lint
FIELDALIGNMENT_BIN := $(HOME)/go/bin/fieldalignment
API_CMD := ./services/api/cmd
BUILD_OUT := bin/api

.PHONY: help run build test fmt setup-lint lint setup-fieldalignment fieldalignment \
	fieldalignment-fix setup-migrate migrate-create migrate-up migrate-down \
	migrate-version docker-up docker-down docker-logs

## help: show this help
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /make /' | sort

## run: run the api service locally. Requires DATABASE_DSN (see services/api/.env.example)
run:
	go run $(API_CMD)

## build: build the api service binary to bin/api
build:
	go build -o $(BUILD_OUT) $(API_CMD)

## test: run all tests
test:
	go test ./...

setup-lint:
	@if [ ! -x "$(GOLANGCI_LINT_BIN)" ]; then \
		echo "Installing golangci-lint CLI..."; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
	fi

## fmt: format all Go source files
fmt: setup-lint
	$(GOLANGCI_LINT_BIN) fmt ./...

## lint: run the configured linters
lint: setup-lint
	$(GOLANGCI_LINT_BIN) run ./...

setup-fieldalignment:
	@if [ ! -x "$(FIELDALIGNMENT_BIN)" ]; then \
		echo "Installing fieldalignment CLI..."; \
		go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest; \
	fi

## fieldalignment: report structs that could shrink by reordering fields
fieldalignment: setup-fieldalignment
	$(FIELDALIGNMENT_BIN) ./...

## fieldalignment-fix: rewrite structs in place to minimize padding
fieldalignment-fix: setup-fieldalignment
	$(FIELDALIGNMENT_BIN) -fix ./...

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

## docker-up: start the local stack (Postgres + api) with Docker Compose
docker-up:
	docker compose up -d --build

## docker-down: stop the local Docker Compose stack
docker-down:
	docker compose down

## docker-logs: tail logs from the Docker Compose stack
docker-logs:
	docker compose logs -f
