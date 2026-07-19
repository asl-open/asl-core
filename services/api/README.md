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

## Layout

```
cmd/                         server entry point (fx.New(internal.Module).Run())
internal/
├── gateway/                  clients for external services (not implemented yet)
├── repository/                domain entities and data access (not implemented yet)
├── services/                  business logic and application errors (not implemented yet)
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
│       └── ping/                    example handler: Handler interface, struct, constructor,
│                                    its own fx.Provide module, and a testify MockHandler
└── module.go                  aggregates config.Module + logger.Module +
                              database.Module + http.Module
```

The PostgreSQL pool (`pkg/database`, not `internal/database` - shared
infra like `pkg/config`/`pkg/logger`) is wired eagerly via `fx.Invoke` so a
bad `DATABASE_DSN` or an unreachable database fails startup immediately.

## Conventions

- **Layering**: each concern (`gateway`, `repository`, `services`, `http`)
  is its own package. Handlers only implement their request logic (e.g.
  `Ping(c *gin.Context)`) and know nothing about `*gin.Engine` - route
  registration is centralized in `internal/http/routes.go`.
- **Adding a handler**: create a subpackage under `internal/http/handlers/`
  with a `Handler` interface, a struct implementing it, a constructor, and
  `var Module = fx.Provide(New)` (see `handlers/ping` as a template). Add it
  to `handlers/module.go`'s `fx.Options(...)` list and wire its route in
  `routes.go`.
- **Fx modules**: every package that needs to be wired into the app exposes
  its own `Module` (`fx.Provide`/`fx.Options`/`fx.Module`), aggregated one
  level up. `cmd/main.go` only ever calls `fx.New(internal.Module).Run()`.
- **Config**: typed via `pkg/config` (env vars only, no config file). See
  that package for how to add a new setting.
- **Shared infra lives in `pkg/`, not `internal/`**: `pkg/config`, `pkg/logger`
  and `pkg/database` aren't specific to the `api` service - anything a future
  service would also need goes there instead of under `services/api/internal`.
