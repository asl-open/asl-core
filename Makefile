# Development commands for the ASL Core monorepo. Run `make help` for a
# list of available targets.

.DEFAULT_GOAL := help

MIGRATE_BIN := $(HOME)/go/bin/migrate
MIGRATIONS_DIR := migrations/api
GOLANGCI_LINT_BIN := $(HOME)/go/bin/golangci-lint
FIELDALIGNMENT_BIN := $(HOME)/go/bin/fieldalignment
MODERNIZE_BIN := $(HOME)/go/bin/modernize
SWAG_BIN := $(HOME)/go/bin/swag
VACUUM_BIN := $(HOME)/go/bin/vacuum
GITLEAKS_BIN := $(HOME)/go/bin/gitleaks
API_CMD := ./services/api/cmd
BUILD_OUT := bin/api
OPENAPI_DIR := services/api/internal/http/docs
OPENAPI_SPEC := $(OPENAPI_DIR)/swagger.yaml
SWAG_GENERAL := main.go
SWAG_DIRS := services/api/cmd,services/api/internal/http/handlers/health

.PHONY: help run build test test-race fmt fmt-check setup-lint lint setup-fieldalignment \
	fieldalignment fieldalignment-fix go-fix go-fix-check setup-modernize modernize \
	modernize-fix setup-swag openapi openapi-check setup-vacuum openapi-validate \
	setup-gitleaks secret-scan \
	setup-migrate migrate-create migrate-up migrate-down migrate-version \
	docker-up docker-down docker-logs

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

## test-race: run all tests with the race detector enabled
test-race:
	go test -race ./...

setup-lint:
	@if [ ! -x "$(GOLANGCI_LINT_BIN)" ]; then \
		echo "Installing golangci-lint CLI..."; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
	fi

## fmt: format all Go source files
fmt: setup-lint
	$(GOLANGCI_LINT_BIN) fmt ./...

## fmt-check: fail if any Go source file is not formatted, without changing files
fmt-check: setup-lint
	$(GOLANGCI_LINT_BIN) fmt ./... --diff

## lint: run the configured linters
lint: setup-lint
	$(GOLANGCI_LINT_BIN) run ./...

setup-gitleaks:
	@if [ ! -x "$(GITLEAKS_BIN)" ]; then \
		echo "Installing gitleaks CLI..."; \
		go install github.com/zricethezav/gitleaks/v8@latest; \
	fi

## secret-scan: scan the git history for committed secrets (gitleaks)
secret-scan: setup-gitleaks
	$(GITLEAKS_BIN) detect --source . --config .gitleaks.toml --redact --verbose

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

## go-fix: apply the standard library's suggested API fixes (go tool fix)
go-fix:
	go fix ./...

## go-fix-check: fail if go fix would change any file, without changing files
go-fix-check:
	go fix -diff ./...

setup-modernize:
	@if [ ! -x "$(MODERNIZE_BIN)" ]; then \
		echo "Installing modernize CLI..."; \
		go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest; \
	fi

## modernize: report code that could use newer Go language/library features
modernize: setup-modernize
	$(MODERNIZE_BIN) -diff ./...

## modernize-fix: rewrite code in place to use newer Go language/library features
modernize-fix: setup-modernize
	$(MODERNIZE_BIN) -fix ./...

setup-swag:
	@if [ ! -x "$(SWAG_BIN)" ]; then \
		echo "Installing swag CLI..."; \
		go install github.com/swaggo/swag/v2/cmd/swag@latest; \
	fi

## openapi: regenerate the OpenAPI spec from swag annotations
openapi: setup-swag
	$(SWAG_BIN) init -g $(SWAG_GENERAL) -d $(SWAG_DIRS) -o $(OPENAPI_DIR) --ot yaml --v3.1 --parseInternal

## openapi-check: fail if the generated OpenAPI spec is out of date, without changing tracked files
openapi-check: setup-swag
	@tmp=$$(mktemp -d); \
	$(SWAG_BIN) init -g $(SWAG_GENERAL) -d $(SWAG_DIRS) -o $$tmp --ot yaml --v3.1 --parseInternal > /dev/null 2>&1; \
	if diff -u $(OPENAPI_SPEC) $$tmp/swagger.yaml; then \
		echo "OpenAPI spec is up to date"; \
		rm -rf $$tmp; \
	else \
		echo "OpenAPI spec is out of date - run 'make openapi' to regenerate"; \
		rm -rf $$tmp; \
		exit 1; \
	fi

setup-vacuum:
	@if [ ! -x "$(VACUUM_BIN)" ]; then \
		echo "Installing vacuum CLI..."; \
		go install github.com/daveshanley/vacuum@latest; \
	fi

## openapi-validate: validate the OpenAPI spec
openapi-validate: setup-vacuum
	$(VACUUM_BIN) lint $(OPENAPI_SPEC)

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
