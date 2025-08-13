package notedownls

import (
	"encoding/json"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
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
