// Package internal contains the api service's layered implementation.
package internal

import (
	"go.uber.org/fx"

	"github.com/asl-open/asl-core/pkg/config"
	"github.com/asl-open/asl-core/pkg/database"
	"github.com/asl-open/asl-core/pkg/logger"
	apihttp "github.com/asl-open/asl-core/services/api/internal/http"
	"github.com/asl-open/asl-core/services/api/internal/repository"
	"github.com/asl-open/asl-core/services/api/internal/services"
)

// Module aggregates all Fx modules for the api service.
var Module = fx.Module("api",
	config.Module,
	logger.Module,
	database.Module,
	repository.Module,
	services.Module,
	apihttp.Module,
)
