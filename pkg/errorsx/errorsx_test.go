package errorsx_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/errorsx"
)

func TestAppError_Builders(t *testing.T) {
	err := errorsx.New().
		WithStatusCode(http.StatusNotFound).
		WithCode(404).
		WithMessage("not found")

	require.Equal(t, http.StatusNotFound, err.StatusCode())
	require.Equal(t, 404, err.Code())
	require.Equal(t, "not found", err.Message())
	require.False(t, err.Internal())
}

func TestAppError_WithInternal(t *testing.T) {
	err := errorsx.New().WithInternal()

	require.True(t, err.Internal())
}

func TestAppError_BuildersDoNotMutateReceiver(t *testing.T) {
	base := errorsx.New().WithMessage("base")

	_ = base.WithMessage("customized for this call site")

	require.Equal(t, "base", base.Message(),
		"With... must return a copy - a shared canonical error must not be mutated by one caller's customization")
}

func TestAppError_Error(t *testing.T) {
	err := errorsx.New().WithMessage("boom").WithCode(500)

	require.Contains(t, err.Error(), "boom")
}

func TestAs(t *testing.T) {
	t.Run("matches an AppError", func(t *testing.T) {
		appErr := errorsx.New().WithMessage("boom")

		got, ok := errorsx.As(appErr)

		require.True(t, ok)
		require.Same(t, appErr, got)
	})

	t.Run("matches an AppError wrapped with %w", func(t *testing.T) {
		appErr := errorsx.New().WithMessage("boom")
		wrapped := fmt.Errorf("context: %w", appErr)

		got, ok := errorsx.As(wrapped)

		require.True(t, ok)
		require.Same(t, appErr, got)
	})

	t.Run("does not match string-concatenated text", func(t *testing.T) {
		appErr := errorsx.New().WithMessage("boom")
		notWrapped := errors.New("context: " + appErr.Error())

		_, ok := errorsx.As(notWrapped)

		require.False(t, ok, "string concatenation is not errors.Is/As-compatible wrapping")
	})

	t.Run("rejects a plain error", func(t *testing.T) {
		got, ok := errorsx.As(errors.New("plain"))

		require.False(t, ok)
		require.Nil(t, got)
	})
}
