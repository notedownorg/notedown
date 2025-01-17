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

func TestSingleLineBasicLists(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "bullet list item",
			Input:         "- Hello\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list item EOF",
			Input:         "- Hello",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list item with double newline",
			Input:         "- Hello\n\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("Hello"), Bln)),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list item with *",
			Input:         "* Hello\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "*", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list item with +",
			Input:         "+ Hello\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "+", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "ordered list item",
			Input:         "1. Hello\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ol(1, Oli("", "1.", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "ordered list item EOF",
			Input:         "1. Hello",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ol(1, Oli("", "1.", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "ordered list item with ) terminator",
			Input:         "1) Hello\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ol(1, Oli("", "1)", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "ordered list item with large number",
			Input:         "123456789. Hello\n", // longer than commonmark supports
			Parser:        ListParser(ctx),
			ExpectedMatch: Ol(123456789, Oli("", "123456789.", " ", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list item with external spaces",
			Input:         "   - Hello World\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("   ", "-", " ", P("Hello World"))),
			ExpectedOK:    true,
		},
		{
			Name:           "bullet list item with too many external spaces",
			Input:          "    - Hello World\n",
			Parser:         ListParser(ctx),
			RemainingInput: "    - Hello World\n",
		},
		{
			Name:          "bullet list item with three internal spaces",
			Input:         "-   Hello World\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", "   ", P("Hello World"))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list item with indented code block",
			Input:         "-\t   some code\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", Cbi("      some code"))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list item with internal tabs",
			Input:         "-\tHello\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", "\t", P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet empty list item",
			Input:         "-\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ")),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet empty list item EOF",
			Input:         "-",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ")),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet empty list item with whitespace",
			Input:         "- \n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ")),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet empty list item with more whitespace",
			Input:         "-   \n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("  "))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet empty list item with whitespace EOF",
			Input:         "-   ",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("  "))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet empty list item paragraph",
			Input:         "-\n  Hello\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", Bln, P("Hello"))),
			ExpectedOK:    true,
		},
		{
			Name:           "two blank lines not allowed in empty list item",
			Input:          "-\n\n  nope\n",
			Parser:         ListParser(ctx),
			ExpectedMatch:  Ul(Uli("", "-", " ", Bln), Bln),
			ExpectedOK:     true,
			RemainingInput: "  nope\n",
		},
	}
	RunTests(t, tests)
}

func TestNestedBasicLists(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "bullet list with nested child",
			Input:         "- Hello\n  - World\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("Hello"), Ul(Uli("", "-", " ", P("World"))))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list with nested child EOF",
			Input:         "- Hello\n  - World",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("Hello"), Ul(Uli("", "-", " ", P("World"))))),
			ExpectedOK:    true,
		},
		{
			Name:          "bullet list with nested child and sibling",
			Input:         "- Hello\n  - World\n- Goodbye\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("Hello"), Ul(Uli("", "-", " ", P("World")))), Uli("", "-", " ", P("Goodbye"))),
			ExpectedOK:    true,
		},
		{
			Name:           "bullet list followed by different marker",
			Input:          "- Hello\n* World\n",
			Parser:         ListParser(ctx),
			ExpectedMatch:  Ul(Uli("", "-", " ", P("Hello"))),
			ExpectedOK:     true,
			RemainingInput: "* World\n",
		},
		{
			Name:           "bullet list followed by ordered list",
			Input:          "- Hello\n1. World\n",
			Parser:         ListParser(ctx),
			ExpectedMatch:  Ul(Uli("", "-", " ", P("Hello"))),
			ExpectedOK:     true,
			RemainingInput: "1. World\n",
		},
		{
			Name:           "ordered list followed by different marker",
			Input:          "1. Hello\n2) World\n",
			Parser:         ListParser(ctx),
			ExpectedMatch:  Ol(1, Oli("", "1.", " ", P("Hello"))),
			ExpectedOK:     true,
			RemainingInput: "2) World\n",
		},
		{
			Name:          "bullet list with nested indented code block",
			Input:         "- paragraph\n\n      code\n",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("paragraph"), Bln, Cbi("    code"))),
			ExpectedOK:    true,
		},
	}
	RunTests(t, tests)
}
