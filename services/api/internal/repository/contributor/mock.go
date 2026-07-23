package contributor

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, c *Contributor) (Contributor, error) {
	called := m.Called(ctx, c)
	out, _ := called.Get(0).(Contributor)
	return out, called.Error(1)
}

func (m *MockRepo) Get(ctx context.Context, id uuid.UUID) (Contributor, error) {
	called := m.Called(ctx, id)
	out, _ := called.Get(0).(Contributor)
	return out, called.Error(1)
}

func (m *MockRepo) List(ctx context.Context) ([]Contributor, error) {
	called := m.Called(ctx)
	out, _ := called.Get(0).([]Contributor)
	return out, called.Error(1)
}

var _ Repo = (*MockRepo)(nil)
