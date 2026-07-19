package logger

import "context"

type MockLogger struct {
}

func (m *MockLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
}

func (m *MockLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
}

func (m *MockLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
}

func (m *MockLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
}

var _ Logger = (*MockLogger)(nil)
