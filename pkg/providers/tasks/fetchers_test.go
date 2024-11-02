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
	"testing"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllTasks(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)
	tasks := c.ListTasks(tasks.FetchAllTasks())
	wantTasks := append(eventTasks["zero.md"], eventTasks["one.md"]...)
	assert.ElementsMatch(t, wantTasks, tasks)
}

func TestFetchTasksForDocument(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)
	tasks := c.ListTasks(tasks.FetchTasksForDocument("one.md"))
	wantTasks := eventTasks["one.md"]
	assert.ElementsMatch(t, wantTasks, tasks)
}

func TestFetchAllDocuments(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)
	documents := c.ListDocuments(tasks.FetchAllDocuments())
	wantDocuments := map[string]reader.Document{
		"zero.md":  events[0].Document,
		"one.md":   events[1].Document,
		"two.md":   events[2].Document,
		"three.md": events[3].Document,
	}
	assert.Equal(t, wantDocuments, documents)
}
