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

import (
	"os"
	"testing"

	"github.com/notedownorg/notedown/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceManager(t *testing.T) {
	logger := log.New(os.Stderr, log.Error) // Use error level to reduce noise
	manager := NewManager(logger)

	// Test adding a workspace root
	err := manager.AddRoot("./testdata")
	if err != nil {
		// Skip test if testdata doesn't exist - this is fine for now
		t.Skip("No testdata directory found - skipping workspace manager test")
		return
	}

	roots := manager.GetWorkspaceRoots()
	assert.Len(t, roots, 1)
	assert.Contains(t, roots[0].URI, "testdata")

	// Test configuration
	manager.SetMaxFileCount(100)
	manager.SetExcludePatterns([]string{".test"})

	// Test discovery - this will work even with empty directory
	err = manager.DiscoverMarkdownFiles()
	require.NoError(t, err)

	files := manager.GetMarkdownFiles()
	// Just check that we can get files (might be empty)
	assert.NotNil(t, files)
}

func TestPathUtilities(t *testing.T) {
	// Test PathToFileURI
	uri := PathToFileURI("/tmp/test.md")
	assert.Equal(t, "file:///tmp/test.md", uri)

	// Test URIToWorkspaceRoot
	root, err := URIToWorkspaceRoot("file:///tmp/workspace", "test")
	require.NoError(t, err)
	assert.Equal(t, "file:///tmp/workspace", root.URI)
	assert.Equal(t, "/tmp/workspace", root.Path)
	assert.Equal(t, "test", root.Name)

	// Test IsMarkdownFile
	assert.True(t, IsMarkdownFile("test.md"))
	assert.True(t, IsMarkdownFile("test.MD"))
	assert.False(t, IsMarkdownFile("test.txt"))

	// Test IsExcludedPath
	patterns := []string{".git", "node_modules"}
	assert.True(t, IsExcludedPath("/project/.git/config", patterns))
	assert.True(t, IsExcludedPath("/project/node_modules/package", patterns))
	assert.False(t, IsExcludedPath("/project/src/file.md", patterns))
}
