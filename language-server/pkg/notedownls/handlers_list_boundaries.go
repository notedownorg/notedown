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

	"github.com/notedownorg/notedown/language-server/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/parser"
)

// ListItemFinder helps find list items at specific positions using the parser AST
type ListItemFinder struct {
	targetLine int // 0-based line number we're looking for
	foundItem  *parser.ListItem
}

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

// Visit implements the Visitor interface to find list items at target position
func (f *ListItemFinder) Visit(node parser.Node) error {
	if node.Type() != parser.NodeListItem {
		return nil
	}

	listItem, ok := node.(*parser.ListItem)
	if !ok {
		return nil
	}

	// Check if the target line falls within this list item's range
	// Convert from 1-based to 0-based line numbers from parser
	startLine := listItem.Range().Start.Line - 1
	endLine := listItem.Range().End.Line - 1

	if f.targetLine >= startLine && f.targetLine <= endLine {
		f.foundItem = listItem
	}

	return nil
}

// findListItemAtPosition finds a list item at the given position using the parser AST
func (s *Server) findListItemAtPosition(content string, position lsp.Position) (*parser.ListItem, error) {
	// Use the existing parser to parse the document
	p := parser.NewParser()
	doc, err := p.ParseString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	// Create a finder visitor to locate the list item at the target position
	finder := &ListItemFinder{
		targetLine: position.Line, // position.Line is already 0-based
	}

	// Walk the AST to find the list item
	walker := parser.NewWalker(finder)
	if err := walker.Walk(doc); err != nil {
		return nil, fmt.Errorf("failed to walk AST: %w", err)
	}

	return finder.foundItem, nil
}

// findLastChildLine finds the end line of a list item including all its children
// This traverses the AST to find nested lists that are children of the current item
func (s *Server) findLastChildLine(item *parser.ListItem) int {
	// Start with the item's own end line
	lastLine := item.Range().End.Line - 1 // Convert to 0-based

	// The key insight: nested lists in markdown AST are structured as:
	// ListItem -> [TextBlock/Paragraph] -> [nested List] -> [nested ListItems]
	// We need to find child List nodes and then their ListItem children
	for _, child := range item.Children() {
		childLastLine := s.findLastLineInNode(child)
		if childLastLine > lastLine {
			lastLine = childLastLine
		}
	}

	return lastLine
}

// findLastLineInNode recursively finds the last line in any AST node
func (s *Server) findLastLineInNode(node parser.Node) int {
	// Start with this node's end line
	lastLine := node.Range().End.Line - 1 // Convert to 0-based

	// Recursively check all children
	for _, child := range node.Children() {
		childLastLine := s.findLastLineInNode(child)
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

	// Find the list item at the cursor position using the parser
	item, err := s.findListItemAtPosition(doc.Content, position)
	if err != nil {
		s.logger.Error("failed to find list item at position", "error", err)
		return &BoundaryResponse{Found: false}, nil
	}

	if item == nil {
		s.logger.Debug("no list item found at position for boundaries", "line", position.Line, "character", position.Character)
		return &BoundaryResponse{Found: false}, nil
	}

	// Calculate the boundaries including all children
	startLine := item.Range().Start.Line - 1 // Convert to 0-based
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
