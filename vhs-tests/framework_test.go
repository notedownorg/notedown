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

package vhstests

import (
	"testing"
	"time"

	"github.com/notedownorg/notedown/vhs-tests/pkg/notedown"
)

// vhsTests defines all VHS test cases to run.
var vhsTests = []notedown.VHSTest{
	{Name: "plugin-initialization", Workspace: "plugin-init-test", Timeout: 300 * time.Second},
	{Name: "wikilink-navigation", Workspace: "plugin-init-test", Timeout: 300 * time.Second},
	{Name: "wikilink-completion", Workspace: "completion-test", Timeout: 300 * time.Second},
	{Name: "wikilink-diagnostics", Workspace: "diagnostics-test", Timeout: 300 * time.Second},
}

// TestVHSFramework runs all VHS tests using the clean runner approach.
func TestVHSFramework(t *testing.T) {
	runner := notedown.NewNotedownVHSRunner()

	for _, test := range vhsTests {
		test := test // capture loop variable
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			runner.RunTest(t, test)
		})
	}
}
