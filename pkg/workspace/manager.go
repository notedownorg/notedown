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

package workspace

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/notedownorg/notedown/pkg/log"
)

// Manager manages workspace-wide file discovery and indexing
type Manager struct {
	roots     []WorkspaceRoot
	fileIndex map[string]*FileInfo // URI -> lightweight file info
	mutex     sync.RWMutex
	logger    *log.Logger

	// Configuration
	maxFileCount    int
	excludePatterns []string
}

// NewManager creates a new workspace manager
func NewManager(logger *log.Logger) *Manager {
	return &Manager{
		roots:     make([]WorkspaceRoot, 0),
		fileIndex: make(map[string]*FileInfo),
		logger:    logger,

		// Default configuration
		maxFileCount:    10000,
		excludePatterns: DefaultExcludePatterns,
	}
}

// SetExcludePatterns sets the patterns for directories to exclude from scanning
func (wm *Manager) SetExcludePatterns(patterns []string) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	wm.excludePatterns = patterns
}

// SetMaxFileCount sets the maximum number of files to index
func (wm *Manager) SetMaxFileCount(count int) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()
	wm.maxFileCount = count
}

// AddRoot adds a workspace root by path
func (wm *Manager) AddRoot(path string) error {
	return wm.AddRootWithName(path, "")
}

// AddRootWithName adds a workspace root with a custom name
func (wm *Manager) AddRootWithName(path, name string) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	// Convert path to file:// URI
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	uri := PathToFileURI(absPath)
	root, err := URIToWorkspaceRoot(uri, name)
	if err != nil {
		return err
	}

	wm.roots = append(wm.roots, root)
	wm.logger.Debug("added workspace root", "path", path, "uri", root.URI, "name", root.Name)
	return nil
}

// AddRootByURI adds a workspace root by URI
func (wm *Manager) AddRootByURI(uri, name string) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	root, err := URIToWorkspaceRoot(uri, name)
	if err != nil {
		return err
	}

	wm.roots = append(wm.roots, root)
	wm.logger.Debug("added workspace root from uri", "uri", root.URI, "name", root.Name)
	return nil
}

// GetWorkspaceRoots returns the current workspace roots
func (wm *Manager) GetWorkspaceRoots() []WorkspaceRoot {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	roots := make([]WorkspaceRoot, len(wm.roots))
	copy(roots, wm.roots)
	return roots
}

// GetFileIndex returns a copy of the file index
func (wm *Manager) GetFileIndex() map[string]*FileInfo {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	index := make(map[string]*FileInfo)
	for uri, info := range wm.fileIndex {
		// Create a copy of the FileInfo
		indexCopy := *info
		index[uri] = &indexCopy
	}
	return index
}

// GetMarkdownFiles returns all indexed Markdown files
func (wm *Manager) GetMarkdownFiles() []*FileInfo {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	files := make([]*FileInfo, 0, len(wm.fileIndex))
	for _, info := range wm.fileIndex {
		// Create a copy of the FileInfo
		infoCopy := *info
		files = append(files, &infoCopy)
	}
	return files
}

// DiscoverMarkdownFiles scans all workspace roots for Markdown files
func (wm *Manager) DiscoverMarkdownFiles() error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	// Clear existing index
	oldCount := len(wm.fileIndex)
	wm.fileIndex = make(map[string]*FileInfo)
	if oldCount > 0 {
		wm.logger.Debug("cleared workspace cache for rediscovery", "previousFiles", oldCount)
	}

	totalFiles := 0
	for _, root := range wm.roots {
		wm.logger.Info("scanning workspace root for Markdown files", "root", root.Name, "path", root.Path)

		files, err := DiscoverMarkdownFiles(root, wm.excludePatterns, wm.maxFileCount-totalFiles)
		if err != nil {
			wm.logger.Error("failed to scan workspace root", "root", root.Name, "path", root.Path, "error", err)
			continue
		}

		// Add files to index
		for _, fileInfo := range files {
			wm.fileIndex[fileInfo.URI] = fileInfo
			totalFiles++
			wm.logger.Debug("added file to workspace cache", "uri", fileInfo.URI, "path", fileInfo.Path, "size", fileInfo.Size)

			// Check if we've reached the limit
			if wm.maxFileCount > 0 && totalFiles >= wm.maxFileCount {
				wm.logger.Warn("reached maximum file count limit, stopping scan",
					"limit", wm.maxFileCount, "root", root.Name)
				break
			}
		}

		if wm.maxFileCount > 0 && totalFiles >= wm.maxFileCount {
			break
		}
	}

	wm.logger.Info("workspace Markdown file discovery completed",
		"totalFiles", totalFiles, "roots", len(wm.roots))

	return nil
}

// RefreshFileIndex performs an incremental refresh of the file index
func (wm *Manager) RefreshFileIndex() error {
	// For now, perform a full rescan
	// TODO: Implement incremental refresh based on file modification times
	return wm.DiscoverMarkdownFiles()
}

// AddFileToIndex adds a single file to the index
func (wm *Manager) AddFileToIndex(uri string) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	// Parse URI to get path
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return err
	}

	if parsedURI.Scheme != "file" {
		return nil // Skip non-file URIs
	}

	path := parsedURI.Path

	// Check if it's a Markdown file
	if !IsMarkdownFile(path) {
		return nil
	}

	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Find the appropriate workspace root
	var relPath string
	for _, root := range wm.roots {
		if strings.HasPrefix(path, root.Path) {
			var err error
			relPath, err = filepath.Rel(root.Path, path)
			if err != nil {
				relPath = path
			}
			break
		}
	}
	if relPath == "" {
		relPath = path // Fallback if no root matches
	}

	// Add to index
	fileInfo := &FileInfo{
		URI:     uri,
		Path:    relPath,
		ModTime: info.ModTime(),
		Size:    info.Size(),
	}

	wm.fileIndex[uri] = fileInfo
	wm.logger.Debug("added file to workspace cache", "uri", uri, "path", relPath, "size", info.Size())

	return nil
}

// RemoveFileFromIndex removes a file from the index
func (wm *Manager) RemoveFileFromIndex(uri string) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	if fileInfo, exists := wm.fileIndex[uri]; exists {
		delete(wm.fileIndex, uri)
		wm.logger.Debug("removed file from workspace cache", "uri", uri, "path", fileInfo.Path)
	}
}

// SetTestFileIndex sets the file index for testing purposes
// This allows tests to mock the workspace file index without requiring actual files on disk
func (wm *Manager) SetTestFileIndex(files map[string]*FileInfo) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	// Create a new map and copy files into it
	wm.fileIndex = make(map[string]*FileInfo)
	for uri, info := range files {
		// Create a copy of the FileInfo to avoid shared references
		infoCopy := *info
		wm.fileIndex[uri] = &infoCopy
	}
}
