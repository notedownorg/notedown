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

func TestCommonmarkComplianceBlocksCodeBlockIndented(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{

		{
			Name:          "Example 107",
			Input:         spec(107),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    a simple\n      indented code block")),
		},
		{
			Name:          "Example 108",
			Input:         spec(108),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("  ", "-", " ", P("foo"), Bln, P("bar")))),
		},
		{
			Name:          "Example 109",
			Input:         spec(109),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ol(1, Oli("", "1.", "  ", P("foo"), Bln, Ul(Uli("", "-", " ", P("bar")))))),
		},
		{
			Name:          "Example 110",
			Input:         spec(110),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    <a/>\n    *hi*\n\n    - one")),
		},
		{
			Name:          "Example 111",
			Input:         spec(111),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    chunk1\n\n    chunk2\n  \n \n \n    chunk3")),
		},
		{
			Name:          "Example 112",
			Input:         spec(112),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    chunk1\n      \n      chunk2")),
		},
		{
			Name:          "Example 113",
			Input:         spec(113),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\n    bar"), Bln),
		},
		{
			Name:          "Example 114",
			Input:         spec(114),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    foo"), P("bar")),
		},
		{
			Name:          "Example 115",
			Input:         spec(115),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha1(0, " Heading", Cbi("    foo")), Hs2("Heading", "------", Cbi("    foo"), Tb("-", "----"))),
		},
		{
			Name:          "Example 116",
			Input:         spec(116),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("        foo\n    bar")),
		},
		{
			Name:          "Example 117",
			Input:         spec(117),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bln, Cbi("    \n    foo\n    \n")),
		},
		{
			Name:          "Example 118",
			Input:         spec(118),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    foo  ")),
		},
	}
	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
