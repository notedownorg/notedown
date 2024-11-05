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
	"fmt"
	"testing"
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	client, _ := buildClient([]reader.Event{{Op: reader.SubscriberLoadComplete}},

		// Create
		func(doc writer.Document, mutations ...writer.LineMutation) error {
			assert.Equal(t, writer.Document{Path: "path"}, doc)
			lines := []string{"line 1", "line 2", "line 3"}
			for _, mutation := range mutations {
				lines, _ = mutation("", lines)
			}
			assert.Equal(t, []string{"line 1", "line 2", "line 3", "- [ ] Task"}, lines)
			return nil
		},

		// Update
		func(doc writer.Document, mutations ...writer.LineMutation) error {
			assert.Equal(t, writer.Document{Path: "path", Checksum: "version"}, doc)
			lines := []string{"line 1", "line 2", "line 3"}
			for _, mutation := range mutations {
				lines, _ = mutation("version", lines)
			}
			assert.Equal(t, []string{"line 1", "line 2", "- [ ] Task"}, lines)
			return nil
		},

		// Update with recurrence completion
		func(doc writer.Document, mutations ...writer.LineMutation) error {
			assert.Equal(t, writer.Document{Path: "path", Checksum: "version"}, doc)
			lines := []string{"line 1", "line 2", "- [ ] Task every:day"}
			for _, mutation := range mutations {
				lines, _ = mutation("version", lines)
			}
			completed := fmt.Sprintf("- [x] Task every:day completed:%s", time.Now().Format("2006-01-02"))
			assert.Equal(t, []string{"line 1", "line 2", "- [ ] Task every:day", completed}, lines)
			return nil
		},

		// Delete
		func(doc writer.Document, mutations ...writer.LineMutation) error {
			assert.Equal(t, writer.Document{Path: "path", Checksum: "version"}, doc)
			lines := []string{"line 1", "- [ ] Task", "line 3"}
			for _, mutation := range mutations {
				lines, _ = mutation("version", lines)
			}
			assert.Equal(t, []string{"line 1", "line 3"}, lines)
			return nil
		},
	)

	assert.NoError(t, client.Create("path", writer.AT_END, "Task", tasks.Todo))
	assert.NoError(t, client.Update(tasks.NewTask(tasks.NewIdentifier("path", "version", 3), "Task", tasks.Todo)))

	every, _ := tasks.NewEvery("day")
	original := tasks.NewTask(tasks.NewIdentifier("path", "version", 3), "Task", tasks.Todo, tasks.WithEvery(every))
	completed := tasks.NewTaskFromTask(original, tasks.WithStatus(tasks.Done))
	assert.NoError(t, client.Update(completed))

	assert.NoError(t, client.Delete(tasks.NewTask(tasks.NewIdentifier("path", "version", 2), "Task", tasks.Todo)))

}
