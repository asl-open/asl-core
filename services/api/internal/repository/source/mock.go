package source

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, s *Source) (Source, error) {
	called := m.Called(ctx, s)
	out, _ := called.Get(0).(Source)
	return out, called.Error(1)
}

func (m *MockRepo) Get(ctx context.Context, id uuid.UUID) (Source, error) {
	called := m.Called(ctx, id)
	out, _ := called.Get(0).(Source)
	return out, called.Error(1)
}

func (m *MockRepo) List(ctx context.Context) ([]Source, error) {
	called := m.Called(ctx)
	out, _ := called.Get(0).([]Source)
	return out, called.Error(1)
}

var _ Repo = (*MockRepo)(nil)
