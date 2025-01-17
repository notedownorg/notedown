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
	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
)

func ListParser(ctx Context) Parser[ast.Block] {
	return func(in Input) (ast.Block, bool, error) {
		start := in.Checkpoint()

		// Attempt to read the first list item as this determines if we are ordered/unordered and the list marker
		initial, found, err := ListItemParser(ctx)(in)
		if err != nil || !found {
			return nil, false, nil
		}
		first, ok := initial.(listItem)
		if !ok {
			return nil, false, nil // not sure this would ever happen but just in case
		}

		// Now keep reading list items until we can't find any more,
		// or we come across a thematic break with the same bullet
		// or the list type/marker changes
		listItems := []ast.Block{first}
		for {
			iteration := in.Checkpoint()

			// Thematic breaks that have the same bullet close the list
			tbBlock, found, err := ThematicBreakParser(in)
			if err != nil {
				in.Restore(start)
				return nil, false, err
			}
			if found {
				in.Restore(iteration)
				tb, _ := tbBlock.(*thematicBreak)
				if first.SameType(NewListItemUnordered("", tb.char, "")) {
					break
				}
			}

			next, found, err := ListItemParser(ctx)(in)
			if err != nil {
				in.Restore(start)
				return nil, false, err
			}
			if !found {
				// As lists can have as many blank lines as they want between items but empty list items cant
				// we need to read off any blank lines and continue if we find any
				if ch, _ := in.Peek(1); ch == "\n" {
					in.Take(1)
					listItems = append(listItems, NewBlankLine())
					continue
				}
				break
			}

			// Check that the list item is of the same type as the first list item
			if m, ok := next.(listItem); !ok || !m.SameType(first) {
				in.Restore(iteration)
				break
			}

			listItems = append(listItems, next)
		}

		// Build the list based on the first list item
		if li, ok := ListItemOrdered(first); ok {
			return NewOrderedList(li.marker.num, listItems...), true, nil
		}
		if _, ok := first.(*ListItemTask); ok {
			return NewTaskList(listItems...), true, nil
		}
		return NewUnorderedList(listItems...), true, nil
	}
}
