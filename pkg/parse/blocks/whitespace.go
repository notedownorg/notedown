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

// Not typically considered a block but we need it to be able to roundtrip
// This isn't always used, sometimes we just include the newline in the block content (e.g. paragraph)
var _ ast.Block = &BlankLine{}

const BlankLineBlockType ast.BlockType = "blank_line"

type BlankLine struct{}

func NewBlankLine() *BlankLine {
	return &BlankLine{}
}

func (b *BlankLine) Type() ast.BlockType {
	return BlankLineBlockType
}

func (b *BlankLine) Children() []ast.Block {
	return []ast.Block{} // Blank lines are leaf nodes
}

func (b *BlankLine) Markdown() string {
	return ""
}

func (b *BlankLine) Modified() bool {
	return false // Blank lines cannot be modified, only removed
}
