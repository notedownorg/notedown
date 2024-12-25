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

package commonmark

import (
	"testing"

	. "github.com/liamawhite/parse/test"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/test"
)

func TestCommonmarkComplianceBlocksParagraphs(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 219",
			Input:         spec(219),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("aaa"), Bln, P("bbb")),
		},
		{
			Name:          "Example 220",
			Input:         spec(220),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("aaa\nbbb"), Bln, P("ccc\nddd")),
		},
		{
			Name:          "Example 221",
			Input:         spec(221),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("aaa"), Bln, Bln, P("bbb")),
		},
		{
			Name:          "Example 222",
			Input:         spec(222),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("  aaa\n bbb")),
		},
		{
			Name:          "Example 223",
			Input:         spec(223),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("aaa\n             bbb\n                                       ccc")),
		},
		{
			Name:          "Example 224",
			Input:         spec(224),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("   aaa\nbbb")),
		},
		{
			Name:          "Example 225",
			Input:         spec(225),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    aaa"), P("bbb")),
		},
		{
			Name:          "Example 226",
			Input:         spec(226),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("aaa     \nbbb     ")),
		},

		// 4.9 Blank lines
		{
			Name:          "Example 227",
			Input:         spec(227),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("  "), Bln, P("aaa\n  "), Bln, Ha1(0, " aaa", Bln, P("  "))),
		},
	}
	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
