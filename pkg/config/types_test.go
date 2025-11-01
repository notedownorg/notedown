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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskState_HasValue(t *testing.T) {
	state := TaskState{
		Value:   "x",
		Aliases: []string{"X", "done", "completed"},
	}

	tests := []struct {
		value    string
		expected bool
	}{
		{"x", true},         // main value
		{"X", true},         // alias
		{"done", true},      // alias
		{"completed", true}, // alias
		{"✓", false},        // not configured
		{"y", false},        // not configured
		{"", false},         // empty
		{"todo", false},     // not configured
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			result := state.HasValue(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTasksConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      TasksConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: TasksConfig{
				States: []TaskState{
					{Value: " ", Name: "todo"},
					{Value: "x", Name: "done"},
				},
			},
			expectError: false,
		},
		{
			name: "empty states",
			config: TasksConfig{
				States: []TaskState{},
			},
			expectError: true,
			errorMsg:    "at least one task state must be defined",
		},
		{
			name: "empty value",
			config: TasksConfig{
				States: []TaskState{
					{Value: "", Name: "todo"},
				},
			},
			expectError: true,
			errorMsg:    "value cannot be empty",
		},
		{
			name: "empty name",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: ""},
				},
			},
			expectError: true,
			errorMsg:    "name cannot be empty",
		},
		{
			name: "value with reserved character",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x]", Name: "invalid"},
				},
			},
			expectError: true,
			errorMsg:    "value cannot contain ']' character",
		},
		{
			name: "duplicate values",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: "done"},
					{Value: "x", Name: "complete"},
				},
			},
			expectError: true,
			errorMsg:    "value \"x\" conflicts",
		},
		{
			name: "valid aliases",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: "done", Aliases: []string{"X", "completed"}},
					{Value: "wip", Name: "work-in-progress", Aliases: []string{"working"}},
				},
			},
			expectError: false,
		},
		{
			name: "alias conflicts with value",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: "done"},
					{Value: "wip", Name: "work-in-progress", Aliases: []string{"x"}},
				},
			},
			expectError: true,
			errorMsg:    "alias \"x\" conflicts",
		},
		{
			name: "alias conflicts with another alias",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: "done", Aliases: []string{"completed"}},
					{Value: "wip", Name: "work-in-progress", Aliases: []string{"completed"}},
				},
			},
			expectError: true,
			errorMsg:    "alias \"completed\" conflicts",
		},
		{
			name: "alias same as main value",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: "done", Aliases: []string{"x"}},
				},
			},
			expectError: true,
			errorMsg:    "alias \"x\" cannot be the same as the main value",
		},
		{
			name: "empty alias",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: "done", Aliases: []string{""}},
				},
			},
			expectError: true,
			errorMsg:    "alias 0 cannot be empty",
		},
		{
			name: "alias with reserved character",
			config: TasksConfig{
				States: []TaskState{
					{Value: "x", Name: "done", Aliases: []string{"x]"}},
				},
			},
			expectError: true,
			errorMsg:    "alias \"x]\" cannot contain ']' character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				Tasks: TasksConfig{
					States: []TaskState{
						{Value: " ", Name: "todo"},
						{Value: "x", Name: "done"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid tasks config",
			config: Config{
				Tasks: TasksConfig{
					States: []TaskState{}, // empty states
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "tasks configuration error")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
