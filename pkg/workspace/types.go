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

package workspace

import "time"

// WorkspaceFile represents basic file information for workspace operations
type WorkspaceFile interface {
	GetURI() string
	GetPath() string
}

// WorkspaceRoot represents a workspace root directory
type WorkspaceRoot struct {
	URI  string // file:// URI
	Path string // local filesystem path
	Name string // display name
}

// FileInfo contains lightweight metadata about a Markdown file
type FileInfo struct {
	URI     string    // file:// URI
	Path    string    // relative path from workspace root
	ModTime time.Time // last modification time
	Size    int64     // file size in bytes
}

// GetURI returns the file URI (implements WorkspaceFile)
func (f *FileInfo) GetURI() string {
	return f.URI
}

// GetPath returns the file path (implements WorkspaceFile)
func (f *FileInfo) GetPath() string {
	return f.Path
}
