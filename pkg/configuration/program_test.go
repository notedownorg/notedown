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

package configuration_test

import (
	_ "embed"
	"testing"

	. "github.com/notedownorg/notedown/pkg/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewProgramConfiguration(t *testing.T) {
	tests := []struct {
		file    string
		want    *ProgramConfiguration
		wantErr bool
	}{
		{
			file: "testdata/full.yaml",
			want: &ProgramConfiguration{
				Workspaces: map[string]WorkspaceConfiguration{
					"personal": {
						Location: "~/notes",
					},
					"work": {
						Location: "~/worknotes",
					},
				},
				DefaultWorkspace: "personal",
			},
		},
		{
			file: "testdata/minimal.yaml",
			want: &ProgramConfiguration{
				Workspaces: map[string]WorkspaceConfiguration{
					"personal": {
						Location: "~/notes",
					},
				},
				DefaultWorkspace: "personal",
			},
		},
		{
			file:    "testdata/empty.yaml",
			wantErr: true,
		},
		{
			file:    "testdata/invalid.yaml",
			wantErr: true,
		},
		{
			file:    "testdata/missing.yaml",
			wantErr: true,
		},
		{
			file:    "testdata/defaultworkspace_notexist.yaml",
			wantErr: true,
		},
		{
			file: "program_default.yaml", // ensure we dont break the init command
			want: &ProgramConfiguration{
				Workspaces: map[string]WorkspaceConfiguration{
					"personal": {
						Location: "~/notes",
					},
				},
				DefaultWorkspace: "personal",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			got, err := NewProgramConfiguration(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
