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
	. "github.com/liamawhite/parse/core"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	"sigs.k8s.io/yaml"
)

var (
	frontMatterFence = String("---")

	frontMatterOpen  = SequenceOf3(frontMatterFence, OptionalInlineWhitespace, NewLine)
	frontMatterClose = SequenceOf4(NewLine, frontMatterFence, Optional(InlineWhitespace), Optional(NewLine))

	emptyFrontMatter = SequenceOf5(frontMatterFence, AtLeast(1, SequenceOf2(Optional(InlineWhitespace), NewLine)), frontMatterFence, OptionalInlineWhitespace, Optional(NewLine))
)

func FrontmatterParser(in Input) (ast.Block, bool, error) {
	start := in.Checkpoint()

	// Front matter can only be at the start of a document so return false if we're not at the start
	if int(start) != 0 {
		return nil, false, nil
	}

	// Check if we have an empty or whitespace only front frontmatter
	_, found, err := emptyFrontMatter(in)
	if err != nil {
		return nil, false, err
	}
	if found {
		return NewFrontmatter(nil), true, nil
	}

	// Otherwise check for a full front matter
	if _, found, err := frontMatterOpen(in); err != nil || !found {
		return nil, false, err
	}

	// Read the contents until the close front matter
	contents, found, err := StringWhileNot(frontMatterClose)(in)
	if err != nil || !found {
		in.Restore(start)
		return nil, false, err
	}

	// Parse the contents as yaml
	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(contents), &metadata); err != nil {
		in.Restore(start)
		return nil, false, nil
	}

	// Pop off the close front matter
	frontMatterClose(in)

	return NewFrontmatter(metadata), true, nil
}
