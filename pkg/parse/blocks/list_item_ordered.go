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
	"strconv"
	"strings"

	"github.com/notedownorg/notedown/pkg/parse/ast"
)

func ListItemOrdered(li ast.Block) (*listItemOrdered, bool) {
	l, ok := li.(*listItemOrdered)
	return l, ok
}

var _ ast.Block = &listItemOrdered{}
var _ listItem = &listItemOrdered{}

const ListItemOrderedBlockType = "list_item_ordered"

type listItemOrdered struct {
	*tracker
	external string
	internal string
	marker   orderedListItemMarker
	children []ast.Block
}

type orderedListItemMarker struct {
	original   string
	num        int
	terminator string
}

func newOrderedListItemMarker(original string) orderedListItemMarker {
	number := original[:len(original)-1]
	num, err := strconv.Atoi(number)
	if err != nil {
		num = 1
	}
	terminator := original[len(original)-1:]
	return orderedListItemMarker{original: original, num: num, terminator: terminator}
}

func (o orderedListItemMarker) String() string {
	return o.original
}

func NewListItemOrdered(external string, marker string, internal string, children ...ast.Block) *listItemOrdered {
	li := &listItemOrdered{external: external, marker: newOrderedListItemMarker(marker), internal: internal, children: append([]ast.Block{}, children...)}
	li.tracker = newTracker(li)
	return li
}

func (l *listItemOrdered) Type() ast.BlockType {
	return ListItemOrderedBlockType
}

func (l *listItemOrdered) Children() []ast.Block {
	return l.children
}

func (l *listItemOrdered) Markdown() string {
	// If there are no children, we still need to add the marker
	if len(l.children) == 0 {
		return l.external + l.marker.String() + l.internal
	}

	lines := make([]string, 0)
	for i, child := range l.children {
		md := child.Markdown()
		childLines := strings.Split(md, "\n")
		for j, line := range childLines {
			// If first line of first child we need the full marker
			if i == 0 && j == 0 {
				lines = append(lines, l.external+l.marker.String()+l.internal+line)
				continue
			}
			// If the line is empty, dont add any indentation
			if line == "" {
				lines = append(lines, "")
				continue
			}

			// Otherwise, add the correct amount of indentation
			lines = append(lines, strings.Repeat(" ", len(l.external)+len(l.marker.String())+len(l.internal))+line)
		}
	}
	return strings.Join(lines, "\n")
}

func (l *listItemOrdered) SameType(li listItem) bool {
	or, ok := ListItemOrdered(li)
	if !ok {
		return false
	}
	return or.marker.terminator == l.marker.terminator
}

func (l *listItemOrdered) Modified() bool {
	return l.tracker.Modified(l)
}
