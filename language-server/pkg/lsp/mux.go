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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/notedownorg/notedown/language-server/pkg/jsonrpc"
	"github.com/notedownorg/notedown/pkg/log"
)

type NotificationHandler func(params json.RawMessage) error
type MethodHandler func(params json.RawMessage) (any, error)

// formatRequestID formats a request ID for logging purposes
func formatRequestID(id *json.RawMessage) string {
	if id == nil {
		return "null"
	}
	return string(*id)
}

type Mux struct {
	reader               *bufio.Reader
	writer               *bufio.Writer
	notificationHandlers map[method]NotificationHandler
	methodHandlers       map[method]MethodHandler
	writeMutex           *sync.Mutex

	version string
	logger  *log.Logger
	server  LanguageServer
}

func NewMux(reader *bufio.Reader, writer *bufio.Writer, version string, logger *log.Logger) *Mux {
	return &Mux{
		reader:               reader,
		writer:               writer,
		notificationHandlers: make(map[method]NotificationHandler),
		methodHandlers:       make(map[method]MethodHandler),
		writeMutex:           &sync.Mutex{},
		version:              version,
		logger:               logger.WithScope("lsp/pkg/lsp"),
	}
}

func (m *Mux) RegisterNotification(method method, handler NotificationHandler) {
	m.notificationHandlers[method] = handler
}

func (m *Mux) RegisterMethod(method method, handler MethodHandler) {
	m.methodHandlers[method] = handler
}

func (m *Mux) SetServer(server LanguageServer) {
	m.server = server
}

func (m *Mux) write(response jsonrpc.Message) error {
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()
	return jsonrpc.Write(m.writer, response)
}

// PublishNotification sends a notification to the client
func (m *Mux) PublishNotification(method string, params any) error {
	notification := jsonrpc.NewNotification(method, params)
	return m.write(notification)
}

func (m *Mux) process() error {
	request, err := jsonrpc.Read(m.reader)
	if err != nil {
		m.logger.Error("failed to read JSON-RPC request", "error", err)
		return err
	}

	go func(request *jsonrpc.Request) {
		if request.IsNotification() {
			m.logger.Debug("processing notification", "method", request.Method)
			if handler, ok := m.notificationHandlers[method(request.Method)]; ok {
				if err := handler(request.Params); err != nil {
					m.logger.Error("notification handler failed", "method", request.Method, "error", err)
				}
			} else {
				m.logger.Warn("no handler for notification", "method", request.Method)
			}
		} else {
			m.logger.Debug("processing request", "method", request.Method, "id", formatRequestID(request.ID))
			handler, ok := m.methodHandlers[method(request.Method)]
			if !ok {
				m.logger.Warn("method not found", "method", request.Method, "id", formatRequestID(request.ID))
				if err := m.write(jsonrpc.NewMethodNotFoundError(request.ID, request.Method)); err != nil {
					m.logger.Error("failed to write method not found error", "error", err)
				}
				return
			}
			result, err := handler(request.Params)
			if err != nil {
				m.logger.Error("method handler failed", "method", request.Method, "id", formatRequestID(request.ID), "error", err)
				if writeErr := m.write(jsonrpc.NewInternalError(request.ID, err)); writeErr != nil {
					m.logger.Error("failed to write internal error", "error", writeErr)
				}
				return
			}
			m.logger.Debug("method completed successfully", "method", request.Method, "id", formatRequestID(request.ID))
			if err := m.write(jsonrpc.NewResponse(request.ID, result)); err != nil {
				m.logger.Error("failed to write response", "error", err)
			}
		}
	}(request)
	return nil
}

func (m *Mux) Run() error {
	if m.server == nil {
		m.logger.Error("no lsp server set")
		return fmt.Errorf("no lsp server set")
	}

	// Register initialize handler
	m.RegisterMethod(MethodInitialize, func(params json.RawMessage) (any, error) {
		var initParams InitializeParams
		if err := json.Unmarshal(params, &initParams); err != nil {
			m.logger.Error("failed to unmarshal initialize params", "error", err)
			return nil, err
		}

		result, err := m.server.Initialize(initParams)
		if err != nil {
			m.logger.Error("server initialization failed", "error", err)
			return nil, err
		}

		// Register all other handlers after successful initialization
		if err := m.server.RegisterHandlers(m); err != nil {
			m.logger.Error("failed to register server handlers", "error", err)
			return nil, err
		}

		return result, nil
	})

	m.logger.Info("starting lsp message processing loop")
	for {
		if err := m.process(); err != nil {
			m.logger.Error("lsp processing error", "error", err)
			return err
		}
	}
}

// SendRequest sends a request to the client and waits for a response
func (m *Mux) SendRequest(method string, params any) (any, error) {
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()

	// Create a unique request ID
	requestID := json.RawMessage(`1`) // Simple ID for now

	// Marshal params to JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		m.logger.Error("failed to marshal request params", "method", method, "error", err)
		return nil, err
	}

	// Create the request
	request := &jsonrpc.Request{
		ProtocolVersion: "2.0",
		Method:          method,
		Params:          paramsJSON,
		ID:              &requestID,
	}

	// Write the request
	if err := jsonrpc.Write(m.writer, request); err != nil {
		m.logger.Error("failed to send request", "method", method, "error", err)
		return nil, err
	}

	// For now, we'll return immediately as implementing full request/response
	// handling requires more complex state management. The client will receive
	// the request and should handle it appropriately.
	m.logger.Debug("sent request to client", "method", method)
	return nil, nil
}
