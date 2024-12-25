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

var (
	dashBreakPart       = StringFrom(Rune('-'), OptionalInlineWhitespace)
	underscoreBreakPart = StringFrom(Rune('_'), OptionalInlineWhitespace)
	asteriskBreakPart   = StringFrom(Rune('*'), OptionalInlineWhitespace)

	dashBreak       = StringFrom(indentSlice, AtLeast(3, dashBreakPart))
	underscoreBreak = StringFrom(indentSlice, AtLeast(3, underscoreBreakPart))
	asteriskBreak   = StringFrom(indentSlice, AtLeast(3, asteriskBreakPart))

	ThematicBreakParser = func(in Input) (ast.Block, bool, error) {
		parsed, found, err := SequenceOf2(dashBreak, StringFrom(Or(NewLine, EOF[string]())))(in)
		dash, _ := parsed.Values()
		if err != nil {
			return nil, false, err
		}
		if found {
			return NewThematicBreak("-", dash), true, nil
		}
		parsed, found, err = SequenceOf2(asteriskBreak, StringFrom(Or(NewLine, EOF[string]())))(in)
		asterisk, _ := parsed.Values()
		if err != nil {
			return nil, false, err
		}
		if found {
			return NewThematicBreak("*", asterisk), true, nil
		}
		parsed, found, err = SequenceOf2(underscoreBreak, StringFrom(Or(NewLine, EOF[string]())))(in)
		underscore, _ := parsed.Values()
		if err != nil {
			return nil, false, err
		}
		if found {
			return NewThematicBreak("_", underscore), true, nil
		}
		return nil, false, nil
	}
)
