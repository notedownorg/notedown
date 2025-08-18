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
	"fmt"
	"strings"
)

// TaskState represents a single task state configuration
type TaskState struct {
	Value       string   `yaml:"value" json:"value"`
	Name        string   `yaml:"name" json:"name"`
	Description *string  `yaml:"description,omitempty" json:"description,omitempty"`
	Aliases     []string `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	Conceal     *string  `yaml:"conceal,omitempty" json:"conceal,omitempty"`
}

// TasksConfig holds the configuration for task states
type TasksConfig struct {
	States []TaskState `yaml:"states" json:"states"`
}

// Config represents the complete workspace configuration
type Config struct {
	Tasks TasksConfig `yaml:"tasks" json:"tasks"`
}

// Validate checks the configuration for consistency and conflicts
func (c *Config) Validate() error {
	if err := c.Tasks.Validate(); err != nil {
		return fmt.Errorf("tasks configuration error: %w", err)
	}
	return nil
}

// Validate checks the tasks configuration for conflicts and consistency
func (tc *TasksConfig) Validate() error {
	if len(tc.States) == 0 {
		return fmt.Errorf("at least one task state must be defined")
	}

	// Track all values and aliases to check for conflicts
	valueMap := make(map[string]string) // value/alias -> state name for error reporting

	for i, state := range tc.States {
		// Validate required fields
		if state.Value == "" {
			return fmt.Errorf("state %d: value cannot be empty", i)
		}
		if state.Name == "" {
			return fmt.Errorf("state %d: name cannot be empty", i)
		}

		// Check for reserved characters in value
		if strings.Contains(state.Value, "]") {
			return fmt.Errorf("state %q: value cannot contain ']' character", state.Name)
		}

		// Check if value conflicts with existing values or aliases
		if existing, exists := valueMap[state.Value]; exists {
			return fmt.Errorf("state %q: value %q conflicts with state %q", state.Name, state.Value, existing)
		}
		valueMap[state.Value] = state.Name

		// Validate and check aliases
		for j, alias := range state.Aliases {
			if alias == "" {
				return fmt.Errorf("state %q: alias %d cannot be empty", state.Name, j)
			}

			// Check for reserved characters in alias
			if strings.Contains(alias, "]") {
				return fmt.Errorf("state %q: alias %q cannot contain ']' character", state.Name, alias)
			}

			// Check if alias conflicts with the main value of this state
			if alias == state.Value {
				return fmt.Errorf("state %q: alias %q cannot be the same as the main value", state.Name, alias)
			}

			// Check if alias conflicts with existing values or aliases
			if existing, exists := valueMap[alias]; exists {
				return fmt.Errorf("state %q: alias %q conflicts with state %q", state.Name, alias, existing)
			}
			valueMap[alias] = state.Name
		}
	}

	return nil
}

// HasValue checks if the given value matches this task state or any of its aliases
func (ts *TaskState) HasValue(value string) bool {
	if ts.Value == value {
		return true
	}

	// Check aliases
	for _, alias := range ts.Aliases {
		if alias == value {
			return true
		}
	}

	return false
}

// GetConcealText returns the conceal text or the value if conceal is not set
func (ts *TaskState) GetConcealText() string {
	if ts.Conceal != nil {
		return *ts.Conceal
	}
	return "[" + ts.Value + "]"
}
