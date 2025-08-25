// Copyright 2024 Notedown Authors
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

// Package notedown provides a clean VHS testing runner for Notedown plugin tests.
package notedown

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"text/template"
	"time"

	"github.com/stretchr/testify/assert"
)

// VHSTest defines a single VHS test case.
type VHSTest struct {
	Name      string
	Workspace string
	Timeout   time.Duration
}

// NotedownVHSRunner handles all VHS testing boilerplate for Notedown plugin tests.
type NotedownVHSRunner struct {
	templateDir      string
	lspBuilder       func() (string, error)
	pluginInstaller  func(string) error
	workspaceCreator func(string, string) error
}

// NewNotedownVHSRunner creates a new runner with default Notedown-specific setup.
func NewNotedownVHSRunner() *NotedownVHSRunner {
	return &NotedownVHSRunner{
		templateDir:      "testdata/templates",
		lspBuilder:       ensureLSPBinary,
		pluginInstaller:  installPlugin,
		workspaceCreator: createWorkspace,
	}
}

// RunTest executes a VHS test with all necessary setup and cleanup.
func (r *NotedownVHSRunner) RunTest(t *testing.T, test VHSTest) {
	// Cleanup old files
	cleanupTestFiles(test.Name)

	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "vhs-test-"+test.Name)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to cleanup temp dir: %v", err)
		}
	}()

	// Build LSP server
	lspBinary, err := r.lspBuilder()
	if err != nil {
		t.Fatalf("Failed to build LSP server: %v", err)
	}

	// Install plugin
	pluginDir := filepath.Join(tmpDir, "plugin")
	if err := r.pluginInstaller(pluginDir); err != nil {
		t.Fatalf("Failed to install plugin: %v", err)
	}

	// Create workspace
	workspaceDir := filepath.Join(tmpDir, "workspace")
	if err := r.workspaceCreator(test.Workspace, workspaceDir); err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Render template
	outputFile := filepath.Join(tmpDir, test.Name+".ascii")
	gifFile := filepath.Join("gifs", test.Name+".gif")
	configFile := filepath.Join(pluginDir, "init.lua")

	templateData := map[string]interface{}{
		"OutputFile":   outputFile,
		"WorkspaceDir": workspaceDir,
		"ConfigFile":   configFile,
		"TmpDir":       tmpDir,
		"LSPBinary":    lspBinary,
	}

	tapeFile, err := r.renderTemplate(test.Name, templateData, tmpDir)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Execute VHS
	timeout := test.Timeout
	if timeout == 0 {
		timeout = 300 * time.Second
	}

	result, err := r.executeVHS(tapeFile, timeout)
	if err != nil {
		t.Fatalf("VHS execution failed: %v", err)
	}

	// Assert golden file match
	goldenFile := filepath.Join("golden", test.Name+".ascii")
	r.assertGoldenMatch(t, goldenFile, result)

	// Copy GIF to gifs directory for visual inspection
	if srcGif := filepath.Join(tmpDir, test.Name+".gif"); fileExists(srcGif) {
		if err := os.MkdirAll("gifs", 0750); err != nil {
			t.Logf("Failed to create gifs directory: %v", err)
			return
		}
		if err := copyFile(srcGif, gifFile); err != nil {
			t.Logf("Failed to copy GIF file: %v", err)
		}
	}
}

// renderTemplate renders a VHS template with the given data.
func (r *NotedownVHSRunner) renderTemplate(name string, data map[string]interface{}, outputDir string) (string, error) {
	templateFile := filepath.Join(r.templateDir, name+".tape.tmpl")
	// #nosec G304 - templateFile path is controlled within test framework
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templateFile, err)
	}

	tmpl, err := template.New(name).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	tapeFile := filepath.Join(outputDir, name+".tape")
	// #nosec G304 - tapeFile path is controlled within test framework
	f, err := os.Create(tapeFile)
	if err != nil {
		return "", fmt.Errorf("failed to create tape file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			return // template execution error will be returned instead
		}
	}()

	if err := tmpl.Execute(f, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return tapeFile, nil
}

// executeVHS runs VHS on the given tape file.
func (r *NotedownVHSRunner) executeVHS(tapeFile string, timeout time.Duration) ([]byte, error) {
	cmd := exec.Command("vhs", tapeFile)
	cmd.Dir = filepath.Dir(tapeFile)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("vhs failed: %w\nStderr: %s", err, stderr.String())
		}
	case <-time.After(timeout):
		if err := cmd.Process.Kill(); err != nil {
			// Process may have already exited, ignore error
			_ = err
		}
		return nil, fmt.Errorf("vhs timed out after %v", timeout)
	}

	// Read the output file
	outputFile := filepath.Join(filepath.Dir(tapeFile), filepath.Base(tapeFile[:len(tapeFile)-5])+".ascii")
	// #nosec G304 - outputFile path is controlled within test framework
	return os.ReadFile(outputFile)
}

// assertGoldenMatch compares actual output with golden file.
func (r *NotedownVHSRunner) assertGoldenMatch(t *testing.T, goldenFile string, actual []byte) {
	// #nosec G304 - goldenFile path is controlled within test framework
	if expected, err := os.ReadFile(goldenFile); err == nil {
		assert.Equal(t, string(expected), string(actual), "Output should match golden file %s", goldenFile)
	} else {
		// Create golden file if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(goldenFile), 0750); err != nil {
			t.Logf("Failed to create golden directory: %v", err)
			return
		}
		if err := os.WriteFile(goldenFile, actual, 0600); err != nil {
			t.Logf("Failed to write golden file: %v", err)
		}
		t.Logf("Created golden file: %s", goldenFile)
	}
}

// Shared LSP binary building to avoid redundant builds in parallel tests
var (
	sharedLSPBinary string
	buildOnce       sync.Once
	buildError      error
)

// ensureLSPBinary builds the LSP server once and returns the path to the shared binary.
func ensureLSPBinary() (string, error) {
	buildOnce.Do(func() {
		// Create temporary location for shared binary
		tmpDir, err := os.MkdirTemp("/tmp", "vhs-shared")
		if err != nil {
			buildError = fmt.Errorf("failed to create shared temp dir: %w", err)
			return
		}
		sharedLSPBinary = filepath.Join(tmpDir, "notedown-language-server")

		// Get project root (vhs directory)
		wd, err := os.Getwd()
		if err != nil {
			buildError = err
			return
		}
		projectRoot := filepath.Dir(wd) // Go from vhs-tests to vhs

		// Build LSP server once
		// #nosec G204 - subprocess execution with controlled arguments for test framework
		cmd := exec.Command("go", "build",
			"-ldflags", "-w -s -X github.com/notedownorg/notedown/pkg/version.version=test",
			"-o", sharedLSPBinary,
			"./language-server/")
		cmd.Dir = projectRoot

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			buildError = fmt.Errorf("shared LSP build failed: %w\nStderr: %s", err, stderr.String())
			return
		}

		// Make sure binary is executable
		// #nosec G302 - binary needs to be executable for testing
		buildError = os.Chmod(sharedLSPBinary, 0700)
	})
	return sharedLSPBinary, buildError
}

// installPlugin installs the Notedown Neovim plugin.
func installPlugin(pluginDir string) error {
	// Get project root (vhs directory)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	projectRoot := filepath.Dir(wd) // Go from vhs-tests to vhs
	neovimSrc := filepath.Join(projectRoot, "neovim")

	// Copy plugin files to isolated directory
	// #nosec G204 - subprocess execution with controlled arguments for test framework
	cmd := exec.Command("cp", "-r", neovimSrc, pluginDir)
	return cmd.Run()
}

// createWorkspace creates a workspace in the specified directory.
func createWorkspace(workspaceName string, outputDir string) error {
	// Get current working directory to find workspaces
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Find workspaces directory (could be in current dir or parent vhs dir)
	var srcWorkspace string

	// Try current directory first
	candidate := filepath.Join(wd, "workspaces", workspaceName)
	if _, err := os.Stat(candidate); err == nil {
		srcWorkspace = candidate
	} else {
		// Try parent directory (vhs/workspaces)
		parent := filepath.Dir(wd) // Go from vhs-tests to vhs
		candidate = filepath.Join(parent, "workspaces", workspaceName)
		if _, err := os.Stat(candidate); err == nil {
			srcWorkspace = candidate
		} else {
			return fmt.Errorf("workspace %s not found in workspaces directory", workspaceName)
		}
	}

	// Copy workspace to output directory
	// #nosec G204 - subprocess execution with controlled arguments for test framework
	cmd := exec.Command("cp", "-r", srcWorkspace, outputDir)
	return cmd.Run()
}

// cleanupTestFiles removes old ASCII and GIF files for a specific test.
func cleanupTestFiles(testName string) {
	// Create directories if they don't exist
	_ = os.MkdirAll("golden", 0750) // Best effort
	_ = os.MkdirAll("gifs", 0750)   // Best effort

	// Remove old golden file for this test
	goldenFile := filepath.Join("golden", testName+".ascii")
	_ = os.Remove(goldenFile) // Best effort cleanup

	// Remove old GIF file for this test
	gifFile := filepath.Join("gifs", testName+".gif")
	_ = os.Remove(gifFile) // Best effort cleanup
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	// #nosec G304 - src path is controlled within test framework
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close() // Defer close, error handled by main function
	}()

	// Create destination directory if needed
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0750); err != nil {
		return err
	}

	// #nosec G304 - dst path is controlled within test framework
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		if err := dstFile.Close(); err != nil && err.Error() != "file already closed" {
			// Log error but don't override main function error
			_ = err
		}
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
