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
	"regexp"
	"strings"
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
	// Area and Feature support for hierarchical testing
	Area    string
	Feature string
}

// NotedownVHSRunner handles all VHS testing boilerplate for Notedown plugin tests.
type NotedownVHSRunner struct {
	templateDir      string
	lspBuilder       func() (string, error)
	pluginInstaller  func(string) error
	workspaceCreator func(string, string) error
	generateGIF      bool
}

// NewNotedownVHSRunner creates a new runner with default Notedown-specific setup.
func NewNotedownVHSRunner() *NotedownVHSRunner {
	return &NotedownVHSRunner{
		templateDir:      "", // Will be determined by area/feature structure
		lspBuilder:       ensureLSPBinary,
		pluginInstaller:  installPlugin,
		workspaceCreator: createWorkspace,
		generateGIF:      true, // Default to generating GIFs
	}
}

// SetGenerateGIF configures whether GIF files should be generated during testing.
func (r *NotedownVHSRunner) SetGenerateGIF(generate bool) {
	r.generateGIF = generate
}

// RunFeatureTest executes a feature test with area/feature directory structure.
func (r *NotedownVHSRunner) RunFeatureTest(t *testing.T, test VHSTest) {
	r.runTestWithPaths(t, test, true)
}

// RunTest executes a VHS test with all necessary setup and cleanup (legacy support).
func (r *NotedownVHSRunner) RunTest(t *testing.T, test VHSTest) {
	r.runTestWithPaths(t, test, false)
}

// runTestWithPaths is the internal implementation that handles both legacy and feature-based paths.
func (r *NotedownVHSRunner) runTestWithPaths(t *testing.T, test VHSTest, useFeatureStructure bool) {
	// Cleanup old files and VHS processes (skip in CI to avoid termination)
	if !isCI() {
		cleanupTestFiles(test.Name)
		cleanupVHSProcesses()
	}

	// Create temporary directory for test (replace slashes to avoid path issues)
	safeName := strings.ReplaceAll(test.Name, "/", "-")
	tmpDir, err := os.MkdirTemp("", "vhs-test-"+safeName)
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

	// Install plugin with specific LSP binary path
	pluginDir := filepath.Join(tmpDir, "plugin")
	if err := installPluginWithLSP(pluginDir, lspBinary); err != nil {
		t.Fatalf("Failed to install plugin: %v", err)
	}

	// Create workspace - resolve paths based on structure
	var workspaceDir string
	var workspaceSrc string
	if useFeatureStructure {
		// For feature structure: area/feature/workspace
		featureDir := filepath.Join(test.Area, test.Feature)
		workspaceDir = filepath.Join(tmpDir, "workspace")
		workspaceSrc = filepath.Join(featureDir, test.Workspace)
	} else {
		// Legacy: workspaces/workspace-name
		workspaceDir = filepath.Join(tmpDir, "workspace")
		workspaceSrc = test.Workspace
	}

	if err := r.createWorkspaceFromPath(workspaceSrc, workspaceDir); err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Render template - resolve template path based on structure
	var outputFile, gifFile, templatePath string
	if useFeatureStructure {
		featureDir := filepath.Join(test.Area, test.Feature)
		outputFile = filepath.Join(tmpDir, safeName+".ascii") // Use safe name for file paths
		gifFile = filepath.Join(featureDir, "demo.gif")
		templatePath = filepath.Join(featureDir, "demo.tape.tmpl")
	} else {
		outputFile = filepath.Join(tmpDir, test.Name+".ascii")
		gifFile = filepath.Join("gifs", test.Name+".gif")
		templatePath = filepath.Join("testdata/templates", test.Name+".tape.tmpl")
	}

	configFile := filepath.Join(pluginDir, "init.lua")

	// Create consistent template data with safe names (reuse safeName from earlier)
	templateData := map[string]interface{}{
		"OutputFile":   outputFile,
		"WorkspaceDir": workspaceDir,
		"ConfigFile":   configFile,
		"TmpDir":       tmpDir,
		"LSPBinary":    lspBinary,
		"TestName":     safeName, // Use safe name in templates
		"GenerateGIF":  r.generateGIF, // Control GIF generation
	}

	tapeFile, err := r.renderTemplateFromPath(templatePath, test.Name, templateData, tmpDir)
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

	// Assert golden file match - resolve golden file path based on structure
	var goldenFile string
	if useFeatureStructure {
		featureDir := filepath.Join(test.Area, test.Feature)
		goldenFile = filepath.Join(featureDir, "expected.ascii")
	} else {
		goldenFile = filepath.Join("golden", test.Name+".ascii")
	}
	r.assertGoldenMatch(t, goldenFile, result)

	// Copy GIF to appropriate location for visual inspection (only if generation enabled)
	if r.generateGIF {
		srcGif := filepath.Join(tmpDir, safeName+".gif")
		if fileExists(srcGif) {
			if useFeatureStructure {
				// Copy to feature directory
				if err := os.MkdirAll(filepath.Dir(gifFile), 0750); err != nil {
					t.Logf("Failed to create feature directory: %v", err)
					return
				}
			} else {
				// Copy to gifs directory (legacy)
				if err := os.MkdirAll("gifs", 0750); err != nil {
					t.Logf("Failed to create gifs directory: %v", err)
					return
				}
			}
			if err := copyFile(srcGif, gifFile); err != nil {
				t.Logf("Failed to copy GIF file: %v", err)
			}
		}
	}
}

// renderTemplateFromPath renders a VHS template from a specific path.
func (r *NotedownVHSRunner) renderTemplateFromPath(templatePath string, name string, data map[string]interface{}, outputDir string) (string, error) {
	// #nosec G304 - templatePath is controlled within test framework
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New(name).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Create a safe filename for the tape file (avoid subdirectories)
	safeName := strings.ReplaceAll(name, "/", "-")
	tapeFile := filepath.Join(outputDir, safeName+".tape")
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

// createWorkspaceFromPath creates a workspace from a specific path.
func (r *NotedownVHSRunner) createWorkspaceFromPath(workspaceSrc string, outputDir string) error {
	// Get current working directory to find workspaces
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Check if workspaceSrc is relative path or absolute
	var srcWorkspace string
	if filepath.IsAbs(workspaceSrc) {
		srcWorkspace = workspaceSrc
	} else {
		// Try as relative path from current directory
		candidate := filepath.Join(wd, workspaceSrc)
		if _, err := os.Stat(candidate); err == nil {
			srcWorkspace = candidate
		} else {
			return fmt.Errorf("workspace %s not found at %s", workspaceSrc, candidate)
		}
	}

	// Copy workspace to output directory
	// #nosec G204 - subprocess execution with controlled arguments for test framework
	cmd := exec.Command("cp", "-r", srcWorkspace, outputDir)
	return cmd.Run()
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

	// Read the output file (VHS creates it with same base name as tape file)
	baseName := strings.TrimSuffix(filepath.Base(tapeFile), ".tape")
	outputFile := filepath.Join(filepath.Dir(tapeFile), baseName+".ascii")
	// #nosec G304 - outputFile path is controlled within test framework
	return os.ReadFile(outputFile)
}

// normalizeOutput removes dynamic content from VHS output for consistent golden file testing
func normalizeOutput(content string) string {
	// Replace temporary directory paths with a normalized placeholder
	// Handles multiple patterns:
	// - macOS: /var/folders/.../T/vhs-test-...
	// - macOS private: /private/tmp/vhs-test-...
	// - nix-shell: /tmp/nix-shell.../vhs-test-...
	// - nix develop: /tmp/nix-shell.../vhs-test-...
	// - Linux CI: /home/runner/work/_temp/nix-shell.../vhs-test-...
	// - generic tmp: /tmp/.../vhs-test-...
	re := regexp.MustCompile(`/(?:var/folders/[^/]+/[^/]+/T|private/tmp|tmp/[^/]*nix-shell[^/]*|home/runner/work/_temp/[^/]*nix-shell[^/]*|tmp)/vhs-test-[^/]+/`)
	content = re.ReplaceAllString(content, "/tmp/vhs-test-normalized/")

	// Also handle /private/tmp directly since the above might not catch all cases
	privateRe := regexp.MustCompile(`/private(/tmp/vhs-test-normalized/)`)
	content = privateRe.ReplaceAllString(content, "$1")

	// Remove shell prompts and terminal escape sequences that might be inconsistent
	// These can appear differently between test runs
	shellPromptRe := regexp.MustCompile(`\\\[\\\]> \\\[\\\]`)
	content = shellPromptRe.ReplaceAllString(content, "")

	// Normalize terminal escape sequences for colors/formatting
	escapeRe := regexp.MustCompile(`\x1b\[[0-9;]*[mGKH]`)
	content = escapeRe.ReplaceAllString(content, "")

	// Normalize carriage returns and extra whitespace
	crRe := regexp.MustCompile(`\r`)
	content = crRe.ReplaceAllString(content, "")

	// Remove LSP server error messages that might be inconsistent across environments
	lspErrorRe := regexp.MustCompile(`(?m)^.*language server.*failed.*The language server is either not installed.*$\n?`)
	content = lspErrorRe.ReplaceAllString(content, "")

	// Remove any diagnostic print statements from VHS test
	diagnosticRe := regexp.MustCompile(`(?m)^VHS Test: LSP Binary configured as:.*$\n?`)
	content = diagnosticRe.ReplaceAllString(content, "")

	return content
}

// assertGoldenMatch compares actual output with golden file.
func (r *NotedownVHSRunner) assertGoldenMatch(t *testing.T, goldenFile string, actual []byte) {
	normalizedActual := normalizeOutput(string(actual))

	// #nosec G304 - goldenFile path is controlled within test framework
	if expected, err := os.ReadFile(goldenFile); err == nil {
		// Golden file exists - compare for regression testing
		normalizedExpected := normalizeOutput(string(expected))
		if !assert.Equal(t, normalizedExpected, normalizedActual, "Output should match golden file %s", goldenFile) {
			t.Logf("Golden file mismatch! Expected content from: %s", goldenFile)
			t.Logf("To update golden files, delete %s and re-run the test", goldenFile)
		}
	} else {
		// Golden file doesn't exist - this might be a new test or missing file
		t.Logf("Golden file %s does not exist", goldenFile)

		// Auto-create golden file only in local development, fail in CI
		if isCI() {
			t.Fatalf("Golden file %s missing in CI environment - please generate locally first", goldenFile)
		}

		// Create golden file for local development
		if err := os.MkdirAll(filepath.Dir(goldenFile), 0750); err != nil {
			t.Fatalf("Failed to create golden directory: %v", err)
		}
		// Store normalized content in golden file
		normalizedContent := normalizeOutput(string(actual))
		if err := os.WriteFile(goldenFile, []byte(normalizedContent), 0600); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Logf("Created golden file: %s (commit this file to git)", goldenFile)
	}
}

// isCI detects if we're running in a CI environment.
func isCI() bool {
	// Common CI environment variables
	return os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != ""
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

		// Get project root (find the directory containing language-server)
		wd, err := os.Getwd()
		if err != nil {
			buildError = err
			return
		}

		// Go up directories until we find language-server
		projectRoot := wd
		for {
			lsPath := filepath.Join(projectRoot, "language-server")
			if _, err := os.Stat(lsPath); err == nil {
				break // Found language-server directory
			}
			parent := filepath.Dir(projectRoot)
			if parent == projectRoot {
				buildError = fmt.Errorf("could not find project root with language-server directory")
				return
			}
			projectRoot = parent
		}

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
	return installPluginWithLSP(pluginDir, "notedown-language-server")
}

// installPluginWithLSP installs the Notedown Neovim plugin with a specific LSP binary path.
func installPluginWithLSP(pluginDir string, lspBinary string) error {
	// Get project root (find the directory containing neovim)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Go up directories until we find neovim (skip features/neovim, find actual neovim)
	projectRoot := wd
	for {
		nvimPath := filepath.Join(projectRoot, "neovim")
		// Make sure it's not the features/neovim directory by checking for lua/ subdirectory
		luaPath := filepath.Join(nvimPath, "lua")
		if _, err := os.Stat(nvimPath); err == nil {
			if _, err := os.Stat(luaPath); err == nil {
				break // Found neovim directory with lua subdirectory
			}
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return fmt.Errorf("could not find project root with neovim directory containing lua files")
		}
		projectRoot = parent
	}
	neovimSrc := filepath.Join(projectRoot, "neovim")

	// Create the neovim subdirectory in plugin dir
	neovimDest := filepath.Join(pluginDir, "neovim")
	if err := os.MkdirAll(neovimDest, 0750); err != nil {
		return err
	}

	// Copy plugin files to isolated directory using Go's filepath for safety
	err = filepath.Walk(neovimSrc, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(neovimSrc, srcPath)
		if err != nil {
			return err
		}
		destPath := filepath.Join(neovimDest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		return copyFile(srcPath, destPath)
	})
	if err != nil {
		return fmt.Errorf("failed to copy neovim files from %s to %s: %w", neovimSrc, neovimDest, err)
	}

	// Verify the copy worked by checking for config.lua
	configPath := filepath.Join(neovimDest, "lua", "notedown", "config.lua")
	if _, err := os.Stat(configPath); err != nil {
		return fmt.Errorf("config.lua not found at %s after copy: %w", configPath, err)
	}

	// Create minimal init.lua that properly loads the plugin for VHS testing
	initLua := `-- VHS test configuration with plugin loading
vim.opt.runtimepath:prepend("` + pluginDir + `/neovim")

-- Basic Neovim settings
vim.opt.compatible = false
vim.opt.number = true
vim.opt.termguicolors = true
vim.opt.timeout = false
vim.opt.ttimeout = false

-- Force consistent terminal dimensions for status line consistency
vim.opt.columns = 120
vim.opt.lines = 30

-- Netrw settings for deterministic behavior
vim.g.netrw_sort_by = "name"
vim.g.netrw_sort_direction = "normal" 
vim.g.netrw_banner = 1
vim.g.netrw_browse_split = 0
vim.g.netrw_winsize = 25
vim.g.netrw_liststyle = 3
vim.g.netrw_home = "/tmp"  -- Use /tmp for netrw state (avoid persistence)

-- Add Lua package path for the plugin
package.path = package.path .. ";` + pluginDir + `/neovim/lua/?.lua;` + pluginDir + `/neovim/lua/?/init.lua"

-- Set up custom LSP binary path before loading plugin
local notedown_config = require("notedown.config")
notedown_config.defaults.server.cmd = { "` + lspBinary + `", "serve", "--log-level", "debug", "--log-file", "/tmp/notedown.log" }

-- Load the plugin file which contains user commands and calls setup()
pcall(dofile, "` + pluginDir + `/neovim/plugin/notedown.lua")

-- Diagnostic output to verify the configuration
print("VHS Test: LSP Binary configured as: " .. "` + lspBinary + `")
`

	initFile := filepath.Join(pluginDir, "init.lua")
	return os.WriteFile(initFile, []byte(initLua), 0600)
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

// cleanupTestFiles removes old output files but preserves golden files for comparison.
func cleanupTestFiles(testName string) {
	// Create directories if they don't exist
	_ = os.MkdirAll("golden", 0750) // Best effort
	_ = os.MkdirAll("gifs", 0750)   // Best effort

	// Remove old GIF file for this test (golden files should be preserved!)
	gifFile := filepath.Join("gifs", testName+".gif")
	_ = os.Remove(gifFile) // Best effort cleanup
}

// cleanupVHSProcesses kills any existing VHS processes to prevent network connection conflicts.
func cleanupVHSProcesses() {
	// Kill any existing VHS processes that might be holding network connections
	// #nosec G204 - subprocess execution with controlled arguments for test cleanup
	cmd := exec.Command("pkill", "-f", "vhs")
	_ = cmd.Run() // Best effort cleanup, ignore errors

	// Brief delay to ensure processes are fully terminated
	time.Sleep(100 * time.Millisecond)
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
