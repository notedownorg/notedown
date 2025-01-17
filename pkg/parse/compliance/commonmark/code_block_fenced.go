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

func TestCommonmarkComplianceBlocksCodeBlockFenced(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 119",
			Input:         spec(119),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "<\n >", "```")),
		},
		{
			Name:          "Example 120",
			Input:         spec(120),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("~~~", "", "<\n >", "~~~")),
		},
		{
			Name:          "Example 121",
			Input:         spec(121),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("``\nfoo\n``")),
		},
		{
			Name:          "Example 122",
			Input:         spec(122),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "aaa\n~~~", "```")),
		},
		{
			Name:          "Example 123",
			Input:         spec(123),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("~~~", "", "aaa\n```", "~~~")),
		},
		{
			Name:          "Example 124",
			Input:         spec(124),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("````", "", "aaa\n```", "``````")),
		},
		{
			Name:          "Example 125",
			Input:         spec(125),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("~~~~", "", "aaa\n~~~", "~~~~")),
		},
		{
			Name:          "Example 126",
			Input:         spec(126),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "", "")),
		},
		{
			Name:          "Example 127",
			Input:         spec(127),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("`````", "", "\n```\naaa", "")),
		},
		{
			Name:          "Example 128",
			Input:         spec(128),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Cbf("```", "", "aaa", "")), P("bbb")),
		},
		{
			Name:          "Example 129",
			Input:         spec(129),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "\n  ", "```")),
		},
		{
			Name:          "Example 130",
			Input:         spec(130),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "", "```")),
		},
		{
			Name:          "Example 131",
			Input:         spec(131),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf(" ```", "", " aaa\naaa", "```")),
		},
		{
			Name:          "Example 132",
			Input:         spec(132),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("  ```", "", "aaa\n  aaa\naaa", "  ```")),
		},
		{
			Name:          "Example 133",
			Input:         spec(133),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("   ```", "", "   aaa\n    aaa\n  aaa", "   ```")),
		},
		{
			Name:          "Example 134",
			Input:         spec(134),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    ```\n    aaa\n    ```")),
		},
		{
			Name:          "Example 135",
			Input:         spec(135),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "aaa", "  ```")),
		},
		{
			Name:          "Example 136",
			Input:         spec(136),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("   ```", "", "aaa", "  ```")),
		},
		{
			Name:          "Example 137",
			Input:         spec(137),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "aaa\n    ```", "")),
		},
		{
			Name:          "Example 138",
			Input:         spec(138),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("``` ```\naaa")),
		},
		{
			Name:          "Example 139",
			Input:         spec(139),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("~~~~~~", "", "aaa\n~~~ ~~", "")),
		},
		{
			Name:          "Example 140",
			Input:         spec(140),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("foo"), Cbf("```", "", "bar", "```"), P("baz")),
		},
		{
			Name:          "Example 141",
			Input:         spec(141),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Hs1("foo", "---", Cbf("~~~", "", "bar", "~~~")), Ha1(0, " baz")),
		},
		{
			Name:          "Example 142",
			Input:         spec(142),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "ruby", "def foo(x)\n  return 3\nend", "```")),
		},
		{
			Name:          "Example 143",
			Input:         spec(143),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("~~~~", "    ruby startline=3 $%@#$", "def foo(x)\n  return 3\nend", "~~~~~~~")),
		},
		{
			Name:          "Example 144",
			Input:         spec(144),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("````", ";", "", "````")),
		},
		{
			Name:          "Example 145",
			Input:         spec(145),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("``` aa ```\nfoo")),
		},
		{
			Name:          "Example 146",
			Input:         spec(146),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("~~~", " aa ``` ~~~", "foo", "~~~")),
		},
		{
			Name:          "Example 147",
			Input:         spec(147),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbf("```", "", "``` aaa", "```")),
		},
	}
	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
