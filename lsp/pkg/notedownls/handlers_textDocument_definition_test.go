package notedownls

import (
	"encoding/json"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestHandleDefinition(t *testing.T) {
	tests := []struct {
		name           string
		documentURI    string
		documentContent string
		position       lsp.Position
		workspaceFiles map[string]*FileInfo
		expectLocation bool
		expectError    bool
		expectedURI    string
		description    string
	}{
		{
			name:            "existing file - simple target",
			documentURI:     "file:///test/document.md",
			documentContent: "Link to [[existing-file]] here",
			position:        lsp.Position{Line: 0, Character: 15}, // Inside wikilink
			workspaceFiles: map[string]*FileInfo{
				"file:///test/existing-file.md": {
					Path: "/test/existing-file.md",
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
				"file:///test/existing-file.md": {
					Path: "/test/existing-file.md",
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
				"file:///test/docs/api-reference.md": {
					Path: "/test/docs/api-reference.md",
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
				"file:///test/target-file.md": {
					Path: "/test/target-file.md",
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

			// Set up workspace files and roots
			server.workspace.fileIndex = tt.workspaceFiles
			server.workspace.roots = []WorkspaceRoot{
				{URI: "file:///test", Path: "/test", Name: "test"},
			}

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

func TestFindFileForTarget(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Set up test workspace files  
	server.workspace.fileIndex = map[string]*FileInfo{
		"file:///test/simple-file.md": {
			Path: "/test/simple-file.md",
			URI:  "file:///test/simple-file.md",
		},
		"file:///test/docs/api-guide.md": {
			Path: "/test/docs/api-guide.md", 
			URI:  "file:///test/docs/api-guide.md",
		},
		"file:///test/projects/project-alpha.md": {
			Path: "/test/projects/project-alpha.md",
			URI:  "file:///test/projects/project-alpha.md",
		},
	}

	tests := []struct {
		name     string
		target   string
		expected *FileInfo
	}{
		{
			name:   "exact match",
			target: "simple-file",
			expected: &FileInfo{
				Path: "/test/simple-file.md",
				URI:  "file:///test/simple-file.md",
			},
		},
		{
			name:   "exact match with path",
			target: "docs/api-guide",
			expected: &FileInfo{
				Path: "/test/docs/api-guide.md",
				URI:  "file:///test/docs/api-guide.md",
			},
		},
		{
			name:   "target with .md extension",
			target: "simple-file.md",
			expected: &FileInfo{
				Path: "/test/simple-file.md",
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

func TestResolveTargetPath(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Set up workspace root
	server.workspace.roots = []WorkspaceRoot{
		{URI: "file:///test/workspace", Path: "/test/workspace", Name: "test"},
	}

	tests := []struct {
		name        string
		target      string
		expectedPath string
		expectedURI  string
	}{
		{
			name:        "simple target",
			target:      "simple-file",
			expectedPath: "/test/workspace/simple-file.md",
			expectedURI:  "file:///test/workspace/simple-file.md",
		},
		{
			name:        "target with extension - should double add extension (current behavior)",
			target:      "file.md",
			expectedPath: "/test/workspace/file.md.md",
			expectedURI:  "file:///test/workspace/file.md.md",
		},
		{
			name:        "path-based target",
			target:      "docs/api-reference",
			expectedPath: "/test/workspace/docs/api-reference.md",
			expectedURI:  "file:///test/workspace/docs/api-reference.md",
		},
		{
			name:        "nested path target",
			target:      "projects/alpha/readme",
			expectedPath: "/test/workspace/projects/alpha/readme.md",
			expectedURI:  "file:///test/workspace/projects/alpha/readme.md",
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