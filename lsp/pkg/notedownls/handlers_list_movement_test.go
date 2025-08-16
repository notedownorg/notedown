package notedownls

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

func TestParseListItem(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	tests := []struct {
		name        string
		line        string
		lineNum     int
		expected    *ListItem
		shouldExist bool
	}{
		{
			name:    "task list item",
			line:    "- [ ] This is a task",
			lineNum: 0,
			expected: &ListItem{
				StartLine:   0,
				EndLine:     0,
				IndentLevel: 0,
				Marker:      "- [ ]",
				Content:     "This is a task",
			},
			shouldExist: true,
		},
		{
			name:    "completed task",
			line:    "- [x] Completed task",
			lineNum: 1,
			expected: &ListItem{
				StartLine:   1,
				EndLine:     1,
				IndentLevel: 0,
				Marker:      "- [x]",
				Content:     "Completed task",
			},
			shouldExist: true,
		},
		{
			name:    "bullet list item",
			line:    "- Simple bullet",
			lineNum: 2,
			expected: &ListItem{
				StartLine:   2,
				EndLine:     2,
				IndentLevel: 0,
				Marker:      "-",
				Content:     "Simple bullet",
			},
			shouldExist: true,
		},
		{
			name:    "indented bullet",
			line:    "  - Indented bullet",
			lineNum: 3,
			expected: &ListItem{
				StartLine:   3,
				EndLine:     3,
				IndentLevel: 2,
				Marker:      "-",
				Content:     "Indented bullet",
			},
			shouldExist: true,
		},
		{
			name:    "numbered list",
			line:    "1. First item",
			lineNum: 4,
			expected: &ListItem{
				StartLine:   4,
				EndLine:     4,
				IndentLevel: 0,
				Marker:      "1.",
				Content:     "First item",
			},
			shouldExist: true,
		},
		{
			name:    "numbered list with multiple digits",
			line:    "123. Item 123",
			lineNum: 5,
			expected: &ListItem{
				StartLine:   5,
				EndLine:     5,
				IndentLevel: 0,
				Marker:      "123.",
				Content:     "Item 123",
			},
			shouldExist: true,
		},
		{
			name:        "not a list item",
			line:        "Just regular text",
			lineNum:     6,
			expected:    nil,
			shouldExist: false,
		},
		{
			name:        "heading",
			line:        "# This is a heading",
			lineNum:     7,
			expected:    nil,
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.parseListItem(tt.line, tt.lineNum)

			if tt.shouldExist {
				require.NotNil(t, result, "Expected to parse list item")
				assert.Equal(t, tt.expected.StartLine, result.StartLine)
				assert.Equal(t, tt.expected.EndLine, result.EndLine)
				assert.Equal(t, tt.expected.IndentLevel, result.IndentLevel)
				assert.Equal(t, tt.expected.Marker, result.Marker)
				assert.Equal(t, tt.expected.Content, result.Content)
			} else {
				assert.Nil(t, result, "Expected no list item to be parsed")
			}
		})
	}
}

func TestParseListHierarchy(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	content := `# Test Document

- Top level 1
  - Nested 1.1
  - Nested 1.2
- Top level 2
- [ ] Task item
  - [ ] Subtask 1
  - [x] Subtask 2

1. Numbered item 1
2. Numbered item 2
   - Bullet under numbered
3. Numbered item 3`

	hierarchy, err := server.parseListHierarchy(content)
	require.NoError(t, err, "Failed to parse list hierarchy")

	// Check top-level items count
	expectedTopLevel := 6 // "Top level 1", "Top level 2", "Task item", "Numbered item 1", "Numbered item 2", "Numbered item 3"
	assert.Len(t, hierarchy.Items, expectedTopLevel, "Unexpected number of top-level items")

	// Check that "Top level 1" has 2 children
	if len(hierarchy.Items) > 0 {
		topLevel1 := hierarchy.Items[0]
		assert.Len(t, topLevel1.Children, 2, "Expected 'Top level 1' to have 2 children")
		assert.Equal(t, "Top level 1", topLevel1.Content)
	}

	// Check that task item has 2 children
	taskItemIndex := -1
	for i, item := range hierarchy.Items {
		if item.Marker == "- [ ]" {
			taskItemIndex = i
			break
		}
	}
	require.GreaterOrEqual(t, taskItemIndex, 0, "Expected to find task item")

	taskItem := hierarchy.Items[taskItemIndex]
	assert.Len(t, taskItem.Children, 2, "Expected task item to have 2 children")
}

func TestFindItemAtPosition(t *testing.T) {
	hierarchy := &ListHierarchy{
		Items: []*ListItem{
			{
				StartLine:   2,
				EndLine:     2,
				IndentLevel: 0,
				Marker:      "-",
				Content:     "First item",
				Children: []*ListItem{
					{
						StartLine:   3,
						EndLine:     3,
						IndentLevel: 2,
						Marker:      "-",
						Content:     "Nested item",
					},
				},
			},
			{
				StartLine:   4,
				EndLine:     4,
				IndentLevel: 0,
				Marker:      "-",
				Content:     "Second item",
			},
		},
		LineToItem: map[int]*ListItem{},
	}

	// Populate line to item mapping
	hierarchy.LineToItem[2] = hierarchy.Items[0]
	hierarchy.LineToItem[3] = hierarchy.Items[0].Children[0]
	hierarchy.LineToItem[4] = hierarchy.Items[1]

	tests := []struct {
		name     string
		position lsp.Position
		expected *ListItem
	}{
		{
			name:     "find first item",
			position: lsp.Position{Line: 2, Character: 5},
			expected: hierarchy.Items[0],
		},
		{
			name:     "find nested item",
			position: lsp.Position{Line: 3, Character: 5},
			expected: hierarchy.Items[0].Children[0],
		},
		{
			name:     "find second item",
			position: lsp.Position{Line: 4, Character: 5},
			expected: hierarchy.Items[1],
		},
		{
			name:     "position not on list item",
			position: lsp.Position{Line: 1, Character: 5},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hierarchy.findItemAtPosition(tt.position)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateListItemMove(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	// Create a simple hierarchy for testing
	hierarchy := &ListHierarchy{
		Items: []*ListItem{
			{
				StartLine:     2,
				EndLine:       2,
				IndentLevel:   0,
				Marker:        "-",
				Content:       "First item",
				OriginalLines: []string{"- First item"},
			},
			{
				StartLine:     3,
				EndLine:       3,
				IndentLevel:   0,
				Marker:        "-",
				Content:       "Second item",
				OriginalLines: []string{"- Second item"},
			},
			{
				StartLine:     4,
				EndLine:       4,
				IndentLevel:   0,
				Marker:        "-",
				Content:       "Third item",
				OriginalLines: []string{"- Third item"},
			},
		},
	}

	tests := []struct {
		name          string
		itemIndex     int
		moveUp        bool
		expectSuccess bool
	}{
		{
			name:          "move second item up",
			itemIndex:     1,
			moveUp:        true,
			expectSuccess: true,
		},
		{
			name:          "move first item up (boundary)",
			itemIndex:     0,
			moveUp:        true,
			expectSuccess: false,
		},
		{
			name:          "move second item down",
			itemIndex:     1,
			moveUp:        false,
			expectSuccess: true,
		},
		{
			name:          "move last item down (boundary)",
			itemIndex:     2,
			moveUp:        false,
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := hierarchy.Items[tt.itemIndex]
			documentURI := "file:///test.md"

			workspaceEdit, err := server.calculateListItemMove(hierarchy, item, tt.moveUp, documentURI)

			if tt.expectSuccess {
				assert.NoError(t, err, "Expected successful move calculation")
				assert.NotNil(t, workspaceEdit, "Expected workspace edit")
				assert.NotEmpty(t, workspaceEdit.Changes, "Expected workspace edit to have changes")

				edits, exists := workspaceEdit.Changes[documentURI]
				assert.True(t, exists, "Expected edits for document URI")
				assert.NotEmpty(t, edits, "Expected text edits")
			} else {
				// Boundary case should return nil workspace edit (no move possible)
				assert.NoError(t, err, "Should not error on boundary")
				assert.Nil(t, workspaceEdit, "Expected nil workspace edit for boundary case")
			}
		})
	}
}

func TestCreateSwapTextEdits(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	item1 := &ListItem{
		StartLine:     2,
		EndLine:       2,
		OriginalLines: []string{"- First item"},
	}

	item2 := &ListItem{
		StartLine:     3,
		EndLine:       3,
		OriginalLines: []string{"- Second item"},
	}

	edits := server.createSwapTextEdits(item1, item2)

	// Should create 2 edits (one for each item)
	assert.Len(t, edits, 2, "Expected 2 text edits for swap")

	// Verify edit ranges
	assert.Equal(t, 2, edits[0].Range.Start.Line, "First edit should start at line 2")
	assert.Equal(t, 3, edits[1].Range.Start.Line, "Second edit should start at line 3")

	// Verify content swap
	assert.Equal(t, "- Second item", edits[0].NewText, "First edit should contain second item content")
	assert.Equal(t, "- First item", edits[1].NewText, "Second edit should contain first item content")
}

func TestOrderedListRenumbering(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	parentItems := []*ListItem{
		{
			StartLine:     2,
			EndLine:       2,
			IndentLevel:   0,
			Marker:        "1.",
			OriginalLines: []string{"1. First item"},
		},
		{
			StartLine:     3,
			EndLine:       3,
			IndentLevel:   0,
			Marker:        "2.",
			OriginalLines: []string{"2. Second item"},
		},
		{
			StartLine:     4,
			EndLine:       4,
			IndentLevel:   0,
			Marker:        "3.",
			OriginalLines: []string{"3. Third item"},
		},
	}

	// Initial edits (swap edits)
	initialEdits := []lsp.TextEdit{
		{
			Range: lsp.Range{
				Start: lsp.Position{Line: 2, Character: 0},
				End:   lsp.Position{Line: 3, Character: 0},
			},
			NewText: "2. Second item",
		},
		{
			Range: lsp.Range{
				Start: lsp.Position{Line: 3, Character: 0},
				End:   lsp.Position{Line: 4, Character: 0},
			},
			NewText: "1. First item",
		},
	}

	// Test renumbering after swapping items 0 and 1
	edits := server.handleOrderedListRenumbering(initialEdits, parentItems, 0, 1)

	// Should have original swap edits plus renumbering edits
	assert.GreaterOrEqual(t, len(edits), 2, "Expected at least the original swap edits")

	// Look for renumbering edits
	hasRenumbering := false
	for _, edit := range edits {
		if edit.NewText == "2. First item" || edit.NewText == "1. Second item" {
			hasRenumbering = true
			break
		}
	}
	assert.True(t, hasRenumbering, "Expected to find renumbering edits")
}

// Integration tests

func TestListMovementIntegration(t *testing.T) {
	// Create a test server
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	// Test document content
	testContent := `# Test Document

- First item
- Second item  
- Third item

1. Numbered first
2. Numbered second
3. Numbered third

- [ ] Task one
- [x] Task two
- [ ] Task three
`

	testURI := "file:///test.md"

	// Add document to server and set content
	doc, err := server.AddDocument(testURI)
	require.NoError(t, err, "Failed to add document")
	doc.Content = testContent

	tests := []struct {
		name          string
		line          int
		character     int
		moveUp        bool
		expectSuccess bool
	}{
		{
			name:          "move second bullet item up",
			line:          3, // "Second item" (0-based)
			character:     0,
			moveUp:        true,
			expectSuccess: true,
		},
		{
			name:          "move first bullet item up (boundary)",
			line:          2, // "First item"
			character:     0,
			moveUp:        true,
			expectSuccess: false, // Should fail - already at top
		},
		{
			name:          "move task item down",
			line:          10, // "Task one" (corrected line number)
			character:     0,
			moveUp:        false,
			expectSuccess: true,
		},
		{
			name:          "move numbered item up",
			line:          6, // "Numbered second"
			character:     0,
			moveUp:        true,
			expectSuccess: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare command arguments
			args := []any{
				testURI,
				map[string]any{
					"line":      float64(tt.line),
					"character": float64(tt.character),
				},
			}

			var result any
			var err error

			if tt.moveUp {
				result, err = server.handleMoveListItemUp(args)
			} else {
				result, err = server.handleMoveListItemDown(args)
			}

			if tt.expectSuccess {
				assert.NoError(t, err, "Expected successful list movement")

				// Check that result is a workspace edit
				workspaceEdit, ok := result.(*lsp.WorkspaceEdit)
				require.True(t, ok, "Expected WorkspaceEdit result")

				// Verify workspace edit has changes
				assert.NotEmpty(t, workspaceEdit.Changes, "Expected workspace edit to have changes")

				edits, exists := workspaceEdit.Changes[testURI]
				assert.True(t, exists, "Expected workspace edit to have changes for test URI")
				assert.NotEmpty(t, edits, "Expected workspace edit to have text edits")

				t.Logf("Successfully generated %d text edits", len(edits))

			} else {
				assert.Error(t, err, "Expected error for boundary case")
			}
		})
	}
}

func TestListMovementWithComplexHierarchy(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	complexContent := `# Complex List Test

- Main item 1
  - Sub item 1.1
  - Sub item 1.2
    - Sub-sub item 1.2.1
    - Sub-sub item 1.2.2
  - Sub item 1.3
- Main item 2
  - Sub item 2.1
- Main item 3

1. Ordered main 1
   - Mixed bullet under ordered
   - Another mixed bullet
2. Ordered main 2
3. Ordered main 3
`

	testURI := "file:///complex.md"
	doc, err := server.AddDocument(testURI)
	require.NoError(t, err, "Failed to add document")
	doc.Content = complexContent

	// Test moving a main item with nested children
	args := []any{
		testURI,
		map[string]any{
			"line":      float64(2), // "Main item 1"
			"character": float64(0),
		},
	}

	result, err := server.handleMoveListItemDown(args)
	require.NoError(t, err, "Failed to move complex list item")

	workspaceEdit, ok := result.(*lsp.WorkspaceEdit)
	require.True(t, ok, "Expected WorkspaceEdit result")

	edits := workspaceEdit.Changes[testURI]
	assert.NotEmpty(t, edits, "Expected text edits for complex hierarchy move")

	t.Logf("Complex hierarchy move generated %d edits", len(edits))
}

func TestListMovementValidation(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	testContent := `# Test Document

Regular paragraph

- List item 1
- List item 2

More regular text
`

	testURI := "file:///test.md"
	doc, err := server.AddDocument(testURI)
	require.NoError(t, err)
	doc.Content = testContent

	t.Run("position not on list item", func(t *testing.T) {
		args := []any{
			testURI,
			map[string]any{
				"line":      float64(1), // Regular paragraph
				"character": float64(0),
			},
		}

		result, err := server.handleMoveListItemDown(args)
		assert.Error(t, err, "Expected error when position is not on list item")
		assert.Nil(t, result)
	})

	t.Run("invalid document URI", func(t *testing.T) {
		args := []any{
			"file:///nonexistent.md",
			map[string]any{
				"line":      float64(0),
				"character": float64(0),
			},
		}

		result, err := server.handleMoveListItemDown(args)
		assert.Error(t, err, "Expected error for nonexistent document")
		assert.Nil(t, result)
	})
}

func TestListMovementDebugRealWorld(t *testing.T) {
	logger := log.New(os.Stderr, log.Error)
	server := NewServer("test", logger)

	// Use the exact content from the user's lists.md
	testContent := `# Task Lists Test

## Simple Tasks

- [ ] Unchecked task
- [x] Checked task  
- [ ] Another unchecked task

## Nested Tasks

- [ ] Main task
  - [ ] Subtask 1
  - [x] Completed subtask
  - [ ] Subtask 3
    - [ ] Sub-subtask
- [x] Another main task

## Bullet Lists

- First item
- Second item
  - Nested item
  - Another nested item
- Third item

## Code Block Test

` + "```" + `javascript
function test() {
  console.log("Testing folding");
}
` + "```" + `

End of file.`

	testURI := "file:///test.md"

	// Add document to server and set content
	doc, err := server.AddDocument(testURI)
	require.NoError(t, err, "Failed to add document")
	doc.Content = testContent

	// Test moving the "Main task" down (line 10, 0-based)
	args := []any{
		testURI,
		map[string]any{
			"line":      float64(10), // "Main task" line (0-based)
			"character": float64(0),
		},
	}

	t.Logf("=== BEFORE ===")
	lines := strings.Split(testContent, "\n")
	for i, line := range lines {
		t.Logf("%2d: %s", i, line)
	}

	result, err := server.handleMoveListItemDown(args)
	require.NoError(t, err, "Failed to move list item")

	workspaceEdit, ok := result.(*lsp.WorkspaceEdit)
	require.True(t, ok, "Expected WorkspaceEdit result")

	t.Logf("=== TEXT EDITS ===")
	edits := workspaceEdit.Changes[testURI]
	for i, edit := range edits {
		t.Logf("Edit %d:", i)
		t.Logf("  Range: %d:%d to %d:%d",
			edit.Range.Start.Line, edit.Range.Start.Character,
			edit.Range.End.Line, edit.Range.End.Character)
		t.Logf("  NewText: %q", edit.NewText)
	}

	// Apply the edits manually to see the result
	finalContent := applyTextEditsManually(testContent, edits)

	t.Logf("=== AFTER ===")
	finalLines := strings.Split(finalContent, "\n")
	for i, line := range finalLines {
		t.Logf("%2d: %s", i, line)
	}

	// Verify the expected outcome
	expectedAfter := `# Task Lists Test

## Simple Tasks

- [ ] Unchecked task
- [x] Checked task  
- [ ] Another unchecked task

## Nested Tasks

- [x] Another main task
- [ ] Main task
  - [ ] Subtask 1
  - [x] Completed subtask
  - [ ] Subtask 3
    - [ ] Sub-subtask

## Bullet Lists

- First item
- Second item
  - Nested item
  - Another nested item
- Third item

## Code Block Test

` + "```" + `javascript
function test() {
  console.log("Testing folding");
}
` + "```" + `

End of file.`

	assert.Equal(t, expectedAfter, finalContent, "Content after move should match expected")
}

// applyTextEditsManually manually applies text edits to see the result
func applyTextEditsManually(content string, edits []lsp.TextEdit) string {
	lines := strings.Split(content, "\n")

	// Apply edits in reverse order to avoid index shifting issues
	for i := len(edits) - 1; i >= 0; i-- {
		edit := edits[i]
		startLine := edit.Range.Start.Line
		endLine := edit.Range.End.Line

		// Split new text into lines
		newLines := strings.Split(edit.NewText, "\n")

		// Replace the range
		before := lines[:startLine]
		after := lines[endLine:]

		// Combine before + newLines + after
		lines = append(before, append(newLines, after...)...)
	}

	return strings.Join(lines, "\n")
}
