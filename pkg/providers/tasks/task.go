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

package tasks

import (
	"fmt"
	"strings"
	"time"

	"github.com/a-h/parse"
	"github.com/teambition/rrule-go"
)

type Identifier struct {
	path    string
	line    int // line is 1-indexed, not 0
	version string
}

// Line is 1-indexed, not 0-indexed
func NewIdentifier(path string, version string, line int) Identifier {
	return Identifier{path: path, version: version, line: line}
}

func (i Identifier) String() string {
	// Pipe separators are good enough for now but may need to be changed as pipes
	// are technically valid (although unlikely to actually be used) in unix file paths
	// We may want to consider an actual encoding scheme for this in the future.
	var builder strings.Builder
	builder.WriteString(i.path)
	builder.WriteString("|")
	builder.WriteString(i.version)
	builder.WriteString("|")
	builder.WriteString(fmt.Sprintf("%v", i.line))
	return builder.String()
}

type Status string

const (
	Todo      Status = " "
	Blocked   Status = "b"
	Doing     Status = "/"
	Done      Status = "x"
	Abandoned Status = "a"
)

type Task struct {
	identifier Identifier
	name       string
	status     Status
	due        *time.Time
	scheduled  *time.Time
	completed  *time.Time
	priority   *int
	every      *Every

	// This is used to track if the task has been mutated to done and has an every set
	// but not yet written back to the file. This is so that we can handle the repeat.
	uncommittedRepeat bool
}

type Every struct {
	rrule *rrule.RRule
	text  string // maintain the original text for every so we can write it back out
}

func NewEvery(text string) (Every, error) {
	// Handle e:<text> vs every:<text> vs <text>
	if !strings.HasPrefix(text, "e:") && !strings.HasPrefix(text, "every:") {
		text = "e:" + text
	}

	e, ok, err := everyParser(time.Now()).Parse(parse.NewInput(text))
	if err != nil {
		return Every{}, fmt.Errorf("failed to parse every text: %w", err)
	}
	if !ok {
		return Every{}, fmt.Errorf("failed to parse every text: unable to find a valid recurrence")
	}
	return e, nil
}

func (e Every) String() string {
	return e.text
}

type TaskOption func(*Task)

// Used to create new tasks. For mutating tasks, use NewTaskFromTask.
func NewTask(identifier Identifier, name string, status Status, options ...TaskOption) Task {
	task := Task{
		identifier: identifier,
		name:       name,
		status:     status,
	}
	for _, option := range options {
		option(&task)
	}
	return task
}

// Used if you want to mutate/update a task, crucially this does not expose any method to update the identifier
// This means that if this new task is passed to the update method, we will know which task to update.
func NewTaskFromTask(t Task, options ...TaskOption) Task {
	task := Task{
		identifier:        t.identifier,
		name:              t.name,
		status:            t.status,
		due:               t.due,
		scheduled:         t.scheduled,
		completed:         t.completed,
		priority:          t.priority,
		every:             t.every,
		uncommittedRepeat: t.uncommittedRepeat,
	}
	for _, option := range options {
		option(&task)
	}
	return task
}

func WithName(name string) TaskOption {
	return func(t *Task) {
		t.name = name
	}
}

func WithStatus(status Status) TaskOption {
	return func(t *Task) {

		// If the task is being marked as done and it wasn't done before...
		if status == Done && t.status != Done {

			// If there is no completed time, set it to now
			if t.completed == nil {
				WithCompleted(time.Now())(t)
			}

			// If the task has a repeat, mark the task so we know we need to handle it when persisting
			if t.every != nil {
				t.uncommittedRepeat = true
			}
		}
		t.status = status
	}
}

func WithDue(due time.Time) TaskOption {
	return func(t *Task) {
		t.due = &due
	}
}

func WithScheduled(scheduled time.Time) TaskOption {
	return func(t *Task) {
		t.scheduled = &scheduled
	}
}

func WithCompleted(completed time.Time) TaskOption {
	return func(t *Task) {
		t.completed = &completed
	}
}

func WithPriority(priority int) TaskOption {
	return func(t *Task) {
		t.priority = &priority
	}
}

func WithEvery(every Every) TaskOption {
	return func(t *Task) {
		t.every = &every
	}
}

func (t Task) Identifier() Identifier {
	return t.identifier
}

func (t Task) Line() int {
	return t.Identifier().line
}

func (t Task) Path() string {
	return t.Identifier().path
}

func (t Task) Version() string {
	return t.Identifier().version
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Status() Status {
	return t.status
}

func (t Task) Due() *time.Time {
	if t.due == nil {
		return nil
	}
	res := *t.due
	return &res
}

func (t Task) Scheduled() *time.Time {
	if t.scheduled == nil {
		return nil
	}
	res := *t.scheduled
	return &res
}

func (t Task) Completed() *time.Time {
	if t.completed == nil {
		return nil
	}
	res := *t.completed
	return &res
}

func (t Task) Priority() *int {
	if t.priority == nil {
		return nil
	}
	res := *t.priority
	return &res
}

func (t Task) Every() *Every {
	if t.every == nil {
		return nil
	}
	res := *t.every
	return &res
}

func (t Task) String() string {
	return fmt.Sprintf("- [%v] %v", t.status, t.Body())
}

func (t Task) Body() string {
	var b strings.Builder
	b.WriteString(t.name)
	if t.due != nil {
		b.WriteString(fmt.Sprintf(" due:%v", t.due.Format("2006-01-02")))
	}
	if t.scheduled != nil {
		b.WriteString(fmt.Sprintf(" scheduled:%v", t.scheduled.Format("2006-01-02")))
	}
	if t.priority != nil {
		b.WriteString(fmt.Sprintf(" priority:%v", *t.priority))
	}
	if t.every != nil {
		b.WriteString(fmt.Sprintf(" every:%s", t.every))
	}
	if t.completed != nil {
		b.WriteString(fmt.Sprintf(" completed:%v", t.completed.Format("2006-01-02")))
	}
	return b.String()
}
