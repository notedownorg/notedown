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
