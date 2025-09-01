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

// Package notedown provides a containerized VHS testing runner using testcontainers-go.
package notedown

import (
	"context"
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

	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
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

// ContainerVHSRunner handles VHS testing using Docker containers for better isolation and parallelism.
// Assumes that the notedown-vhs:latest image has already been built.
type ContainerVHSRunner struct {
	generateGIF bool
	imageName   string
}

// NewContainerVHSRunner creates a new containerized VHS runner.
func NewContainerVHSRunner() *ContainerVHSRunner {
	return &ContainerVHSRunner{
		generateGIF: true,
		imageName:   "notedown-vhs:latest",
	}
}

// SetGenerateGIF configures whether GIF files should be generated during testing.
func (r *ContainerVHSRunner) SetGenerateGIF(generate bool) {
	r.generateGIF = generate
}

// RunFeatureTest executes a feature test using Docker containers with testcontainers-go.
func (r *ContainerVHSRunner) RunFeatureTest(t *testing.T, test VHSTest) {
	// Create temporary directory for test data
	safeName := strings.ReplaceAll(test.Name, "/", "-")
	tmpDir, err := os.MkdirTemp("", "notedown-container-"+safeName)
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to cleanup temp dir: %v", err)
		}
	}()

	// Set up workspace and template
	workspaceDir := filepath.Join(tmpDir, "workspace")
	featureDir := filepath.Join(test.Area, test.Feature)
	workspaceSrc := filepath.Join(featureDir, test.Workspace)

	if err := r.createWorkspaceFromPath(workspaceSrc, workspaceDir); err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// No need to create plugin config - using built-in docker-init.lua

	// Render VHS template
	templatePath := filepath.Join(featureDir, "demo.tape.tmpl")

	templateData := map[string]interface{}{
		"OutputFile":   "/vhs/" + safeName + ".ascii",             // Container path
		"WorkspaceDir": "/vhs/workspace",                          // Container path
		"TmpDir":       "/vhs",                                    // Container path
		"LSPBinary":    "/usr/local/bin/notedown-language-server", // Container built-in binary
		"TestName":     safeName,
		"GenerateGIF":  r.generateGIF,
	}

	tapeFile, err := r.renderTemplateFromPath(templatePath, test.Name, templateData, tmpDir)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Run VHS in container with neovim available
	result, err := r.runVHSContainer(context.Background(), t, tmpDir, filepath.Base(tapeFile), test.Timeout)
	if err != nil {
		t.Fatalf("Container VHS execution failed: %v", err)
	}

	// Check for golden file and assert match
	goldenFile := filepath.Join(featureDir, "expected.ascii")
	r.assertGoldenMatch(t, goldenFile, result)

	// Copy GIF if generation enabled
	if r.generateGIF {
		srcGif := filepath.Join(tmpDir, safeName+".gif")
		if fileExists(srcGif) {
			// Check if GIF file has content (non-zero size)
			if stat, err := os.Stat(srcGif); err == nil {
				if stat.Size() == 0 {
					// List all files in tmp directory for debugging
					if files, err := os.ReadDir(tmpDir); err == nil {
						t.Logf("Files in temp directory %s:", tmpDir)
						for _, file := range files {
							if info, err := file.Info(); err == nil {
								t.Logf("  %s (%d bytes)", file.Name(), info.Size())
							}
						}
					}
					t.Fatalf("Generated GIF file %s is empty (0 bytes) - this indicates VHS generation failed", srcGif)
				}
			} else {
				t.Fatalf("Failed to stat GIF file %s: %v", srcGif, err)
			}

			gifFile := filepath.Join(featureDir, "demo.gif")
			if err := os.MkdirAll(filepath.Dir(gifFile), 0750); err != nil {
				t.Fatalf("Failed to create feature directory: %v", err)
			}
			if err := copyFile(srcGif, gifFile); err != nil {
				t.Fatalf("Failed to copy GIF file: %v", err)
			}
		} else {
			// List all files in tmp directory for debugging
			if files, err := os.ReadDir(tmpDir); err == nil {
				t.Logf("Files in temp directory %s:", tmpDir)
				for _, file := range files {
					if info, err := file.Info(); err == nil {
						t.Logf("  %s (%d bytes)", file.Name(), info.Size())
					}
				}
			}
			t.Fatalf("Expected GIF file %s was not generated by VHS", srcGif)
		}
	}
}

// runVHSContainer executes VHS using the comprehensive Notedown VHS image.
func (r *ContainerVHSRunner) runVHSContainer(ctx context.Context, t *testing.T, mountDir, tapeFile string, timeout time.Duration) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Get absolute path for bind mount
	absPath, err := filepath.Abs(mountDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image: r.imageName, // Use pre-built image
		Cmd:   []string{tapeFile},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			// Mount the test directory only
			hostConfig.Binds = append(hostConfig.Binds, absPath+":/vhs")
		},
		// Don't wait for exit during startup - we'll poll manually
	}

	// Start container
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}
	defer func() {
		if err := c.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	// Wait for container to finish by checking its state periodically
	for {
		state, err := c.State(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get container state: %w", err)
		}

		if state.Status == "exited" {
			// Always capture logs for debugging
			logs, logErr := c.Logs(ctx)
			var logOutput string
			if logErr == nil {
				defer func() { _ = logs.Close() }()
				logBytes, _ := io.ReadAll(logs)
				logOutput = string(logBytes)
			}

			if state.ExitCode != 0 {
				return nil, fmt.Errorf("container exited with non-zero code: %d\nContainer logs:\n%s", state.ExitCode, logOutput)
			}

			// Store logs for potential debugging (even on success)
			if len(logOutput) > 0 {
				t.Logf("VHS Container logs:\n%s", logOutput)
			}
			break
		}

		// Check timeout
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled")
		case <-time.After(time.Second):
			// Continue polling
		}
	}

	// Read the output file
	outputFile := filepath.Join(mountDir, strings.ReplaceAll(tapeFile, ".tape", ".ascii"))
	result, err := os.ReadFile(outputFile) // #nosec G304 - controlled test file paths
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	return result, nil
}

// findProjectRoot finds the project root directory (where go.mod is located)

// Helper methods (reusing from original runner)

func (r *ContainerVHSRunner) createWorkspaceFromPath(srcWorkspace, outputDir string) error {
	// #nosec G204 - subprocess execution with controlled arguments for test framework
	cmd := exec.Command("cp", "-r", srcWorkspace, outputDir)
	return cmd.Run()
}

// Build coordination - ensures VHS image is built once per test session
var (
	imageBuildOnce     sync.Once
	imageBuildErr      error
	imageBuildComplete = make(chan struct{})
)

// EnsureVHSImage builds the notedown-vhs Docker image once per test session.
// Subsequent calls will wait for the initial build to complete.
func EnsureVHSImage(t *testing.T) error {
	imageBuildOnce.Do(func() {
		defer close(imageBuildComplete)
		imageBuildErr = buildVHSImage(t)
	})

	// Wait for build to complete if it's in progress
	<-imageBuildComplete
	return imageBuildErr
}

// buildVHSImage performs the actual Docker image build
func buildVHSImage(t *testing.T) error {
	t.Logf("Building notedown-vhs Docker image...")

	// Find project root (where go.mod is located)
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	// Build the Docker image with Dockerfile in features/neovim directory
	// #nosec G204 - subprocess execution with controlled arguments for test framework
	cmd := exec.Command("docker", "build", "-t", "notedown-vhs:latest", "-f", "features/neovim/Dockerfile", projectRoot)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Docker build failed: %s", string(output))
		return fmt.Errorf("failed to build Docker image: %w", err)
	}

	t.Logf("Successfully built notedown-vhs:latest")
	return nil
}

// findProjectRoot finds the project root directory (where go.mod is located)
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root (go.mod not found)")
}

func (r *ContainerVHSRunner) renderTemplateFromPath(templatePath, name string, data map[string]interface{}, tmpDir string) (string, error) {
	// Read template file
	templateContent, err := os.ReadFile(templatePath) // #nosec G304 - controlled template file paths
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Parse template
	tmpl, err := template.New(name).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	safeName := strings.ReplaceAll(name, "/", "-")
	outputPath := filepath.Join(tmpDir, safeName+".tape")
	outputFile, err := os.Create(outputPath) // #nosec G304 - controlled output file paths
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() { _ = outputFile.Close() }()

	// Execute template
	if err := tmpl.Execute(outputFile, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return outputPath, nil
}

func (r *ContainerVHSRunner) assertGoldenMatch(t *testing.T, goldenFile string, actual []byte) {
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
			t.Fatalf("Failed to create golden file directory: %v", err)
		}
		if err := os.WriteFile(goldenFile, []byte(normalizedActual), 0600); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Logf("Created golden file: %s (commit this file to git)", goldenFile)
	}
}

// Helper functions from the original runner

// normalizeOutput removes dynamic content from VHS output for consistent golden file testing
func normalizeOutput(content string) string {
	// Replace temporary directory paths with normalized placeholder
	re := regexp.MustCompile(`/(?:var/folders/[^/]+/[^/]+/T|private/tmp|tmp/[^/]*nix-shell[^/]*|home/runner/work/_temp/[^/]*nix-shell[^/]*|tmp)/notedown[^/]*/`)
	content = re.ReplaceAllString(content, "/tmp/notedown-normalized/")

	// Handle /private/tmp directly
	privateRe := regexp.MustCompile(`/private(/tmp/notedown-normalized/)`)
	content = privateRe.ReplaceAllString(content, "$1")

	// Remove shell prompts and terminal escape sequences
	shellPromptRe := regexp.MustCompile(`\\\[\\\]> \\\[\\\]`)
	content = shellPromptRe.ReplaceAllString(content, "")

	standalonePromptRe := regexp.MustCompile(`(?m)^>\n`)
	content = standalonePromptRe.ReplaceAllString(content, "\n")

	// Normalize terminal escape sequences
	escapeRe := regexp.MustCompile(`\x1b\[[0-9;]*[mGKH]`)
	content = escapeRe.ReplaceAllString(content, "")

	crRe := regexp.MustCompile(`\r`)
	content = crRe.ReplaceAllString(content, "")

	// Fix line wrapping differences
	wrapRe := regexp.MustCompile(`([a-zA-Z0-9/_-]+)\n([a-zA-Z0-9])`)
	content = wrapRe.ReplaceAllString(content, "$1$2")

	// Normalize blank lines
	content = strings.TrimLeft(content, "\n")
	content = strings.TrimRight(content, "\n") + "\n"

	// Remove starting separators
	startingSeparatorRe := regexp.MustCompile(`^[\s\n]*────────────────────────────────────────────────────────────────────────────────\n`)
	content = startingSeparatorRe.ReplaceAllString(content, "")

	// Remove empty screens
	emptyScreenRe := regexp.MustCompile(`────────────────────────────────────────────────────────────────────────────────\n\n+────────────────────────────────────────────────────────────────────────────────\n`)
	content = emptyScreenRe.ReplaceAllString(content, "")

	// Remove blank lines before separators
	blankBeforeSeparatorRe := regexp.MustCompile(`\n+────────────────────────────────────────────────────────────────────────────────\n`)
	content = blankBeforeSeparatorRe.ReplaceAllString(content, "\n────────────────────────────────────────────────────────────────────────────────\n")

	// Normalize multiple blank lines
	multiBlankRe := regexp.MustCompile(`\n{3,}`)
	content = multiBlankRe.ReplaceAllString(content, "\n\n")

	// Normalize workspace status output
	workspaceStatusRe := regexp.MustCompile(`(Matched Workspace: [^\n]+)\n+\s*(Detection Method:)`)
	content = workspaceStatusRe.ReplaceAllString(content, "$1\n  $2")

	// Remove LSP error messages
	lspErrorRe := regexp.MustCompile(`(?m)^.*language server.*failed.*The language server is either not installed.*$\n?`)
	content = lspErrorRe.ReplaceAllString(content, "")

	// Final normalization
	content = strings.Trim(content, "\n") + "\n"

	return content
}

// isCI detects if we're running in a CI environment
func isCI() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != ""
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// #nosec G304 - src path is controlled within test framework
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = srcFile.Close() }()

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
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
