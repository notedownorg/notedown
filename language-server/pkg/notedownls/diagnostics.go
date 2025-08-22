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
	"regexp"
	"strings"

	"github.com/notedownorg/notedown/language-server/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/config"
	"github.com/notedownorg/notedown/pkg/parser"
)

// generateWikilinkDiagnostics generates diagnostics for wikilink conflicts in a document
func (s *Server) generateWikilinkDiagnostics(uri, content string) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic

	// Regular expression to find wikilinks and their positions
	wikilinkRegex := regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	matches := wikilinkRegex.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		// Extract the target from the match
		targetStart := match[2]
		targetEnd := match[3]
		if targetStart == -1 || targetEnd == -1 {
			continue
		}

		target := content[targetStart:targetEnd]
		target = strings.TrimSpace(target)

		// Get target info from index
		allTargets := s.wikilinkIndex.GetAllTargets()
		targetInfo, exists := allTargets[target]

		if exists && targetInfo.IsAmbiguous {
			// Calculate line and character positions
			line, char := s.positionFromOffset(content, match[0])
			endLine, endChar := s.positionFromOffset(content, match[1])

			// Create diagnostic for ambiguous wikilink
			severity := lsp.DiagnosticSeverityWarning
			source := "notedown"
			message := fmt.Sprintf("Ambiguous wikilink '%s' matches multiple files: %s",
				target, strings.Join(targetInfo.MatchingFiles, ", "))

			diagnostic := lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: line, Character: char},
					End:   lsp.Position{Line: endLine, Character: endChar},
				},
				Severity: &severity,
				Source:   &source,
				Message:  message,
				Code:     "ambiguous-wikilink",
			}

			// Add related information for each matching file
			for _, filePath := range targetInfo.MatchingFiles {
				fileURI := "file://" + filePath
				diagnostic.RelatedInformation = append(diagnostic.RelatedInformation,
					lsp.DiagnosticRelatedInformation{
						Location: lsp.Location{
							URI: fileURI,
							Range: lsp.Range{
								Start: lsp.Position{Line: 0, Character: 0},
								End:   lsp.Position{Line: 0, Character: 0},
							},
						},
						Message: fmt.Sprintf("Matches file: %s", filePath),
					})
			}

			diagnostics = append(diagnostics, diagnostic)
		}
	}

	return diagnostics
}

// generateTaskDiagnostics generates diagnostics for invalid task states in a document
//
// Uses a hybrid approach combining regex pattern matching with AST validation:
// 1. Regex finds all potential task checkbox patterns in list items
// 2. Parser identifies which list items are actually recognized as valid tasks
// 3. Cross-reference: checkboxes that look like tasks but aren't valid = diagnostics
//
// This catches malformed tasks like [invalid] that the parser would ignore,
// while avoiding false positives on regular text that contains [brackets].
func (s *Server) generateTaskDiagnostics(uri, content string) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic

	// Load workspace configuration for task state validation
	cfg, err := s.loadWorkspaceConfig()
	if err != nil {
		s.logger.Error("failed to load workspace config for task diagnostics", "error", err)
		return diagnostics
	}

	// Parse the document to get the AST
	p := parser.NewParser()
	doc, err := p.ParseString(content)
	if err != nil {
		s.logger.Error("failed to parse document for task diagnostics", "error", err)
		return diagnostics
	}

	// Find invalid task checkboxes using hybrid regex+parser approach
	diagnostics = append(diagnostics, s.findInvalidTaskCheckboxes(doc, content, cfg)...)

	return diagnostics
}

// findInvalidTaskCheckboxes finds potential task checkboxes that failed validation
//
// This implements the hybrid approach:
// 1. Regex finds all checkbox patterns that LOOK like tasks in list items
// 2. Parser AST tells us which list items are ACTUALLY recognized as valid tasks
// 3. Invalid = checkbox pattern found + not recognized as valid task + invalid state
//
// Examples:
//   - [x] task        ← Regex finds, Parser recognizes as task → No diagnostic
//   - [invalid] task  ← Regex finds, Parser ignores (not task) → Diagnostic
//     Text [invalid]    ← Regex ignores (not in list) → No diagnostic
func (s *Server) findInvalidTaskCheckboxes(doc *parser.Document, content string, cfg *config.Config) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic

	// Step 1: Use regex to find all potential task checkbox patterns in list items
	// Pattern: line start + list marker (-, *, +, or number.) + whitespace + [text] + whitespace
	// The trailing whitespace requirement distinguishes task checkboxes from wikilinks/markdown links
	taskRegex := regexp.MustCompile(`(?m)^(\s*[-*+]|\s*\d+\.)\s*(\[([^\]]*)\])\s`)
	matches := taskRegex.FindAllStringSubmatchIndex(content, -1)

	if len(matches) == 0 {
		return diagnostics
	}

	// Step 2: Walk the AST to get all list items and their parser-determined task status
	listItems := s.collectListItems(doc)

	// Step 3: Cross-reference regex matches with parser results
	for _, match := range matches {
		// Extract the checkbox state from the regex match
		stateStart := match[6] // Start of the state content inside [] (group 3)
		stateEnd := match[7]   // End of the state content inside [] (group 3)
		if stateStart == -1 || stateEnd == -1 {
			continue
		}

		state := content[stateStart:stateEnd]
		checkboxStart := match[4] // Start of the '[' character (group 2)
		checkboxEnd := match[5]   // End of the ']' character + 1 (group 2)

		// Convert byte offsets to LSP line/character positions
		line, char := s.positionFromOffset(content, checkboxStart)
		endLine, endChar := s.positionFromOffset(content, checkboxEnd)

		// Check if the parser recognized this checkbox as a valid task
		// The parser only sets TaskList=true for list items with valid task states
		isValidTask := false
		for _, item := range listItems {
			// Check if this checkbox is within this list item's range
			if s.isPositionInRange(line, char, item.Range()) {
				// If the parser marked this list item as a task, the checkbox was valid
				if item.TaskList {
					isValidTask = true
				}
				break
			}
		}

		// Create diagnostic if: checkbox pattern found + parser didn't recognize it + invalid state
		// This means someone tried to create a task but used an invalid state value
		if !isValidTask && !s.isValidTaskState(state, cfg) {
			// Create list of valid states for the message
			validStates := s.getValidTaskStates(cfg)
			validStatesList := strings.Join(validStates, "', '")

			severity := lsp.DiagnosticSeverityWarning
			source := "notedown-task"
			message := fmt.Sprintf("Invalid task state '%s'. Valid states: '%s'", state, validStatesList)

			diagnostic := lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: line, Character: char},
					End:   lsp.Position{Line: endLine, Character: endChar},
				},
				Severity: &severity,
				Source:   &source,
				Message:  message,
				Code:     "invalid-task-state",
			}

			diagnostics = append(diagnostics, diagnostic)
		}
	}

	return diagnostics
}

// collectListItems walks the AST and collects all list items with their ranges
//
// This gives us the parser's view of which list items are actual tasks.
// The parser sets TaskList=true only for list items with valid task checkbox states.
// We use this to distinguish between:
//   - Valid tasks: [x] item    → Parser creates ListItem with TaskList=true
//   - Invalid tasks: [bad] item → Parser creates ListItem with TaskList=false
//   - Regular items: No checkbox → Parser creates ListItem with TaskList=false
func (s *Server) collectListItems(doc *parser.Document) []*parser.ListItem {
	var items []*parser.ListItem

	// Use a walk function to collect all list items from the AST
	walkFunc := parser.WalkFunc(func(node parser.Node) error {
		if listItem, ok := node.(*parser.ListItem); ok {
			items = append(items, listItem)
		}
		return nil
	})

	// Walk the entire document tree
	walker := parser.NewWalker(walkFunc)
	if err := walker.Walk(doc); err != nil {
		s.logger.Error("failed to walk document tree for list items", "error", err)
		return items
	}

	return items
}

// isPositionInRange checks if a line/character position is within a parser range
func (s *Server) isPositionInRange(line, char int, rng parser.Range) bool {
	// Convert 0-based LSP positions to 1-based parser positions
	parserLine := line + 1
	parserChar := char + 1

	// Check if position is within range
	if parserLine < rng.Start.Line || parserLine > rng.End.Line {
		return false
	}

	// If on start line, check character position
	if parserLine == rng.Start.Line && parserChar < rng.Start.Column {
		return false
	}

	// If on end line, check character position
	if parserLine == rng.End.Line && parserChar > rng.End.Column {
		return false
	}

	return true
}

// isValidTaskState checks if the given state value is valid according to configuration
func (s *Server) isValidTaskState(stateValue string, cfg *config.Config) bool {
	// Check each configured task state (including aliases)
	for _, state := range cfg.Tasks.States {
		if state.HasValue(stateValue) {
			return true
		}
	}
	return false
}

// getValidTaskStates returns a list of all valid task state values (including aliases)
func (s *Server) getValidTaskStates(cfg *config.Config) []string {
	var states []string
	for _, state := range cfg.Tasks.States {
		states = append(states, state.Value)
		states = append(states, state.Aliases...)
	}
	return states
}

// positionFromOffset converts a byte offset to line and character position
func (s *Server) positionFromOffset(content string, offset int) (int, int) {
	lines := strings.Split(content[:offset], "\n")
	line := len(lines) - 1
	char := len(lines[line])
	if line > 0 {
		char = len(lines[line])
	}
	return line, char
}
