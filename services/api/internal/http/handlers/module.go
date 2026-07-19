package handlers

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/services/api/internal/http/handlers/health"
)

var Module = fx.Options(
	health.Module,
)
