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
go run ./services/api/cmd/migrate create <name>  # scaffold a new up/down pair
go run ./services/api/cmd/migrate up             # apply all pending migrations
go run ./services/api/cmd/migrate down           # roll back the last migration
go run ./services/api/cmd/migrate version        # show the current version
```

`create` only touches the filesystem (`MIGRATION_SOURCE`), no database
needed. `up`/`down`/`version` need `DATABASE_DSN` set (see
`services/api/.env.example`). These will get wrapped into
`make migrate-up`/`make migrate-down`/etc. once the Makefile lands (#13).

## Adding a migration

```
go run ./services/api/cmd/migrate create add_users_table
```

This creates `NNNNNN_add_users_table.up.sql`/`.down.sql` in the service's
migration directory, continuing its existing sequence. Fill in the SQL by
hand.
