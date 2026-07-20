package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/asl-open/asl-core/services/api/internal/http/docs"
)

// registerRoutes is the single place listing every HTTP route and the
// handler method it maps to.
func registerRoutes(engine *gin.Engine, p *Params) {
	engine.GET("/health", p.HealthHandler.Health)
	engine.GET("/ready", p.HealthHandler.Ready)

	engine.GET("/openapi.yaml", serveOpenAPISpec)
	engine.GET("/docs", serveAPIDocs)
}

// serveOpenAPISpec serves the generated OpenAPI spec (see
// services/api/internal/http/docs, regenerated via `make openapi`).
func serveOpenAPISpec(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml", docs.OpenAPISpec)
}

// serveAPIDocs serves a RapiDoc page rendering the spec from
// /openapi.yaml, for browsing the API locally.
func serveAPIDocs(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", docs.RapiDocPage)
}
