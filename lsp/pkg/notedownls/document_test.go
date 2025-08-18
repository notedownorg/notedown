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

package notedownls

import (
	"testing"
)

func TestNewDocument(t *testing.T) {
	tests := []struct {
		name        string
		uri         string
		wantPath    string
		shouldError bool
	}{
		{
			name:     "file URI",
			uri:      "file:///path/to/document.md",
			wantPath: "document",
		},
		{
			name:     "file URI with relative path",
			uri:      "file://./README.md",
			wantPath: "README",
		},
		{
			name:     "non-file URI",
			uri:      "http://example.com/doc",
			wantPath: "",
		},
		{
			name:     "file with multiple extensions",
			uri:      "file:///path/to/archive.tar.gz",
			wantPath: "archive.tar",
		},
		{
			name:     "file without extension",
			uri:      "file:///path/to/LICENSE",
			wantPath: "LICENSE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument(tt.uri)
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if doc.URI != tt.uri {
				t.Errorf("Expected URI %s, got %s", tt.uri, doc.URI)
			}

			if doc.Basepath != tt.wantPath {
				t.Errorf("Expected basepath %s, got %s", tt.wantPath, doc.Basepath)
			}
		})
	}
}
