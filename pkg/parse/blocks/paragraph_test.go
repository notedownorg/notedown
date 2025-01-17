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

func TestParagraph(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:           "Blank line interupts paragraph",
			Input:          "This is a paragraph.\n\nThis is another paragraph.",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."), // keep the newline so we can roundtrip
			ExpectedOK:     true,
			RemainingInput: "\nThis is another paragraph.",
		},
		{
			Name:          "Paragraph with newline",
			Input:         "This is a paragraph.\n",
			Parser:        ParagraphParser(ctx),
			ExpectedMatch: P("This is a paragraph."),
			ExpectedOK:    true,
		},
		{
			Name:           "Paragraph with two newlines",
			Input:          "This is a paragraph.\n\n",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "\n",
		},
		{
			Name:           "Paragraph with three newlines",
			Input:          "This is a paragraph.\n\n\n",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "\n\n",
		},
		{
			Name:           "Thematic break interupts paragraph",
			Input:          "This is a paragraph.\n---\nThis is another paragraph.",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "---\nThis is another paragraph.",
		},
		{
			Name:           "Atx heading interupts paragraph",
			Input:          "This is a paragraph.\n# This is a heading",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "# This is a heading",
		},
		{
			Name:           "Fenced code block interupts paragraph",
			Input:          "This is a paragraph.\n```\nThis is a code block\n```",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "```\nThis is a code block\n```",
		},
		{
			Name:           "HTML1 block interupts paragraph",
			Input:          "This is a paragraph.\n<pre></pre>",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "<pre></pre>",
		},
		{
			Name:           "HTML2 block interupts paragraph",
			Input:          "This is a paragraph.\n<!-- comment -->",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "<!-- comment -->",
		},
		{
			Name:           "HTML3 block interupts paragraph",
			Input:          "This is a paragraph.\n<?php echo 'Hello World'; ?>",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "<?php echo 'Hello World'; ?>",
		},
		{
			Name:           "HTML4 block interupts paragraph",
			Input:          "This is a paragraph.\n<!DOCTYPE html>",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "<!DOCTYPE html>",
		},
		{
			Name:           "HTML5 block interupts paragraph",
			Input:          "This is a paragraph.\n<![CDATA[<sender>John Smith</sender>]]>",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "<![CDATA[<sender>John Smith</sender>]]>",
		},
		{
			Name:           "HTML6 block interupts paragraph",
			Input:          "This is a paragraph.\n<div></div>",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "<div></div>",
		},
		{
			Name:           "HTML7 block doesn't interrupt paragraph",
			Input:          "This is a paragraph.\n<foo>\n\n",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph.\n<foo>"),
			ExpectedOK:     true,
			RemainingInput: "\n",
		},
		{
			Name:           "Block quote interupts paragraph",
			Input:          "This is a paragraph.\n> This is a block quote",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "> This is a block quote",
		},
		{
			Name:          "EOF match",
			Input:         "This is a paragraph.",
			Parser:        ParagraphParser(ctx),
			ExpectedMatch: P("This is a paragraph."),
			ExpectedOK:    true,
		},
		{
			Name:          "Newline then EOF match doesn't have a newline",
			Input:         "This is a paragraph.\n",
			Parser:        ParagraphParser(ctx),
			ExpectedMatch: P("This is a paragraph."),
			ExpectedOK:    true,
		},
		{
			Name:          "Leading newline persists the newline",
			Input:         "\nThis is a paragraph.",
			Parser:        ParagraphParser(ctx),
			ExpectedMatch: P("\nThis is a paragraph."),
			ExpectedOK:    true,
		},
		{
			Name:           "Leading newline for a second paragraph is unaffected",
			Input:          "This is a paragraph.\n\n\nThis is another paragraph.",
			Parser:         ParagraphParser(ctx),
			ExpectedMatch:  P("This is a paragraph."),
			ExpectedOK:     true,
			RemainingInput: "\n\nThis is another paragraph.",
		},
	}
	RunTests(t, tests)
}
