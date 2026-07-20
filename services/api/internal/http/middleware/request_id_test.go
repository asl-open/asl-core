package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/logger"
	"github.com/asl-open/asl-core/pkg/requestid"
	"github.com/asl-open/asl-core/services/api/internal/http/middleware"
)

func newRequestIDEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	mw := middleware.New(middleware.Params{Logger: &logger.MockLogger{}})

	engine := gin.New()
	engine.Use(mw.RequestID())
	engine.GET("/", func(c *gin.Context) {
		id, ok := requestid.FromContext(c.Request.Context())
		require.True(t, ok, "RequestID must attach the ID to the request context")
		c.String(http.StatusOK, "%s", id)
	})

	return engine
}

func doRequest(t *testing.T, engine *gin.Engine, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", http.NoBody)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	return rec
}

func TestRequestID_PreservesExistingValidHeader(t *testing.T) {
	engine := newRequestIDEngine(t)

	rec := doRequest(t, engine, map[string]string{"X-Request-Id": "upstream-correlation-id-42"})

	require.Equal(t, "upstream-correlation-id-42", rec.Header().Get("X-Request-Id"))
	require.Equal(t, "upstream-correlation-id-42", rec.Body.String())
}

func TestRequestID_GeneratesWhenMissing(t *testing.T) {
	engine := newRequestIDEngine(t)

	rec := doRequest(t, engine, nil)

	require.NotEmpty(t, rec.Header().Get("X-Request-Id"))
	require.Equal(t, rec.Header().Get("X-Request-Id"), rec.Body.String())
}

func TestRequestID_GeneratedIDsAreUnique(t *testing.T) {
	engine := newRequestIDEngine(t)

	first := doRequest(t, engine, nil).Header().Get("X-Request-Id")
	second := doRequest(t, engine, nil).Header().Get("X-Request-Id")

	require.NotEqual(t, first, second)
}

func TestRequestID_RejectsUnsafeHeaderValues(t *testing.T) {
	cases := map[string]string{
		"empty":                    "",
		"contains a newline":       "abc\ndef",
		"contains a control byte":  "abc\x7fdef",
		"exceeds the length limit": strings.Repeat("a", 129),
	}

	for name, value := range cases {
		t.Run(name, func(t *testing.T) {
			engine := newRequestIDEngine(t)

			rec := doRequest(t, engine, map[string]string{"X-Request-Id": value})

			got := rec.Header().Get("X-Request-Id")
			require.NotEqual(t, value, got, "an unsafe header value must not be echoed back as-is")
			require.NotEmpty(t, got)
		})
	}
}
