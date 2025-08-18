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

// ExecuteCommandParams represents the parameters for workspace/executeCommand
type ExecuteCommandParams struct {
	// The identifier of the actual command handler.
	Command string `json:"command"`
	// Arguments that the command handler should be invoked with.
	Arguments []any `json:"arguments,omitempty"`
}

// WorkspaceEdit represents changes to many resources managed in the workspace
type WorkspaceEdit struct {
	// Holds changes to existing resources.
	Changes map[string][]TextEdit `json:"changes,omitempty"`
	// Depending on the client capability `workspace.workspaceEdit.resourceOperations` document changes
	// are either an array of `TextEdit`s to apply to the document or an array of document edits.
	DocumentChanges []any `json:"documentChanges,omitempty"` // (TextDocumentEdit | CreateFile | RenameFile | DeleteFile)[]
}

// ApplyWorkspaceEditParams represents the parameters for workspace/applyEdit
type ApplyWorkspaceEditParams struct {
	// An optional label of the workspace edit. This label is
	// presented in the user interface for example on an undo
	// stack to undo the workspace edit.
	Label *string `json:"label,omitempty"`
	// The edits to apply.
	Edit WorkspaceEdit `json:"edit"`
}
