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

func TestHtml(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "HTML1: basic html pre",
			Input:         "<pre>some code</pre>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht1("<pre>some code</pre>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML1: basic html script",
			Input:         "<script>some code</script>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht1("<script>some code</script>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML1: basic html style",
			Input:         "<style>some code</style>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht1("<style>some code</style>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML1: basic html textarea",
			Input:         "<textarea>some code</textarea>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht1("<textarea>some code</textarea>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML1: handles new lines",
			Input:         "<pre>\nsome code\n</pre>\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht1("<pre>\nsome code\n</pre>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML1: handles trailing text",
			Input:         "<pre>some code</pre> \ttrailing text\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht1("<pre>some code</pre> \ttrailing text"),
			ExpectedOK:    true,
		},

		{
			Name:          "HTML2: basic html match",
			Input:         "<!-- some comment -->",
			Parser:        HTMLParser,
			ExpectedMatch: Ht2("<!-- some comment -->"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML2: handles new lines",
			Input:         "<!-- some\ncomment\n -->\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht2("<!-- some\ncomment\n -->"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML2: handles trailing text",
			Input:         "<!-- some comment --> \ttrailing text\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht2("<!-- some comment --> \ttrailing text"),
			ExpectedOK:    true,
		},

		{
			Name:          "HTML3: basic html match",
			Input:         "<? some code ?>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht3("<? some code ?>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML3: handles new lines",
			Input:         "<?\n some\ncode\n ?>\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht3("<?\n some\ncode\n ?>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML3: handles trailing text",
			Input:         "<? some code ?> \ttrailing text\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht3("<? some code ?> \ttrailing text"),
			ExpectedOK:    true,
		},

		{
			Name:          "HTML4: basic html match",
			Input:         "<!some code >",
			Parser:        HTMLParser,
			ExpectedMatch: Ht4("<!some code >"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML4: handles new lines",
			Input:         "<!some\ncode\n >\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht4("<!some\ncode\n >"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML4: handles trailing text",
			Input:         "<!some code > \ttrailing text\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht4("<!some code > \ttrailing text"),
			ExpectedOK:    true,
		},

		{
			Name:          "HTML5: basic html match",
			Input:         "<![CDATA[some code]]>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht5("<![CDATA[some code]]>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML5: handles new lines",
			Input:         "<![CDATA[\nsome\ncode\n]]>\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht5("<![CDATA[\nsome\ncode\n]]>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML5: handles trailing text",
			Input:         "<![CDATA[some code]]> \ttrailing text\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht5("<![CDATA[some code]]> \ttrailing text"),
			ExpectedOK:    true,
		},

		{
			Name:           "HTML6: basic html match",
			Input:          "<address>\n\nparagraph",
			Parser:         HTMLParser,
			ExpectedMatch:  Ht6("<address>\n"),
			ExpectedOK:     true,
			RemainingInput: "paragraph",
		},
		{
			Name:          "HTML6: closing html tag",
			Input:         "</article>\n\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("</article>\n"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML6: open/close html tag",
			Input:         "<aside/>\n\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("<aside/>\n"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML6: technically valid in the commonmark spec?",
			Input:         "</base/>\n\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("</base/>\n"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML6: eof",
			Input:         "<basefont>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("<basefont>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML6: opener with space",
			Input:         "<blockquote >",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("<blockquote >"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML6: opener with tab",
			Input:         "<body\t>",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("<body\t>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML6: opener with newline",
			Input:         "<caption\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("<caption\n"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML6: opener with two newlines and EOF", // first part of opener, second part of content, then EOF to close
			Input:         "<center\n\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht6("<center\n"),
			ExpectedOK:    true,
		},
		{
			Name:           "HTML6: opener with three newlines",
			Input:          "<col\n\n\nparagraph",
			Parser:         HTMLParser,
			ExpectedMatch:  Ht6("<col\n\n"),
			ExpectedOK:     true,
			RemainingInput: "paragraph",
		},

		{
			Name:           "HTML7: basic match",
			Input:          "<tag>\n\nparagraph",
			Parser:         HTMLParser,
			ExpectedMatch:  Ht7("<tag>\n"),
			ExpectedOK:     true,
			RemainingInput: "paragraph",
		},
		{
			Name:           "HTML7: invalid eof", // The spec requires at least one newline in the start condition
			Input:          "<tag>",
			Parser:         HTMLParser,
			RemainingInput: "<tag>",
		},
		{
			Name:          "HTML7: valid eof",
			Input:         "<tag>\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht7("<tag>"), // no newline to support roundtripping
			ExpectedOK:    true,
		},
		{
			Name:          "HTML7: html closing tag",
			Input:         "</tag>\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht7("</tag>"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML7: with spaces",
			Input:         "<tag \n\t>\n\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht7("<tag \n\t>\n"),
			ExpectedOK:    true,
		},
		{
			Name:          "HTML7: opener with two newlines then EOF", // first part of opener, second part of content
			Input:         "<tag>\n\n",
			Parser:        HTMLParser,
			ExpectedMatch: Ht7("<tag>\n"),
			ExpectedOK:    true,
		},
		{
			Name:           "HTML7: opener with three newlines and filler", // first part of opener, second + third part of closer
			Input:          "<tag>\nsome more html\n\nparagraph",
			Parser:         HTMLParser,
			ExpectedMatch:  Ht7("<tag>\nsome more html\n"),
			ExpectedOK:     true,
			RemainingInput: "paragraph",
		},
		{
			Name:           "HTML7: opener with three newlines",
			Input:          "<tag>\n\n\nparagraph",
			Parser:         HTMLParser,
			ExpectedMatch:  Ht7("<tag>\n"),
			ExpectedOK:     true,
			RemainingInput: "\nparagraph",
		},
	}

	RunTests(t, tests)
}
