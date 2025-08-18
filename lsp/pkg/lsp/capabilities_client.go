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

// ClientInfo contains information about the client.
type ClientInfo struct {
	// The name of the client as defined by the client.
	Name string `json:"name"`
	// The client's version as defined by the client.
	Version string `json:"version"`
}

// ClientCapabilities defines capabilities provided by the client.
type ClientCapabilities struct {
	// Workspace specific client capabilities.
	Workspace *WorkspaceClientCapabilities `json:"workspace,omitempty"`
	// Text document specific client capabilities.
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	// Capabilities specific to the notebook document support.
	NotebookDocument *NotebookDocumentClientCapabilities `json:"notebookDocument,omitempty"`
	// Window specific client capabilities.
	Window *WindowClientCapabilities `json:"window,omitempty"`
	// General client capabilities.
	General *GeneralClientCapabilities `json:"general,omitempty"`
	// Experimental client capabilities.
	Experimental any `json:"experimental,omitempty"`
}

// WorkspaceClientCapabilities defines workspace specific client capabilities.
type WorkspaceClientCapabilities struct {
	// The client supports applying batch edits to the workspace by supporting the request 'workspace/applyEdit'
	ApplyEdit *bool `json:"applyEdit,omitempty"`
	// Capabilities specific to `WorkspaceEdit`s.
	WorkspaceEdit *WorkspaceEditClientCapabilities `json:"workspaceEdit,omitempty"`
	// Capabilities specific to the `workspace/didChangeConfiguration` notification.
	DidChangeConfiguration *DidChangeConfigurationClientCapabilities `json:"didChangeConfiguration,omitempty"`
	// Capabilities specific to the `workspace/didChangeWatchedFiles` notification.
	DidChangeWatchedFiles *DidChangeWatchedFilesClientCapabilities `json:"didChangeWatchedFiles,omitempty"`
	// Capabilities specific to the `workspace/symbol` request.
	Symbol *WorkspaceSymbolClientCapabilities `json:"symbol,omitempty"`
	// Capabilities specific to the `workspace/executeCommand` request.
	ExecuteCommand *ExecuteCommandClientCapabilities `json:"executeCommand,omitempty"`
	// The client has support for workspace folders.
	WorkspaceFolders *bool `json:"workspaceFolders,omitempty"`
	// The client supports `workspace/configuration` requests.
	Configuration *bool `json:"configuration,omitempty"`
	// Capabilities specific to the semantic token requests scoped to the workspace.
	SemanticTokens *SemanticTokensWorkspaceClientCapabilities `json:"semanticTokens,omitempty"`
	// Capabilities specific to the code lens requests scoped to the workspace.
	CodeLens *CodeLensWorkspaceClientCapabilities `json:"codeLens,omitempty"`
	// The client has support for file requests/notifications.
	FileOperations *WorkspaceFileOperationsClientCapabilities `json:"fileOperations,omitempty"`
	// Client workspace capabilities specific to inline values.
	InlineValue *InlineValueWorkspaceClientCapabilities `json:"inlineValue,omitempty"`
	// Client workspace capabilities specific to inlay hints.
	InlayHint *InlayHintWorkspaceClientCapabilities `json:"inlayHint,omitempty"`
	// Client workspace capabilities specific to diagnostics.
	Diagnostics *DiagnosticWorkspaceClientCapabilities `json:"diagnostics,omitempty"`
}

// TextDocumentClientCapabilities defines text document specific client capabilities.
type TextDocumentClientCapabilities struct {
	// Capabilities specific to the `textDocument/didOpen`, `textDocument/didChange` and `textDocument/didClose` notifications.
	Synchronization *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
	// Capabilities specific to the `textDocument/completion` request.
	Completion *CompletionClientCapabilities `json:"completion,omitempty"`
	// Capabilities specific to the `textDocument/hover` request.
	Hover *HoverClientCapabilities `json:"hover,omitempty"`
	// Capabilities specific to the `textDocument/signatureHelp` request.
	SignatureHelp *SignatureHelpClientCapabilities `json:"signatureHelp,omitempty"`
	// Capabilities specific to the `textDocument/declaration` request.
	Declaration *DeclarationClientCapabilities `json:"declaration,omitempty"`
	// Capabilities specific to the `textDocument/definition` request.
	Definition *DefinitionClientCapabilities `json:"definition,omitempty"`
	// Capabilities specific to the `textDocument/typeDefinition` request.
	TypeDefinition *TypeDefinitionClientCapabilities `json:"typeDefinition,omitempty"`
	// Capabilities specific to the `textDocument/implementation` request.
	Implementation *ImplementationClientCapabilities `json:"implementation,omitempty"`
	// Capabilities specific to the `textDocument/references` request.
	References *ReferenceClientCapabilities `json:"references,omitempty"`
	// Capabilities specific to the `textDocument/documentHighlight` request.
	DocumentHighlight *DocumentHighlightClientCapabilities `json:"documentHighlight,omitempty"`
	// Capabilities specific to the `textDocument/documentSymbol` request.
	DocumentSymbol *DocumentSymbolClientCapabilities `json:"documentSymbol,omitempty"`
	// Capabilities specific to the `textDocument/codeAction` request.
	CodeAction *CodeActionClientCapabilities `json:"codeAction,omitempty"`
	// Capabilities specific to the `textDocument/codeLens` request.
	CodeLens *CodeLensClientCapabilities `json:"codeLens,omitempty"`
	// Capabilities specific to the `textDocument/documentLink` request.
	DocumentLink *DocumentLinkClientCapabilities `json:"documentLink,omitempty"`
	// Capabilities specific to the `textDocument/documentColor` and the `textDocument/colorPresentation` request.
	ColorProvider *DocumentColorClientCapabilities `json:"colorProvider,omitempty"`
	// Capabilities specific to the `textDocument/formatting` request.
	Formatting *DocumentFormattingClientCapabilities `json:"formatting,omitempty"`
	// Capabilities specific to the `textDocument/rangeFormatting` request.
	RangeFormatting *DocumentRangeFormattingClientCapabilities `json:"rangeFormatting,omitempty"`
	// Capabilities specific to the `textDocument/onTypeFormatting` request.
	OnTypeFormatting *DocumentOnTypeFormattingClientCapabilities `json:"onTypeFormatting,omitempty"`
	// Capabilities specific to the `textDocument/rename` request.
	Rename *RenameClientCapabilities `json:"rename,omitempty"`
	// Capabilities specific to the `textDocument/publishDiagnostics` notification.
	PublishDiagnostics *PublishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`
	// Capabilities specific to the `textDocument/foldingRange` request.
	FoldingRange *FoldingRangeClientCapabilities `json:"foldingRange,omitempty"`
	// Capabilities specific to the `textDocument/selectionRange` request.
	SelectionRange *SelectionRangeClientCapabilities `json:"selectionRange,omitempty"`
	// Capabilities specific to the `textDocument/linkedEditingRange` request.
	LinkedEditingRange *LinkedEditingRangeClientCapabilities `json:"linkedEditingRange,omitempty"`
	// Capabilities specific to the various call hierarchy requests.
	CallHierarchy *CallHierarchyClientCapabilities `json:"callHierarchy,omitempty"`
	// Capabilities specific to the various semantic token requests.
	SemanticTokens *SemanticTokensClientCapabilities `json:"semanticTokens,omitempty"`
	// Capabilities specific to the `textDocument/moniker` request.
	Moniker *MonikerClientCapabilities `json:"moniker,omitempty"`
	// Capabilities specific to the various type hierarchy requests.
	TypeHierarchy *TypeHierarchyClientCapabilities `json:"typeHierarchy,omitempty"`
	// Capabilities specific to the `textDocument/inlineValue` request.
	InlineValue *InlineValueClientCapabilities `json:"inlineValue,omitempty"`
	// Capabilities specific to the `textDocument/inlayHint` request.
	InlayHint *InlayHintClientCapabilities `json:"inlayHint,omitempty"`
	// Capabilities specific to the diagnostic pull model.
	Diagnostic *DiagnosticClientCapabilities `json:"diagnostic,omitempty"`
}

// NotebookDocumentClientCapabilities defines notebook document specific client capabilities.
type NotebookDocumentClientCapabilities struct {
	// Capabilities specific to notebook document synchronization.
	Synchronization *NotebookDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
}

// WindowClientCapabilities defines window specific client capabilities.
type WindowClientCapabilities struct {
	// It indicates whether the client supports server initiated progress using the `window/workDoneProgress/create` request.
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
	// Capabilities specific to the showMessage request.
	ShowMessage *ShowMessageRequestClientCapabilities `json:"showMessage,omitempty"`
	// Client capabilities for the show document request.
	ShowDocument *ShowDocumentClientCapabilities `json:"showDocument,omitempty"`
}

// GeneralClientCapabilities defines general client capabilities.
type GeneralClientCapabilities struct {
	// Client capability that signals how the client handles stale requests (e.g. a request in a notebook cell that got edited is not cancelled when the cell is executed again).
	StaleRequestSupport *StaleRequestSupportClientCapabilities `json:"staleRequestSupport,omitempty"`
	// Client capabilities specific to regular expressions.
	RegularExpressions *RegularExpressionsClientCapabilities `json:"regularExpressions,omitempty"`
	// Client capabilities specific to the client's markdown parser.
	Markdown *MarkdownClientCapabilities `json:"markdown,omitempty"`
	// The position encodings supported by the client. Client and server have to agree on the same position encoding to ensure that offsets (e.g. character position in a line) are interpreted the same on both sides.
	PositionEncodings []string `json:"positionEncodings,omitempty"`
}

// WorkspaceFileOperationsClientCapabilities defines client capabilities for workspace file operations.
type WorkspaceFileOperationsClientCapabilities struct {
	// Whether the client supports dynamic registration for file requests/notifications.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client has support for sending didCreateFiles notifications.
	DidCreate *bool `json:"didCreate,omitempty"`
	// The client has support for sending willCreateFiles requests.
	WillCreate *bool `json:"willCreate,omitempty"`
	// The client has support for sending didRenameFiles notifications.
	DidRename *bool `json:"didRename,omitempty"`
	// The client has support for sending willRenameFiles requests.
	WillRename *bool `json:"willRename,omitempty"`
	// The client has support for sending didDeleteFiles notifications.
	DidDelete *bool `json:"didDelete,omitempty"`
	// The client has support for sending willDeleteFiles requests.
	WillDelete *bool `json:"willDelete,omitempty"`
}

// StaleRequestSupportClientCapabilities defines client capabilities for stale request support.
type StaleRequestSupportClientCapabilities struct {
	// The client will actively cancel the request.
	Cancel bool `json:"cancel"`
	// The list of requests for which the client will retry the request if it receives a response with error code `ContentModified`
	RetryOnContentModified []string `json:"retryOnContentModified"`
}

// Placeholder types for capability sub-structures
// These would be fully defined based on specific LSP feature requirements when implementing specific features

// WorkspaceEditClientCapabilities defines client capabilities for workspace edits
type WorkspaceEditClientCapabilities struct{}

// DidChangeConfigurationClientCapabilities defines client capabilities for workspace configuration changes
type DidChangeConfigurationClientCapabilities struct{}

// DidChangeWatchedFilesClientCapabilities defines client capabilities for watching files
type DidChangeWatchedFilesClientCapabilities struct{}

// WorkspaceSymbolClientCapabilities defines client capabilities for workspace symbol requests
type WorkspaceSymbolClientCapabilities struct{}

// ExecuteCommandClientCapabilities defines client capabilities for execute command requests
type ExecuteCommandClientCapabilities struct{}

// SemanticTokensWorkspaceClientCapabilities defines client capabilities for workspace semantic tokens
type SemanticTokensWorkspaceClientCapabilities struct{}

// CodeLensWorkspaceClientCapabilities defines client capabilities for workspace code lens
type CodeLensWorkspaceClientCapabilities struct{}

// InlineValueWorkspaceClientCapabilities defines client capabilities for workspace inline values
type InlineValueWorkspaceClientCapabilities struct{}

// InlayHintWorkspaceClientCapabilities defines client capabilities for workspace inlay hints
type InlayHintWorkspaceClientCapabilities struct{}

// DiagnosticWorkspaceClientCapabilities defines client capabilities for workspace diagnostics
type DiagnosticWorkspaceClientCapabilities struct{}

// TextDocumentSyncClientCapabilities defines client capabilities for text document synchronization
type TextDocumentSyncClientCapabilities struct{}

// CompletionClientCapabilities defines client capabilities for completion requests
type CompletionClientCapabilities struct{}

// HoverClientCapabilities defines client capabilities for hover requests
type HoverClientCapabilities struct{}

// SignatureHelpClientCapabilities defines client capabilities for signature help requests
type SignatureHelpClientCapabilities struct{}

// DeclarationClientCapabilities defines client capabilities for declaration requests
type DeclarationClientCapabilities struct{}

// DefinitionClientCapabilities defines client capabilities for definition requests
type DefinitionClientCapabilities struct{}

// TypeDefinitionClientCapabilities defines client capabilities for type definition requests
type TypeDefinitionClientCapabilities struct{}

// ImplementationClientCapabilities defines client capabilities for implementation requests
type ImplementationClientCapabilities struct{}

// ReferenceClientCapabilities defines client capabilities for reference requests
type ReferenceClientCapabilities struct{}

// DocumentHighlightClientCapabilities defines client capabilities for document highlight requests
type DocumentHighlightClientCapabilities struct{}

// DocumentSymbolClientCapabilities defines client capabilities for document symbol requests
type DocumentSymbolClientCapabilities struct{}

// CodeActionClientCapabilities defines client capabilities for code action requests
type CodeActionClientCapabilities struct{}

// CodeLensClientCapabilities defines client capabilities for code lens requests
type CodeLensClientCapabilities struct{}

// DocumentLinkClientCapabilities defines client capabilities for document link requests
type DocumentLinkClientCapabilities struct{}

// DocumentColorClientCapabilities defines client capabilities for document color requests
type DocumentColorClientCapabilities struct{}

// DocumentFormattingClientCapabilities defines client capabilities for document formatting requests
type DocumentFormattingClientCapabilities struct{}

// DocumentRangeFormattingClientCapabilities defines client capabilities for document range formatting requests
type DocumentRangeFormattingClientCapabilities struct{}

// DocumentOnTypeFormattingClientCapabilities defines client capabilities for document on-type formatting requests
type DocumentOnTypeFormattingClientCapabilities struct{}

// RenameClientCapabilities defines client capabilities for rename requests
type RenameClientCapabilities struct{}

// PublishDiagnosticsClientCapabilities defines client capabilities for publish diagnostics notifications
type PublishDiagnosticsClientCapabilities struct{}

// FoldingRangeClientCapabilities defines client capabilities for folding range requests
type FoldingRangeClientCapabilities struct{}

// SelectionRangeClientCapabilities defines client capabilities for selection range requests
type SelectionRangeClientCapabilities struct{}

// LinkedEditingRangeClientCapabilities defines client capabilities for linked editing range requests
type LinkedEditingRangeClientCapabilities struct{}

// CallHierarchyClientCapabilities defines client capabilities for call hierarchy requests
type CallHierarchyClientCapabilities struct{}

// SemanticTokensClientCapabilities defines client capabilities for semantic token requests
type SemanticTokensClientCapabilities struct{}

// MonikerClientCapabilities defines client capabilities for moniker requests
type MonikerClientCapabilities struct{}

// TypeHierarchyClientCapabilities defines client capabilities for type hierarchy requests
type TypeHierarchyClientCapabilities struct{}

// InlineValueClientCapabilities defines client capabilities for inline value requests
type InlineValueClientCapabilities struct{}

// InlayHintClientCapabilities defines client capabilities for inlay hint requests
type InlayHintClientCapabilities struct{}

// DiagnosticClientCapabilities defines client capabilities for diagnostic requests
type DiagnosticClientCapabilities struct{}

// NotebookDocumentSyncClientCapabilities defines client capabilities for notebook document synchronization
type NotebookDocumentSyncClientCapabilities struct{}

// ShowMessageRequestClientCapabilities defines client capabilities for show message requests
type ShowMessageRequestClientCapabilities struct{}

// ShowDocumentClientCapabilities defines client capabilities for show document requests
type ShowDocumentClientCapabilities struct{}

// RegularExpressionsClientCapabilities defines client capabilities for regular expressions
type RegularExpressionsClientCapabilities struct{}

// MarkdownClientCapabilities defines client capabilities for markdown parsing
type MarkdownClientCapabilities struct{}
