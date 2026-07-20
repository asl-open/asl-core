package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/logger"
	"github.com/asl-open/asl-core/services/api/internal/http/apierrors"
	"github.com/asl-open/asl-core/services/api/internal/http/middleware"
)

func newTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	mw := middleware.New(middleware.Params{Logger: &logger.MockLogger{}})

	engine := gin.New()
	engine.Use(mw.Errors())

	engine.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})
	engine.GET("/bad-request", func(c *gin.Context) {
		_ = c.Error(apierrors.ErrBadRequest.WithMessage("field x is required"))
	})
	engine.GET("/not-found", func(c *gin.Context) {
		_ = c.Error(apierrors.ErrNotFound.WithMessage("thing 42 not found"))
	})
	engine.GET("/unexpected", func(c *gin.Context) {
		_ = c.Error(errors.New("dial tcp 10.0.0.5:5432: connection refused"))
	})

	return engine
}

func do(t *testing.T, engine *gin.Engine, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, path, http.NoBody)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

func TestErrors_Panic(t *testing.T) {
	engine := newTestEngine(t)

	rec := do(t, engine, "/panic", nil)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"message":"Internal Server Error","code":500}`, rec.Body.String())
	require.NotContains(t, rec.Body.String(), "boom",
		"must not leak the panic value to the client")
}

func TestErrors_BadRequest(t *testing.T) {
	engine := newTestEngine(t)

	rec := do(t, engine, "/bad-request", nil)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.JSONEq(t, `{"message":"field x is required","code":400}`, rec.Body.String())
}

func TestErrors_NotFound(t *testing.T) {
	engine := newTestEngine(t)

	rec := do(t, engine, "/not-found", nil)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.JSONEq(t, `{"message":"thing 42 not found","code":404}`, rec.Body.String())
}

func TestErrors_UnexpectedError(t *testing.T) {
	engine := newTestEngine(t)

	rec := do(t, engine, "/unexpected", nil)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.JSONEq(t, `{"message":"Internal Server Error","code":500}`, rec.Body.String())
	require.NotContains(t, rec.Body.String(), "10.0.0.5",
		"must not leak internal error details (host, driver, etc.) to the client")
}

func TestErrors_RequestID(t *testing.T) {
	t.Run("echoes the X-Request-Id header when present", func(t *testing.T) {
		engine := newTestEngine(t)

		rec := do(t, engine, "/not-found", map[string]string{"X-Request-Id": "req-123"})

		require.JSONEq(t, `{"message":"thing 42 not found","code":404,"request_id":"req-123"}`, rec.Body.String())
	})

	t.Run("omits request_id when not available", func(t *testing.T) {
		engine := newTestEngine(t)

		rec := do(t, engine, "/not-found", nil)

		require.JSONEq(t, `{"message":"thing 42 not found","code":404}`, rec.Body.String())
	})
}
