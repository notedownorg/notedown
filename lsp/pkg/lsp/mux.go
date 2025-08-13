package lsp

import (
	"bufio"
	"encoding/json"
	"sync"

	"github.com/notedownorg/notedown/lsp/pkg/jsonrpc"
	"github.com/notedownorg/notedown/pkg/log"
)

type NotificationHandler func(params json.RawMessage) error
type MethodHandler func(params json.RawMessage) (any, error)

type Mux struct {
	reader               *bufio.Reader
	writer               *bufio.Writer
	notificationHandlers map[method]NotificationHandler
	methodHandlers       map[method]MethodHandler
	writeMutex           *sync.Mutex

	version string
	logger  *log.Logger
}

func NewMux(reader *bufio.Reader, writer *bufio.Writer, version string, logger *log.Logger) *Mux {
	return &Mux{
		reader:               reader,
		writer:               writer,
		notificationHandlers: make(map[method]NotificationHandler),
		methodHandlers:       make(map[method]MethodHandler),
		writeMutex:           &sync.Mutex{},
		version:              version,
		logger:               logger,
	}
}

func (m *Mux) RegisterNotification(method method, handler NotificationHandler) {
	m.notificationHandlers[method] = handler
}

func (m *Mux) RegisterMethod(method method, handler MethodHandler) {
	m.methodHandlers[method] = handler
}

func (m *Mux) write(response jsonrpc.Message) error {
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()
	return jsonrpc.Write(m.writer, response)
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
			m.logger.Debug("processing request", "method", request.Method, "id", request.ID)
			handler, ok := m.methodHandlers[method(request.Method)]
			if !ok {
				m.logger.Warn("method not found", "method", request.Method, "id", request.ID)
				m.write(jsonrpc.NewMethodNotFoundError(request.ID, request.Method))
				return
			}
			result, err := handler(request.Params)
			if err != nil {
				m.logger.Error("method handler failed", "method", request.Method, "id", request.ID, "error", err)
				m.write(jsonrpc.NewInternalError(request.ID, err))
				return
			}
			m.logger.Debug("method completed successfully", "method", request.Method, "id", request.ID)
			m.write(jsonrpc.NewResponse(request.ID, result))
		}
	}(request)
	return nil
}
