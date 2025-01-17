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

func TestCommonmarkComplianceBlocksATXHeadings(t *testing.T) {
	tests := []ParserTest[[]ast.Block]{
		{
			Name:          "Example 62",
			Input:         spec(62),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha1(0, " foo", Ha2(0, " foo", Ha3(0, " foo", Ha4(0, " foo", Ha5(0, " foo", Ha6(0, " foo"))))))),
		},
		{
			Name:          "Example 63",
			Input:         spec(63),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("####### foo")),
		},
		{
			Name:          "Example 64",
			Input:         spec(64),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("#5 bolt"), Bln, P("#hashtag")),
		},
		{
			Name:          "Example 65",
			Input:         spec(65),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P(`\## foo`)),
		},
		{
			Name:          "Example 66",
			Input:         spec(66),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha1(0, ` foo *bar* \*baz\*`)),
		},
		{
			Name:          "Example 67",
			Input:         spec(67),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha1(0, `                  foo                     `)),
		},
		{
			Name:          "Example 68",
			Input:         spec(68),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha3(1, " foo"), Ha2(2, " foo"), Ha1(3, " foo")),
		},
		{
			Name:          "Example 69",
			Input:         spec(69),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Cbi("    # foo")),
		},
		{
			Name:          "Example 70",
			Input:         spec(70),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("foo\n    # bar")),
		},
		{
			Name:          "Example 71",
			Input:         spec(71),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha2(0, " foo ##", Ha3(2, "   bar    ###"))),
		},
		{
			Name:          "Example 72",
			Input:         spec(72),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha1(0, " foo ##################################", Ha5(0, " foo ##"))),
		},
		{
			Name:          "Example 73",
			Input:         spec(73),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha3(0, " foo ###     ")),
		},
		{
			Name:          "Example 74",
			Input:         spec(74),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha3(0, " foo ### b")),
		},
		{
			Name:          "Example 75",
			Input:         spec(75),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha1(0, " foo#")),
		},
		{
			Name:          "Example 76",
			Input:         spec(76),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha3(0, ` foo \###`), Ha2(0, ` foo #\##`), Ha1(0, ` foo \#`)),
		},
		{
			Name:          "Example 77",
			Input:         spec(77),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Tb("*", "****"), Ha2(0, " foo", Tb("*", "****"))),
		},
		{
			Name:          "Example 78",
			Input:         spec(78),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(P("Foo bar"), Ha1(0, " baz", P("Bar foo"))),
		},
		{
			Name:          "Example 79",
			Input:         spec(79),
			Parser:        Blocks,
			ExpectedOK:    true,
			ExpectedMatch: Bl(Ha2(0, " "), Ha1(0, "", Ha3(0, " ###"))),
		},
	}
	RunTests(t, tests)
	VerifyRoundTrip(t, tests)
}
