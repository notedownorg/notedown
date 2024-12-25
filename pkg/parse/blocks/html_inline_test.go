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

	. "github.com/liamawhite/parse/core"
	. "github.com/liamawhite/parse/test"
	. "github.com/notedownorg/notedown/pkg/parse/blocks"
)

func inline[T any](p Parser[T]) Parser[string] {
	return StringFrom(AtLeast(1, p))
}

func TestHtmlInline(t *testing.T) {
	tests := []ParserTest[string]{
		// Commonmark 6.6 Specification Testcases adjusted to make them testable at this level
		// Tests with a, b, c etc. are required as a failure on the first element will stop the test so we need to test individually
		{
			Name:          "Commonmark: Example 613",
			Input:         "<a><bab><c2c>",
			Parser:        inline(HTMLOpenTag),
			ExpectedMatch: "<a><bab><c2c>",
			ExpectedOK:    true,
		},
		{
			Name:          "Commonmark: Example 614",
			Input:         "<a/><b2/>",
			Parser:        inline(HTMLOpenTag),
			ExpectedMatch: "<a/><b2/>",
			ExpectedOK:    true,
		},
		{
			Name:          "Commonmark: Example 615",
			Input:         "<a  /><b2\ndata=\"foo\" >",
			Parser:        inline(HTMLOpenTag),
			ExpectedMatch: "<a  /><b2\ndata=\"foo\" >",
			ExpectedOK:    true,
		},
		{
			Name:          "Commonmark: Example 616",
			Input:         "<a foo=\"bar\" bam = 'baz <em>\"</em>'\n_boolean zoop:33=zoop:33 />",
			Parser:        inline(HTMLOpenTag),
			ExpectedMatch: "<a foo=\"bar\" bam = 'baz <em>\"</em>'\n_boolean zoop:33=zoop:33 />",
			ExpectedOK:    true,
		},
		{
			Name:          "Commonmark: Example 617",
			Input:         "<responsive-image src=\"foo.jpg\" />",
			Parser:        inline(HTMLOpenTag),
			ExpectedMatch: "<responsive-image src=\"foo.jpg\" />",
			ExpectedOK:    true,
		},
		{
			Name:           "Commonmark: Example 618a",
			Input:          "<33>",
			Parser:         inline(HTMLOpenTag),
			ExpectedMatch:  "",
			RemainingInput: "<33>",
		},
		{
			Name:           "Commonmark: Example 618b",
			Input:          "<__>",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<__>",
		},
		{
			Name:           "Commonmark: Example 619",
			Input:          "<a h*#ref=\"hi\">",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<a h*#ref=\"hi\">",
		},
		{
			Name:           "Commonmark: Example 620a",
			Input:          "<a href=\"hi'>",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<a href=\"hi'>",
		},
		{
			Name:           "Commonmark: Example 620b",
			Input:          "<a href=hi'>",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<a href=hi'>",
		},
		{
			Name:           "Commonmark: Example 621a",
			Input:          "< a>",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "< a>",
		},
		{
			Name:           "Commonmark: Example 621b",
			Input:          "<\nfoo>",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<\nfoo>",
		},
		{
			Name:           "Commonmark: Example 621c",
			Input:          "<bar/ >",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<bar/ >",
		},
		{
			Name:           "Commonmark: Example 621d",
			Input:          "<foo bar=baz\nbim!bop />",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<foo bar=baz\nbim!bop />",
		},
		{
			Name:           "Commonmark: Example 622",
			Input:          "<a href='bar'title=title>",
			Parser:         inline(HTMLOpenTag),
			RemainingInput: "<a href='bar'title=title>",
		},
		{
			Name:          "Commonmark: Example 623",
			Input:         "</a></foo >",
			Parser:        inline(HTMLClosingTag),
			ExpectedMatch: "</a></foo >",
			ExpectedOK:    true,
		},
		{
			Name:           "Commonmark: Example 624",
			Input:          "</a href=\"foo\">",
			Parser:         inline(HTMLClosingTag),
			RemainingInput: "</a href=\"foo\">",
		},

		// Additional testcases
		{
			Name:          "HTML inline: basic tag with attributes",
			Input:         "<a href=\"foo\">",
			Parser:        inline(HTMLOpenTag),
			ExpectedMatch: "<a href=\"foo\">",
			ExpectedOK:    true,
		},
	}
	RunTests(t, tests)
}
