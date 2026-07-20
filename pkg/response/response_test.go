package response_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/errorsx"
	"github.com/asl-open/asl-core/pkg/response"
)

func TestErr_AppError(t *testing.T) {
	appErr := errorsx.New().
		WithStatusCode(http.StatusBadRequest).
		WithCode(400).
		WithMessage("field x is required")

	resp := response.Err(appErr)

	require.Equal(t, http.StatusBadRequest, resp.HeaderCode)
	require.Equal(t, 400, resp.Code)
	require.Equal(t, "field x is required", resp.Message)
}

func TestErr_UnknownError(t *testing.T) {
	resp := response.Err(errors.New("dial tcp 10.0.0.5:5432: connection refused"))

	require.Equal(t, http.StatusInternalServerError, resp.HeaderCode)
	require.Equal(t, "Internal Server Error", resp.Message,
		"an unrecognized error's real content must never reach the response")
}

func TestInternal(t *testing.T) {
	resp := response.Internal()

	require.Equal(t, http.StatusInternalServerError, resp.HeaderCode)
	require.Equal(t, "Internal Server Error", resp.Message)
}

func TestResponse_WithRequestID(t *testing.T) {
	base := response.Internal()

	withID := base.WithRequestID("req-123")

	require.Empty(t, base.RequestID, "WithRequestID must not mutate the receiver")
	require.Equal(t, "req-123", withID.RequestID)
}
