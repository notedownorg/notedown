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

package blocks_test

import (
	"testing"

	. "github.com/liamawhite/parse/test"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/blocks"
	. "github.com/notedownorg/notedown/pkg/parse/test"
)

func TestFrontMatter(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "valid frontmatter with trailing newline",
			Input:         "---\ntitle: \"Hello, World!\"\n---\n",
			Parser:        FrontmatterParser,
			ExpectedMatch: Fm(map[string]interface{}{"title": "Hello, World!"}),
			ExpectedOK:    true,
		},
		{
			Name:          "valid frontmatter with EOF",
			Input:         "---\ntitle: \"Hello, World!\"\n---",
			Parser:        FrontmatterParser,
			ExpectedMatch: Fm(map[string]interface{}{"title": "Hello, World!"}),
			ExpectedOK:    true,
		},
		{
			Name:           "valid frontmatter with double newline",
			Input:          "---\ntitle: \"Hello, World!\"\n---\n\n",
			Parser:         FrontmatterParser,
			ExpectedMatch:  Fm(map[string]interface{}{"title": "Hello, World!"}),
			ExpectedOK:     true,
			RemainingInput: "\n",
		},
		{
			Name:           "invalid yaml in frontmatter",
			Input:          "---\nnope\n---\n",
			Parser:         FrontmatterParser,
			RemainingInput: "---\nnope\n---\n",
		},
		{
			Name:           "not at start of input",
			Input:          `x---\ntitle: "Hello, World!"\n---\n`,
			Parser:         FrontmatterParser,
			RemainingInput: `x---\ntitle: "Hello, World!"\n---\n`,
		},
		{
			Name:           "no frontmatter",
			Input:          `# Hello, World!`,
			Parser:         FrontmatterParser,
			RemainingInput: `# Hello, World!`,
		},
		{
			Name:          "empty frontmatter",
			Input:         "---\n---\n",
			Parser:        FrontmatterParser,
			ExpectedMatch: Fm(nil),
			ExpectedOK:    true,
		},
		{
			Name:          "empty frontmatter with newlines content",
			Input:         "---\n\n---\n",
			Parser:        FrontmatterParser,
			ExpectedMatch: Fm(nil),
			ExpectedOK:    true,
		},
		{
			Name:          "empty frontmatter with whitespace content",
			Input:         "---\n  \t  \n---\n",
			Parser:        FrontmatterParser,
			ExpectedMatch: Fm(nil),
			ExpectedOK:    true,
		},
		{
			Name:          "frontmatter yaml with leading and trailing newlines",
			Input:         "---\n\ntitle: \"Hello, World!\"\n\n---\n",
			Parser:        FrontmatterParser,
			ExpectedMatch: Fm(map[string]interface{}{"title": "Hello, World!"}),
			ExpectedOK:    true,
		},
		{
			Name:           "frontmatter yaml without newline before closing fence",
			Input:          "---\ntitle: \"Hello, World!\"---\n",
			Parser:         FrontmatterParser,
			ExpectedOK:     false,
			RemainingInput: "---\ntitle: \"Hello, World!\"---\n",
		},
	}
	RunTests(t, tests)
}
