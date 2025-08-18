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

package parser

import (
	"testing"

	"github.com/notedownorg/notedown/pkg/config"
	"github.com/notedownorg/notedown/pkg/parser/extensions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
)

func TestConfigurableTaskStates(t *testing.T) {
	// Create a test configuration with custom task states
	cfg := &config.Config{
		Tasks: config.TasksConfig{
			States: []config.TaskState{
				{Value: " ", Name: "Unchecked"},
				{Value: "x", Name: "Checked"},
				{Value: "wip", Name: "Work in Progress"},
				{Value: "?", Name: "Question"},
				{Value: "!", Name: "Important"},
			},
		},
	}

	// Create parser with this configuration
	parser := &NotedownParser{
		goldmark: createTestGoldmarkWithConfig(cfg),
	}

	tests := []struct {
		name        string
		markdown    string
		expectedLen int
		checkStates func(t *testing.T, doc *Document)
	}{
		{
			name: "basic task states",
			markdown: `- [ ] Unchecked task
- [x] Checked task
- [wip] Work in progress task
- [?] Question task
- [!] Important task`,
			expectedLen: 1, // One list
			checkStates: func(t *testing.T, doc *Document) {
				list := findFirstList(t, doc)
				items := list.GetListItems()
				require.Len(t, items, 5)

				// Check each task state
				assert.True(t, items[0].TaskList)
				assert.Equal(t, " ", items[0].TaskState)

				assert.True(t, items[1].TaskList)
				assert.Equal(t, "x", items[1].TaskState)

				assert.True(t, items[2].TaskList)
				assert.Equal(t, "wip", items[2].TaskState)

				assert.True(t, items[3].TaskList)
				assert.Equal(t, "?", items[3].TaskState)

				assert.True(t, items[4].TaskList)
				assert.Equal(t, "!", items[4].TaskState)
			},
		},
		{
			name: "invalid task states should not be parsed as tasks",
			markdown: `- [invalid] Invalid state
- [xyz] Unknown state
- [123] Numeric state
- [ ] Valid unchecked
- [x] Valid checked`,
			expectedLen: 1,
			checkStates: func(t *testing.T, doc *Document) {
				list := findFirstList(t, doc)
				items := list.GetListItems()
				require.Len(t, items, 5)

				// First three should not be task list items
				assert.False(t, items[0].TaskList, "Invalid state should not create task")
				assert.False(t, items[1].TaskList, "Unknown state should not create task")
				assert.False(t, items[2].TaskList, "Numeric state should not create task")

				// Last two should be valid task list items
				assert.True(t, items[3].TaskList, "Valid unchecked should be task")
				assert.Equal(t, " ", items[3].TaskState)

				assert.True(t, items[4].TaskList, "Valid checked should be task")
				assert.Equal(t, "x", items[4].TaskState)
			},
		},
		{
			name: "mixed regular and task list items",
			markdown: `- Regular list item
- [x] Task item
- Another regular item
- [wip] Another task item
- Final regular item`,
			expectedLen: 1,
			checkStates: func(t *testing.T, doc *Document) {
				list := findFirstList(t, doc)
				items := list.GetListItems()
				require.Len(t, items, 5)

				// Check which items are tasks vs regular
				assert.False(t, items[0].TaskList, "Regular item should not be task")
				assert.True(t, items[1].TaskList, "Should be task item")
				assert.Equal(t, "x", items[1].TaskState)
				assert.False(t, items[2].TaskList, "Regular item should not be task")
				assert.True(t, items[3].TaskList, "Should be task item")
				assert.Equal(t, "wip", items[3].TaskState)
				assert.False(t, items[4].TaskList, "Regular item should not be task")
			},
		},
		{
			name: "nested task lists",
			markdown: `- [x] Parent task
  - [wip] Child task
  - [ ] Another child task
- [?] Another parent task`,
			expectedLen: 2, // Parent list and nested list
			checkStates: func(t *testing.T, doc *Document) {
				// Find all lists in the document
				var lists []*List
				walker := NewWalker(WalkFunc(func(node Node) error {
					if list, ok := node.(*List); ok {
						lists = append(lists, list)
					}
					return nil
				}))
				_ = walker.Walk(doc)

				require.Len(t, lists, 2, "Should have parent and nested list")

				// Check parent list
				parentItems := lists[0].GetListItems()
				require.Len(t, parentItems, 2)
				assert.Equal(t, "x", parentItems[0].TaskState)
				assert.Equal(t, "?", parentItems[1].TaskState)

				// Check nested list
				nestedItems := lists[1].GetListItems()
				require.Len(t, nestedItems, 2)
				assert.Equal(t, "wip", nestedItems[0].TaskState)
				assert.Equal(t, " ", nestedItems[1].TaskState)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.ParseString(tt.markdown)
			require.NoError(t, err)
			require.NotNil(t, doc)

			tt.checkStates(t, doc)
		})
	}
}

func TestTaskStateValues(t *testing.T) {
	// Test that task state values are correctly preserved
	cfg := &config.Config{
		Tasks: config.TasksConfig{
			States: []config.TaskState{
				{Value: " ", Name: "Unchecked"},
				{Value: "x", Name: "Checked"},
				{Value: "done", Name: "Done"},
				{Value: "cancelled", Name: "Cancelled"},
			},
		},
	}

	parser := &NotedownParser{
		goldmark: createTestGoldmarkWithConfig(cfg),
	}

	markdown := `- [ ] Unchecked task
- [x] Checked task
- [done] Done task
- [cancelled] Cancelled task`

	doc, err := parser.ParseString(markdown)
	require.NoError(t, err)

	list := findFirstList(t, doc)
	items := list.GetListItems()
	require.Len(t, items, 4)

	// Test that state values are preserved correctly
	assert.Equal(t, " ", items[0].TaskState)
	assert.Equal(t, "x", items[1].TaskState)
	assert.Equal(t, "done", items[2].TaskState)
	assert.Equal(t, "cancelled", items[3].TaskState)
}

func TestDefaultConfigurationTaskStates(t *testing.T) {
	// Test with default configuration (should only support " " and "x")
	parser := NewParser()

	markdown := `- [ ] Unchecked task
- [x] Checked task
- [wip] Should not be parsed as task
- [done] Should not be parsed as task`

	doc, err := parser.ParseString(markdown)
	require.NoError(t, err)

	list := findFirstList(t, doc)
	items := list.GetListItems()
	require.Len(t, items, 4)

	// Only first two should be task items
	assert.True(t, items[0].TaskList)
	assert.Equal(t, " ", items[0].TaskState)

	assert.True(t, items[1].TaskList)
	assert.Equal(t, "x", items[1].TaskState)

	// Last two should not be task items (invalid states)
	assert.False(t, items[2].TaskList)
	assert.False(t, items[3].TaskList)
}

// Helper function to find the first list in a document
func findFirstList(t *testing.T, doc *Document) *List {
	var list *List
	walker := NewWalker(WalkFunc(func(node Node) error {
		if l, ok := node.(*List); ok && list == nil {
			list = l
		}
		return nil
	}))
	_ = walker.Walk(doc)
	require.NotNil(t, list, "Should find at least one list")
	return list
}

// Helper function to create goldmark instance with test configuration
func createTestGoldmarkWithConfig(cfg *config.Config) goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extensions.NewTaskListExtension(cfg),
		),
	)
}
