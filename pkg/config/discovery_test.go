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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindWorkspaceRoot(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()

	// Create nested directories: tempDir/project/src/deep
	projectDir := filepath.Join(tempDir, "project")
	srcDir := filepath.Join(projectDir, "src")
	deepDir := filepath.Join(srcDir, "deep")

	err := os.MkdirAll(deepDir, 0755)
	require.NoError(t, err)

	// Create .notedown directory in project root
	notedownDir := filepath.Join(projectDir, ".notedown")
	err = os.MkdirAll(notedownDir, 0755)
	require.NoError(t, err)

	tests := []struct {
		name        string
		startPath   string
		expected    string
		expectError bool
	}{
		{
			name:      "find from project root",
			startPath: projectDir,
			expected:  projectDir,
		},
		{
			name:      "find from src directory",
			startPath: srcDir,
			expected:  projectDir,
		},
		{
			name:      "find from deep directory",
			startPath: deepDir,
			expected:  projectDir,
		},
		{
			name:      "not found from temp root",
			startPath: tempDir,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FindWorkspaceRoot(tt.startPath)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFindConfigFile(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()

	// Create project with .notedown directory
	projectDir := filepath.Join(tempDir, "project")
	notedownDir := filepath.Join(projectDir, ".notedown")
	err := os.MkdirAll(notedownDir, 0755)
	require.NoError(t, err)

	// Create a subdirectory
	srcDir := filepath.Join(projectDir, "src")
	err = os.MkdirAll(srcDir, 0755)
	require.NoError(t, err)

	tests := []struct {
		name         string
		setupFunc    func()
		startPath    string
		expectedFile string
	}{
		{
			name: "yaml file found",
			setupFunc: func() {
				yamlPath := filepath.Join(notedownDir, "settings.yaml")
				err := os.WriteFile(yamlPath, []byte("tasks:\n  states: []"), 0644)
				require.NoError(t, err)
			},
			startPath:    srcDir,
			expectedFile: filepath.Join(notedownDir, "settings.yaml"),
		},
		{
			name: "json file found when yaml absent",
			setupFunc: func() {
				// Remove yaml file if exists
				yamlPath := filepath.Join(notedownDir, "settings.yaml")
				os.Remove(yamlPath)

				jsonPath := filepath.Join(notedownDir, "settings.json")
				err := os.WriteFile(jsonPath, []byte(`{"tasks":{"states":[]}}`), 0644)
				require.NoError(t, err)
			},
			startPath:    srcDir,
			expectedFile: filepath.Join(notedownDir, "settings.json"),
		},
		{
			name: "yaml preferred over json",
			setupFunc: func() {
				yamlPath := filepath.Join(notedownDir, "settings.yaml")
				err := os.WriteFile(yamlPath, []byte("tasks:\n  states: []"), 0644)
				require.NoError(t, err)

				jsonPath := filepath.Join(notedownDir, "settings.json")
				err = os.WriteFile(jsonPath, []byte(`{"tasks":{"states":[]}}`), 0644)
				require.NoError(t, err)
			},
			startPath:    srcDir,
			expectedFile: filepath.Join(notedownDir, "settings.yaml"),
		},
		{
			name: "no config file found",
			setupFunc: func() {
				// Remove all config files
				yamlPath := filepath.Join(notedownDir, "settings.yaml")
				jsonPath := filepath.Join(notedownDir, "settings.json")
				os.Remove(yamlPath)
				os.Remove(jsonPath)
			},
			startPath:    srcDir,
			expectedFile: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupFunc()

			result, err := FindConfigFile(tt.startPath)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedFile, result)
		})
	}
}

func TestHasWorkspaceConfig(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()

	// Project with config
	projectWithConfig := filepath.Join(tempDir, "with-config")
	notedownDir := filepath.Join(projectWithConfig, ".notedown")
	err := os.MkdirAll(notedownDir, 0755)
	require.NoError(t, err)

	yamlPath := filepath.Join(notedownDir, "settings.yaml")
	err = os.WriteFile(yamlPath, []byte("tasks:\n  states: []"), 0644)
	require.NoError(t, err)

	// Project without config
	projectWithoutConfig := filepath.Join(tempDir, "without-config")
	err = os.MkdirAll(projectWithoutConfig, 0755)
	require.NoError(t, err)

	tests := []struct {
		name      string
		startPath string
		expected  bool
	}{
		{
			name:      "has config",
			startPath: projectWithConfig,
			expected:  true,
		},
		{
			name:      "no config",
			startPath: projectWithoutConfig,
			expected:  false,
		},
		{
			name:      "no workspace",
			startPath: tempDir,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HasWorkspaceConfig(tt.startPath)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
