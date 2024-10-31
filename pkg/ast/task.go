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

package ast

import (
	"fmt"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
)

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
}

type Every struct {
	RRule *rrule.RRule
	Text  string // maintain the original text for every so we can write it back out
}

type TaskOption func(*Task)

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

// Used if you want to mutate/update a task
func NewTaskFromTask(t Task, options ...TaskOption) Task {
	task := Task{
		identifier: t.identifier,
		name:       t.name,
		status:     t.status,
		due:        t.due,
		scheduled:  t.scheduled,
		completed:  t.completed,
		priority:   t.priority,
		every:      t.every,
	}
	for _, option := range options {
		option(&task)
	}
	return task
}

func WithLine(line int) TaskOption {
	return func(t *Task) {
		t.identifier.line = line
	}
}

func WithName(name string) TaskOption {
	return func(t *Task) {
		t.name = name
	}
}

func WithStatus(status Status) TaskOption {
	return func(t *Task) {
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
		b.WriteString(fmt.Sprintf(" every:%v", t.every.Text))
	}
	if t.completed != nil {
		b.WriteString(fmt.Sprintf(" completed:%v", t.completed.Format("2006-01-02")))
	}
	return b.String()
}
