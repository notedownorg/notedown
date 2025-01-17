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

func TestThematicBreak(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "match: dashes",
			Input:         "---\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("-", "---"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: underscores",
			Input:         "___\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("_", "___"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: asterisks",
			Input:         "***\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("*", "***"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: eof",
			Input:         "* * *",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("*", "* * *"),
			ExpectedOK:    true,
		},
		{
			Name:           "no match: plus",
			Input:          "+++\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "+++\n",
		},
		{
			Name:           "no match: equals",
			Input:          "===\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "===\n",
		},
		{
			Name:           "no match: too few dashes",
			Input:          "--\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "--\n",
		},
		{
			Name:           "no match: too few asterisks",
			Input:          "**\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "**\n",
		},
		{
			Name:           "no match: too few underscores",
			Input:          "__\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "__\n",
		},
		{
			Name:          "match: one space prefix",
			Input:         " ***\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("*", " ***"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: two space prefixes",
			Input:         "  ***\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("*", "  ***"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: three space prefixes",
			Input:         "   ***\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("*", "   ***"),
			ExpectedOK:    true,
		},
		{
			Name:           "no match: four space prefixes",
			Input:          "    ***\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "    ***\n",
		},
		{
			Name:          "match: long dashes",
			Input:         "-------------------------------------\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("-", "-------------------------------------"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: spaced dashes",
			Input:         " - - -\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("-", " - - -"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: lots of spaced asterisks",
			Input:         " **  * ** * ** * **\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("*", " **  * ** * ** * **"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: wider spaced dashes",
			Input:         "-     -      -      -\n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("-", "-     -      -      -"),
			ExpectedOK:    true,
		},
		{
			Name:          "match: spaced dashes with suffixed whitespace",
			Input:         "- - - -    \n",
			Parser:        ThematicBreakParser,
			ExpectedMatch: Tb("-", "- - - -    "),
			ExpectedOK:    true,
		},
		{
			Name:           "no match: trailing characters",
			Input:          "_ _ _ _ a\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "_ _ _ _ a\n",
		},
		{
			Name:           "no match: leading characters",
			Input:          "a------\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "a------\n",
		},
		{
			Name:           "no match: surrounded characters",
			Input:          "---a---\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "---a---\n",
		},
		{
			Name:           "no match: mix and match valid characters",
			Input:          "*-*\n",
			Parser:         ThematicBreakParser,
			ExpectedOK:     false,
			RemainingInput: "*-*\n",
		},
	}
	RunTests(t, tests)
}
