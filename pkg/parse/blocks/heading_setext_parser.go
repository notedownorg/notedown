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
	. "github.com/notedownorg/notedown/pkg/parse/traversal"
)

// We split this out to avoid cyclic dependencies
func headingSetextParser(min, max int) func(in Input) (*HeadingSetext, bool, error) {
	return func(in Input) (*HeadingSetext, bool, error) {

		// These aren't valid values for min and max so we can return early
		if min < 1 || min > 2 || max < 1 || max > 2 {
			return nil, false, nil
		}

		start := in.Checkpoint()

		// Avoid reading a blank line
		if _, found, err := NewLine(in); err != nil || found {
			in.Restore(start)
			return nil, false, err
		}

		// Read off any leading whitespace, can be up to 3 spaces
		indent, found, err := indent(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		underlineParser := SequenceOf5(
			NewLine,
			indentSlice,
			StringFrom(Or(AtLeast(1, Rune('=')), AtLeast(1, Rune('-')))),
			OptionalInlineWhitespace,
			Or(NewLine, EOF[string]()),
		)

		// Read off the title, can be >= 1 lines of text not interrupted by a blank line or an underline sequence
		title, found, err := StringWhileNotEOFOr(Or(underlineParser, Times(2, NewLine)))(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		underlineSequence, found, err := underlineParser(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		_, underlinePrefix, underlinePart, underlineTrailing, _ := underlineSequence.Values()
		underline := strings.Join(append(underlinePrefix, underlinePart, underlineTrailing), "")

		// Verify that the underline sequence is at the correct level
		// It might be possible to do this as part of the underlineParser but my brain started to melt
		if min == 1 && max == 1 {
			if strings.Contains(underline, "-") {
				in.Restore(start)
				return nil, false, nil
			}
		}
		if min == 2 && max == 2 {
			if strings.Contains(underline, "=") {
				in.Restore(start)
				return nil, false, nil
			}
		}

		return NewHeadingSetext(indent+title, underline), true, nil
	}
}

func underlinePartParser(min, max int) Parser[string] {
	return func(in Input) (string, bool, error) {
		// If both min and max are 1, we are looking for `=`
		if min == 1 && max == 1 {
			return StringFrom(AtLeast(1, Rune('=')))(in)
		}
		// If both min and max are 2, we are looking for `-`
		if min == 2 && max == 2 {
			return StringFrom(AtLeast(1, Rune('-')))(in)
		}
		// Otherwise, we are looking for either `=` or `-`
		return StringFrom(Or(AtLeast(1, Rune('=')), AtLeast(1, Rune('-'))))(in)
	}
}

func HeadingSetextParser(ctx Context, min, max int, closers ...Parser[ast.Block]) Parser[ast.Block] {
	return func(in Input) (ast.Block, bool, error) {
		// Handle the heading
		h, found, err := headingSetextParser(min, max)(in)
		if err != nil || !found {
			return nil, false, err
		}

		// We exit on any of our parents closers, ourself or EOF
		// If we come across an atx heading we also close the assumption here is that if an author is
		// using a different heading type they are signalling they want to close the current heading
		closers = append(closers, HeadingSetextParser(ctx, h.level, h.level), HeadingAtxParser(ctx, 1, 6), EOF[ast.Block]())

		children, _, err := DepthFirstParser([]Parser[ast.Block]{
			BlankLineParser,
			HeadingSetextParser(ctx, h.level+1, 2, closers...),
			CodeBlockFencedParser,
			CodeBlockIndentedParser,
			ThematicBreakParser,
			ParagraphParser(ctx),
		}, closers...)(in)

		return NewHeadingSetext(h.title, h.underline, children...), true, nil
	}
}
