package notedownls

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestCacheClearDebugLogging(t *testing.T) {
	// Create a debug logger to see cache operations
	logger := log.New(os.Stderr, log.Debug)
	wm := NewWorkspaceManager(logger)

	// Create a temporary workspace with multiple files
	tempDir := t.TempDir()
	testFiles := []string{"file1.md", "file2.md", "file3.md"}

	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(filePath, []byte("# "+fileName), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fileName, err)
		}
	}

	// Initialize workspace
	params := lsp.InitializeParams{
		RootUri: &[]string{pathToFileURI(tempDir)}[0],
	}
	err := wm.InitializeFromParams(params)
	if err != nil {
		t.Fatalf("Failed to initialize workspace: %v", err)
	}

	t.Log("=== Initial discovery (populating cache) ===")

	// Initial discovery
	err = wm.DiscoverMarkdownFiles()
	if err != nil {
		t.Fatalf("Failed to discover files: %v", err)
	}

	// Verify files are in cache
	files := wm.GetMarkdownFiles()
	if len(files) != len(testFiles) {
		t.Errorf("Expected %d files in cache, got %d", len(testFiles), len(files))
	}

	t.Log("=== Rediscovery (should show 'cleared workspace cache for rediscovery') ===")

	// Rediscovery - should show cache clearing debug log
	err = wm.DiscoverMarkdownFiles()
	if err != nil {
		t.Fatalf("Failed to rediscover files: %v", err)
	}

	t.Log("=== Cache clear test completed ===")
}
