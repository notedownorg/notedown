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

func TestHeadingsSetext(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "Setext H1: basic h1 with EOF",
			Input:         "H1\n========",
			Parser:        HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch: Hs1("H1", "========"),
			ExpectedOK:    true,
		},
		{
			Name:          "Setext H1: basic h1 with paragraph",
			Input:         "H1\n========\nparagraph",
			Parser:        HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch: Hs1("H1", "========", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "Setext H1: multiple h1s should return before the second",
			Input:          "H1\n========\nH1\n========",
			Parser:         HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch:  Hs1("H1", "========"),
			ExpectedOK:     true,
			RemainingInput: "H1\n========",
		},
		{
			Name:           "Setext H1: h1 with a blank line shouldn't match",
			Input:          "H1\n\n========",
			Parser:         HeadingSetextParser(ctx, 1, 2),
			ExpectedOK:     false,
			RemainingInput: "H1\n\n========",
		},
		{
			Name:          "Setext H1: cascading headers",
			Input:         "H1\n========\nH2\n--------",
			Parser:        HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch: Hs1("H1", "========", Hs2("H2", "--------")),
			ExpectedOK:    true,
		},
		{
			Name:           "Setext H1: mutliple h1s with a h2 in between",
			Input:          "H1\n========\nH2\n--------\nH1\n========",
			Parser:         HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch:  Hs1("H1", "========", Hs2("H2", "--------")),
			ExpectedOK:     true,
			RemainingInput: "H1\n========",
		},
		{
			Name:          "Setext H1: leading whitespace",
			Input:         "   H1\n========",
			Parser:        HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch: Hs1("   H1", "========"),
			ExpectedOK:    true,
		},
		{
			Name:           "Setext H1: too much leading whitespace",
			Input:          "    H1\n========",
			Parser:         HeadingSetextParser(ctx, 1, 2),
			ExpectedOK:     false,
			RemainingInput: "    H1\n========",
		},
		{
			Name:          "Setext H1: trailing whitespace",
			Input:         "H1\n========   ",
			Parser:        HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch: Hs1("H1", "========   "),
			ExpectedOK:    true,
		},
		{
			Name:          "Setext H1: underline leading whitespace",
			Input:         "H1\n   ========",
			Parser:        HeadingSetextParser(ctx, 1, 2),
			ExpectedMatch: Hs1("H1", "   ========"),
			ExpectedOK:    true,
		},
		{
			Name:           "Setext H1: underline too much whitespace",
			Input:          "H1\n    ========",
			Parser:         HeadingSetextParser(ctx, 1, 2),
			ExpectedOK:     false,
			RemainingInput: "H1\n    ========",
		},
		{
			Name:           "Setext H1: h1 only doesnt eat h2",
			Input:          "H2\n--------\nH1\n========",
			Parser:         HeadingSetextParser(ctx, 1, 1),
			ExpectedOK:     false,
			RemainingInput: "H2\n--------\nH1\n========",
		},
		{
			Name:           "Setext H1: h1 only doesnt eat paragraph",
			Input:          "paragraph\n\nH1\n========",
			Parser:         HeadingSetextParser(ctx, 1, 1),
			ExpectedOK:     false,
			RemainingInput: "paragraph\n\nH1\n========",
		},

		// H2
		{
			Name:          "Setext H2: basic h2 with EOF",
			Input:         "H2\n--------",
			Parser:        HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedMatch: Hs2("H2", "--------"),
			ExpectedOK:    true,
		},
		{
			Name:          "Setext H2: basic h2 with paragraph",
			Input:         "H2\n--------\nparagraph",
			Parser:        HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedMatch: Hs2("H2", "--------", P("paragraph")),
			ExpectedOK:    true,
		},
		{
			Name:           "Setext H2: multiple h2s should return before the second",
			Input:          "H2\n--------\nH2\n--------",
			Parser:         HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedMatch:  Hs2("H2", "--------"),
			ExpectedOK:     true,
			RemainingInput: "H2\n--------",
		},
		{
			Name:           "Setext H2: h2 with a blank line shouldn't match",
			Input:          "H2\n\n--------",
			Parser:         HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedOK:     false,
			RemainingInput: "H2\n\n--------",
		},
		{
			Name:          "Setext H2: leading whitespace",
			Input:         "   H2\n--------",
			Parser:        HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedMatch: Hs2("   H2", "--------"),
			ExpectedOK:    true,
		},
		{
			Name:           "Setext H2: too much leading whitespace",
			Input:          "    H2\n--------",
			Parser:         HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedOK:     false,
			RemainingInput: "    H2\n--------",
		},
		{
			Name:          "Setext H2: trailing whitespace",
			Input:         "H2\n--------   ",
			Parser:        HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedMatch: Hs2("H2", "--------   "),
			ExpectedOK:    true,
		},
		{
			Name:          "Setext H2: underline leading whitespace",
			Input:         "H2\n   --------",
			Parser:        HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedMatch: Hs2("H2", "   --------"),
			ExpectedOK:    true,
		},
		{
			Name:           "Setext H2: underline too much whitespace",
			Input:          "H2\n    --------",
			Parser:         HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedOK:     false,
			RemainingInput: "H2\n    --------",
		},
		{
			Name:           "Setext H2: h2 only doesnt eat h1",
			Input:          "H1\n========\nH2\n--------",
			Parser:         HeadingSetextParser(ctx, 2, 2, HeadingSetextParser(ctx, 1, 1)),
			ExpectedOK:     false,
			RemainingInput: "H1\n========\nH2\n--------",
		},
	}
	RunTests(t, tests)
}
