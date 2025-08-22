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
	"encoding/json"
	"testing"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleFoldingRange(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	// Add a test document
	testURI := "file:///test.md"
	doc, err := server.AddDocument(testURI)
	require.NoError(t, err)

	// Set document content with various foldable sections
	doc.Content = `# Header 1

Some content under header 1.

## Header 2

More content under header 2.

- [ ] Task 1
- [x] Task 2
  - [ ] Subtask 1
  - [ ] Subtask 2
- [ ] Task 3

Regular list:
- Item 1
- Item 2
  - Nested item
- Item 3

` + "```" + `javascript
function test() {
  console.log("hello");
}
` + "```" + `

# Another Header

Final content.`

	// Prepare folding range request
	params := lsp.FoldingRangeParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: testURI},
	}
	jsonParams, err := json.Marshal(params)
	require.NoError(t, err)

	// Call the folding range handler
	result, err := server.handleFoldingRange(jsonParams)
	require.NoError(t, err)

	// Verify the result
	ranges, ok := result.([]lsp.FoldingRange)
	require.True(t, ok, "Expected []lsp.FoldingRange")
	assert.NotEmpty(t, ranges, "Should have generated folding ranges")

	// Debug: print all ranges
	t.Logf("Generated %d folding ranges:", len(ranges))
	for i, r := range ranges {
		kind := "none"
		if r.Kind != nil {
			kind = string(*r.Kind)
		}
		t.Logf("  Range %d: lines %d-%d, kind: %s", i, r.StartLine, r.EndLine, kind)
	}

	// Verify we have ranges for different types of content
	hasHeaderRange := false
	hasListRange := false
	hasCodeRange := false

	for _, r := range ranges {
		// Check based on line numbers what type of range this likely is
		if r.StartLine == 0 { // "# Header 1" starts at line 0
			hasHeaderRange = true
		}
		if r.StartLine >= 8 && r.StartLine <= 15 { // Task list area
			hasListRange = true
		}
		if r.StartLine >= 19 && r.StartLine <= 23 { // Code block area
			hasCodeRange = true
		}
	}

	assert.True(t, hasHeaderRange, "Should have header folding range")
	assert.True(t, hasListRange, "Should have list folding range")
	assert.True(t, hasCodeRange, "Should have code block folding range")
}

func TestGenerateHeaderFoldingRanges(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	_ = `# Header 1
Content under header 1.

## Header 2
Content under header 2.

### Header 3
Content under header 3.

## Another Header 2
More content.

# Header 1 Again
Final content.`

	lines := []string{
		"# Header 1",
		"Content under header 1.",
		"",
		"## Header 2",
		"Content under header 2.",
		"",
		"### Header 3",
		"Content under header 3.",
		"",
		"## Another Header 2",
		"More content.",
		"",
		"# Header 1 Again",
		"Final content.",
	}

	ranges := server.generateHeaderFoldingRanges(lines)

	// Should have ranges for each header section
	assert.NotEmpty(t, ranges, "Should generate header folding ranges")

	// Verify ranges are reasonable
	for _, r := range ranges {
		assert.Greater(t, r.EndLine, r.StartLine, "End line should be after start line")
		assert.Equal(t, lsp.FoldingRangeKindRegion, *r.Kind, "Header ranges should be region kind")
	}
}

func TestGenerateListFoldingRanges(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	lines := []string{
		"- [ ] Task 1",
		"- [x] Task 2",
		"  - [ ] Subtask 1",
		"  - [ ] Subtask 2",
		"- [ ] Task 3",
		"",
		"Regular list:",
		"- Item 1",
		"- Item 2",
		"  - Nested item",
		"- Item 3",
	}

	ranges := server.generateListFoldingRanges(lines)

	// Should generate some list folding ranges
	assert.NotEmpty(t, ranges, "Should generate list folding ranges")

	// Verify ranges are reasonable
	for _, r := range ranges {
		assert.GreaterOrEqual(t, r.EndLine, r.StartLine, "End line should be at or after start line")
	}
}

func TestGenerateCodeBlockFoldingRanges(t *testing.T) {
	logger := log.NewDefault()
	server := NewServer("test", logger)

	lines := []string{
		"Some text",
		"```javascript",
		"function test() {",
		"  console.log('hello');",
		"}",
		"```",
		"More text",
		"```",
		"plain code block",
		"```",
		"End",
	}

	ranges := server.generateCodeBlockFoldingRanges(lines)

	// Should generate folding ranges for both code blocks
	assert.Len(t, ranges, 2, "Should generate 2 code block folding ranges")

	// Verify ranges
	assert.Equal(t, 1, ranges[0].StartLine, "First code block should start at line 1")
	assert.Equal(t, 5, ranges[0].EndLine, "First code block should end at line 5")
	assert.Equal(t, lsp.FoldingRangeKindRegion, *ranges[0].Kind, "Code block should be region kind")

	assert.Equal(t, 7, ranges[1].StartLine, "Second code block should start at line 7")
	assert.Equal(t, 9, ranges[1].EndLine, "Second code block should end at line 9")
	assert.Equal(t, lsp.FoldingRangeKindRegion, *ranges[1].Kind, "Code block should be region kind")
}
