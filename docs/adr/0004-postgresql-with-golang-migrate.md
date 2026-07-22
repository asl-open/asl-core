# 0004. PostgreSQL with golang-migrate up/down migrations

- **Status:** Accepted
- **Date:** 2026-07-22

## Context

ASL Core stores structured, versioned content with strong relational
integrity (sources, entries, references, revisions, reviews, releases).
This calls for a relational database with solid constraint support, and a
disciplined, reviewable way to evolve the schema over time. Schema changes
must be reproducible across local, CI and production environments, and
reversible so a bad migration can be rolled back.

## Decision

We will use **PostgreSQL**, with schema changes managed by
[**golang-migrate**](https://github.com/golang-migrate/migrate).

- Migrations live under `migrations/<service>/`, one directory per service
  that owns a database (currently `migrations/api/`).
- Each migration is a matched **up/down** pair:

  ```
  NNNNNN_description.up.sql
  NNNNNN_description.down.sql
  ```

  `NNNNNN` is a zero-padded sequence number, exactly what
  `migrate create -seq -digits 6` produces. Every `up` has a `down` that
  reverses it.
- Migrations are plain SQL, reviewed like any other code. UUID primary
  keys are generated with the `pgcrypto` extension (enabled by the first
  migration).
- The `migrate` CLI is driven through `make` targets
  (`migrate-create`, `migrate-up`, `migrate-down`, `migrate-version`);
  `DATABASE_DSN` configures the target database.

## Consequences

- Schema evolution is explicit, ordered, version-controlled and
  reviewable; the same migrations run in local, CI and production.
- Every change is reversible, which supports safe rollbacks and is
  required by the domain issues' acceptance criteria (`migrate up` then
  `down` must succeed on a clean database).
- Writing a correct `down` for every `up` is extra work and discipline,
  and destructive down-migrations can lose data — authors must consider
  this explicitly.
- Committing to PostgreSQL-specific features (extensions such as
  `pgcrypto`, and pgx driver behaviour) is acceptable; cross-database
  portability is a non-goal (see
  [ADR-0001](0001-language-and-framework-stack.md)).
