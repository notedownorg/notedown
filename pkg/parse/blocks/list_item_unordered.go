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
	"strings"

	"github.com/notedownorg/notedown/pkg/parse/ast"
)

func ListItemUnordered(li ast.Block) (*listItemUnordered, bool) {
	l, ok := li.(*listItemUnordered)
	return l, ok
}

var _ ast.Block = &listItemUnordered{}
var _ listItem = &listItemUnordered{}

const ListItemUnorderedBlockType = "list_item_unordered"

type listItemUnordered struct {
	*tracker
	external string
	internal string
	marker   bulletListItemMarker
	children []ast.Block
}

type bulletListItemMarker int

const (
	asterisk bulletListItemMarker = iota
	plus
	minus
)

func NewListItemUnordered(external string, marker string, internal string, children ...ast.Block) *listItemUnordered {
	li := &listItemUnordered{external: external, marker: newBulletListItemMarker(marker), internal: internal, children: append([]ast.Block{}, children...)}
	li.tracker = newTracker(li)
	return li
}

func (l *listItemUnordered) Type() ast.BlockType {
	return ListItemUnorderedBlockType
}

func (l *listItemUnordered) Children() []ast.Block {
	return l.children
}

func (l *listItemUnordered) Markdown() string {

	// If there are no children, we still need to add the marker
	if len(l.children) == 0 {
		return l.external + marker(l.marker) + l.internal
	}

	lines := make([]string, 0)
	for i, child := range l.children {
		md := child.Markdown()
		childLines := strings.Split(md, "\n")
		for j, line := range childLines {
			// If first line of first child we need the full marker
			if i == 0 && j == 0 {
				lines = append(lines, l.external+marker(l.marker)+l.internal+line)
				continue
			}
			// If the line is empty, dont add any indentation
			if line == "" {
				lines = append(lines, "")
				continue
			}

			// Otherwise, add the correct amount of indentation
			lines = append(lines, strings.Repeat(" ", len(l.external)+len(marker(l.marker))+len(l.internal))+line)
		}
	}
	return strings.Join(lines, "\n")
}

func (l *listItemUnordered) SameType(li listItem) bool {
	unor, ok := ListItemUnordered(li)
	if !ok {
		return false
	}
	return l.marker == unor.marker
}

func (l *listItemUnordered) Modified() bool {
	return l.tracker.Modified(l)
}

func newBulletListItemMarker(s string) bulletListItemMarker {
	switch s {
	case "*":
		return asterisk
	case "+":
		return plus
	case "-":
		return minus
	}
	return minus
}

func marker(b bulletListItemMarker) string {
	switch b {
	case asterisk:
		return "*"
	case plus:
		return "+"
	case minus:
		return "-"
	}
	return "-"
}
