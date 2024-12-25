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

package workspace

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/notedownorg/notedown/pkg/parse/test"
	"github.com/stretchr/testify/assert"
)

var inputs = map[string]string{
	"empty": "",

	"frontmatter": `---
title: "Hello, World!"
---`,

	"frontmatter and words": `---
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
			want:  Document{Blocks: Bl()},
		},
		{
			name:  "frontmatter",
			input: inputs["frontmatter"],
			want: Document{
				Metadata: map[string]interface{}{
					"title": "Hello, World!",
				},
				Blocks: Bl(),
			},
		},
		{
			name:  "frontmatter and content",
			input: inputs["frontmatter and words"],
			want: Document{
				Metadata: map[string]interface{}{"title": "Hello, World!"},
				Blocks:   Bl(P("This is some text"), Bln, P("Some more text!"), Bln, P("EVEN MOAR!@!@")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := os.TempDir()
			filename := uuid.NewString() + ".md"
			fullPath := filepath.Join(tmpDir, filename)
			os.WriteFile(fullPath, []byte(tt.input), 0644)
			defer os.Remove(fullPath)
			got, err := LoadDocument(tmpDir, filename, time.Now())
			assert.NoError(t, err)
			assert.Equal(t, tt.want.Metadata, got.Metadata)
			assert.Equal(t, tt.want.Blocks, got.Blocks)
			assert.NotZero(t, got.lastModified)
			assert.NotZero(t, got.creationHash)
		})
	}
}
