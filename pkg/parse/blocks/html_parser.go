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

package blocks

import (
	"strings"

	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
)

var HTMLParser = Any(Html1Parser, Html2Parser, Html3Parser, Html4Parser, Html5Parser, Html6Parser, Html7Parser)

// Start condition: line begins with the string <pre, <script, <style, or <textarea (case-insensitive), followed by a space, a tab, the string >, or the end of the line.
// End condition: line contains an end tag </pre>, </script>, </style>, or </textarea> (case-insensitive; it need not match the start tag).
var Html1Parser = func(in Input) (ast.Block, bool, error) {
	opener := StringFrom(SequenceOf2(
		Any(String("<pre"), String("<script"), String("<style"), String("<textarea")),
		Any(String(" "), String("\t"), String(">"), NewLine),
	))
	closer := func(in Input) (string, bool, error) {
		tag, found, err := Any(String("</pre>"), String("</script>"), String("</style>"), String("</textarea>"))(in)
		if err != nil || !found {
			return "", false, err
		}

		// Everything left on the line is considered part of the block so read it
		rem, _, err := StringWhileNot(NewLine)(in)

		// Need to consume the newline to make roundtripping work
		NewLine(in)

		return tag + rem, true, nil
	}
	return htmlBuilder(HtmlOne, opener, closer)(in)
}

// Start condition: line begins with the string <!--.
// End condition: line contains the string -->.
var Html2Parser = func(in Input) (ast.Block, bool, error) {
	opener := String("<!--")
	closer := func(in Input) (string, bool, error) {
		tag, found, err := String("-->")(in)
		if err != nil || !found {
			return "", false, err
		}

		// Everything left on the line is considered part of the block so read it
		rem, _, err := StringWhileNot(NewLine)(in)

		// Need to consume the newline to make roundtripping work
		NewLine(in)

		return tag + rem, true, nil
	}
	return htmlBuilder(HtmlTwo, opener, closer)(in)
}

// Start condition: line begins with the string <?.
// End condition: line contains the string ?>.
var Html3Parser = func(in Input) (ast.Block, bool, error) {
	opener := String("<?")
	closer := func(in Input) (string, bool, error) {
		tag, found, err := String("?>")(in)
		if err != nil || !found {
			return "", false, err
		}

		// Everything left on the line is considered part of the block so read it
		rem, _, err := StringWhileNot(NewLine)(in)

		// Need to consume the newline to make roundtripping work
		NewLine(in)

		return tag + rem, true, nil
	}
	return htmlBuilder(HtmlThree, opener, closer)(in)
}

// Start condition: line begins with the string <! followed by an ASCII letter.
// End condition: line contains the character >.
var Html4Parser = func(in Input) (ast.Block, bool, error) {
	opener := StringFrom(SequenceOf2(String("<!"), asciiLetter))
	closer := func(in Input) (string, bool, error) {
		tag, found, err := String(">")(in)
		if err != nil || !found {
			return "", false, err
		}

		// Everything left on the line is considered part of the block so read it
		rem, _, err := StringWhileNot(NewLine)(in)

		// Need to consume the newline to make roundtripping work
		NewLine(in)

		return tag + rem, true, nil
	}
	return htmlBuilder(HtmlFour, opener, closer)(in)
}

// Start condition: line begins with the string <![CDATA[.
// End condition: line contains the string ]]>.
var Html5Parser = func(in Input) (ast.Block, bool, error) {
	opener := String("<![CDATA[")
	closer := func(in Input) (string, bool, error) {
		tag, found, err := String("]]>")(in)
		if err != nil || !found {
			return "", false, err
		}

		// Everything left on the line is considered part of the block so read it
		rem, _, err := StringWhileNot(NewLine)(in)

		// Need to consume the newline to make roundtripping work
		NewLine(in)

		return tag + rem, true, nil
	}
	return htmlBuilder(HtmlFive, opener, closer)(in)
}

// Each block should be able to be concatenated with the next block with a newline to form the original markdown
// Therefore we need to look for two newlines to end the block but only include one in the content
var htmlNewlineCloser = func(in Input) (string, bool, error) {
	start := in.Checkpoint()

	// We need to control whether or not we're at the end of the file to handle newlines correctly
	// If we're at the EOF return empty but true
	_, eof, err := EOF[bool]()(in)
	if err != nil {
		return "", false, err
	}
	if eof {
		return "", true, nil
	}

	// Now we either need a newline followed by EOF or another newline
	_, found, err := NewLine(in)
	if err != nil || !found {
		in.Restore(start)
		return "", false, err
	}

	// If we're at EOF we need to return true but empty
	_, eof, err = EOF[bool]()(in)
	if err != nil {
		in.Restore(start)
		return "", false, err
	}
	if eof {
		return "", true, nil
	}

	// If we have a second newline we need to return true + \n
	_, nl, err := NewLine(in)
	if err != nil {
		in.Restore(start)
		return "", false, err
	}
	if nl {
		return "\n", true, nil
	}

	// Otherwise we tell the parser to continue
	in.Restore(start)
	return "", false, nil
}

// Start condition: line begins with the string < or </ followed by one of the strings (case-insensitive) address, article, aside, base, basefont, blockquote, body, caption, center, col, colgroup, dd, details, dialog, dir, div, dl, dt, fieldset, figcaption, figure, footer, form, frame, frameset, h1, h2, h3, h4, h5, h6, head, header, hr, html, iframe, legend, li, link, main, menu, menuitem, nav, noframes, ol, optgroup, option, p, param, search, section, summary, table, tbody, td, tfoot, th, thead, title, tr, track, ul, followed by a space, a tab, the end of the line, the string >, or the string />.
// End condition: line is followed by a blank line.
var Html6Parser = func(in Input) (ast.Block, bool, error) {
	opener := StringFrom(SequenceOf3(
		StringFrom(SequenceOf2(String("<"), Optional(String("/")))),
		Any(
			// More specific matches have to be first
			StringInsensitive("address"), StringInsensitive("article"), StringInsensitive("aside"), StringInsensitive("basefont"), StringInsensitive("base"), StringInsensitive("blockquote"), StringInsensitive("body"), StringInsensitive("caption"), StringInsensitive("center"), StringInsensitive("colgroup"), StringInsensitive("col"), StringInsensitive("dd"), StringInsensitive("details"), StringInsensitive("dialog"), StringInsensitive("dir"), StringInsensitive("div"), StringInsensitive("dl"), StringInsensitive("dt"), StringInsensitive("fieldset"), StringInsensitive("figcaption"), StringInsensitive("figure"), StringInsensitive("footer"), StringInsensitive("form"), StringInsensitive("frame"), StringInsensitive("frameset"), StringInsensitive("h1"), StringInsensitive("h2"), StringInsensitive("h3"), StringInsensitive("h4"), StringInsensitive("h5"), StringInsensitive("h6"), StringInsensitive("header"), StringInsensitive("head"), StringInsensitive("hr"), StringInsensitive("html"), StringInsensitive("iframe"), StringInsensitive("legend"), StringInsensitive("link"), StringInsensitive("li"), StringInsensitive("main"), StringInsensitive("menuitem"), StringInsensitive("menu"), StringInsensitive("nav"), StringInsensitive("noframes"), StringInsensitive("ol"), StringInsensitive("optgroup"), StringInsensitive("option"), StringInsensitive("param"), StringInsensitive("p"), StringInsensitive("search"), StringInsensitive("section"), StringInsensitive("summary"), StringInsensitive("table"), StringInsensitive("tbody"), StringInsensitive("td"), StringInsensitive("tfoot"), StringInsensitive("thead"), StringInsensitive("th"), StringInsensitive("title"), StringInsensitive("track"), StringInsensitive("tr"), StringInsensitive("ul"),
		),
		Any(String(" "), String("\t"), StringFrom(SequenceOf2(Optional(String("/")), String(">"))), String("\n")),
	))

	return htmlBuilder(HtmlSix, opener, htmlNewlineCloser)(in)
}

// Start condition: line begins with a complete open tag (with any tag name other than pre, script, style, or textarea) or a complete closing tag, followed by zero or more spaces and tabs, followed by the end of the line.
// End condition: line is followed by a blank line.
var Html7Parser = func(in Input) (ast.Block, bool, error) {
	opener := func(in Input) (string, bool, error) {
		open := func(in Input) (string, bool, error) {
			tag, found, err := HTMLOpenTag(in)
			if err != nil || !found {
				return "", false, err
			}
			if strings.HasPrefix(tag, "<pre") || strings.HasPrefix(tag, "<script") || strings.HasPrefix(tag, "<style") || strings.HasPrefix(tag, "<textarea") {
				return "", false, nil
			}
			return tag, true, nil
		}
		parsed, found, err := StringFrom(
			SequenceOf2(
				Or(open, HTMLClosingTag),
				OptionalInlineWhitespace,
			),
		)(in)
		if err != nil || !found {
			return "", false, err
		}

		// Has to be followed by a newline but we don't want to include it or
		// consume it to make roundtripping easier
		final := in.Checkpoint()
		_, found, err = NewLine(in)
		if err != nil || !found {
			return "", false, err
		}
		in.Restore(final)
		return parsed, true, nil
	}

	return htmlBuilder(HtmlSeven, opener, htmlNewlineCloser)(in)
}

func htmlBuilder(kind HtmlKind, opener, closer Parser[string]) Parser[ast.Block] {
	return func(in Input) (ast.Block, bool, error) {
		start := in.Checkpoint()

		// Can be indented up to 3 spaces
		ind, _, err := indent(in)
		if err != nil {
			return nil, false, err
		}

		open, found, err := opener(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		filling, found, err := StringWhileNotEOFOr(closer)(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		// If we did hit an EOF we need to trim a newline from the returned block
		_, eof, err := EOF[bool]()(in)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}
		if eof {
			return NewHtml(kind, ind+open+strings.TrimSuffix(filling, "\n")), true, nil
		}

		// Otherwise return with the closer
		close, _, err := closer(in)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}

		return NewHtml(kind, ind+open+filling+close), true, nil
	}
}
