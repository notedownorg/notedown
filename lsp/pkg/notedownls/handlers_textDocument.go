package notedownls

import (
	"encoding/json"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
)

// handleDidOpen handles textDocument/didOpen notifications
func (s *Server) handleDidOpen(params json.RawMessage) error {
	var didOpenParams lsp.DidOpenTextDocumentParams
	if err := json.Unmarshal(params, &didOpenParams); err != nil {
		s.logger.Error("failed to unmarshal didOpen params", "error", err)
		return err
	}

	uri := didOpenParams.TextDocument.URI
	_, err := s.AddDocument(uri)
	if err != nil {
		s.logger.Error("failed to add document", "uri", uri, "error", err)
		return err
	}

	s.logger.Info("document opened", "uri", uri, "languageId", didOpenParams.TextDocument.LanguageID)
	return nil
}

// handleDidClose handles textDocument/didClose notifications
func (s *Server) handleDidClose(params json.RawMessage) error {
	var didCloseParams lsp.DidCloseTextDocumentParams
	if err := json.Unmarshal(params, &didCloseParams); err != nil {
		s.logger.Error("failed to unmarshal didClose params", "error", err)
		return err
	}

	uri := didCloseParams.TextDocument.URI
	s.RemoveDocument(uri)
	s.logger.Info("document closed", "uri", uri)
	return nil
}

// handleDidChange handles textDocument/didChange notifications
func (s *Server) handleDidChange(params json.RawMessage) error {
	var didChangeParams lsp.DidChangeTextDocumentParams
	if err := json.Unmarshal(params, &didChangeParams); err != nil {
		s.logger.Error("failed to unmarshal didChange params", "error", err)
		return err
	}

	uri := didChangeParams.TextDocument.URI
	version := didChangeParams.TextDocument.Version
	changeCount := len(didChangeParams.ContentChanges)

	s.logger.Debug("document changed", "uri", uri, "version", version, "changes", changeCount)
	return nil
}
