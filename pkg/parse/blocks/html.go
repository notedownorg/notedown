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

var _ ast.Block = &Html{}

const HtmlBlockType ast.BlockType = "html"

type HtmlKind int

const (
	HtmlOne HtmlKind = iota + 1
	HtmlTwo
	HtmlThree
	HtmlFour
	HtmlFive
	HtmlSix
	HtmlSeven
)

type Html struct {
	*tracker
	kind    HtmlKind
	content string
}

func NewHtml(kind HtmlKind, content string) *Html {
	html := &Html{kind: kind, content: content}
	html.tracker = newTracker(html)
	return html
}

func (b *Html) Type() ast.BlockType {
	return HtmlBlockType
}

func (b *Html) Children() []ast.Block {
	return []ast.Block{} // HTML blocks are leaf nodes
}

func (b *Html) Markdown() string {
	return b.content
}

func (b *Html) Modified() bool {
	return b.tracker.Modified(b)
}
