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
	"regexp"
	"strings"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
)

// handleFoldingRange handles textDocument/foldingRange requests
func (s *Server) handleFoldingRange(params json.RawMessage) (any, error) {
	var foldingParams lsp.FoldingRangeParams
	if err := json.Unmarshal(params, &foldingParams); err != nil {
		s.logger.Error("failed to unmarshal folding range params", "error", err)
		return nil, err
	}

	s.logger.Debug("folding range request received", "uri", foldingParams.TextDocument.URI)

	// Get the document
	doc, exists := s.GetDocument(foldingParams.TextDocument.URI)
	if !exists {
		s.logger.Debug("document not found for folding range", "uri", foldingParams.TextDocument.URI)
		return []lsp.FoldingRange{}, nil
	}

	// Generate folding ranges
	ranges := s.generateFoldingRanges(doc.Content)
	s.logger.Debug("generated folding ranges", "count", len(ranges))

	return ranges, nil
}

// generateFoldingRanges analyzes document content and generates folding ranges
func (s *Server) generateFoldingRanges(content string) []lsp.FoldingRange {
	var ranges []lsp.FoldingRange
	lines := strings.Split(content, "\n")

	// Generate header folding ranges
	ranges = append(ranges, s.generateHeaderFoldingRanges(lines)...)

	// Generate list folding ranges
	ranges = append(ranges, s.generateListFoldingRanges(lines)...)

	// Generate code block folding ranges
	ranges = append(ranges, s.generateCodeBlockFoldingRanges(lines)...)

	return ranges
}

// generateHeaderFoldingRanges creates folding ranges for markdown headers
func (s *Server) generateHeaderFoldingRanges(lines []string) []lsp.FoldingRange {
	var ranges []lsp.FoldingRange
	headerRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)`)

	for i, line := range lines {
		if match := headerRegex.FindStringSubmatch(line); match != nil {
			headerLevel := len(match[1])
			startLine := i

			// Find the end of this header section
			endLine := len(lines) - 1
			for j := i + 1; j < len(lines); j++ {
				if nextMatch := headerRegex.FindStringSubmatch(lines[j]); nextMatch != nil {
					nextLevel := len(nextMatch[1])
					if nextLevel <= headerLevel {
						endLine = j - 1
						break
					}
				}
			}

			// Only create folding range if there's content to fold
			if endLine > startLine {
				regionKind := lsp.FoldingRangeKindRegion
				ranges = append(ranges, lsp.FoldingRange{
					StartLine: startLine,
					EndLine:   endLine,
					Kind:      &regionKind,
				})
			}
		}
	}

	return ranges
}

// generateListFoldingRanges creates folding ranges for nested lists and task lists
func (s *Server) generateListFoldingRanges(lines []string) []lsp.FoldingRange {
	var ranges []lsp.FoldingRange
	listItemRegex := regexp.MustCompile(`^(\s*)([-*+]|\d+\.)\s+(.*)`)
	taskItemRegex := regexp.MustCompile(`^(\s*)-\s+\[[^\]]*\]\s+(.*)`)

	for i, line := range lines {
		var match []string
		var isTask bool

		// Check for task list first, then regular list
		if taskMatch := taskItemRegex.FindStringSubmatch(line); taskMatch != nil {
			match = []string{taskMatch[0], taskMatch[1], "-", taskMatch[2]}
			isTask = true
		} else if listMatch := listItemRegex.FindStringSubmatch(line); listMatch != nil {
			match = listMatch
			isTask = false
		}

		if match != nil {
			indent := len(match[1])
			startLine := i

			// Find the end of this list item (including nested items)
			endLine := startLine
			for j := i + 1; j < len(lines); j++ {
				nextLine := lines[j]

				// Empty lines continue the list item
				if strings.TrimSpace(nextLine) == "" {
					continue
				}

				var nextMatch []string
				if nextTaskMatch := taskItemRegex.FindStringSubmatch(nextLine); nextTaskMatch != nil {
					nextMatch = []string{nextTaskMatch[0], nextTaskMatch[1], "-", nextTaskMatch[2]}
				} else if nextListMatch := listItemRegex.FindStringSubmatch(nextLine); nextListMatch != nil {
					nextMatch = nextListMatch
				}

				if nextMatch != nil {
					nextIndent := len(nextMatch[1])
					// If next item is at same or lower level, this item ends
					if nextIndent <= indent {
						endLine = j - 1
						break
					}
				} else {
					// Check if this line is continuation of the list item (indented properly)
					if len(nextLine) > 0 && nextLine[0] == ' ' {
						lineIndent := len(nextLine) - len(strings.TrimLeft(nextLine, " "))
						if lineIndent > indent {
							// This is continuation content
							continue
						}
					}
					// Non-list content at same or lower indent level ends the list item
					endLine = j - 1
					break
				}
			}

			// Set endLine to last line if we reached end of document
			if endLine == startLine {
				endLine = len(lines) - 1
			}

			// Only create folding range if there's content to fold
			if endLine > startLine {
				var kind *lsp.FoldingRangeKind
				if isTask {
					regionKind := lsp.FoldingRangeKindRegion
					kind = &regionKind
				}

				ranges = append(ranges, lsp.FoldingRange{
					StartLine: startLine,
					EndLine:   endLine,
					Kind:      kind,
				})
			}
		}
	}

	return ranges
}

// generateCodeBlockFoldingRanges creates folding ranges for code blocks
func (s *Server) generateCodeBlockFoldingRanges(lines []string) []lsp.FoldingRange {
	var ranges []lsp.FoldingRange
	codeBlockRegex := regexp.MustCompile(`^` + "```" + `(\w+)?`)

	inCodeBlock := false
	startLine := 0

	for i, line := range lines {
		if codeBlockRegex.MatchString(line) {
			if !inCodeBlock {
				// Start of code block
				inCodeBlock = true
				startLine = i
			} else {
				// End of code block
				inCodeBlock = false
				if i > startLine {
					regionKind := lsp.FoldingRangeKindRegion
					ranges = append(ranges, lsp.FoldingRange{
						StartLine: startLine,
						EndLine:   i,
						Kind:      &regionKind,
					})
				}
			}
		}
	}

	return ranges
}
