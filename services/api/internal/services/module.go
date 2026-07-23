// Package services aggregates every business-logic service's fx module.
package services

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/services/api/internal/services/contributor"
	"github.com/asl-open/asl-core/services/api/internal/services/health"
	"github.com/asl-open/asl-core/services/api/internal/services/source"
)

var Module = fx.Options(
	health.Module,
	source.Module,
	contributor.Module,
)
