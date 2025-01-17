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

package test

import (
	"strings"
	"time"

	"github.com/notedownorg/notedown/pkg/parse/ast"
	"github.com/notedownorg/notedown/pkg/parse/blocks"
	"github.com/teambition/rrule-go"
)

// Provide a DSL like API for creating test cases.
func Bl(blocks ...ast.Block) []ast.Block {
	return append([]ast.Block{}, blocks...)
}

func Tb(char, text string) ast.Block {
	return blocks.NewThematicBreak(char, text)
}

func P(text string) ast.Block {
	return blocks.NewParagraph(text)
}

var Bln = blocks.NewBlankLine()

func Ha1(indent int, value string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingAtx(indent, 1, value, append([]ast.Block{}, children...)...)
}

func Ha2(indent int, value string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingAtx(indent, 2, value, append([]ast.Block{}, children...)...)
}

func Ha3(indent int, value string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingAtx(indent, 3, value, append([]ast.Block{}, children...)...)
}

func Ha4(indent int, value string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingAtx(indent, 4, value, append([]ast.Block{}, children...)...)
}

func Ha5(indent int, value string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingAtx(indent, 5, value, append([]ast.Block{}, children...)...)
}

func Ha6(indent int, value string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingAtx(indent, 6, value, append([]ast.Block{}, children...)...)
}

func Hs1(value, underline string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingSetext(value, underline, append([]ast.Block{}, children...)...)
}

func Hs2(value, underline string, children ...ast.Block) ast.Block {
	return blocks.NewHeadingSetext(value, underline, append([]ast.Block{}, children...)...)
}

func Cbf(open, info, code, close string) ast.Block {
	return blocks.NewCodeBlockFenced(open, info, code, close)
}

func Cbi(code string) ast.Block {
	return blocks.NewCodeBlockIndented(strings.Split(code, "\n"))
}

func Ht1(text string) ast.Block {
	return blocks.NewHtml(blocks.HtmlOne, text)
}

func Ht2(text string) ast.Block {
	return blocks.NewHtml(blocks.HtmlTwo, text)
}

func Ht3(text string) ast.Block {
	return blocks.NewHtml(blocks.HtmlThree, text)
}

func Ht4(text string) ast.Block {
	return blocks.NewHtml(blocks.HtmlFour, text)
}

func Ht5(text string) ast.Block {
	return blocks.NewHtml(blocks.HtmlFive, text)
}

func Ht6(text string) ast.Block {
	return blocks.NewHtml(blocks.HtmlSix, text)
}

func Ht7(text string) ast.Block {
	return blocks.NewHtml(blocks.HtmlSeven, text)
}

func Bq(indent string, children ...ast.Block) ast.Block {
	return blocks.NewBlockQuote(indent, append([]ast.Block{}, children...)...)
}

func Ol(start int, children ...ast.Block) ast.Block {
	return blocks.NewOrderedList(start, append([]ast.Block{}, children...)...)
}

func Ul(children ...ast.Block) ast.Block {
	return blocks.NewUnorderedList(append([]ast.Block{}, children...)...)
}

func Tl(children ...ast.Block) ast.Block {
	return blocks.NewTaskList(append([]ast.Block{}, children...)...)
}

func Uli(external, marker, internal string, children ...ast.Block) ast.Block {
	return blocks.NewListItemUnordered(external, marker, internal, append([]ast.Block{}, children...)...)
}

func Oli(external, marker, internal string, children ...ast.Block) ast.Block {
	return blocks.NewListItemOrdered(external, marker, internal, append([]ast.Block{}, children...)...)
}

type task struct {
	status    blocks.TaskStatus
	text      string
	due       *time.Time
	scheduled *time.Time
	completed *time.Time
	priority  *int
	every     *rrule.RRule
	everyText string
}

type taskOption func(*task)

func Due(year, month, day int) taskOption {
	return func(t *task) {
		ti := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		t.due = &ti
	}
}

func Scheduled(year, month, day int) taskOption {
	return func(t *task) {
		ti := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		t.scheduled = &ti
	}
}

func Completed(year, month, day int) taskOption {
	return func(t *task) {
		ti := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		t.completed = &ti
	}
}

func Priority(p int) taskOption {
	return func(t *task) {
		t.priority = &p
	}
}

func Every(opts rrule.ROption, original string) taskOption {
	return func(t *task) {
		t.every, _ = rrule.NewRRule(opts)
		t.everyText = original
	}
}

func T(status blocks.TaskStatus, text string, opts ...taskOption) task {
	t := task{status: status, text: text}
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

func Tli(external, marker, internal string, task task, children ...ast.Block) ast.Block {
	opts := []blocks.ListItemTaskOption{blocks.TaskWithChildren(children...)}
	if task.due != nil {
		opts = append(opts, blocks.TaskWithDue(*task.due))
	}
	if task.scheduled != nil {
		opts = append(opts, blocks.TaskWithScheduled(*task.scheduled))
	}
	if task.completed != nil {
		opts = append(opts, blocks.TaskWithCompleted(*task.completed))
	}
	if task.priority != nil {
		opts = append(opts, blocks.TaskWithPriority(*task.priority))
	}
	if task.every != nil {
		opts = append(opts, blocks.TaskWithEvery(blocks.NewEvery(task.every, task.everyText)))
	}
	return blocks.NewListItemTask(external, marker, internal, task.status, task.text, opts...)
}

func Fm(metadata map[string]interface{}) ast.Block {
	return blocks.NewFrontmatter(metadata)
}
