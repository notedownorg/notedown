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

// CompletionParams represents the parameters for a textDocument/completion request
type CompletionParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The position inside the text document.
	Position Position `json:"position"`
	// The completion context. This is only available if the client specifies
	// to send this using the client capability `textDocument.completion.contextSupport === true`
	Context *CompletionContext `json:"context,omitempty"`
}

// CompletionContext contains additional information about the context in which a completion request is triggered
type CompletionContext struct {
	// How the completion was triggered.
	TriggerKind CompletionTriggerKind `json:"triggerKind"`
	// The trigger character (a single character) that has triggered the completion request.
	// Is undefined if `triggerKind !== CompletionTriggerKind.TriggerCharacter`
	TriggerCharacter *string `json:"triggerCharacter,omitempty"`
}

// CompletionTriggerKind defines how a completion was triggered
type CompletionTriggerKind int

const (
	// Completion was triggered by typing an identifier (24x7 code complete), manual invocation (e.g Ctrl+Space) or via API.
	CompletionTriggerKindInvoked CompletionTriggerKind = 1
	// Completion was triggered by a trigger character specified by the `triggerCharacters` properties of the `CompletionRegistrationOptions`.
	CompletionTriggerKindTriggerCharacter CompletionTriggerKind = 2
	// Completion was re-triggered as the current completion list is incomplete.
	CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// CompletionList represents a collection of completion items to be presented in the editor
type CompletionList struct {
	// This list is not complete. Further typing should result in re-computing this list.
	IsIncomplete bool `json:"isIncomplete"`
	// The completion items.
	Items []CompletionItem `json:"items"`
}

// CompletionItem represents a completion item to be presented in the editor
type CompletionItem struct {
	// The label of this completion item. By default this is also the text that is inserted when selecting this completion.
	Label string `json:"label"`
	// The kind of this completion item. Based on the kind an icon is chosen by the editor.
	Kind *CompletionItemKind `json:"kind,omitempty"`
	// A human-readable string with additional information about this item, like type or symbol information.
	Detail *string `json:"detail,omitempty"`
	// A human-readable string that represents a doc-comment.
	Documentation any `json:"documentation,omitempty"` // string | MarkupContent
	// Indicates if this item is deprecated.
	Deprecated *bool `json:"deprecated,omitempty"`
	// Select this item when showing.
	Preselect *bool `json:"preselect,omitempty"`
	// A string that should be used when comparing this item with other items.
	SortText *string `json:"sortText,omitempty"`
	// A string that should be used when filtering a set of completion items.
	FilterText *string `json:"filterText,omitempty"`
	// A string that should be inserted into a document when selecting this completion.
	InsertText *string `json:"insertText,omitempty"`
	// The format of the insert text. The format applies to both the `insertText` property
	// and the `newText` property of a provided `textEdit`. If omitted defaults to `InsertTextFormat.PlainText`.
	InsertTextFormat *InsertTextFormat `json:"insertTextFormat,omitempty"`
	// How whitespace and indentation is handled during completion item insertion.
	InsertTextMode *InsertTextMode `json:"insertTextMode,omitempty"`
	// An edit which is applied to a document when selecting this completion.
	TextEdit any `json:"textEdit,omitempty"` // TextEdit | InsertReplaceEdit
	// An optional array of additional text edits that are applied when selecting this completion.
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`
	// An optional set of characters that when pressed while this completion is active will accept it first and
	// then type that character.
	CommitCharacters []string `json:"commitCharacters,omitempty"`
	// An optional command that is executed *after* inserting this completion.
	Command *Command `json:"command,omitempty"`
	// A data entry field that is preserved between a completion and a completion resolve request.
	Data any `json:"data,omitempty"`
}

// CompletionItemKind defines the kind of a completion entry
type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

// InsertTextFormat defines whether the insert text in a completion item should be interpreted as plain text or a snippet
type InsertTextFormat int

const (
	// The primary text to be inserted is treated as a plain string.
	InsertTextFormatPlainText InsertTextFormat = 1
	// The primary text to be inserted is treated as a snippet.
	InsertTextFormatSnippet InsertTextFormat = 2
)

// InsertTextMode defines how whitespace and indentation is handled during completion item insertion
type InsertTextMode int

const (
	// The insertion or replace strings is taken as is.
	InsertTextModeAsIs InsertTextMode = 1
	// The editor adjusts leading whitespace of a text edit or insert replace edit.
	InsertTextModeAdjustIndentation InsertTextMode = 2
)

// TextEdit represents a textual edit applicable to a text document
type TextEdit struct {
	// The range of the text document to be manipulated.
	Range Range `json:"range"`
	// The string to be inserted. For delete operations use an empty string.
	NewText string `json:"newText"`
}

// Command represents a reference to a command
type Command struct {
	// Title of the command, like `save`.
	Title string `json:"title"`
	// The identifier of the actual command handler.
	Command string `json:"command"`
	// Arguments that the command handler should be invoked with.
	Arguments []any `json:"arguments,omitempty"`
}
