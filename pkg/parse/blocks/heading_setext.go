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

var _ ast.Block = &HeadingSetext{}

const HeadingSetextBlockType ast.BlockType = "heading_setext"

type HeadingSetext struct {
	*tracker
	level     int
	title     string
	underline string
	children  []ast.Block
}

func NewHeadingSetext(title string, underline string, children ...ast.Block) *HeadingSetext {
	level := 1
	if strings.Contains(underline, "-") {
		level = 2
	}
	h := &HeadingSetext{level: level, title: title, underline: underline, children: children}
	h.tracker = newTracker(h)
	return h
}

func (b *HeadingSetext) Type() ast.BlockType {
	return HeadingSetextBlockType
}

func (b *HeadingSetext) Children() []ast.Block {
	return b.children
}

func (b *HeadingSetext) Markdown() string {
	var s strings.Builder
	s.WriteString(b.title)
	s.WriteString("\n")
	s.WriteString(b.underline)
	for _, child := range b.children {
		s.WriteString("\n")
		s.WriteString(child.Markdown())
	}
	return s.String()
}

func (b *HeadingSetext) Modified() bool {
	return b.tracker.Modified(b)
}
