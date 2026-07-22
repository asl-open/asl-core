# 0003. Per-resource packages and the HTTP-layer boundary

- **Status:** Accepted
- **Date:** 2026-07-22

## Context

As domain resources are added (sources, knowledge entries, revisions,
reviews, releases), we need a package layout that keeps each resource
self-contained and keeps the domain independent of delivery and
infrastructure concerns. Two anti-patterns to avoid: a single shared
"repository" or "service" package for every table, and framework types
(Gin, Fx, pgx) leaking into the domain, which makes the domain hard to
test and hard to reuse.

## Decision

We will give **each resource its own packages**, split by concern, within a
service's `internal/` tree:

```
internal/repository/<resource>    domain entity + data access
internal/services/<resource>      business logic
internal/http/handlers/<resource> HTTP transport
```

For example: `internal/repository/source`, `internal/services/source`,
`internal/http/handlers/source`.

We will enforce these boundaries:

- **Gin types stay inside `internal/http`.** Domain entities and services
  never import Gin. Route registration is centralised in
  `internal/http/routes.go`; handlers only implement their request logic
  and depend on `internal/services/*`, never directly on `pkg/database` or
  other infrastructure.
- **Domain entities do not depend on Gin, Fx or pgx.** Business and
  dependency logic lives in the service layer.
- **One package per concern per resource** — no single repository or
  service package shared across all tables.
- **Gateways are created only for real external systems.** We do not add
  empty layers for structural symmetry; a `gateway` package appears only
  when a resource actually integrates with an external service.
- Every package wired into the app exposes its own Fx `Module`, aggregated
  one level up; `cmd/main.go` only calls `fx.New(internal.Module).Run()`.

## Consequences

- Each resource is understandable and testable in isolation; the domain
  can be unit-tested without HTTP or a database.
- The domain is portable and not coupled to the current frameworks, so
  Gin, Fx or pgx could be swapped at the edges without touching business
  logic.
- More packages and more boilerplate per resource (entity, service,
  handler, module) — accepted as the cost of clear boundaries.
- Contributors must respect the layering: handlers above services,
  services above repositories, no infrastructure imports in the domain.
  This is enforced by review and by the linters.
