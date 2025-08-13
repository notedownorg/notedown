package notedownls

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/notedownorg/notedown/lsp/pkg/lsp"
	"github.com/notedownorg/notedown/pkg/log"
)

// WorkspaceManager manages workspace-wide file discovery and indexing
type WorkspaceManager struct {
	roots      []WorkspaceRoot
	fileIndex  map[string]*FileInfo // URI -> lightweight file info
	loadedDocs map[string]*Document // URI -> full document (for opened files)
	mutex      sync.RWMutex
	logger     *log.Logger

	// Configuration
	maxFileCount    int
	excludePatterns []string
}

// WorkspaceRoot represents a workspace root directory
type WorkspaceRoot struct {
	URI  string // file:// URI
	Path string // local filesystem path
	Name string // display name
}

// FileInfo contains lightweight metadata about a Markdown file
type FileInfo struct {
	URI     string    // file:// URI
	Path    string    // relative path from workspace root
	ModTime time.Time // last modification time
	Size    int64     // file size in bytes
}

// NewWorkspaceManager creates a new workspace manager
func NewWorkspaceManager(logger *log.Logger) *WorkspaceManager {
	return &WorkspaceManager{
		roots:      make([]WorkspaceRoot, 0),
		fileIndex:  make(map[string]*FileInfo),
		loadedDocs: make(map[string]*Document),
		logger:     logger,

		// Default configuration
		maxFileCount: 10000,
		excludePatterns: []string{
			".git",
			".vscode",
			".idea",
			"node_modules",
			".next",
			".nuxt",
			"dist",
			"build",
			"target",
			"__pycache__",
			".pytest_cache",
			"coverage",
			".coverage",
			"venv",
			".venv",
			"env",
			".env",
		},
	}
}

// InitializeFromParams initializes workspace roots from LSP InitializeParams
func (wm *WorkspaceManager) InitializeFromParams(params lsp.InitializeParams) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	// Clear existing roots
	wm.roots = make([]WorkspaceRoot, 0)

	// Priority 1: WorkspaceFolders (modern approach)
	if len(params.WorkspaceFolders) > 0 {
		for _, folder := range params.WorkspaceFolders {
			root, err := wm.workspaceFolderToRoot(folder)
			if err != nil {
				wm.logger.Warn("failed to convert workspace folder", "uri", folder.Uri, "error", err)
				continue
			}
			wm.roots = append(wm.roots, root)
			wm.logger.Debug("added workspace root from folder", "uri", root.URI, "name", root.Name)
		}
		return nil
	}

	// Priority 2: RootUri (deprecated but still supported)
	if params.RootUri != nil && *params.RootUri != "" {
		root, err := wm.uriToWorkspaceRoot(*params.RootUri, "")
		if err != nil {
			wm.logger.Error("failed to convert rootUri to workspace root", "rootUri", *params.RootUri, "error", err)
			return err
		}
		wm.roots = append(wm.roots, root)
		wm.logger.Debug("added workspace root from rootUri", "uri", root.URI)
		return nil
	}

	// Priority 3: RootPath (deprecated but still supported)
	if params.RootPath != nil && *params.RootPath != "" {
		// Convert path to file:// URI
		absPath, err := filepath.Abs(*params.RootPath)
		if err != nil {
			wm.logger.Error("failed to get absolute path", "rootPath", *params.RootPath, "error", err)
			return err
		}

		uri := pathToFileURI(absPath)
		root, err := wm.uriToWorkspaceRoot(uri, "")
		if err != nil {
			wm.logger.Error("failed to convert rootPath to workspace root", "rootPath", *params.RootPath, "error", err)
			return err
		}
		wm.roots = append(wm.roots, root)
		wm.logger.Debug("added workspace root from rootPath", "path", *params.RootPath, "uri", root.URI)
		return nil
	}

	wm.logger.Info("no workspace roots found in initialize params")
	return nil
}

// workspaceFolderToRoot converts an LSP WorkspaceFolder to WorkspaceRoot
func (wm *WorkspaceManager) workspaceFolderToRoot(folder lsp.WorkspaceFolder) (WorkspaceRoot, error) {
	return wm.uriToWorkspaceRoot(folder.Uri, folder.Name)
}

// uriToWorkspaceRoot converts a URI to WorkspaceRoot
func (wm *WorkspaceManager) uriToWorkspaceRoot(uri, name string) (WorkspaceRoot, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return WorkspaceRoot{}, err
	}

	if parsedURI.Scheme != "file" {
		return WorkspaceRoot{}, nil // Skip non-file URIs
	}

	path := parsedURI.Path
	if name == "" {
		name = filepath.Base(path)
	}

	return WorkspaceRoot{
		URI:  uri,
		Path: path,
		Name: name,
	}, nil
}

// pathToFileURI converts a local path to a file:// URI
func pathToFileURI(path string) string {
	// Ensure path is absolute and clean
	absPath, _ := filepath.Abs(filepath.Clean(path))

	// Convert to URL path format
	urlPath := filepath.ToSlash(absPath)

	// Add file:// scheme
	return "file://" + urlPath
}

// GetWorkspaceRoots returns the current workspace roots
func (wm *WorkspaceManager) GetWorkspaceRoots() []WorkspaceRoot {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	roots := make([]WorkspaceRoot, len(wm.roots))
	copy(roots, wm.roots)
	return roots
}

// GetFileIndex returns a copy of the file index
func (wm *WorkspaceManager) GetFileIndex() map[string]*FileInfo {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	index := make(map[string]*FileInfo)
	for uri, info := range wm.fileIndex {
		// Create a copy of the FileInfo
		indexCopy := *info
		index[uri] = &indexCopy
	}
	return index
}

// GetMarkdownFiles returns all indexed Markdown files
func (wm *WorkspaceManager) GetMarkdownFiles() []*FileInfo {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	files := make([]*FileInfo, 0, len(wm.fileIndex))
	for _, info := range wm.fileIndex {
		// Create a copy of the FileInfo
		infoCopy := *info
		files = append(files, &infoCopy)
	}
	return files
}

// isExcludedPath checks if a path should be excluded from indexing
func (wm *WorkspaceManager) isExcludedPath(path string) bool {
	// Check against exclusion patterns
	for _, pattern := range wm.excludePatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	// Skip hidden directories (starting with .)
	if strings.HasPrefix(filepath.Base(path), ".") {
		return true
	}

	return false
}

// isMarkdownFile checks if a file is a Markdown file
func (wm *WorkspaceManager) isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md"
}

// DiscoverMarkdownFiles scans all workspace roots for Markdown files
func (wm *WorkspaceManager) DiscoverMarkdownFiles() error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	// Clear existing index
	oldCount := len(wm.fileIndex)
	wm.fileIndex = make(map[string]*FileInfo)
	if oldCount > 0 {
		wm.logger.Debug("cleared workspace cache for rediscovery", "previousFiles", oldCount)
	}

	fileCount := 0
	for _, root := range wm.roots {
		rootPath := root.Path
		wm.logger.Info("scanning workspace root for Markdown files", "root", root.Name, "path", rootPath)

		err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				wm.logger.Warn("error accessing path during workspace scan", "path", path, "error", err)
				return nil // Continue walking despite errors
			}

			// Skip directories
			if d.IsDir() {
				// Check if this directory should be excluded
				if wm.isExcludedPath(path) {
					wm.logger.Debug("skipping excluded directory", "path", path)
					return filepath.SkipDir
				}
				return nil
			}

			// Check if this is a Markdown file
			if !wm.isMarkdownFile(path) {
				return nil
			}

			// Check file count limit
			if fileCount >= wm.maxFileCount {
				wm.logger.Warn("reached maximum file count limit, stopping scan",
					"limit", wm.maxFileCount, "root", root.Name)
				return filepath.SkipAll
			}

			// Get file info
			info, err := d.Info()
			if err != nil {
				wm.logger.Warn("failed to get file info", "path", path, "error", err)
				return nil
			}

			// Create file URI
			absPath, err := filepath.Abs(path)
			if err != nil {
				wm.logger.Warn("failed to get absolute path", "path", path, "error", err)
				return nil
			}
			uri := pathToFileURI(absPath)

			// Create relative path from workspace root
			relPath, err := filepath.Rel(rootPath, path)
			if err != nil {
				relPath = path // Fallback to absolute path
			}

			// Add to index
			fileInfo := &FileInfo{
				URI:     uri,
				Path:    relPath,
				ModTime: info.ModTime(),
				Size:    info.Size(),
			}

			wm.fileIndex[uri] = fileInfo
			fileCount++

			wm.logger.Debug("added file to workspace cache", "uri", uri, "path", relPath, "size", info.Size())
			return nil
		})

		if err != nil {
			wm.logger.Error("failed to walk workspace root", "root", root.Name, "path", rootPath, "error", err)
			continue
		}
	}

	wm.logger.Info("workspace Markdown file discovery completed",
		"totalFiles", fileCount, "roots", len(wm.roots))

	return nil
}

// RefreshFileIndex performs an incremental refresh of the file index
func (wm *WorkspaceManager) RefreshFileIndex() error {
	// For now, perform a full rescan
	// TODO: Implement incremental refresh based on file modification times
	return wm.DiscoverMarkdownFiles()
}

// AddFileToIndex adds a single file to the index
func (wm *WorkspaceManager) AddFileToIndex(uri string) error {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	// Parse URI to get path
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return err
	}

	if parsedURI.Scheme != "file" {
		return nil // Skip non-file URIs
	}

	path := parsedURI.Path

	// Check if it's a Markdown file
	if !wm.isMarkdownFile(path) {
		return nil
	}

	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Find the appropriate workspace root
	var relPath string
	for _, root := range wm.roots {
		if strings.HasPrefix(path, root.Path) {
			var err error
			relPath, err = filepath.Rel(root.Path, path)
			if err != nil {
				relPath = path
			}
			break
		}
	}
	if relPath == "" {
		relPath = path // Fallback if no root matches
	}

	// Add to index
	fileInfo := &FileInfo{
		URI:     uri,
		Path:    relPath,
		ModTime: info.ModTime(),
		Size:    info.Size(),
	}

	wm.fileIndex[uri] = fileInfo
	wm.logger.Debug("added file to workspace cache", "uri", uri, "path", relPath, "size", info.Size())

	return nil
}

// RemoveFileFromIndex removes a file from the index
func (wm *WorkspaceManager) RemoveFileFromIndex(uri string) {
	wm.mutex.Lock()
	defer wm.mutex.Unlock()

	if fileInfo, exists := wm.fileIndex[uri]; exists {
		delete(wm.fileIndex, uri)
		wm.logger.Debug("removed file from workspace cache", "uri", uri, "path", fileInfo.Path)
	}
}
