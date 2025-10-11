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
	"strings"

	"github.com/notedownorg/notedown/pkg/parser/extensions"
	"github.com/yuin/goldmark/ast"
)

// WikilinkInfo contains information about an extracted wikilink
type WikilinkInfo struct {
	Target      string // The wikilink target (e.g., "project-alpha", "docs/api")
	DisplayText string // Display text if pipe notation is used
	Line        int    // 1-based line number
	Column      int    // 1-based column number
}

// TaskInfo contains information about an extracted task
type TaskInfo struct {
	State  string // Raw state value (e.g., " ", "x", "wip", "in-progress")
	Text   string // Task description text
	Line   int    // 1-based line number
	Column int    // 1-based column number
}

// ExtractWikilinks extracts all wikilinks from a parsed document
func ExtractWikilinks(doc *Document) []WikilinkInfo {
	var wikilinks []WikilinkInfo

	walker := NewWalker(WalkFunc(func(node Node) error {
		if wikilink, ok := node.(*Wikilink); ok {
			info := WikilinkInfo{
				Target:      strings.TrimSpace(wikilink.Target),
				DisplayText: wikilink.DisplayText,
				Line:        wikilink.Range().Start.Line,
				Column:      wikilink.Range().Start.Column,
			}
			wikilinks = append(wikilinks, info)
		}
		return nil
	}))

	if err := walker.Walk(doc); err != nil {
		// Log error but don't fail extraction
		return wikilinks
	}

	return wikilinks
}

// ExtractTasks extracts all tasks from a parsed document
// This function works with the goldmark AST since tasks are handled at the goldmark level
func ExtractTasks(content []byte, markdownAST ast.Node) []TaskInfo {
	var tasks []TaskInfo

	ast.Walk(markdownAST, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if taskBox, ok := node.(*extensions.TaskCheckBox); ok {
			// Find the parent list item to extract the task text
			var taskText string
			if parent := taskBox.Parent(); parent != nil {
				// Walk the siblings after the checkbox to get the text
				for sibling := taskBox.NextSibling(); sibling != nil; sibling = sibling.NextSibling() {
					if text, ok := sibling.(*ast.Text); ok {
						taskText += string(text.Text(content))
					} else if textContent := extractTextFromNode(sibling, content); textContent != "" {
						taskText += textContent
					}
				}
			}

			// Calculate position from segment
			line, column := 1, 1
			if taskBox.HasChildren() {
				// Use the first segment to determine position
				if segment := taskBox.Lines().At(0); segment.Len() > 0 {
					line, column = calculateLineColumn(content, segment.Start)
				}
			}

			task := TaskInfo{
				State:  taskBox.State,
				Text:   strings.TrimSpace(taskText),
				Line:   line,
				Column: column,
			}
			tasks = append(tasks, task)
		}

		return ast.WalkContinue, nil
	})

	return tasks
}

// extractTextFromNode extracts text content from an AST node recursively
func extractTextFromNode(node ast.Node, content []byte) string {
	var result string

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if text, ok := n.(*ast.Text); ok {
			result += string(text.Text(content))
		}

		return ast.WalkContinue, nil
	})

	return result
}

// calculateLineColumn calculates 1-based line and column from byte offset
func calculateLineColumn(content []byte, offset int) (int, int) {
	if offset < 0 || offset >= len(content) {
		return 1, 1
	}

	line, column := 1, 1
	for i := 0; i < offset && i < len(content); i++ {
		if content[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}

	return line, column
}
