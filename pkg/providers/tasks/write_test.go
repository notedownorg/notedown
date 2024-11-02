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

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {

	client, _ := buildClient([]reader.Event{{Op: reader.SubscriberLoadComplete}},
		// Create
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "add", method)
			assert.Equal(t, writer.Document{Path: "path"}, doc)
			assert.Equal(t, writer.AT_END, line)
			return nil
		},

		// Update
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "update", method)
			assert.Equal(t, writer.Document{Path: "path", Checksum: "version"}, doc)
			assert.Equal(t, 7, line)
			return nil
		},

		// Delete
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "remove", method)
			assert.Equal(t, writer.Document{Path: "path", Checksum: "version"}, doc)
			assert.Equal(t, 1, line)
			return nil
		},

		// Move
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "add", method)
			assert.Equal(t, writer.Document{Path: "newPath"}, doc)
			assert.Equal(t, writer.AT_END, line)
			return nil
		},
		func(method string, doc writer.Document, line int, obj fmt.Stringer) error {
			assert.Equal(t, "remove", method)
			assert.Equal(t, writer.Document{Path: "path", Checksum: "version"}, doc)
			assert.Equal(t, 7, line)
			return nil
		},
	)

	assert.NoError(t, client.Create("path", "Task", tasks.Todo, tasks.WithLine(writer.AT_END)))
	assert.NoError(t, client.Update(tasks.NewTask(tasks.NewIdentifier("path", "version"), "Task", tasks.Todo, tasks.WithLine(7))))
	assert.NoError(t, client.Delete(tasks.NewTask(tasks.NewIdentifier("path", "version"), "Task", tasks.Todo, tasks.WithLine(1))))
	assert.NoError(t, client.Move(tasks.NewTask(tasks.NewIdentifier("path", "version"), "Task", tasks.Todo, tasks.WithLine(7)), "newPath"))

}
