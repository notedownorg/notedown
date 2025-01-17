// Copyright 2025 Notedown Authors
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

func TestBlockQuotes(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "Single line blockquote no indent",
			Input:         "> quote!",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("", P("quote!")),
			ExpectedOK:    true,
		},
		{
			Name:          "Single line blockquote with valid indent",
			Input:         "   > quote!",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("   ", P("quote!")),
			ExpectedOK:    true,
		},
		{
			Name:          "Single line empty blockquote",
			Input:         ">",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq(""),
			ExpectedOK:    true,
		},
		{
			Name:          "Single line empty blockquote with two spaces",
			Input:         ">  ",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("", P(" ")),
			ExpectedOK:    true,
		},
		{
			Name:          "Single line empty blockquote with single space",
			Input:         "> ",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq(""),
			ExpectedOK:    true,
		},
		{
			Name:           "Single line blockquote with invalid indent",
			Input:          "    > quote!",
			Parser:         BlockQuoteParser(ctx),
			RemainingInput: "    > quote!",
		},
		{
			Name:          "Multi line blockquote",
			Input:         "> quote!\n> quote!\n",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("", P("quote!\nquote!")),
			ExpectedOK:    true,
		},
		{
			Name:          "Multi line blockquote with valid indents",
			Input:         "   > quote!\n  > quote!\n > quote!\n",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("   ", P("quote!\nquote!\nquote!")),
			ExpectedOK:    true,
		},
		{
			Name:          "Multi line blockquote with internal indents",
			Input:         ">   quote!\n>\t\tquote!\n>          quote!\n",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("", P("  quote!\n\t\tquote!\n         quote!")),
			ExpectedOK:    true,
		},
		{
			Name:          "Multi line empty blockquote with various spaces",
			Input:         ">\n>  \n> \n",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("", Bln, P(" "), Bln),
			ExpectedOK:    true,
		},
		{
			Name:          "Nested blockquotes",
			Input:         "> quote!\n> > quote!\n",
			Parser:        BlockQuoteParser(ctx),
			ExpectedMatch: Bq("", P("quote!"), Bq("", P("quote!"))),
			ExpectedOK:    true,
		},
	}
	RunTests(t, tests)
}
