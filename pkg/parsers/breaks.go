package parsers

import (
	"github.com/a-h/parse"
)

// 4.1 Thematic breaks
//
// A line consisting of optionally up to three spaces of indentation, followed by a sequence of three or more matching
// -, _, or * characters, each followed optionally by any number of spaces or tabs, forms a thematic break.

var whitespacePrefix = parse.AtMost(3, parse.RuneIn(" "))

var dashBreakPart = parse.StringFrom(parse.Rune('-'), remainingInlineWhitespace)
var underscoreBreakPart = parse.StringFrom(parse.Rune('_'), remainingInlineWhitespace)
var asteriskBreakPart = parse.StringFrom(parse.Rune('*'), remainingInlineWhitespace)

var dashBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, dashBreakPart), parse.Times(1, parse.NewLine))
var underscoreBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, underscoreBreakPart), parse.Times(1, parse.NewLine))
var asteriskBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, asteriskBreakPart), parse.Times(1, parse.NewLine))

var ThematicBreak = parse.Any(dashBreak, underscoreBreak, asteriskBreak)
