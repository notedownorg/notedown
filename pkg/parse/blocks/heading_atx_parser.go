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
	. "github.com/notedownorg/notedown/pkg/parse/traversal"
)

// We split this out to avoid cyclic dependencies
func headingAtxParser(min, max int) func(in Input) (*HeadingAtx, bool, error) {
	return func(in Input) (*HeadingAtx, bool, error) {

		// These aren't valid values for min and max so we can return early
		if min < 1 || min > 6 || max < 1 || max > 6 {
			return nil, false, nil
		}

		start := in.Checkpoint()

		// Read off any leading whitespace, can be up to 3 spaces
		indent, found, err := indent(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		// Next we read the hashes
		hashes, found, err := Between(min, max, Rune('#'))(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		// Read off the text until the newline
		title, found, err := StringWhileNotEOFOr(NewLine)(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		// If the line is empty or the first character is a space, we have an empty heading so read the newline and return
		if title == "" || title[0] == ' ' {
			NewLine(in)
			return NewHeadingAtx(len(indent), len(hashes), title), true, nil
		}

		// If not we dont have a header so we roll back
		in.Restore(start)
		return nil, false, nil
	}
}

func HeadingAtxParser(ctx Context, min, max int, closers ...Parser[ast.Block]) Parser[ast.Block] {
	return func(in Input) (ast.Block, bool, error) {

		// Handle the heading itself
		h, found, err := headingAtxParser(min, max)(in)
		if err != nil || !found {
			return nil, false, err
		}

		// We exit on any of our parents closers, ourself or EOF
		// If we come across a setext heading we also close the assumption here is that if an
		// author is using a different heading type they are signalling they want to close the current heading
		closers = append(closers, HeadingAtxParser(ctx, h.level, h.level), EOF[ast.Block](), HeadingSetextParser(ctx, 1, 2))

		children, _, err := DepthFirstParser([]Parser[ast.Block]{
			BlankLineParser,
			HeadingAtxParser(ctx, h.level+1, 6, closers...),
			CodeBlockFencedParser,
			CodeBlockIndentedParser,
			ThematicBreakParser,
			ParagraphParser(ctx),
		}, closers...)(in)

		return NewHeadingAtx(h.indent, h.level, h.title, children...), true, nil
	}
}
