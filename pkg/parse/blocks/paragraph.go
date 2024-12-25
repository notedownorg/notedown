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
	"github.com/notedownorg/notedown/pkg/parse/ast"
)

func Paragraph(b ast.Block) (*paragraph, bool) {
	p, ok := b.(*paragraph)
	return p, ok
}

var _ ast.Block = &paragraph{}

const ParagraphBlockType ast.BlockType = "paragraph"

type paragraph struct {
	*tracker
	text string
}

func NewParagraph(text string) *paragraph {
	p := &paragraph{text: text}
	p.tracker = newTracker(p)
	return p
}

func (p *paragraph) Type() ast.BlockType {
	return ParagraphBlockType
}

func (p *paragraph) Children() []ast.Block {
	return []ast.Block{} // paragraph is a leaf block so cannot have children
}

func (p *paragraph) Markdown() string {
	return p.text
}

func (p *paragraph) Modified() bool {
	return p.tracker.Modified(p)
}
