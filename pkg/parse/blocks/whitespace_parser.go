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

// Only use when you have already consumed a the preceeding newline and are looking for a blank line
var BlankLineParser = func(in Input) (ast.Block, bool, error) {
	_, found, err := NewLine(in)
	if err != nil || !found {
		return nil, false, err
	}
	return &BlankLine{}, true, nil
}

// Looks for at most 3 spaces preceeding a non-space character or rolls back.
// Useful for checking indentation for things that arent indented code blocks.
func indent(in Input) (string, bool, error) {
	parsed, found, err := indentSlice(in)
	if err != nil || !found {
		return "", false, err
	}
	return strings.Join(parsed, ""), true, nil
}

// Looks for at most 3 spaces preceeing a non-space character or rolls back.
// Useful for checking indentation for things that arent indented code blocks.
func indentSlice(in Input) ([]string, bool, error) {
	start := in.Checkpoint()
	parsed, found, err := AtMost(3, RuneIn(" "))(in)
	if err != nil || !found {
		return nil, false, err
	}

	// Check that the spaces are followed by a non-space character.
	final := in.Checkpoint()
	_, found, err = RuneNotIn(" ")(in)
	if err != nil || !found {
		in.Restore(start)
		return nil, false, err
	}

	// If we reach this point we have valid leading whitespace
	// Undo the RuneNotIn(" ") check and return the parsed whitespace
	in.Restore(final)
	return parsed, true, nil
}

// Calculate the number of spaces in a string containing both spaces and tabs
// Tab is considered to be 4 spaces
// Characters other than space and tab are not counted but also do not stop the count
func spaces(in string) int {
	var count int
	for _, r := range in {
		if r == ' ' {
			count++
		} else if r == '\t' {
			count += 4
		}
	}
	return count
}
