package handlers

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/services/api/internal/http/handlers/ping"
)

var Module = fx.Options(
	ping.Module,
)
