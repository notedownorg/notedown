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

package parsers_test

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/parsers"
	"github.com/stretchr/testify/assert"
)

func TestFrontMatter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected parsers.FrontMatter
		notFound bool
	}{
		{
			name: "valid frontmatter",
			input: `---
title: "Hello, World!"
---`,
			expected: parsers.FrontMatter(`title: "Hello, World!"`),
		},
		{
			name: "invalid yaml in frontmatter",
			input: `---
title:
Hello, World!
---`,
			notFound: true,
		},
		{
			name:     "no frontmatter",
			input:    `# Hello, World!`,
			notFound: true,
		},
		{
			name: "empty frontmatter",
			input: `---
---`,
			expected: parsers.FrontMatter(""),
		},
		{
			name: "empty frontmatter with whitespace",
			input: `---
      
---`,
			expected: parsers.FrontMatter("      "), // there are 6 spaces in the input
		},
		{
			name: "empty frontmatter with newline",
			input: `---

---`,
			expected: parsers.FrontMatter(""),
		},
		{
			name: "frontmatter yaml with leading and trailing newlines",
			input: `---


title: "Hello, World!"


---`,
			expected: parsers.FrontMatter(`title: "Hello, World!"`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			fm, ok, _ := parsers.Frontmatter.Parse(in)

			if test.notFound {
				if ok {
					t.Fatalf("expected not found, content: %s", string(fm))
				}
				return
			}
			if !ok {
				t.Fatalf("expected found")
			}
			assert.Equal(t, string(test.expected), string(fm))
		})
	}
}
