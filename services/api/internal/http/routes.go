package http

import "github.com/gin-gonic/gin"

// registerRoutes is the single place listing every HTTP route and the
// handler method it maps to.
func registerRoutes(engine *gin.Engine, p Params) {
	engine.GET("/health", p.HealthHandler.Health)
	engine.GET("/ready", p.HealthHandler.Ready)
}
