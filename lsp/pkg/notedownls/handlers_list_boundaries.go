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
	"fmt"
	"regexp"
	"strings"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
)

// ListItem represents a markdown list item with its position and content
type ListItem struct {
	// The line number where this list item starts (0-based)
	StartLine int
	// The line number where this list item ends (0-based, inclusive)
	EndLine int
	// The indentation level (number of spaces or tabs)
	IndentLevel int
	// The list marker (e.g., "-", "*", "1.", "- [ ]", "- [x]")
	Marker string
	// The content of the list item (without marker and indentation)
	Content string
	// Child list items
	Children []*ListItem
	// The original text lines for this item and its children
	OriginalLines []string
}

// ListHierarchy represents the complete list structure of a document
type ListHierarchy struct {
	// All list items found in the document
	Items []*ListItem
	// Map from line number to list item for quick lookup
	LineToItem map[int]*ListItem
}

// Regular expressions for different list item types
var (
	// Matches task lists: - [ ] or - [x] or - [X]
	taskListRegex = regexp.MustCompile(`^(\s*)(- \[[xX ]?\])(.*)$`)
	// Matches bullet lists: - or *
	bulletListRegex = regexp.MustCompile(`^(\s*)([-*])(\s+.*)$`)
	// Matches numbered lists: 1. or 123.
	numberedListRegex = regexp.MustCompile(`^(\s*)(\d+\.)(\s+.*)$`)
)

// handleExecuteCommand handles workspace/executeCommand requests for list boundaries
func (s *Server) handleExecuteCommand(params json.RawMessage) (any, error) {
	var executeParams lsp.ExecuteCommandParams
	if err := json.Unmarshal(params, &executeParams); err != nil {
		s.logger.Error("failed to unmarshal execute command params", "error", err)
		return nil, err
	}

	s.logger.Debug("execute command request received", "command", executeParams.Command)

	switch executeParams.Command {
	case "notedown.getListItemBoundaries":
		return s.handleGetListItemBoundaries(executeParams.Arguments)
	default:
		return nil, fmt.Errorf("unknown command: %s", executeParams.Command)
	}
}

// parseListHierarchy parses the document content and builds a list hierarchy
func (s *Server) parseListHierarchy(content string) (*ListHierarchy, error) {
	lines := strings.Split(content, "\n")
	hierarchy := &ListHierarchy{
		Items:      make([]*ListItem, 0),
		LineToItem: make(map[int]*ListItem),
	}

	var currentStack []*ListItem // Stack to track nesting levels

	for lineNum, line := range lines {
		item := s.parseListItem(line, lineNum)
		if item == nil {
			// Not a list item, but might be continuation of previous item
			// Only include lines that are indented more than the current list item
			// or blank lines immediately following a list item
			if len(currentStack) > 0 {
				lastItem := currentStack[len(currentStack)-1]
				trimmedLine := strings.TrimSpace(line)

				// Only include as continuation if:
				// 1. It's indented more than the list item (content continuation), OR
				// 2. It's not a heading or other structural element
				// Note: We're being more restrictive about blank lines to avoid including too much
				if len(line) > lastItem.IndentLevel &&
					len(line)-len(strings.TrimLeft(line, " \t")) > lastItem.IndentLevel &&
					!strings.HasPrefix(trimmedLine, "#") {
					lastItem.EndLine = lineNum
					lastItem.OriginalLines = append(lastItem.OriginalLines, line)
					continue
				}
			}
			// If we reach here, this line is not part of any list item
			// Clear the stack as we've moved past the list
			currentStack = currentStack[:0]
			continue
		}

		// Determine where this item fits in the hierarchy
		item.OriginalLines = []string{line}
		item.EndLine = lineNum

		// Find the appropriate parent level
		var parent *ListItem
		for len(currentStack) > 0 {
			potentialParent := currentStack[len(currentStack)-1]
			if item.IndentLevel > potentialParent.IndentLevel {
				// This item is a child of the potential parent
				parent = potentialParent
				break
			}
			// Pop from stack - this item is at same or higher level
			currentStack = currentStack[:len(currentStack)-1]
		}

		// Add to parent's children or top-level items
		if parent != nil {
			parent.Children = append(parent.Children, item)
		} else {
			hierarchy.Items = append(hierarchy.Items, item)
		}

		// Add to line mapping
		hierarchy.LineToItem[lineNum] = item

		// Push to stack
		currentStack = append(currentStack, item)
	}

	return hierarchy, nil
}

// parseListItem parses a single line to determine if it's a list item
func (s *Server) parseListItem(line string, lineNum int) *ListItem {
	// Try task list first
	if matches := taskListRegex.FindStringSubmatch(line); matches != nil {
		return &ListItem{
			StartLine:   lineNum,
			EndLine:     lineNum,
			IndentLevel: len(matches[1]),
			Marker:      matches[2],
			Content:     strings.TrimSpace(matches[3]),
		}
	}

	// Try bullet list
	if matches := bulletListRegex.FindStringSubmatch(line); matches != nil {
		return &ListItem{
			StartLine:   lineNum,
			EndLine:     lineNum,
			IndentLevel: len(matches[1]),
			Marker:      matches[2],
			Content:     strings.TrimSpace(matches[3]),
		}
	}

	// Try numbered list
	if matches := numberedListRegex.FindStringSubmatch(line); matches != nil {
		return &ListItem{
			StartLine:   lineNum,
			EndLine:     lineNum,
			IndentLevel: len(matches[1]),
			Marker:      matches[2],
			Content:     strings.TrimSpace(matches[3]),
		}
	}

	return nil
}

// findItemAtPosition finds the list item at the given position
func (h *ListHierarchy) findItemAtPosition(position lsp.Position) *ListItem {
	item, exists := h.LineToItem[position.Line]
	if exists {
		return item
	}

	// Check if position is within a multi-line list item
	for _, item := range h.getAllItems() {
		if position.Line >= item.StartLine && position.Line <= item.EndLine {
			return item
		}
	}

	return nil
}

// getAllItems returns all list items in the hierarchy (flattened)
func (h *ListHierarchy) getAllItems() []*ListItem {
	var allItems []*ListItem

	var traverse func([]*ListItem)
	traverse = func(items []*ListItem) {
		for _, item := range items {
			allItems = append(allItems, item)
			traverse(item.Children)
		}
	}

	traverse(h.Items)
	return allItems
}

// findLastChildLine recursively finds the last line of a list item and all its children
func (s *Server) findLastChildLine(item *ListItem) int {
	lastLine := item.EndLine

	for _, child := range item.Children {
		childLastLine := s.findLastChildLine(child)
		if childLastLine > lastLine {
			lastLine = childLastLine
		}
	}

	return lastLine
}

// BoundaryResponse represents the response for list item boundary requests
type BoundaryResponse struct {
	Start lsp.Position `json:"start"`
	End   lsp.Position `json:"end"`
	Found bool         `json:"found"`
}

// handleGetListItemBoundaries returns the boundaries of a list item and all its children
func (s *Server) handleGetListItemBoundaries(arguments []any) (any, error) {
	if len(arguments) < 2 {
		return nil, fmt.Errorf("getListItemBoundaries requires document URI and position arguments")
	}

	// Extract document URI
	documentURI, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("first argument must be document URI (string)")
	}

	// Extract position
	positionMap, ok := arguments[1].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("second argument must be position object")
	}

	line, ok := positionMap["line"].(float64)
	if !ok {
		return nil, fmt.Errorf("position must have line number")
	}

	character, ok := positionMap["character"].(float64)
	if !ok {
		return nil, fmt.Errorf("position must have character number")
	}

	position := lsp.Position{
		Line:      int(line),
		Character: int(character),
	}

	s.logger.Debug("getting list item boundaries", "uri", documentURI, "line", position.Line, "character", position.Character)

	// Get the document
	doc, exists := s.GetDocument(documentURI)
	if !exists {
		s.logger.Error("document not found for boundary request", "uri", documentURI)
		return &BoundaryResponse{Found: false}, nil
	}

	// Parse list hierarchy
	hierarchy, err := s.parseListHierarchy(doc.Content)
	if err != nil {
		s.logger.Error("failed to parse list hierarchy for boundaries", "error", err)
		return &BoundaryResponse{Found: false}, nil
	}

	// Find the list item at the cursor position
	item := hierarchy.findItemAtPosition(position)
	if item == nil {
		s.logger.Debug("no list item found at position for boundaries", "line", position.Line, "character", position.Character)
		return &BoundaryResponse{Found: false}, nil
	}

	// Calculate the boundaries including all children
	startLine := item.StartLine
	endLine := s.findLastChildLine(item)

	// Create boundary response
	response := &BoundaryResponse{
		Start: lsp.Position{
			Line:      startLine,
			Character: 0, // Start at beginning of line
		},
		End: lsp.Position{
			Line:      endLine + 1, // Include next line for proper text object behavior
			Character: 0,
		},
		Found: true,
	}

	s.logger.Debug("calculated list item boundaries",
		"start_line", response.Start.Line,
		"end_line", response.End.Line,
		"item_start", startLine,
		"item_end", endLine)

	return response, nil
}