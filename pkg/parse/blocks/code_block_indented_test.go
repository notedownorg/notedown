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

func TestCodeBlockIndented(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "Indented code block: basic code block",
			Input:         "    a simple\n      indented code block\n",
			Parser:        CodeBlockIndentedParser,
			ExpectedMatch: Cbi("    a simple\n      indented code block"),
			ExpectedOK:    true,
		},
		{
			Name:          "Indented code block: basic code block with tabs",
			Input:         "\ta simple\n\t  indented code block\n",
			Parser:        CodeBlockIndentedParser,
			ExpectedMatch: Cbi("\ta simple\n\t  indented code block"),
			ExpectedOK:    true,
		},
		{
			Name:          "Indented code block: handles blank lines",
			Input:         "    a simple\n    \n\n\n      indented code block\n",
			Parser:        CodeBlockIndentedParser,
			ExpectedMatch: Cbi("    a simple\n    \n\n\n      indented code block"),
			ExpectedOK:    true,
		},
		{
			Name:          "Indented code block: trailing whitespace",
			Input:         "    a simple\n      indented code block   \n",
			Parser:        CodeBlockIndentedParser,
			ExpectedMatch: Cbi("    a simple\n      indented code block   "),
			ExpectedOK:    true,
		},
	}
	RunTests(t, tests)
}
