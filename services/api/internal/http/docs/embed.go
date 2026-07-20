// Package docs embeds the generated OpenAPI spec and its RapiDoc viewer
// page so both can be served at runtime without depending on the
// filesystem.
package docs

import _ "embed"

//go:embed swagger.yaml
var OpenAPISpec []byte

//go:embed rapidoc.html
var RapiDocPage []byte
