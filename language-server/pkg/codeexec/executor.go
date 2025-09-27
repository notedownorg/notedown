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

package codeexec

import (
	"fmt"
	"time"

	"github.com/notedownorg/notedown/pkg/parser"
)

// CodeBlock represents a code block with its metadata
type CodeBlock struct {
	Language string
	Content  string
	Range    parser.Range
}

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Success       bool          `json:"success"`
	Stdout        string        `json:"stdout"`
	Stderr        string        `json:"stderr"`
	ExitCode      int           `json:"exitCode"`
	ExecutionTime time.Duration `json:"executionTime"`
	Error         string        `json:"error,omitempty"`
}

// Executor defines the interface for language-specific code execution
type Executor interface {
	// Execute runs the provided code blocks and returns the result
	Execute(blocks []CodeBlock, workspaceRoot string) (*ExecutionResult, error)

	// Language returns the language identifier this executor handles
	Language() string

	// IsAvailable checks if the necessary runtime is available
	IsAvailable() bool
}

// ExecutorFactory creates executors for different languages
type ExecutorFactory struct {
	executors map[string]Executor
}

// NewExecutorFactory creates a new executor factory with default executors
func NewExecutorFactory() *ExecutorFactory {
	factory := &ExecutorFactory{
		executors: make(map[string]Executor),
	}

	// Register the Go executor
	factory.Register(NewGoExecutor())

	return factory
}

// Register registers an executor for a specific language
func (f *ExecutorFactory) Register(executor Executor) {
	f.executors[executor.Language()] = executor
}

// GetExecutor returns an executor for the specified language
func (f *ExecutorFactory) GetExecutor(language string) (Executor, error) {
	executor, exists := f.executors[language]
	if !exists {
		return nil, fmt.Errorf("no executor available for language: %s", language)
	}

	if !executor.IsAvailable() {
		return nil, fmt.Errorf("executor for language %s is not available (missing runtime)", language)
	}

	return executor, nil
}

// GetSupportedLanguages returns a list of supported languages
func (f *ExecutorFactory) GetSupportedLanguages() []string {
	var languages []string
	for lang, executor := range f.executors {
		if executor.IsAvailable() {
			languages = append(languages, lang)
		}
	}
	return languages
}
