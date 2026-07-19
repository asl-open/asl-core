// Package internal contains the api service's layered implementation.
package internal

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/logger"
	apihttp "github.com/asl-open/asl-core/services/api/internal/http"
)

// Module aggregates all Fx modules for the api service.
var Module = fx.Module("api",
	config.Module,
	logger.Module,
	apihttp.Module,
)
