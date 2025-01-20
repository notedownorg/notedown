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

package configuration

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultWorkspace = &WorkspaceConfiguration{
		Sources: Sources{
			DefaultDirectory: "sources",
		},
	}
	fullWorkspace = &WorkspaceConfiguration{
		Sources: Sources{
			DefaultDirectory: "library",
		},
	}
)

func TestNewWorkspaceConfiguration(t *testing.T) {
	tests := []struct {
		file    string
		want    *WorkspaceConfiguration
		wantErr bool
	}{
		{
			file: "testdata/workspace/full.yaml",
			want: fullWorkspace,
		},
		{
			file: "testdata/workspace/minimal.yaml",
			want: &WorkspaceConfiguration{
				Sources: Sources{
					DefaultDirectory: "sources",
				},
			},
		},
		{
			file: "testdata/workspace/empty.yaml",
			want: &WorkspaceConfiguration{
				Sources: Sources{
					DefaultDirectory: "",
				},
			},
		},
		{
			file:    "testdata/workspace/invalid.yaml",
			wantErr: true,
		},
		{
			file:    "testdata/workspace/missing.yaml",
			wantErr: true,
		},
		{
			file: "testdata/workspace/source_defaultdir_empty.yaml", // valid, it means in workspace root
			want: &WorkspaceConfiguration{
				Sources: Sources{
					DefaultDirectory: "",
				},
			},
		},
		{
			file: "workspace_default.yaml", // dont accidentally create an invalid configuration
			want: defaultWorkspace,
		},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			got, err := loadWorkspaceConfiguration(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnsureWorkspaceConfiguration(t *testing.T) {
	tmpdir := t.TempDir()

	// Test that the default configuration is created
	config, err := EnsureWorkspaceConfiguration(tmpdir)
	assert.NoError(t, err)
	assert.Equal(t, defaultWorkspace, config)

	// Overwrite the configuration file with a custom one and check we don't overwrite it
	full, _ := os.ReadFile("testdata/workspace/full.yaml")
	assert.NoError(t, os.WriteFile(filepath.Join(tmpdir, workspaceConfigurationPath), full, 0644))
	config, err = EnsureWorkspaceConfiguration(tmpdir)
	assert.NoError(t, err)
	assert.Equal(t, fullWorkspace, config)

}
