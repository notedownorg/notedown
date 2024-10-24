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

package tasks_test

import (
	"time"

	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/documents/reader"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
)

func buildClient(events []reader.Event, validators ...validator) (*tasks.Client, chan reader.Event) {
	feed := make(chan reader.Event)
	go func() {
		for _, event := range events {
			feed <- event
		}
	}()

	client := tasks.NewClient(
		&MockLineWriter{validators: validators},
		feed,
		tasks.WithInitialLoadWaiter(100*time.Millisecond),
	)
	return client, feed
}

func loadEvents() []reader.Event {
	return []reader.Event{
		// Project with tasks
		{
			Op:  reader.Load,
			Key: "one.md",
			Document: reader.Document{
				Document: ast.Document{
					Metadata: ast.Metadata{ast.MetadataType: "project"},
					Tasks: []ast.Task{
						ast.NewTask(ast.NewIdentifier("initial-one.md", "version"), "Task 1", ast.Todo, ast.WithLine(1)),
						ast.NewTask(ast.NewIdentifier("initial-one.md", "version"), "Task 2", ast.Doing, ast.WithLine(2), ast.WithPriority(1)),
					},
				},
			},
		},
		// Document with tasks
		{
			Op:  reader.Load,
			Key: "two.md",
			Document: reader.Document{
				Document: ast.Document{
					Tasks: []ast.Task{
						ast.NewTask(ast.NewIdentifier("intitial-two.md", "version"), "Task 3", ast.Done, ast.WithLine(1), ast.WithPriority(2)),
						ast.NewTask(ast.NewIdentifier("intitial-two.md", "version"), "Task 4", ast.Abandoned, ast.WithLine(2), ast.WithPriority(3)),
					},
				},
			},
		},
		// Document with no tasks
		{
			Op:       reader.Load,
			Key:      "three.md",
			Document: reader.Document{Document: ast.Document{}},
		},
		// Project with no tasks
		{
			Op:  reader.Load,
			Key: "four.md",
			Document: reader.Document{
				Document: ast.Document{
					Metadata: ast.Metadata{ast.MetadataType: "project"},
				},
			},
		},
		// Load complete
		{
			Op: reader.SubscriberLoadComplete,
		},
	}
}

func defaultEvents() []reader.Event {
	return []reader.Event{
		// Project with tasks
		{
			Op:  reader.Change,
			Key: "one.md",
			Document: reader.Document{
				Document: ast.Document{
					Metadata: ast.Metadata{ast.MetadataType: "project"},
					Tasks: []ast.Task{
						ast.NewTask(ast.NewIdentifier("one.md", "version"), "Task 1", ast.Todo, ast.WithLine(1)),
						ast.NewTask(ast.NewIdentifier("one.md", "version"), "Task 2", ast.Doing, ast.WithLine(2), ast.WithPriority(1)),
					},
				},
			},
		},
		// Document with tasks
		{
			Op:  reader.Change,
			Key: "two.md",
			Document: reader.Document{
				Document: ast.Document{
					Tasks: []ast.Task{
						ast.NewTask(ast.NewIdentifier("two.md", "version"), "Task 3", ast.Done, ast.WithLine(1), ast.WithPriority(2)),
						ast.NewTask(ast.NewIdentifier("two.md", "version"), "Task 4", ast.Abandoned, ast.WithLine(2), ast.WithPriority(3)),
						ast.NewTask(ast.NewIdentifier("two.md", "version"), "Task 5", ast.Blocked, ast.WithLine(3), ast.WithPriority(10)),
					},
				},
			},
		},
		// Document with no tasks
		{
			Op:       reader.Change,
			Key:      "three.md",
			Document: reader.Document{Document: ast.Document{}},
		},
		// Project with no tasks
		{
			Op:  reader.Change,
			Key: "four.md",
			Document: reader.Document{
				Document: ast.Document{
					Metadata: ast.Metadata{ast.MetadataType: "project"},
				},
			},
		},
	}
}
