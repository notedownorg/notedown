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

package indexes

import (
	"testing"
	"time"

	"github.com/notedownorg/notedown/pkg/log"
)

// mockFile implements the WorkspaceFile interface for testing
type mockFile struct {
	uri  string
	path string
}

func (m *mockFile) GetURI() string  { return m.uri }
func (m *mockFile) GetPath() string { return m.path }

func TestWikilinkIndex_ConflictDetection(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Create workspace files with conflicting basenames
	workspaceFiles := map[string]WorkspaceFile{
		"api/user-endpoints.md":  &mockFile{uri: "file:///api/user-endpoints.md", path: "api/user-endpoints.md"},
		"docs/user-endpoints.md": &mockFile{uri: "file:///docs/user-endpoints.md", path: "docs/user-endpoints.md"},
		"simple-file.md":         &mockFile{uri: "file:///simple-file.md", path: "simple-file.md"},
	}

	// Test conflicting target
	target := "user-endpoints"
	sourceURI := "file:///test.md"

	// Add target with conflict detection
	exists, matchingFiles := index.targetExistsInWorkspace(target, workspaceFiles)
	if !exists {
		t.Errorf("Expected target '%s' to exist", target)
	}

	expectedMatches := 2 // Should match both api/user-endpoints.md and docs/user-endpoints.md
	if len(matchingFiles) != expectedMatches {
		t.Errorf("Expected %d matching files, got %d: %v", expectedMatches, len(matchingFiles), matchingFiles)
	}

	// Add to index
	index.AddTargetWithMatches(target, sourceURI, exists, matchingFiles)

	// Verify target info shows ambiguity
	allTargets := index.GetAllTargets()
	targetInfo, found := allTargets[target]
	if !found {
		t.Fatalf("Target '%s' not found in index", target)
	}

	if !targetInfo.IsAmbiguous {
		t.Errorf("Expected target '%s' to be marked as ambiguous", target)
	}

	if len(targetInfo.MatchingFiles) != expectedMatches {
		t.Errorf("Expected %d matching files in target info, got %d", expectedMatches, len(targetInfo.MatchingFiles))
	}

	// Test non-conflicting target
	target2 := "simple-file"
	exists2, matchingFiles2 := index.targetExistsInWorkspace(target2, workspaceFiles)
	if !exists2 {
		t.Errorf("Expected target '%s' to exist", target2)
	}

	if len(matchingFiles2) != 1 {
		t.Errorf("Expected 1 matching file for non-conflicting target, got %d", len(matchingFiles2))
	}

	index.AddTargetWithMatches(target2, sourceURI, exists2, matchingFiles2)

	// Verify non-conflicting target is not marked as ambiguous
	allTargets2 := index.GetAllTargets()
	targetInfo2, found2 := allTargets2[target2]
	if !found2 {
		t.Fatalf("Target '%s' not found in index", target2)
	}

	if targetInfo2.IsAmbiguous {
		t.Errorf("Expected target '%s' to NOT be marked as ambiguous", target2)
	}
}

func TestWikilinkIndex_GetAmbiguousTargets(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Add an ambiguous target
	index.AddTargetWithMatches("config", "file:///test.md", true, []string{"api/config.md", "docs/config.md"})

	// Add a non-ambiguous target
	index.AddTargetWithMatches("readme", "file:///test.md", true, []string{"readme.md"})

	// Add an unreferenced ambiguous target (should not be returned)
	index.targets["unreferenced"] = &WikilinkTargetInfo{
		Target:        "unreferenced",
		Exists:        true,
		ReferencedBy:  make(map[string]bool), // Empty - no references
		LastSeen:      time.Now(),
		MatchingFiles: []string{"api/unreferenced.md", "docs/unreferenced.md"},
		IsAmbiguous:   true,
	}

	// Get ambiguous targets
	ambiguousTargets := index.GetAmbiguousTargets()

	// Should only return the "config" target (has references and is ambiguous)
	if len(ambiguousTargets) != 1 {
		t.Errorf("Expected 1 ambiguous target, got %d", len(ambiguousTargets))
	}

	configTarget, found := ambiguousTargets["config"]
	if !found {
		t.Errorf("Expected 'config' target to be in ambiguous targets")
	}

	if !configTarget.IsAmbiguous {
		t.Errorf("Expected 'config' target to be marked as ambiguous")
	}

	if len(configTarget.MatchingFiles) != 2 {
		t.Errorf("Expected 2 matching files for 'config', got %d", len(configTarget.MatchingFiles))
	}

	// Verify unreferenced target is not included
	_, found = ambiguousTargets["unreferenced"]
	if found {
		t.Errorf("Expected 'unreferenced' target to NOT be in ambiguous targets (no references)")
	}
}

func TestWikilinkIndex_ExtractWikilinksFromDocument_WithConflicts(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Create workspace files with conflicts
	workspaceFiles := map[string]WorkspaceFile{
		"api/config.md":  &mockFile{uri: "file:///api/config.md", path: "api/config.md"},
		"docs/config.md": &mockFile{uri: "file:///docs/config.md", path: "docs/config.md"},
	}

	// Document content with wikilink to conflicting target
	content := `# Test Document

This references [[config]] which should create a conflict.

Also references [[api/config]] which should be specific.
`

	documentURI := "file:///test.md"

	// Extract wikilinks
	targets := index.ExtractWikilinksFromDocument(content, documentURI, workspaceFiles)

	// Should extract both targets
	expectedTargets := []string{"config", "api/config"}
	if len(targets) != len(expectedTargets) {
		t.Errorf("Expected %d targets, got %d: %v", len(expectedTargets), len(targets), targets)
	}

	// Check that "config" is marked as ambiguous
	allTargets := index.GetAllTargets()
	configTarget, found := allTargets["config"]
	if !found {
		t.Fatalf("Target 'config' not found in index")
	}

	if !configTarget.IsAmbiguous {
		t.Errorf("Expected 'config' target to be marked as ambiguous")
	}

	// Check that "api/config" is not ambiguous (exact path match)
	apiConfigTarget, found := allTargets["api/config"]
	if !found {
		t.Fatalf("Target 'api/config' not found in index")
	}

	if apiConfigTarget.IsAmbiguous {
		t.Errorf("Expected 'api/config' target to NOT be marked as ambiguous")
	}
}
