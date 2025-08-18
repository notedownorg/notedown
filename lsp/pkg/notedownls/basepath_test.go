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

	"github.com/notedownorg/notedown/pkg/log"
)

func TestBasepathCorrectness(t *testing.T) {
	// Create a debug logger to see the corrected basepath logging
	logger := log.New(os.Stderr, log.Debug)
	server := NewServer("test", logger)

	// Create a temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "my-awesome-document.md")
	if err := os.WriteFile(testFile, []byte("# My Awesome Document"), 0600); err != nil {
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
