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
			err := client.Create(tt.path, tt.metadata, tt.content)
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

const (
	basicChecksum                = "6c8e08c7544069890a42303050e57b0b46a4c0bb2c5dd55b0d0f7929eb0f9c51"
	emptyChecksum                = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	frontmatterChecksum          = "71946b820432d103bd6071ba7eaf80ec129295a868a9546132a8ed1ed0291a56"
	frontmatterNoContentChecksum = "cc0e785e1411c118e1341d0d69f786fca07ceb405b492ad07d26d83d142fe353"
)

func TestUpdateContentDocument(t *testing.T) {
	tests := []struct {
		name      string
		doc       writer.Document
		mutations []writer.LineMutation
		wantErr   bool
		wantFinal []byte
	}{
		{
			name: "Basic file",
			doc:  writer.Document{Path: "basic.md", Checksum: basicChecksum},
			mutations: []writer.LineMutation{
				writer.RemoveLine(5),
				writer.AddLine(1, Text("This line was added")),
				writer.UpdateLine(7, Text("This line was updated")),
			},
			wantFinal: []byte(`This line was added
This is a basic document

It has no front matter


This line was updated
`),
		},
		{
			name: "File with front matter and content",
			doc:  writer.Document{Path: "frontmatter.md", Checksum: frontmatterChecksum},
			mutations: []writer.LineMutation{
				writer.AddLine(1, Text("This line was added")),
				writer.RemoveLine(1),
				writer.UpdateLine(2, Text("This line was updated")),
			},
			wantFinal: []byte(`---
key: value
---

This line was updated
`),
		},
		{
			name: "File with front matter and no content",
			doc:  writer.Document{Path: "frontmatter_no_content.md", Checksum: frontmatterNoContentChecksum},
			mutations: []writer.LineMutation{
				writer.AddLine(1, Text("This line was added")),
				writer.UpdateLine(1, Text("This line was updated")),
			},
			wantFinal: []byte(`---
key: value
---
This line was updated
`),
		},
		{
			name: "Empty file",
			doc:  writer.Document{Path: "empty.md", Checksum: emptyChecksum},
			mutations: []writer.LineMutation{
				writer.AddLine(1, Text("This line was added")),
				writer.AddLine(2, Text("This line was added as well")),
				writer.UpdateLine(1, Text("This line was updated")),
				writer.RemoveLine(2),
			},
			wantFinal: []byte("This line was updated\n"),
		},
		{
			name:      "File no longer exists",
			doc:       writer.Document{Path: "does_not_exist.md", Checksum: ""},
			mutations: []writer.LineMutation{},
			wantErr:   true,
		},
		{
			name:      "File has been modified since last read",
			doc:       writer.Document{Path: "basic.md", Checksum: "bad_checksum"},
			mutations: []writer.LineMutation{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		dir, err := copyTestData(t.Name())
		if err != nil {
			t.Fatalf("failed to copy test data: %v", err)
		}
		client := writer.NewClient(dir)

		t.Run(tt.name, func(t *testing.T) {
			err := client.UpdateContent(tt.doc, tt.mutations...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check the file contents
			contents, err := os.ReadFile(filepath.Join(dir, tt.doc.Path))
			assert.NoError(t, err)
			assert.Equal(t, string(tt.wantFinal), string(contents))
		})
	}
}

func TestDeleteDocument(t *testing.T) {
	tests := []struct {
		name    string
		doc     writer.Document
		wantErr bool
	}{
		{
			name: "File exists",
			doc:  writer.Document{Path: "basic.md", Checksum: basicChecksum},
		},
		{
			name: "File does not exist",
			doc:  writer.Document{Path: "does_not_exist.md", Checksum: ""},
		},
		{
			name: "File has been modified since last read",
			doc:  writer.Document{Path: "basic.md", Checksum: "bad_checksum"},
		},
	}

	for _, tt := range tests {
		dir, err := copyTestData(t.Name())
		if err != nil {
			t.Fatalf("failed to copy test data: %v", err)
		}
		client := writer.NewClient(dir)

		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteDocument(tt.doc)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check the file does not exist
			_, err = os.Stat(filepath.Join(dir, tt.doc.Path))
			assert.Error(t, err)
			assert.True(t, os.IsNotExist(err))
		})
	}
}
