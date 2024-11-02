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

package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var inputs = map[string]string{
	"empty": "",

	"frontmatter": `---
title: "Hello, World!"
---`,

	"frontmatter and tasks": `---
title: "Hello, World!"
---
This is some text

Some more text!

EVEN MOAR!@!@
`,
}

func TestDocument(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Document
	}{
		{
			name:  "empty",
			input: inputs["empty"],
			want:  Document{Contents: []byte("")},
		},
		{
			name:  "frontmatter",
			input: inputs["frontmatter"],
			want: Document{
				Metadata: map[string]interface{}{
					"title": "Hello, World!",
				},
				Contents: []byte(""),
			},
		},
		{
			name:  "frontmatter and content",
			input: inputs["frontmatter and tasks"],
			want: Document{
				Metadata: map[string]interface{}{"title": "Hello, World!"},
				Contents: []byte("This is some text\n\nSome more text!\n\nEVEN MOAR!@!@\n"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDocument()(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
