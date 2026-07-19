# Migrations

Each service that owns a database gets its own directory here
(`migrations/<service>/`), managed with
[golang-migrate](https://github.com/golang-migrate/migrate).

## Naming convention

```
NNNNNN_description.up.sql
NNNNNN_description.down.sql
```

`NNNNNN` is a zero-padded sequence number (`000001`, `000002`, ...). Every
`up` migration has a matching `down` migration that reverses it.

## Commands

There's no `migrate` CLI dependency - each service has its own
`cmd/migrate` binary. For the `api` service, run from the repository root:

```
go run ./services/api/cmd/migrate up       # apply all pending migrations
go run ./services/api/cmd/migrate down     # roll back the last migration
go run ./services/api/cmd/migrate version  # show the current version
```

`DATABASE_DSN` must be set (see `services/api/.env.example`). These will
get wrapped into `make migrate-up`/`make migrate-down` once the Makefile
lands (#13).

## Adding a migration

Create a new `NNNNNN_description.up.sql`/`.down.sql` pair in the
service's migration directory, continuing the sequence.
