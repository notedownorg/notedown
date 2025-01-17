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

func TestCommonmarkComplianceBlocksSetextHeadings(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 80",
			Input:         spec(80),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs1("Foo *bar*", "=========", Bln, Hs2("Foo *bar*", "---------"))),
		},
		{
			Name:          "Example 81",
			Input:         spec(81),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs1("Foo *bar\nbaz*", "====")),
		},
		{
			Name:          "Example 82",
			Input:         spec(82),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs1("  Foo *bar\nbaz*\t", "====")),
		},
		{
			Name:          "Example 83",
			Input:         spec(83),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2("Foo", "-------------------------", Bln), Hs1("Foo", "=")),
		},
		{
			Name:          "Example 84",
			Input:         spec(84),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2("   Foo", "---", Bln), Hs2("  Foo", "-----", Bln), Hs1("  Foo", "  ===")),
		},
		{
			Name:          "Example 85",
			Input:         spec(85),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    Foo\n    ---\n\n    Foo"), Tb("-", "---")),
		},
		{
			Name:          "Example 86",
			Input:         spec(86),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2("Foo", "   ----      ")),
		},
		{
			Name:          "Example 87",
			Input:         spec(87),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\n    ---")),
		},
		{
			Name:          "Example 88",
			Input:         spec(88),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\n= ="), Bln, P("Foo"), Tb("-", "--- -")),
		},
		{
			Name:          "Example 89",
			Input:         spec(89),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs1("Foo  ", "-----")),
		},
		{
			Name:          "Example 90",
			Input:         spec(90),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2(`Foo\`, "----")),
		},
		{
			Name:          "Example 91",
			Input:         spec(91),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2("`Foo", "----", P("`"), Bln), Hs2(`<a title="a lot`, "---", P(`of dashes"/>`))),
		},
		{
			Name:          "Example 92",
			Input:         spec(92),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("Foo")), Tb("-", "---")),
		},
		{
			Name:          "Example 93",
			Input:         spec(93),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("foo")), Hs1("bar", "===")),
			// Commonmark would interpret this as Bl(Bq("", P("foo\nbar\n==="))), but we don't support lazy continuation
		},
		{
			Name:          "Example 94",
			Input:         spec(94),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("Foo"))), Tb("-", "---")),
		},
		{
			Name:          "Example 95",
			Input:         spec(95),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2("Foo\nBar", "---")),
		},
		{
			Name:          "Example 96",
			Input:         spec(96),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("-", "---"), Hs2("Foo", "---"), Hs2("Bar", "---", P("Baz"))),
		},
		{
			Name:          "Example 97",
			Input:         spec(97),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bln, P("====")),
		},
		{
			Name:       "Example 98",
			Input:      spec(98),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as... Bl(Tb("-", "---"), Tb("-", "---")),
			// But we support frontmatter so this happens to conflict
			ExpectedMatch: Bl(Fm(nil)),
		},
		{
			Name:          "Example 99",
			Input:         spec(99),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo"))), Tb("-", "-----")),
		},
		{
			Name:          "Example 100",
			Input:         spec(100),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    foo"), Tb("-", "---")),
		},
		{
			Name:          "Example 101",
			Input:         spec(101),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("foo")), Tb("-", "-----")),
		},
		{
			Name:          "Example 102",
			Input:         spec(102),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs2(`\> foo`, "------")),
		},
		{
			Name:          "Example 103",
			Input:         spec(103),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo"), Bln, Hs2("bar", "---", P("baz"))),
		},
		{
			Name:          "Example 104",
			Input:         spec(104),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\nbar"), Bln, Tb("-", "---"), Bln, P("baz")),
		},
		{
			Name:          "Example 105",
			Input:         spec(105),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\nbar"), Tb("*", "* * *"), P("baz")),
		},
		{
			Name:          "Example 106",
			Input:         spec(106),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\nbar\n\\---\nbaz")),
		},
	}
	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
