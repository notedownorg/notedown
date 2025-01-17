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
	"strings"

	"github.com/notedownorg/notedown/pkg/parse/ast"
)

var _ ast.Block = &blockQuote{}

const BlockQuoteBlockType = "blockquote"

type blockQuote struct {
	*tracker
	identation string
	parent     ast.Block
	children   []ast.Block
}

func NewBlockQuote(indentation string, children ...ast.Block) *blockQuote {
	bq := &blockQuote{identation: indentation, children: append([]ast.Block{}, children...)}
	bq.tracker = newTracker(bq)
	return bq
}

func (b *blockQuote) Type() ast.BlockType {
	return BlockQuoteBlockType
}

func (b *blockQuote) Parent() ast.Block {
	return b.parent
}

func (b *blockQuote) Children() []ast.Block {
	return b.children
}

// For each child, get the markdown and prepend every line with the identation and a >
// This violates the rule of not modifying the original document, but I couldn't find
// a way to avoid this without massively without massively increasing the complexity of the parser.
func (b *blockQuote) Markdown() string {
	lines := make([]string, 0)
	for _, child := range b.children {
		md := child.Markdown()
		childLines := strings.Split(md, "\n")
		for _, line := range childLines {
			lines = append(lines, fmt.Sprintf("%s> %s", b.identation, line))
		}
	}
	return strings.Join(lines, "\n")
}

func (b *blockQuote) Modified() bool {
	return b.tracker.Modified(b)
}
