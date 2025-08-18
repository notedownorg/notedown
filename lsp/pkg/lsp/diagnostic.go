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

// DiagnosticSeverity defines the severity of a diagnostic
type DiagnosticSeverity int

const (
	DiagnosticSeverityError       DiagnosticSeverity = 1
	DiagnosticSeverityWarning     DiagnosticSeverity = 2
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	DiagnosticSeverityHint        DiagnosticSeverity = 4
)

// DiagnosticTag provides additional metadata about the nature of a diagnostic
type DiagnosticTag int

const (
	DiagnosticTagUnnecessary DiagnosticTag = 1
	DiagnosticTagDeprecated  DiagnosticTag = 2
)

// DiagnosticRelatedInformation represents a related message and source code location for a diagnostic
type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// CodeDescription represents an identifier for a diagnostic code
type CodeDescription struct {
	Href string `json:"href"`
}

// Diagnostic represents a diagnostic message from the language server
type Diagnostic struct {
	Range              Range                          `json:"range"`
	Severity           *DiagnosticSeverity            `json:"severity,omitempty"`
	Code               any                            `json:"code,omitempty"` // string | number
	CodeDescription    *CodeDescription               `json:"codeDescription,omitempty"`
	Source             *string                        `json:"source,omitempty"`
	Message            string                         `json:"message"`
	Tags               []DiagnosticTag                `json:"tags,omitempty"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
	Data               any                            `json:"data,omitempty"`
}

// PublishDiagnosticsParams represents the parameters of a publishDiagnostics notification
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     *int         `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// DocumentDiagnosticReport represents the result of a textDocument/diagnostic request
type DocumentDiagnosticReport struct {
	Kind  string       `json:"kind"`
	Items []Diagnostic `json:"items"`
}
