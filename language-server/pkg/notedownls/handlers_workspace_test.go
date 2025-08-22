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
	"slices"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestHandleDidChangeWatchedFiles(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Create test parameters with multiple file changes
	params := lsp.DidChangeWatchedFilesParams{
		Changes: []lsp.FileEvent{
			{
				URI:  "file:///test/created.md",
				Type: lsp.FileChangeTypeCreated,
			},
			{
				URI:  "file:///test/modified.md",
				Type: lsp.FileChangeTypeChanged,
			},
			{
				URI:  "file:///test/deleted.md",
				Type: lsp.FileChangeTypeDeleted,
			},
		},
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	// Test handler
	err = server.handleDidChangeWatchedFiles(paramBytes)
	if err != nil {
		t.Errorf("handleDidChangeWatchedFiles failed: %v", err)
	}
}

func TestHandleFileSystemChange(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	testURI := "file:///test/document.md"

	tests := []struct {
		name       string
		changeType lsp.FileChangeType
		setupDoc   bool // whether to add document before test
		expectDoc  bool // whether document should exist after test
	}{
		{
			name:       "external file created",
			changeType: lsp.FileChangeTypeCreated,
			setupDoc:   false,
			expectDoc:  false, // External creation doesn't auto-add to document store
		},
		{
			name:       "external tracked file changed",
			changeType: lsp.FileChangeTypeChanged,
			setupDoc:   true,
			expectDoc:  true, // Document should still exist
		},
		{
			name:       "external untracked file changed",
			changeType: lsp.FileChangeTypeChanged,
			setupDoc:   false,
			expectDoc:  false, // Document wasn't tracked, so still not tracked
		},
		{
			name:       "external tracked file deleted",
			changeType: lsp.FileChangeTypeDeleted,
			setupDoc:   true,
			expectDoc:  false, // Document should be removed
		},
		{
			name:       "external untracked file deleted",
			changeType: lsp.FileChangeTypeDeleted,
			setupDoc:   false,
			expectDoc:  false, // Document wasn't tracked anyway
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset server state
			server.documents = make(map[string]*Document)

			// Setup document if needed
			if tt.setupDoc {
				_, err := server.AddDocument(testURI)
				if err != nil {
					t.Fatalf("Failed to setup document: %v", err)
				}
			}

			// Create and handle file event
			event := lsp.FileEvent{
				URI:  testURI,
				Type: tt.changeType,
			}

			server.handleFileSystemChange(event)

			// Check document state
			hasDoc := server.HasDocument(testURI)
			if hasDoc != tt.expectDoc {
				t.Errorf("Expected document existence %v, got %v", tt.expectDoc, hasDoc)
			}
		})
	}
}

func TestHandleExternalFileCreated(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	testURI := "file:///test/newfile.md"

	// Handle external file creation
	server.handleExternalFileCreated(testURI)

	// External creation should not automatically add document to store
	if server.HasDocument(testURI) {
		t.Error("External file creation should not automatically add document to store")
	}
}

func TestHandleExternalFileChanged(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	testURI := "file:///test/document.md"

	// Test with tracked file
	_, err := server.AddDocument(testURI)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}

	server.handleExternalFileChanged(testURI)

	// Document should still exist
	if !server.HasDocument(testURI) {
		t.Error("Document should still exist after external change")
	}

	// Test with untracked file
	untrackedURI := "file:///test/untracked.md"
	server.handleExternalFileChanged(untrackedURI)

	// Untracked file should remain untracked
	if server.HasDocument(untrackedURI) {
		t.Error("Untracked file should remain untracked after external change")
	}
}

func TestHandleExternalFileDeleted(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	testURI := "file:///test/document.md"

	// Test with tracked file
	_, err := server.AddDocument(testURI)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}

	server.handleExternalFileDeleted(testURI)

	// Document should be removed from store
	if server.HasDocument(testURI) {
		t.Error("Document should be removed after external deletion")
	}

	// Test with untracked file
	untrackedURI := "file:///test/untracked.md"
	server.handleExternalFileDeleted(untrackedURI)

	// Should not cause any issues (no-op for untracked files)
	if server.HasDocument(untrackedURI) {
		t.Error("Untracked file should not be added to store during deletion handling")
	}
}

func TestRegisterFileWatcher(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Mock client register function
	var registeredParams lsp.RegistrationParams
	mockRegister := func(params lsp.RegistrationParams) error {
		registeredParams = params
		return nil
	}

	err := server.RegisterFileWatcher(mockRegister)
	if err != nil {
		t.Errorf("RegisterFileWatcher failed: %v", err)
	}

	// Verify registration was called with correct parameters
	if len(registeredParams.Registrations) != 1 {
		t.Errorf("Expected 1 registration, got %d", len(registeredParams.Registrations))
	}

	registration := registeredParams.Registrations[0]
	if registration.ID != "notedown-file-watcher" {
		t.Errorf("Expected registration ID 'notedown-file-watcher', got %s", registration.ID)
	}

	if registration.Method != "workspace/didChangeWatchedFiles" {
		t.Errorf("Expected method 'workspace/didChangeWatchedFiles', got %s", registration.Method)
	}

	// Verify register options
	options, ok := registration.RegisterOptions.(lsp.DidChangeWatchedFilesRegistrationOptions)
	if !ok {
		t.Error("RegisterOptions should be of type DidChangeWatchedFilesRegistrationOptions")
	}

	if len(options.Watchers) != 1 {
		t.Errorf("Expected 1 watcher, got %d", len(options.Watchers))
	}

	// Check that we're watching .md files only (Markdown-focused)
	patterns := make([]string, len(options.Watchers))
	for i, watcher := range options.Watchers {
		patterns[i] = watcher.GlobPattern
	}

	expectedPatterns := []string{"**/*.md"}
	for _, expected := range expectedPatterns {
		found := slices.Contains(patterns, expected)
		if !found {
			t.Errorf("Expected pattern %s not found in watchers", expected)
		}
	}
}
