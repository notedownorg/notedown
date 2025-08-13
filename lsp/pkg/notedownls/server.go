package notedownls

import (
	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

// Server implements the LSPServer interface for Notedown
type Server struct {
	version string
	logger  *log.Logger
}

// NewServer creates a new Notedown LSP server
func NewServer(version string, logger *log.Logger) *Server {
	return &Server{
		version: version,
		logger:  logger,
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
	// TODO: Register document lifecycle handlers
	// mux.RegisterNotification(lsp.MethodTextDocumentDidOpen, s.handleDidOpen)
	// mux.RegisterNotification(lsp.MethodTextDocumentDidChange, s.handleDidChange)
	// mux.RegisterNotification(lsp.MethodTextDocumentDidClose, s.handleDidClose)

	return nil
}

// Shutdown handles cleanup when the server is shutting down
func (s *Server) Shutdown() error {
	s.logger.Info("shutting down Notedown language server", "version", s.version)
	return nil
}
