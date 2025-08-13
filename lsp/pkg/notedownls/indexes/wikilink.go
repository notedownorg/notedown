package indexes

import (
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/notedownorg/notedown/pkg/log"
)

// WikilinkTargetInfo contains information about a wikilink target
type WikilinkTargetInfo struct {
	Target       string          // The wikilink target (e.g., "project-alpha", "docs/api")
	Exists       bool            // Whether the target file actually exists
	ReferencedBy map[string]bool // Set of document URIs that reference this target
	LastSeen     time.Time       // When this target was last seen during scanning
	SuggestedURI string          // Suggested file URI if this target were to be created
}

// WikilinkIndex manages all wikilink targets across the workspace
type WikilinkIndex struct {
	targets map[string]*WikilinkTargetInfo // target -> info
	mutex   sync.RWMutex
	logger  *log.Logger
}

// NewWikilinkIndex creates a new wikilink index
func NewWikilinkIndex(logger *log.Logger) *WikilinkIndex {
	return &WikilinkIndex{
		targets: make(map[string]*WikilinkTargetInfo),
		logger:  logger,
	}
}

// AddTarget adds or updates a wikilink target in the index
func (wi *WikilinkIndex) AddTarget(target, sourceURI string, exists bool) {
	wi.mutex.Lock()
	defer wi.mutex.Unlock()

	targetInfo, found := wi.targets[target]
	if !found {
		targetInfo = &WikilinkTargetInfo{
			Target:       target,
			Exists:       exists,
			ReferencedBy: make(map[string]bool),
			LastSeen:     time.Now(),
		}
		wi.targets[target] = targetInfo
	}

	// Update existence status (prioritize true if any source claims it exists)
	if exists {
		targetInfo.Exists = true
	}

	// Add source document to references
	targetInfo.ReferencedBy[sourceURI] = true
	targetInfo.LastSeen = time.Now()

	// Generate suggested URI for non-existent targets
	if !targetInfo.Exists {
		targetInfo.SuggestedURI = wi.generateSuggestedURI(target)
	}

	wi.logger.Debug("added wikilink target",
		"target", target,
		"exists", targetInfo.Exists,
		"references", len(targetInfo.ReferencedBy))
}

// RemoveTargetReference removes a reference to a target from a specific document
func (wi *WikilinkIndex) RemoveTargetReference(target, sourceURI string) {
	wi.mutex.Lock()
	defer wi.mutex.Unlock()

	if targetInfo, found := wi.targets[target]; found {
		delete(targetInfo.ReferencedBy, sourceURI)

		// If no references remain and target doesn't exist, remove it
		if len(targetInfo.ReferencedBy) == 0 && !targetInfo.Exists {
			delete(wi.targets, target)
			wi.logger.Debug("removed unused wikilink target", "target", target)
		}
	}
}

// GetAllTargets returns all wikilink targets, optionally filtered by existence
func (wi *WikilinkIndex) GetAllTargets() map[string]*WikilinkTargetInfo {
	wi.mutex.RLock()
	defer wi.mutex.RUnlock()

	// Return a copy to avoid concurrent access issues
	result := make(map[string]*WikilinkTargetInfo)
	for target, info := range wi.targets {
		// Create a copy of the target info
		infoCopy := &WikilinkTargetInfo{
			Target:       info.Target,
			Exists:       info.Exists,
			ReferencedBy: make(map[string]bool),
			LastSeen:     info.LastSeen,
			SuggestedURI: info.SuggestedURI,
		}
		for uri := range info.ReferencedBy {
			infoCopy.ReferencedBy[uri] = true
		}
		result[target] = infoCopy
	}
	return result
}

// GetNonExistentTargets returns all targets that don't correspond to existing files
func (wi *WikilinkIndex) GetNonExistentTargets() map[string]*WikilinkTargetInfo {
	wi.mutex.RLock()
	defer wi.mutex.RUnlock()

	result := make(map[string]*WikilinkTargetInfo)
	for target, info := range wi.targets {
		if !info.Exists {
			// Create a copy
			infoCopy := &WikilinkTargetInfo{
				Target:       info.Target,
				Exists:       info.Exists,
				ReferencedBy: make(map[string]bool),
				LastSeen:     info.LastSeen,
				SuggestedURI: info.SuggestedURI,
			}
			for uri := range info.ReferencedBy {
				infoCopy.ReferencedBy[uri] = true
			}
			result[target] = infoCopy
		}
	}
	return result
}

// UpdateTargetExistence updates whether a target exists based on file system changes
func (wi *WikilinkIndex) UpdateTargetExistence(target string, exists bool) {
	wi.mutex.Lock()
	defer wi.mutex.Unlock()

	if targetInfo, found := wi.targets[target]; found {
		targetInfo.Exists = exists
		if exists {
			targetInfo.SuggestedURI = "" // Clear suggested URI since it now exists
		} else {
			targetInfo.SuggestedURI = wi.generateSuggestedURI(target)
		}
		wi.logger.Debug("updated wikilink target existence", "target", target, "exists", exists)
	}
}

// generateSuggestedURI generates a suggested file URI for a wikilink target
func (wi *WikilinkIndex) generateSuggestedURI(target string) string {
	// If target already has an extension, use as-is
	if filepath.Ext(target) != "" {
		return target
	}

	// Add .md extension for markdown files
	suggestedPath := target + ".md"

	// Ensure proper path separators
	suggestedPath = strings.ReplaceAll(suggestedPath, "\\", "/")

	return suggestedPath
}

// GetReferenceCount returns the number of documents that reference a target
func (wi *WikilinkIndex) GetReferenceCount(target string) int {
	wi.mutex.RLock()
	defer wi.mutex.RUnlock()

	if targetInfo, found := wi.targets[target]; found {
		return len(targetInfo.ReferencedBy)
	}
	return 0
}

// Clear removes all targets from the index
func (wi *WikilinkIndex) Clear() {
	wi.mutex.Lock()
	defer wi.mutex.Unlock()

	wi.targets = make(map[string]*WikilinkTargetInfo)
	wi.logger.Debug("cleared wikilink index")
}

// GetTargetsByPrefix returns targets that start with the given prefix
func (wi *WikilinkIndex) GetTargetsByPrefix(prefix string) map[string]*WikilinkTargetInfo {
	wi.mutex.RLock()
	defer wi.mutex.RUnlock()

	result := make(map[string]*WikilinkTargetInfo)
	lowerPrefix := strings.ToLower(prefix)

	for target, info := range wi.targets {
		if strings.HasPrefix(strings.ToLower(target), lowerPrefix) {
			// Create a copy
			infoCopy := &WikilinkTargetInfo{
				Target:       info.Target,
				Exists:       info.Exists,
				ReferencedBy: make(map[string]bool),
				LastSeen:     info.LastSeen,
				SuggestedURI: info.SuggestedURI,
			}
			for uri := range info.ReferencedBy {
				infoCopy.ReferencedBy[uri] = true
			}
			result[target] = infoCopy
		}
	}
	return result
}

// WorkspaceFile represents basic file information for workspace operations
type WorkspaceFile interface {
	GetURI() string
	GetPath() string
}

// ExtractWikilinksFromDocument parses a document and extracts all wikilink targets
func (wi *WikilinkIndex) ExtractWikilinksFromDocument(content, documentURI string, workspaceFiles map[string]WorkspaceFile) []string {
	// Use regex to extract wikilinks for now (simpler approach while parser issues are resolved)
	wikilinkRegex := regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)
	matches := wikilinkRegex.FindAllStringSubmatch(content, -1)

	var targets []string

	for _, match := range matches {
		if len(match) > 1 {
			target := strings.TrimSpace(match[1])
			targets = append(targets, target)

			// Check if this target corresponds to an existing file
			exists := wi.targetExistsInWorkspace(target, workspaceFiles)

			// Add to index
			wi.AddTarget(target, documentURI, exists)

			wi.logger.Debug("extracted wikilink",
				"target", target,
				"exists", exists,
				"source", documentURI)
		}
	}

	wi.logger.Debug("extracted wikilinks from document",
		"uri", documentURI,
		"count", len(targets))

	return targets
}

// targetExistsInWorkspace checks if a wikilink target corresponds to an existing file
func (wi *WikilinkIndex) targetExistsInWorkspace(target string, workspaceFiles map[string]WorkspaceFile) bool {
	// Direct match: target matches a file's path without extension
	for _, fileInfo := range workspaceFiles {
		path := fileInfo.GetPath()
		// Check if target matches file path without extension
		pathWithoutExt := strings.TrimSuffix(path, filepath.Ext(path))
		if target == pathWithoutExt {
			return true
		}

		// Check if target matches just the filename without extension
		baseWithoutExt := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		if target == baseWithoutExt {
			return true
		}
	}

	return false
}

// RefreshDocumentWikilinks removes old wikilink references for a document and re-extracts them
func (wi *WikilinkIndex) RefreshDocumentWikilinks(content, documentURI string, workspaceFiles map[string]WorkspaceFile) {
	// Remove existing references from this document
	wi.removeDocumentReferences(documentURI)

	// Extract new wikilinks
	wi.ExtractWikilinksFromDocument(content, documentURI, workspaceFiles)
}

// removeDocumentReferences removes all wikilink references from a specific document
func (wi *WikilinkIndex) removeDocumentReferences(documentURI string) {
	wi.mutex.Lock()
	defer wi.mutex.Unlock()

	// Find all targets that reference this document and remove the reference
	targetsToRemove := make([]string, 0)

	for target, info := range wi.targets {
		if info.ReferencedBy[documentURI] {
			delete(info.ReferencedBy, documentURI)

			// If no references remain and target doesn't exist, mark for removal
			if len(info.ReferencedBy) == 0 && !info.Exists {
				targetsToRemove = append(targetsToRemove, target)
			}
		}
	}

	// Remove targets with no references
	for _, target := range targetsToRemove {
		delete(wi.targets, target)
		wi.logger.Debug("removed unreferenced wikilink target", "target", target)
	}

	wi.logger.Debug("removed document references",
		"uri", documentURI,
		"removedTargets", len(targetsToRemove))
}
