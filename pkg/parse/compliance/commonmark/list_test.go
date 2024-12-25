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

func TestCommonmarkComplianceBlocksList(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 301",
			Input:         spec(301),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo")), Uli("", "-", " ", P("bar"))), Ul(Uli("", "+", " ", P("baz")))),
		},
		{
			Name:          "Example 302",
			Input:         spec(302),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", P("foo")), Oli("", "2.", " ", P("bar"))), Ol(3, Oli("", "3)", " ", P("baz")))),
		},
		{
			Name:          "Example 303",
			Input:         spec(303),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo"), Ul(Uli("", "-", " ", P("bar")), Uli("", "-", " ", P("baz")))),
		},
		{
			Name:          "Example 304",
			Input:         spec(304),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("The number of windows in my house is\n14.  The number of doors is 6.")),
		},
		{
			Name:          "Example 305",
			Input:         spec(305),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("The number of windows in my house is"), Ol(1, Oli("", "1.", "  ", P("The number of doors is 6.")))),
		},
		{
			Name:          "Example 306",
			Input:         spec(306),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo"), Bln), Uli("", "-", " ", P("bar"), Bln, Bln), Uli("", "-", " ", P("baz")))),
		},
		{
			Name:          "Example 307",
			Input:         spec(307),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo"), Ul(Uli("", "-", " ", P("bar"), Ul(Uli("", "-", " ", P("baz"), Bln, Bln, P("bim")))))))),
		},
		{
			Name:          "Example 308",
			Input:         spec(308),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo")), Uli("", "-", " ", P("bar"), Bln)), Ht2("<!-- -->"), Bln, Ul(Uli("", "-", " ", P("baz")), Uli("", "-", " ", P("bim")))),
		},
		{
			Name:          "Example 309",
			Input:         spec(309),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", "   ", P("foo"), Bln, P("notcode"), Bln), Uli("", "-", "   ", P("foo"), Bln)), Ht2("<!-- -->"), Bln, Cbi("    code")),
		},
		{
			Name:          "Example 310",
			Input:         spec(310),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a")), Uli(" ", "-", " ", P("b")), Uli("  ", "-", " ", P("c")), Uli("   ", "-", " ", P("d")), Uli("  ", "-", " ", P("e")), Uli(" ", "-", " ", P("f")), Uli("", "-", " ", P("g")))),
		},
		{
			Name:          "Example 311",
			Input:         spec(311),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", P("a"), Bln), Oli("  ", "2.", " ", P("b"), Bln), Oli("   ", "3.", " ", P("c")))),
		},
		{
			Name:       "Example 312",
			Input:      spec(312),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a")), Uli(" ", "-", " ", P("b")), Uli("  ", "-", " ", P("c")), Uli("   ", "-", " ", P("d"))), Cbi("    - e")),
		},
		{
			Name:          "Example 313",
			Input:         spec(313),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", P("a"), Bln), Oli("  ", "2.", " ", P("b"), Bln)), Cbi("    3. c")),
		},
		{
			Name:          "Example 314",
			Input:         spec(314),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a")), Uli("", "-", " ", P("b"), Bln), Uli("", "-", " ", P("c")))),
		},
		{
			Name:          "Example 315",
			Input:         spec(315),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "*", " ", P("a")), Uli("", "*", " ", Bln), Bln, Uli("", "*", " ", P("c")))),
		},
		{
			Name:          "Example 316",
			Input:         spec(316),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a")), Uli("", "-", " ", P("b"), Bln, P("c")), Uli("", "-", " ", P("d")))),
		},
		{
			Name:       "Example 317",
			Input:      spec(317),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as a link reference definition
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a")), Uli("", "-", " ", P("b"), Bln, P("[ref]: /url")), Uli("", "-", " ", P("d")))),
		},
		{
			Name:          "Example 318",
			Input:         spec(318),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a")), Uli("", "-", " ", Cbf("```", "", "b\n\n", "```")), Uli("", "-", " ", P("c")))),
		},
		{
			Name:          "Example 319",
			Input:         spec(319),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a"), Ul(Uli("", "-", " ", P("b"), Bln, P("c")))), Uli("", "-", " ", P("d")))),
		},
		{
			Name:          "Example 320",
			Input:         spec(320),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "*", " ", P("a"), Bq("", P("b"), Bln)), Uli("", "*", " ", P("c")))),
		},
		{
			Name:          "Example 321",
			Input:         spec(321),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a"), Bq("", P("b")), Cbf("```", "", "c", "```")), Uli("", "-", " ", P("d")))),
		},
		{
			Name:          "Example 322",
			Input:         spec(322),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a")))),
		},
		{
			Name:          "Example 323",
			Input:         spec(323),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a"), Ul(Uli("", "-", " ", P("b")))))),
		},
		{
			Name:          "Example 324",
			Input:         spec(324),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", Cbf("```", "", "foo", "```"), Bln, P("bar")))),
		},
		{
			Name:          "Example 325",
			Input:         spec(325),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "*", " ", P("foo"), Ul(Uli("", "*", " ", P("bar"), Bln)), P("baz")))),
		},
		{
			Name:          "Example 326",
			Input:         spec(326),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("a"), Ul(Uli("", "-", " ", P("b")), Uli("", "-", " ", P("c"), Bln))), Uli("", "-", " ", P("d"), Ul(Uli("", "-", " ", P("e")), Uli("", "-", " ", P("f")))))),
		},
	}
	RunTests(t, tests)

	// For certain roundtrips we do actually want a different result so overide the input
	// This is nearly always because we want to avoid non-blockquote parsers from knowing they are in a blockquote
	for i, test := range tests {
		switch test.Name {
		case "Example 315":
			// Add the space after an empty list item
			// Doing so allows us to keep the other parsers unaware that they are in a list
			tests[i].Input = "* a\n* \n\n* c\n"
		case "Example 320":
			// Add the space after a blockquote
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "* a\n  > b\n  > \n* c\n"
		}
	}

	VerifyRoundTrip(t, tests)
}
