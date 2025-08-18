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

	"github.com/notedownorg/notedown/pkg/log"
)

// MockWorkspaceFile implements the WorkspaceFile interface for testing
type MockWorkspaceFile struct {
	URI  string
	Path string
}

func (m *MockWorkspaceFile) GetURI() string {
	return m.URI
}

func (m *MockWorkspaceFile) GetPath() string {
	return m.Path
}

func TestWikilinkIndex_AddTarget(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Test adding an existing target
	index.AddTarget("project-alpha", "file:///test1.md", true)

	targets := index.GetAllTargets()
	if len(targets) != 1 {
		t.Errorf("expected 1 target, got %d", len(targets))
	}

	target := targets["project-alpha"]
	if target == nil {
		t.Fatal("target 'project-alpha' not found")
	}

	if !target.Exists {
		t.Error("target should exist")
	}

	if len(target.ReferencedBy) != 1 {
		t.Errorf("expected 1 reference, got %d", len(target.ReferencedBy))
	}

	if !target.ReferencedBy["file:///test1.md"] {
		t.Error("reference to test1.md not found")
	}
}

func TestWikilinkIndex_NonExistentTargets(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Add mix of existing and non-existing targets
	index.AddTarget("existing-file", "file:///test1.md", true)
	index.AddTarget("non-existent", "file:///test1.md", false)
	index.AddTarget("also-missing", "file:///test2.md", false)

	// Test getting non-existent targets
	nonExistent := index.GetNonExistentTargets()
	if len(nonExistent) != 2 {
		t.Errorf("expected 2 non-existent targets, got %d", len(nonExistent))
	}

	if _, found := nonExistent["non-existent"]; !found {
		t.Error("'non-existent' target should be in non-existent list")
	}

	if _, found := nonExistent["also-missing"]; !found {
		t.Error("'also-missing' target should be in non-existent list")
	}

	if _, found := nonExistent["existing-file"]; found {
		t.Error("'existing-file' should not be in non-existent list")
	}
}

func TestWikilinkIndex_ExtractWikilinksFromDocument(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Create mock workspace files
	workspaceFiles := map[string]WorkspaceFile{
		"file:///existing.md": &MockWorkspaceFile{
			URI:  "file:///existing.md",
			Path: "existing.md",
		},
	}

	content := `# Test Document

This document contains [[existing]] and [[non-existent]] wikilinks.
Also has [[docs/new-guide]] and [[existing|display text]].`

	targets := index.ExtractWikilinksFromDocument(content, "file:///test.md", workspaceFiles)

	expectedTargets := []string{"existing", "non-existent", "docs/new-guide", "existing"}
	if len(targets) != len(expectedTargets) {
		t.Errorf("expected %d targets, got %d", len(expectedTargets), len(targets))
	}

	// Check that targets were added to index
	allTargets := index.GetAllTargets()

	// Should have 3 unique targets (existing appears twice)
	expectedUniqueTargets := map[string]bool{
		"existing":       true,
		"non-existent":   false,
		"docs/new-guide": false,
	}

	for target, shouldExist := range expectedUniqueTargets {
		targetInfo, found := allTargets[target]
		if !found {
			t.Errorf("target '%s' not found in index", target)
			continue
		}

		if targetInfo.Exists != shouldExist {
			t.Errorf("target '%s' existence: expected %v, got %v", target, shouldExist, targetInfo.Exists)
		}

		if !targetInfo.ReferencedBy["file:///test.md"] {
			t.Errorf("target '%s' should be referenced by test.md", target)
		}
	}
}

func TestWikilinkIndex_RemoveTargetReference(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Add a target referenced by two documents
	index.AddTarget("shared-target", "file:///doc1.md", false)
	index.AddTarget("shared-target", "file:///doc2.md", false)

	// Add a target referenced by only one document
	index.AddTarget("single-target", "file:///doc1.md", false)

	// Remove reference from doc1
	index.RemoveTargetReference("shared-target", "file:///doc1.md")
	index.RemoveTargetReference("single-target", "file:///doc1.md")

	allTargets := index.GetAllTargets()

	// shared-target should still exist (referenced by doc2)
	if sharedTarget, found := allTargets["shared-target"]; found {
		if len(sharedTarget.ReferencedBy) != 1 {
			t.Errorf("shared-target should have 1 reference, got %d", len(sharedTarget.ReferencedBy))
		}
		if !sharedTarget.ReferencedBy["file:///doc2.md"] {
			t.Error("shared-target should still be referenced by doc2.md")
		}
	} else {
		t.Error("shared-target should still exist")
	}

	// single-target should be removed (no references and doesn't exist)
	if _, found := allTargets["single-target"]; found {
		t.Error("single-target should have been removed")
	}
}

func TestWikilinkIndex_GetTargetsByPrefix(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	// Add various targets
	index.AddTarget("project-alpha", "file:///test.md", true)
	index.AddTarget("project-beta", "file:///test.md", false)
	index.AddTarget("docs/api", "file:///test.md", false)
	index.AddTarget("notes/ideas", "file:///test.md", false)

	// Test prefix matching
	tests := []struct {
		prefix   string
		expected []string
	}{
		{"proj", []string{"project-alpha", "project-beta"}},
		{"project-", []string{"project-alpha", "project-beta"}},
		{"docs", []string{"docs/api"}},
		{"", []string{"project-alpha", "project-beta", "docs/api", "notes/ideas"}}, // empty prefix returns all
		{"xyz", []string{}}, // no matches
	}

	for _, tt := range tests {
		t.Run("prefix_"+tt.prefix, func(t *testing.T) {
			results := index.GetTargetsByPrefix(tt.prefix)

			if len(results) != len(tt.expected) {
				t.Errorf("prefix '%s': expected %d results, got %d", tt.prefix, len(tt.expected), len(results))
			}

			for _, expectedTarget := range tt.expected {
				if _, found := results[expectedTarget]; !found {
					t.Errorf("prefix '%s': expected target '%s' not found", tt.prefix, expectedTarget)
				}
			}
		})
	}
}

func TestWikilinkIndex_RefreshDocumentWikilinks(t *testing.T) {
	logger := log.NewDefault()
	index := NewWikilinkIndex(logger)

	workspaceFiles := map[string]WorkspaceFile{
		"file:///existing.md": &MockWorkspaceFile{
			URI:  "file:///existing.md",
			Path: "existing.md",
		},
	}

	// Initial content with some wikilinks
	initialContent := `[[old-link]] and [[shared-link]]`
	index.ExtractWikilinksFromDocument(initialContent, "file:///test.md", workspaceFiles)

	// Verify initial state
	allTargets := index.GetAllTargets()
	if len(allTargets) != 2 {
		t.Errorf("expected 2 initial targets, got %d", len(allTargets))
	}

	// Update content with different wikilinks
	newContent := `[[new-link]] and [[shared-link]]`
	index.RefreshDocumentWikilinks(newContent, "file:///test.md", workspaceFiles)

	// Verify updated state
	allTargets = index.GetAllTargets()

	// Should have new-link and shared-link, but not old-link
	if _, found := allTargets["new-link"]; !found {
		t.Error("new-link should be present after refresh")
	}

	if _, found := allTargets["shared-link"]; !found {
		t.Error("shared-link should still be present after refresh")
	}

	if _, found := allTargets["old-link"]; found {
		t.Error("old-link should be removed after refresh")
	}
}
