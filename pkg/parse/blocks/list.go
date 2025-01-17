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

var _ ast.Block = &list{}

const ListOrderedBlockType = "list_ordered"
const ListUnorderedBlockType = "list_unordered"
const ListTaskBlockType = "list_task"

type kind int

const (
	ordered kind = iota
	unordered
	task
)

type list struct {
	*tracker
	kind     kind
	start    int // only used if kind is ordered
	children []ast.Block
}

func NewOrderedList(start int, children ...ast.Block) *list {
	list := &list{kind: ordered, start: start, children: append([]ast.Block{}, children...)}
	list.tracker = newTracker(list)
	return list
}

func NewUnorderedList(children ...ast.Block) *list {
	return &list{kind: unordered, children: append([]ast.Block{}, children...)}
}

func NewTaskList(children ...ast.Block) *list {
	return &list{kind: task, children: append([]ast.Block{}, children...)}
}

func (l *list) Type() ast.BlockType {
	switch l.kind {
	case ordered:
		return ListOrderedBlockType
	case task:
		return ListTaskBlockType
	default:
		return ListUnorderedBlockType
	}
}

func (l *list) Children() []ast.Block {
	return l.children
}

func (l *list) Markdown() string {
	var s strings.Builder
	for i, child := range l.children {
		s.WriteString(child.Markdown())
		if i < len(l.children)-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

func (l *list) Modified() bool {
	return false
}

type listItem interface {
	ast.Block
	SameType(listItem) bool
}
