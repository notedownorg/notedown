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
)

// DefaultExcludePatterns provides default directory patterns to exclude from scanning
var DefaultExcludePatterns = []string{
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
}

// PathToFileURI converts a local path to a file:// URI
func PathToFileURI(path string) string {
	// Ensure path is absolute and clean
	absPath, _ := filepath.Abs(filepath.Clean(path))

	// Convert to URL path format
	urlPath := filepath.ToSlash(absPath)

	// Add file:// scheme
	return "file://" + urlPath
}

// URIToWorkspaceRoot converts a URI to WorkspaceRoot
func URIToWorkspaceRoot(uri, name string) (WorkspaceRoot, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return WorkspaceRoot{}, err
	}

	if parsedURI.Scheme != "file" {
		return WorkspaceRoot{}, nil // Skip non-file URIs
	}

	path := parsedURI.Path
	if name == "" {
		name = filepath.Base(path)
	}

	return WorkspaceRoot{
		URI:  uri,
		Path: path,
		Name: name,
	}, nil
}

// IsExcludedPath checks if a path should be excluded from indexing
func IsExcludedPath(path string, excludePatterns []string) bool {
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

// IsMarkdownFile checks if a file is a Markdown file
func IsMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md"
}

// DiscoverMarkdownFiles scans a workspace root for Markdown files
func DiscoverMarkdownFiles(root WorkspaceRoot, excludePatterns []string, maxFileCount int) ([]*FileInfo, error) {
	if excludePatterns == nil {
		excludePatterns = DefaultExcludePatterns
	}

	var files []*FileInfo
	fileCount := 0

	err := filepath.WalkDir(root.Path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Continue walking despite errors
		}

		// Skip directories
		if d.IsDir() {
			// Check if this directory should be excluded
			if IsExcludedPath(path, excludePatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if this is a Markdown file
		if !IsMarkdownFile(path) {
			return nil
		}

		// Check file count limit
		if maxFileCount > 0 && fileCount >= maxFileCount {
			return filepath.SkipAll
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Create file URI
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil
		}
		uri := PathToFileURI(absPath)

		// Create relative path from workspace root
		relPath, err := filepath.Rel(root.Path, path)
		if err != nil {
			relPath = path // Fallback to absolute path
		}

		// Add to results
		fileInfo := &FileInfo{
			URI:     uri,
			Path:    relPath,
			ModTime: info.ModTime(),
			Size:    info.Size(),
		}

		files = append(files, fileInfo)
		fileCount++

		return nil
	})

	return files, err
}
