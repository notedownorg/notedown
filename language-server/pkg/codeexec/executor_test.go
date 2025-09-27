// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codeexec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutorFactory_NewExecutorFactory(t *testing.T) {
	factory := NewExecutorFactory()
	require.NotNil(t, factory)
	assert.NotNil(t, factory.executors)
}

func TestExecutorFactory_GetExecutor_Go(t *testing.T) {
	factory := NewExecutorFactory()

	executor, err := factory.GetExecutor("go")
	if err != nil {
		// Go might not be available on the test system
		assert.Contains(t, err.Error(), "not available")
		return
	}

	require.NotNil(t, executor)
	assert.Equal(t, "go", executor.Language())
}

func TestExecutorFactory_GetExecutor_Unsupported(t *testing.T) {
	factory := NewExecutorFactory()

	executor, err := factory.GetExecutor("python")
	assert.Nil(t, executor)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no executor available for language: python")
}

func TestExecutorFactory_GetSupportedLanguages(t *testing.T) {
	factory := NewExecutorFactory()

	languages := factory.GetSupportedLanguages()
	assert.NotNil(t, languages)

	// Check that the list contains Go if it's available
	goExecutor := NewGoExecutor()
	if goExecutor.IsAvailable() {
		assert.Contains(t, languages, "go")
	}
}

func TestExecutorFactory_Register(t *testing.T) {
	factory := NewExecutorFactory()

	// Create a mock executor
	mockExecutor := &MockExecutor{
		language:  "test",
		available: true,
	}

	factory.Register(mockExecutor)

	executor, err := factory.GetExecutor("test")
	require.NoError(t, err)
	assert.Equal(t, mockExecutor, executor)
}

// MockExecutor is a simple mock executor for testing
type MockExecutor struct {
	language  string
	available bool
}

func (m *MockExecutor) Execute(blocks []CodeBlock, workspaceRoot string) (*ExecutionResult, error) {
	return &ExecutionResult{
		Success:  true,
		ExitCode: 0,
		Stdout:   "mock output",
	}, nil
}

func (m *MockExecutor) Language() string {
	return m.language
}

func (m *MockExecutor) IsAvailable() bool {
	return m.available
}

func TestExecutorFactory_GetExecutor_NotAvailable(t *testing.T) {
	factory := NewExecutorFactory()

	// Create a mock executor that's not available
	mockExecutor := &MockExecutor{
		language:  "unavailable",
		available: false,
	}

	factory.Register(mockExecutor)

	executor, err := factory.GetExecutor("unavailable")
	assert.Nil(t, executor)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not available (missing runtime)")
}
