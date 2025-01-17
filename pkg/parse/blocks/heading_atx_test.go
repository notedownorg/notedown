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

package blocks_test

import (
	"testing"

	. "github.com/liamawhite/parse/test"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/blocks"
	. "github.com/notedownorg/notedown/pkg/parse/test"
)

func TestHeadingsAtx(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		// H1
		{
			Name:          "H1: basic h1",
			Input:         "# Hello, World!\nparagraph",
			Parser:        HeadingAtxParser(ctx, 1, 6),
			ExpectedMatch: Ha1(0, " Hello, World!", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "H1: multiple h1s should return before the second h1",
			Input:          "# Hello, World!\n# Hello, World!\n",
			Parser:         HeadingAtxParser(ctx, 1, 6),
			ExpectedMatch:  Ha1(0, " Hello, World!"),
			ExpectedOK:     true,
			RemainingInput: "# Hello, World!\n",
		},
		{
			Name:          "H1: h1 with blank line",
			Input:         "# Hello, World!\n\nparagraph",
			Parser:        HeadingAtxParser(ctx, 1, 6),
			ExpectedMatch: Ha1(0, " Hello, World!", Bln, P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:          "H1: cascading headers",
			Input:         "# H1\n## H2\n### H3\n#### H4\n##### H5\n###### H6\n",
			Parser:        HeadingAtxParser(ctx, 1, 6),
			ExpectedMatch: Ha1(0, " H1", Ha2(0, " H2", Ha3(0, " H3", Ha4(0, " H4", Ha5(0, " H5", Ha6(0, " H6")))))),
			ExpectedOK:    true,
		},
		{
			Name:           "H1: multiple h1s with a h2 in between",
			Input:          "# H1\n## H2\n# H1\n",
			Parser:         HeadingAtxParser(ctx, 1, 6),
			ExpectedMatch:  Ha1(0, " H1", Ha2(0, " H2")),
			ExpectedOK:     true,
			RemainingInput: "# H1\n",
		},

		// H2
		{
			Name:          "H2: basic h2",
			Input:         "## Hello, World!\nparagraph",
			Parser:        HeadingAtxParser(ctx, 2, 6, HeadingAtxParser(ctx, 1, 6)),
			ExpectedMatch: Ha2(0, " Hello, World!", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "H2: multiple h2s should return before the second h2",
			Input:          "## Hello, World!\n## Hello, World!\n",
			Parser:         HeadingAtxParser(ctx, 2, 6, HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha2(0, " Hello, World!"),
			ExpectedOK:     true,
			RemainingInput: "## Hello, World!\n",
		},
		{
			Name:          "H2: h2 with blank line",
			Input:         "## Hello, World!\n\nparagraph",
			Parser:        HeadingAtxParser(ctx, 2, 6, HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha2(0, " Hello, World!", Bln, P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:          "H2: cascading headers",
			Input:         "## H2\n### H3\n#### H4\n##### H5\n###### H6\n",
			Parser:        HeadingAtxParser(ctx, 2, 6, HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha2(0, " H2", Ha3(0, " H3", Ha4(0, " H4", Ha5(0, " H5", Ha6(0, " H6"))))),
			ExpectedOK:    true,
		},
		{
			Name:           "H2: multiple h2s with a h3 in between",
			Input:          "## H2\n### H3\n## H2\n",
			Parser:         HeadingAtxParser(ctx, 2, 6, HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha2(0, " H2", Ha3(0, " H3")),
			ExpectedOK:     true,
			RemainingInput: "## H2\n",
		},
		{
			Name:           "H2: h2 cancelled by h1",
			Input:          "## H2\n# H1\n",
			Parser:         HeadingAtxParser(ctx, 2, 6, HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha2(0, " H2"),
			ExpectedOK:     true,
			RemainingInput: "# H1\n",
		},

		// H3
		{
			Name:          "H3: basic h3",
			Input:         "### Hello, World!\nparagraph",
			Parser:        HeadingAtxParser(ctx, 3, 6, HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha3(0, " Hello, World!", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "H3: multiple h3s should return before the second h3",
			Input:          "### Hello, World!\n### Hello, World!\n",
			Parser:         HeadingAtxParser(ctx, 3, 6, HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha3(0, " Hello, World!"),
			ExpectedOK:     true,
			RemainingInput: "### Hello, World!\n",
		},
		{
			Name:          "H3: h3 with blank line",
			Input:         "### Hello, World!\n\nparagraph",
			Parser:        HeadingAtxParser(ctx, 3, 6, HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha3(0, " Hello, World!", Bln, P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:          "H3: cascading headers",
			Input:         "### H3\n#### H4\n##### H5\n###### H6\n",
			Parser:        HeadingAtxParser(ctx, 3, 6, HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha3(0, " H3", Ha4(0, " H4", Ha5(0, " H5", Ha6(0, " H6")))),
			ExpectedOK:    true,
		},
		{
			Name:           "H3: multiple h3s with a h4 in between",
			Input:          "### H3\n#### H4\n### H3\n",
			Parser:         HeadingAtxParser(ctx, 3, 6, HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha3(0, " H3", Ha4(0, " H4")),
			ExpectedOK:     true,
			RemainingInput: "### H3\n",
		},
		{
			Name:           "H3: h3 cancelled by h2",
			Input:          "### H3\n## H2\n",
			Parser:         HeadingAtxParser(ctx, 3, 6, HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha3(0, " H3"),
			ExpectedOK:     true,
			RemainingInput: "## H2\n",
		},
		{
			Name:           "H3: h3 cancelled by h1",
			Input:          "### H3\n# H1\n",
			Parser:         HeadingAtxParser(ctx, 3, 6, HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha3(0, " H3"),
			ExpectedOK:     true,
			RemainingInput: "# H1\n",
		},

		// H4
		{
			Name:          "H4: basic h4",
			Input:         "#### Hello, World!\nparagraph",
			Parser:        HeadingAtxParser(ctx, 4, 6, HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha4(0, " Hello, World!", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "H4: multiple h4s should return before the second h4",
			Input:          "#### Hello, World!\n#### Hello, World!\n",
			Parser:         HeadingAtxParser(ctx, 4, 6, HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha4(0, " Hello, World!"),
			ExpectedOK:     true,
			RemainingInput: "#### Hello, World!\n",
		},
		{
			Name:          "H4: h4 with blank line",
			Input:         "#### Hello, World!\n\nparagraph",
			Parser:        HeadingAtxParser(ctx, 4, 6, HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha4(0, " Hello, World!", Bln, P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:          "H4: cascading headers",
			Input:         "#### H4\n##### H5\n###### H6\n",
			Parser:        HeadingAtxParser(ctx, 4, 6, HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha4(0, " H4", Ha5(0, " H5", Ha6(0, " H6"))),
			ExpectedOK:    true,
		},
		{
			Name:           "H4: multiple h4s with a h5 in between",
			Input:          "#### H4\n##### H5\n#### H4\n",
			Parser:         HeadingAtxParser(ctx, 4, 6, HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha4(0, " H4", Ha5(0, " H5")),
			ExpectedOK:     true,
			RemainingInput: "#### H4\n",
		},
		{
			Name:           "H4: h4 cancelled by h3",
			Input:          "#### H4\n### H3\n",
			Parser:         HeadingAtxParser(ctx, 4, 6, HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha4(0, " H4"),
			ExpectedOK:     true,
			RemainingInput: "### H3\n",
		},

		// H5
		{
			Name:          "H5: basic h5",
			Input:         "##### Hello, World!\nparagraph",
			Parser:        HeadingAtxParser(ctx, 5, 6, HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha5(0, " Hello, World!", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "H5: multiple h5s should return before the second h5",
			Input:          "##### Hello, World!\n##### Hello, World!\n",
			Parser:         HeadingAtxParser(ctx, 5, 6, HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha5(0, " Hello, World!"),
			ExpectedOK:     true,
			RemainingInput: "##### Hello, World!\n",
		},
		{
			Name:          "H5: h5 with blank line",
			Input:         "##### Hello, World!\n\nparagraph",
			Parser:        HeadingAtxParser(ctx, 5, 6, HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha5(0, " Hello, World!", Bln, P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:          "H5: cascading headers",
			Input:         "##### H5\n###### H6\n",
			Parser:        HeadingAtxParser(ctx, 5, 6, HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha5(0, " H5", Ha6(0, " H6")),
			ExpectedOK:    true,
		},
		{
			Name:           "H5: multiple h5s with a h6 in between",
			Input:          "##### H5\n###### H6\n##### H5\n",
			Parser:         HeadingAtxParser(ctx, 5, 6, HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha5(0, " H5", Ha6(0, " H6")),
			ExpectedOK:     true,
			RemainingInput: "##### H5\n",
		},
		{
			Name:           "H5: h5 cancelled by h4",
			Input:          "##### H5\n#### H4\n",
			Parser:         HeadingAtxParser(ctx, 5, 6, HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha5(0, " H5"),
			ExpectedOK:     true,
			RemainingInput: "#### H4\n",
		},

		// H6
		{
			Name:          "H6: basic h6",
			Input:         "###### Hello, World!\nparagraph",
			Parser:        HeadingAtxParser(ctx, 6, 6, HeadingAtxParser(ctx, 5, 5), HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha6(0, " Hello, World!", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "H6: multiple h6s should return before the second h6",
			Input:          "###### Hello, World!\n###### Hello, World!\n",
			Parser:         HeadingAtxParser(ctx, 6, 6, HeadingAtxParser(ctx, 5, 5), HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha6(0, " Hello, World!"),
			ExpectedOK:     true,
			RemainingInput: "###### Hello, World!\n",
		},
		{
			Name:          "H6: h6 with blank line",
			Input:         "###### Hello, World!\n\nparagraph",
			Parser:        HeadingAtxParser(ctx, 6, 6, HeadingAtxParser(ctx, 5, 5), HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch: Ha6(0, " Hello, World!", Bln, P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "H6: h6 cancelled by h5",
			Input:          "###### H6\n##### H5\n",
			Parser:         HeadingAtxParser(ctx, 6, 6, HeadingAtxParser(ctx, 5, 5), HeadingAtxParser(ctx, 4, 4), HeadingAtxParser(ctx, 3, 3), HeadingAtxParser(ctx, 2, 2), HeadingAtxParser(ctx, 1, 1)),
			ExpectedMatch:  Ha6(0, " H6"),
			ExpectedOK:     true,
			RemainingInput: "##### H5\n",
		},
	}
	RunTests(t, tests)
}
