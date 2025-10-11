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
	"testing"

	"github.com/notedownorg/notedown/language-server/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestHandleDidOpen(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Create test parameters
	params := lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        "file:///test/document.md",
			LanguageID: "markdown",
			Version:    1,
			Text:       "# Test Document",
		},
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	// Test handler
	err = server.handleDidOpen(paramBytes)
	if err != nil {
		t.Errorf("handleDidOpen failed: %v", err)
	}

	// Verify document was added
	if !server.HasDocument(params.TextDocument.URI) {
		t.Error("Document should be added after handleDidOpen")
	}
}

func TestHandleDidClose(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	testURI := "file:///test/document.md"

	// Add document first
	_, err := server.AddDocument(testURI)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}

	// Create test parameters
	params := lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{
			URI: testURI,
		},
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	// Test handler
	err = server.handleDidClose(paramBytes)
	if err != nil {
		t.Errorf("handleDidClose failed: %v", err)
	}

	// Verify document was removed
	if server.HasDocument(testURI) {
		t.Error("Document should be removed after handleDidClose")
	}
}

func TestHandleDidChange(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	testURI := "file:///test/document.md"

	// Add document first
	_, err := server.AddDocument(testURI)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}

	// Create test parameters
	params := lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{
				URI: testURI,
			},
			Version: &[]int{2}[0],
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{
			{
				Text: "# Updated Document",
			},
		},
	}

	paramBytes, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	// Test handler
	err = server.handleDidChange(paramBytes)
	if err != nil {
		t.Errorf("handleDidChange failed: %v", err)
	}

	// Document should still exist (didChange doesn't remove it)
	if !server.HasDocument(testURI) {
		t.Error("Document should still exist after handleDidChange")
	}
}

func TestGetCompleteWikilinkTarget(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	tests := []struct {
		name      string
		content   string
		line      int
		character int
		expected  string
	}{
		{
			name:      "simple wikilink - cursor at start",
			content:   "Here is a [[meeting-notes]] link",
			line:      0,
			character: 11, // At the 'm' in meeting-notes
			expected:  "meeting-notes",
		},
		{
			name:      "simple wikilink - cursor in middle",
			content:   "Here is a [[meeting-notes]] link",
			line:      0,
			character: 18, // At the '-' in meeting-notes
			expected:  "meeting-notes",
		},
		{
			name:      "simple wikilink - cursor at end",
			content:   "Here is a [[meeting-notes]] link",
			line:      0,
			character: 24, // At the ']'
			expected:  "meeting-notes",
		},
		{
			name:      "wikilink with pipe separator",
			content:   "Check out [[project-alpha|Alpha Project]] for details",
			line:      0,
			character: 15, // In the target part
			expected:  "project-alpha",
		},
		{
			name:      "wikilink with pipe separator - cursor in display text",
			content:   "Check out [[project-alpha|Alpha Project]] for details",
			line:      0,
			character: 30,              // In the display part
			expected:  "project-alpha", // Should still return the target part
		},
		{
			name:      "path-based wikilink",
			content:   "See [[docs/api-reference]] for more info",
			line:      0,
			character: 10, // In the target
			expected:  "docs/api-reference",
		},
		{
			name:      "wikilink with spaces",
			content:   "Link to [[ spaced target ]] here",
			line:      0,
			character: 15, // In the target
			expected:  "spaced target",
		},
		{
			name:      "multiple wikilinks - first one",
			content:   "Links: [[first-link]] and [[second-link]]",
			line:      0,
			character: 12, // In first link
			expected:  "first-link",
		},
		{
			name:      "multiple wikilinks - second one",
			content:   "Links: [[first-link]] and [[second-link]]",
			line:      0,
			character: 32, // In second link
			expected:  "second-link",
		},
		{
			name:      "not in wikilink - before",
			content:   "Some text [[wikilink]] here",
			line:      0,
			character: 5, // Before the wikilink
			expected:  "",
		},
		{
			name:      "not in wikilink - after",
			content:   "Some text [[wikilink]] here",
			line:      0,
			character: 25, // After the wikilink
			expected:  "",
		},
		{
			name:      "incomplete wikilink - no closing",
			content:   "Incomplete [[wikilink without closing",
			line:      0,
			character: 15, // In the incomplete link
			expected:  "",
		},
		{
			name:      "empty wikilink",
			content:   "Empty [[]] link",
			line:      0,
			character: 8, // In empty link
			expected:  "",
		},
		{
			name:      "multiline content - wikilink on second line",
			content:   "First line\nSecond line with [[target]] here",
			line:      1,
			character: 20, // In the target on second line
			expected:  "target",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a document with the test content
			doc := &Document{
				URI:     "file:///test.md",
				Content: tt.content,
				Version: 1,
			}

			// Test the method
			result := server.getCompleteWikilinkTarget(doc, lsp.Position{
				Line:      tt.line,
				Character: tt.character,
			})

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Definition provider tests moved to handlers_textDocument_definition_test.go

func TestHandleDefinitionDeprecated(t *testing.T) {
	tests := []struct {
		name            string
		documentURI     string
		documentContent string
		position        lsp.Position
		workspaceFiles  map[string]*FileInfo
		expectLocation  bool
		expectError     bool
		expectedURI     string
		description     string
	}{
		{
			name:            "existing file - simple target",
			documentURI:     "file:///test/document.md",
			documentContent: "Link to [[existing-file]] here",
			position:        lsp.Position{Line: 0, Character: 15}, // Inside wikilink
			workspaceFiles: map[string]*FileInfo{
				"existing-file": {
					Path: "existing-file.md",
					URI:  "file:///test/existing-file.md",
				},
			},
			expectLocation: true,
			expectError:    false,
			expectedURI:    "file:///test/existing-file.md",
			description:    "Should find existing file and return its location",
		},
		{
			name:            "existing file - with extension",
			documentURI:     "file:///test/document.md",
			documentContent: "Link to [[existing-file.md]] here",
			position:        lsp.Position{Line: 0, Character: 15}, // Inside wikilink
			workspaceFiles: map[string]*FileInfo{
				"existing-file": {
					Path: "existing-file.md",
					URI:  "file:///test/existing-file.md",
				},
			},
			expectLocation: true,
			expectError:    false,
			expectedURI:    "file:///test/existing-file.md",
			description:    "Should find existing file even when target includes .md extension",
		},
		{
			name:            "existing file - path-based target",
			documentURI:     "file:///test/document.md",
			documentContent: "Link to [[docs/api-reference]] here",
			position:        lsp.Position{Line: 0, Character: 20}, // Inside wikilink
			workspaceFiles: map[string]*FileInfo{
				"docs/api-reference": {
					Path: "docs/api-reference.md",
					URI:  "file:///test/docs/api-reference.md",
				},
			},
			expectLocation: true,
			expectError:    false,
			expectedURI:    "file:///test/docs/api-reference.md",
			description:    "Should find existing file with path-based target",
		},
		{
			name:            "wikilink with pipe separator",
			documentURI:     "file:///test/document.md",
			documentContent: "Link to [[target-file|Display Text]] here",
			position:        lsp.Position{Line: 0, Character: 15}, // Inside target part
			workspaceFiles: map[string]*FileInfo{
				"target-file": {
					Path: "target-file.md",
					URI:  "file:///test/target-file.md",
				},
			},
			expectLocation: true,
			expectError:    false,
			expectedURI:    "file:///test/target-file.md",
			description:    "Should extract target from wikilink with display text",
		},
		{
			name:            "not in wikilink context",
			documentURI:     "file:///test/document.md",
			documentContent: "Regular text with [[wikilink]] here",
			position:        lsp.Position{Line: 0, Character: 5}, // Outside wikilink
			workspaceFiles:  map[string]*FileInfo{},
			expectLocation:  false,
			expectError:     false,
			description:     "Should return nil when cursor not in wikilink",
		},
		{
			name:            "empty target",
			documentURI:     "file:///test/document.md",
			documentContent: "Empty [[]] link",
			position:        lsp.Position{Line: 0, Character: 8}, // Inside empty brackets
			workspaceFiles:  map[string]*FileInfo{},
			expectLocation:  false,
			expectError:     false,
			description:     "Should return nil for empty wikilink target",
		},
		{
			name:            "document not found",
			documentURI:     "file:///test/nonexistent.md",
			documentContent: "", // Won't be used since document doesn't exist
			position:        lsp.Position{Line: 0, Character: 0},
			workspaceFiles:  map[string]*FileInfo{},
			expectLocation:  false,
			expectError:     false,
			description:     "Should return nil when document doesn't exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.NewDefault()
			server := NewServer("test", logger)

			// Set up workspace
			roots := []WorkspaceRoot{
				{URI: "file:///test", Path: "/test", Name: "test"},
			}
			server.workspace.SetupTestWorkspace(roots, tt.workspaceFiles)

			// Add the test document if it exists
			if tt.documentURI == "file:///test/document.md" {
				doc, err := server.AddDocument(tt.documentURI)
				if err != nil {
					t.Fatalf("Failed to add document: %v", err)
				}
				doc.Content = tt.documentContent
			}

			// Create definition request
			params := lsp.DefinitionParams{
				TextDocumentPositionParams: lsp.TextDocumentPositionParams{
					TextDocument: lsp.TextDocumentIdentifier{
						URI: tt.documentURI,
					},
					Position: tt.position,
				},
			}

			// Marshal params to json.RawMessage for the handler
			paramBytes, err := json.Marshal(params)
			if err != nil {
				t.Fatalf("Failed to marshal params: %v", err)
			}

			// Test the handler
			result, err := server.handleDefinition(paramBytes)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check location expectation
			if tt.expectLocation {
				if result == nil {
					t.Errorf("Expected location but got nil")
				} else {
					location := result.(lsp.Location)
					if location.URI != tt.expectedURI {
						t.Errorf("Expected URI %s, got %s", tt.expectedURI, location.URI)
					}
				}
			} else {
				if result != nil {
					t.Errorf("Expected nil but got location: %v", result)
				}
			}

			t.Logf("Test passed: %s", tt.description)
		})
	}
}

func TestFindFileForTargetDeprecated(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Set up test workspace
	roots := []WorkspaceRoot{
		{URI: "file:///test", Path: "/test", Name: "test"},
	}
	workspaceFiles := map[string]*FileInfo{
		"simple-file": {
			Path: "simple-file.md",
			URI:  "file:///test/simple-file.md",
		},
		"docs/api-guide": {
			Path: "docs/api-guide.md",
			URI:  "file:///test/docs/api-guide.md",
		},
		"project-alpha": {
			Path: "projects/project-alpha.md",
			URI:  "file:///test/projects/project-alpha.md",
		},
	}
	server.workspace.SetupTestWorkspace(roots, workspaceFiles)

	tests := []struct {
		name     string
		target   string
		expected *FileInfo
	}{
		{
			name:   "exact match",
			target: "simple-file",
			expected: &FileInfo{
				Path: "simple-file.md",
				URI:  "file:///test/simple-file.md",
			},
		},
		{
			name:   "exact match with path",
			target: "docs/api-guide",
			expected: &FileInfo{
				Path: "docs/api-guide.md",
				URI:  "file:///test/docs/api-guide.md",
			},
		},
		{
			name:   "target with .md extension",
			target: "simple-file.md",
			expected: &FileInfo{
				Path: "simple-file.md",
				URI:  "file:///test/simple-file.md",
			},
		},
		{
			name:     "non-existent target",
			target:   "non-existent",
			expected: nil,
		},
		{
			name:     "empty target",
			target:   "",
			expected: nil,
		},
		{
			name:     "directory traversal attempt - should be rejected",
			target:   "../parent-file",
			expected: nil,
		},
		{
			name:     "directory traversal in subdirectory - should be rejected",
			target:   "docs/../secret",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.findFileForTarget(tt.target)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("Expected %v, got nil", tt.expected)
				} else {
					if result.Path != tt.expected.Path || result.URI != tt.expected.URI {
						t.Errorf("Expected %v, got %v", tt.expected, result)
					}
				}
			}
		})
	}
}

func TestResolveTargetPathDeprecated(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Set up workspace root
	roots := []WorkspaceRoot{
		{URI: "file:///test/workspace", Path: "/test/workspace", Name: "test"},
	}
	server.workspace.SetupTestWorkspace(roots, nil)

	tests := []struct {
		name         string
		target       string
		expectedPath string
		expectedURI  string
	}{
		{
			name:         "simple target",
			target:       "simple-file",
			expectedPath: "/test/workspace/simple-file.md",
			expectedURI:  "file:///test/workspace/simple-file.md",
		},
		{
			name:         "target with extension",
			target:       "file.md",
			expectedPath: "/test/workspace/file.md.md",
			expectedURI:  "file:///test/workspace/file.md.md",
		},
		{
			name:         "path-based target",
			target:       "docs/api-reference",
			expectedPath: "/test/workspace/docs/api-reference.md",
			expectedURI:  "file:///test/workspace/docs/api-reference.md",
		},
		{
			name:         "nested path target",
			target:       "projects/alpha/readme",
			expectedPath: "/test/workspace/projects/alpha/readme.md",
			expectedURI:  "file:///test/workspace/projects/alpha/readme.md",
		},
		{
			name:         "directory traversal attempt - should be rejected",
			target:       "../parent-doc",
			expectedPath: "",
			expectedURI:  "",
		},
		{
			name:         "directory traversal in subdirectory - should be rejected",
			target:       "docs/../secret",
			expectedPath: "",
			expectedURI:  "",
		},
		{
			name:         "multiple directory traversal - should be rejected",
			target:       "../../outside-workspace",
			expectedPath: "",
			expectedURI:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, uri := server.resolveTargetPath(tt.target)

			if path != tt.expectedPath {
				t.Errorf("Expected path %s, got %s", tt.expectedPath, path)
			}
			if uri != tt.expectedURI {
				t.Errorf("Expected URI %s, got %s", tt.expectedURI, uri)
			}
		})
	}
}
