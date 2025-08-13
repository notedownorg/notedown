package lsp

import "github.com/stretchr/testify/mock"

// MockServer is a mock implementation of the Server interface for testing
type MockServer struct {
	mock.Mock
}

func (m *MockServer) Initialize(params InitializeParams) (InitializeResult, error) {
	args := m.Called(params)
	return args.Get(0).(InitializeResult), args.Error(1)
}

func (m *MockServer) RegisterHandlers(mux *Mux) error {
	args := m.Called(mux)
	return args.Error(0)
}

func (m *MockServer) Shutdown() error {
	args := m.Called()
	return args.Error(0)
}
