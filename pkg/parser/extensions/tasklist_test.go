package extensions

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

func TestTaskCheckBox_NewTaskCheckBox(t *testing.T) {
	tests := []struct {
		name    string
		checked bool
	}{
		{
			name:    "unchecked checkbox",
			checked: false,
		},
		{
			name:    "checked checkbox",
			checked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTaskCheckBox(tt.checked)

			if node == nil {
				t.Fatal("Expected node to be created, got nil")
			}

			if node.IsChecked != tt.checked {
				t.Errorf("Expected IsChecked = %v, got %v", tt.checked, node.IsChecked)
			}

			if node.Kind() != KindTaskCheckBox {
				t.Errorf("Expected kind %v, got %v", KindTaskCheckBox, node.Kind())
			}
		})
	}
}

func TestTaskCheckBox_Dump(t *testing.T) {
	node := NewTaskCheckBox(true)

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
	parser := NewTaskListParser()
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
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension()))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var checkboxes []*TaskCheckBox
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
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
				if cb.IsChecked != tt.wantChecked[i] {
					t.Errorf("Checkbox %d: expected IsChecked = %v, got %v", i, tt.wantChecked[i], cb.IsChecked)
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
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension()))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var count int
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func TestTaskListHTMLRenderer_RegisterFuncs(t *testing.T) {
	renderer := NewTaskListHTMLRenderer()

	// Create a mock registerer to test function registration
	registered := make(map[ast.NodeKind]bool)
	mockReg := &mockRegisterer{registered: registered}

	renderer.RegisterFuncs(mockReg)

	if !registered[KindTaskCheckBox] {
		t.Error("Expected TaskCheckBox renderer to be registered")
	}
}

func TestTaskListHTMLRenderer_RenderTaskCheckBox(t *testing.T) {
	tests := []struct {
		name     string
		checked  bool
		wantHTML string
	}{
		{
			name:     "checked checkbox",
			checked:  true,
			wantHTML: `<input checked="" disabled="" type="checkbox">`,
		},
		{
			name:     "unchecked checkbox",
			checked:  false,
			wantHTML: `<input disabled="" type="checkbox">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewTaskListHTMLRenderer()
			node := NewTaskCheckBox(tt.checked)
			buf := &bytes.Buffer{}
			writer := &testBufWriter{buf}

			r := renderer.(*taskListHTMLRenderer)
			status, err := r.renderTaskCheckBox(writer, []byte{}, node, true)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if status != ast.WalkContinue {
				t.Errorf("Expected WalkContinue, got %v", status)
			}

			html := buf.String()
			if html != tt.wantHTML {
				t.Errorf("Expected HTML %q, got %q", tt.wantHTML, html)
			}
		})
	}
}

func TestTaskListHTMLRenderer_RenderTaskCheckBoxNotEntering(t *testing.T) {
	// Test that renderer does nothing when not entering
	renderer := NewTaskListHTMLRenderer()
	node := NewTaskCheckBox(true)
	buf := &bytes.Buffer{}
	writer := &testBufWriter{buf}

	r := renderer.(*taskListHTMLRenderer)
	status, err := r.renderTaskCheckBox(writer, []byte{}, node, false)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if status != ast.WalkContinue {
		t.Errorf("Expected WalkContinue, got %v", status)
	}

	html := buf.String()
	if html != "" {
		t.Errorf("Expected empty output when not entering, got %q", html)
	}
}

func TestTaskListIntegration(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantHTML string
	}{
		{
			name:     "simple task list",
			markdown: "- [x] Completed task\n- [ ] Incomplete task",
			wantHTML: `<ul>
<li><input checked="" disabled="" type="checkbox">Completed task</li>
<li><input disabled="" type="checkbox">Incomplete task</li>
</ul>
`,
		},
		{
			name:     "mixed task and regular list",
			markdown: "- [x] Task item\n- Regular item\n- [ ] Another task",
			wantHTML: `<ul>
<li><input checked="" disabled="" type="checkbox">Task item</li>
<li>Regular item</li>
<li><input disabled="" type="checkbox">Another task</li>
</ul>
`,
		},
		{
			name:     "ordered task list",
			markdown: "1. [x] First task\n2. [ ] Second task",
			wantHTML: `<ol>
<li><input checked="" disabled="" type="checkbox">First task</li>
<li><input disabled="" type="checkbox">Second task</li>
</ol>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(NewTaskListExtension()),
			)

			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.markdown), &buf); err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			html := buf.String()
			if html != tt.wantHTML {
				t.Errorf("Expected HTML:\n%s\nGot HTML:\n%s", tt.wantHTML, html)
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
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension()))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var checkboxes []*TaskCheckBox
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
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
				if i < len(tt.wantStates) && cb.IsChecked != tt.wantStates[i] {
					t.Errorf("Checkbox %d: expected IsChecked = %v, got %v", i, tt.wantStates[i], cb.IsChecked)
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
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension()))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var count int
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
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
			md := goldmark.New(goldmark.WithExtensions(NewTaskListExtension()))
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var count int
			ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
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

// Helper functions for testing

func createNodeOfKind(kind ast.NodeKind) ast.Node {
	switch kind {
	case ast.KindDocument:
		return ast.NewDocument()
	case ast.KindParagraph:
		return ast.NewParagraph()
	case ast.KindTextBlock:
		return ast.NewTextBlock()
	case ast.KindListItem:
		return ast.NewListItem(1)
	case ast.KindList:
		return ast.NewList('-')
	default:
		return ast.NewParagraph() // fallback
	}
}

// mockRegisterer implements renderer.NodeRendererFuncRegisterer for testing
type mockRegisterer struct {
	registered map[ast.NodeKind]bool
}

func (r *mockRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	r.registered[kind] = true
}

// testBufWriter implements util.BufWriter for testing
type testBufWriter struct {
	*bytes.Buffer
}

func (w *testBufWriter) Buffered() int {
	return w.Len()
}

func (w *testBufWriter) Available() int {
	return 0 // Not used in our tests
}

func (w *testBufWriter) WriteString(s string) (int, error) {
	return w.Buffer.WriteString(s)
}

func (w *testBufWriter) WriteByte(c byte) error {
	return w.Buffer.WriteByte(c)
}

func (w *testBufWriter) WriteRune(r rune) (int, error) {
	return w.Buffer.WriteRune(r)
}

func (w *testBufWriter) Flush() error {
	return nil // Buffer doesn't need flushing
}
