package notedownls

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/parser"
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

// handleExecuteCommand handles workspace/executeCommand requests
func (s *Server) handleExecuteCommand(params json.RawMessage) (any, error) {
	var executeParams lsp.ExecuteCommandParams
	if err := json.Unmarshal(params, &executeParams); err != nil {
		s.logger.Error("failed to unmarshal execute command params", "error", err)
		return nil, err
	}

	s.logger.Debug("execute command request received", "command", executeParams.Command)

	switch executeParams.Command {
	case "notedown.moveListItemUp":
		return s.handleMoveListItemUp(executeParams.Arguments)
	case "notedown.moveListItemDown":
		return s.handleMoveListItemDown(executeParams.Arguments)
	default:
		return nil, fmt.Errorf("unknown command: %s", executeParams.Command)
	}
}

// handleMoveListItemUp moves a list item and its children up
func (s *Server) handleMoveListItemUp(arguments []any) (any, error) {
	return s.handleMoveListItem(arguments, true)
}

// handleMoveListItemDown moves a list item and its children down
func (s *Server) handleMoveListItemDown(arguments []any) (any, error) {
	return s.handleMoveListItem(arguments, false)
}

// handleMoveListItem handles the core logic for moving list items
func (s *Server) handleMoveListItem(arguments []any, moveUp bool) (any, error) {
	if len(arguments) < 2 {
		return nil, fmt.Errorf("moveListItem requires document URI and position arguments")
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

	s.logger.Debug("moving list item", "uri", documentURI, "line", position.Line, "character", position.Character, "up", moveUp)

	// Get the document
	doc, exists := s.GetDocument(documentURI)
	if !exists {
		s.logger.Error("document not found for list movement", "uri", documentURI)
		return nil, fmt.Errorf("document not found: %s", documentURI)
	}

	s.logger.Debug("found document for list movement", "uri", documentURI, "content_length", len(doc.Content))

	// Validate document structure using parser library
	parser := parser.NewParser()
	parsedDoc, err := parser.ParseString(doc.Content)
	if err != nil {
		s.logger.Error("failed to parse document", "error", err)
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	// Verify that there's a list item at the cursor position
	listItem := parsedDoc.FindListItemAtLine(int(position.Line) + 1)
	if listItem == nil {
		s.logger.Debug("no list item found at position", "line", position.Line, "character", position.Character)
		return nil, fmt.Errorf("no list item found at position %d:%d", position.Line, position.Character)
	}

	s.logger.Debug("validated list item at position", "line", position.Line, "range", listItem.Range())

	// Use existing regex-based parsing for text editing operations
	hierarchy, err := s.parseListHierarchy(doc.Content)
	if err != nil {
		s.logger.Error("failed to parse list hierarchy", "error", err)
		return nil, fmt.Errorf("failed to parse list hierarchy: %w", err)
	}

	// Find the list item at the cursor position using the existing logic
	item := hierarchy.findItemAtPosition(position)
	if item == nil {
		s.logger.Debug("no list item found at position in hierarchy", "line", position.Line, "character", position.Character)
		return nil, fmt.Errorf("no list item found at position %d:%d", position.Line, position.Character)
	}

	// Find the target position for the move
	workspaceEdit, err := s.calculateListItemMove(hierarchy, item, moveUp, documentURI)
	if err != nil {
		s.logger.Error("failed to calculate list item move", "error", err)
		return nil, fmt.Errorf("failed to calculate list item move: %w", err)
	}

	if workspaceEdit == nil {
		// No move possible (already at boundary)
		s.logger.Debug("list item already at boundary, no move performed")
		return nil, fmt.Errorf("cannot move list item: already at boundary")
	}

	s.logger.Debug("calculated workspace edit", "changes_count", len(workspaceEdit.Changes[documentURI]))
	for i, edit := range workspaceEdit.Changes[documentURI] {
		s.logger.Debug("text edit", "index", i, "start_line", edit.Range.Start.Line, "end_line", edit.Range.End.Line, "new_text_length", len(edit.NewText))
	}

	// Return the workspace edit for the client to apply
	// This is the standard way for executeCommand to handle workspace edits
	s.logger.Debug("returning workspace edit for client to apply")
	return workspaceEdit, nil
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
			if len(currentStack) > 0 {
				// Add to the last item in the stack
				lastItem := currentStack[len(currentStack)-1]
				lastItem.EndLine = lineNum
				lastItem.OriginalLines = append(lastItem.OriginalLines, line)
			}
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

// calculateListItemMove calculates the workspace edit needed to move a list item
func (s *Server) calculateListItemMove(hierarchy *ListHierarchy, item *ListItem, moveUp bool, documentURI string) (*lsp.WorkspaceEdit, error) {
	// Find the parent container and sibling items
	parentItems, targetIndex := s.findParentAndIndex(hierarchy, item)
	if targetIndex == -1 {
		return nil, fmt.Errorf("could not find item in hierarchy")
	}

	// Check if move is possible
	if moveUp && targetIndex == 0 {
		// Already at the top
		return nil, nil
	}
	if !moveUp && targetIndex == len(parentItems)-1 {
		// Already at the bottom
		return nil, nil
	}

	// Calculate new position
	var swapIndex int
	if moveUp {
		swapIndex = targetIndex - 1
	} else {
		swapIndex = targetIndex + 1
	}

	// Get the item to swap with
	swapItem := parentItems[swapIndex]

	// Create text edits to swap the items
	textEdits := s.createSwapTextEdits(item, swapItem)

	// Handle renumbering for ordered lists
	textEdits = s.handleOrderedListRenumbering(textEdits, parentItems, targetIndex, swapIndex)

	workspaceEdit := &lsp.WorkspaceEdit{
		Changes: map[string][]lsp.TextEdit{
			documentURI: textEdits,
		},
	}

	return workspaceEdit, nil
}

// findParentAndIndex finds the parent container and index of the given item
func (s *Server) findParentAndIndex(hierarchy *ListHierarchy, targetItem *ListItem) ([]*ListItem, int) {
	// Check top-level items first
	for i, item := range hierarchy.Items {
		if item == targetItem {
			return hierarchy.Items, i
		}
	}

	// Search recursively in children
	var searchInChildren func([]*ListItem) ([]*ListItem, int)
	searchInChildren = func(items []*ListItem) ([]*ListItem, int) {
		for i, item := range items {
			if item == targetItem {
				return items, i
			}
			if parentItems, index := searchInChildren(item.Children); index != -1 {
				return parentItems, index
			}
		}
		return nil, -1
	}

	for _, item := range hierarchy.Items {
		if parentItems, index := searchInChildren(item.Children); index != -1 {
			return parentItems, index
		}
	}

	return nil, -1
}

// createSwapTextEdits creates text edits to swap two list items
func (s *Server) createSwapTextEdits(item1, item2 *ListItem) []lsp.TextEdit {
	var edits []lsp.TextEdit

	// Get the full text ranges for both items (including all children)
	item1Range := s.getItemFullRange(item1)
	item2Range := s.getItemFullRange(item2)

	// Get the full text content for both items
	item1Text := strings.Join(item1.OriginalLines, "\n")
	item2Text := strings.Join(item2.OriginalLines, "\n")

	// Add all children text for item1
	item1Text += s.getChildrenText(item1)
	// Add all children text for item2
	item2Text += s.getChildrenText(item2)

	// Create edits to swap the content
	if item1Range.Start.Line < item2Range.Start.Line {
		// item1 comes before item2
		edits = append(edits, lsp.TextEdit{
			Range:   item1Range,
			NewText: item2Text,
		})
		edits = append(edits, lsp.TextEdit{
			Range:   item2Range,
			NewText: item1Text,
		})
	} else {
		// item2 comes before item1
		edits = append(edits, lsp.TextEdit{
			Range:   item2Range,
			NewText: item1Text,
		})
		edits = append(edits, lsp.TextEdit{
			Range:   item1Range,
			NewText: item2Text,
		})
	}

	return edits
}

// getItemFullRange returns the full range of a list item including all its children
func (s *Server) getItemFullRange(item *ListItem) lsp.Range {
	startLine := item.StartLine
	endLine := item.EndLine

	// Find the last line of the deepest child
	endLine = s.findLastChildLine(item)

	return lsp.Range{
		Start: lsp.Position{Line: startLine, Character: 0},
		End:   lsp.Position{Line: endLine + 1, Character: 0}, // Include the line after for proper replacement
	}
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

// getChildrenText returns the text content of all children of a list item
func (s *Server) getChildrenText(item *ListItem) string {
	var text strings.Builder

	var traverse func(*ListItem)
	traverse = func(item *ListItem) {
		for _, child := range item.Children {
			text.WriteString("\n")
			text.WriteString(strings.Join(child.OriginalLines, "\n"))
			traverse(child)
		}
	}

	traverse(item)
	return text.String()
}

// handleOrderedListRenumbering updates the numbering for ordered lists after a move
func (s *Server) handleOrderedListRenumbering(edits []lsp.TextEdit, parentItems []*ListItem, index1, index2 int) []lsp.TextEdit {
	// Check if we're dealing with numbered lists
	if len(parentItems) == 0 {
		return edits
	}

	// Check if ALL items in this specific level are numbered lists
	// We only renumber if we're moving within a numbered list
	allNumbered := true
	for _, item := range parentItems {
		if matched, _ := regexp.MatchString(`^\d+\.`, item.Marker); !matched {
			allNumbered = false
			break
		}
	}

	if !allNumbered {
		return edits
	}

	s.logger.Debug("renumbering ordered list after move", "parent_items_count", len(parentItems))

	// Renumber all items in this level, accounting for the swap
	for i, item := range parentItems {
		// Calculate what the new number should be after the swap
		var newNumber int
		if i == index1 {
			// This item is moving to index2's position
			newNumber = index2 + 1
		} else if i == index2 {
			// This item is moving to index1's position
			newNumber = index1 + 1
		} else {
			// This item stays in place
			newNumber = i + 1
		}

		newMarker := strconv.Itoa(newNumber) + "."

		// Update the marker in the first line of the item
		if len(item.OriginalLines) > 0 {
			oldLine := item.OriginalLines[0]
			newLine := regexp.MustCompile(`^\s*\d+\.`).ReplaceAllString(oldLine, strings.Repeat(" ", item.IndentLevel)+newMarker)

			// Only add edit if the line actually changes
			if newLine != oldLine {
				s.logger.Debug("updating numbered list marker", "line", item.StartLine, "old_marker", item.Marker, "new_marker", newMarker)
				edits = append(edits, lsp.TextEdit{
					Range: lsp.Range{
						Start: lsp.Position{Line: item.StartLine, Character: 0},
						End:   lsp.Position{Line: item.StartLine, Character: len(oldLine)},
					},
					NewText: newLine,
				})
			}
		}
	}

	return edits
}