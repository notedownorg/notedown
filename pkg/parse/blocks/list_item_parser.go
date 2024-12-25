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

	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/traversal"
)

func newListItem(ctx Context) func(external string, marker string, internal string, children ...ast.Block) ast.Block {
	return func(external string, marker string, internal string, children ...ast.Block) ast.Block {
		if marker == "*" || marker == "-" || marker == "+" {
			// If the first child is a paragraph we can potentially have a task
			// So attempt to parse the first child as a task
			if len(children) > 0 {
				if p, ok := Paragraph(children[0]); ok {
					in := NewInput(p.text)
					taskBuilder, found, _ := listItemTaskParser(ctx)(in)
					if found {
						// Consume the rest of the original paragraph and set it as child 0 if there is anything left
						// Otherwise just remove the first child
						remaining, _, _ := StringUntil(EOF[string]())(in)
						if remaining != "" {
							children[0] = NewParagraph(remaining)
						} else {
							children = children[1:]
						}
						opts := append(taskBuilder.opts, TaskWithChildren(children...))
						return NewListItemTask(external, marker, internal, taskBuilder.status, taskBuilder.text, opts...)
					}
				}
			}

			return NewListItemUnordered(external, marker, internal, children...)
		}
		return NewListItemOrdered(external, marker, internal, children...)
	}
}

func ListItemParser(ctx Context) Parser[ast.Block] {
	return func(in Input) (ast.Block, bool, error) {
		start := in.Checkpoint()

		// Can this list item be interpreted as a thematic break? If so the thematic break parser takes precedence.
		// This can't be done in the top level parser as the correct ordering required is not possible due to other constraints.
		_, found, err := ThematicBreakParser(in)
		if err != nil || found {
			in.Restore(start)
			return nil, false, err
		}

		opener, found, err := Or(emptyListItemOpenParser, listItemOpenParser)(in)
		if err != nil || !found {
			return nil, false, err
		}

		// If the opener is empty and we have a newline next, we're done as you can have at most one empty line to open
		empty, notempty := opener.Values()
		if empty.Ok() {
			if ch, _ := in.Peek(1); ch == "\n" {
				o := empty.Values()
				var child ast.Block = &BlankLine{}
				if o.initialContent != "" {
					child = NewParagraph(o.initialContent)
				}
				return newListItem(ctx)(o.indentation, o.marker.String(), o.internal, child), true, nil
			}
		}

		// Otherwise we need to fill in the list item then handle children
		external, marker, internal, initialContent := "", "", "", ""
		if empty.Ok() {
			o := empty.Values()
			external, marker, internal, initialContent = o.indentation, o.marker.String(), o.internal, o.initialContent
		} else {
			// If we're not empty
			o := notempty.Values()
			external, marker, internal, initialContent = o.indentation, o.marker.String(), o.internal, o.initialContent
		}

		requiredIndentation := len(external) + len(marker) + spaces(internal)

		// Put the content from the first line into the string builder
		var s strings.Builder
		s.WriteString(initialContent)

		// Read the rest of the list item line by line until we reach a line that doesn't have enough indentation
		// Strip the indentation from each line before passing the input to the children
		// Otherwise we would have to make every parser aware of the list item indentation and that would be a mess
		for {
			// If we are at the end of the file, break
			if _, found, _ := EOF[bool]()(in); found {
				break
			}

			// If we're not at the end of the file, pop the newline and add one to the string
			NewLine(in)
			s.WriteString("\n")

			// Now we can actually parse the line
			iteration := in.Checkpoint()
			line, _, err := SequenceOf2(StringFrom(ZeroOrMore(RuneIn(" \t"))), StringWhileNotEOFOr(NewLine))(in)
			if err != nil {
				in.Restore(start)
				return nil, false, err
			}
			indentation, content := line.Values()

			// If there isnt enough indentation we're done unless we have a blank line
			// We need to restore the input to the start of the line so the next parser can pick up where we left off
			if spaces(indentation) < requiredIndentation && indentation+content != "" {
				in.Restore(iteration)
				break
			}

			// Make sure we persist the excess indentation for things like indented code blocks
			s.WriteString(strings.TrimPrefix(indentation, strings.Repeat(" ", requiredIndentation)))
			s.WriteString(content)
		}

		// Build a new input from the string we've built and parse the children
		liInput := NewInput(s.String())

		children, _, err := DepthFirstParser(TopLevelParser(ctx), EOF[ast.Block]())(liInput)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}

		return newListItem(ctx)(external, marker, internal, children...), true, nil
	}
}

type openListItem struct {
	indentation    string
	marker         listItemMarker
	internal       string
	initialContent string
}

// Empty list items refer to list items that start with a blank or whitespace only line
// There can only be a single blank line (on the initial line), anything beyond that is considered a new block
func emptyListItemOpenParser(in Input) (openListItem, bool, error) {
	start := in.Checkpoint()

	prefix, found, err := SequenceOf2(indent, listMarkerParser)(in)
	if err != nil || !found {
		return openListItem{}, false, err
	}
	indentation, marker := prefix.Values()

	// Line is empty if there is only whitespace between the marker and the newline or EOF
	rem, found, err := SequenceOf2(OptionalInlineWhitespace, StringFrom(Or(NewLine, EOF[string]())))(in)
	if err != nil || !found {
		in.Restore(start)
		return openListItem{}, false, err
	}

	inline, _ := rem.Values()
	internal, initialContent := " ", ""
	if spaces(inline) > 0 {
		initialContent = strings.Repeat(" ", spaces(inline)-1)
	}

	return openListItem{indentation: indentation, marker: marker, internal: internal, initialContent: initialContent}, true, nil
}

// Standard list item
func listItemOpenParser(in Input) (openListItem, bool, error) {
	// Ensure that we are even at the start of a list item and stop short if not
	first, found, err := SequenceOf4(indent, listMarkerParser, StringFrom(AtLeast(1, RuneIn(" \t"))), StringWhileNotEOFOr(NewLine))(in)
	if err != nil || !found {
		return openListItem{}, false, nil
	}

	indentation, marker, internal, initialContent := first.Values()

	// If the internal indentation is large enough to be the start of an indented code block we keep one space and add the
	// rest to the front of the content. Ensuring that we still parse the indented code block correctly.
	if spaces(internal) >= 5 { // 4 spaces for code block + 1 space for list item
		initialContent = strings.Repeat(" ", spaces(internal)-1) + initialContent
		internal = " "
	}

	return openListItem{indentation: indentation, marker: marker, internal: internal, initialContent: initialContent}, true, nil
}
