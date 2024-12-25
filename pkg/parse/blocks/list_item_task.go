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

package blocks

import (
	"fmt"
	"strings"
	"time"

	"github.com/notedownorg/notedown/pkg/parse/ast"
	"github.com/teambition/rrule-go"
)

var _ ast.Block = &ListItemTask{}
var _ listItem = &ListItemTask{}

const ListItemTaskBlockType = "list_item_task"

type TaskStatus string

func (s TaskStatus) String() string {
	return fmt.Sprintf("[%s]", string(s))
}

const (
	Todo      TaskStatus = " "
	Blocked   TaskStatus = "b"
	Doing     TaskStatus = "/"
	Done      TaskStatus = "x"
	Abandoned TaskStatus = "a"
)

type TaskEvery struct {
	rrule *rrule.RRule
	text  string // maintain the original text for every so we can write it back out
}

func NewEvery(rrule *rrule.RRule, text string) TaskEvery {
	return TaskEvery{
		rrule: rrule,
		text:  text,
	}
}

type ListItemTask struct {
	*tracker
	*listItemUnordered

	Text   string
	Status TaskStatus

	Due       *time.Time
	Scheduled *time.Time
	Completed *time.Time
	Priority  *int
	Every     *TaskEvery
}

type ListItemTaskOption func(*ListItemTask)

func TaskWithChildren(children ...ast.Block) ListItemTaskOption {
	return func(l *ListItemTask) {
		l.children = append(l.children, children...)
	}
}

func TaskWithDue(due time.Time) ListItemTaskOption {
	return func(l *ListItemTask) {
		l.Due = &due
	}
}

func TaskWithScheduled(scheduled time.Time) ListItemTaskOption {
	return func(l *ListItemTask) {
		l.Scheduled = &scheduled
	}
}

func TaskWithCompleted(completed time.Time) ListItemTaskOption {
	return func(l *ListItemTask) {
		l.Completed = &completed
	}
}

func TaskWithPriority(priority int) ListItemTaskOption {
	return func(l *ListItemTask) {
		l.Priority = &priority
	}
}

func TaskWithEvery(every TaskEvery) ListItemTaskOption {
	return func(l *ListItemTask) {
		l.Every = &every
	}
}

func NewListItemTask(external, marker, internal string, status TaskStatus, text string, opts ...ListItemTaskOption) *ListItemTask {
	lit := &ListItemTask{
		listItemUnordered: &listItemUnordered{
			external: external,
			marker:   newBulletListItemMarker(marker),
			internal: internal,
		},
		Text:   text,
		Status: status,
	}

	for _, opt := range opts {
		opt(lit)
	}

	lit.tracker = newTracker(lit)
	return lit
}

func (l *ListItemTask) Type() ast.BlockType {
	return ListItemTaskBlockType
}

func (l *ListItemTask) Children() []ast.Block {
	return l.children
}

func (l *ListItemTask) Markdown() string {
	first := l.external + marker(l.marker) + l.internal + l.Status.String() + " " + l.Text
	if len(l.children) == 0 {
		return first
	}

	lines := []string{first}
	for _, child := range l.children {
		md := child.Markdown()
		childLines := strings.Split(md, "\n")
		for _, line := range childLines {
			// If the line is empty, dont add any indentation
			if line == "" {
				lines = append(lines, "")
				continue
			}

			// Otherwise, add the correct amount of indentation
			lines = append(lines, strings.Repeat(" ", len(l.external)+len(marker(l.marker))+len(l.internal))+line)
		}
	}
	return strings.Join(lines, "\n")

}

func (l *ListItemTask) Modified() bool {
	return l.tracker.Modified(l)
}

func (l *ListItemTask) SameType(li listItem) bool {
	tsk, ok := li.(*ListItemTask)
	if !ok {
		return false
	}
	return l.marker == tsk.marker
}
