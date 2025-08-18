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
)

// FindWorkspaceRoot searches for a .notedown directory starting from the given path
// and walking up the directory tree. Returns the path containing .notedown directory
// or empty string if not found.
func FindWorkspaceRoot(startPath string) (string, error) {
	// Convert to absolute path to ensure consistent behavior
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", err
	}

	// Walk up the directory tree
	currentPath := absPath
	for {
		notedownDir := filepath.Join(currentPath, ".notedown")

		// Check if .notedown directory exists
		if info, err := os.Stat(notedownDir); err == nil && info.IsDir() {
			return currentPath, nil
		}

		// Move to parent directory
		parentPath := filepath.Dir(currentPath)

		// If we've reached the root and haven't found .notedown, stop
		if parentPath == currentPath {
			break
		}

		currentPath = parentPath
	}

	return "", nil // Not found
}

// FindConfigFile searches for settings.yaml or settings.json in the .notedown directory
// starting from the given path. Returns the full path to the config file or empty string if not found.
func FindConfigFile(startPath string) (string, error) {
	workspaceRoot, err := FindWorkspaceRoot(startPath)
	if err != nil {
		return "", err
	}

	if workspaceRoot == "" {
		return "", nil // No workspace found
	}

	notedownDir := filepath.Join(workspaceRoot, ".notedown")

	// Check for settings.yaml first (preferred)
	yamlPath := filepath.Join(notedownDir, "settings.yaml")
	if _, err := os.Stat(yamlPath); err == nil {
		return yamlPath, nil
	}

	// Check for settings.json as fallback
	jsonPath := filepath.Join(notedownDir, "settings.json")
	if _, err := os.Stat(jsonPath); err == nil {
		return jsonPath, nil
	}

	return "", nil // No config file found
}

// HasWorkspaceConfig checks if a workspace configuration exists for the given path
func HasWorkspaceConfig(startPath string) (bool, error) {
	configPath, err := FindConfigFile(startPath)
	if err != nil {
		return false, err
	}
	return configPath != "", nil
}
