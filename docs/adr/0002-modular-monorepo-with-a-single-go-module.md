# 0002. Modular monorepo with a single Go module

- **Status:** Accepted
- **Date:** 2026-07-22

## Context

ASL Core is expected to grow into several concerns (sources, content,
editorial, publishing, a public API) that may eventually become
independently deployable services. We needed to decide how to organise the
code: one repository or many, and one Go module or several.

A single module keeps dependency versions and tooling uniform and makes
cross-cutting refactors trivial, at the cost of coupling everything to one
`go.mod`. Multiple modules or repositories allow independent versioning but
add friction early, when the boundaries are still moving.

## Decision

We will use a **modular monorepo with a single root `go.mod`**
(`github.com/asl-open/asl-core`).

- Deployable services live under `services/<name>/` (currently
  `services/api`). Each service owns its `internal/` tree and is built and
  run independently (e.g. its own `cmd/` and Dockerfile).
- Code shared across services lives under `pkg/` (e.g. `pkg/config`,
  `pkg/logger`, `pkg/database`), not under any single service's
  `internal/`.
- Database migrations live under `migrations/<service>/`.
- One `go.mod`, one `Makefile` and one CI pipeline cover the whole repo.

Services are structured to be independently runnable and deployable now,
so that any of them can be split into its own module or repository later
without a rewrite.

## Consequences

- Uniform dependency versions and tooling; a single `go test ./...`,
  `make lint` and CI run cover everything.
- Cross-service refactoring and shared-code changes are atomic in one
  commit.
- Anything genuinely reusable must go in `pkg/`, not a service's
  `internal/`, or other services cannot import it; this is a discipline
  contributors must follow.
- A single module means all services share one dependency graph and are
  versioned together; independent module versioning is deferred until a
  service actually needs to split out.
