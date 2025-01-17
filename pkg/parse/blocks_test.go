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

package parse_test

import (
	"testing"
	"time"

	. "github.com/liamawhite/parse/test"
	"github.com/notedownorg/notedown/pkg/parse"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/test"
)

// Test for edge cases not covered in the compliance testing.
// This could be for various reasons but is typically due to a difference in what
// Notedown is trying to achieve compared to the CommonMark spec.
// e.e. Headings are leaf blocks in commonmark but not in Notedown.
func TestBlocks(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Setext heading closes ATX heading",
			Input:         "# ATX\nSetext\n---\n",
			Parser:        parse.Blocks(time.Now()),
			ExpectedMatch: Bl(Ha1(0, " ATX"), Hs2("Setext", "---")),
			ExpectedOK:    true,
		},
		{
			Name:          "ATX heading closes Setext heading",
			Input:         "Setext\n---\n# ATX\n",
			Parser:        parse.Blocks(time.Now()),
			ExpectedMatch: Bl(Hs1("Setext", "---"), Ha1(0, " ATX")),
			ExpectedOK:    true,
		},
	}
	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
