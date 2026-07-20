package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/logger"
	"github.com/asl-open/asl-core/services/api/internal/http/handlers/health"
	apihttpmiddleware "github.com/asl-open/asl-core/services/api/internal/http/middleware"
	servicehealth "github.com/asl-open/asl-core/services/api/internal/services/health"
)

func newTestEngine(checker *servicehealth.MockChecker) *gin.Engine {
	gin.SetMode(gin.TestMode)

	mw := apihttpmiddleware.New(apihttpmiddleware.Params{Logger: &logger.MockLogger{}})

	engine := gin.New()
	engine.Use(mw.Errors())

	h := health.New(health.Params{Checker: checker, Logger: &logger.MockLogger{}})
	registerRoutes(engine, &Params{HealthHandler: h})

	return engine
}

func TestRegisterRoutes_Health(t *testing.T) {
	engine := newTestEngine(&servicehealth.MockChecker{})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/health", http.NoBody)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestRegisterRoutes_Ready(t *testing.T) {
	t.Run("dependency available", func(t *testing.T) {
		checker := &servicehealth.MockChecker{}
		checker.On("Ready", mock.Anything).Return(nil)

		engine := newTestEngine(checker)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/ready", http.NoBody)
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
	})

	t.Run("dependency unavailable", func(t *testing.T) {
		checker := &servicehealth.MockChecker{}
		checker.On("Ready", mock.Anything).Return(errors.New("connection refused"))

		engine := newTestEngine(checker)

		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/ready", http.NoBody)
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)
		require.JSONEq(t, `{"status":"unavailable"}`, rec.Body.String())
		require.NotContains(t, rec.Body.String(), "connection refused",
			"must not leak internal error details to the client")
	})
}

// recordingLogger wraps MockLogger's no-op Info/Debug/Warn with an Error
// that records what it was called with.
type recordingLogger struct {
	*logger.MockLogger
	errors []string
}

func (r *recordingLogger) Error(_ context.Context, msg string, _ ...any) {
	r.errors = append(r.errors, msg)
}

func TestShutdownHook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hook := shutdownHook(&logger.MockLogger{}, func(context.Context) error {
			return nil
		})

		require.NoError(t, hook(t.Context()))
	})

	t.Run("logs and wraps the error on failure", func(t *testing.T) {
		recLogger := &recordingLogger{MockLogger: &logger.MockLogger{}}
		shutdownErr := errors.New("connection reset")

		hook := shutdownHook(recLogger, func(context.Context) error {
			return shutdownErr
		})

		err := hook(t.Context())

		require.Error(t, err)
		require.ErrorIs(t, err, shutdownErr)
		require.NotEmpty(t, recLogger.errors, "a shutdown failure must be logged")
	})
}
