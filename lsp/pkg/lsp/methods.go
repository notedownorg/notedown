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

// method represents a JSON-RPC method name (private type)
type method string

// String returns the string representation of the method
func (m method) String() string {
	return string(m)
}

// JSON-RPC 2.0 and Language Server Protocol (LSP) 3.17 methods
// Reference: https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/
const (
	// Server lifecycle methods
	MethodInitialize  method = "initialize"
	MethodInitialized method = "initialized"
	MethodShutdown    method = "shutdown"
	MethodExit        method = "exit"

	// Window methods
	MethodWindowShowMessage        method = "window/showMessage"
	MethodWindowShowMessageRequest method = "window/showMessageRequest"
	MethodWindowLogMessage         method = "window/logMessage"
	MethodWindowWorkDoneProgress   method = "window/workDoneProgress/create"

	// Telemetry methods
	MethodTelemetryEvent method = "telemetry/event"

	// Client methods (sent by client to server)
	MethodClientRegisterCapability   method = "client/registerCapability"
	MethodClientUnregisterCapability method = "client/unregisterCapability"

	// Workspace methods
	MethodWorkspaceDidChangeWorkspaceFolders method = "workspace/didChangeWorkspaceFolders"
	MethodWorkspaceDidChangeConfiguration    method = "workspace/didChangeConfiguration"
	MethodWorkspaceDidChangeWatchedFiles     method = "workspace/didChangeWatchedFiles"
	MethodWorkspaceSymbol                    method = "workspace/symbol"
	MethodWorkspaceExecuteCommand            method = "workspace/executeCommand"
	MethodWorkspaceApplyEdit                 method = "workspace/applyEdit"

	// Text synchronization methods
	MethodTextDocumentDidOpen   method = "textDocument/didOpen"
	MethodTextDocumentDidChange method = "textDocument/didChange"
	MethodTextDocumentDidClose  method = "textDocument/didClose"
	MethodTextDocumentDidSave   method = "textDocument/didSave"
	MethodTextDocumentWillSave  method = "textDocument/willSave"

	// Language features methods
	MethodTextDocumentCompletion              method = "textDocument/completion"
	MethodCompletionItemResolve               method = "completionItem/resolve"
	MethodTextDocumentHover                   method = "textDocument/hover"
	MethodTextDocumentSignatureHelp           method = "textDocument/signatureHelp"
	MethodTextDocumentDeclaration             method = "textDocument/declaration"
	MethodTextDocumentDefinition              method = "textDocument/definition"
	MethodTextDocumentTypeDefinition          method = "textDocument/typeDefinition"
	MethodTextDocumentImplementation          method = "textDocument/implementation"
	MethodTextDocumentReferences              method = "textDocument/references"
	MethodTextDocumentDocumentHighlight       method = "textDocument/documentHighlight"
	MethodTextDocumentDocumentSymbol          method = "textDocument/documentSymbol"
	MethodTextDocumentCodeAction              method = "textDocument/codeAction"
	MethodCodeActionResolve                   method = "codeAction/resolve"
	MethodTextDocumentCodeLens                method = "textDocument/codeLens"
	MethodCodeLensResolve                     method = "codeLens/resolve"
	MethodTextDocumentDocumentLink            method = "textDocument/documentLink"
	MethodDocumentLinkResolve                 method = "documentLink/resolve"
	MethodTextDocumentDocumentColor           method = "textDocument/documentColor"
	MethodTextDocumentColorPresentation       method = "textDocument/colorPresentation"
	MethodTextDocumentFormatting              method = "textDocument/formatting"
	MethodTextDocumentRangeFormatting         method = "textDocument/rangeFormatting"
	MethodTextDocumentOnTypeFormatting        method = "textDocument/onTypeFormatting"
	MethodTextDocumentRename                  method = "textDocument/rename"
	MethodTextDocumentPrepareRename           method = "textDocument/prepareRename"
	MethodTextDocumentFoldingRange            method = "textDocument/foldingRange"
	MethodTextDocumentSelectionRange          method = "textDocument/selectionRange"
	MethodTextDocumentPrepareCallHierarchy    method = "textDocument/prepareCallHierarchy"
	MethodCallHierarchyIncomingCalls          method = "callHierarchy/incomingCalls"
	MethodCallHierarchyOutgoingCalls          method = "callHierarchy/outgoingCalls"
	MethodTextDocumentSemanticTokensFull      method = "textDocument/semanticTokens/full"
	MethodTextDocumentSemanticTokensFullDelta method = "textDocument/semanticTokens/full/delta"
	MethodTextDocumentSemanticTokensRange     method = "textDocument/semanticTokens/range"
	MethodWorkspaceSemanticTokensRefresh      method = "workspace/semanticTokens/refresh"

	// Diagnostic methods
	MethodTextDocumentPublishDiagnostics method = "textDocument/publishDiagnostics"
	MethodTextDocumentDiagnostic         method = "textDocument/diagnostic"

	// Progress methods
	MethodProgress method = "$/progress"

	// Cancellation methods
	MethodCancelRequest method = "$/cancelRequest"

	// LSP 3.17 specific methods
	MethodTextDocumentLinkedEditingRange   method = "textDocument/linkedEditingRange"
	MethodTextDocumentMoniker              method = "textDocument/moniker"
	MethodTextDocumentPrepareTypeHierarchy method = "textDocument/prepareTypeHierarchy"
	MethodTypeHierarchySupertypes          method = "typeHierarchy/supertypes"
	MethodTypeHierarchySubtypes            method = "typeHierarchy/subtypes"
	MethodTextDocumentInlineValue          method = "textDocument/inlineValue"
	MethodTextDocumentInlayHint            method = "textDocument/inlayHint"
	MethodInlayHintResolve                 method = "inlayHint/resolve"
	MethodWorkspaceInlayHintRefresh        method = "workspace/inlayHintRefresh"
	MethodNotebookDocumentDidOpen          method = "notebookDocument/didOpen"
	MethodNotebookDocumentDidChange        method = "notebookDocument/didChange"
	MethodNotebookDocumentDidSave          method = "notebookDocument/didSave"
	MethodNotebookDocumentDidClose         method = "notebookDocument/didClose"
)
