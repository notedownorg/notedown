package parsers

import (
	"fmt"
	"testing"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     api.Document
		notFound bool
	}{
		{
			name:  "empty",
			input: "",
			want:  api.Document{},
		},
		{
			name: "frontmatter",
			input: `---
title: "Hello, World!"
---`,
			want: api.Document{
				Metadata: map[string]interface{}{
					"title": "Hello, World!",
				},
			},
		},
		{
			name: "frontmatter and tasks",
			input: `---
title: "Hello, World!"
---
- [ ] Task 1
- [/] Task 2
`,
			want: api.Document{
				Metadata: map[string]interface{}{"title": "Hello, World!"},
				Tasks: []api.Task{
					{Name: "Task 1", Status: api.Todo, Line: 3},
					{Name: "Task 2", Status: api.Doing, Line: 4},
				},
			},
		},
		{
			name: "tasks with interleaved text",
			input: `- [ ] Task 1
This is some text
- [/] Task 2
This is some more text`,
			want: api.Document{
				Tasks: []api.Task{
					{Name: "Task 1", Status: api.Todo, Line: 0},
					{Name: "Task 2", Status: api.Doing, Line: 2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			fmt.Println(len(tt.input))
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
