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

func TestCommonmarkComplianceBlocksThematicBreaks(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 43",
			Input:         spec(43),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("*", "***"), Tb("-", "---"), Tb("_", "___")),
		},
		{
			Name:          "Example 44",
			Input:         spec(44),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("+++")),
		},
		{
			Name:          "Example 45",
			Input:         spec(45),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("===")),
		},
		{
			Name:          "Example 46",
			Input:         spec(46),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("--\n**\n__")),
		},
		{
			Name:          "Example 47",
			Input:         spec(47),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("*", " ***"), Tb("*", "  ***"), Tb("*", "   ***")),
		},
		{
			Name:          "Example 48",
			Input:         spec(48),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    ***")),
		},
		{
			Name:          "Example 49",
			Input:         spec(49),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\n    ***")),
		},
		{
			Name:          "Example 50",
			Input:         spec(50),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("_", "_____________________________________")),
		},
		{
			Name:          "Example 51",
			Input:         spec(51),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("-", " - - -")),
		},
		{
			Name:          "Example 52",
			Input:         spec(52),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("*", " **  * ** * ** * **")),
		},
		{
			Name:          "Example 53",
			Input:         spec(53),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("-", "-     -      -      -")),
		},
		{
			Name:          "Example 54",
			Input:         spec(54),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("-", "- - - -    ")),
		},
		{
			Name:          "Example 55",
			Input:         spec(55),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("_ _ _ _ a"), Bln, P("a------"), Bln, P("---a---")),
		},
		{
			Name:          "Example 56",
			Input:         spec(56),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P(" *-*")),
		},
		{
			Name:          "Example 57",
			Input:         spec(57),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo"))), Tb("*", "***"), Ul(Uli("", "-", " ", P("bar")))),
		},
		{
			Name:          "Example 58",
			Input:         spec(58),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo"), Tb("*", "***"), P("bar")),
		},
		{
			Name:          "Example 59",
			Input:         spec(59),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2("Foo", "---", P("bar"))),
		},
		{
			Name:          "Example 60",
			Input:         spec(60),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "*", " ", P("Foo"))), Tb("*", "* * *"), Ul(Uli("", "*", " ", P("Bar")))),
		},
		{
			Name:          "Example 61",
			Input:         spec(61),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("Foo")), Uli("", "-", " ", Tb("*", "* * *")))),
		},
	}

	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
