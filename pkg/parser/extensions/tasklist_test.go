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

package extensions

import (
	"testing"

	"github.com/notedownorg/notedown/pkg/config"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestTaskCheckBox_NewTaskCheckBox(t *testing.T) {
	tests := []struct {
		name  string
		state string
	}{
		{
			name:  "unchecked checkbox",
			state: " ",
		},
		{
			name:  "checked checkbox",
			state: "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTaskCheckBox(tt.state)

			if node == nil {
				t.Fatal("Expected node to be created, got nil")
			}

			if node.State != tt.state {
				t.Errorf("Expected State = %q, got %q", tt.state, node.State)
			}

			if node.Kind() != KindTaskCheckBox {
				t.Errorf("Expected kind %v, got %v", KindTaskCheckBox, node.Kind())
			}
		})
	}
}

func TestTaskCheckBox_Dump(t *testing.T) {
	node := NewTaskCheckBox("x")

	// Test that Dump doesn't panic - we can't easily test the output
	// since it writes to internal goldmark structures
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump panicked: %v", r)
		}
	}()

	node.Dump([]byte("test source"), 0)
}

func TestKindTaskCheckBox(t *testing.T) {
	// Test that KindTaskCheckBox is properly initialized
	if KindTaskCheckBox.String() != "TaskCheckBox" {
		t.Errorf("Expected KindTaskCheckBox to be 'TaskCheckBox', got %q", KindTaskCheckBox.String())
	}
}

func TestTaskListParser_Trigger(t *testing.T) {
	parser := NewTaskListParser(config.GetDefaultConfig())
	triggers := parser.Trigger()

	if len(triggers) != 1 || triggers[0] != '[' {
		t.Errorf("Expected trigger '[', got %v", triggers)
	}
}

func TestTaskListParser_Parse(t *testing.T) {
	// Test basic parsing functionality with integration tests
	// This avoids the complex goldmark parser context API issues
	tests := []struct {
		name        string
		markdown    string
		wantChecked []bool
	}{
		{
			name:        "checked checkbox",
			markdown:    "- [x] Task",
			wantChecked: []bool{true},
		},
		{
			name:        "unchecked checkbox",
			markdown:    "- [ ] Task",
			wantChecked: []bool{false},
		},
		{
			name:        "uppercase X checkbox",
			markdown:    "- [X] Task",
			wantChecked: []bool{true},
		},
		{
			name:        "multiple checkboxes",
			markdown:    "- [x] Done\n- [ ] Todo",
			wantChecked: []bool{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension(config.GetDefaultConfig())))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var checkboxes []*TaskCheckBox
			_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if cb, ok := node.(*TaskCheckBox); ok {
						checkboxes = append(checkboxes, cb)
					}
				}
				return ast.WalkContinue, nil
			})

			if len(checkboxes) != len(tt.wantChecked) {
				t.Errorf("Expected %d checkboxes, got %d", len(tt.wantChecked), len(checkboxes))
				return
			}

			for i, cb := range checkboxes {
				isChecked := cb.State != " " && cb.State != ""
				if isChecked != tt.wantChecked[i] {
					t.Errorf("Checkbox %d: expected checked = %v, got %v (state: %q)", i, tt.wantChecked[i], isChecked, cb.State)
				}
			}
		})
	}
}

func TestTaskListParser_ParseWithExistingChildren(t *testing.T) {
	// Test basic negative cases through integration tests
	tests := []struct {
		name      string
		markdown  string
		wantCount int
	}{
		{
			name:      "not a checkbox in paragraph",
			markdown:  "[x] Not in list",
			wantCount: 0,
		},
		{
			name:      "invalid checkbox chars",
			markdown:  "- [y] Invalid",
			wantCount: 0,
		},
		{
			name:      "malformed checkbox",
			markdown:  "- [x Not closed",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension(config.GetDefaultConfig())))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var count int
			_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if _, ok := node.(*TaskCheckBox); ok {
						count++
					}
				}
				return ast.WalkContinue, nil
			})

			if count != tt.wantCount {
				t.Errorf("Expected %d checkboxes, got %d", tt.wantCount, count)
			}
		})
	}
}

func TestTaskListASTParsing(t *testing.T) {
	// Test that task checkboxes are correctly parsed into AST
	tests := []struct {
		name       string
		markdown   string
		wantCount  int
		wantStates []bool
	}{
		{
			name:       "single checked task",
			markdown:   "- [x] Task",
			wantCount:  1,
			wantStates: []bool{true},
		},
		{
			name:       "single unchecked task",
			markdown:   "- [ ] Task",
			wantCount:  1,
			wantStates: []bool{false},
		},
		{
			name:       "mixed tasks",
			markdown:   "- [x] Done\n- [ ] Todo\n- [X] Also done",
			wantCount:  3,
			wantStates: []bool{true, false, true},
		},
		{
			name:      "no tasks in regular list",
			markdown:  "- Regular item\n- Another item",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension(config.GetDefaultConfig())))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var checkboxes []*TaskCheckBox
			_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if cb, ok := node.(*TaskCheckBox); ok {
						checkboxes = append(checkboxes, cb)
					}
				}
				return ast.WalkContinue, nil
			})

			if len(checkboxes) != tt.wantCount {
				t.Errorf("Expected %d checkboxes, got %d", tt.wantCount, len(checkboxes))
				return
			}

			for i, cb := range checkboxes {
				if i < len(tt.wantStates) {
					isChecked := cb.State != " " && cb.State != ""
					if isChecked != tt.wantStates[i] {
						t.Errorf("Checkbox %d: expected checked = %v, got %v (state: %q)", i, tt.wantStates[i], isChecked, cb.State)
					}
				}
			}
		})
	}
}

func TestTaskListContextualParsing(t *testing.T) {
	// Test that checkboxes are only parsed in appropriate contexts
	tests := []struct {
		name      string
		markdown  string
		wantCount int
	}{
		{
			name:      "checkbox in list item",
			markdown:  "- [x] Task",
			wantCount: 1,
		},
		{
			name:      "checkbox in nested list",
			markdown:  "- Item\n  - [x] Subtask",
			wantCount: 1,
		},
		{
			name:      "checkbox not in list - paragraph",
			markdown:  "[x] Not a task",
			wantCount: 0,
		},
		{
			name:      "checkbox not in list - heading",
			markdown:  "# [x] Not a task",
			wantCount: 0,
		},
		{
			name:      "checkbox in code block",
			markdown:  "```\n- [x] Not a task\n```",
			wantCount: 0,
		},
		{
			name:      "checkbox in inline code",
			markdown:  "Code: `[x] not a task`",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension(config.GetDefaultConfig())))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var count int
			_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if _, ok := node.(*TaskCheckBox); ok {
						count++
					}
				}
				return ast.WalkContinue, nil
			})

			if count != tt.wantCount {
				t.Errorf("Expected %d checkboxes, got %d", tt.wantCount, count)
			}
		})
	}
}

func TestTaskListEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		markdown  string
		wantCount int
	}{
		{
			name:      "multiple spaces in checkbox",
			markdown:  "- [  ] Not a valid checkbox",
			wantCount: 0,
		},
		{
			name:      "tab in checkbox",
			markdown:  "- [\t] Not a valid checkbox",
			wantCount: 0,
		},
		{
			name:      "checkbox with extra characters",
			markdown:  "- [xx] Not a valid checkbox",
			wantCount: 0,
		},
		{
			name:      "checkbox case sensitivity",
			markdown:  "- [x] Lower\n- [X] Upper",
			wantCount: 2,
		},
		{
			name:      "checkbox at end of line",
			markdown:  "- [x]",
			wantCount: 1,
		},
		{
			name:      "checkbox with tabs and spaces",
			markdown:  "- [x]\t\tTab after\n- [ ]  Space after",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension(config.GetDefaultConfig())))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var count int
			_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					if _, ok := node.(*TaskCheckBox); ok {
						count++
					}
				}
				return ast.WalkContinue, nil
			})

			if count != tt.wantCount {
				t.Errorf("Expected %d checkboxes, got %d", tt.wantCount, count)
			}
		})
	}
}
