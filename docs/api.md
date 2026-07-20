# API documentation

## OpenAPI spec

`services/api/internal/http/docs/swagger.yaml` is generated from `swag`
annotations on the handlers (`@Summary`, `@Router`, ...) and the general
API info block in `services/api/cmd/main.go` - it isn't hand-written, and
shouldn't be edited directly. Regenerate it after changing a handler's
annotations or adding a new route:

```
make openapi          # regenerate services/api/internal/http/docs/swagger.yaml
make openapi-check    # fail if the committed spec is out of date, without changing it
make openapi-validate # lint the spec (structure, missing descriptions, etc.)
```

The spec is OpenAPI 3.1.

## Browsing it locally

The running app serves the spec itself, embedded via `go:embed` (no
filesystem dependency, works the same in the Docker image):

- `GET /openapi.yaml` - the raw spec
- `GET /docs` - a [RapiDoc](https://rapidocweb.com/) page rendering it
  (dark theme, try-it enabled, loaded from a CDN script)

`make docker-up` (or `go run ./cmd`/`make run`), then open
`http://localhost:8080/docs` (or whatever `HTTP_ADDR`/`API_PORT` you're
using).

## Adding a new endpoint

Add the `@Summary`/`@Description`/`@Tags`/`@Produce`/`@Success`/
`@Failure`/`@Router`/`@ID` annotations directly above the handler method
(see `services/api/internal/http/handlers/health/health.go` and
`ready.go` for the pattern), then run `make openapi`. If the handler's
directory isn't already covered, add it to `SWAG_DIRS` in the root
`Makefile`.

## Error responses

Every HTTP error uses the same JSON envelope
(`{"message": "...", "code": ..., "request_id": "..."}`, see
`pkg/response.Response`). To report one from a handler:

```go
c.Error(apierrors.ErrNotFound.WithMessage("user 42 not found"))
```

`apierrors` (`services/api/internal/http/apierrors`) holds the canonical
errors (`ErrBadRequest` -> 400, `ErrNotFound` -> 404); `.With...` returns
a customized copy without mutating the shared value. `middleware.Errors()`
(registered ahead of every route) turns whatever was passed to `c.Error`
into the response, and separately recovers panics into a generic 500 -
neither ever sends the real error text to the client for anything that
isn't a deliberately constructed `*errorsx.AppError`. See
`pkg/errorsx`, `pkg/response`, and
`services/api/internal/http/middleware/errors.go` for the full mapping.

## Versioning convention

Business/resource endpoints, once they exist, will be mounted under
`/api/v1`. Bump to `/api/v2` etc. only for breaking changes to an
existing resource's contract, not for additions.

`/health` and `/ready` are the one deliberate exception: they stay
unversioned at the root, because they're operational endpoints consumed
by load balancers and orchestrators (e.g. Kubernetes probes), which need
a stable path across API versions - versioning them would break every
external health check on every API version bump.
