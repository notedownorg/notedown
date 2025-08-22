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
	"strings"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/config"
	"github.com/notedownorg/notedown/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTaskDiagnostics(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	tests := []struct {
		name           string
		content        string
		expectedCount  int
		expectedCodes  []string
		expectedRanges []lsp.Range
	}{
		{
			name: "valid task states",
			content: `# Todo List

- [ ] Valid todo task
- [x] Valid done task
- [X] Valid done task (alias)

This is not a task.`,
			expectedCount:  0,
			expectedCodes:  []string{},
			expectedRanges: []lsp.Range{},
		},
		{
			name: "invalid task states",
			content: `# Todo List

- [invalid] Invalid task state
- [?] Another invalid state
- [wip] Work in progress (invalid in default config)

This is not a task.`,
			expectedCount: 3,
			expectedCodes: []string{"invalid-task-state", "invalid-task-state", "invalid-task-state"},
			expectedRanges: []lsp.Range{
				{
					Start: lsp.Position{Line: 2, Character: 2},
					End:   lsp.Position{Line: 2, Character: 11},
				},
				{
					Start: lsp.Position{Line: 3, Character: 2},
					End:   lsp.Position{Line: 3, Character: 5},
				},
				{
					Start: lsp.Position{Line: 4, Character: 2},
					End:   lsp.Position{Line: 4, Character: 7},
				},
			},
		},
		{
			name: "mixed task states",
			content: `# Mixed Tasks

- [ ] Valid todo
- [x] Valid done
- [invalid] Invalid state
- [X] Valid done alias

1. [x] Numbered list with valid task
2. [bad] Numbered list with invalid task`,
			expectedCount: 2,
			expectedCodes: []string{"invalid-task-state", "invalid-task-state"},
			expectedRanges: []lsp.Range{
				{
					Start: lsp.Position{Line: 4, Character: 2},
					End:   lsp.Position{Line: 4, Character: 11},
				},
				{
					Start: lsp.Position{Line: 8, Character: 3},
					End:   lsp.Position{Line: 8, Character: 8},
				},
			},
		},
		{
			name: "tasks with various list formats",
			content: `# Various Lists

* [x] Asterisk list valid
* [invalid] Asterisk list invalid
+ [ ] Plus list valid
+ [bad] Plus list invalid
  - [x] Indented dash valid
  - [wrong] Indented dash invalid`,
			expectedCount: 3,
			expectedCodes: []string{"invalid-task-state", "invalid-task-state", "invalid-task-state"},
			expectedRanges: []lsp.Range{
				{
					Start: lsp.Position{Line: 3, Character: 2},
					End:   lsp.Position{Line: 3, Character: 11},
				},
				{
					Start: lsp.Position{Line: 5, Character: 2},
					End:   lsp.Position{Line: 5, Character: 7},
				},
				{
					Start: lsp.Position{Line: 7, Character: 4},
					End:   lsp.Position{Line: 7, Character: 11},
				},
			},
		},
		{
			name: "non-task checkboxes ignored",
			content: `# Not Tasks

This is [invalid] but not a task.
[invalid] at line start but not a list item.

Regular list:
- This is a regular list item [invalid] in middle
- [invalid] this should be detected as task

Paragraph with [x] checkbox not in list.`,
			expectedCount: 1,
			expectedCodes: []string{"invalid-task-state"},
			expectedRanges: []lsp.Range{
				{
					Start: lsp.Position{Line: 7, Character: 2},
					End:   lsp.Position{Line: 7, Character: 11},
				},
			},
		},
		{
			name: "wikilinks and markdown links should be ignored",
			content: `# Mixed Content

- [[wikilink]] Regular list item with wikilink
- [markdown-link](http://example.com) Regular list item with markdown link
- [[docs/architecture]] Full path wikilink
- [[target|display]] Wikilink with display text
- [x] Valid task
- [invalid] Invalid task that should be flagged
- [ ] Another valid task

Regular paragraph with [brackets] that should be ignored.
`,
			expectedCount: 1,
			expectedCodes: []string{"invalid-task-state"},
			expectedRanges: []lsp.Range{
				{
					Start: lsp.Position{Line: 7, Character: 2},
					End:   lsp.Position{Line: 7, Character: 11},
				},
			},
		},
		{
			name: "edge cases with complex brackets and spacing",
			content: `# Edge Cases

- [x]Valid task without space after bracket should not be detected
- [ ] Valid task with space
- [invalid]No space after bracket should not be detected  
- [wrong] Space after bracket should be detected as invalid
- [[wikilink]] Should be ignored
- [[complex/path/file]] Should be ignored
- [[target|display text]] Should be ignored
- [link text](https://example.com) Should be ignored
- [complex link text](path/to/file.md) Should be ignored
- Text with [brackets] in middle should be ignored

1. [x] Valid numbered task
2. [bad] Invalid numbered task
3. [[numbered/wikilink]] Should be ignored
`,
			expectedCount: 2,
			expectedCodes: []string{"invalid-task-state", "invalid-task-state"},
			expectedRanges: []lsp.Range{
				{
					Start: lsp.Position{Line: 5, Character: 2},
					End:   lsp.Position{Line: 5, Character: 9},
				},
				{
					Start: lsp.Position{Line: 14, Character: 3},
					End:   lsp.Position{Line: 14, Character: 8},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagnostics := server.generateTaskDiagnostics("file:///test.md", tt.content)

			assert.Equal(t, tt.expectedCount, len(diagnostics), "Unexpected number of diagnostics")

			for i, diagnostic := range diagnostics {
				if i < len(tt.expectedCodes) {
					assert.Equal(t, tt.expectedCodes[i], diagnostic.Code, "Unexpected diagnostic code")
				}
				if i < len(tt.expectedRanges) {
					assert.Equal(t, tt.expectedRanges[i], diagnostic.Range, "Unexpected diagnostic range")
				}

				// Verify common properties
				assert.Equal(t, lsp.DiagnosticSeverityWarning, *diagnostic.Severity, "Should be warning severity")
				assert.Equal(t, "notedown-task", *diagnostic.Source, "Should have notedown-task source")
				assert.Contains(t, diagnostic.Message, "Invalid task state", "Should contain error message")
				assert.Contains(t, diagnostic.Message, "Valid states:", "Should list valid states")
			}
		})
	}
}

func TestGenerateWikilinkDiagnostics(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Test that the function exists and doesn't crash
	diagnostics := server.generateWikilinkDiagnostics("file:///test.md", "# Test\n\n[[example]] wikilink")
	// Wikilink diagnostics require workspace setup, so we just verify it doesn't crash
	_ = diagnostics
}

func TestIsValidTaskState(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Test with default config
	defaultConfig := config.GetDefaultConfig()

	tests := []struct {
		name     string
		state    string
		config   *config.Config
		expected bool
	}{
		{"empty state valid", " ", defaultConfig, true},
		{"x state valid", "x", defaultConfig, true},
		{"X alias valid", "X", defaultConfig, true},
		{"completed alias valid", "completed", defaultConfig, true},
		{"invalid state", "invalid", defaultConfig, false},
		{"question mark invalid", "?", defaultConfig, false},
		{"wip invalid in default", "wip", defaultConfig, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.isValidTaskState(tt.state, tt.config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetValidTaskStates(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Test with default config
	defaultConfig := config.GetDefaultConfig()
	validStates := server.getValidTaskStates(defaultConfig)

	// Should contain the main values and all aliases
	expectedStates := []string{" ", "x", "X", "completed"}

	assert.Equal(t, len(expectedStates), len(validStates), "Should have correct number of valid states")

	for _, expected := range expectedStates {
		assert.Contains(t, validStates, expected, "Should contain expected state: %s", expected)
	}
}

func TestTaskDiagnosticsWithCustomConfig(t *testing.T) {
	// Create a custom config with additional task states
	customConfig := &config.Config{
		Tasks: config.TasksConfig{
			States: []config.TaskState{
				{
					Value:   " ",
					Name:    "todo",
					Aliases: []string{},
				},
				{
					Value:   "x",
					Name:    "done",
					Aliases: []string{"X"},
				},
				{
					Value:   "wip",
					Name:    "in-progress",
					Aliases: []string{"~"},
				},
			},
		},
	}

	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Since generateTaskDiagnostics loads config internally, we need to test isValidTaskState directly
	tests := []struct {
		state    string
		expected bool
	}{
		{" ", true},
		{"x", true},
		{"X", true},
		{"wip", true},
		{"~", true},
		{"invalid", false},
		{"?", false},
	}

	for _, tt := range tests {
		t.Run("state_"+tt.state, func(t *testing.T) {
			result := server.isValidTaskState(tt.state, customConfig)
			assert.Equal(t, tt.expected, result, "State %q should be %v", tt.state, tt.expected)
		})
	}
}

func TestTaskDiagnosticsEmptyContent(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	diagnostics := server.generateTaskDiagnostics("file:///empty.md", "")
	assert.Equal(t, 0, len(diagnostics), "Empty content should produce no diagnostics")
}

func TestTaskDiagnosticsNoTasks(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	content := `# Regular Markdown

This is a regular paragraph with no tasks.

## Another heading

- Regular list item
- Another regular item

1. Numbered list
2. Another numbered item`

	diagnostics := server.generateTaskDiagnostics("file:///notasks.md", content)
	assert.Equal(t, 0, len(diagnostics), "Content with no tasks should produce no diagnostics")
}

func TestPositionFromOffset(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	content := `line 0
line 1
line 2`

	tests := []struct {
		offset       int
		expectedLine int
		expectedChar int
	}{
		{0, 0, 0},  // Start of first line
		{3, 0, 3},  // Middle of first line
		{7, 1, 0},  // Start of second line (after \n)
		{10, 1, 3}, // Middle of second line
		{14, 2, 0}, // Start of third line
		{20, 2, 6}, // End of content
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			line, char := server.positionFromOffset(content, tt.offset)
			assert.Equal(t, tt.expectedLine, line, "Unexpected line for offset %d", tt.offset)
			assert.Equal(t, tt.expectedChar, char, "Unexpected character for offset %d", tt.offset)
		})
	}
}

func TestParserBasedTaskDiagnostics(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Test that parser correctly identifies invalid vs valid tasks
	content := `# Task List

- [ ] Valid todo task
- [x] Valid done task  
- [invalid] This should be invalid
- [wip] This should also be invalid
- [X] Valid done task with alias

Some text in between.

1. [ ] Valid numbered task
2. [bad] Invalid numbered task`

	diagnostics := server.generateTaskDiagnostics("file:///test.md", content)

	// Should find exactly 3 invalid task states
	assert.Equal(t, 3, len(diagnostics), "Should find 3 invalid task states")

	// Verify the specific invalid states were detected
	invalidStates := make([]string, len(diagnostics))
	for i, diag := range diagnostics {
		// Extract state from message: "Invalid task state 'STATE'..."
		start := strings.Index(diag.Message, "'") + 1
		end := strings.Index(diag.Message[start:], "'") + start
		invalidStates[i] = diag.Message[start:end]
	}

	expectedInvalid := []string{"invalid", "wip", "bad"}
	for _, expected := range expectedInvalid {
		assert.Contains(t, invalidStates, expected, "Should detect invalid state: %s", expected)
	}

	// Verify all diagnostics have correct properties
	for _, diag := range diagnostics {
		assert.Equal(t, lsp.DiagnosticSeverityWarning, *diag.Severity, "Should be warning severity")
		assert.Equal(t, "notedown-task", *diag.Source, "Should have notedown-task source")
		assert.Equal(t, "invalid-task-state", diag.Code, "Should have invalid-task-state code")
		assert.Contains(t, diag.Message, "Invalid task state", "Should contain error message")
		assert.Contains(t, diag.Message, "Valid states:", "Should list valid states")
	}
}
