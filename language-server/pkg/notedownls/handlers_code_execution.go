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
	"strings"
	"time"

	"github.com/notedownorg/notedown/language-server/pkg/codeexec"
	"github.com/notedownorg/notedown/language-server/pkg/lsp"
)

// ExecuteCodeBlocksParams represents the parameters for the executeCodeBlocks command
type ExecuteCodeBlocksParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Language     *string                `json:"language,omitempty"` // optional, defaults to "go"
}

// TextDocumentIdentifier represents a text document identifier
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// ExecuteCodeBlocksResult represents the result of code execution with workspace edits
type ExecuteCodeBlocksResult struct {
	Applied bool               `json:"applied"`
	Changes *lsp.WorkspaceEdit `json:"changes,omitempty"`
	Error   string             `json:"error,omitempty"`
}

// handleExecuteCodeBlocks handles the notedown.executeCodeBlocks command
func (s *Server) handleExecuteCodeBlocks(args []json.RawMessage) (*ExecuteCodeBlocksResult, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("executeCodeBlocks requires parameters")
	}

	// Parse the parameters
	var params ExecuteCodeBlocksParams
	if err := json.Unmarshal(args[0], &params); err != nil {
		s.logger.Error("failed to unmarshal executeCodeBlocks params", "error", err)
		return nil, fmt.Errorf("invalid parameters: %v", err)
	}

	// Get the document content
	doc, exists := s.GetDocument(params.TextDocument.URI)
	if !exists {
		return &ExecuteCodeBlocksResult{
			Applied: false,
			Error:   fmt.Sprintf("document not found: %s", params.TextDocument.URI),
		}, nil
	}

	// Determine languages to execute
	var languagesToExecute []string
	if params.Language != nil && *params.Language != "" {
		languagesToExecute = []string{*params.Language}
	} else {
		// Detect all supported languages with code blocks in document
		detectedLanguages, err := s.detectLanguagesInDocument(doc.Content)
		if err != nil {
			return &ExecuteCodeBlocksResult{
				Applied: false,
				Error:   fmt.Sprintf("failed to detect languages: %v", err),
			}, nil
		}
		languagesToExecute = detectedLanguages
	}

	if len(languagesToExecute) == 0 {
		return &ExecuteCodeBlocksResult{
			Applied: false,
			Error:   "no supported languages with code blocks found in document",
		}, nil
	}

	s.logger.Debug("executing code blocks", "uri", params.TextDocument.URI, "languages", languagesToExecute)

	// Get workspace root for execution
	workspaceRoot, err := s.getWorkspaceRoot(params.TextDocument.URI)
	if err != nil {
		s.logger.Error("failed to get workspace root", "uri", params.TextDocument.URI, "error", err)
		return &ExecuteCodeBlocksResult{
			Applied: false,
			Error:   fmt.Sprintf("failed to determine workspace root: %v", err),
		}, nil
	}

	// Generate workspace edit with execution results
	workspaceEdit, err := s.generateExecutionWorkspaceEdit(params.TextDocument.URI, doc.Content, languagesToExecute, workspaceRoot)
	if err != nil {
		s.logger.Error("failed to generate workspace edit", "error", err)
		return &ExecuteCodeBlocksResult{
			Applied: false,
			Error:   fmt.Sprintf("failed to generate workspace edit: %v", err),
		}, nil
	}

	s.logger.Debug("code execution workspace edit generated", "languages", languagesToExecute)
	return &ExecuteCodeBlocksResult{
		Applied: true,
		Changes: workspaceEdit,
	}, nil
}

// detectLanguagesInDocument scans document content to find supported languages with code blocks
func (s *Server) detectLanguagesInDocument(content string) ([]string, error) {
	factory := codeexec.NewExecutorFactory()
	supportedLanguages := factory.GetSupportedLanguages()

	var foundLanguages []string
	languageMap := make(map[string]bool)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "```") && len(trimmedLine) > 3 {
			language := strings.ToLower(strings.TrimSpace(trimmedLine[3:]))
			if language != "" && !languageMap[language] {
				// Check if this is a supported language
				for _, supportedLang := range supportedLanguages {
					if language == supportedLang {
						foundLanguages = append(foundLanguages, language)
						languageMap[language] = true
						break
					}
				}
			}
		}
	}

	return foundLanguages, nil
}

// generateExecutionWorkspaceEdit executes code and generates workspace edit with results
func (s *Server) generateExecutionWorkspaceEdit(uri, content string, languages []string, workspaceRoot string) (*lsp.WorkspaceEdit, error) {
	// Parse existing output blocks to find what needs to be removed
	blocksToRemove := s.findExistingOutputBlocks(content, languages)

	// Execute code for each language and collect results
	var newOutputBlocks []OutputBlock

	factory := codeexec.NewExecutorFactory()
	for _, language := range languages {
		// Collect code blocks for this language
		codeBlocks, err := s.collectCodeBlocks(uri, language)
		if err != nil {
			s.logger.Error("failed to collect code blocks", "language", language, "error", err)
			continue
		}

		if len(codeBlocks) == 0 {
			continue // Skip languages with no code blocks
		}

		// Get executor and execute
		executor, err := factory.GetExecutor(language)
		if err != nil {
			s.logger.Error("failed to get executor", "language", language, "error", err)
			continue
		}

		result, err := executor.Execute(codeBlocks, workspaceRoot)
		if err != nil {
			s.logger.Error("code execution failed", "language", language, "error", err)
			continue
		}

		// Generate output blocks for this execution
		timestamp := time.Now().Format("2006-01-02T15:04:05")
		duration := s.formatDuration(result.ExecutionTime)

		if result.Stdout != "" {
			newOutputBlocks = append(newOutputBlocks, OutputBlock{
				Language: language,
				Type:     "stdout",
				Content:  result.Stdout,
				Metadata: timestamp + " " + duration,
			})
		}

		if result.Stderr != "" {
			newOutputBlocks = append(newOutputBlocks, OutputBlock{
				Language: language,
				Type:     "stderr",
				Content:  result.Stderr,
				Metadata: timestamp + " " + duration,
			})
		}
	}

	// Generate text edits
	textEdits := s.generateTextEdits(content, blocksToRemove, newOutputBlocks)

	// Create workspace edit
	workspaceEdit := &lsp.WorkspaceEdit{
		Changes: map[string][]lsp.TextEdit{
			uri: textEdits,
		},
	}

	return workspaceEdit, nil
}

// OutputBlock represents an output block to be inserted
type OutputBlock struct {
	Language string
	Type     string // "stdout" or "stderr"
	Content  string
	Metadata string // timestamp and duration
}

// BlockRange represents a range of lines to be removed
type BlockRange struct {
	StartLine int
	EndLine   int
}

// findExistingOutputBlocks finds existing output blocks that need to be removed
func (s *Server) findExistingOutputBlocks(content string, languages []string) []BlockRange {
	var blocks []BlockRange
	lines := strings.Split(content, "\n")

	languageMap := make(map[string]bool)
	for _, lang := range languages {
		languageMap[lang] = true
	}

	inOutputBlock := false
	var blockStartLine int

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check for output block start
		if strings.HasPrefix(trimmedLine, "```output:") && !inOutputBlock {
			// Parse the language from the output block header
			parts := strings.Split(trimmedLine[10:], ":")
			if len(parts) >= 2 {
				language := parts[0]
				if languageMap[language] {
					inOutputBlock = true
					blockStartLine = i
				}
			}
		} else if strings.HasPrefix(trimmedLine, "```") && inOutputBlock {
			// End of output block
			blocks = append(blocks, BlockRange{
				StartLine: blockStartLine,
				EndLine:   i,
			})
			inOutputBlock = false
		}
	}

	return blocks
}

// generateTextEdits creates the text edits to remove old blocks and add new ones
func (s *Server) generateTextEdits(content string, blocksToRemove []BlockRange, newBlocks []OutputBlock) []lsp.TextEdit {
	lines := strings.Split(content, "\n")
	var edits []lsp.TextEdit

	// Sort blocks to remove in reverse order to maintain line numbers
	for i := len(blocksToRemove) - 1; i >= 0; i-- {
		block := blocksToRemove[i]
		edit := lsp.TextEdit{
			Range: lsp.Range{
				Start: lsp.Position{Line: block.StartLine, Character: 0},
				End:   lsp.Position{Line: block.EndLine + 1, Character: 0}, // Include the line after
			},
			NewText: "", // Delete the block
		}
		edits = append(edits, edit)
	}

	// Add new output blocks at the end of the document
	if len(newBlocks) > 0 {
		var newContent strings.Builder
		newContent.WriteString("\n") // Add spacing before output blocks

		for _, block := range newBlocks {
			newContent.WriteString(fmt.Sprintf("```output:%s:%s %s\n", block.Language, block.Type, block.Metadata))
			newContent.WriteString(block.Content)
			if !strings.HasSuffix(block.Content, "\n") {
				newContent.WriteString("\n")
			}
			newContent.WriteString("```\n\n")
		}

		// Insert at end of document
		endLine := len(lines)
		edit := lsp.TextEdit{
			Range: lsp.Range{
				Start: lsp.Position{Line: endLine, Character: 0},
				End:   lsp.Position{Line: endLine, Character: 0},
			},
			NewText: newContent.String(),
		}
		edits = append(edits, edit)
	}

	return edits
}

// formatDuration formats execution duration in human-readable format
func (s *Server) formatDuration(duration time.Duration) string {
	if duration < time.Second {
		return fmt.Sprintf("%dms", duration.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", duration.Seconds())
}

// getWorkspaceRoot determines the workspace root for the given document URI
func (s *Server) getWorkspaceRoot(uri string) (string, error) {
	// Try to get the workspace root from the workspace manager
	roots := s.workspace.GetWorkspaceRoots()
	if len(roots) == 0 {
		return "", fmt.Errorf("no workspace roots available")
	}

	// For now, use the first workspace root
	// In the future, we could be smarter about selecting the correct root
	// based on the document URI
	return roots[0].Path, nil
}
