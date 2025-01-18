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
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"sigs.k8s.io/yaml"
)

type ProgramConfiguration struct {
	Workspaces       map[string]WorkspaceConfiguration `json:"workspaces" validate:"required,min=1"`
	DefaultWorkspace string                            `json:"default_workspace" validate:"required"`
}

type WorkspaceConfiguration struct {
	Location string `json:"location" validate:"required"`
}

const programConfigurationPath = ".notedown/config.yaml"

var validate *validator.Validate

func LoadProgramConfiguration() (*ProgramConfiguration, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	path := filepath.Join(home, programConfigurationPath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("program configuration file does not exist at %s: %w", path, err)
	}
	return NewProgramConfiguration(path)
}

func NewProgramConfiguration(path string) (*ProgramConfiguration, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read program configuration file at %s: %w", path, err)
	}

	var config ProgramConfiguration
	if err := yaml.Unmarshal(contents, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal program configuration: %w", err)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&config); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf("program configuration validation error: %s", err)
		}
	}

	return &config, nil
}
