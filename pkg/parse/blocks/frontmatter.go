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
	"sigs.k8s.io/yaml"
)

func Frontmatter(fm ast.Block) (*frontMatter, bool) {
	f, ok := fm.(*frontMatter)
	return f, ok
}

var _ ast.Block = &frontMatter{}

const FrontMatterBlockType = "frontmatter"

type frontMatter struct {
	*tracker
	metadata map[string]interface{}
}

func NewFrontmatter(metadata map[string]interface{}) *frontMatter {
	fm := &frontMatter{metadata: metadata}
	fm.tracker = newTracker(fm)
	return fm
}

func (f *frontMatter) Type() ast.BlockType {
	return FrontMatterBlockType
}

func (f *frontMatter) Children() []ast.Block {
	return []ast.Block{} // Front matter has no children
}

func (f *frontMatter) Markdown() string {
	var s strings.Builder
	s.WriteString("---\n")
	if f.metadata != nil && len(f.metadata) > 0 {
		content, _ := yaml.Marshal(f.metadata)
		s.WriteString(string(content))
	}
	s.WriteString("---")
	return s.String()
}

func (f *frontMatter) Metadata() map[string]interface{} {
	return f.metadata
}

func (f *frontMatter) Modified() bool {
	return f.tracker.Modified(f)
}
