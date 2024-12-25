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

var _ ast.Block = &CodeBlockFenced{}

const CodeBlockFencedType ast.BlockType = "code_block_fenced"

type CodeBlockFenced struct {
	*tracker
	open       string // includes indent
	infostring string
	code       string
	close      string // includes indent
}

func NewCodeBlockFenced(open, infostring, code, close string) *CodeBlockFenced {
	cb := &CodeBlockFenced{open: open, infostring: infostring, code: code, close: close}
	cb.tracker = newTracker(cb)
	return cb
}

func (b *CodeBlockFenced) Type() ast.BlockType {
	return CodeBlockFencedType
}

func (b *CodeBlockFenced) Children() []ast.Block {
	return []ast.Block{} // Code blocks are leaf nodes
}

func (b *CodeBlockFenced) Markdown() string {
	var s strings.Builder
	s.WriteString(b.open)
	s.WriteString(b.infostring)
	// code and close could be empty depending on where we encounter EOF
	if b.code != "" {
		s.WriteString("\n")
		s.WriteString(b.code)
	}
	if b.close != "" {
		s.WriteString("\n")
		s.WriteString(b.close)
	}
	return s.String()
}

func (b *CodeBlockFenced) Modified() bool {
	return b.tracker.Modified(b)
}
