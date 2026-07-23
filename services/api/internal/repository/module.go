package repository

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/services/api/internal/repository/contributor"
	"github.com/asl-open/asl-core/services/api/internal/repository/source"
)

// Module aggregates every resource repository's fx module.
var Module = fx.Options(
	source.Module,
	contributor.Module,
)
