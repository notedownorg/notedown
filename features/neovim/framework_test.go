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
}

// TestFeatures runs all feature tests using the area/feature structure.
func TestFeatures(t *testing.T) {
	t.Logf("=== Starting feature test suite with %d tests (GIF generation: %t) ===", len(featureTests), *generateGIF)
	runner := notedown.NewNotedownVHSRunner()
	runner.SetGenerateGIF(*generateGIF)

	for _, test := range featureTests {
		test := test // capture loop variable
		testName := test.Area + "/" + test.Feature
		t.Run(testName, func(t *testing.T) {
			t.Parallel()
			// Convert FeatureTest to VHSTest for runner
			vhsTest := notedown.VHSTest{
				Name:      testName,
				Workspace: test.Workspace,
				Timeout:   test.Timeout,
				Area:      test.Area,
				Feature:   test.Feature,
			}
			runner.RunFeatureTest(t, vhsTest)
		})
	}
	t.Logf("=== Feature test suite completed ===")
}
