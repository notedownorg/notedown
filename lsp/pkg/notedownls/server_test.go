package notedownls

import (
	"testing"

	"github.com/notedownorg/notedown/pkg/log"
)

func TestServerDocumentOperations(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	testURI := "file:///test/document.md"

	// Initially should not have document
	if server.HasDocument(testURI) {
		t.Error("Server should not have document initially")
	}

	// Add document
	doc, err := server.AddDocument(testURI)
	if err != nil {
		t.Errorf("Failed to add document: %v", err)
	}

	if doc.URI != testURI {
		t.Errorf("Expected URI %s, got %s", testURI, doc.URI)
	}

	// Should now have document
	if !server.HasDocument(testURI) {
		t.Error("Server should have document after adding")
	}

	// Get document
	retrievedDoc, exists := server.GetDocument(testURI)
	if !exists {
		t.Error("Document should exist")
	}

	if retrievedDoc.URI != testURI {
		t.Errorf("Expected URI %s, got %s", testURI, retrievedDoc.URI)
	}

	// Remove document
	server.RemoveDocument(testURI)

	// Should no longer have document
	if server.HasDocument(testURI) {
		t.Error("Server should not have document after removal")
	}

	// Get should return false
	_, exists = server.GetDocument(testURI)
	if exists {
		t.Error("Document should not exist after removal")
	}
}
