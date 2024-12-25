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

func TestCommonmarkComplianceBlocksHtml(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 148",
			Input:         spec(148),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<table><tr><td>\n<pre>\n**Hello**,\n"), P("_world_.\n</pre>"), Ht6("</td></tr></table>")),
		},
		{
			Name:          "Example 149",
			Input:         spec(149),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<table>\n  <tr>\n    <td>\n           hi\n    </td>\n  </tr>\n</table>\n"), P("okay.")),
		},
		{
			Name:          "Example 150",
			Input:         spec(150),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6(" <div>\n  *hello*\n         <foo><a>")),
		},
		{
			Name:          "Example 151",
			Input:         spec(151),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("</div>\n*foo*")),
		},
		{
			Name:          "Example 152",
			Input:         spec(152),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<DIV CLASS=\"foo\">\n"), P("*Markdown*"), Bln, Ht6("</DIV>")),
		},
		{
			Name:          "Example 153",
			Input:         spec(153),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div id=\"foo\"\n  class=\"bar\">\n</div>")),
		},
		{
			Name:          "Example 154",
			Input:         spec(154),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div id=\"foo\" class=\"bar\n  baz\">\n</div>")),
		},
		{
			Name:          "Example 155",
			Input:         spec(155),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div>\n*foo*\n"), P("*bar*")),
		},
		{
			Name:          "Example 156",
			Input:         spec(156),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div id=\"foo\"\n*hi*")),
		},
		{
			Name:          "Example 157",
			Input:         spec(157),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div class\nfoo")),
		},
		{
			Name:          "Example 158",
			Input:         spec(158),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div *???-&&&-<---\n*foo*")),
		},
		{
			Name:          "Example 159",
			Input:         spec(159),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div><a href=\"bar\">*foo*</a></div>")),
		},
		{
			Name:          "Example 160",
			Input:         spec(160),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<table><tr><td>\nfoo\n</td></tr></table>")),
		},
		{
			Name:          "Example 161",
			Input:         spec(161),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div></div>\n``` c\nint x = 33;\n```")),
		},
		{
			Name:          "Example 162",
			Input:         spec(162),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht7("<a href=\"foo\">\n*bar*\n</a>")),
		},
		{
			Name:          "Example 163",
			Input:         spec(163),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht7("<Warning>\n*bar*\n</Warning>")),
		},
		{
			Name:          "Example 164",
			Input:         spec(164),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht7("<i class=\"foo\">\n*bar*\n</i>")),
		},
		{
			Name:          "Example 165",
			Input:         spec(165),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht7("</ins>\n*bar*")),
		},
		{
			Name:          "Example 166",
			Input:         spec(166),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht7("<del>\n*foo*\n</del>")),
		},
		{
			Name:          "Example 167",
			Input:         spec(167),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht7("<del>\n"), P("*foo*"), Bln, Ht7("</del>")),
		},
		{
			Name:          "Example 168",
			Input:         spec(168),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("<del>*foo*</del>")),
		},
		{
			Name:          "Example 169",
			Input:         spec(169),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht1("<pre language=\"haskell\"><code>\nimport Text.HTML.TagSoup\n\nmain :: IO ()\nmain = print $ parseTags tags\n</code></pre>"), P("okay")),
		},
		{
			Name:          "Example 170",
			Input:         spec(170),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht1("<script type=\"text/javascript\">\n// JavaScript example\n\ndocument.getElementById(\"demo\").innerHTML = \"Hello JavaScript!\";\n</script>"), P("okay")),
		},
		{
			Name:          "Example 171",
			Input:         spec(171),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht1("<textarea>\n\n*foo*\n\n_bar_\n\n</textarea>")),
		},
		{
			Name:          "Example 172",
			Input:         spec(172),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht1("<style\n  type=\"text/css\">\nh1 {color:red;}\n\np {color:blue;}\n</style>"), P("okay")),
		},
		{
			Name:          "Example 173",
			Input:         spec(173),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht1("<style\n  type=\"text/css\">\n\nfoo")),
		},
		{
			Name:          "Example 174",
			Input:         spec(174),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Bq("", Ht6("<div>\nfoo")), Bln, P("bar")),
		},
		{
			Name:          "Example 175",
			Input:         spec(175),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ul(Uli("", "-", " ", Ht6("<div>")), Uli("", "-", " ", P("foo")))),
		},
		{
			Name:          "Example 176",
			Input:         spec(176),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht1("<style>p{color:red;}</style>"), P("*foo*")),
		},
		{
			Name:          "Example 177",
			Input:         spec(177),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht2("<!-- foo -->*bar*"), P("*baz*")),
		},
		{
			Name:          "Example 178",
			Input:         spec(178),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht1("<script>\nfoo\n</script>1. *bar*")),
		},
		{
			Name:          "Example 179",
			Input:         spec(179),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht2("<!-- Foo\n\nbar\n   baz -->"), P("okay")),
		},
		{
			Name:          "Example 180",
			Input:         spec(180),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht3("<?php\n\n  echo '>';\n\n?>"), P("okay")),
		},
		{
			Name:          "Example 181",
			Input:         spec(181),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht4("<!DOCTYPE html>")),
		},
		{
			Name:          "Example 182",
			Input:         spec(182),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht5("<![CDATA[\nfunction matchwo(a,b)\n{\n  if (a < b && a < 0) then {\n    return 1;\n\n  } else {\n\n    return 0;\n  }\n}\n]]>"), P("okay")),
		},
		{
			Name:          "Example 183",
			Input:         spec(183),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht2("  <!-- foo -->"), Bln, Cbi("    <!-- foo -->")),
		},
		{
			Name:          "Example 184",
			Input:         spec(184),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("  <div>\n"), Cbi("    <div>")),
		},
		{
			Name:          "Example 185",
			Input:         spec(185),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo"), Ht6("<div>\nbar\n</div>")),
		},
		{
			Name:          "Example 186",
			Input:         spec(186),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div>\nbar\n</div>\n*foo*")),
		},
		{
			Name:          "Example 187",
			Input:         spec(187),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo\n<a href=\"bar\">\nbaz")),
		},
		{
			Name:          "Example 188",
			Input:         spec(188),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div>\n"), P("*Emphasized* text."), Bln, Ht6("</div>")),
		},
		{
			Name:          "Example 189",
			Input:         spec(189),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<div>\n*Emphasized* text.\n</div>")),
		},
		{
			Name:          "Example 190",
			Input:         spec(190),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<table>\n"), Ht6("<tr>\n"), Ht6("<td>\nHi\n</td>\n"), Ht6("</tr>\n"), Ht6("</table>")),
		},
		{
			Name:          "Example 191",
			Input:         spec(191),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ht6("<table>\n"), Ht6("  <tr>\n"), Cbi("    <td>\n      Hi\n    </td>\n"), Ht6("  </tr>\n"), Ht6("</table>")),
		},
	}
	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
