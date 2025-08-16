package notedownls

import (
	"encoding/json"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/lsp/pkg/notedownls/indexes"
	"github.com/notedownorg/notedown/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleCodeAction(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Mock workspace with ambiguous files
	server.workspace.roots = []WorkspaceRoot{
		{URI: "file:///test", Path: "/test", Name: "test"},
	}

	// Add files with same base name to create ambiguity
	server.workspace.fileIndex = map[string]*FileInfo{
		"file:///test/config.md": {
			URI:  "file:///test/config.md",
			Path: "config.md",
		},
		"file:///test/docs/config.md": {
			URI:  "file:///test/docs/config.md",
			Path: "docs/config.md",
		},
		"file:///test/project/config.md": {
			URI:  "file:///test/project/config.md",
			Path: "project/config.md",
		},
	}

	// Create test document with ambiguous wikilink
	testDoc := &Document{
		URI:     "file:///test/main.md",
		Content: "This is a [[config]] link that is ambiguous.",
	}
	server.documents["file:///test/main.md"] = testDoc

	// Set up wikilink index with ambiguous target
	workspaceFiles := map[string]indexes.WorkspaceFile{
		"file:///test/config.md": &FileInfo{
			URI:  "file:///test/config.md",
			Path: "config.md",
		},
		"file:///test/docs/config.md": &FileInfo{
			URI:  "file:///test/docs/config.md",
			Path: "docs/config.md",
		},
		"file:///test/project/config.md": &FileInfo{
			URI:  "file:///test/project/config.md",
			Path: "project/config.md",
		},
	}
	server.wikilinkIndex.ExtractWikilinksFromDocument(testDoc.Content, testDoc.URI, workspaceFiles)

	tests := []struct {
		name             string
		params           lsp.CodeActionParams
		expectActions    bool
		expectedCount    int
		checkFirstAction func(t *testing.T, action lsp.CodeAction)
	}{
		{
			name: "code actions for ambiguous wikilink",
			params: lsp.CodeActionParams{
				TextDocument: lsp.TextDocumentIdentifier{
					URI: "file:///test/main.md",
				},
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 10},
					End:   lsp.Position{Line: 0, Character: 20},
				},
				Context: lsp.CodeActionContext{
					Diagnostics: []lsp.Diagnostic{
						{
							Range: lsp.Range{
								Start: lsp.Position{Line: 0, Character: 10},
								End:   lsp.Position{Line: 0, Character: 20},
							},
							Message: "Ambiguous wikilink 'config' matches multiple files: config.md, docs/config.md, project/config.md",
							Code:    "ambiguous-wikilink",
						},
					},
				},
			},
			expectActions: true,
			expectedCount: 3, // One for each matching file
			checkFirstAction: func(t *testing.T, action lsp.CodeAction) {
				assert.Contains(t, action.Title, "Link to")
				assert.NotNil(t, action.Kind)
				assert.Equal(t, lsp.CodeActionKindQuickFix, *action.Kind)
				assert.NotNil(t, action.Edit)
				assert.Len(t, action.Edit.Changes, 1)

				// Check that the text edit contains a qualified path and display text
				changes := action.Edit.Changes["file:///test/main.md"]
				require.Len(t, changes, 1)

				newText := changes[0].NewText
				// Should be in format [[qualified-path|config]]
				assert.Contains(t, newText, "|config]]")
				assert.True(t,
					newText == "[[./config|config]]" ||
						newText == "[[docs/config|config]]" ||
						newText == "[[project/config|config]]",
					"Expected qualified wikilink format, got: %s", newText)
			},
		},
		{
			name: "no code actions for non-ambiguous range",
			params: lsp.CodeActionParams{
				TextDocument: lsp.TextDocumentIdentifier{
					URI: "file:///test/main.md",
				},
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 0},
					End:   lsp.Position{Line: 0, Character: 5},
				},
				Context: lsp.CodeActionContext{
					Diagnostics: []lsp.Diagnostic{}, // No ambiguous wikilink diagnostics
				},
			},
			expectActions: false,
			expectedCount: 0,
		},
		{
			name: "no code actions for non-existent document",
			params: lsp.CodeActionParams{
				TextDocument: lsp.TextDocumentIdentifier{
					URI: "file:///test/nonexistent.md",
				},
				Range: lsp.Range{
					Start: lsp.Position{Line: 0, Character: 0},
					End:   lsp.Position{Line: 0, Character: 5},
				},
				Context: lsp.CodeActionContext{
					Diagnostics: []lsp.Diagnostic{},
				},
			},
			expectActions: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paramsJSON, err := json.Marshal(tt.params)
			require.NoError(t, err)

			result, err := server.handleCodeAction(paramsJSON)
			require.NoError(t, err)

			actions, ok := result.([]lsp.CodeAction)
			require.True(t, ok, "Expected []lsp.CodeAction, got %T", result)

			if tt.expectActions {
				assert.Len(t, actions, tt.expectedCount)
				if len(actions) > 0 && tt.checkFirstAction != nil {
					tt.checkFirstAction(t, actions[0])
				}
			} else {
				assert.Len(t, actions, 0)
			}
		})
	}
}

func TestGenerateQualifiedPath(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	tests := []struct {
		name     string
		fileInfo *FileInfo
		expected string
	}{
		{
			name: "root file gets ./ prefix",
			fileInfo: &FileInfo{
				Path: "config.md",
			},
			expected: "./config",
		},
		{
			name: "subdirectory file keeps relative path",
			fileInfo: &FileInfo{
				Path: "docs/config.md",
			},
			expected: "docs/config",
		},
		{
			name: "nested subdirectory file",
			fileInfo: &FileInfo{
				Path: "project/docs/api.md",
			},
			expected: "project/docs/api",
		},
		{
			name: "windows-style paths converted",
			fileInfo: &FileInfo{
				Path: "docs\\config.md",
			},
			expected: "docs/config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.generateQualifiedPath(tt.fileInfo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildWikilinkWithDisplayText(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	tests := []struct {
		name           string
		qualifiedPath  string
		originalTarget string
		expected       string
	}{
		{
			name:           "root file with display text",
			qualifiedPath:  "./config",
			originalTarget: "config",
			expected:       "[[./config|config]]",
		},
		{
			name:           "subdirectory file with display text",
			qualifiedPath:  "docs/config",
			originalTarget: "config",
			expected:       "[[docs/config|config]]",
		},
		{
			name:           "nested path with display text",
			qualifiedPath:  "project/docs/api",
			originalTarget: "api",
			expected:       "[[project/docs/api|api]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.buildWikilinkWithDisplayText(tt.qualifiedPath, tt.originalTarget)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractTargetFromDiagnostic(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "extracts target from ambiguous wikilink message",
			message:  "Ambiguous wikilink 'config' matches multiple files: config.md, docs/config.md",
			expected: "config",
		},
		{
			name:     "extracts complex target name",
			message:  "Ambiguous wikilink 'project-alpha' matches multiple files: project-alpha.md, docs/project-alpha.md",
			expected: "project-alpha",
		},
		{
			name:     "returns empty for malformed message",
			message:  "Some other diagnostic message",
			expected: "",
		},
		{
			name:     "handles message without quotes",
			message:  "Ambiguous wikilink config matches multiple files",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.extractTargetFromDiagnostic(tt.message)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRangesOverlap(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	tests := []struct {
		name     string
		range1   lsp.Range
		range2   lsp.Range
		expected bool
	}{
		{
			name: "ranges overlap on same line",
			range1: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 5},
				End:   lsp.Position{Line: 0, Character: 15},
			},
			range2: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 10},
				End:   lsp.Position{Line: 0, Character: 20},
			},
			expected: true,
		},
		{
			name: "ranges don't overlap on same line",
			range1: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 5},
				End:   lsp.Position{Line: 0, Character: 10},
			},
			range2: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 15},
				End:   lsp.Position{Line: 0, Character: 20},
			},
			expected: false,
		},
		{
			name: "ranges on different lines don't overlap",
			range1: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 5},
				End:   lsp.Position{Line: 0, Character: 10},
			},
			range2: lsp.Range{
				Start: lsp.Position{Line: 1, Character: 5},
				End:   lsp.Position{Line: 1, Character: 10},
			},
			expected: false,
		},
		{
			name: "multiline ranges overlap",
			range1: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 5},
				End:   lsp.Position{Line: 2, Character: 10},
			},
			range2: lsp.Range{
				Start: lsp.Position{Line: 1, Character: 5},
				End:   lsp.Position{Line: 1, Character: 10},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.rangesOverlap(tt.range1, tt.range2)
			assert.Equal(t, tt.expected, result)
		})
	}
}
