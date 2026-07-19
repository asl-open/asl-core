# Migrations

Each service that owns a database gets its own directory here
(`migrations/<service>/`), managed with the
[golang-migrate](https://github.com/golang-migrate/migrate) CLI.

## Naming convention

```
NNNNNN_description.up.sql
NNNNNN_description.down.sql
```

`NNNNNN` is a zero-padded sequence number (`000001`, `000002`, ...). Every
`up` migration has a matching `down` migration that reverses it. This is
exactly what `migrate create -seq -digits 6` produces, see below.

## Commands

Run from the repository root. `DATABASE_DSN` must be set for
`up`/`down`/`version` (see `services/api/.env.example`); `migrate-create`
only touches the filesystem, no database needed.

```
make migrate-create name=add_users_table  # scaffold a new up/down pair
make migrate-up                           # apply all pending migrations
make migrate-down                         # roll back the last migration
make migrate-version                      # show the current version
```

The Makefile installs the `migrate` CLI itself the first time it's needed
(`go install ... github.com/golang-migrate/migrate/v4/cmd/migrate`) - no
manual setup required. These targets currently only cover the `api`
service; the rest of the Makefile (run/build/test/lint/docker) lands with
#13.
