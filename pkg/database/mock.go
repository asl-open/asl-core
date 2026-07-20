package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

type MockConn struct {
	mock.Mock
}

func (m *MockConn) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	called := m.Called(ctx, sql, args)
	tag, _ := called.Get(0).(pgconn.CommandTag)
	return tag, called.Error(1)
}

func (m *MockConn) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	called := m.Called(ctx, sql, args)
	rows, _ := called.Get(0).(pgx.Rows)
	return rows, called.Error(1)
}

func (m *MockConn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	called := m.Called(ctx, sql, args)
	row, _ := called.Get(0).(pgx.Row)
	return row
}

func (m *MockConn) Ping(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

var _ Conn = (*MockConn)(nil)
