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

package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/notedownorg/notedown/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMuxSetServer(t *testing.T) {
	tests := []struct {
		name   string
		server LanguageServer
	}{
		{
			name:   "set valid mock server",
			server: &MockServer{},
		},
		{
			name:   "set nil server",
			server: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(bytes.NewReader([]byte{}))
			writer := bufio.NewWriter(bytes.NewBuffer([]byte{}))
			logger := log.NewLsp(log.Info, log.FormatText)

			mux := NewMux(reader, writer, "test", logger)
			mux.SetServer(tt.server)

			assert.Equal(t, tt.server, mux.server)
		})
	}
}

func TestMuxRunWithoutServer(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte{}))
	writer := bufio.NewWriter(bytes.NewBuffer([]byte{}))
	logger := log.NewLsp(log.Info, log.FormatText)

	mux := NewMux(reader, writer, "test", logger)
	// Don't set a server

	err := mux.Run()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no lsp server set")
}

func TestMuxInitializeHandling(t *testing.T) {
	tests := []struct {
		name           string
		initParams     InitializeParams
		mockResult     InitializeResult
		mockError      error
		expectError    bool
		expectResponse bool
	}{
		{
			name: "successful initialization",
			initParams: InitializeParams{
				ClientInfo:   &ClientInfo{Name: "test-client"},
				Capabilities: ClientCapabilities{},
			},
			mockResult: InitializeResult{
				ServerInfo:   &ServerInfo{Name: "Test Server", Version: "1.0.0"},
				Capabilities: ServerCapabilities{},
			},
			mockError:      nil,
			expectError:    false,
			expectResponse: true,
		},
		{
			name: "initialization with nil client info",
			initParams: InitializeParams{
				ClientInfo:   nil,
				Capabilities: ClientCapabilities{},
			},
			mockResult: InitializeResult{
				ServerInfo:   &ServerInfo{Name: "Test Server", Version: "1.0.0"},
				Capabilities: ServerCapabilities{},
			},
			mockError:      nil,
			expectError:    false,
			expectResponse: true,
		},
		{
			name: "server initialization fails",
			initParams: InitializeParams{
				ClientInfo:   &ClientInfo{Name: "test-client"},
				Capabilities: ClientCapabilities{},
			},
			mockResult:     InitializeResult{},
			mockError:      errors.New("initialization failed"),
			expectError:    true,
			expectResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := &MockServer{}
			mockServer.On("Initialize", tt.initParams).Return(tt.mockResult, tt.mockError)

			if !tt.expectError {
				// If initialization succeeds, RegisterHandlers should be called
				mockServer.On("RegisterHandlers", mock.AnythingOfType("*lsp.Mux")).Return(nil)
			}

			// Set up mux
			reader := bufio.NewReader(bytes.NewReader([]byte{}))
			var responseBuffer bytes.Buffer
			writer := bufio.NewWriter(&responseBuffer)
			logger := log.NewLsp(log.Info, log.FormatText)

			mux := NewMux(reader, writer, "test", logger)
			mux.SetServer(mockServer)

			// Test the initialize handler directly
			paramsJSON, err := json.Marshal(tt.initParams)
			require.NoError(t, err)

			// Get the initialize handler that would be registered
			handler := mux.methodHandlers[MethodInitialize]
			if handler == nil {
				// Handler isn't registered yet, need to trigger Run to register it
				// For this test, we'll register it manually
				mux.RegisterMethod(MethodInitialize, func(params json.RawMessage) (any, error) {
					var initParams InitializeParams
					if err := json.Unmarshal(params, &initParams); err != nil {
						return nil, err
					}

					result, err := mux.server.Initialize(initParams)
					if err != nil {
						return nil, err
					}

					// Register all other handlers after successful initialization
					if err := mux.server.RegisterHandlers(mux); err != nil {
						return nil, err
					}

					return result, nil
				})
				handler = mux.methodHandlers[MethodInitialize]
			}

			// Call the handler
			result, err := handler(json.RawMessage(paramsJSON))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Verify the result matches expected
				initResult, ok := result.(InitializeResult)
				assert.True(t, ok, "Expected InitializeResult type")
				assert.Equal(t, tt.mockResult, initResult)
			}

			// Verify mock expectations
			mockServer.AssertExpectations(t)
		})
	}
}

func TestMuxHandlerRegistration(t *testing.T) {
	tests := []struct {
		name                string
		registerHandlersErr error
		expectError         bool
	}{
		{
			name:                "successful handler registration",
			registerHandlersErr: nil,
			expectError:         false,
		},
		{
			name:                "handler registration fails",
			registerHandlersErr: errors.New("failed to register handlers"),
			expectError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := &MockServer{}
			initParams := InitializeParams{
				ClientInfo:   &ClientInfo{Name: "test-client"},
				Capabilities: ClientCapabilities{},
			}
			initResult := InitializeResult{
				ServerInfo:   &ServerInfo{Name: "Test Server", Version: "1.0.0"},
				Capabilities: ServerCapabilities{},
			}

			mockServer.On("Initialize", initParams).Return(initResult, nil)
			mockServer.On("RegisterHandlers", mock.AnythingOfType("*lsp.Mux")).Return(tt.registerHandlersErr)

			// Set up mux
			reader := bufio.NewReader(bytes.NewReader([]byte{}))
			var responseBuffer bytes.Buffer
			writer := bufio.NewWriter(&responseBuffer)
			logger := log.NewLsp(log.Info, log.FormatText)

			mux := NewMux(reader, writer, "test", logger)
			mux.SetServer(mockServer)

			// Register the initialize handler manually for testing
			mux.RegisterMethod(MethodInitialize, func(params json.RawMessage) (any, error) {
				var initParams InitializeParams
				if err := json.Unmarshal(params, &initParams); err != nil {
					return nil, err
				}

				result, err := mux.server.Initialize(initParams)
				if err != nil {
					return nil, err
				}

				// Register all other handlers after successful initialization
				if err := mux.server.RegisterHandlers(mux); err != nil {
					return nil, err
				}

				return result, nil
			})

			// Test the initialize handler
			paramsJSON, err := json.Marshal(initParams)
			require.NoError(t, err)

			handler := mux.methodHandlers[MethodInitialize]
			result, err := handler(json.RawMessage(paramsJSON))

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			// Verify mock expectations
			mockServer.AssertExpectations(t)
		})
	}
}

func TestMuxMethodAndNotificationRegistration(t *testing.T) {
	reader := bufio.NewReader(bytes.NewReader([]byte{}))
	writer := bufio.NewWriter(bytes.NewBuffer([]byte{}))
	logger := log.NewLsp(log.Info, log.FormatText)

	mux := NewMux(reader, writer, "test", logger)

	// Test method registration
	methodCalled := false
	testMethod := method("test/method")
	testHandler := func(params json.RawMessage) (any, error) {
		methodCalled = true
		return "test result", nil
	}

	mux.RegisterMethod(testMethod, testHandler)

	// Verify handler was registered
	assert.Contains(t, mux.methodHandlers, testMethod)

	// Call the handler
	result, err := mux.methodHandlers[testMethod](json.RawMessage(`{}`))
	assert.NoError(t, err)
	assert.Equal(t, "test result", result)
	assert.True(t, methodCalled)

	// Test notification registration
	notificationCalled := false
	testNotification := method("test/notification")
	testNotificationHandler := func(params json.RawMessage) error {
		notificationCalled = true
		return nil
	}

	mux.RegisterNotification(testNotification, testNotificationHandler)

	// Verify handler was registered
	assert.Contains(t, mux.notificationHandlers, testNotification)

	// Call the handler
	err = mux.notificationHandlers[testNotification](json.RawMessage(`{}`))
	assert.NoError(t, err)
	assert.True(t, notificationCalled)
}
