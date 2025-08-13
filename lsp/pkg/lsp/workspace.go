package lsp

// FileChangeType represents the type of change to a file
type FileChangeType int

const (
	// FileChangeTypeCreated indicates the file was created
	FileChangeTypeCreated FileChangeType = 1
	// FileChangeTypeChanged indicates the file was changed
	FileChangeTypeChanged FileChangeType = 2
	// FileChangeTypeDeleted indicates the file was deleted
	FileChangeTypeDeleted FileChangeType = 3
)

// FileEvent represents an event describing a file change
type FileEvent struct {
	// The file's URI
	URI string `json:"uri"`
	// The change type
	Type FileChangeType `json:"type"`
}

// DidChangeWatchedFilesParams represents the parameters for workspace/didChangeWatchedFiles
type DidChangeWatchedFilesParams struct {
	// The actual file events
	Changes []FileEvent `json:"changes"`
}

// FileSystemWatcher describes a file system watcher registration
type FileSystemWatcher struct {
	// The glob pattern to watch
	GlobPattern string `json:"globPattern"`
	// The kind of events of interest
	Kind *WatchKind `json:"kind,omitempty"`
}

// WatchKind represents the kind of file system events to watch
type WatchKind int

const (
	// WatchKindCreate indicates interest in create events
	WatchKindCreate WatchKind = 1
	// WatchKindChange indicates interest in change events
	WatchKindChange WatchKind = 2
	// WatchKindDelete indicates interest in delete events
	WatchKindDelete WatchKind = 4
)

// RegistrationParams represents the parameters for client/registerCapability
type RegistrationParams struct {
	Registrations []Registration `json:"registrations"`
}

// Registration represents a capability registration
type Registration struct {
	// The id used to register the request
	ID string `json:"id"`
	// The method / capability to register for
	Method string `json:"method"`
	// Options necessary for the registration
	RegisterOptions any `json:"registerOptions,omitempty"`
}

// DidChangeWatchedFilesRegistrationOptions represents options for file watching registration
type DidChangeWatchedFilesRegistrationOptions struct {
	// The watchers to register
	Watchers []FileSystemWatcher `json:"watchers"`
}
