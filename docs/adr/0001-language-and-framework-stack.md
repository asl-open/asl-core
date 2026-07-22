# 0001. Language and framework stack

- **Status:** Accepted
- **Date:** 2026-07-22

## Context

ASL Core is a backend platform for storing, verifying, versioning and
publishing structured Islamic content, exposed primarily through a public
HTTP API. It needs a stack that offers a strong standard library and
static typing, good HTTP performance, straightforward deployment as small
container images, and mature libraries for HTTP routing, dependency
injection and PostgreSQL access.

The team is proficient in Go, and Go's single-binary builds and small
runtime images fit the "independently runnable and deployable services"
goal. Within Go, several HTTP and application-composition approaches were
available (net/http with a router, Gin, Echo, chi; manual wiring, Wire,
Uber Fx).

## Decision

We will build ASL Core in **Go**, using:

- **Gin** (`github.com/gin-gonic/gin`) as the HTTP framework — routing,
  middleware and request handling.
- **Uber Fx** (`go.uber.org/fx`) for dependency injection, application
  composition and lifecycle management (start/stop hooks).
- **PostgreSQL** as the primary datastore, accessed via **pgx**
  (`github.com/jackc/pgx/v5`).
- **Zap** (`go.uber.org/zap`) for structured logging.

The specific decisions about how these are structured (monorepo layout,
package boundaries, migrations) are recorded in
[ADR-0002](0002-modular-monorepo-with-a-single-go-module.md),
[ADR-0003](0003-per-resource-packages-and-http-layer-boundary.md) and
[ADR-0004](0004-postgresql-with-golang-migrate.md).

## Consequences

- Small, dependency-free container images and fast startup, suited to the
  deployability goal.
- Fx gives explicit, testable composition and clean lifecycle hooks
  (used, for example, for graceful shutdown), at the cost of some
  indirection and a learning curve for contributors new to Fx.
- Gin keeps HTTP handling ergonomic, but its types must not leak into the
  domain — see [ADR-0003](0003-per-resource-packages-and-http-layer-boundary.md).
- Committing to PostgreSQL/pgx means database-specific features (e.g.
  `pgcrypto`) are fair game, but portability to other databases is not a
  goal.
- Contributors must be comfortable with Go and these libraries; versions
  are pinned in `go.mod` and upgraded deliberately.
