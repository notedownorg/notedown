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
	"fmt"
	"strings"

	"github.com/notedownorg/notedown/language-server/pkg/codeexec"
	"github.com/notedownorg/notedown/pkg/parser"
)

// CodeBlockCollector implements the visitor pattern to collect code blocks
type CodeBlockCollector struct {
	language   string               // target language to collect (empty = all)
	codeBlocks []codeexec.CodeBlock // collected code blocks
}

// NewCodeBlockCollector creates a new code block collector
func NewCodeBlockCollector(language string) *CodeBlockCollector {
	return &CodeBlockCollector{
		language:   strings.ToLower(strings.TrimSpace(language)),
		codeBlocks: make([]codeexec.CodeBlock, 0),
	}
}

// Visit implements the visitor pattern to collect code blocks
func (c *CodeBlockCollector) Visit(node parser.Node) error {
	if codeBlock, ok := node.(*parser.CodeBlock); ok {
		// Only collect fenced code blocks (not indented code blocks)
		if codeBlock.Fenced {
			blockLanguage := strings.ToLower(strings.TrimSpace(codeBlock.Language))

			// If no specific language requested, collect all fenced code blocks
			// If specific language requested, only collect matching blocks
			if c.language == "" || blockLanguage == c.language {
				c.codeBlocks = append(c.codeBlocks, codeexec.CodeBlock{
					Language: blockLanguage,
					Content:  codeBlock.Content,
					Range:    codeBlock.Range(),
				})
			}
		}
	}

	return nil
}

// GetCodeBlocks returns the collected code blocks
func (c *CodeBlockCollector) GetCodeBlocks() []codeexec.CodeBlock {
	return c.codeBlocks
}

// GetCodeBlocksByLanguage returns code blocks filtered by language
func (c *CodeBlockCollector) GetCodeBlocksByLanguage(language string) []codeexec.CodeBlock {
	targetLanguage := strings.ToLower(strings.TrimSpace(language))
	var filtered []codeexec.CodeBlock

	for _, block := range c.codeBlocks {
		if block.Language == targetLanguage {
			filtered = append(filtered, block)
		}
	}

	return filtered
}

// collectCodeBlocks is a helper function to collect code blocks from a document
func (s *Server) collectCodeBlocks(uri, language string) ([]codeexec.CodeBlock, error) {
	// Get the document
	doc, exists := s.GetDocument(uri)
	if !exists {
		return nil, fmt.Errorf("document not found: %s", uri)
	}

	// Parse the document content to extract code blocks
	blocks, err := s.parseCodeBlocksFromContent(doc.Content, language)
	if err != nil {
		s.logger.Error("failed to parse code blocks from document", "uri", uri, "error", err)
		return nil, err
	}

	s.logger.Debug("collected code blocks", "uri", uri, "language", language, "count", len(blocks))

	return blocks, nil
}

// parseCodeBlocksFromContent parses code blocks from document content using regex
func (s *Server) parseCodeBlocksFromContent(content, targetLanguage string) ([]codeexec.CodeBlock, error) {
	var blocks []codeexec.CodeBlock
	lines := strings.Split(content, "\n")

	var inCodeBlock bool
	var currentLanguage string
	var currentContent strings.Builder
	var blockStartLine int

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check for fenced code block start
		if strings.HasPrefix(trimmedLine, "```") && !inCodeBlock {
			// Starting a new code block
			inCodeBlock = true
			blockStartLine = i

			// Extract language identifier
			currentLanguage = strings.ToLower(strings.TrimSpace(trimmedLine[3:]))
			currentContent.Reset()
			continue
		}

		// Check for fenced code block end
		if strings.HasPrefix(trimmedLine, "```") && inCodeBlock {
			// Ending code block
			inCodeBlock = false

			// Only collect if language matches (or collecting all languages)
			if targetLanguage == "" || currentLanguage == targetLanguage {
				blocks = append(blocks, codeexec.CodeBlock{
					Language: currentLanguage,
					Content:  currentContent.String(),
					Range: parser.Range{
						Start: parser.Position{Line: blockStartLine + 1, Column: 1},
						End:   parser.Position{Line: i + 1, Column: 1},
					},
				})
			}

			currentLanguage = ""
			currentContent.Reset()
			continue
		}

		// If we're inside a code block, collect the content
		if inCodeBlock {
			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(line)
		}
	}

	return blocks, nil
}
