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
	"github.com/notedownorg/notedown/pkg/config"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// TaskCheckBox represents a task list checkbox node with configurable states
type TaskCheckBox struct {
	ast.BaseInline
	State string // The actual state value (e.g., " ", "x", "wip", "in-progress")
}

// NewTaskCheckBox creates a new TaskCheckBox node with a state string
func NewTaskCheckBox(state string) *TaskCheckBox {
	return &TaskCheckBox{
		State: state,
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

// taskListParser is a parser for task list items with configurable states
type taskListParser struct {
	config *config.Config
}

// NewTaskListParser creates a new task list parser with configuration
func NewTaskListParser(cfg *config.Config) parser.InlineParser {
	return &taskListParser{
		config: cfg,
	}
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

	// Must start with [
	if line[0] != '[' {
		return nil
	}

	// Find the closing ]
	closingBracket := -1
	for i := 1; i < len(line); i++ {
		if line[i] == ']' {
			closingBracket = i
			break
		}
	}

	if closingBracket == -1 {
		return nil
	}

	// Extract the content between brackets
	stateValue := string(line[1:closingBracket])

	// Check if this state value is valid according to configuration
	if !s.isValidTaskState(stateValue) {
		return nil
	}

	// Consume the entire checkbox [state]
	block.Advance(closingBracket + 1)

	// Skip optional space after checkbox
	if closingBracket+1 < len(line) && line[closingBracket+1] == ' ' {
		block.Advance(1)
	}

	return NewTaskCheckBox(stateValue)
}

// isValidTaskState checks if the given state value is valid according to configuration
func (s *taskListParser) isValidTaskState(stateValue string) bool {
	// Use default configuration if none provided
	cfg := s.config
	if cfg == nil {
		cfg = config.GetDefaultConfig()
	}

	// Check each configured task state (including aliases)
	for _, state := range cfg.Tasks.States {
		if state.HasValue(stateValue) {
			return true
		}
	}

	return false
}

// TaskListExtension is an extension that adds support for configurable task lists
type TaskListExtension struct {
	config *config.Config
}

// NewTaskListExtension creates a new task list extension with configuration
func NewTaskListExtension(cfg *config.Config) goldmark.Extender {
	return &TaskListExtension{
		config: cfg,
	}
}

// Extend implements goldmark.Extender.Extend
func (e *TaskListExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTaskListParser(e.config), 0),
	))
}
