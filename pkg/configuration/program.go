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
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

//go:embed program_default.yaml
var defaultProgramConfiguration []byte

type ProgramConfiguration struct {
	Workspaces       map[string]Workspace `json:"workspaces"`
	DefaultWorkspace string               `json:"default_workspace"`
}

type Workspace struct {
	Location string `json:"location"`
}

// Handle expansion of ~ if present
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("failed to get home directory, defaulting to current directory", "error", err)
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

const programConfigurationPath = ".config/notedown/config.yaml"

func EnsureProgramConfiguration() (*ProgramConfiguration, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	path := filepath.Join(home, programConfigurationPath)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for program configuration file: %w", err)
		}

		if err := os.WriteFile(path, defaultProgramConfiguration, 0644); err != nil {
			return nil, fmt.Errorf("failed to write program configuration file: %w", err)
		}
		slog.Info("initialized program configuration", "path", path)
	}
	return loadProgramConfiguration(path)
}

func loadProgramConfiguration(path string) (*ProgramConfiguration, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read program configuration file at %s: %w", path, err)
	}

	var config ProgramConfiguration
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal program configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("program configuration is invalid: %w", err)
	}

	return &config, nil
}

func DefaultWorkspace() (*Workspace, error) {
	programConfig, err := EnsureProgramConfiguration()
	if err != nil {
		return nil, err
	}
	// Ensure errors if the default workspace is not set or does not exist so this is safe
	ws := programConfig.Workspaces[programConfig.DefaultWorkspace]
	return &ws, nil
}

func (c ProgramConfiguration) Validate() error {
	// Workspaces
	if c.Workspaces == nil || len(c.Workspaces) == 0 {
		return fmt.Errorf("atleast one workspace must be configured")
	}
	for name, workspace := range c.Workspaces {
		if workspace.Location == "" {
			return fmt.Errorf("location is required for workspace \"%s\"", name)
		}
	}
	if c.DefaultWorkspace == "" {
		return fmt.Errorf("default workspace is required")
	}
	if _, ok := c.Workspaces[c.DefaultWorkspace]; !ok {
		return fmt.Errorf("default_workspace \"%s\" does not exist", c.DefaultWorkspace)
	}

	return nil
}
