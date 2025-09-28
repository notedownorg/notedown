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
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GoExecutor implements the Executor interface for Go code
type GoExecutor struct {
	timeout time.Duration
}

// NewGoExecutor creates a new Go executor with default settings
func NewGoExecutor() *GoExecutor {
	return &GoExecutor{
		timeout: 30 * time.Second, // Default 30 second timeout as specified
	}
}

// Execute runs Go code blocks by merging them into a single main.go file
func (g *GoExecutor) Execute(blocks []CodeBlock, workspaceRoot string) (*ExecutionResult, error) {
	if len(blocks) == 0 {
		return &ExecutionResult{
			Success:  false,
			ExitCode: 1,
			Error:    "no Go code blocks found",
		}, nil
	}

	// Generate the merged Go program
	program, err := g.generateProgram(blocks)
	if err != nil {
		return &ExecutionResult{
			Success:  false,
			ExitCode: 1,
			Error:    fmt.Sprintf("failed to generate Go program: %v", err),
		}, nil
	}

	// Create temporary directory for execution
	tempDir, err := os.MkdirTemp("", "notedown-go-exec-*")
	if err != nil {
		return &ExecutionResult{
			Success:  false,
			ExitCode: 1,
			Error:    fmt.Sprintf("failed to create temp directory: %v", err),
		}, nil
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Write the program to main.go
	mainFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(program), 0600); err != nil {
		return &ExecutionResult{
			Success:  false,
			ExitCode: 1,
			Error:    fmt.Sprintf("failed to write Go program: %v", err),
		}, nil
	}

	// Execute the program
	return g.runProgram(mainFile, workspaceRoot)
}

// generateProgram merges code blocks into a single Go program
func (g *GoExecutor) generateProgram(blocks []CodeBlock) (string, error) {
	var program strings.Builder

	// Add package declaration
	program.WriteString("package main\n\n")

	// Concatenate all code blocks in document order
	for _, block := range blocks {
		program.WriteString(block.Content)
		program.WriteString("\n\n")
	}

	return program.String(), nil
}

// runProgram executes the Go program and captures output
func (g *GoExecutor) runProgram(mainFile, workspaceRoot string) (*ExecutionResult, error) {
	startTime := time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()

	// Run 'go run main.go' with workspace as working directory
	cmd := exec.CommandContext(ctx, "go", "run", mainFile)
	cmd.Dir = workspaceRoot

	// Capture stdout and stderr
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	executionTime := time.Since(startTime)

	result := &ExecutionResult{
		Stdout:        stdout.String(),
		Stderr:        stderr.String(),
		ExecutionTime: executionTime,
	}

	if err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			result.Success = false
			result.ExitCode = 1
			result.Error = fmt.Sprintf("execution timed out after %v", g.timeout)
		} else if exitError, ok := err.(*exec.ExitError); ok {
			// Command ran but exited with non-zero status
			result.Success = false
			result.ExitCode = exitError.ExitCode()
		} else {
			// Command failed to run
			result.Success = false
			result.ExitCode = 1
			result.Error = fmt.Sprintf("failed to execute Go program: %v", err)
		}
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	return result, nil
}

// Language returns the language identifier
func (g *GoExecutor) Language() string {
	return "go"
}

// IsAvailable checks if the Go runtime is available
func (g *GoExecutor) IsAvailable() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

// SetTimeout sets the execution timeout
func (g *GoExecutor) SetTimeout(timeout time.Duration) {
	g.timeout = timeout
}
