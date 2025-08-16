package extensions

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// TaskCheckBox represents a task list checkbox node ala GFM
type TaskCheckBox struct {
	ast.BaseInline
	IsChecked bool
}

// NewTaskCheckBox creates a new TaskCheckBox node
func NewTaskCheckBox(checked bool) *TaskCheckBox {
	return &TaskCheckBox{
		IsChecked: checked,
	}
}

// Dump implements ast.Node.Dump
func (n *TaskCheckBox) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindTaskCheckBox is a NodeKind of the TaskCheckBox node
var KindTaskCheckBox = ast.NewNodeKind("TaskCheckBox")

// Kind implements ast.Node.Kind
func (n *TaskCheckBox) Kind() ast.NodeKind {
	return KindTaskCheckBox
}

// taskListParser is a parser for task list items
type taskListParser struct{}

// NewTaskListParser creates a new task list parser
func NewTaskListParser() parser.InlineParser {
	return &taskListParser{}
}

// Trigger implements parser.InlineParser.Trigger
func (s *taskListParser) Trigger() []byte {
	return []byte{'['}
}

// Parse implements parser.InlineParser.Parse
func (s *taskListParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 3 {
		return nil
	}

	// Check if we're at the beginning of a list item
	if parent.Kind() != ast.KindTextBlock {
		return nil
	}

	// Check if parent's parent is a list item
	if parent.Parent() == nil || parent.Parent().Kind() != ast.KindListItem {
		return nil
	}

	// Check if this is the first child of the text block
	if parent.FirstChild() != nil {
		return nil
	}

	// Look for [x], [ ], or [X] pattern
	if line[0] != '[' {
		return nil
	}

	var checked bool
	if len(line) >= 3 && line[2] == ']' {
		switch line[1] {
		case ' ':
			checked = false
		case 'x', 'X':
			checked = true
		default:
			return nil
		}
	} else {
		return nil
	}

	// Consume the checkbox
	block.Advance(3)

	// Skip optional space after checkbox
	if len(line) > 3 && line[3] == ' ' {
		block.Advance(1)
	}

	return NewTaskCheckBox(checked)
}

// taskListHTMLRenderer is a renderer for task list checkboxes
type taskListHTMLRenderer struct{}

// NewTaskListHTMLRenderer creates a new task list HTML renderer
func NewTaskListHTMLRenderer() renderer.NodeRenderer {
	return &taskListHTMLRenderer{}
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs
func (r *taskListHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindTaskCheckBox, r.renderTaskCheckBox)
}

func (r *taskListHTMLRenderer) renderTaskCheckBox(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	checkbox := n.(*TaskCheckBox)
	if checkbox.IsChecked {
		_, _ = w.WriteString(`<input checked="" disabled="" type="checkbox">`)
	} else {
		_, _ = w.WriteString(`<input disabled="" type="checkbox">`)
	}

	return ast.WalkContinue, nil
}

// TaskListExtension is an extension that adds support for task lists
type TaskListExtension struct{}

// NewTaskListExtension creates a new task list extension
func NewTaskListExtension() goldmark.Extender {
	return &TaskListExtension{}
}

// Extend implements goldmark.Extender.Extend
func (e *TaskListExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTaskListParser(), 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewTaskListHTMLRenderer(), 500),
	))
}
