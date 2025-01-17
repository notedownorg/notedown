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

var CodeBlockFencedParser = func(in Input) (ast.Block, bool, error) {
	start := in.Checkpoint()

	// Must start with up to 3 spaces (optional) followed by a fence
	open, found, err := SequenceOf2(indent, StringFrom(Or(AtLeast(3, Rune('`')), AtLeast(3, Rune('~')))))(in)
	if err != nil || !found {
		return nil, false, err
	}
	openIndent, fence := open.Values()

	// Look for the infostring, ignore whether or not it's found as it's optional
	infostring, _, err := infostringParser(fence)(in)
	if err != nil {
		in.Restore(start)
		return nil, false, err
	}

	// If we have a newline next we know have a valid code block
	// We just need to read the code and determine the closing fence/EOF
	if ch, _ := in.Peek(1); ch != "\n" {
		in.Restore(start)
		return nil, false, nil
	}

	// Read until we hit the closing fence
	code, found, err := StringWhileNotEOFOr(closingFence(fence))(in)
	if err != nil || !found {
		in.Restore(start)
		return nil, false, err
	}

	// Read the closing fence, if its not found we've hit EOF
	close, found, err := closingFence(fence)(in)
	if err != nil {
		in.Restore(start)
		return nil, false, err
	}

	// If we've hit EOF we need to trim a newline from the code if it exists
	if !found {
		code = strings.TrimSuffix(code, "\n")
	}

	// We also need to trim a newline from the start of the code if it exists
	code = strings.TrimPrefix(code, "\n")

	return NewCodeBlockFenced(openIndent+fence, infostring, code, close), true, nil
}

func infostringParser(fence string) Parser[string] {
	return func(in Input) (string, bool, error) {
		start := in.Checkpoint()

		// Read until we hit a NewLine
		info, found, err := StringWhileNot(NewLine)(in)
		if err != nil || !found {
			in.Restore(start)
			return "", false, err
		}

		// If the opening fence is backticks and the parsed infostring contains backticks
		// return false because this is not a valid infostring
		if strings.Contains(fence, "`") && strings.Contains(info, "`") {
			in.Restore(start)
			return "", false, err
		}

		return info, true, nil
	}
}

func closingFence(open string) Parser[string] {
	return func(in Input) (string, bool, error) {
		// Closing fence must have at least the same number of backticks or tildes
		// as the opening fence but it can have more so we need to read them all
		fence := AtLeast(len(open), Rune(rune(open[0])))

		// At least according to the examples provided in the common mark spec the
		// the closing fence must be followed by a newline
		closing, found, err := SequenceOf4(NewLine, indent, StringFrom(fence), NewLine)(in)
		if err != nil || !found {
			return "", false, err
		}

		_, closingIndent, closingFence, _ := closing.Values()
		return closingIndent + closingFence, true, nil
	}
}
