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
	"path/filepath"
	"strings"
	"sync"

	"github.com/notedownorg/notedown/language-server/pkg/lsp"
	"github.com/notedownorg/notedown/language-server/pkg/notedownls/indexes"
	"github.com/notedownorg/notedown/pkg/config"
	"github.com/notedownorg/notedown/pkg/log"
)

// Server implements the LSPServer interface for Notedown
type Server struct {
	version string
	logger  *log.Logger

	// Document storage
	documents      map[string]*Document
	documentsMutex sync.RWMutex

	// Workspace management
	workspace *WorkspaceManager

	// Wikilink management
	wikilinkIndex *indexes.WikilinkIndex

	// Diagnostic publishing
	diagnosticPublisher func(params lsp.PublishDiagnosticsParams) error

	// Client request sender (for workspace/applyEdit)
	clientRequestSender func(method string, params any) (any, error)
}

// NewServer creates a new Notedown LSP server
func NewServer(version string, logger *log.Logger) *Server {
	scopedLogger := logger.WithScope("lsp/pkg/notedownls")
	return &Server{
		version:       version,
		logger:        scopedLogger,
		documents:     make(map[string]*Document),
		workspace:     NewWorkspaceManager(scopedLogger.WithScope("workspace")),
		wikilinkIndex: indexes.NewWikilinkIndex(scopedLogger.WithScope("wikilink")),
	}
}

// Initialize handles the LSP initialize request
func (s *Server) Initialize(params lsp.InitializeParams) (lsp.InitializeResult, error) {
	clientName := "unknown"
	if params.ClientInfo != nil {
		clientName = params.ClientInfo.Name
	}

	s.logger.Info("lsp client initialized", "client", clientName, "server_version", s.version)

	// Initialize workspace from parameters
	if err := s.workspace.InitializeFromParams(params); err != nil {
		s.logger.Error("failed to initialize workspace", "error", err)
		return lsp.InitializeResult{}, err
	}

	// Start workspace file discovery in background
	go func() {
		if err := s.workspace.DiscoverMarkdownFiles(); err != nil {
			s.logger.Error("workspace file discovery failed", "error", err)
		}
	}()

	syncKind := lsp.TextDocumentSyncKindFull
	result := lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "Notedown Language Server", Version: s.version},
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: &[]bool{true}[0],
				Change:    &syncKind,
			},
			CompletionProvider: &lsp.CompletionOptions{
				TriggerCharacters: []string{"["},
				ResolveProvider:   &[]bool{false}[0],
			},
			DefinitionProvider:   &[]bool{true}[0],
			FoldingRangeProvider: &[]bool{true}[0],
			CodeActionProvider:   &lsp.CodeActionOptions{},
			ExecuteCommandProvider: &lsp.ExecuteCommandOptions{
				Commands: []string{
					"notedown.getListItemBoundaries",
					"notedown.getConcealRanges",
				},
			},
			// DiagnosticProvider: &lsp.DiagnosticOptions{}, // Disabled to avoid pull diagnostic conflicts
			Workspace: &lsp.WorkspaceServerCapabilities{
				WorkspaceFolders: &lsp.WorkspaceFoldersServerCapabilities{
					Supported:           &[]bool{true}[0],
					ChangeNotifications: true,
				},
			},
		},
	}
	return result, nil
}

// RegisterHandlers registers all method and notification handlers
func (s *Server) RegisterHandlers(mux *lsp.Mux) error {
	// Register document lifecycle handlers
	mux.RegisterNotification(lsp.MethodTextDocumentDidOpen, s.handleDidOpen)
	mux.RegisterNotification(lsp.MethodTextDocumentDidChange, s.handleDidChange)
	mux.RegisterNotification(lsp.MethodTextDocumentDidClose, s.handleDidClose)

	// Register workspace handlers
	mux.RegisterNotification(lsp.MethodWorkspaceDidChangeWatchedFiles, s.handleDidChangeWatchedFiles)

	// Register completion handler
	mux.RegisterMethod(lsp.MethodTextDocumentCompletion, s.handleCompletion)

	// Register definition handler
	mux.RegisterMethod(lsp.MethodTextDocumentDefinition, s.handleDefinition)

	// Register folding range handler
	mux.RegisterMethod(lsp.MethodTextDocumentFoldingRange, s.handleFoldingRange)

	// Register code action handler
	mux.RegisterMethod(lsp.MethodTextDocumentCodeAction, s.handleCodeAction)

	// Register workspace execute command handler
	mux.RegisterMethod(lsp.MethodWorkspaceExecuteCommand, s.handleExecuteCommand)

	// Register diagnostic handler (pull diagnostics) - disabled to avoid conflicts
	// mux.RegisterMethod(lsp.MethodTextDocumentDiagnostic, s.handleDiagnostic)

	// Set up diagnostic publishing
	s.SetDiagnosticPublisher(func(params lsp.PublishDiagnosticsParams) error {
		return mux.PublishNotification(string(lsp.MethodTextDocumentPublishDiagnostics), params)
	})

	// Set up client request sending (for workspace/applyEdit)
	s.SetClientRequestSender(func(method string, params any) (any, error) {
		return mux.SendRequest(method, params)
	})

	s.logger.Debug("registered document lifecycle, workspace, completion, definition, folding range, code action, execute command handlers, and diagnostic publishing")
	return nil
}

// GetDocument retrieves a document by URI
func (s *Server) GetDocument(uri string) (*Document, bool) {
	s.documentsMutex.RLock()
	defer s.documentsMutex.RUnlock()
	doc, exists := s.documents[uri]
	return doc, exists
}

// AddDocument adds or updates a document in storage
func (s *Server) AddDocument(uri string) (*Document, error) {
	doc, err := NewDocument(uri)
	if err != nil {
		return nil, err
	}

	s.documentsMutex.Lock()
	defer s.documentsMutex.Unlock()
	s.documents[uri] = doc

	s.logger.Debug("document added", "uri", uri, "basepath", doc.Basepath)
	return doc, nil
}

// RemoveDocument removes a document from storage
func (s *Server) RemoveDocument(uri string) {
	s.documentsMutex.Lock()
	defer s.documentsMutex.Unlock()

	if doc, exists := s.documents[uri]; exists {
		delete(s.documents, uri)
		s.logger.Debug("document removed", "uri", uri, "basepath", doc.Basepath)
	}
}

// HasDocument checks if a document exists in storage
func (s *Server) HasDocument(uri string) bool {
	s.documentsMutex.RLock()
	defer s.documentsMutex.RUnlock()
	_, exists := s.documents[uri]
	return exists
}

// GetWorkspace returns the workspace manager
func (s *Server) GetWorkspace() *WorkspaceManager {
	return s.workspace
}

// GetWorkspaceFiles returns all indexed Markdown files
func (s *Server) GetWorkspaceFiles() []*FileInfo {
	return s.workspace.GetMarkdownFiles()
}

// GetWorkspaceRoots returns the workspace roots
func (s *Server) GetWorkspaceRoots() []WorkspaceRoot {
	return s.workspace.GetWorkspaceRoots()
}

// handleCompletion provides wikilink and task state completion suggestions
func (s *Server) handleCompletion(params json.RawMessage) (any, error) {
	var completionParams lsp.CompletionParams
	if err := json.Unmarshal(params, &completionParams); err != nil {
		s.logger.Error("failed to unmarshal completion params", "error", err)
		return nil, err
	}

	s.logger.Debug("completion request received",
		"uri", completionParams.TextDocument.URI,
		"line", completionParams.Position.Line,
		"character", completionParams.Position.Character)

	// Get the document
	doc, exists := s.GetDocument(completionParams.TextDocument.URI)
	if !exists {
		s.logger.Debug("document not found for completion", "uri", completionParams.TextDocument.URI)
		return &lsp.CompletionList{IsIncomplete: false, Items: []lsp.CompletionItem{}}, nil
	}

	// Check if we're in a task state context first
	taskContext := s.getTaskContext(doc, completionParams.Position)
	if taskContext != nil {
		s.logger.Debug("detected task context", "prefix", taskContext.Prefix, "isComplete", taskContext.IsComplete)

		// Load workspace configuration for task states
		cfg, err := s.loadWorkspaceConfig()
		if err != nil {
			s.logger.Error("failed to load workspace config for task completion", "error", err)
			return &lsp.CompletionList{IsIncomplete: false, Items: []lsp.CompletionItem{}}, nil
		}

		// Get task state completion items
		items := s.getTaskStateCompletions(taskContext.Prefix, cfg, !taskContext.IsComplete)
		s.logger.Debug("generated task state completion items", "count", len(items))
		return &lsp.CompletionList{
			IsIncomplete: false,
			Items:        items,
		}, nil
	}

	// Check if we're in a wikilink context
	wikilinkInfo := s.getWikilinkContext(doc, completionParams.Position)
	if wikilinkInfo == nil {
		s.logger.Debug("not in wikilink or task context, returning empty completion")
		return &lsp.CompletionList{IsIncomplete: false, Items: []lsp.CompletionItem{}}, nil
	}

	s.logger.Debug("detected wikilink context", "prefix", wikilinkInfo.Prefix, "isComplete", wikilinkInfo.IsComplete)

	// Get completion items based on workspace files
	items := s.getWikilinkCompletions(wikilinkInfo.Prefix, completionParams.TextDocument.URI, !wikilinkInfo.IsComplete)

	s.logger.Debug("generated completion items", "count", len(items))
	return &lsp.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

// WikilinkContext contains information about the current wikilink being edited
type WikilinkContext struct {
	Prefix     string    // The partial wikilink text before cursor
	IsComplete bool      // Whether the wikilink is complete (has closing ]])
	Range      lsp.Range // The range of the wikilink prefix
}

// TaskContext contains information about the current task state being edited
type TaskContext struct {
	Prefix     string    // The partial task state text before cursor
	IsComplete bool      // Whether the task state is complete (has closing ])
	Range      lsp.Range // The range of the task state prefix
}

// getWikilinkContext analyzes the cursor position to determine if we're in a wikilink
func (s *Server) getWikilinkContext(doc *Document, position lsp.Position) *WikilinkContext {
	lines := strings.Split(doc.Content, "\n")
	if position.Line >= len(lines) {
		return nil
	}

	line := lines[position.Line]
	if position.Character > len(line) {
		return nil
	}

	// Look for [[ before the cursor position
	beforeCursor := line[:position.Character]

	// Find the last occurrence of [[ before cursor
	lastWikilinkStart := strings.LastIndex(beforeCursor, "[[")
	if lastWikilinkStart == -1 {
		return nil
	}

	// Check if there's a ]] between [[ and cursor
	afterWikilinkStart := beforeCursor[lastWikilinkStart+2:]
	if strings.Contains(afterWikilinkStart, "]]") {
		return nil // We're not in an open wikilink
	}

	// Extract the prefix (text between [[ and cursor)
	// For wikilinks with pipe separators ([[target|display]]), we only want the target part
	prefix := afterWikilinkStart
	if pipeIndex := strings.Index(prefix, "|"); pipeIndex != -1 {
		prefix = prefix[:pipeIndex]
	}

	// Check if the wikilink is complete (look for ]] after cursor)
	afterCursor := line[position.Character:]
	isComplete := strings.HasPrefix(afterCursor, "]]") || strings.Contains(afterCursor, "]]")

	// Calculate the end position based on the prefix length
	endCharacter := lastWikilinkStart + 2 + len(prefix)

	return &WikilinkContext{
		Prefix:     prefix,
		IsComplete: isComplete,
		Range: lsp.Range{
			Start: lsp.Position{Line: position.Line, Character: lastWikilinkStart + 2},
			End:   lsp.Position{Line: position.Line, Character: endCharacter},
		},
	}
}

// getTaskContext analyzes the cursor position to determine if we're in a task state
func (s *Server) getTaskContext(doc *Document, position lsp.Position) *TaskContext {
	lines := strings.Split(doc.Content, "\n")
	if position.Line >= len(lines) {
		return nil
	}

	line := lines[position.Line]
	if position.Character > len(line) {
		return nil
	}

	// Look for task list pattern: - [ at start of line (with optional whitespace)
	// This matches patterns like:
	// - [ ]
	// - [x]
	// - [wip]
	// etc.

	// First, check if this line looks like a task list item
	trimmedLine := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(trimmedLine, "- [") {
		return nil
	}

	// Find the position of the opening [ for the task state
	taskStart := strings.Index(line, "- [")
	if taskStart == -1 {
		return nil
	}

	// The task state starts right after "- ["
	stateStart := taskStart + 3

	// Check if cursor is within the task state bracket
	if position.Character < stateStart {
		return nil // cursor is before the task state
	}

	// Look for the closing ] to determine if the task state is complete
	afterStateStart := line[stateStart:]
	closeIndex := strings.Index(afterStateStart, "]")

	// If there's no closing bracket, or cursor is after it, we're not in task context
	if closeIndex == -1 {
		// No closing bracket yet - we might be typing the state
		if position.Character <= len(line) {
			prefix := line[stateStart:position.Character]
			return &TaskContext{
				Prefix:     prefix,
				IsComplete: false,
				Range: lsp.Range{
					Start: lsp.Position{Line: position.Line, Character: stateStart},
					End:   lsp.Position{Line: position.Line, Character: position.Character},
				},
			}
		}
		return nil
	}

	stateEnd := stateStart + closeIndex

	// Check if cursor is within the task state brackets
	if position.Character > stateEnd {
		return nil // cursor is after the closing bracket
	}

	// Extract the current prefix up to the cursor position
	prefixEnd := position.Character
	if prefixEnd > stateEnd {
		prefixEnd = stateEnd
	}

	prefix := line[stateStart:prefixEnd]

	return &TaskContext{
		Prefix:     prefix,
		IsComplete: true,
		Range: lsp.Range{
			Start: lsp.Position{Line: position.Line, Character: stateStart},
			End:   lsp.Position{Line: position.Line, Character: stateEnd},
		},
	}
}

// getWikilinkCompletions generates completion items for wikilinks
func (s *Server) getWikilinkCompletions(prefix, currentDocURI string, needsClosing bool) []lsp.CompletionItem {
	var items []lsp.CompletionItem

	// Priority 1: Existing files
	items = append(items, s.getExistingFileCompletions(prefix, currentDocURI, needsClosing)...)

	// Priority 2: Non-existent targets referenced in other documents
	items = append(items, s.getNonExistentTargetCompletions(prefix, currentDocURI, needsClosing)...)

	// Priority 3: Directory path suggestions
	items = append(items, s.getDirectoryPathCompletions(prefix, currentDocURI, needsClosing)...)

	return items
}

// getExistingFileCompletions generates completions for existing workspace files
func (s *Server) getExistingFileCompletions(prefix, currentDocURI string, needsClosing bool) []lsp.CompletionItem {
	var items []lsp.CompletionItem

	// Get all markdown files from workspace
	files := s.GetWorkspaceFiles()

	// Track targets to detect conflicts
	targetToFiles := make(map[string][]*FileInfo)

	// First pass: collect all targets and their matching files
	for _, fileInfo := range files {
		// Skip the current document
		if fileInfo.URI == currentDocURI {
			continue
		}

		// Generate possible wikilink targets for this file
		targets := s.generateWikilinkTargets(fileInfo, currentDocURI)

		for _, target := range targets {
			// Filter based on prefix
			if prefix == "" || strings.HasPrefix(strings.ToLower(target.Link), strings.ToLower(prefix)) {
				targetToFiles[target.Link] = append(targetToFiles[target.Link], fileInfo)
			}
		}
	}

	// Second pass: generate completion items with conflict information
	addedTargets := make(map[string]bool)

	for target, matchingFiles := range targetToFiles {
		if addedTargets[target] {
			continue
		}
		addedTargets[target] = true

		insertText := target
		if needsClosing {
			insertText += "]]"
		}

		kind := lsp.CompletionItemKindFile

		if len(matchingFiles) == 1 {
			// Single match - use existing logic
			fileInfo := matchingFiles[0]
			detail := fmt.Sprintf("Link to %s", fileInfo.Path)
			sortKey := fmt.Sprintf("0_%s", target)

			items = append(items, lsp.CompletionItem{
				Label:      target,
				Kind:       &kind,
				Detail:     &detail,
				InsertText: &insertText,
				FilterText: &target,
				SortText:   &sortKey,
			})
		} else {
			// Multiple matches - show ambiguous target with warning
			var filePaths []string
			for _, fileInfo := range matchingFiles {
				filePaths = append(filePaths, fileInfo.Path)
			}

			detail := fmt.Sprintf("⚠️ Ambiguous: %s", strings.Join(filePaths, ", "))
			sortKey := fmt.Sprintf("0_%s_ambiguous", target)

			items = append(items, lsp.CompletionItem{
				Label:      target + " (ambiguous)",
				Kind:       &kind,
				Detail:     &detail,
				InsertText: &insertText,
				FilterText: &target,
				SortText:   &sortKey,
			})

			// Also add specific path-based completions for each match
			for i, fileInfo := range matchingFiles {
				pathWithoutExt := strings.TrimSuffix(fileInfo.Path, ".md")
				pathWithoutExt = strings.ReplaceAll(pathWithoutExt, "\\", "/")

				pathInsertText := pathWithoutExt
				if needsClosing {
					pathInsertText += "]]"
				}

				pathDetail := fmt.Sprintf("Link to %s (disambiguated)", fileInfo.Path)
				pathSortKey := fmt.Sprintf("0_%s_path_%d", target, i)

				items = append(items, lsp.CompletionItem{
					Label:      pathWithoutExt,
					Kind:       &kind,
					Detail:     &pathDetail,
					InsertText: &pathInsertText,
					FilterText: &pathWithoutExt,
					SortText:   &pathSortKey,
				})
			}
		}
	}

	return items
}

// getNonExistentTargetCompletions generates completions for non-existent wikilink targets
func (s *Server) getNonExistentTargetCompletions(prefix, currentDocURI string, needsClosing bool) []lsp.CompletionItem {
	var items []lsp.CompletionItem

	// Get non-existent targets from the wikilink index
	nonExistentTargets := s.wikilinkIndex.GetNonExistentTargets()

	for _, targetInfo := range nonExistentTargets {
		target := targetInfo.Target

		// Filter based on prefix
		if prefix == "" || strings.HasPrefix(strings.ToLower(target), strings.ToLower(prefix)) {
			refCount := len(targetInfo.ReferencedBy)

			insertText := target
			if needsClosing {
				insertText += "]]"
			}

			kind := lsp.CompletionItemKindReference
			detail := fmt.Sprintf("Referenced in %d file(s) (create new)", refCount)
			sortKey := fmt.Sprintf("2_%s", target) // Lower priority than existing files

			items = append(items, lsp.CompletionItem{
				Label:      target,
				Kind:       &kind,
				Detail:     &detail,
				InsertText: &insertText,
				FilterText: &target,
				SortText:   &sortKey,
			})
		}
	}

	return items
}

// getDirectoryPathCompletions generates completions based on existing directory structure
func (s *Server) getDirectoryPathCompletions(prefix, currentDocURI string, needsClosing bool) []lsp.CompletionItem {
	var items []lsp.CompletionItem

	// Get all workspace files to analyze directory structure
	files := s.GetWorkspaceFiles()
	directorySet := make(map[string]bool)

	// Extract all directory paths from existing files
	for _, fileInfo := range files {
		dir := filepath.Dir(fileInfo.Path)
		if dir != "." {
			// Add the directory and all parent directories
			parts := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
			currentPath := ""
			for _, part := range parts {
				if currentPath == "" {
					currentPath = part
				} else {
					currentPath += "/" + part
				}
				directorySet[currentPath] = true
			}
		}
	}

	// Generate completions for directories that match the prefix
	for dirPath := range directorySet {
		// Complete directory paths (e.g., "docs/", "projects/")
		dirCompletion := dirPath + "/"

		if prefix == "" || strings.HasPrefix(strings.ToLower(dirCompletion), strings.ToLower(prefix)) {
			kind := lsp.CompletionItemKindFolder
			detail := "Directory path completion"
			sortKey := fmt.Sprintf("3_%s", dirCompletion) // Lower priority than files and referenced targets

			// Don't include closing ]] for directories - user might want to continue typing
			items = append(items, lsp.CompletionItem{
				Label:      dirCompletion,
				Kind:       &kind,
				Detail:     &detail,
				InsertText: &dirCompletion,
				FilterText: &dirCompletion,
				SortText:   &sortKey,
			})
		}

		// If prefix matches the directory, also suggest common file patterns within it
		if prefix != "" && strings.HasPrefix(strings.ToLower(dirPath), strings.ToLower(prefix)) {
			// Suggest a generic file in this directory
			suggestedFile := dirPath + "/new-file"

			if !s.targetAlreadyExists(suggestedFile, files) {
				insertText := suggestedFile
				if needsClosing {
					insertText += "]]"
				}

				kind := lsp.CompletionItemKindValue
				detail := fmt.Sprintf("Create new file in %s/", dirPath)
				sortKey := fmt.Sprintf("4_%s", suggestedFile)

				items = append(items, lsp.CompletionItem{
					Label:      suggestedFile,
					Kind:       &kind,
					Detail:     &detail,
					InsertText: &insertText,
					FilterText: &suggestedFile,
					SortText:   &sortKey,
				})
			}
		}
	}

	return items
}

// getTaskStateCompletions generates completion items for task states
func (s *Server) getTaskStateCompletions(prefix string, cfg *config.Config, needsClosing bool) []lsp.CompletionItem {
	var items []lsp.CompletionItem

	// Generate completions for each configured task state
	for i, state := range cfg.Tasks.States {
		// Collect all possible values (main value + aliases)
		allValues := []string{state.Value}
		allValues = append(allValues, state.Aliases...)

		for j, value := range allValues {
			// Filter based on prefix
			if prefix == "" || strings.HasPrefix(strings.ToLower(value), strings.ToLower(prefix)) {
				insertText := value
				if needsClosing {
					insertText += "]"
				}

				kind := lsp.CompletionItemKindEnum

				// Build detail: name + conceal (if set)
				detail := state.Name
				if state.Conceal != nil && *state.Conceal != "" {
					detail += fmt.Sprintf(" %s", *state.Conceal)
				}

				// Sort by config order, then group main value with its aliases
				// Format: configOrder_aliasIndex_value
				// This groups: value1, value1alias1, value1alias2, value2, value2alias1, etc.
				sortKey := fmt.Sprintf("%02d_%02d_%s", i, j, value)

				item := lsp.CompletionItem{
					Label:      value,
					Kind:       &kind,
					Detail:     &detail,
					InsertText: &insertText,
					FilterText: &value,
					SortText:   &sortKey,
				}

				// Build concise documentation (Option 2 format)
				var docParts []string

				if j > 0 {
					// For aliases: "Alias for 'x' (main value). See also: X, completed"
					docParts = append(docParts, fmt.Sprintf("Alias for '%s' (main value)", state.Value))
					if len(state.Aliases) > 1 {
						// Show other aliases (excluding current one)
						var otherAliases []string
						for _, alias := range state.Aliases {
							if alias != value {
								otherAliases = append(otherAliases, alias)
							}
						}
						if len(otherAliases) > 0 {
							docParts = append(docParts, fmt.Sprintf("See also: %s", strings.Join(otherAliases, ", ")))
						}
					}
				} else {
					// For main values: description + see also (aliases)
					if state.Description != nil && *state.Description != "" {
						docParts = append(docParts, *state.Description)
					}
					if len(state.Aliases) > 0 {
						docParts = append(docParts, fmt.Sprintf("See also: %s", strings.Join(state.Aliases, ", ")))
					}
				}

				if len(docParts) > 0 {
					item.Documentation = strings.Join(docParts, "\n\n")
				}

				items = append(items, item)
			}
		}
	}

	return items
}

// targetAlreadyExists checks if a target would conflict with existing files or targets
func (s *Server) targetAlreadyExists(target string, files []*FileInfo) bool {
	// Check existing files
	for _, fileInfo := range files {
		pathWithoutExt := strings.TrimSuffix(fileInfo.Path, filepath.Ext(fileInfo.Path))
		baseWithoutExt := strings.TrimSuffix(filepath.Base(fileInfo.Path), filepath.Ext(fileInfo.Path))

		if target == pathWithoutExt || target == baseWithoutExt {
			return true
		}
	}

	// Check existing wikilink targets
	allTargets := s.wikilinkIndex.GetAllTargets()
	_, exists := allTargets[target]

	return exists
}

// WikilinkTarget represents a possible wikilink target
type WikilinkTarget struct {
	Link    string // The wikilink text (without [[ ]])
	Detail  string // Additional information to show
	SortKey string // Key for sorting (prioritizes shorter paths)
}

// generateWikilinkTargets generates possible wikilink targets for a file
func (s *Server) generateWikilinkTargets(fileInfo *FileInfo, currentDocURI string) []WikilinkTarget {
	var targets []WikilinkTarget

	// Remove .md extension for the base name
	baseName := strings.TrimSuffix(filepath.Base(fileInfo.Path), ".md")

	// Add base name as primary target
	targets = append(targets, WikilinkTarget{
		Link:    baseName,
		Detail:  fmt.Sprintf("Link to %s", fileInfo.Path),
		SortKey: fmt.Sprintf("0_%s", baseName), // High priority
	})

	// If the file is in a subdirectory, also add the relative path version
	if dir := filepath.Dir(fileInfo.Path); dir != "." {
		pathWithoutExt := strings.TrimSuffix(fileInfo.Path, ".md")

		// Convert backslashes to forward slashes for consistency
		pathWithoutExt = strings.ReplaceAll(pathWithoutExt, "\\", "/")

		targets = append(targets, WikilinkTarget{
			Link:    pathWithoutExt,
			Detail:  fmt.Sprintf("Link to %s", fileInfo.Path),
			SortKey: fmt.Sprintf("1_%s", pathWithoutExt), // Lower priority
		})
	}

	return targets
}

// extractWikilinksFromDocument extracts wikilinks from a document and adds them to the index
func (s *Server) extractWikilinksFromDocument(documentURI, content string) {
	// Get workspace files as a map for the interface
	workspaceFiles := s.getWorkspaceFilesMap()

	// Extract wikilinks
	s.wikilinkIndex.ExtractWikilinksFromDocument(content, documentURI, workspaceFiles)
}

// refreshWikilinksFromDocument refreshes wikilinks for a document (removes old, adds new)
func (s *Server) refreshWikilinksFromDocument(documentURI, content string) {
	// Get workspace files as a map for the interface
	workspaceFiles := s.getWorkspaceFilesMap()

	// Refresh wikilinks
	s.wikilinkIndex.RefreshDocumentWikilinks(content, documentURI, workspaceFiles)
}

// removeWikilinksFromDocument removes all wikilink references from a document
func (s *Server) removeWikilinksFromDocument(documentURI string) {
	// Remove references by refreshing with empty content
	s.wikilinkIndex.RefreshDocumentWikilinks("", documentURI, nil)
}

// getWorkspaceFilesMap converts workspace files to the interface map
func (s *Server) getWorkspaceFilesMap() map[string]indexes.WorkspaceFile {
	workspaceFiles := s.GetWorkspaceFiles()
	result := make(map[string]indexes.WorkspaceFile)

	for _, fileInfo := range workspaceFiles {
		result[fileInfo.URI] = fileInfo
	}

	return result
}

// SetDiagnosticPublisher sets the diagnostic publishing function
func (s *Server) SetDiagnosticPublisher(publisher func(params lsp.PublishDiagnosticsParams) error) {
	s.diagnosticPublisher = publisher
}

// SetClientRequestSender sets the client request sending function
func (s *Server) SetClientRequestSender(sender func(method string, params any) (any, error)) {
	s.clientRequestSender = sender
}

// publishDiagnostics publishes diagnostics for a document
func (s *Server) publishDiagnostics(uri string, diagnostics []lsp.Diagnostic) {
	if s.diagnosticPublisher == nil {
		return
	}

	// Ensure diagnostics is never nil to avoid JSON serialization issues
	if diagnostics == nil {
		diagnostics = []lsp.Diagnostic{}
	}

	params := lsp.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	s.logger.Debug("preparing to publish diagnostics", "uri", uri, "count", len(diagnostics), "diagnostics_is_nil", diagnostics == nil)

	if err := s.diagnosticPublisher(params); err != nil {
		s.logger.Error("failed to publish diagnostics", "uri", uri, "error", err)
	} else {
		s.logger.Debug("published diagnostics", "uri", uri, "count", len(diagnostics))
	}
}

// refreshAllDocumentDiagnostics regenerates and publishes diagnostics for all open documents
func (s *Server) refreshAllDocumentDiagnostics() {
	s.documentsMutex.RLock()
	defer s.documentsMutex.RUnlock()

	// Get workspace files as a map for the indexing
	workspaceFiles := s.getWorkspaceFilesMap()

	for uri, doc := range s.documents {
		// Refresh wikilinks for this document to detect new conflicts
		s.wikilinkIndex.RefreshDocumentWikilinks(doc.Content, uri, workspaceFiles)

		// Generate and publish updated diagnostics
		wikilinkDiagnostics := s.generateWikilinkDiagnostics(uri, doc.Content)
		taskDiagnostics := s.generateTaskDiagnostics(uri, doc.Content)

		// Combine all diagnostics
		allDiagnostics := append(wikilinkDiagnostics, taskDiagnostics...)
		s.publishDiagnostics(uri, allDiagnostics)
	}

	s.logger.Debug("refreshed diagnostics for all open documents", "count", len(s.documents))
}

// loadWorkspaceConfig loads the configuration for the current workspace
func (s *Server) loadWorkspaceConfig() (*config.Config, error) {
	// Get workspace roots
	workspaceRoots := s.GetWorkspaceRoots()
	if len(workspaceRoots) == 0 {
		// No workspace, use default config
		return config.GetDefaultConfig(), nil
	}

	// Try to load config from the first workspace root
	cfg, err := config.LoadConfig(workspaceRoots[0].Path)
	if err != nil {
		s.logger.Debug("failed to load workspace config, using default", "error", err, "path", workspaceRoots[0].Path)
		return config.GetDefaultConfig(), nil
	}

	return cfg, nil
}

// Shutdown handles cleanup when the server is shutting down
func (s *Server) Shutdown() error {
	s.logger.Info("shutting down Notedown language server", "version", s.version)
	return nil
}
