package notedownls

import (
	"sync"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

// Server implements the LSPServer interface for Notedown
type Server struct {
	version string
	logger  *log.Logger

	// Document storage
	documents      map[string]*Document
	documentsMutex sync.RWMutex
}

// NewServer creates a new Notedown LSP server
func NewServer(version string, logger *log.Logger) *Server {
	return &Server{
		version:   version,
		logger:    logger,
		documents: make(map[string]*Document),
	}
}

// Initialize handles the LSP initialize request
func (s *Server) Initialize(params lsp.InitializeParams) (lsp.InitializeResult, error) {
	clientName := "unknown"
	if params.ClientInfo != nil {
		clientName = params.ClientInfo.Name
	}

	s.logger.Info("lsp client initialized", "client", clientName, "server_version", s.version)

	syncKind := lsp.TextDocumentSyncKindFull
	result := lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "Notedown Language Server", Version: s.version},
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: &[]bool{true}[0],
				Change:    &syncKind,
			},
		},
	}
	return result, nil
}

// RegisterHandlers registers all method and notification handlers
func (s *Server) RegisterHandlers(mux *lsp.Mux) error {
	// Register document lifecycle handlers
	mux.RegisterNotification(lsp.MethodTextDocumentDidOpen, s.handleDidOpen)
	mux.RegisterNotification(lsp.MethodTextDocumentDidChange, s.handleDidChange)
	mux.RegisterNotification(lsp.MethodTextDocumentDidClose, s.handleDidClose)

	s.logger.Debug("registered document lifecycle handlers")
	return nil
}

// GetDocument retrieves a document by URI
func (s *Server) GetDocument(uri string) (*Document, bool) {
	s.documentsMutex.RLock()
	defer s.documentsMutex.RUnlock()
	doc, exists := s.documents[uri]
	return doc, exists
}

// AddDocument adds or updates a document in storage
func (s *Server) AddDocument(uri string) (*Document, error) {
	doc, err := NewDocument(uri)
	if err != nil {
		return nil, err
	}

	s.documentsMutex.Lock()
	defer s.documentsMutex.Unlock()
	s.documents[uri] = doc

	s.logger.Debug("document added", "uri", uri, "basepath", doc.Basepath)
	return doc, nil
}

// RemoveDocument removes a document from storage
func (s *Server) RemoveDocument(uri string) {
	s.documentsMutex.Lock()
	defer s.documentsMutex.Unlock()

	if doc, exists := s.documents[uri]; exists {
		delete(s.documents, uri)
		s.logger.Debug("document removed", "uri", uri, "basepath", doc.Basepath)
	}
}

// HasDocument checks if a document exists in storage
func (s *Server) HasDocument(uri string) bool {
	s.documentsMutex.RLock()
	defer s.documentsMutex.RUnlock()
	_, exists := s.documents[uri]
	return exists
}

// Shutdown handles cleanup when the server is shutting down
func (s *Server) Shutdown() error {
	s.logger.Info("shutting down Notedown language server", "version", s.version)
	return nil
}
