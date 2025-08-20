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
	"io"
	"testing"

	"github.com/notedownorg/notedown/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper to create a server instance
func createTestServerForBoundaries() *Server {
	logger := log.New(io.Discard, log.Error) // Suppress logs during tests
	return &Server{
		logger:    logger,
		documents: make(map[string]*Document),
	}
}

func TestHandleGetListItemBoundaries_SimpleList(t *testing.T) {
	server := createTestServerForBoundaries()

	// Test content - simple list
	content := `# Test List

- First item
- Second item
- Third item
- Fourth item

Some text after the list.`

	// Add document to server
	doc := &Document{
		URI:     "file:///test.md",
		Content: content,
	}
	server.documents[doc.URI] = doc

	// Test getting boundaries for second item (line 4, 0-based index 3)
	position := map[string]any{
		"line":      float64(3), // Second item is on line 3 (0-based)
		"character": float64(2), // Position in "Second item"
	}

	arguments := []any{
		"file:///test.md",
		position,
	}

	result, err := server.handleGetListItemBoundaries(arguments)
	require.NoError(t, err)

	// Assert the result is correct
	response, ok := result.(*BoundaryResponse)
	require.True(t, ok)
	assert.True(t, response.Found)
	assert.Equal(t, 3, response.Start.Line) // Line 3 (0-based)
	assert.Equal(t, 0, response.Start.Character)
	assert.Equal(t, 4, response.End.Line) // Next line (exclusive)
	assert.Equal(t, 0, response.End.Character)
}

func TestHandleGetListItemBoundaries_NestedList(t *testing.T) {
	server := createTestServerForBoundaries()

	// Test content - nested list
	content := `# Nested List

- Level 1 Item A
  - Level 2 Item A.1
    - Level 3 Item A.1.a
      - Level 4 Item A.1.a.i
  - Level 2 Item A.2
- Level 1 Item B`

	// Add document to server
	doc := &Document{
		URI:     "file:///test.md",
		Content: content,
	}
	server.documents[doc.URI] = doc

	// Test getting boundaries for "Level 1 Item A" (should include all children)
	position := map[string]any{
		"line":      float64(2), // "Level 1 Item A" is on line 2 (0-based)
		"character": float64(2), // Position in "Level 1 Item A"
	}

	arguments := []any{
		"file:///test.md",
		position,
	}

	result, err := server.handleGetListItemBoundaries(arguments)
	require.NoError(t, err)

	// Assert the result includes all children
	response, ok := result.(*BoundaryResponse)
	require.True(t, ok)
	assert.True(t, response.Found)
	assert.Equal(t, 2, response.Start.Line) // Line 2 (0-based)
	assert.Equal(t, 0, response.Start.Character)
	assert.Equal(t, 7, response.End.Line) // After "Level 2 Item A.2" (exclusive)
	assert.Equal(t, 0, response.End.Character)
}

func TestHandleGetListItemBoundaries_DeepNestedItem(t *testing.T) {
	server := createTestServerForBoundaries()

	// Test content - nested list
	content := `# Nested List

- Level 1 Item A
  - Level 2 Item A.1
    - Level 3 Item A.1.a
      - Level 4 Item A.1.a.i
  - Level 2 Item A.2
- Level 1 Item B`

	// Add document to server
	doc := &Document{
		URI:     "file:///test.md",
		Content: content,
	}
	server.documents[doc.URI] = doc

	// Test getting boundaries for "Level 4 Item A.1.a.i" (deepest item, no children)
	position := map[string]any{
		"line":      float64(5), // "Level 4 Item A.1.a.i" is on line 5 (0-based)
		"character": float64(8), // Position in the text
	}

	arguments := []any{
		"file:///test.md",
		position,
	}

	result, err := server.handleGetListItemBoundaries(arguments)
	require.NoError(t, err)

	// Assert the result is just the single item (no children)
	response, ok := result.(*BoundaryResponse)
	require.True(t, ok)
	assert.True(t, response.Found)
	assert.Equal(t, 5, response.Start.Line) // Line 5 (0-based)
	assert.Equal(t, 0, response.Start.Character)
	assert.Equal(t, 6, response.End.Line) // Next line (exclusive)
	assert.Equal(t, 0, response.End.Character)
}

func TestHandleGetListItemBoundaries_TaskList(t *testing.T) {
	server := createTestServerForBoundaries()

	// Test content - task list
	content := `# Task List

- [x] Completed task
  - [ ] Subtask A
    - [ ] Sub-subtask A.1.a
    - [x] Sub-subtask A.1.b
  - [x] Subtask B
- [ ] Incomplete task`

	// Add document to server
	doc := &Document{
		URI:     "file:///test.md",
		Content: content,
	}
	server.documents[doc.URI] = doc

	// Test getting boundaries for completed task (should include all subtasks)
	position := map[string]any{
		"line":      float64(2), // "Completed task" is on line 2 (0-based)
		"character": float64(6), // Position after "[x] "
	}

	arguments := []any{
		"file:///test.md",
		position,
	}

	result, err := server.handleGetListItemBoundaries(arguments)
	require.NoError(t, err)

	// Assert the result includes all subtasks
	response, ok := result.(*BoundaryResponse)
	require.True(t, ok)
	assert.True(t, response.Found)
	assert.Equal(t, 2, response.Start.Line) // Line 2 (0-based)
	assert.Equal(t, 0, response.Start.Character)
	assert.Equal(t, 7, response.End.Line) // After "Subtask B" (exclusive)
	assert.Equal(t, 0, response.End.Character)
}

func TestHandleGetListItemBoundaries_NotFound(t *testing.T) {
	server := createTestServerForBoundaries()

	// Test content with non-list text
	content := `# Test Document

This is some regular text.
Not a list item.

- First item
- Second item`

	// Add document to server
	doc := &Document{
		URI:     "file:///test.md",
		Content: content,
	}
	server.documents[doc.URI] = doc

	// Test getting boundaries for non-list text
	position := map[string]any{
		"line":      float64(3), // "Not a list item." is on line 3 (0-based)
		"character": float64(5), // Position in text
	}

	arguments := []any{
		"file:///test.md",
		position,
	}

	result, err := server.handleGetListItemBoundaries(arguments)
	require.NoError(t, err)

	// Assert no list item found
	response, ok := result.(*BoundaryResponse)
	require.True(t, ok)
	assert.False(t, response.Found)
}

func TestHandleGetListItemBoundaries_DocumentNotFound(t *testing.T) {
	server := createTestServerForBoundaries()

	// Test with non-existent document
	position := map[string]any{
		"line":      float64(0),
		"character": float64(0),
	}

	arguments := []any{
		"file:///nonexistent.md",
		position,
	}

	result, err := server.handleGetListItemBoundaries(arguments)
	require.NoError(t, err)

	// Assert document not found returns Found: false
	response, ok := result.(*BoundaryResponse)
	require.True(t, ok)
	assert.False(t, response.Found)
}

func TestHandleGetListItemBoundaries_InvalidArguments(t *testing.T) {
	server := createTestServerForBoundaries()

	tests := []struct {
		name      string
		arguments []any
		expectErr bool
	}{
		{
			name:      "too few arguments",
			arguments: []any{"file:///test.md"},
			expectErr: true,
		},
		{
			name:      "invalid URI type",
			arguments: []any{123, map[string]any{"line": float64(0), "character": float64(0)}},
			expectErr: true,
		},
		{
			name:      "invalid position type",
			arguments: []any{"file:///test.md", "not a position"},
			expectErr: true,
		},
		{
			name: "missing line in position",
			arguments: []any{
				"file:///test.md",
				map[string]any{"character": float64(0)},
			},
			expectErr: true,
		},
		{
			name: "missing character in position",
			arguments: []any{
				"file:///test.md",
				map[string]any{"line": float64(0)},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := server.handleGetListItemBoundaries(tt.arguments)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestHandleGetListItemBoundaries_NumberedList(t *testing.T) {
	server := createTestServerForBoundaries()

	// Test content - numbered list
	content := `# Numbered List

1. First numbered item
2. Second numbered item
   a. Sub item a
   b. Sub item b
3. Third numbered item`

	// Add document to server
	doc := &Document{
		URI:     "file:///test.md",
		Content: content,
	}
	server.documents[doc.URI] = doc

	// Test getting boundaries for second numbered item (with sub items)
	position := map[string]any{
		"line":      float64(3), // "Second numbered item" is on line 3 (0-based)
		"character": float64(3), // Position in text
	}

	arguments := []any{
		"file:///test.md",
		position,
	}

	result, err := server.handleGetListItemBoundaries(arguments)
	require.NoError(t, err)

	// Assert the result includes sub items
	response, ok := result.(*BoundaryResponse)
	require.True(t, ok)
	assert.True(t, response.Found)
	assert.Equal(t, 3, response.Start.Line) // Line 3 (0-based)
	assert.Equal(t, 0, response.Start.Character)
	assert.Equal(t, 6, response.End.Line) // After "Sub item b" (exclusive)
	assert.Equal(t, 0, response.End.Character)
}
