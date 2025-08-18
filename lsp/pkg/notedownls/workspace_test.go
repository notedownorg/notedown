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
	"os"
	"path/filepath"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestNewWorkspaceManager(t *testing.T) {
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	if len(wm.roots) != 0 {
		t.Errorf("Expected 0 roots, got %d", len(wm.roots))
	}

	if len(wm.fileIndex) != 0 {
		t.Errorf("Expected empty file index, got %d files", len(wm.fileIndex))
	}

	if wm.maxFileCount != 10000 {
		t.Errorf("Expected maxFileCount 10000, got %d", wm.maxFileCount)
	}
}

func TestInitializeFromParams(t *testing.T) {
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	tests := []struct {
		name           string
		params         lsp.InitializeParams
		expectedRoots  int
		shouldHaveRoot bool
	}{
		{
			name: "with workspace folders",
			params: lsp.InitializeParams{
				WorkspaceFolders: []lsp.WorkspaceFolder{
					{Uri: "file:///test/workspace", Name: "test-workspace"},
				},
			},
			expectedRoots:  1,
			shouldHaveRoot: true,
		},
		{
			name: "with rootUri",
			params: lsp.InitializeParams{
				RootUri: &[]string{"file:///test/workspace"}[0],
			},
			expectedRoots:  1,
			shouldHaveRoot: true,
		},
		{
			name: "with rootPath",
			params: lsp.InitializeParams{
				RootPath: &[]string{"/test/workspace"}[0],
			},
			expectedRoots:  1,
			shouldHaveRoot: true,
		},
		{
			name:           "no workspace info",
			params:         lsp.InitializeParams{},
			expectedRoots:  0,
			shouldHaveRoot: false,
		},
		{
			name: "workspace folders take priority over rootUri",
			params: lsp.InitializeParams{
				WorkspaceFolders: []lsp.WorkspaceFolder{
					{Uri: "file:///test/workspace1", Name: "workspace1"},
				},
				RootUri: &[]string{"file:///test/workspace2"}[0],
			},
			expectedRoots:  1,
			shouldHaveRoot: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := wm.InitializeFromParams(tt.params)
			if err != nil {
				t.Errorf("InitializeFromParams failed: %v", err)
			}

			roots := wm.GetWorkspaceRoots()
			if len(roots) != tt.expectedRoots {
				t.Errorf("Expected %d roots, got %d", tt.expectedRoots, len(roots))
			}

			if tt.shouldHaveRoot && len(roots) > 0 {
				root := roots[0]
				if root.URI == "" {
					t.Error("Root URI should not be empty")
				}
				if root.Path == "" {
					t.Error("Root path should not be empty")
				}
			}
		})
	}
}

func TestPathToFileURI(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "absolute path",
			path:     "/test/path",
			expected: "file:///test/path",
		},
		{
			name:     "relative path gets converted to absolute",
			path:     "test/path",
			expected: "file://" + func() string { wd, _ := os.Getwd(); return filepath.Join(wd, "test/path") }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pathToFileURI(tt.path)
			// For relative paths, we just check that it starts with file://
			// and contains our path, since the absolute path depends on current directory
			if tt.path == "test/path" {
				if !filepath.IsAbs(result[7:]) { // Remove "file://" prefix
					t.Errorf("Expected absolute path in URI, got %s", result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestIsMarkdownFile(t *testing.T) {
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"markdown file", "README.md", true},
		{"uppercase markdown", "DOC.MD", true},
		{"text file", "notes.txt", false},
		{"no extension", "README", false},
		{"other extension", "config.json", false},
		{"markdown in path but not extension", "markdown/file.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wm.isMarkdownFile(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.path, result)
			}
		})
	}
}

func TestIsExcludedPath(t *testing.T) {
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"git directory", "/project/.git/config", true},
		{"node_modules", "/project/node_modules/lib", true},
		{"vscode directory", "/project/.vscode/settings.json", true},
		{"hidden file", "/project/.hidden", true},
		{"normal file", "/project/README.md", false},
		{"normal directory", "/project/docs/guide.md", false},
		{"git in filename but not directory", "/project/git-info.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wm.isExcludedPath(tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v for %s, got %v", tt.expected, tt.path, result)
			}
		})
	}
}

func TestDiscoverMarkdownFilesIntegration(t *testing.T) {
	// Create a temporary workspace for testing
	tempDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"README.md":               "# Main README",
		"docs/guide.md":           "# User Guide",
		"docs/api.md":             "# API Documentation",
		"src/README.md":           "# Source Code",
		"config.json":             `{"key": "value"}`, // Should be ignored
		"notes.txt":               "Some notes",       // Should be ignored
		".git/config":             "git config",       // Should be excluded
		"node_modules/lib/pkg.js": "module code",      // Should be excluded
	}

	// Create directories and files
	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)

		if err := os.MkdirAll(dir, 0750); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Test workspace discovery
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	// Initialize with temp directory as workspace root
	params := lsp.InitializeParams{
		RootUri: &[]string{pathToFileURI(tempDir)}[0],
	}

	err := wm.InitializeFromParams(params)
	if err != nil {
		t.Fatalf("Failed to initialize workspace: %v", err)
	}

	// Discover files
	err = wm.DiscoverMarkdownFiles()
	if err != nil {
		t.Fatalf("Failed to discover files: %v", err)
	}

	// Check results
	files := wm.GetMarkdownFiles()

	// Should find exactly 4 markdown files (excluding those in excluded directories)
	expectedFiles := []string{"README.md", "docs/guide.md", "docs/api.md", "src/README.md"}
	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
		for _, file := range files {
			t.Logf("Found file: %s", file.Path)
		}
	}

	// Verify file information
	pathMap := make(map[string]*FileInfo)
	for _, file := range files {
		pathMap[file.Path] = file
	}

	for _, expectedPath := range expectedFiles {
		if file, exists := pathMap[expectedPath]; !exists {
			t.Errorf("Expected file %s not found in results", expectedPath)
		} else {
			// Verify file info
			if file.URI == "" {
				t.Errorf("File %s has empty URI", expectedPath)
			}
			if file.Size <= 0 {
				t.Errorf("File %s has invalid size: %d", expectedPath, file.Size)
			}
			if file.ModTime.IsZero() {
				t.Errorf("File %s has zero mod time", expectedPath)
			}
		}
	}
}

func TestWorkspaceManagerConcurrency(t *testing.T) {
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	// Initialize with a dummy workspace
	params := lsp.InitializeParams{
		RootUri: &[]string{"file:///test"}[0],
	}
	_ = wm.InitializeFromParams(params)

	// Test concurrent access
	done := make(chan bool)

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			_ = wm.GetMarkdownFiles()
			_ = wm.GetWorkspaceRoots()
			_ = wm.GetFileIndex()
		}
		done <- true
	}()

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			testURI := "file:///test/doc.md"
			_ = wm.AddFileToIndex(testURI)
			wm.RemoveFileFromIndex(testURI)
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// Test should complete without race conditions or deadlocks
}

func TestAddRemoveFileFromIndex(t *testing.T) {
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	// Initialize with a workspace root
	params := lsp.InitializeParams{
		RootUri: &[]string{"file:///test"}[0],
	}
	_ = wm.InitializeFromParams(params)

	// Create a temporary file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.md")
	if err := os.WriteFile(tempFile, []byte("# Test"), 0600); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	testURI := pathToFileURI(tempFile)

	// Test adding file to index
	err := wm.AddFileToIndex(testURI)
	if err != nil {
		t.Errorf("Failed to add file to index: %v", err)
	}

	// Verify file was added
	files := wm.GetMarkdownFiles()
	found := false
	for _, file := range files {
		if file.URI == testURI {
			found = true
			break
		}
	}
	if !found {
		t.Error("File was not added to index")
	}

	// Test removing file from index
	wm.RemoveFileFromIndex(testURI)

	// Verify file was removed
	files = wm.GetMarkdownFiles()
	for _, file := range files {
		if file.URI == testURI {
			t.Error("File was not removed from index")
			break
		}
	}
}

func TestMaxFileCountLimit(t *testing.T) {
	logger := log.NewDefault()
	wm := NewWorkspaceManager(logger)

	// Set a low limit for testing
	wm.maxFileCount = 2

	// Create a temporary workspace with more files than the limit
	tempDir := t.TempDir()
	testFiles := []string{"file1.md", "file2.md", "file3.md", "file4.md"}

	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(filePath, []byte("# "+fileName), 0600); err != nil {
			t.Fatalf("Failed to create file %s: %v", fileName, err)
		}
	}

	// Initialize workspace
	params := lsp.InitializeParams{
		RootUri: &[]string{pathToFileURI(tempDir)}[0],
	}
	_ = wm.InitializeFromParams(params)

	// Discover files - should be limited by maxFileCount
	err := wm.DiscoverMarkdownFiles()
	if err != nil {
		t.Fatalf("Failed to discover files: %v", err)
	}

	files := wm.GetMarkdownFiles()
	if len(files) > wm.maxFileCount {
		t.Errorf("Expected at most %d files due to limit, got %d", wm.maxFileCount, len(files))
	}
}
