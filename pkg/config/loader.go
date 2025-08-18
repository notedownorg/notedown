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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from the workspace config file.
// If no config file is found, returns the default configuration.
func LoadConfig(startPath string) (*Config, error) {
	configPath, err := FindConfigFile(startPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find config file: %w", err)
	}

	// If no config file found, return default configuration
	if configPath == "" {
		return GetDefaultConfig(), nil
	}

	return LoadConfigFromFile(configPath)
}

// LoadConfigFromFile loads configuration from a specific file path
func LoadConfigFromFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath) // #nosec G304 - configPath is from trusted config discovery
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config Config

	// Determine file format by extension
	ext := strings.ToLower(filepath.Ext(configPath))
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config file %s: %w", configPath, err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config file %s: %w", configPath, err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s (expected .yaml, .yml, or .json)", ext)
	}

	// Validate the loaded configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %w", configPath, err)
	}

	return &config, nil
}

// GetDefaultConfig returns the default configuration with standard task states
func GetDefaultConfig() *Config {
	todoDesc := "A task that needs to be completed"
	doneDesc := "A completed task"

	return &Config{
		Tasks: TasksConfig{
			States: []TaskState{
				{
					Value:       " ",
					Name:        "todo",
					Description: &todoDesc,
					Aliases:     []string{},
				},
				{
					Value:       "x",
					Name:        "done",
					Description: &doneDesc,
					Aliases:     []string{"X", "completed"},
				},
			},
		},
	}
}

// SaveConfig saves configuration to a file. The format is determined by the file extension.
func SaveConfig(config *Config, configPath string) error {
	// Validate before saving
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	var data []byte
	var err error

	// Determine file format by extension
	ext := strings.ToLower(filepath.Ext(configPath))
	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(config)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML config: %w", err)
		}
	case ".json":
		data, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s (expected .yaml, .yml, or .json)", ext)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", configPath, err)
	}

	return nil
}
