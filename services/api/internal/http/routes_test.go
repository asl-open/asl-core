package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/logger"
	"github.com/asl-open/asl-core/services/api/internal/http/handlers/health"
	servicehealth "github.com/asl-open/asl-core/services/api/internal/services/health"
)

func newTestEngine(checker *servicehealth.MockChecker) *gin.Engine {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	h := health.New(health.Params{Checker: checker, Logger: &logger.MockLogger{}})
	registerRoutes(engine, Params{HealthHandler: h})

	return engine
}

func TestRegisterRoutes_Health(t *testing.T) {
	engine := newTestEngine(&servicehealth.MockChecker{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
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

		req := httptest.NewRequest(http.MethodGet, "/ready", nil)
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
	})

	t.Run("dependency unavailable", func(t *testing.T) {
		checker := &servicehealth.MockChecker{}
		checker.On("Ready", mock.Anything).Return(errors.New("connection refused"))

		engine := newTestEngine(checker)

		req := httptest.NewRequest(http.MethodGet, "/ready", nil)
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)

		require.Equal(t, http.StatusServiceUnavailable, rec.Code)
		require.JSONEq(t, `{"status":"unavailable"}`, rec.Body.String())
		require.NotContains(t, rec.Body.String(), "connection refused",
			"must not leak internal error details to the client")
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}
