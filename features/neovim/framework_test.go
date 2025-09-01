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

package features

import (
	"flag"
	"testing"
	"time"

	"github.com/notedownorg/notedown/features/neovim/pkg/notedown"
)

// Command line flags
var (
	generateGIF = flag.Bool("gif", true, "Generate GIF files during testing (default: true)")
)

// FeatureTest defines a test for a specific feature within an area.
type FeatureTest struct {
	Area      string
	Feature   string
	Workspace string
	Timeout   time.Duration
}

// featureTests defines all feature tests to run.
var featureTests = []FeatureTest{
	{Area: "initialization", Feature: "workspace-status-command", Workspace: "workspace", Timeout: 300 * time.Second},
	{Area: "wikilinks", Feature: "completion", Workspace: "workspace", Timeout: 300 * time.Second},
	{Area: "wikilinks", Feature: "syntax", Workspace: "workspace", Timeout: 300 * time.Second},
	{Area: "wikilinks", Feature: "diagnostics-and-code-actions", Workspace: "workspace", Timeout: 300 * time.Second},
}

// TestFeatures runs all feature tests using build-once + optimized parallel execution.
func TestFeatures(t *testing.T) {
	t.Logf("=== Starting containerized feature test suite with %d tests (GIF: %t) ===",
		len(featureTests), *generateGIF)

	// Phase 1: Build Docker image once before running any tests
	t.Logf("Phase 1: Building VHS Docker image...")
	if err := notedown.EnsureVHSImage(t); err != nil {
		t.Fatalf("Failed to build VHS Docker image: %v", err)
	}

	// Phase 2: Run tests with Go's native parallelism (t.Parallel)
	t.Logf("Phase 2: Running tests with parallel execution...")
	runTestsWithNativeParallelism(t)

	t.Logf("=== Feature test suite completed ===")
}

// runTestsWithNativeParallelism runs tests using Go's built-in t.Parallel() after image is built
func runTestsWithNativeParallelism(t *testing.T) {
	for _, test := range featureTests {
		test := test // capture loop variable
		testName := test.Area + "/" + test.Feature
		t.Run(testName, func(t *testing.T) {
			t.Parallel() // Image is already built, safe to run in parallel

			// Create runner instance for this test
			runner := notedown.NewContainerVHSRunner()
			runner.SetGenerateGIF(*generateGIF)

			// Convert FeatureTest to VHSTest for runner
			vhsTest := notedown.VHSTest{
				Name:      testName,
				Workspace: test.Workspace,
				Timeout:   test.Timeout,
				Area:      test.Area,
				Feature:   test.Feature,
			}

			// Run the test
			runner.RunFeatureTest(t, vhsTest)
		})
	}
}
