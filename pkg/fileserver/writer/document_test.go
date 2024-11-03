// Copyright 2024 Notedown Authors
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

package writer_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/stretchr/testify/assert"
)

func unique(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.md", uuid.New().String()))
}

func TestAddDocument(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		metadata  reader.Metadata
		content   []byte
		wantErr   bool
		wantFinal []byte
	}{
		{
			name:      "Parent exists and file does not exist",
			path:      unique("parent"),
			metadata:  reader.Metadata{"type": "note"},
			content:   []byte("Hello, world!"),
			wantErr:   false,
			wantFinal: []byte("---\ntype: note\n---\nHello, world!"),
		},
		{
			name:    "Parent exists and file exists",
			path:    "parent/existing.md",
			wantErr: true,
		},
		{
			name:      "Parent does not exist",
			path:      "newparent/new.md",
			metadata:  reader.Metadata{"type": "note"},
			content:   []byte("Hello, world!"),
			wantErr:   false,
			wantFinal: []byte("---\ntype: note\n---\nHello, world!"),
		},
		{
			name:      "Metadata is nil",
			path:      unique("parent"),
			content:   []byte("Hello, world!"),
			wantErr:   false,
			wantFinal: []byte("Hello, world!"),
		},
		{
			name:      "Metadata is empty",
			path:      unique("parent"),
			metadata:  reader.Metadata{},
			content:   []byte("Hello, world!"),
			wantErr:   false,
			wantFinal: []byte("Hello, world!"),
		},
		{
			name:      "Content is nil",
			path:      unique("parent"),
			metadata:  reader.Metadata{"type": "note"},
			wantErr:   false,
			wantFinal: []byte("---\ntype: note\n---\n"),
		},
		{
			name:      "Content and metadata are nil",
			path:      unique("parent"),
			wantErr:   false,
			wantFinal: []byte(""),
		},
	}
	for _, tt := range tests {
		dir, err := copyTestData(t.Name())
		if err != nil {
			t.Fatalf("failed to copy test data: %v", err)
		}
		client := writer.NewClient(dir)

		t.Run(tt.name, func(t *testing.T) {
			err := client.AddDocument(tt.path, tt.metadata, tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check the file contents
			contents, err := os.ReadFile(filepath.Join(dir, tt.path))
			assert.NoError(t, err)
			assert.Equal(t, tt.wantFinal, contents)
		})
	}

}
