package health

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockChecker struct {
	mock.Mock
}

func (m *MockChecker) Ready(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

var _ Checker = (*MockChecker)(nil)
