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

	"sigs.k8s.io/yaml"
)

//go:embed workspace_default.yaml
var defaultWorkspaceConfiguration []byte

type WorkspaceConfiguration struct {
	Sources Sources `json:"sources"`
}

type Sources struct {
	DefaultDirectory string `json:"default_directory"`
}

const workspaceConfigurationPath = ".notedown/config.yaml"

func EnsureWorkspaceConfiguration(location string) (*WorkspaceConfiguration, error) {
	path := filepath.Join(location, workspaceConfigurationPath)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for workspace configuration file: %w", err)
		}

		if err := os.WriteFile(path, defaultWorkspaceConfiguration, 0644); err != nil {
			return nil, fmt.Errorf("failed to write workspace configuration file: %w", err)
		}
		slog.Info("initialized workspace configuration", "path", path)
	}
	return loadWorkspaceConfiguration(path)
}

func loadWorkspaceConfiguration(path string) (*WorkspaceConfiguration, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read workspace configuration file at %s: %w", path, err)
	}

	var config WorkspaceConfiguration
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal workspace configuration: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("workspace configuration is invalid: %w", err)
	}

	return &config, nil
}

func (c WorkspaceConfiguration) Validate() error {
	return nil
}
