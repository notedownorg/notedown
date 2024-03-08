package leaf

import (
	"github.com/a-h/parse"
)

// 4.1 Thematic breaks
//
// A line consisting of optionally up to three spaces of indentation, followed by a sequence of three or more matching
// -, _, or * characters, each followed optionally by any number of spaces or tabs, forms a thematic break.

type thematicBreakChar rune

const (
    ThematicBreakCharacterDash       thematicBreakChar = '-'
    ThematicBreakCharacterUnderscore thematicBreakChar = '_'
    ThematicBreakCharacterAsterisk   thematicBreakChar = '*'
)

type ThematicBreak struct {
    Character thematicBreakChar
}

var whitespacePrefix = parse.AtMost(3, parse.RuneIn(" "))

var dashBreakPart = parse.StringFrom(parse.Rune('-'), parse.StringFrom(parse.AtLeast(0, parse.RuneIn(" \t"))))
var underscoreBreakPart = parse.StringFrom(parse.Rune('_'), parse.StringFrom(parse.AtLeast(0, parse.RuneIn(" \t"))))
var asteriskBreakPart = parse.StringFrom(parse.Rune('*'), parse.StringFrom(parse.AtLeast(0, parse.RuneIn(" \t"))))


var dashBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, dashBreakPart), parse.Times(1, parse.NewLine))
var underscoreBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, underscoreBreakPart), parse.Times(1, parse.NewLine))
var asteriskBreak = parse.StringFrom(whitespacePrefix, parse.AtLeast(3, asteriskBreakPart), parse.Times(1, parse.NewLine))


var ThematicBreakParser parse.Parser[ThematicBreak] = parse.Func(func(in *parse.Input) (ThematicBreak, bool, error) {
    if _, ok, error := dashBreak.Parse(in); ok {
        if error != nil {
            return ThematicBreak{}, false, error
        }
        return ThematicBreak{Character: ThematicBreakCharacterDash}, true, nil
    }
    
    if _, ok, error := underscoreBreak.Parse(in); ok {
        if error != nil {
            return ThematicBreak{}, false, error
        }
        return ThematicBreak{Character: ThematicBreakCharacterUnderscore}, true, nil
    }

    if _, ok, error := asteriskBreak.Parse(in); ok {
        if error != nil {
            return ThematicBreak{}, false, error
        }
        return ThematicBreak{Character: ThematicBreakCharacterAsterisk}, true, nil
    }

    return ThematicBreak{}, false, nil
})
