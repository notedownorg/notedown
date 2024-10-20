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

	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllTasks(t *testing.T) {
	events := defaultEvents()
	c, _ := buildClient(events)
	tasks, err := c.ListTasks(tasks.FetchAllTasks())
	wantTasks := append(events[0].Document.Tasks, events[1].Document.Tasks...)

	assert.NoError(t, err)
	assert.ElementsMatch(t, wantTasks, tasks)
}

func TestFetchTasksForDocument(t *testing.T) {
	events := defaultEvents()
	c, _ := buildClient(events)
	tasks, err := c.ListTasks(tasks.FetchTasksForDocument("two.md"))
	wantTasks := events[1].Document.Tasks

	assert.NoError(t, err)
	assert.ElementsMatch(t, wantTasks, tasks)
}
