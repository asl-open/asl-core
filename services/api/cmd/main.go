// Package main is the api service's entry point.
//
//	@title			ASL Core API
//	@description	Backend API for ASL Core. Operational endpoints
//	@description	(/health, /ready) are intentionally unversioned - they
//	@description	must stay stable for load balancers/orchestrators
//	@description	across API versions. Business/resource endpoints, once
//	@description	added, will live under /api/v1.
//	@version		0.1.0
//	@servers.url	/
//	@tag.name		health
//	@tag.description	Operational liveness/readiness endpoints. Unversioned by design.
package main

import (
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/services/api/internal"
)

func main() {
	fx.New(
		internal.Module,
		fx.RecoverFromPanics(),
		fx.StopTimeout(shutdownTimeout()),
	).Run()
}

// shutdownTimeout reads SHUTDOWN_TIMEOUT before the Fx app exists, since
// fx.StopTimeout is an Option evaluated at fx.New - it can't come from
// something the app's own DI graph provides. config.New is cheap and
// side-effect-free (just reads env vars), so the app's normal DI graph
// still constructs its own Config independently; this is a second,
// throwaway read of the same environment.
func shutdownTimeout() time.Duration {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	d, err := time.ParseDuration(cfg.GetString("shutdown.timeout"))
	if err != nil {
		log.Fatalf("invalid SHUTDOWN_TIMEOUT: %v", err)
	}

	return d
}
