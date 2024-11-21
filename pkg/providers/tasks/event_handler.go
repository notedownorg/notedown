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
	"log/slog"
	"time"

	"github.com/a-h/parse"
	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	. "github.com/notedownorg/notedown/pkg/parsers"
	"github.com/notedownorg/notedown/pkg/providers/pkg/traits"
)

func onLoad(c *TaskClient) traits.EventHandler {
	return func(event reader.Event) {
		if c.handleChanges(event) {
			c.publisher.Events <- Event{Op: Load}
		}
	}
}

func onChange(c *TaskClient) traits.EventHandler {
	return func(event reader.Event) {
		if c.handleChanges(event) {
			c.publisher.Events <- Event{Op: Change}
		}
	}
}

func onDelete(c *TaskClient) traits.EventHandler {
	return func(event reader.Event) {
		c.tasksMutex.Lock()
		delete(c.tasks, event.Key)
		c.tasksMutex.Unlock()
		c.publisher.Events <- Event{Op: Delete}
	}
}

func (c *TaskClient) handleChanges(event reader.Event) bool {
	tasks := make(map[int]Task)

	// Go through the contents block by block in search of tasks
	in := parse.NewInput(string(event.Document.Contents))
	blocks, ok, err := parse.Until(parseBlock(event.Key, event.Document.Checksum, time.Now()), parse.EOF[string]()).Parse(in)
	if err != nil {
		slog.Error("failed to parse blocks", slog.String("file", event.Key), slog.String("error", err.Error()))
		return false
	}
	if !ok {
		slog.Debug("no blocks found", slog.String("file", event.Key))
		return false
	}
	for _, block := range blocks {
		for _, task := range block {
			tasks[task.Line()] = task
		}
	}

	c.tasksMutex.Lock()
	c.tasks[event.Key] = tasks
	c.tasksMutex.Unlock()
	slog.Debug("updated task in cache", "path", event.Key)
	return true
}

var parseBlock = func(path, version string, relativeTo time.Time) parse.Parser[[]Task] {
	return parse.Func(func(in *parse.Input) ([]Task, bool, error) {
		var res []Task

		// Drop any leading newline
		_, _, err := parse.NewLine.Parse(in)

		for {
			task, ok, err := ParseTask(path, version, relativeTo).Parse(in)
			if err != nil {
				return nil, false, err
			}
			if !ok {
				break
			}
			res = append(res, task)

		}

		// Process the input until the next newline or EOF as the current line isnt a task
		_, _, err = parse.StringUntil(NewLineOrEOF).Parse(in)
		if err != nil {
			return nil, false, err
		}

		return res, true, nil
	})
}
