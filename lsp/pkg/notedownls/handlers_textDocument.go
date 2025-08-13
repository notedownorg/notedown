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
	content := didOpenParams.TextDocument.Text
	version := didOpenParams.TextDocument.Version

	doc, err := NewDocumentWithContent(uri, content, version)
	if err != nil {
		s.logger.Error("failed to create document", "uri", uri, "error", err)
		return err
	}

	s.documentsMutex.Lock()
	s.documents[uri] = doc
	s.documentsMutex.Unlock()

	// Extract wikilinks from the document content
	s.extractWikilinksFromDocument(uri, content)

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
	
	// Remove wikilink references from this document
	s.removeWikilinksFromDocument(uri)
	
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
	version := *didChangeParams.TextDocument.Version
	changeCount := len(didChangeParams.ContentChanges)

	// For full text sync, we expect a single change with the full content
	if changeCount > 0 {
		// Get the document
		s.documentsMutex.Lock()
		doc, exists := s.documents[uri]
		if exists && len(didChangeParams.ContentChanges) > 0 {
			// Update with the new content (assuming full text sync)
			newContent := didChangeParams.ContentChanges[0].Text
			doc.UpdateContent(newContent, version)
			
			// Update wikilinks for this document
			s.refreshWikilinksFromDocument(uri, newContent)
		}
		s.documentsMutex.Unlock()
	}

	s.logger.Debug("document changed", "uri", uri, "version", version, "changes", changeCount)
	return nil
}
