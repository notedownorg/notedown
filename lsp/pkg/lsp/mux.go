package lsp

import (
	"bufio"
	"encoding/json"
	"sync"

	"github.com/notedownorg/notedown/lsp/pkg/jsonrpc"
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
}

func NewMux(reader *bufio.Reader, writer *bufio.Writer, version string) *Mux {
	return &Mux{
		reader:               reader,
		writer:               writer,
		notificationHandlers: make(map[method]NotificationHandler),
		methodHandlers:       make(map[method]MethodHandler),
		writeMutex:           &sync.Mutex{},
		version:              version,
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
		return err
	}
	go func(request *jsonrpc.Request) {
		if request.IsNotification() {
			if handler, ok := m.notificationHandlers[method(request.Method)]; ok {
				handler(request.Params)
			}
		} else {
			handler, ok := m.methodHandlers[method(request.Method)]
			if !ok {
				m.write(jsonrpc.NewMethodNotFoundError(request.ID, request.Method))
				return
			}
			result, err := handler(request.Params)
			if err != nil {
				m.write(jsonrpc.NewInternalError(request.ID, err))
				return
			}
			m.write(jsonrpc.NewResponse(request.ID, result))
		}
	}(request)
	return nil
}
