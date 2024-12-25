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
	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
)

// A paragraph is a sequence of non-blank lines that cannot be interpreted as other kinds of blocks.
// Therefore it acts as a catch-all for any text that doesn't match any other block type and must be called last.
// At parse time we maintain the whitespace (other than the blank line) so a renderer can choose to render it as it sees fit.
func ParagraphParser(ctx Context) Parser[ast.Block] {
	return func(in Input) (ast.Block, bool, error) {
		text, found, err := StringWhileNotEOFOr(paragraphCloser(ctx))(in)
		if err != nil || !found {
			return nil, false, err
		}

		// Pop the newline if not at EOF
		NewLine(in)

		return NewParagraph(text), true, nil
	}
}

func paragraphCloser(ctx Context) Parser[bool] {
	return func(in Input) (bool, bool, error) {
		start := in.Checkpoint()
		defer in.Restore(start)

		// Must start with a newline
		_, found, err := NewLine(in)
		if err != nil || !found {
			return false, false, err
		}

		// Now either a second newline or EOF terminate
		_, found, err = Or(NewLine, EOF[string]())(in)
		if err != nil || found {
			return true, true, err
		}

		// Or if we don't have a paragraph continuation
		_, found, err = paragraphContinuation(ctx)(in)
		if err != nil || !found {
			return true, true, err
		}

		return false, false, nil
	}
}

func paragraphContinuation(ctx Context) Parser[bool] {
	return func(in Input) (bool, bool, error) {
		start := in.Checkpoint()
		defer in.Restore(start)

		// Thematic break
		_, found, err := ThematicBreakParser(in)
		if err != nil || found {
			return false, false, err
		}

		// Or an atx heading
		_, found, err = headingAtxParser(1, 6)(in)
		if err != nil || found {
			return false, false, err
		}

		// Or a list item, but only a non-empty one
		// And if we have an ordered list then it must start at 1
		opener, found, err := listItemOpenParser(in)
		if err != nil || found {
			marker, ok := opener.marker.(orderedListMarker)
			if ok {
				if marker.number == 1 { // ordered list starting at 1
					return false, false, nil
				}
			} else { // unordered list
				return false, false, nil
			}
		}

		// Or a fenced code block
		_, found, err = CodeBlockFencedParser(in)
		if err != nil || found {
			return false, false, err
		}

		// Or html blocks type 1-6 (but not 7)
		_, found, err = Any(Html1Parser, Html2Parser, Html3Parser, Html4Parser, Html5Parser, Html6Parser)(in)
		if err != nil || found {
			return false, false, err
		}

		// Or a block quote
		_, found, err = BlockQuoteParser(ctx)(in)
		if err != nil || found {
			return false, false, err
		}

		return true, true, nil
	}
}
