# api service

The first runnable ASL Core service. No business logic yet — this is
foundation/bootstrap work (Fx wiring, HTTP transport, configuration).

## Running

```
go run ./cmd
```

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

Multi-stage build (`golang:1.25-alpine` builder, `gcr.io/distroless/static-debian12:nonroot`
runtime) - no shell, no Go toolchain, no package manager in the final
image, runs as `nonroot`. The server is PID 1 and exits on `SIGTERM`
(what `docker stop` sends), so shutdown is graceful the same way it is
outside a container.

## Endpoints

- `GET /health` — always `200 {"status":"ok"}` while the process is
  running. Does not check dependencies.
- `GET /ready` — `200 {"status":"ok"}` if dependencies (currently just the
  database) are reachable, `503 {"status":"unavailable"}` otherwise.
  Dependency checks have a 2s timeout. The underlying error is logged
  server-side, never returned to the client.

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
│   ├── docs/                      OpenAPI spec (not added yet)
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
