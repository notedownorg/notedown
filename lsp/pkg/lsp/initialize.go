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

type InitializeParams struct {
	// The process Id of the parent process that started the server.
	// Is null if the process has not been started by another process.
	// If the parent process is not alive then the server should exit.
	ProcessId *int `json:"processId"`

	// Information about the client
	ClientInfo *ClientInfo `json:"clientInfo"`

	// The locale the client is currently showing the user interface in.
	// This must not necessarily be the locale of the operating system.
	Locale *string `json:"locale,omitempty"`

	// The rootPath of the workspace. Is null if no folder is open.
	// @deprecated in favour of rootUri.
	RootPath *string `json:"rootPath,omitempty"`

	// The rootUri of the workspace. Is null if no folder is open.
	// If both `rootPath` and `rootUri` are set `rootUri` wins.
	// @deprecated in favour of workspaceFolders.
	RootUri *string `json:"rootUri,omitempty"`

	// The workspace folders configured in the client when the server starts.
	// This property is only available if the client supports workspace folders.
	// It can be `null` if the client supports workspace folders but none are
	// configured.
	// @since 3.6.0
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders,omitempty"`

	// User provided initialization options.
	InitializationOptions any `json:"initializationOptions,omitempty"`

	// The capabilities provided by the client (editor or tool)
	Capabilities ClientCapabilities `json:"capabilities"`

	// The initial trace setting. If omitted trace is disabled ('off').
	Trace *string `json:"trace,omitempty"`
}

// WorkspaceFolder represents a workspace folder in the client
type WorkspaceFolder struct {
	// The associated URI for this workspace folder.
	Uri string `json:"uri"`

	// The name of the workspace folder. Used to refer to this
	// workspace folder in the user interface.
	Name string `json:"name"`
}

type InitializeResult struct {
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
	Capabilities ServerCapabilities `json:"capabilities"`
}
