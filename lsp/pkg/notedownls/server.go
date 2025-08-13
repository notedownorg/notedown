package notedownls

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
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
}

// NewServer creates a new Notedown LSP server
func NewServer(version string, logger *log.Logger) *Server {
	scopedLogger := logger.WithScope("lsp/pkg/notedownls")
	return &Server{
		version:   version,
		logger:    scopedLogger,
		documents: make(map[string]*Document),
		workspace: NewWorkspaceManager(scopedLogger.WithScope("workspace")),
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

	s.logger.Debug("registered document lifecycle, workspace, and completion handlers")
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

// handleCompletion provides wikilink completion suggestions
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

	// Check if we're in a wikilink context
	wikilinkInfo := s.getWikilinkContext(doc, completionParams.Position)
	if wikilinkInfo == nil {
		s.logger.Debug("not in wikilink context, returning empty completion")
		return &lsp.CompletionList{IsIncomplete: false, Items: []lsp.CompletionItem{}}, nil
	}

	s.logger.Debug("detected wikilink context", "prefix", wikilinkInfo.Prefix, "isComplete", wikilinkInfo.IsComplete)

	// Get completion items based on workspace files
	items := s.getWikilinkCompletions(wikilinkInfo.Prefix, completionParams.TextDocument.URI)

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

// getWikilinkCompletions generates completion items for wikilinks
func (s *Server) getWikilinkCompletions(prefix, currentDocURI string) []lsp.CompletionItem {
	var items []lsp.CompletionItem

	// Get all markdown files from workspace
	files := s.GetWorkspaceFiles()

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
				kind := lsp.CompletionItemKindFile
				items = append(items, lsp.CompletionItem{
					Label:      target.Link,
					Kind:       &kind,
					Detail:     &target.Detail,
					InsertText: &target.Link,
					FilterText: &target.Link,
					SortText:   &target.SortKey,
				})
			}
		}
	}

	return items
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

// Shutdown handles cleanup when the server is shutting down
func (s *Server) Shutdown() error {
	s.logger.Info("shutting down Notedown language server", "version", s.version)
	return nil
}
