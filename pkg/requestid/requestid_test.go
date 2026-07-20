package requestid_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asl-open/asl-core/pkg/requestid"
)

func TestNewContext_FromContext(t *testing.T) {
	ctx := requestid.NewContext(context.Background(), "req-123")

	id, ok := requestid.FromContext(ctx)

	require.True(t, ok)
	require.Equal(t, "req-123", id)
}

func TestFromContext_NotSet(t *testing.T) {
	id, ok := requestid.FromContext(context.Background())

	require.False(t, ok)
	require.Empty(t, id)
}
