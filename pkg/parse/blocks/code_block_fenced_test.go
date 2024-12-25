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

func TestCodeBlockFenced(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "Fenced code block: basic backticks",
			Input:         "```\ncodeblock\n```\n",
			Parser:        CodeBlockFencedParser,
			ExpectedMatch: Cbf("```", "", "codeblock", "```"),
			ExpectedOK:    true,
		},
		{
			Name:          "Fenced code block: basic tilde",
			Input:         "~~~\ncodeblock\n~~~\n",
			Parser:        CodeBlockFencedParser,
			ExpectedMatch: Cbf("~~~", "", "codeblock", "~~~"),
			ExpectedOK:    true,
		},
		{
			Name:          "Fenced code block: basic backticks with language",
			Input:         "```go\ncodeblock\n```\n",
			Parser:        CodeBlockFencedParser,
			ExpectedMatch: Cbf("```", "go", "codeblock", "```"),
			ExpectedOK:    true,
		},
		{
			Name:          "Fenced code block: real world example",
			Input:         "```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```\n",
			Parser:        CodeBlockFencedParser,
			ExpectedMatch: Cbf("```", "go", "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}", "```"),
			ExpectedOK:    true,
		},
	}
	RunTests(t, tests)
}
