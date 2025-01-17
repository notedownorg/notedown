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

func TestCommonmarkComplianceBlocksListItem(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 253",
			Input:         spec(253),
			Parser:        Blocks,
			ExpectedMatch: Bl(P("A paragraph\nwith two lines."), Bln, Cbi("    indented code\n"), Bq("", P("A block quote."))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 254",
			Input:         spec(254),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", "  ", P("A paragraph\nwith two lines."), Bln, Cbi("    indented code\n"), Bq("", P("A block quote."))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 255",
			Input:         spec(255),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("one"), Bln)), P(" two")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 256",
			Input:         spec(256),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("one"), Bln, P("two")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 257",
			Input:         spec(257),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli(" ", "-", "    ", P("one"), Bln)), Cbi("     two")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 258",
			Input:         spec(258),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli(" ", "-", "    ", P("one"), Bln, P("two")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 259",
			Input:         spec(259),
			Parser:        Blocks,
			ExpectedMatch: Bl(Bq("   ", Bq("", Ol(1, Oli("", "1.", "  ", P("one"), Bln, P("two")))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 260",
			Input:         spec(260),
			Parser:        Blocks,
			ExpectedMatch: Bl(Bq("", Bq("", Ul(Uli("", "-", " ", P("one"), Bln)), P("two")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 261",
			Input:         spec(261),
			Parser:        Blocks,
			ExpectedMatch: Bl(P("-one"), Bln, P("2.two")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 262",
			Input:         spec(262),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo"), Bln, Bln, P("bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 263",
			Input:         spec(263),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", "  ", P("foo"), Bln, Cbf("```", "", "bar", "```"), Bln, P("baz"), Bln, Bq("", P("bam"))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 264",
			Input:         spec(264),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("Foo"), Bln, Cbi("    bar\n\n\n    baz")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 265",
			Input:         spec(265),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(123456789, Oli("", "123456789.", " ", P("ok")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 266",
			Input:         spec(266),
			Parser:        Blocks,
			ExpectedMatch: Bl(P("1234567890. not ok")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 267",
			Input:         spec(267),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(0, Oli("", "0.", " ", P("ok")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 268",
			Input:         spec(268),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(3, Oli("", "003.", " ", P("ok")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 269",
			Input:         spec(269),
			Parser:        Blocks,
			ExpectedMatch: Bl(P("-1. not ok")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 270",
			Input:         spec(270),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo"), Bln, Cbi("    bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 271",
			Input:         spec(271),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(10, Oli("  ", "10.", "  ", P("foo"), Bln, Cbi("    bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 272",
			Input:         spec(272),
			Parser:        Blocks,
			ExpectedMatch: Bl(Cbi("    indented code\n"), P("paragraph"), Bln, Cbi("    more code")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 273",
			Input:         spec(273),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", Cbi("    indented code\n"), P("paragraph"), Bln, Cbi("    more code")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 274",
			Input:         spec(274),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", Cbi("     indented code\n"), P("paragraph"), Bln, Cbi("    more code")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 275",
			Input:         spec(275),
			Parser:        Blocks,
			ExpectedMatch: Bl(P("   foo"), Bln, P("bar")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 276",
			Input:         spec(276),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", "    ", P("foo"), Bln)), P("  bar")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 277",
			Input:         spec(277),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", "  ", P("foo"), Bln, P("bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 278",
			Input:         spec(278),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", Bln, P("foo")), Uli("", "-", " ", Bln, Cbf("```", "", "bar", "```")), Uli("", "-", " ", Bln, Cbi("    baz")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 279",
			Input:         spec(279),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("  \nfoo")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 280",
			Input:         spec(280),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", Bln), Bln), P("  foo")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 281",
			Input:         spec(281),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo")), Uli("", "-", " ", Bln), Uli("", "-", " ", P("bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 282",
			Input:         spec(282),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo")), Uli("", "-", " ", P("  ")), Uli("", "-", " ", P("bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 283",
			Input:         spec(283),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", P("foo")), Oli("", "2.", " ", Bln), Oli("", "3.", " ", P("bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 284",
			Input:         spec(284),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "*", " "))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 285",
			Input:         spec(285),
			Parser:        Blocks,
			ExpectedMatch: Bl(P("foo\n*"), Bln, P("foo\n1.")),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 286",
			Input:         spec(286),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli(" ", "1.", "  ", P("A paragraph\nwith two lines."), Bln, Cbi("    indented code\n"), Bq("", P("A block quote."))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 287",
			Input:         spec(287),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("  ", "1.", "  ", P("A paragraph\nwith two lines."), Bln, Cbi("    indented code\n"), Bq("", P("A block quote."))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 288",
			Input:         spec(288),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("   ", "1.", "  ", P("A paragraph\nwith two lines."), Bln, Cbi("    indented code\n"), Bq("", P("A block quote."))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 289",
			Input:         spec(289),
			Parser:        Blocks,
			ExpectedMatch: Bl(Cbi("    1.  A paragraph\n        with two lines.\n\n            indented code\n\n        > A block quote.")),
			ExpectedOK:    true,
		},
		{
			Name:   "Example 290",
			Input:  spec(290),
			Parser: Blocks,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Ol(1, Oli("  ", "1.", "  ", P("A paragraph"))), P("with two lines."), Bln, Cbi("          indented code\n\n      > A block quote.")),
			ExpectedOK:    true,
		},
		{
			Name:   "Example 291",
			Input:  spec(291),
			Parser: Blocks,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Ol(1, Oli("  ", "1.", "  ", P("A paragraph"))), Cbi("    with two lines.")),
			ExpectedOK:    true,
		},
		{
			Name:   "Example 292",
			Input:  spec(292),
			Parser: Blocks,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", Ol(1, Oli("", "1.", " ", Bq("", P("Blockquote"))))), P("continued here.")),
			ExpectedOK:    true,
		},
		{
			Name:   "Example 293",
			Input:  spec(293),
			Parser: Blocks,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", Ol(1, Oli("", "1.", " ", Bq("", P("Blockquote")))), P("continued here."))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 294",
			Input:         spec(294),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo"), Ul(Uli("", "-", " ", P("bar"), Ul(Uli("", "-", " ", P("baz"), Ul(Uli("", "-", " ", P("boo")))))))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 295",
			Input:         spec(295),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", P("foo")), Uli(" ", "-", " ", P("bar")), Uli("  ", "-", " ", P("baz")), Uli("   ", "-", " ", P("boo")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 296",
			Input:         spec(296),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(10, Oli("", "10)", " ", P("foo"), Ul(Uli("", "-", " ", P("bar")))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 297",
			Input:         spec(297),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(10, Oli("", "10)", " ", P("foo"))), Ul(Uli("   ", "-", " ", P("bar")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 298",
			Input:         spec(298),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", Ul(Uli("", "-", " ", P("foo")))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 299",
			Input:         spec(299),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", " ", Ul(Uli("", "-", " ", Ol(2, Oli("", "2.", " ", P("foo")))))))),
			ExpectedOK:    true,
		},
		{
			Name:          "Example 300",
			Input:         spec(300),
			Parser:        Blocks,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", Ha1(0, " Foo")), Uli("", "-", " ", Hs2("Bar", "---", P("baz"))))),
			ExpectedOK:    true,
		},
	}
	RunTests(t, tests)

	// For certain roundtrips we do actually want a different result so overide the input
	// This is nearly always because we want to avoid non-blockquote parsers from knowing they are in a blockquote
	for i, test := range tests {
		switch test.Name {
		case "Example 259":
			// Add the space after the >
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "   > > 1.  one\n   > > \n   > >     two\n"
		case "Example 260":
			// Add the space after the >
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> > - one\n> > \n> > two\n"
		case "Example 278":
			// Add the space after an empty list
			// Doing so allows us to keep the other parsers unaware that they are in a list
			tests[i].Input = "- \n  foo\n- \n  ```\n  bar\n  ```\n- \n      baz\n"
		case "Example 280":
			// Add the space after an empty list
			// Doing so allows us to keep the other parsers unaware that they are in a list
			tests[i].Input = "- \n\n  foo\n"
		case "Example 281":
			// Add the space after an empty list
			// Doing so allows us to keep the other parsers unaware that they are in a list
			tests[i].Input = "- foo\n- \n- bar\n"
		case "Example 283":
			// Add the space after an empty list
			// Doing so allows us to keep the other parsers unaware that they are in a list
			tests[i].Input = "1. foo\n2. \n3. bar\n"
		case "Example 284":
			// Add the space after an empty list
			// Doing so allows us to keep the other parsers unaware that they are in a list
			tests[i].Input = "* \n"
		}
	}

	VerifyRoundTrip(t, tests)
}
