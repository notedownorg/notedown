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

package parsers

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/notedownorg/notedown/pkg/ast"
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
- [ ] Task 1 due:2021-01-01
- [/] Task 2
`,

	"tasks with interleaved text": `- [ ] Task 1
This is some text
- [/] Task 2
This is some more text`,

	"lots of newlines": `

- [ ] Task 1

- [/] Task 2 due:2021-01-01



This is some text`,
}

func TestDocumentParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     ast.Document
		notFound bool
	}{
		{
			name:  "empty",
			input: inputs["empty"],
			want:  ast.Document{},
		},
		{
			name:  "frontmatter",
			input: inputs["frontmatter"],
			want: ast.Document{
				Metadata: map[string]interface{}{
					"title": "Hello, World!",
				},
				Markers: ast.Markers{ContentStart: 3},
			},
		},
		{
			name:  "frontmatter and tasks",
			input: inputs["frontmatter and tasks"],
			want: ast.Document{
				Metadata: map[string]interface{}{"title": "Hello, World!"},
				Tasks: []ast.Task{
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 1", ast.Todo, ast.WithDue(date(2021, 1, 1)), ast.WithLine(3)),
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 2", ast.Doing, ast.WithLine(4)),
				},
				Markers: ast.Markers{ContentStart: 3},
			},
		},
		{
			name:  "tasks with interleaved text",
			input: inputs["tasks with interleaved text"],
			want: ast.Document{
				Tasks: []ast.Task{
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 1", ast.Todo, ast.WithLine(0)),
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 2", ast.Doing, ast.WithLine(2)),
				},
			},
		},
		{
			name:  "lots of newlines",
			input: inputs["lots of newlines"],
			want: ast.Document{
				Tasks: []ast.Task{
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 1", ast.Todo, ast.WithLine(2)),
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 2", ast.Doing, ast.WithDue(date(2021, 1, 1)), ast.WithLine(4)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			got, found, _ := DocumentParser("p", "v", relativeTo).Parse(input)
			if tt.notFound {
				if found {
					t.Fatalf("expected not found, got %v", got)
				}
				return
			}
			assert.Equal(t, tt.want, got)

			rem, _, _ := parse.StringUntil(parse.EOF[string]()).Parse(input)
			assert.Equal(t, "", rem, "expected input to be consumed")
		})
	}
}

func TestDocument(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ast.Document
	}{
		{
			name:  "empty",
			input: inputs["empty"],
			want:  ast.Document{},
		},
		{
			name:  "frontmatter",
			input: inputs["frontmatter"],
			want: ast.Document{
				Metadata: map[string]interface{}{
					"title": "Hello, World!",
				},
				Markers: ast.Markers{ContentStart: 3},
			},
		},
		{
			name:  "frontmatter and tasks",
			input: inputs["frontmatter and tasks"],
			want: ast.Document{
				Metadata: map[string]interface{}{"title": "Hello, World!"},
				Tasks: []ast.Task{
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 1", ast.Todo, ast.WithLine(3), ast.WithDue(date(2021, 1, 1))),
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 2", ast.Doing, ast.WithLine(4)),
				},
				Markers: ast.Markers{ContentStart: 3},
			},
		},
		{
			name:  "tasks with interleaved text",
			input: inputs["tasks with interleaved text"],
			want: ast.Document{
				Tasks: []ast.Task{
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 1", ast.Todo, ast.WithLine(0)),
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 2", ast.Doing, ast.WithLine(2)),
				},
			},
		},
		{
			name:  "lots of newlines",
			input: inputs["lots of newlines"],
			want: ast.Document{
				Tasks: []ast.Task{
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 1", ast.Todo, ast.WithLine(2)),
					ast.NewTask(ast.NewIdentifier("p", "v"), "Task 2", ast.Doing, ast.WithDue(date(2021, 1, 1)), ast.WithLine(4)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Document("p", "v", relativeTo)(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
