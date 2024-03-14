package parsers

import "github.com/a-h/parse"

var inlineWhitespaceRunes = parse.RuneIn(" \t")

var remainingInlineWhitespace = parse.StringFrom(parse.ZeroOrMore(inlineWhitespaceRunes))

var remainingWhitespace = parse.StringFrom(parse.ZeroOrMore(parse.Any(inlineWhitespaceRunes, parse.NewLine)))

var newLineOrEOF = parse.Any(parse.NewLine, parse.EOF[string]())
