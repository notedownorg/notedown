// Copyright 2025 Notedown Authors
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
	"fmt"
	"strconv"

	. "github.com/liamawhite/parse/core"
)

func listMarkerParser(in Input) (listItemMarker, bool, error) {
	return Any(bulletListMarkerParser, orderedListMarkerParser)(in)
}

type bulletListMarker struct {
	marker string
}

func NewBulletListMarker(marker string) listItemMarker {
	return bulletListMarker{marker}
}

func (b bulletListMarker) String() string {
	return b.marker
}

func (b bulletListMarker) SameType(other listItemMarker) bool {
	otherBullet, ok := other.(bulletListMarker)
	if !ok {
		return false
	}
	return b.marker == otherBullet.marker
}

func bulletListMarkerParser(in Input) (listItemMarker, bool, error) {
	marker, found, err := Any(Rune('-'), Rune('*'), Rune('+'))(in)
	if err != nil || !found {
		return bulletListMarker{}, false, err
	}
	return bulletListMarker{string(marker)}, true, nil
}

type orderedListMarker struct {
	numStr     string
	number     int
	terminator string
}

func NewOrderedListMarker(number string, terminator string) listItemMarker {
	num, _ := strconv.Atoi(number)
	return orderedListMarker{number: num, numStr: number, terminator: terminator}
}

func (o orderedListMarker) String() string {
	return fmt.Sprintf("%s%s", o.numStr, o.terminator)
}

func (o orderedListMarker) SameType(other listItemMarker) bool {
	otherOrdered, ok := other.(orderedListMarker)
	if !ok {
		return false
	}
	return o.terminator == otherOrdered.terminator
}

func orderedListMarkerParser(in Input) (listItemMarker, bool, error) {
	parsed, found, err := SequenceOf2(StringFrom(OneOrMore(Digit)), RuneIn(".)"))(in)
	if err != nil || !found {
		return orderedListMarker{}, false, err
	}

	// Cant be > 9 digits
	numStr, terminator := parsed.Values()
	if len(numStr) > 9 {
		return orderedListMarker{}, false, nil
	}
	number, err := strconv.Atoi(numStr)
	if err != nil {
		return orderedListMarker{}, false, nil
	}
	// or less than 0
	if number < 0 {
		return orderedListMarker{}, false, nil
	}
	return orderedListMarker{numStr, number, terminator}, true, nil
}

type listItemMarker interface {
	fmt.Stringer
	SameType(other listItemMarker) bool
}
