package lsp

// LanguageServer defines the interface that LSP server implementations must satisfy
type LanguageServer interface {
	// Initialize handles the LSP initialize request
	Initialize(params InitializeParams) (InitializeResult, error)

	// RegisterHandlers registers all method and notification handlers with the mux
	RegisterHandlers(mux *Mux) error

	// Shutdown handles cleanup when the server is shutting down
	Shutdown() error
}
