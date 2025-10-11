// Copyright 2025 Notedown Authors
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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrontmatterParsing(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		input    string
		expected map[string]any
	}{
		{
			name: "YAML frontmatter",
			input: `---
title: "Test Document"
author: "John Doe"
tags: ["test", "markdown"]
published: true
---

# Heading

Some content here.
`,
			expected: map[string]any{
				"title":     "Test Document",
				"author":    "John Doe",
				"tags":      []any{"test", "markdown"},
				"published": true,
			},
		},
		{
			name: "empty frontmatter",
			input: `---
---

# Heading

Some content here.
`,
			expected: map[string]any{},
		},
		{
			name: "no frontmatter",
			input: `# Heading

Some content here.
`,
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.ParseString(tt.input)
			require.NoError(t, err)
			require.NotNil(t, doc)

			assert.Equal(t, tt.expected, doc.Metadata)
		})
	}
}
