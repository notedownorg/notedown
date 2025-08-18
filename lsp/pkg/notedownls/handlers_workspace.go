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

package notedownls

import (
	"encoding/json"
	"fmt"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
)

// handleDidChangeWatchedFiles handles workspace/didChangeWatchedFiles notifications
func (s *Server) handleDidChangeWatchedFiles(params json.RawMessage) error {
	var watchedFilesParams lsp.DidChangeWatchedFilesParams
	if err := json.Unmarshal(params, &watchedFilesParams); err != nil {
		s.logger.Error("failed to unmarshal didChangeWatchedFiles params", "error", err)
		return err
	}

	for _, change := range watchedFilesParams.Changes {
		s.handleFileSystemChange(change)
	}

	return nil
}

// handleFileSystemChange processes a single file system change event
func (s *Server) handleFileSystemChange(event lsp.FileEvent) {
	uri := event.URI
	changeType := event.Type

	s.logger.Debug("external file change detected", "uri", uri, "type", changeType)

	switch changeType {
	case lsp.FileChangeTypeCreated:
		s.handleExternalFileCreated(uri)
	case lsp.FileChangeTypeChanged:
		s.handleExternalFileChanged(uri)
	case lsp.FileChangeTypeDeleted:
		s.handleExternalFileDeleted(uri)
	default:
		s.logger.Warn("unknown file change type", "uri", uri, "type", changeType)
	}
}

// handleExternalFileCreated handles external file creation
func (s *Server) handleExternalFileCreated(uri string) {
	// Add to workspace index if it's a Markdown file
	if err := s.workspace.AddFileToIndex(uri); err != nil {
		s.logger.Warn("failed to add created file to workspace index", "uri", uri, "error", err)
	} else {
		s.logger.Info("external Markdown file created and indexed", "uri", uri)

		// File creation might resolve conflicts or create new ones
		s.refreshAllDocumentDiagnostics()
	}

	// Document store handling remains the same - don't auto-add to document store
	if !s.HasDocument(uri) {
		if doc, err := NewDocument(uri); err == nil {
			s.logger.Debug("external file created", "uri", uri, "basepath", doc.Basepath)
		}
	}
}

// handleExternalFileChanged handles external file modification
func (s *Server) handleExternalFileChanged(uri string) {
	// Update workspace index if it's a Markdown file in our index
	if err := s.workspace.AddFileToIndex(uri); err != nil {
		s.logger.Debug("could not update workspace index for changed file", "uri", uri, "error", err)
	}

	if s.HasDocument(uri) {
		// File is being tracked - this means the content has changed externally
		// The LSP client should handle this by sending a didChange notification
		s.logger.Info("tracked file changed externally", "uri", uri)
	} else {
		// File is not being tracked by LSP
		if doc, err := NewDocument(uri); err == nil {
			s.logger.Debug("untracked file changed externally", "uri", uri, "basepath", doc.Basepath)
		} else {
			s.logger.Debug("untracked file changed externally", "uri", uri)
		}
	}
}

// handleExternalFileDeleted handles external file deletion
func (s *Server) handleExternalFileDeleted(uri string) {
	// Remove from workspace index
	s.workspace.RemoveFileFromIndex(uri)

	if s.HasDocument(uri) {
		// File was being tracked but has been deleted externally
		s.RemoveDocument(uri)
		s.logger.Info("tracked file deleted externally and removed from workspace index", "uri", uri)
	} else {
		// File was not being tracked
		s.logger.Debug("untracked file deleted externally and removed from workspace index", "uri", uri)
	}

	// File deletion might resolve conflicts or create new ones
	s.refreshAllDocumentDiagnostics()
}

// RegisterFileWatcher registers file watchers with the LSP client
// This method can be called after initialization to set up file watching
func (s *Server) RegisterFileWatcher(clientRegister func(lsp.RegistrationParams) error) error {
	// Watch for .md files in the workspace (Markdown-focused)
	watchers := []lsp.FileSystemWatcher{
		{
			GlobPattern: "**/*.md",
			Kind:        &[]lsp.WatchKind{lsp.WatchKindCreate | lsp.WatchKindChange | lsp.WatchKindDelete}[0],
		},
	}

	registration := lsp.Registration{
		ID:     "notedown-file-watcher",
		Method: "workspace/didChangeWatchedFiles",
		RegisterOptions: lsp.DidChangeWatchedFilesRegistrationOptions{
			Watchers: watchers,
		},
	}

	params := lsp.RegistrationParams{
		Registrations: []lsp.Registration{registration},
	}

	if err := clientRegister(params); err != nil {
		s.logger.Error("failed to register file watchers", "error", err)
		return fmt.Errorf("failed to register file watchers: %w", err)
	}

	s.logger.Info("file watchers registered successfully", "patterns", []string{"**/*.md"})
	return nil
}
