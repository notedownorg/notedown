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
	defaultProgram = &ProgramConfiguration{
		Workspaces: map[string]Workspace{
			"personal": {
				Location: "~/notes",
			},
		},
		DefaultWorkspace: "personal",
	}
	fullProgram = &ProgramConfiguration{
		Workspaces: map[string]Workspace{
			"personal": {
				Location: "~/notes",
			},
			"work": {
				Location: "~/worknotes",
			},
		},
		DefaultWorkspace: "personal",
	}
)

func TestNewProgramConfiguration(t *testing.T) {
	tests := []struct {
		file    string
		want    *ProgramConfiguration
		wantErr bool
	}{
		{
			file: "testdata/program/full.yaml",
			want: fullProgram,
		},
		{
			file: "testdata/program/minimal.yaml",
			want: &ProgramConfiguration{
				Workspaces: map[string]Workspace{
					"personal": {
						Location: "~/notes",
					},
				},
				DefaultWorkspace: "personal",
			},
		},
		{
			file:    "testdata/program/empty.yaml",
			wantErr: true,
		},
		{
			file:    "testdata/program/invalid.yaml",
			wantErr: true,
		},
		{
			file:    "testdata/program/missing.yaml",
			wantErr: true,
		},
		{
			file:    "testdata/program/workspace_default_notexist.yaml",
			wantErr: true,
		},
		{
			file: "program_default.yaml", // dont accidentally create an invalid configuratio
			want: defaultProgram,
		},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			got, err := loadProgramConfiguration(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnsureProgramConfiguration(t *testing.T) {
	tmpdir := t.TempDir()
	os.Setenv("HOME", tmpdir) // override the home directory so we dont mess with the real one

	// Test that the default configuration is created
	config, err := EnsureProgramConfiguration()
	assert.NoError(t, err)
	assert.Equal(t, defaultProgram, config)

	// Overwrite the configuration file with a custom one and check we don't overwrite it
	full, _ := os.ReadFile("testdata/program/full.yaml")
	assert.NoError(t, os.WriteFile(filepath.Join(tmpdir, programConfigurationPath), full, 0644))
	config, err = EnsureProgramConfiguration()
	assert.NoError(t, err)
	assert.Equal(t, fullProgram, config)

}
