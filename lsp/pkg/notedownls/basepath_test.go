package notedownls

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notedownorg/notedown/pkg/log"
)

func TestBasepathCorrectness(t *testing.T) {
	// Create a debug logger to see the corrected basepath logging
	logger := log.New(os.Stderr, log.Debug)
	server := NewServer("test", logger)

	// Create a temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "my-awesome-document.md")
	if err := os.WriteFile(testFile, []byte("# My Awesome Document"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testURI := pathToFileURI(testFile)

	t.Log("=== Adding document (should show basepath='my-awesome-document') ===")

	// Add document - should now show basepath as just the filename without extension
	_, err := server.AddDocument(testURI)
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}

	// Verify the basepath is correct
	doc, exists := server.GetDocument(testURI)
	if !exists {
		t.Fatal("Document should exist")
	}

	expectedBasepath := "my-awesome-document"
	if doc.Basepath != expectedBasepath {
		t.Errorf("Expected basepath '%s', got '%s'", expectedBasepath, doc.Basepath)
	}

	t.Logf("✅ Basepath correctly set to: '%s'", doc.Basepath)
}
