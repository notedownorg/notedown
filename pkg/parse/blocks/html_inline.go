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
)

// !!!Note!!!
// We don't actually do the inline parsing phase, these are used to parse blocks.

var asciiLetter = RuneIn("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// ASCII letter followed by zero or more ASCII letters, digits, or hyphens (-).
var HTMLTagName = StringFrom(SequenceOf2(asciiLetter, ZeroOrMore(Any(asciiLetter, Digit, Rune('-')))))

var htmlAttribute = StringFrom(SequenceOf2(htmlAttributeName, Optional(htmlAttributeSpecification)))

// ASCII letter, _, or :, followed by zero or more ASCII letters, digits, _, ., :, or -. (Note: This is the XML specification restricted to ASCII. HTML5 is laxer.)
var htmlAttributeName = StringFrom(SequenceOf2(
	Any(asciiLetter, Rune('_'), Rune(':')),
	ZeroOrMore(Any(asciiLetter, Digit, Rune('_'), Rune('.'), Rune(':'), Rune('-'))),
))

// Optional spaces, tabs, and up to one line ending, a = character, optional spaces, tabs, and up to one line ending, and an attribute value.
var htmlAttributeSpecification = StringFrom(SequenceOf8(
	OptionalInlineWhitespace,
	AtMost(1, NewLine),
	OptionalInlineWhitespace,
	Rune('='),
	OptionalInlineWhitespace,
	AtMost(1, NewLine),
	OptionalInlineWhitespace,
	htmlAttributeValue,
))

var htmlAttributeValue = StringFrom(Any(htmlUnquotedAttributeValue, htmlSingleQuotedAttributeValue, htmlDoubleQuotedAttributeValue))

// nonempty string of characters not including spaces, tabs, line endings, ", ', =, <, >, or `.
var htmlUnquotedAttributeValue = StringFrom(OneOrMore(RuneNotIn(" \t\n\r\"'=<>`")))

// ', zero or more characters not including ', and a final '.
var htmlSingleQuotedAttributeValue = StringFrom(SequenceOf3(Rune('\''), ZeroOrMore(RuneNotIn("'")), Rune('\'')))

// ", zero or more characters not including ", and a final ".
var htmlDoubleQuotedAttributeValue = StringFrom(SequenceOf3(Rune('"'), ZeroOrMore(RuneNotIn("\"")), Rune('"')))

// a < character, a tag name, zero or more attributes, optional spaces, tabs, and up to one line ending, an optional / character, and a > character.
var HTMLOpenTag = func(in Input) (string, bool, error) {
	start := in.Checkpoint()

	open, found, error := StringFrom(SequenceOf2(Rune('<'), HTMLTagName))(in)
	if error != nil || !found {
		return "", false, error
	}

	// Keep reading attributes/newlines/whitespace until we find a terminating tag or a second NewLine
	newline := false
	var attributes strings.Builder
	for {
		// Check for the closing tag
		closer, found, error := StringFrom(SequenceOf2(Optional(Rune('/')), Rune('>')))(in)
		if error != nil {
			in.Restore(start)
			return "", false, error
		}
		if found {
			return open + attributes.String() + closer, true, nil
		}

		// Consume the attributes/whitespace
		// Must either be an attr followed by whitespace/newline
		// OR superfluous whitespace that we need to consume
		attr, found, err := StringFrom(Or(
			SequenceOf2(RuneIn(" \t\n"), htmlAttribute), // attrs must be space separated
			ZeroOrMore(RuneIn(" \t\n")),
		))(in)
		if err != nil || !found || attr == "" {
			in.Restore(start)
			return "", false, err
		}

		// Handle the scenario where we have a newline in the attribute
		if attr != "" {
			if strings.Contains(attr, "\n") {
				if newline {
					in.Restore(start)
					return "", false, nil
				}
				newline = true
				attributes.WriteString("\n")
			}
		}

		attributes.WriteString(attr)
	}
}

// </, a tag name, optional spaces, tabs, and up to one line ending, and the character >.
var HTMLClosingTag = StringFrom(SequenceOf6(
	String("</"),
	HTMLTagName,
	OptionalInlineWhitespace,
	AtMost(1, NewLine),
	OptionalInlineWhitespace,
	Rune('>'),
))
