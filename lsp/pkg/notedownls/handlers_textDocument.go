package notedownls

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	// Generate and publish diagnostics for this document
	diagnostics := s.generateWikilinkDiagnostics(uri, content)
	s.publishDiagnostics(uri, diagnostics)

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

	// Clear diagnostics for this document
	s.publishDiagnostics(uri, []lsp.Diagnostic{})

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
			
			// Generate and publish updated diagnostics
			diagnostics := s.generateWikilinkDiagnostics(uri, newContent)
			s.publishDiagnostics(uri, diagnostics)
		}
		s.documentsMutex.Unlock()
	}

	s.logger.Debug("document changed", "uri", uri, "version", version, "changes", changeCount)
	return nil
}

// handleDefinition handles textDocument/definition requests for goto definition
func (s *Server) handleDefinition(params json.RawMessage) (any, error) {
	var definitionParams lsp.DefinitionParams
	if err := json.Unmarshal(params, &definitionParams); err != nil {
		s.logger.Error("failed to unmarshal definition params", "error", err)
		return nil, err
	}

	s.logger.Debug("definition request received",
		"uri", definitionParams.TextDocument.URI,
		"line", definitionParams.Position.Line,
		"character", definitionParams.Position.Character)

	// Get the document
	doc, exists := s.GetDocument(definitionParams.TextDocument.URI)
	if !exists {
		s.logger.Debug("document not found for definition", "uri", definitionParams.TextDocument.URI)
		return nil, nil
	}

	// Extract the complete wikilink target at cursor position
	target := s.getCompleteWikilinkTarget(doc, definitionParams.Position)
	if target == "" {
		s.logger.Debug("not in wikilink context or empty target")
		return nil, nil
	}

	s.logger.Debug("detected wikilink target for definition", "target", target)

	// Try to find existing file for the target
	existingFile := s.findFileForTarget(target)
	if existingFile != nil {
		s.logger.Debug("found existing file for target", "target", target, "file", existingFile.Path)

		// Return location pointing to the existing file
		return lsp.Location{
			URI: existingFile.URI,
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 0, Character: 0},
			},
		}, nil
	}

	// File doesn't exist - offer to create it
	s.logger.Debug("file doesn't exist for target, offering to create", "target", target)

	// Determine the target file path - handle directory structure
	targetPath, targetURI := s.resolveTargetPath(target)
	if targetPath == "" {
		s.logger.Error("failed to resolve target path", "target", target)
		return nil, fmt.Errorf("failed to resolve target path")
	}

	// Create the file
	if err := s.createMarkdownFile(targetPath, target); err != nil {
		s.logger.Error("failed to create file", "path", targetPath, "error", err)
		return nil, err
	}

	s.logger.Info("created new file for wikilink target", "target", target, "path", targetPath)

	// Return location pointing to the newly created file
	return lsp.Location{
		URI: targetURI,
		Range: lsp.Range{
			Start: lsp.Position{Line: 0, Character: 0},
			End:   lsp.Position{Line: 0, Character: 0},
		},
	}, nil
}

// findFileForTarget attempts to find an existing file that matches the wikilink target
func (s *Server) findFileForTarget(target string) *FileInfo {
	// Normalize the target (replace backslashes with forward slashes)
	normalizedTarget := strings.ReplaceAll(target, "\\", "/")

	// Reject targets containing .. sequences to prevent directory traversal
	if strings.Contains(normalizedTarget, "..") {
		s.logger.Debug("target contains directory traversal sequences, rejecting", "target", target)
		return nil
	}

	workspaceFiles := s.GetWorkspaceFiles()

	for _, fileInfo := range workspaceFiles {
		// Check if the file path (without extension) matches the target
		pathWithoutExt := strings.TrimSuffix(fileInfo.Path, filepath.Ext(fileInfo.Path))
		baseWithoutExt := strings.TrimSuffix(filepath.Base(fileInfo.Path), filepath.Ext(fileInfo.Path))

		// Normalize paths for comparison (handle both forward and backward slashes)
		normalizedPath := strings.ReplaceAll(pathWithoutExt, "\\", "/")

		// Match exact path or base name
		if normalizedPath == normalizedTarget || baseWithoutExt == target {
			return fileInfo
		}
	}

	return nil
}

// getCompleteWikilinkTarget extracts the complete wikilink target at the cursor position
func (s *Server) getCompleteWikilinkTarget(doc *Document, position lsp.Position) string {
	lines := strings.Split(doc.Content, "\n")
	if position.Line >= len(lines) {
		return ""
	}

	line := lines[position.Line]
	if position.Character > len(line) {
		return ""
	}

	// Find the wikilink that contains the cursor position
	// Look for [[ before the cursor and ]] after the cursor

	// Find the rightmost [[ before or at cursor position
	beforeCursor := line[:position.Character+1]
	lastWikilinkStart := strings.LastIndex(beforeCursor, "[[")
	if lastWikilinkStart == -1 {
		return ""
	}

	// Find the leftmost ]] after the wikilink start
	afterWikilinkStart := line[lastWikilinkStart:]
	wikilinkEnd := strings.Index(afterWikilinkStart, "]]")
	if wikilinkEnd == -1 {
		return "" // Incomplete wikilink
	}

	// Calculate the absolute position of the wikilink end
	absoluteWikilinkEnd := lastWikilinkStart + wikilinkEnd + 2 // +2 for ]]

	// Check if cursor is actually within the wikilink bounds
	if position.Character < lastWikilinkStart || position.Character >= absoluteWikilinkEnd {
		return "" // Cursor is outside the wikilink
	}

	// Extract the complete wikilink content
	wikilinkContent := afterWikilinkStart[2:wikilinkEnd] // Skip [[ and ]]

	// Handle pipe separators - take only the target part (before |)
	if pipeIndex := strings.Index(wikilinkContent, "|"); pipeIndex != -1 {
		wikilinkContent = wikilinkContent[:pipeIndex]
	}

	target := strings.TrimSpace(wikilinkContent)

	s.logger.Debug("extracted complete wikilink target",
		"target", target,
		"line", position.Line,
		"character", position.Character)

	return target
}

// resolveTargetPath determines the appropriate file path and URI for a wikilink target
func (s *Server) resolveTargetPath(target string) (string, string) {
	// Get workspace root for file creation
	workspaceRoots := s.GetWorkspaceRoots()
	if len(workspaceRoots) == 0 {
		s.logger.Error("no workspace roots available for file creation")
		return "", ""
	}

	workspaceRoot := workspaceRoots[0].Path

	// Normalize the target (replace backslashes with forward slashes)
	normalizedTarget := strings.ReplaceAll(target, "\\", "/")

	// Reject targets containing .. sequences to prevent directory traversal
	if strings.Contains(normalizedTarget, "..") {
		s.logger.Error("target contains directory traversal sequences", "target", target)
		return "", ""
	}

	// If target contains slashes, it might be a path-based wikilink
	var targetPath string
	if strings.Contains(normalizedTarget, "/") {
		// Handle path-based targets like "docs/api" or "projects/ideas"
		targetPath = filepath.Join(workspaceRoot, normalizedTarget+".md")
	} else {
		// Simple filename - place in root
		targetPath = filepath.Join(workspaceRoot, normalizedTarget+".md")
	}

	// Convert to URI format
	targetURI := "file://" + strings.ReplaceAll(targetPath, "\\", "/")

	return targetPath, targetURI
}

// createMarkdownFile creates a new markdown file with basic content
func (s *Server) createMarkdownFile(filePath, title string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create basic markdown content
	content := fmt.Sprintf("# %s\n\n", title)

	// Write the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
