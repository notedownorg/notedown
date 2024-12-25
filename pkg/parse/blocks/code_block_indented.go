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

var _ ast.Block = &CodeBlockIndented{}

const CodeBlockIndentedType ast.BlockType = "code_block_indented"

type CodeBlockIndented struct {
	*tracker
	code []string
}

func NewCodeBlockIndented(code []string) *CodeBlockIndented {
	cb := &CodeBlockIndented{code: code}
	cb.tracker = newTracker(cb)
	return cb
}

func (b *CodeBlockIndented) Type() ast.BlockType {
	return CodeBlockIndentedType
}

func (b *CodeBlockIndented) Children() []ast.Block {
	return []ast.Block{} // Indented code blocks are leaf nodes
}

func (b *CodeBlockIndented) Markdown() string {
	var s strings.Builder
	for i, line := range b.code {
		s.WriteString(line)
		if i < len(b.code)-1 {
			s.WriteString("\n")
		}
	}
	return s.String()
}

func (b *CodeBlockIndented) Modified() bool {
	return b.tracker.Modified(b)
}
