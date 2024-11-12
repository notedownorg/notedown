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

import "github.com/a-h/parse"

var InlineWhitespaceRunes = parse.RuneIn(" \t")

type leadingWhitespaceParser[T any] struct {
	p parse.Parser[T]
}

func (l leadingWhitespaceParser[T]) Parse(in *parse.Input) (T, bool, error) {
	// Read the leading whitespace.
	_, ok, err := parse.OneOrMore(InlineWhitespaceRunes).Parse(in)
	if err != nil || !ok {
		return *new(T), ok, err
	}
	return l.p.Parse(in)
}

func LeadingWhitespace[T any](p parse.Parser[T]) leadingWhitespaceParser[T] {
	return leadingWhitespaceParser[T]{p: p}
}

var RemainingInlineWhitespace = parse.StringFrom(parse.ZeroOrMore(InlineWhitespaceRunes))

var remainingWhitespace = parse.StringFrom(parse.ZeroOrMore(parse.Any(InlineWhitespaceRunes, parse.NewLine)))

var NewLineOrEOF = parse.Any(parse.NewLine, parse.EOF[string]())
