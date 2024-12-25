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

	"github.com/notedownorg/notedown/pkg/parse/ast"
)

var _ ast.Block = &HeadingAtx{}

const HeadingAtxBlockType ast.BlockType = "heading_atx"

type HeadingAtx struct {
	*tracker
	indent   int
	level    int
	title    string
	children []ast.Block
}

func NewHeadingAtx(indent, level int, title string, children ...ast.Block) *HeadingAtx {
	h := &HeadingAtx{indent: indent, level: level, title: title, children: children}
	h.tracker = newTracker(h)
	return h
}

func (b *HeadingAtx) Type() ast.BlockType {
	return HeadingAtxBlockType
}

func (b *HeadingAtx) Children() []ast.Block {
	return b.children
}

func (b *HeadingAtx) Markdown() string {
	var s strings.Builder
	s.WriteString(strings.Repeat(" ", b.indent))
	s.WriteString(strings.Repeat("#", b.level))
	s.WriteString(b.title)
	for _, child := range b.children {
		s.WriteString("\n")
		s.WriteString(child.Markdown())
	}
	return s.String()
}

func (b *HeadingAtx) Modified() bool {
	return b.tracker.Modified(b)
}
