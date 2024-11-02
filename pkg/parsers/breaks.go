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

package parsers

import (
	"github.com/a-h/parse"
)

// 4.1 Thematic breaks
//
// A line consisting of optionally up to three spaces of indentation, followed by a sequence of three or more matching
// -, _, or * characters, each followed optionally by any number of spaces or tabs, forms a thematic break.

var whitespacePrefix = parse.AtMost(3, parse.RuneIn(" "))

var dashBreakPart = parse.StringFrom(parse.Rune('-'), RemainingInlineWhitespace)
var underscoreBreakPart = parse.StringFrom(parse.Rune('_'), RemainingInlineWhitespace)
var asteriskBreakPart = parse.StringFrom(parse.Rune('*'), RemainingInlineWhitespace)

var dashBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, dashBreakPart), parse.Times(1, parse.NewLine))
var underscoreBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, underscoreBreakPart), parse.Times(1, parse.NewLine))
var asteriskBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, asteriskBreakPart), parse.Times(1, parse.NewLine))

var ThematicBreak = parse.Any(dashBreak, underscoreBreak, asteriskBreak)
