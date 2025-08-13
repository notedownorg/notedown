package notedownls

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestCacheDebugLogging(t *testing.T) {
	// Create a debug logger to see cache operations
	logger := log.New(os.Stderr, log.Debug)
	wm := NewWorkspaceManager(logger)

	// Create a temporary workspace
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.md")
	if err := os.WriteFile(testFile, []byte("# Test Document"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Initialize workspace
	params := lsp.InitializeParams{
		RootUri: &[]string{pathToFileURI(tempDir)}[0],
	}
	err := wm.InitializeFromParams(params)
	if err != nil {
		t.Fatalf("Failed to initialize workspace: %v", err)
	}

	t.Log("=== Starting workspace discovery (should show 'added file to workspace cache') ===")

	// Discover files - should show debug log when file is added to cache
	err = wm.DiscoverMarkdownFiles()
	if err != nil {
		t.Fatalf("Failed to discover files: %v", err)
	}

	t.Log("=== Adding external file to cache (should show debug log) ===")

	// Test adding a single file
	testURI := pathToFileURI(testFile)
	err = wm.AddFileToIndex(testURI)
	if err != nil {
		t.Errorf("Failed to add file to index: %v", err)
	}

	t.Log("=== Removing file from cache (should show debug log) ===")

	// Test removing a file
	wm.RemoveFileFromIndex(testURI)

	t.Log("=== Running rediscovery (should show 'cleared workspace cache') ===")

	// Test rediscovery - should show cache clearing debug log
	err = wm.DiscoverMarkdownFiles()
	if err != nil {
		t.Fatalf("Failed to rediscover files: %v", err)
	}

	t.Log("=== Cache operations test completed ===")
}
