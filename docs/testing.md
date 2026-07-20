# Testing

## Running tests

```
make test        # go test ./...
make test-race    # go test -race ./...
```

Both run in CI (`.github/workflows/ci.yml`, `test` job); `test-race` is a
separate command by design - it's slower, so it isn't run on every push,
but should be run locally before merging anything touching concurrency
(goroutines, shared state, `fx` lifecycle hooks).

Tests under `pkg/database` need a real, reachable PostgreSQL instance
(`asl_core_test` database, port 5439 by default - see
`pkg/database/database_test.go`, overridable via `TEST_DATABASE_DSN`).
Locally, run one yourself; in CI, the `test` job starts one via
`.github/actions/postgres`.

## Layout

Test files are colocated with the code they test:
`<file>_test.go` next to `<file>.go`, same directory. No separate `test/`
tree.

## Package: external vs internal test packages

Default to an external test package (`package foo_test`, importing `foo`
like any other consumer). This exercises only the exported API and keeps
tests honest about what's actually usable from outside the package.

Use an internal test package (`package foo`) only when a test needs to
call an unexported identifier directly. Two examples in this repo:

- `pkg/logger/logger_test.go` (`package logger`) tests the unexported
  `isIgnorableSyncError` and `getLevel` functions directly.
- `services/api/internal/http/routes_test.go` (`package http`) tests the
  unexported `registerRoutes` function directly.

Every other test file in the repo (`pkg/config`, `pkg/database`,
`services/api/internal/http/module_test.go`, ...) uses the external
form. If a test doesn't need unexported access, it belongs in the
external package.

## Mocks

Each package that defines an interface consumed elsewhere ships a
hand-written `MockXxx` next to it, in a plain `mock.go` (not `_test.go`,
so other packages' tests can import it - e.g. `logger.MockLogger` is used
from `services/api/internal/http`'s tests). See `pkg/logger/mock.go`,
`pkg/database/mock.go`, `services/api/internal/services/health/mock.go`,
`services/api/internal/http/handlers/health/mock.go`,
`services/api/internal/http/middleware/mock.go`.

These are pure delegation to `testify/mock` (`m.Called(...)`) - no
business logic, no conditionals, nothing that could itself have a bug
worth testing. Keep it that way: if a mock needs a real branch of logic
to be useful, that's a sign the test wants a fake or a real dependency
instead, not a smarter mock.

Beyond mocks, add a shared test helper only when the same setup code
would otherwise be duplicated across multiple test files - see
`newTestEngine`/`newTestMiddleware` in the `http` package's tests for an
example scoped to that package. There's no repo-wide test-helpers
package, and none should be added speculatively; introduce one only when
a second, unrelated package actually needs the same helper.

## Naming

Standard Go: `TestXxx` for top-level tests, `t.Run("description", ...)`
for subtests. Where `Xxx` would otherwise start with a lowercase letter
(testing an unexported identifier from an internal test package), prefix
with an underscore - `Test_getLevel`, `Test_isIgnorableSyncError` - so
the test name doesn't read as an unexported function itself.
