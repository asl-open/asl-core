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
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/services/api/internal"
)

func main() {
	fx.New(
		internal.Module,
		fx.RecoverFromPanics(),
	).Run()
}
