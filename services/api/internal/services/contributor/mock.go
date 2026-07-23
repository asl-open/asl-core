package contributor

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	contributorrepo "github.com/asl-open/asl-core/services/api/internal/repository/contributor"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, in *CreateInput) (contributorrepo.Contributor, error) {
	called := m.Called(ctx, in)
	out, _ := called.Get(0).(contributorrepo.Contributor)
	return out, called.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uuid.UUID) (contributorrepo.Contributor, error) {
	called := m.Called(ctx, id)
	out, _ := called.Get(0).(contributorrepo.Contributor)
	return out, called.Error(1)
}

func (m *MockService) List(ctx context.Context) ([]contributorrepo.Contributor, error) {
	called := m.Called(ctx)
	out, _ := called.Get(0).([]contributorrepo.Contributor)
	return out, called.Error(1)
}

var _ Service = (*MockService)(nil)
