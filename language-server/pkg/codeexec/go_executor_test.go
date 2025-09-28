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
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/notedownorg/notedown/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoExecutor_Language(t *testing.T) {
	executor := NewGoExecutor()
	assert.Equal(t, "go", executor.Language())
}

func TestGoExecutor_IsAvailable(t *testing.T) {
	executor := NewGoExecutor()
	// This test assumes Go is available on the test system
	// If Go is not available, the test will be skipped
	if !executor.IsAvailable() {
		t.Skip("Go is not available on this system")
	}
	assert.True(t, executor.IsAvailable())
}

func TestGoExecutor_Execute_SimpleProgram(t *testing.T) {
	executor := NewGoExecutor()
	if !executor.IsAvailable() {
		t.Skip("Go is not available on this system")
	}

	// Create a temporary workspace
	tempDir, err := os.MkdirTemp("", "notedown-test-workspace-*")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Test simple Hello World program
	blocks := []CodeBlock{
		{
			Language: "go",
			Content: `import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
			Range: parser.Range{},
		},
	}

	result, err := executor.Execute(blocks, tempDir)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Stdout, "Hello, World!")
	assert.Empty(t, result.Stderr)
	assert.Empty(t, result.Error)
	assert.Greater(t, result.ExecutionTime, time.Duration(0))
}

func TestGoExecutor_Execute_MultipleBlocks(t *testing.T) {
	executor := NewGoExecutor()
	if !executor.IsAvailable() {
		t.Skip("Go is not available on this system")
	}

	// Create a temporary workspace
	tempDir, err := os.MkdirTemp("", "notedown-test-workspace-*")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Test multiple code blocks that need to be merged
	blocks := []CodeBlock{
		{
			Language: "go",
			Content: `import (
	"fmt"
	"strings"
)

func greet(name string) string {
	return "Hello, " + strings.Title(name) + "!"
}`,
			Range: parser.Range{},
		},
		{
			Language: "go",
			Content: `func main() {
	message := greet("world")
	fmt.Println(message)
}`,
			Range: parser.Range{},
		},
	}

	result, err := executor.Execute(blocks, tempDir)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Stdout, "Hello, World!")
	assert.Empty(t, result.Stderr)
	assert.Empty(t, result.Error)
}

func TestGoExecutor_Execute_CompilationError(t *testing.T) {
	executor := NewGoExecutor()
	if !executor.IsAvailable() {
		t.Skip("Go is not available on this system")
	}

	// Create a temporary workspace
	tempDir, err := os.MkdirTemp("", "notedown-test-workspace-*")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Test code with compilation error
	blocks := []CodeBlock{
		{
			Language: "go",
			Content: `import "fmt"

func main() {
	fmt.Println("Hello, World!"  // Missing closing parenthesis
}`,
			Range: parser.Range{},
		},
	}

	result, err := executor.Execute(blocks, tempDir)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Success)
	assert.NotEqual(t, 0, result.ExitCode)
	assert.Empty(t, result.Stdout)
	assert.NotEmpty(t, result.Stderr)
	assert.Contains(t, result.Stderr, "syntax error")
}

func TestGoExecutor_Execute_RuntimeError(t *testing.T) {
	executor := NewGoExecutor()
	if !executor.IsAvailable() {
		t.Skip("Go is not available on this system")
	}

	// Create a temporary workspace
	tempDir, err := os.MkdirTemp("", "notedown-test-workspace-*")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Test code with runtime error (panic)
	blocks := []CodeBlock{
		{
			Language: "go",
			Content: `import "fmt"

func main() {
	fmt.Println("Before error")
	panic("test error")
	fmt.Println("After error") // This won't be reached
}`,
			Range: parser.Range{},
		},
	}

	result, err := executor.Execute(blocks, tempDir)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Success)
	assert.NotEqual(t, 0, result.ExitCode)
	assert.Contains(t, result.Stdout, "Before error")
	assert.NotContains(t, result.Stdout, "After error")
	// Panic output goes to stderr
	assert.Contains(t, result.Stderr, "panic")
}

func TestGoExecutor_Execute_NoBlocks(t *testing.T) {
	executor := NewGoExecutor()
	if !executor.IsAvailable() {
		t.Skip("Go is not available on this system")
	}

	// Create a temporary workspace
	tempDir, err := os.MkdirTemp("", "notedown-test-workspace-*")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Test with no code blocks
	var blocks []CodeBlock

	result, err := executor.Execute(blocks, tempDir)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.Success)
	assert.Equal(t, 1, result.ExitCode)
	assert.Contains(t, result.Error, "no Go code blocks found")
}

func TestGoExecutor_Execute_FileSystemAccess(t *testing.T) {
	executor := NewGoExecutor()
	if !executor.IsAvailable() {
		t.Skip("Go is not available on this system")
	}

	// Create a temporary workspace with a test file
	tempDir, err := os.MkdirTemp("", "notedown-test-workspace-*")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("Hello from file!"), 0600)
	require.NoError(t, err)

	// Test code that reads the test file
	blocks := []CodeBlock{
		{
			Language: "go",
			Content: `import (
	"fmt"
	"os"
)

func main() {
	content, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("File content: %s\n", string(content))
}`,
			Range: parser.Range{},
		},
	}

	result, err := executor.Execute(blocks, tempDir)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, result.Success)
	assert.Equal(t, 0, result.ExitCode)
	assert.Contains(t, result.Stdout, "File content: Hello from file!")
}

func TestGoExecutor_GenerateProgram(t *testing.T) {
	executor := NewGoExecutor()

	blocks := []CodeBlock{
		{
			Language: "go",
			Content:  `import "fmt"`,
		},
		{
			Language: "go",
			Content:  `func main() { fmt.Println("test") }`,
		},
	}

	program, err := executor.generateProgram(blocks)
	require.NoError(t, err)

	expected := "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"test\") }\n\n"
	assert.Equal(t, expected, program)
}

func TestGoExecutor_SetTimeout(t *testing.T) {
	executor := NewGoExecutor()

	// Test default timeout
	assert.Equal(t, 30*time.Second, executor.timeout)

	// Test setting custom timeout
	customTimeout := 5 * time.Second
	executor.SetTimeout(customTimeout)
	assert.Equal(t, customTimeout, executor.timeout)
}
