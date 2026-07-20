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

## Adding a new endpoint

Add the `@Summary`/`@Description`/`@Tags`/`@Produce`/`@Success`/
`@Failure`/`@Router`/`@ID` annotations directly above the handler method
(see `services/api/internal/http/handlers/health/health.go` and
`ready.go` for the pattern), then run `make openapi`. If the handler's
directory isn't already covered, add it to `SWAG_DIRS` in the root
`Makefile`.

## Versioning convention

Business/resource endpoints, once they exist, will be mounted under
`/api/v1`. Bump to `/api/v2` etc. only for breaking changes to an
existing resource's contract, not for additions.

`/health` and `/ready` are the one deliberate exception: they stay
unversioned at the root, because they're operational endpoints consumed
by load balancers and orchestrators (e.g. Kubernetes probes), which need
a stable path across API versions - versioning them would break every
external health check on every API version bump.
