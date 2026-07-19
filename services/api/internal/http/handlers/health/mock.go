package health

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Health(c *gin.Context) {
	m.Called(c)
}

func (m *MockHandler) Ready(c *gin.Context) {
	m.Called(c)
}

var _ Handler = (*MockHandler)(nil)

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}
