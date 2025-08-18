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

// ServerInfo contains information about the server.
type ServerInfo struct {
	// The name of the server as defined by the server.
	Name string `json:"name"`
	// The server's version as defined by the server.
	Version string `json:"version"`
}

// ServerCapabilities defines capabilities provided by a language server.
type ServerCapabilities struct {
	// The position encoding the server picked from the encodings offered by the client.
	// If the client didn't provide encodings, only 'utf-16' is valid.
	// Defaults to 'utf-16' if omitted.
	PositionEncoding *string `json:"positionEncoding,omitempty"`

	// Defines how text documents are synced. Can be a detailed structure or
	// a TextDocumentSyncKind number. Defaults to TextDocumentSyncKind.None if omitted.
	TextDocumentSync any `json:"textDocumentSync,omitempty"` // TextDocumentSyncOptions | TextDocumentSyncKind

	// Defines how notebook documents are synced.
	NotebookDocumentSync any `json:"notebookDocumentSync,omitempty"` // NotebookDocumentSyncOptions | NotebookDocumentSyncRegistrationOptions

	// The server provides completion support.
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`

	// The server provides hover support.
	HoverProvider any `json:"hoverProvider,omitempty"` // boolean | HoverOptions

	// The server provides signature help support.
	SignatureHelpProvider *SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`

	// The server provides go to declaration support.
	DeclarationProvider any `json:"declarationProvider,omitempty"` // boolean | DeclarationOptions | DeclarationRegistrationOptions

	// The server provides goto definition support.
	DefinitionProvider any `json:"definitionProvider,omitempty"` // boolean | DefinitionOptions

	// The server provides goto type definition support.
	TypeDefinitionProvider any `json:"typeDefinitionProvider,omitempty"` // boolean | TypeDefinitionOptions | TypeDefinitionRegistrationOptions

	// The server provides goto implementation support.
	ImplementationProvider any `json:"implementationProvider,omitempty"` // boolean | ImplementationOptions | ImplementationRegistrationOptions

	// The server provides find references support.
	ReferencesProvider any `json:"referencesProvider,omitempty"` // boolean | ReferenceOptions

	// The server provides document highlight support.
	DocumentHighlightProvider any `json:"documentHighlightProvider,omitempty"` // boolean | DocumentHighlightOptions

	// The server provides document symbol support.
	DocumentSymbolProvider any `json:"documentSymbolProvider,omitempty"` // boolean | DocumentSymbolOptions

	// The server provides code action support.
	CodeActionProvider any `json:"codeActionProvider,omitempty"` // boolean | CodeActionOptions

	// The server provides code lens support.
	CodeLensProvider *CodeLensOptions `json:"codeLensProvider,omitempty"`

	// The server provides document link support.
	DocumentLinkProvider *DocumentLinkOptions `json:"documentLinkProvider,omitempty"`

	// The server provides color provider support.
	ColorProvider any `json:"colorProvider,omitempty"` // boolean | DocumentColorOptions | DocumentColorRegistrationOptions

	// The server provides document formatting support.
	DocumentFormattingProvider any `json:"documentFormattingProvider,omitempty"` // boolean | DocumentFormattingOptions

	// The server provides document range formatting support.
	DocumentRangeFormattingProvider any `json:"documentRangeFormattingProvider,omitempty"` // boolean | DocumentRangeFormattingOptions

	// The server provides document on type formatting support.
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`

	// The server provides rename support.
	RenameProvider any `json:"renameProvider,omitempty"` // boolean | RenameOptions

	// The server provides folding provider support.
	FoldingRangeProvider any `json:"foldingRangeProvider,omitempty"` // boolean | FoldingRangeOptions | FoldingRangeRegistrationOptions

	// The server provides execute command support.
	ExecuteCommandProvider *ExecuteCommandOptions `json:"executeCommandProvider,omitempty"`

	// The server provides selection range support.
	SelectionRangeProvider any `json:"selectionRangeProvider,omitempty"` // boolean | SelectionRangeOptions | SelectionRangeRegistrationOptions

	// The server provides linked editing range support.
	LinkedEditingRangeProvider any `json:"linkedEditingRangeProvider,omitempty"` // boolean | LinkedEditingRangeOptions | LinkedEditingRangeRegistrationOptions

	// The server provides call hierarchy support.
	CallHierarchyProvider any `json:"callHierarchyProvider,omitempty"` // boolean | CallHierarchyOptions | CallHierarchyRegistrationOptions

	// The server provides semantic tokens support.
	SemanticTokensProvider *SemanticTokensOptions `json:"semanticTokensProvider,omitempty"`

	// The server provides moniker support.
	MonikerProvider any `json:"monikerProvider,omitempty"` // boolean | MonikerOptions | MonikerRegistrationOptions

	// The server provides type hierarchy support.
	TypeHierarchyProvider any `json:"typeHierarchyProvider,omitempty"` // boolean | TypeHierarchyOptions | TypeHierarchyRegistrationOptions

	// The server provides inline value support.
	InlineValueProvider any `json:"inlineValueProvider,omitempty"` // boolean | InlineValueOptions | InlineValueRegistrationOptions

	// The server provides inlay hint support.
	InlayHintProvider any `json:"inlayHintProvider,omitempty"` // boolean | InlayHintOptions | InlayHintRegistrationOptions

	// The server provides pull model diagnostic support.
	DiagnosticProvider any `json:"diagnosticProvider,omitempty"` // DiagnosticOptions | DiagnosticRegistrationOptions

	// The server provides workspace symbol support.
	WorkspaceSymbolProvider any `json:"workspaceSymbolProvider,omitempty"` // boolean | WorkspaceSymbolOptions

	// Workspace specific server capabilities.
	Workspace *WorkspaceServerCapabilities `json:"workspace,omitempty"`

	// Experimental server capabilities.
	Experimental any `json:"experimental,omitempty"`
}

// TextDocumentSyncKind defines how text documents are synced.
type TextDocumentSyncKind int

const (
	// TextDocumentSyncKindNone means documents are not synced at all.
	TextDocumentSyncKindNone TextDocumentSyncKind = iota
	// TextDocumentSyncKindFull means documents are synced by always sending the full content.
	TextDocumentSyncKindFull
	// TextDocumentSyncKindIncremental means documents are synced by sending incremental changes.
	TextDocumentSyncKindIncremental
)

// TextDocumentSyncOptions defines text document synchronization options.
type TextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server. If omitted open close notification should not be sent.
	OpenClose *bool `json:"openClose,omitempty"`
	// Change notifications are sent to the server. See TextDocumentSyncKind.None, TextDocumentSyncKind.Full and TextDocumentSyncKind.Incremental. If omitted it defaults to TextDocumentSyncKind.None.
	Change *TextDocumentSyncKind `json:"change,omitempty"`
	// If present will save notifications are sent to the server. If omitted the notification should not be sent.
	WillSave *bool `json:"willSave,omitempty"`
	// If present will save wait until requests are sent to the server. If omitted the request should not be sent.
	WillSaveWaitUntil *bool `json:"willSaveWaitUntil,omitempty"`
	// If present save notifications are sent to the server. If omitted the notification should not be sent.
	Save *SaveOptions `json:"save,omitempty"`
}

// SaveOptions defines save options.
type SaveOptions struct {
	// The client is supposed to include the content on save.
	IncludeText *bool `json:"includeText,omitempty"`
}

// WorkspaceServerCapabilities defines workspace specific server capabilities.
type WorkspaceServerCapabilities struct {
	// The server supports workspace folder. @since 3.6.0
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
	// The server is interested in file notifications/requests. @since 3.16.0
	FileOperations *FileOperationsServerCapabilities `json:"fileOperations,omitempty"`
}

// WorkspaceFoldersServerCapabilities defines workspace folders server capabilities.
type WorkspaceFoldersServerCapabilities struct {
	// The server has support for workspace folders
	Supported *bool `json:"supported,omitempty"`
	// Whether the server wants to receive workspace folder change notifications.
	ChangeNotifications any `json:"changeNotifications,omitempty"` // string | boolean
}

// FileOperationsServerCapabilities defines file operations server capabilities.
type FileOperationsServerCapabilities struct {
	// The server is interested in receiving didCreateFiles notifications.
	DidCreate *FileOperationRegistrationOptions `json:"didCreate,omitempty"`
	// The server is interested in receiving willCreateFiles requests.
	WillCreate *FileOperationRegistrationOptions `json:"willCreate,omitempty"`
	// The server is interested in receiving didRenameFiles notifications.
	DidRename *FileOperationRegistrationOptions `json:"didRename,omitempty"`
	// The server is interested in receiving willRenameFiles requests.
	WillRename *FileOperationRegistrationOptions `json:"willRename,omitempty"`
	// The server is interested in receiving didDeleteFiles notifications.
	DidDelete *FileOperationRegistrationOptions `json:"didDelete,omitempty"`
	// The server is interested in receiving willDeleteFiles requests.
	WillDelete *FileOperationRegistrationOptions `json:"willDelete,omitempty"`
}

// FileOperationRegistrationOptions defines file operation registration options.
type FileOperationRegistrationOptions struct {
	// The actual filters.
	Filters []FileOperationFilter `json:"filters"`
}

// FileOperationFilter defines a filter to describe in which file operation requests or notifications the server is interested in.
type FileOperationFilter struct {
	// A Uri like file or folder is matched against this. To match all file scheme URI use a string.
	Scheme *string `json:"scheme,omitempty"`
	// The actual file operation pattern.
	Pattern FileOperationPattern `json:"pattern"`
}

// FileOperationPattern defines a file operation pattern.
type FileOperationPattern struct {
	// The glob pattern to match. Glob patterns can have the following syntax:
	// - `*` to match one or more characters in a path segment
	// - `?` to match on one character in a path segment
	// - `**` to match any number of path segments, including none
	// - `{}` to group conditions (e.g. `**/*.{ts,js}` matches all TypeScript and JavaScript files)
	// - `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
	// - `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)
	Glob string `json:"glob"`
	// Whether to match files or folders with this pattern.
	Matches *FileOperationPatternKind `json:"matches,omitempty"`
	// Additional options used during matching.
	Options *FileOperationPatternOptions `json:"options,omitempty"`
}

// FileOperationPatternKind defines the kind of file operation pattern.
type FileOperationPatternKind string

const (
	// The pattern matches a file only.
	FileOperationPatternKindFile FileOperationPatternKind = "file"
	// The pattern matches a folder only.
	FileOperationPatternKindFolder FileOperationPatternKind = "folder"
)

// FileOperationPatternOptions defines additional options used during matching.
type FileOperationPatternOptions struct {
	// The pattern should be matched ignoring casing.
	IgnoreCase *bool `json:"ignoreCase,omitempty"`
}

// Placeholder types for server capability options
// These would be fully defined based on specific LSP feature requirements when implementing specific features

// CompletionOptions defines server capability options for completion
type CompletionOptions struct {
	// The characters that trigger completion automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	// The list of all possible characters that commit a completion. This field can be used
	// if clients don't support individual commit characters per completion item. See
	// `ClientCapabilities.textDocument.completion.completionItem.commitCharactersSupport`
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`
	// The server provides support to resolve additional
	// information for a completion item.
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
	// The server supports the following `CompletionItem` specific
	// capabilities.
	CompletionItem *CompletionItemOptions `json:"completionItem,omitempty"`
}

// CompletionItemOptions defines completion item specific options
type CompletionItemOptions struct {
	// The server has support for completion item label
	// details (see also `CompletionItemLabelDetails`) when
	// receiving a completion item in a resolve call.
	LabelDetailsSupport *bool `json:"labelDetailsSupport,omitempty"`
}

// HoverOptions defines server capability options for hover
type HoverOptions struct{}

// SignatureHelpOptions defines server capability options for signature help
type SignatureHelpOptions struct{}

// DeclarationOptions defines server capability options for go to declaration
type DeclarationOptions struct{}

// DefinitionOptions defines server capability options for go to definition
type DefinitionOptions struct{}

// TypeDefinitionOptions defines server capability options for go to type definition
type TypeDefinitionOptions struct{}

// ImplementationOptions defines server capability options for go to implementation
type ImplementationOptions struct{}

// ReferenceOptions defines server capability options for find references
type ReferenceOptions struct{}

// DocumentHighlightOptions defines server capability options for document highlight
type DocumentHighlightOptions struct{}

// DocumentSymbolOptions defines server capability options for document symbols
type DocumentSymbolOptions struct{}

// CodeActionOptions defines server capability options for code actions
type CodeActionOptions struct{}

// CodeLensOptions defines server capability options for code lens
type CodeLensOptions struct{}

// DocumentLinkOptions defines server capability options for document links
type DocumentLinkOptions struct{}

// DocumentColorOptions defines server capability options for document colors
type DocumentColorOptions struct{}

// DocumentFormattingOptions defines server capability options for document formatting
type DocumentFormattingOptions struct{}

// DocumentRangeFormattingOptions defines server capability options for document range formatting
type DocumentRangeFormattingOptions struct{}

// DocumentOnTypeFormattingOptions defines server capability options for document on-type formatting
type DocumentOnTypeFormattingOptions struct{}

// RenameOptions defines server capability options for rename
type RenameOptions struct{}

// FoldingRangeOptions defines server capability options for folding range
type FoldingRangeOptions struct{}

// ExecuteCommandOptions defines server capability options for execute command
type ExecuteCommandOptions struct {
	// The commands to be executed on the server
	Commands []string `json:"commands"`
}

// SelectionRangeOptions defines server capability options for selection range
type SelectionRangeOptions struct{}

// LinkedEditingRangeOptions defines server capability options for linked editing range
type LinkedEditingRangeOptions struct{}

// CallHierarchyOptions defines server capability options for call hierarchy
type CallHierarchyOptions struct{}

// SemanticTokensOptions defines server capability options for semantic tokens
type SemanticTokensOptions struct{}

// MonikerOptions defines server capability options for moniker
type MonikerOptions struct{}

// TypeHierarchyOptions defines server capability options for type hierarchy
type TypeHierarchyOptions struct{}

// InlineValueOptions defines server capability options for inline value
type InlineValueOptions struct{}

// InlayHintOptions defines server capability options for inlay hint
type InlayHintOptions struct{}

// DiagnosticOptions defines server capability options for diagnostics
type DiagnosticOptions struct{}

// WorkspaceSymbolOptions defines server capability options for workspace symbols
type WorkspaceSymbolOptions struct{}

// Registration option types for dynamic capability registration

// DeclarationRegistrationOptions defines registration options for go to declaration
type DeclarationRegistrationOptions struct{}

// TypeDefinitionRegistrationOptions defines registration options for go to type definition
type TypeDefinitionRegistrationOptions struct{}

// ImplementationRegistrationOptions defines registration options for go to implementation
type ImplementationRegistrationOptions struct{}

// DocumentColorRegistrationOptions defines registration options for document color
type DocumentColorRegistrationOptions struct{}

// FoldingRangeRegistrationOptions defines registration options for folding range
type FoldingRangeRegistrationOptions struct{}

// SelectionRangeRegistrationOptions defines registration options for selection range
type SelectionRangeRegistrationOptions struct{}

// LinkedEditingRangeRegistrationOptions defines registration options for linked editing range
type LinkedEditingRangeRegistrationOptions struct{}

// CallHierarchyRegistrationOptions defines registration options for call hierarchy
type CallHierarchyRegistrationOptions struct{}

// MonikerRegistrationOptions defines registration options for moniker
type MonikerRegistrationOptions struct{}

// TypeHierarchyRegistrationOptions defines registration options for type hierarchy
type TypeHierarchyRegistrationOptions struct{}

// InlineValueRegistrationOptions defines registration options for inline value
type InlineValueRegistrationOptions struct{}

// InlayHintRegistrationOptions defines registration options for inlay hint
type InlayHintRegistrationOptions struct{}

// DiagnosticRegistrationOptions defines registration options for diagnostics
type DiagnosticRegistrationOptions struct{}

// Notebook document synchronization types

// NotebookDocumentSyncOptions defines options for notebook document synchronization
type NotebookDocumentSyncOptions struct{}

// NotebookDocumentSyncRegistrationOptions defines registration options for notebook document synchronization
type NotebookDocumentSyncRegistrationOptions struct{}
