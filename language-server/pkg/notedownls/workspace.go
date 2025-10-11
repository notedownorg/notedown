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
	"sync"

	"github.com/notedownorg/notedown/language-server/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
	"github.com/notedownorg/notedown/pkg/workspace"
)

// WorkspaceManager manages workspace-wide file discovery and indexing for the LSP server
type WorkspaceManager struct {
	manager    *workspace.Manager   // Embedded workspace manager
	loadedDocs map[string]*Document // URI -> full document (for opened files)
	mutex      sync.RWMutex
}

// Use types from the workspace package
type WorkspaceRoot = workspace.WorkspaceRoot
type FileInfo = workspace.FileInfo

// NewWorkspaceManager creates a new workspace manager
func NewWorkspaceManager(logger *log.Logger) *WorkspaceManager {
	return &WorkspaceManager{
		manager:    workspace.NewManager(logger),
		loadedDocs: make(map[string]*Document),
	}
}

// InitializeFromParams initializes workspace roots from LSP InitializeParams
func (wm *WorkspaceManager) InitializeFromParams(params lsp.InitializeParams) error {
	// Priority 1: WorkspaceFolders (modern approach)
	if len(params.WorkspaceFolders) > 0 {
		for _, folder := range params.WorkspaceFolders {
			err := wm.manager.AddRootByURI(folder.Uri, folder.Name)
			if err != nil {
				continue // Skip failed roots
			}
		}
		return nil
	}

	// Priority 2: RootUri (deprecated but still supported)
	if params.RootUri != nil && *params.RootUri != "" {
		err := wm.manager.AddRootByURI(*params.RootUri, "")
		if err != nil {
			return err
		}
		return nil
	}

	// Priority 3: RootPath (deprecated but still supported)
	if params.RootPath != nil && *params.RootPath != "" {
		err := wm.manager.AddRoot(*params.RootPath)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

// GetWorkspaceRoots returns the current workspace roots
func (wm *WorkspaceManager) GetWorkspaceRoots() []WorkspaceRoot {
	return wm.manager.GetWorkspaceRoots()
}

// GetFileIndex returns a copy of the file index
func (wm *WorkspaceManager) GetFileIndex() map[string]*FileInfo {
	return wm.manager.GetFileIndex()
}

// GetMarkdownFiles returns all indexed Markdown files
func (wm *WorkspaceManager) GetMarkdownFiles() []*FileInfo {
	return wm.manager.GetMarkdownFiles()
}

// DiscoverMarkdownFiles scans all workspace roots for Markdown files
func (wm *WorkspaceManager) DiscoverMarkdownFiles() error {
	return wm.manager.DiscoverMarkdownFiles()
}

// RefreshFileIndex performs an incremental refresh of the file index
func (wm *WorkspaceManager) RefreshFileIndex() error {
	return wm.manager.RefreshFileIndex()
}

// AddFileToIndex adds a single file to the index
func (wm *WorkspaceManager) AddFileToIndex(uri string) error {
	return wm.manager.AddFileToIndex(uri)
}

// RemoveFileFromIndex removes a file from the index
func (wm *WorkspaceManager) RemoveFileFromIndex(uri string) {
	wm.manager.RemoveFileFromIndex(uri)
}

// pathToFileURI converts a local path to a file:// URI
// This is kept for backward compatibility with tests
func pathToFileURI(path string) string {
	return workspace.PathToFileURI(path)
}

// SetupTestWorkspace is a test helper that sets up a workspace with mock files
// This allows tests to create file indexes without requiring actual files on disk
func (wm *WorkspaceManager) SetupTestWorkspace(roots []WorkspaceRoot, files map[string]*FileInfo) {
	// Add roots using the public interface
	for _, root := range roots {
		wm.manager.AddRootByURI(root.URI, root.Name)
	}

	// Set up file index using the test utility
	wm.manager.SetTestFileIndex(files)
}
