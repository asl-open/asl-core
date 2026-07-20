package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockMw struct {
	mock.Mock
}

func (m *MockMw) Logging() gin.HandlerFunc {
	args := m.Called()
	if fn, ok := args.Get(0).(gin.HandlerFunc); ok {
		return fn
	}
	return func(c *gin.Context) {}
}

func (m *MockMw) Errors() gin.HandlerFunc {
	args := m.Called()
	if fn, ok := args.Get(0).(gin.HandlerFunc); ok {
		return fn
	}
	return func(c *gin.Context) {}
}

func (m *MockMw) RequestID() gin.HandlerFunc {
	args := m.Called()
	if fn, ok := args.Get(0).(gin.HandlerFunc); ok {
		return fn
	}
	return func(c *gin.Context) {}
}

var _ Middleware = (*MockMw)(nil)
