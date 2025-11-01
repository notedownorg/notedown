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

func TestLoadConfigFromFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name          string
		filename      string
		content       string
		expectedCfg   *Config
		expectError   bool
		errorContains string
	}{
		{
			name:     "valid yaml config",
			filename: "settings.yaml",
			content: `tasks:
  states:
    - value: " "
      name: "todo"
    - value: "x"
      name: "done"`,
			expectedCfg: &Config{
				Tasks: TasksConfig{
					States: []TaskState{
						{Value: " ", Name: "todo"},
						{Value: "x", Name: "done"},
					},
				},
			},
		},
		{
			name:     "valid json config",
			filename: "settings.json",
			content: `{
  "tasks": {
    "states": [
      {"value": " ", "name": "todo"},
      {"value": "x", "name": "done"}
    ]
  }
}`,
			expectedCfg: &Config{
				Tasks: TasksConfig{
					States: []TaskState{
						{Value: " ", Name: "todo"},
						{Value: "x", Name: "done"},
					},
				},
			},
		},
		{
			name:          "invalid yaml",
			filename:      "settings.yaml",
			content:       "invalid: yaml: content:",
			expectError:   true,
			errorContains: "failed to parse YAML",
		},
		{
			name:          "invalid json",
			filename:      "settings.json",
			content:       `{"invalid": json,}`,
			expectError:   true,
			errorContains: "failed to parse JSON",
		},
		{
			name:     "invalid config content",
			filename: "settings.yaml",
			content: `tasks:
  states: []`, // empty states - invalid
			expectError:   true,
			errorContains: "invalid configuration",
		},
		{
			name:          "unsupported file format",
			filename:      "settings.txt",
			content:       "some content",
			expectError:   true,
			errorContains: "unsupported config file format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tempDir, tt.filename)
			err := os.WriteFile(configPath, []byte(tt.content), 0600)
			require.NoError(t, err)

			result, err := LoadConfigFromFile(configPath)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCfg, result)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()

	// Setup project with config
	projectDir := filepath.Join(tempDir, "project")
	notedownDir := filepath.Join(projectDir, ".notedown")
	err := os.MkdirAll(notedownDir, 0750)
	require.NoError(t, err)

	configContent := `tasks:
  states:
    - value: " "
      name: "todo"
    - value: "x"
      name: "done"`

	configPath := filepath.Join(notedownDir, "settings.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0600)
	require.NoError(t, err)

	srcDir := filepath.Join(projectDir, "src")
	err = os.MkdirAll(srcDir, 0750)
	require.NoError(t, err)

	tests := []struct {
		name          string
		startPath     string
		expectDefault bool
	}{
		{
			name:          "load from project with config",
			startPath:     srcDir,
			expectDefault: false,
		},
		{
			name:          "load default when no config",
			startPath:     tempDir,
			expectDefault: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadConfig(tt.startPath)
			require.NoError(t, err)
			require.NotNil(t, result)

			if tt.expectDefault {
				expected := GetDefaultConfig()
				assert.Equal(t, expected, result)
			} else {
				// Should have loaded the custom config
				assert.Len(t, result.Tasks.States, 2)
				assert.Equal(t, " ", result.Tasks.States[0].Value)
				assert.Equal(t, "x", result.Tasks.States[1].Value)
			}
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	require.NotNil(t, config)
	require.Len(t, config.Tasks.States, 3)

	// Check default todo state
	todoState := config.Tasks.States[0]
	assert.Equal(t, " ", todoState.Value)
	assert.Equal(t, "todo", todoState.Name)

	// Check default done state
	doneState := config.Tasks.States[1]
	assert.Equal(t, "x", doneState.Value)
	assert.Equal(t, "done", doneState.Name)

	// Check default wip state
	wipState := config.Tasks.States[2]
	assert.Equal(t, "wip", wipState.Value)
	assert.Equal(t, "work-in-progress", wipState.Name)

	// Ensure config is valid
	err := config.Validate()
	assert.NoError(t, err)
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()

	config := &Config{
		Tasks: TasksConfig{
			States: []TaskState{
				{Value: " ", Name: "todo"},
				{Value: "x", Name: "done"},
			},
		},
	}

	tests := []struct {
		name        string
		filename    string
		expectError bool
	}{
		{
			name:     "save yaml",
			filename: "settings.yaml",
		},
		{
			name:     "save json",
			filename: "settings.json",
		},
		{
			name:        "unsupported format",
			filename:    "settings.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(tempDir, tt.filename)

			err := SaveConfig(config, configPath)

			if tt.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, configPath)

			// Verify we can load it back
			loadedConfig, err := LoadConfigFromFile(configPath)
			require.NoError(t, err)
			assert.Equal(t, config, loadedConfig)
		})
	}
}

func TestSaveConfigCreateDirectory(t *testing.T) {
	tempDir := t.TempDir()

	config := GetDefaultConfig()

	// Try to save to a non-existent directory
	configPath := filepath.Join(tempDir, "subdir", "settings.yaml")

	err := SaveConfig(config, configPath)
	require.NoError(t, err)

	// Verify file and directory were created
	assert.FileExists(t, configPath)
	assert.DirExists(t, filepath.Dir(configPath))
}
