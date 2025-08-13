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
