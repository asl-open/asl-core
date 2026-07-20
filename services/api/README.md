# api service

The first runnable ASL Core service. No business logic yet — this is
foundation/bootstrap work (Fx wiring, HTTP transport, configuration).

## Running

```
go run ./cmd
```

or, from the repository root, `make run` (`make build` builds a binary
to `bin/api`, `make test` runs the test suite, `make test-race` runs it
with the race detector - see `make help` for the full list of targets,
and [`docs/testing.md`](../../docs/testing.md) for testing conventions).

Configuration is loaded from environment variables (see
[`.env.example`](.env.example)). `HTTP_ADDR` overrides the listen address
(defaults to `:8080`), `LOGGER_LEVEL`/`LOGGER_FORMAT` control logging,
`DATABASE_DSN` is required (PostgreSQL connection string):

```
DATABASE_DSN=postgres://postgres:postgres@localhost:5432/asl_core?sslmode=disable go run ./cmd
```

Run migrations first (see [`../../migrations/README.md`](../../migrations/README.md)):

```
make migrate-up
```

### Docker

Build from the repository root (the image needs `pkg/` too, not just this
service):

```
docker build -f services/api/Dockerfile -t asl-core-api .
docker run --rm -e DATABASE_DSN=... -e HTTP_ADDR=:8080 -p 8080:8080 asl-core-api
```

Multi-stage build (`golang:1.26-alpine` builder, `alpine:3.19` runtime) -
no Go toolchain in the final image, runs as the unprivileged `app` user
(~33MB total). The server is PID 1 and exits on `SIGTERM` (what
`docker stop` sends), so shutdown is graceful the same way it is outside
a container.

### Docker Compose

For local development, `docker-compose.yml` at the repository root runs
this service alongside PostgreSQL, with networking and environment
variables preconfigured:

```
docker compose up -d --build
```

or `make docker-up` (`make docker-down` / `make docker-logs`). The API
is reachable at `localhost:8080`, PostgreSQL at `localhost:5432` (data
persists in the `postgres-data` named volume). If either port is already
taken on your machine, override it instead of editing the file:

```
API_PORT=8090 POSTGRES_PORT=5435 docker compose up -d --build
```

Run migrations against the compose database with
`DATABASE_DSN=postgres://postgres:postgres@localhost:5432/asl_core?sslmode=disable make migrate-up`
(adjust the port if overridden). `docker compose down` (or
`make docker-down`) stops everything cleanly; add `-v` to also drop the
named volume.

## Endpoints

- `GET /health` — always `200 {"status":"ok"}` while the process is
  running. Does not check dependencies.
- `GET /ready` — `200 {"status":"ok"}` if dependencies (currently just the
  database) are reachable, `503 {"status":"unavailable"}` otherwise.
  Dependency checks have a 2s timeout. The underlying error is logged
  server-side, never returned to the client.
- `GET /openapi.yaml` — the generated OpenAPI spec (see
  [`../../docs/api.md`](../../docs/api.md)).
- `GET /docs` — a RapiDoc page rendering the spec, for browsing the API
  locally.

## Layout

```
cmd/                         server entry point (fx.New(internal.Module).Run())
internal/
├── gateway/                  clients for external services (not implemented yet)
├── repository/                domain entities and data access (not implemented yet)
├── services/                   business logic, aggregated in module.go
│   └── health/                  Checker.Ready(ctx) error - currently just pings
│                                the database; the one place that knows what
│                                "ready" means as more dependencies are added
├── http/                       HTTP transport - Gin types stay inside this tree
│   ├── module.go                builds the Gin engine, wraps it in *http.Server,
│   │                            registers Fx lifecycle hooks (listen on OnStart,
│   │                            graceful shutdown on OnStop)
│   ├── routes.go                 the single place listing every route -> handler mapping
│   ├── docs/                      generated OpenAPI spec (swagger.yaml) -
│   │                              see docs/api.md at the repo root
│   ├── middleware/                 middleware.Middleware: Logging() logs each
│   │                              completed request (recovery uses gin.Recovery())
│   └── handlers/
│       ├── module.go               aggregates every handler's fx module
│       └── health/                  Handler.Health/Ready - talks to
│                                    internal/services/health, never touches
│                                    pkg/database directly (handlers stay above
│                                    the service layer, not beside it)
└── module.go                  aggregates config.Module + logger.Module +
                              database.Module + services.Module + http.Module
```

The PostgreSQL pool (`pkg/database`, not `internal/database` - shared
infra like `pkg/config`/`pkg/logger`) is wired eagerly via `fx.Invoke` so a
bad `DATABASE_DSN` or an unreachable database fails startup immediately.

## Conventions

- **Layering**: each concern (`gateway`, `repository`, `services`, `http`)
  is its own package. Handlers only implement their request logic (e.g.
  `Ready(c *gin.Context)`) and know nothing about `*gin.Engine` - route
  registration is centralized in `internal/http/routes.go`. Handlers depend
  on `internal/services/*`, never directly on `pkg/database` or other
  infra - business/dependency logic belongs in the service layer.
- **Adding a handler**: create a subpackage under `internal/http/handlers/`
  with a `Handler` interface, a struct implementing it, a constructor, and
  `var Module = fx.Provide(New)` (see `handlers/health` as a template). Add
  it to `handlers/module.go`'s `fx.Options(...)` list and wire its route in
  `routes.go`.
- **Fx modules**: every package that needs to be wired into the app exposes
  its own `Module` (`fx.Provide`/`fx.Options`/`fx.Module`), aggregated one
  level up. `cmd/main.go` only ever calls `fx.New(internal.Module).Run()`.
- **Config**: typed via `pkg/config` (env vars only, no config file). See
  that package for how to add a new setting.
- **Shared infra lives in `pkg/`, not `internal/`**: `pkg/config`, `pkg/logger`
  and `pkg/database` aren't specific to the `api` service - anything a future
  service would also need goes there instead of under `services/api/internal`.
