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

func (c *Client) Create(path string, line int, name string, status Status, options ...TaskOption) error {
	task := NewTask(NewIdentifier(path, "", line), name, status, options...)
	slog.Debug("creating task", "identifier", task.Identifier().String(), "task", task.String())

	mutation := writer.AddLine(task.Line(), task)
	if err := c.writer.UpdateContent(writer.Document{Path: path}, mutation); err != nil {
		return fmt.Errorf("failed to add task: %v: %w", task, err)
	}
	return nil
}

func (c *Client) Update(t Task) error {
	slog.Debug("updating task", "identifier", t.Identifier().String(), "task", t.String())

	// If this task has been flagged as completed with recurrence handle it.
	if t.uncommittedRepeat {
		// Task completion is handled by adding the completed task to the line below.
		// - [ ] Task every:day
		// after:
		// - [ ] Task every:day
		// - [x] Task every:day completed:2024-01-01
		mutation := writer.AddLine(t.Line()+1, t)
		if err := c.writer.UpdateContent(writer.Document{Path: t.Path(), Checksum: t.Version()}, mutation); err != nil {
			return fmt.Errorf("failed to update task: %v: %w", t, err)
		}
		return nil
	}

	mutation := writer.UpdateLine(t.Line(), t)
	if err := c.writer.UpdateContent(writer.Document{Path: t.Path(), Checksum: t.Version()}, mutation); err != nil {
		return fmt.Errorf("failed to update task: %v: %w", t, err)
	}
	return nil
}

func (c *Client) Delete(t Task) error {
	slog.Debug("deleting task", "identifier", t.Identifier().String(), "task", t.String())
	mutation := writer.RemoveLine(t.Line())
	if err := c.writer.UpdateContent(writer.Document{Path: t.Path(), Checksum: t.Version()}, mutation); err != nil {
		return fmt.Errorf("failed to remove task: %v: %w", t, err)
	}
	return nil
}
