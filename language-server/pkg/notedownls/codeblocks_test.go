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
	"testing"

	"github.com/notedownorg/notedown/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCodeBlockCollector_NewCodeBlockCollector(t *testing.T) {
	collector := NewCodeBlockCollector("go")
	require.NotNil(t, collector)
	assert.Equal(t, "go", collector.language)
	assert.Empty(t, collector.codeBlocks)
}

func TestCodeBlockCollector_Visit_FencedCodeBlock(t *testing.T) {
	collector := NewCodeBlockCollector("go")

	// Create a fenced code block
	codeBlock := parser.NewCodeBlock("go", "fmt.Println(\"Hello\")", true, parser.Range{
		Start: parser.Position{Line: 1, Column: 1},
		End:   parser.Position{Line: 3, Column: 1},
	})

	err := collector.Visit(codeBlock)
	require.NoError(t, err)

	blocks := collector.GetCodeBlocks()
	require.Len(t, blocks, 1)
	assert.Equal(t, "go", blocks[0].Language)
	assert.Equal(t, "fmt.Println(\"Hello\")", blocks[0].Content)
}

func TestCodeBlockCollector_Visit_IndentedCodeBlock(t *testing.T) {
	collector := NewCodeBlockCollector("go")

	// Create an indented code block (should be ignored)
	codeBlock := parser.NewCodeBlock("", "    fmt.Println(\"Hello\")", false, parser.Range{
		Start: parser.Position{Line: 1, Column: 1},
		End:   parser.Position{Line: 1, Column: 25},
	})

	err := collector.Visit(codeBlock)
	require.NoError(t, err)

	blocks := collector.GetCodeBlocks()
	assert.Empty(t, blocks) // Indented code blocks should be ignored
}

func TestCodeBlockCollector_Visit_WrongLanguage(t *testing.T) {
	collector := NewCodeBlockCollector("go")

	// Create a Python code block
	codeBlock := parser.NewCodeBlock("python", "print(\"Hello\")", true, parser.Range{
		Start: parser.Position{Line: 1, Column: 1},
		End:   parser.Position{Line: 3, Column: 1},
	})

	err := collector.Visit(codeBlock)
	require.NoError(t, err)

	blocks := collector.GetCodeBlocks()
	assert.Empty(t, blocks) // Should be filtered out
}

func TestCodeBlockCollector_Visit_EmptyLanguageFilter(t *testing.T) {
	collector := NewCodeBlockCollector("") // Empty language = collect all

	// Create code blocks with different languages
	goBlock := parser.NewCodeBlock("go", "fmt.Println(\"Hello\")", true, parser.Range{})
	pythonBlock := parser.NewCodeBlock("python", "print(\"Hello\")", true, parser.Range{})

	err := collector.Visit(goBlock)
	require.NoError(t, err)

	err = collector.Visit(pythonBlock)
	require.NoError(t, err)

	blocks := collector.GetCodeBlocks()
	require.Len(t, blocks, 2)

	// Check that both blocks are collected
	languages := []string{blocks[0].Language, blocks[1].Language}
	assert.Contains(t, languages, "go")
	assert.Contains(t, languages, "python")
}

func TestCodeBlockCollector_Visit_CaseInsensitive(t *testing.T) {
	collector := NewCodeBlockCollector("Go") // Uppercase

	// Create a code block with lowercase language
	codeBlock := parser.NewCodeBlock("go", "fmt.Println(\"Hello\")", true, parser.Range{})

	err := collector.Visit(codeBlock)
	require.NoError(t, err)

	blocks := collector.GetCodeBlocks()
	require.Len(t, blocks, 1)
	assert.Equal(t, "go", blocks[0].Language) // Should be normalized to lowercase
}

func TestCodeBlockCollector_Visit_NonCodeBlock(t *testing.T) {
	collector := NewCodeBlockCollector("go")

	// Create a non-code block node
	textNode := parser.NewText("Some text", parser.Range{})

	err := collector.Visit(textNode)
	require.NoError(t, err)

	blocks := collector.GetCodeBlocks()
	assert.Empty(t, blocks) // Non-code blocks should be ignored
}

func TestCodeBlockCollector_GetCodeBlocksByLanguage(t *testing.T) {
	collector := NewCodeBlockCollector("") // Collect all

	// Add blocks with different languages
	goBlock1 := parser.NewCodeBlock("go", "fmt.Println(\"Hello 1\")", true, parser.Range{})
	goBlock2 := parser.NewCodeBlock("go", "fmt.Println(\"Hello 2\")", true, parser.Range{})
	pythonBlock := parser.NewCodeBlock("python", "print(\"Hello\")", true, parser.Range{})

	_ = collector.Visit(goBlock1)
	_ = collector.Visit(goBlock2)
	_ = collector.Visit(pythonBlock)

	// Test getting only Go blocks
	goBlocks := collector.GetCodeBlocksByLanguage("go")
	require.Len(t, goBlocks, 2)
	assert.Equal(t, "go", goBlocks[0].Language)
	assert.Equal(t, "go", goBlocks[1].Language)

	// Test getting only Python blocks
	pythonBlocks := collector.GetCodeBlocksByLanguage("python")
	require.Len(t, pythonBlocks, 1)
	assert.Equal(t, "python", pythonBlocks[0].Language)

	// Test getting non-existent language
	jsBlocks := collector.GetCodeBlocksByLanguage("javascript")
	assert.Empty(t, jsBlocks)
}

func TestCodeBlockCollector_DocumentOrder(t *testing.T) {
	collector := NewCodeBlockCollector("go")

	// Create blocks in specific order
	block1 := parser.NewCodeBlock("go", "// First block", true, parser.Range{
		Start: parser.Position{Line: 1, Column: 1},
	})
	block2 := parser.NewCodeBlock("go", "// Second block", true, parser.Range{
		Start: parser.Position{Line: 5, Column: 1},
	})
	block3 := parser.NewCodeBlock("go", "// Third block", true, parser.Range{
		Start: parser.Position{Line: 10, Column: 1},
	})

	// Visit in order
	_ = collector.Visit(block1)
	_ = collector.Visit(block2)
	_ = collector.Visit(block3)

	blocks := collector.GetCodeBlocks()
	require.Len(t, blocks, 3)

	// Verify order is preserved
	assert.Contains(t, blocks[0].Content, "First")
	assert.Contains(t, blocks[1].Content, "Second")
	assert.Contains(t, blocks[2].Content, "Third")
}
