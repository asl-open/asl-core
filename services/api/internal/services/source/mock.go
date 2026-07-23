package source

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	sourcerepo "github.com/asl-open/asl-core/services/api/internal/repository/source"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, in *CreateInput) (sourcerepo.Source, error) {
	called := m.Called(ctx, in)
	out, _ := called.Get(0).(sourcerepo.Source)
	return out, called.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (sourcerepo.Source, error) {
	called := m.Called(ctx, id)
	out, _ := called.Get(0).(sourcerepo.Source)
	return out, called.Error(1)
}

func (m *MockService) List(ctx context.Context) ([]sourcerepo.Source, error) {
	called := m.Called(ctx)
	out, _ := called.Get(0).([]sourcerepo.Source)
	return out, called.Error(1)
}

var _ Service = (*MockService)(nil)
