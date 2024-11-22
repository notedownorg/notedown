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

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
)

func buildClient(events []reader.Event, validators ...test.ContentUpdateValidator) (*tasks.TaskClient, chan reader.Event) {
	feed := make(chan reader.Event)
	go func() {
		for _, event := range events {
			feed <- event
		}
	}()

	client := tasks.NewClient(
		&test.MockDocumentWriter{Validators: test.Validators{ContentUpdate: validators}},
		feed,
		tasks.WithInitialLoadWaiter(100*time.Millisecond),
	)
	return client, feed
}

func date(year, month, day int, add time.Duration) *time.Time {
	res := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Add(add)
	return &res
}

var eventTasks = map[string][]tasks.Task{
	"zero.md": {
		tasks.NewTask(tasks.NewIdentifier("zero.md", "version", 1), "Task zero-0", tasks.Abandoned),
		tasks.NewTask(tasks.NewIdentifier("zero.md", "version", 2), "Task zero-1", tasks.Done, tasks.WithPriority(1)),
		tasks.NewTask(tasks.NewIdentifier("zero.md", "version", 3), "Task zero-2", tasks.Doing, tasks.WithPriority(1), tasks.WithDue(*date(1, 1, 1, 0)), tasks.WithCompleted(*date(1, 1, 1, 0)), tasks.WithScheduled(*date(1, 1, 1, 0))),
		tasks.NewTask(tasks.NewIdentifier("zero.md", "version", 4), "Task zero-3", tasks.Doing, tasks.WithDue(*date(1, 1, 2, 0)), tasks.WithCompleted(*date(1, 1, 2, 0)), tasks.WithScheduled(*date(1, 1, 2, 0))),
		tasks.NewTask(tasks.NewIdentifier("zero.md", "version", 5), "Task zero-4", tasks.Doing, tasks.WithDue(*date(1, 1, 3, 0)), tasks.WithCompleted(*date(1, 1, 3, 0)), tasks.WithScheduled(*date(1, 1, 3, 0))),
	},
	"one.md": {
		tasks.NewTask(tasks.NewIdentifier("one.md", "version", 1), "Task one-0", tasks.Doing, tasks.WithPriority(2)),
		tasks.NewTask(tasks.NewIdentifier("one.md", "version", 2), "Task one-1", tasks.Todo, tasks.WithPriority(3)),
		tasks.NewTask(tasks.NewIdentifier("one.md", "version", 3), "Task one-2", tasks.Blocked, tasks.WithPriority(4)),
		tasks.NewTask(tasks.NewIdentifier("one.md", "version", 4), "Task one-3", tasks.Blocked, tasks.WithPriority(4)),
	},
	"two.md":   {},
	"three.md": {},
}

func taskCount(events []reader.Event) int {
	count := 0
	for _, event := range events {
		if event.Op == reader.SubscriberLoadComplete {
			continue
		}
		if ev, ok := eventTasks[event.Key]; ok {
			count += len(ev)
		} else {
			panic("taskCount can only be called with events from loadEvents")
		}
	}
	return count
}

func loadEvents() []reader.Event {
	return []reader.Event{
		// Project with tasks
		{
			Op:  reader.Load,
			Key: "zero.md",
			Document: reader.Document{
				Metadata: reader.Metadata{reader.MetadataTypeKey: "project"},
				Contents: []byte(`- [a] Task zero-0
- [x] Task zero-1 p:1
- [/] Task zero-2 p:1 due:0001-01-01 completed:0001-01-01 scheduled:0001-01-01
- [/] Task zero-3 due:0001-01-02 completed:0001-01-02 scheduled:0001-01-02
- [/] Task zero-4 due:0001-01-03 completed:0001-01-03 scheduled:0001-01-03
`),
				Checksum: "version",
			},
		},
		// Document with tasks
		{
			Op:  reader.Load,
			Key: "one.md",
			Document: reader.Document{
				Contents: []byte(`- [/] Task one-0 p:2
- [ ] Task one-1 p:3
- [b] Task one-2 p:4
- [b] Task one-3 p:4`),
				Checksum: "version",
			},
		},
		// Document with no tasks
		{
			Op:  reader.Load,
			Key: "two.md",
			Document: reader.Document{
				Contents: []byte(``),
				Checksum: "version",
			},
		},
		// Project with no tasks
		{
			Op:  reader.Load,
			Key: "three.md",
			Document: reader.Document{
				Metadata: reader.Metadata{reader.MetadataTypeKey: "project"},
				Contents: []byte(``),
				Checksum: "version",
			},
		},
		// Load complete
		{
			Op: reader.SubscriberLoadComplete,
		},
	}
}
