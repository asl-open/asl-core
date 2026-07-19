// Package services aggregates every business-logic service's fx module.
package services

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/services/api/internal/services/health"
)

var Module = fx.Options(
	health.Module,
)
