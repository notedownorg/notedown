package parsers

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/ast"
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
					{Name: "Task 1", Status: ast.Todo, Line: 3, Due: date(2021, 1, 1)},
					{Name: "Task 2", Status: ast.Doing, Line: 4},
				},
				Markers: ast.Markers{ContentStart: 3},
			},
		},
		{
			name:  "tasks with interleaved text",
			input: inputs["tasks with interleaved text"],
			want: ast.Document{
				Tasks: []ast.Task{
					{Name: "Task 1", Status: ast.Todo, Line: 0},
					{Name: "Task 2", Status: ast.Doing, Line: 2},
				},
			},
		},
		{
			name:  "lots of newlines",
			input: inputs["lots of newlines"],
			want: ast.Document{
				Tasks: []ast.Task{
					{Name: "Task 1", Status: ast.Todo, Line: 2},
					{Name: "Task 2", Status: ast.Doing, Line: 4, Due: date(2021, 1, 1)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			got, found, _ := DocumentParser(relativeTo).Parse(input)
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
			want:  ast.Document{Hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		},
		{
			name:  "frontmatter",
			input: inputs["frontmatter"],
			want: ast.Document{
				Metadata: map[string]interface{}{
					"title": "Hello, World!",
				},
				Markers: ast.Markers{ContentStart: 3},
				Hash:    "e1e3896328c141650e534fae42b886c0ea332f57b48b90492f38ef269596a026",
			},
		},
		{
			name:  "frontmatter and tasks",
			input: inputs["frontmatter and tasks"],
			want: ast.Document{
				Metadata: map[string]interface{}{"title": "Hello, World!"},
				Tasks: []ast.Task{
					{Name: "Task 1", Status: ast.Todo, Line: 3, Due: date(2021, 1, 1)},
					{Name: "Task 2", Status: ast.Doing, Line: 4},
				},
				Markers: ast.Markers{ContentStart: 3},
				Hash:    "1ac564f7760bf58243fb4967ca5aaedb1122c3dbcb2006ed877aa3518d643850",
			},
		},
		{
			name:  "tasks with interleaved text",
			input: inputs["tasks with interleaved text"],
			want: ast.Document{
				Tasks: []ast.Task{
					{Name: "Task 1", Status: ast.Todo, Line: 0},
					{Name: "Task 2", Status: ast.Doing, Line: 2},
				},
				Hash: "1b48378d61de419fde23e6ed991e81c759bd223526dd21bb667a2b273d20e637",
			},
		},
		{
			name:  "lots of newlines",
			input: inputs["lots of newlines"],
			want: ast.Document{
				Tasks: []ast.Task{
					{Name: "Task 1", Status: ast.Todo, Line: 2},
					{Name: "Task 2", Status: ast.Doing, Line: 4, Due: date(2021, 1, 1)},
				},
				Hash: "eabda4800e028ae2ab063083ea2e33e25a1583f12825f15be9080349fd8c1e2a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Document(relativeTo)(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
