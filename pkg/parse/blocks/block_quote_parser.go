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

package blocks

import (
	"strings"

	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/traversal"
)

func BlockQuoteParser(ctx Context) Parser[ast.Block] {
	return func(in Input) (ast.Block, bool, error) {
		start := in.Checkpoint()

		// Final space is optional because the spec allows for block quotes to be empty
		// The spec considers this final space to be part of the prefix
		prefixParser := SequenceOf3(indent, Rune('>'), Optional(Rune(' ')))

		// Check if we are at the start of a block quote and return early if not
		first, found, err := prefixParser(in)
		if err != nil || !found {
			return nil, false, nil
		}
		indentation, _, _ := first.Values()

		// Reset so we can trim the prefix from the first line
		in.Restore(start)

		// Read the block quote line by line until we reach a line that doesn't start with the block quote prefix
		// Strip the block quote prefixes from each line before passing the input to the children
		// Otherwise we would have to make every parser aware of the blockquote prefix and that would be a mess
		var s strings.Builder
		for {
			line, found, err := SequenceOf2(prefixParser, StringWhileNotEOFOr(NewLine))(in)
			if err != nil {
				in.Restore(start)
				return nil, false, err
			}

			// If we didn't find a block quote line we're done
			// We know we have at least one line because we've already checked before entering the loop
			if !found {
				break
			}

			_, content := line.Values()
			s.WriteString(content)

			// If we are at the end of the file, break
			if _, found, _ := EOF[bool]()(in); found {
				break
			}

			// If we're not at the end of the file, pop the newline and add one to the string
			NewLine(in)
			s.WriteString("\n")

		}

		// Build a new input from the string we've built and parse the children
		bqInput := NewInput(s.String())

		children, _, err := DepthFirstParser(TopLevelParser(ctx), EOF[ast.Block]())(bqInput)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}

		return NewBlockQuote(indentation, children...), true, nil
	}
}
