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

func codeBlockIndentedDelimeter(in Input) (bool, bool, error) {
	start := in.Checkpoint()
	defer in.Restore(start)

	// If theres no newline we haven't found a delimiter
	_, found, err := NewLine(in)
	if err != nil || !found {
		return false, false, err
	}

	// Check the next 4 characters to see if we have enough indentation
	for i := 1; i <= 4; i++ {
		ch, ok := in.Take(1)

		// If we reach EOF we have found the delimiter
		if !ok {
			return true, true, nil
		}

		// If any character is a tab we have enough indentation so haven't found the delimiter
		if ch == "\t" {
			return false, false, nil
		}

		// If we come across a newline we have a blank/whitespace only line so haven't found the delimiter
		if ch == "\n" {
			return false, false, nil
		}

		// If we come across a non-space character we have found the delimiter
		if ch != " " {
			return true, true, nil
		}
	}

	// If we've reached this point we've not found the delimiter
	return false, false, nil
}

var CodeBlockIndentedParser = func(in Input) (ast.Block, bool, error) {
	start := in.Checkpoint()

	// We need at least 4 spaces or a tab to start an indented code block
	prefix, found, err := StringFrom(Or(Times(4, Rune(' ')), Times(1, Rune('\t'))))(in)
	if err != nil || !found {
		return nil, false, err
	}

	// Read until we come across a line with between 1 and 3 spaces (blank lines are allowed)
	code, found, err := StringWhileNotEOFOr(codeBlockIndentedDelimeter)(in)
	if err != nil || !found {
		in.Restore(start)
		return nil, false, err
	}

	// Split the code into lines and store them
	lines := strings.Split(code, "\n")

	// Add the prefix to the first line
	if len(lines) > 0 {
		lines[0] = prefix + lines[0]
	}

	// Consume the newline as we don't consume it in the delimiter
	NewLine(in)

	return NewCodeBlockIndented(lines), true, nil
}
