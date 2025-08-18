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

// CodeActionParams represents the parameters for textDocument/codeAction request
type CodeActionParams struct {
	// The document in which the command was invoked.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The range for which the command was invoked.
	Range Range `json:"range"`
	// Context carrying additional information.
	Context CodeActionContext `json:"context"`
}

// CodeActionContext contains additional diagnostic information about the context in which
// a code action is run.
type CodeActionContext struct {
	// An array of diagnostics known on the client side overlapping the range provided to the
	// `textDocument/codeAction` request. They are provided so that the server knows which
	// errors are currently presented to the user for the given range. There is no guarantee
	// that these accurately reflect the error state of the resource. The primary parameter
	// to compute code actions is the provided range.
	Diagnostics []Diagnostic `json:"diagnostics"`
	// Requested kind of actions to return. Actions not of this kind are filtered out by the client before being shown.
	// So servers can omit computing them.
	Only []CodeActionKind `json:"only,omitempty"`
	// The reason why code actions were requested.
	TriggerKind *CodeActionTriggerKind `json:"triggerKind,omitempty"`
}

// CodeActionTriggerKind defines the reason why code actions were requested
type CodeActionTriggerKind int

const (
	// Code actions were explicitly requested by the user or by an extension.
	CodeActionTriggerKindInvoked CodeActionTriggerKind = 1
	// Code actions were requested automatically.
	// This typically happens when current selection in a file changes, but can
	// also be triggered when file content changes.
	CodeActionTriggerKindAutomatic CodeActionTriggerKind = 2
)

// CodeAction represents a change that can be performed in code, e.g. to fix a problem or
// to refactor code.
type CodeAction struct {
	// A short, human-readable, title for this code action.
	Title string `json:"title"`
	// The kind of the code action.
	Kind *CodeActionKind `json:"kind,omitempty"`
	// The diagnostics that this code action resolves.
	Diagnostics []Diagnostic `json:"diagnostics,omitempty"`
	// Marks this as a preferred action. Preferred actions are used by the `auto fix` command and can be targeted
	// by keybindings.
	IsPreferred *bool `json:"isPreferred,omitempty"`
	// Marks that the code action cannot currently be applied.
	Disabled *CodeActionDisabled `json:"disabled,omitempty"`
	// The workspace edit this code action performs.
	Edit *WorkspaceEdit `json:"edit,omitempty"`
	// A command this code action executes. If a code action provides an edit and a command, first the edit is
	// executed and then the command.
	Command *Command `json:"command,omitempty"`
	// A data entry field that is preserved between a code action and a resolve request.
	Data any `json:"data,omitempty"`
}

// CodeActionDisabled represents a disabled code action
type CodeActionDisabled struct {
	// Human readable description of why the code action is currently disabled.
	Reason string `json:"reason"`
}

// CodeActionKind defines the kind of a code action
type CodeActionKind string

const (
	// Base kinds
	CodeActionKindEmpty           CodeActionKind = ""
	CodeActionKindQuickFix        CodeActionKind = "quickfix"
	CodeActionKindRefactor        CodeActionKind = "refactor"
	CodeActionKindRefactorExtract CodeActionKind = "refactor.extract"
	CodeActionKindRefactorInline  CodeActionKind = "refactor.inline"
	CodeActionKindRefactorRewrite CodeActionKind = "refactor.rewrite"
	CodeActionKindSource          CodeActionKind = "source"
	CodeActionKindSourceOrganize  CodeActionKind = "source.organizeImports"
	CodeActionKindSourceFixAll    CodeActionKind = "source.fixAll"
	CodeActionKindNotebook        CodeActionKind = "notebook"
)
