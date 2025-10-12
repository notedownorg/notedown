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

package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// workspaceDiscoverer handles workspace root discovery and file scanning
type workspaceDiscoverer struct {
	workspaceRoot   string
	excludePatterns []string
	mu              sync.RWMutex
}

// DocumentFile represents a discovered markdown file in the workspace
type DocumentFile struct {
	Path     string // Relative path from workspace root
	AbsPath  string // Absolute filesystem path
	Checksum string // SHA-256 hash of content
}

// newWorkspaceDiscoverer creates a new workspace discoverer
func newWorkspaceDiscoverer(workspaceRoot string) *workspaceDiscoverer {
	return &workspaceDiscoverer{
		workspaceRoot: workspaceRoot,
		excludePatterns: []string{
			".git",
			".vscode",
			".idea",
			"node_modules",
			".next",
			".nuxt",
			"dist",
			"build",
			"target",
			"__pycache__",
			".pytest_cache",
			"coverage",
			".coverage",
			"venv",
			".venv",
			"env",
			".env",
		},
	}
}

// discoverDocuments discovers all markdown documents in the workspace and streams them via channels
func (wd *workspaceDiscoverer) discoverDocuments() (<-chan *DocumentFile, <-chan error) {
	wd.mu.RLock()
	workspaceRoot := wd.workspaceRoot
	excludePatterns := make([]string, len(wd.excludePatterns))
	copy(excludePatterns, wd.excludePatterns)
	wd.mu.RUnlock()

	docChan := make(chan *DocumentFile)
	errChan := make(chan error, 1)

	go func() {
		defer close(docChan)
		defer close(errChan)

		err := filepath.WalkDir(workspaceRoot, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // Continue walking despite errors
			}

			// Skip directories
			if d.IsDir() {
				// Check if this directory should be excluded
				if wd.isExcludedPath(path, excludePatterns) {
					return filepath.SkipDir
				}
				return nil
			}

			// Check if this is a Markdown file
			if !wd.isMarkdownFile(path) {
				return nil
			}

			// Create relative path from workspace root
			relPath, err := filepath.Rel(workspaceRoot, path)
			if err != nil {
				relPath = path // Fallback to absolute path
			}

			// Calculate checksum
			checksum, err := wd.calculateChecksum(path)
			if err != nil {
				return nil // Skip files we can't read
			}

			// Send document to channel
			docChan <- &DocumentFile{
				Path:     relPath,
				AbsPath:  path,
				Checksum: checksum,
			}

			return nil
		})

		if err != nil {
			errChan <- fmt.Errorf("failed to walk workspace: %w", err)
		}
	}()

	return docChan, errChan
}

// isExcludedPath checks if a path should be excluded from indexing
func (wd *workspaceDiscoverer) isExcludedPath(path string, excludePatterns []string) bool {
	// Check against exclusion patterns
	for _, pattern := range excludePatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	// Skip hidden directories (starting with .)
	if strings.HasPrefix(filepath.Base(path), ".") {
		return true
	}

	return false
}

// isMarkdownFile checks if a file is a Markdown file
func (wd *workspaceDiscoverer) isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md"
}

// calculateChecksum calculates SHA-256 checksum of file content
func (wd *workspaceDiscoverer) calculateChecksum(path string) (string, error) {
	// #nosec G304 - path is from trusted workspace discovery
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't fail the operation
			_ = err // Explicitly ignore error
		}
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
