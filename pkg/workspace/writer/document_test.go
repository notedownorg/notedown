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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/notedownorg/notedown/pkg/configuration"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/test"
	"github.com/notedownorg/notedown/pkg/workspace"
	"github.com/notedownorg/notedown/pkg/workspace/writer"
	"github.com/stretchr/testify/assert"
)

func unique(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.md", uuid.New().String()))
}

func TestAddDocument(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		metadata  workspace.Metadata
		blocks    []ast.Block
		wantErr   bool
		wantFinal []byte
	}{
		{
			name:      "Parent exists and file does not exist",
			path:      unique("parent"),
			metadata:  workspace.Metadata{"type": "note"},
			blocks:    Bl(P("Hello, world!")),
			wantErr:   false,
			wantFinal: []byte("---\ntype: note\n---\nHello, world!\n"),
		},
		{
			name:    "Parent exists and file exists",
			path:    "parent/existing.md",
			wantErr: true,
		},
		{
			name:      "Parent does not exist",
			path:      "newparent/new.md",
			metadata:  workspace.Metadata{"type": "note"},
			blocks:    Bl(P("Hello, world!")),
			wantErr:   false,
			wantFinal: []byte("---\ntype: note\n---\nHello, world!\n"),
		},
		{
			name:      "Metadata is nil",
			path:      unique("parent"),
			blocks:    Bl(P("Hello, world!")),
			wantErr:   false,
			wantFinal: []byte("Hello, world!\n"),
		},
		{
			name:      "Metadata is empty",
			path:      unique("parent"),
			metadata:  workspace.Metadata{},
			blocks:    Bl(P("Hello, world!")),
			wantErr:   false,
			wantFinal: []byte("---\n---\nHello, world!\n"),
		},
		{
			name:      "Content is nil",
			path:      unique("parent"),
			metadata:  workspace.Metadata{"type": "note"},
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
		ws := &configuration.Workspace{Location: dir}
		client := writer.NewClient(ws)

		t.Run(tt.name, func(t *testing.T) {
			err := client.Create(workspace.NewDocument(tt.path, tt.metadata, tt.blocks...))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check the file contents
			contents, err := os.ReadFile(filepath.Join(dir, tt.path))
			assert.NoError(t, err)
			assert.Equal(t, string(tt.wantFinal), string(contents))
		})
	}
}

func TestUpdateDocument(t *testing.T) {
	dir, err := copyTestData("TestUpdateDocument")
	assert.NoError(t, err)
	ws := &configuration.Workspace{Location: dir}
	client := writer.NewClient(ws)
	fullPath := filepath.Join(dir, "basic.md")
	doc, _ := workspace.LoadDocument(dir, "basic.md", time.Now())
	// Not sleeping here caused GH runners to fail stale check on linux
	// As the whole test runs too quickly for nanosecond precision 🤯
	time.Sleep(time.Millisecond)

	// Test valid update
	doc.Blocks = Bl(P("some new content"))
	err = client.Update(doc)
	assert.NoError(t, err)
	contents, _ := os.ReadFile(fullPath)
	assert.Equal(t, "some new content\n", string(contents))

	// Now that we updated the file the original one is out of date
	// Validate that if we try to update it again we get an error and the file is not modified further
	doc.Blocks = Bl(P("some even newer content"))
	err = client.Update(doc)
	assert.Error(t, err)
	contents, _ = os.ReadFile(fullPath)
	assert.Equal(t, "some new content\n", string(contents))

	// Try to update a file that does not exist
	doc = workspace.NewDocument("does_not_exist.md", nil, nil)
	err = client.Update(doc)
	assert.Error(t, err)
}

func TestRenameDocument(t *testing.T) {
	tests := []struct {
		name    string
		oldPath string
		newPath string
		wantErr bool
	}{
		{
			name:    "File exists",
			oldPath: "basic.md",
			newPath: "renamed.md",
		},
		{
			name:    "File does not exist",
			oldPath: "does_not_exist.md",
			newPath: "renamed.md",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		dir, err := copyTestData(t.Name())
		if err != nil {
			t.Fatalf("failed to copy test data: %v", err)
		}
		ws := &configuration.Workspace{Location: dir}
		client := writer.NewClient(ws)

		t.Run(tt.name, func(t *testing.T) {
			err := client.Rename(tt.oldPath, tt.newPath)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check the old file does not exist
			_, err = os.Stat(filepath.Join(dir, tt.oldPath))
			assert.Error(t, err)
			assert.True(t, os.IsNotExist(err))

			// Check the new file exists
			_, err = os.Stat(filepath.Join(dir, tt.newPath))
			assert.NoError(t, err)
		})
	}
}

func TestDeleteDocument(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "File exists",
			path: "basic.md",
		},
		{
			name: "File does not exist",
			path: "does_not_exist.md",
		},
	}
	for _, tt := range tests {
		dir, err := copyTestData(t.Name())
		if err != nil {
			t.Fatalf("failed to copy test data: %v", err)
		}
		ws := &configuration.Workspace{Location: dir}
		client := writer.NewClient(ws)

		t.Run(tt.name, func(t *testing.T) {
			err := client.Delete(tt.path)
			assert.NoError(t, err)

			// Check the file does not exist
			_, err = os.Stat(filepath.Join(dir, tt.path))
			assert.Error(t, err)
			assert.True(t, os.IsNotExist(err))
		})
	}
}

func deepCopy[T any](src, dst *T) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
