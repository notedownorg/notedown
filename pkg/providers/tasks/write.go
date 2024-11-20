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
	"log/slog"

	"github.com/notedownorg/notedown/pkg/fileserver/writer"
)

func (c *TaskClient) CreateTask(path string, line int, name string, status Status, options ...TaskOption) error {
	task := NewTask(NewIdentifier(path, "", line), name, status, options...)
	slog.Debug("creating task", "identifier", task.Identifier().String(), "task", task.String())

	mutation := writer.AddLine(task.Line(), task)
	if err := c.writer.UpdateContent(writer.Document{Path: path}, mutation); err != nil {
		return fmt.Errorf("failed to add task: %v: %w", task, err)
	}
	return nil
}

// Returns a new repeated task from an uncommitted repeat task
func newForRepeat(t Task) (Task, bool) {
	if !t.uncommittedRepeat || t.completed == nil || t.every == nil {
		return Task{}, false
	}

	// Update the due date/scheduled date to the next recurrence
	completed := *normalizeDate(*t.completed)

	// Set dtstart to the completed date to handle the case where the task is completed in the past
	// Unlikely to happen in practice but I hit this while testing
	t.every.rrule.DTStart(completed)

	if next := t.every.rrule.After(completed, false); next.Unix() != 0 {
		if t.scheduled != nil {
			WithScheduled(next)(&t)
		} else {
			WithDue(next)(&t)
		}
	}

	// Reset the status to todo and clear the uncommitted repeat flag
	t.status = Todo
	t.uncommittedRepeat = false
	t.completed = nil

	return t, true
}

func (c *TaskClient) UpdateTask(t Task) error {
	// If this task has been flagged as completed with recurrence handle it.
	if repeater, repeat := newForRepeat(t); repeat {
		// Task completion is handled by adding the completed task to the line below.
		// - [ ] Task due:2024-01-01 every:day
		// after:
		// - [ ] Task due:2024-01-02 every:day
		// - [x] Task due:2024-01-01 every:day completed:2024-01-01
		mutations := []writer.LineMutation{
			writer.UpdateLine(t.Line(), repeater),
			writer.AddLine(t.Line()+1, t),
		}
		slog.Debug("updating original task", "identifier", t.Identifier().String(), "task", t.String())
		slog.Debug("adding repeated task", "identifier", repeater.Identifier().String(), "task", repeater.String())
		if err := c.writer.UpdateContent(writer.Document{Path: t.Path(), Checksum: t.Version()}, mutations...); err != nil {
			return fmt.Errorf("failed to update task: %v: %w", t, err)
		}
		return nil
	}

	mutation := writer.UpdateLine(t.Line(), t)
	slog.Debug("updating task", "identifier", t.Identifier().String(), "task", t.String())
	if err := c.writer.UpdateContent(writer.Document{Path: t.Path(), Checksum: t.Version()}, mutation); err != nil {
		return fmt.Errorf("failed to update task: %v: %w", t, err)
	}
	return nil
}

func (c *TaskClient) DeleteTask(t Task) error {
	slog.Debug("deleting task", "identifier", t.Identifier().String(), "task", t.String())
	mutation := writer.RemoveLine(t.Line())
	if err := c.writer.UpdateContent(writer.Document{Path: t.Path(), Checksum: t.Version()}, mutation); err != nil {
		return fmt.Errorf("failed to remove task: %v: %w", t, err)
	}
	return nil
}
