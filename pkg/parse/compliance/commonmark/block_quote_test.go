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
	"time"

	. "github.com/liamawhite/parse/test"
	"github.com/notedownorg/notedown/pkg/parse"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/test"
)

var Blocks = parse.Blocks(time.Now())

func TestCommonmarkComplianceBlocksBlockquote(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 228",
			Input:         spec(228),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Ha1(0, " Foo", P("bar\nbaz")))),
		},
		{
			Name:          "Example 229",
			Input:         spec(229),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Ha1(0, " Foo", P("bar\nbaz")))),
		},
		{
			Name:          "Example 230",
			Input:         spec(230),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("   ", Ha1(0, " Foo", P("bar\nbaz")))),
		},
		{
			Name:          "Example 231",
			Input:         spec(231),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    > # Foo\n    > bar\n    > baz")),
		},
		{
			Name:       "Example 232",
			Input:      spec(232),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", Ha1(0, " Foo", P("bar"))), P("baz")),
		},
		{
			Name:       "Example 233",
			Input:      spec(233),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", P("bar")), P("baz"), Bq("", P("foo"))),
		},
		{
			Name:          "Example 234",
			Input:         spec(234),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("foo")), Tb("-", "---")),
		},
		{
			Name:          "Example 235",
			Input:         spec(235),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Ul(Uli("", "-", " ", P("foo")))), Ul(Uli("", "-", " ", P("bar")))),
		},
		{
			Name:          "Example 236",
			Input:         spec(236),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Cbi("    foo")), Cbi("    bar")),
		},
		{
			Name:          "Example 237",
			Input:         spec(237),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Cbf("```", "", "", "")), P("foo"), Cbf("```", "", "", "")),
		},
		{
			Name:       "Example 238",
			Input:      spec(238),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", P("foo")), Cbi("    - bar")),
		},
		{
			Name:          "Example 239",
			Input:         spec(239),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Bln)),
		},
		{
			Name:          "Example 240",
			Input:         spec(240),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Bln, P(" "), Bln)),
		},
		{
			Name:          "Example 241",
			Input:         spec(241),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Bln, P("foo\n "))),
		},
		{
			Name:          "Example 242",
			Input:         spec(242),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("foo")), Bln, Bq("", P("bar"))),
		},
		{
			Name:          "Example 243",
			Input:         spec(243),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("foo\nbar"))),
		},
		{
			Name:          "Example 244",
			Input:         spec(244),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("foo"), Bln, P("bar"))),
		},
		{
			Name:          "Example 245",
			Input:         spec(245),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("foo"), Bq("", P("bar"))),
		},
		{
			Name:          "Example 246",
			Input:         spec(246),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("aaa")), Tb("*", "***"), Bq("", P("bbb"))),
		},
		{
			Name:       "Example 247",
			Input:      spec(247),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", P("bar")), P("baz")),
		},
		{
			Name:          "Example 248",
			Input:         spec(248),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("bar")), Bln, P("baz")),
		},
		{
			Name:          "Example 249",
			Input:         spec(249),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", P("bar"), Bln), P("baz")),
		},
		{
			Name:       "Example 250",
			Input:      spec(250),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", Bq("", Bq("", P("foo")))), P("bar")),
		},
		{
			Name:       "Example 251",
			Input:      spec(251),
			Parser:     Blocks,
			ExpectedOK: true,
			// Commonmark would interpret this as lazy continuation
			ExpectedMatch: Bl(Bq("", Bq("", Bq("", P("foo"))), P("bar"), Bq("", P("baz")))),
		},
		{
			Name:          "Example 252",
			Input:         spec(252),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Cbi("    code")), Bln, Bq("", P("   not code"))),
		},
	}
	RunTests(t, tests)

	// For certain roundtrips we do actually want a different result so overide the input
	// This is nearly always because we want to avoid non-blockquote parsers from knowing they are in a blockquote
	for i, test := range tests {
		switch test.Name {
		case "Example 229":
			// Add the space after the >
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> Foo\n> bar\n> baz\n"
		case "Example 230":
			// Align the blockquote with the indentation of the first line
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "   > Foo\n   > bar\n   > baz\n"
		case "Example 239":
			// Add a space after the >
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> \n"
		case "Example 240":
			// Add a space after the > where needed
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> \n>  \n> \n"
		case "Example 241":
			// Add a space after the > where needed
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> \n> foo\n>  \n"
		case "Example 244":
			// Add a space after the > where needed
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> foo\n> bar\n"
		case "Example 249":
			// Add a space after the > where needed
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> bar\n> \nbaz\n"
		case "Example 251":
			// Add a space after the > where needed
			// Doing so allows us to keep the other parsers unaware that they are in a blockquote
			tests[i].Input = "> > > foo\n> bar\n> > baz\n"
		}
	}

	VerifyRoundTrip(t, tests)
}
