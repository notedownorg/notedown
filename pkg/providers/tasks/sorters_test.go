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

	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/stretchr/testify/assert"
)

func TestSorters(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name      string
		sorter    tasks.TaskSorter
		wantTasks []tasks.Task
	}{
		{
			name:   "Sort by status -> kanban order (then alphabetical)",
			sorter: tasks.SortByStatus(tasks.KanbanOrder()),
			wantTasks: []tasks.Task{
				eventTasks["one.md"][1],
				eventTasks["one.md"][2],
				eventTasks["one.md"][3],
				eventTasks["one.md"][0],
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
				eventTasks["zero.md"][1],
				eventTasks["zero.md"][0],
			},
		},
		{
			name:   "Sort by status -> agenda order (then alphabetical)",
			sorter: tasks.SortByStatus(tasks.AgendaOrder()),
			wantTasks: []tasks.Task{
				eventTasks["one.md"][0],
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
				eventTasks["one.md"][1],
				eventTasks["one.md"][2],
				eventTasks["one.md"][3],
				eventTasks["zero.md"][1],
				eventTasks["zero.md"][0],
			},
		},
		{
			name:   "Sort by priority (then alphabetical)",
			sorter: tasks.SortByPriority(),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][1],
				eventTasks["zero.md"][2],
				eventTasks["one.md"][0],
				eventTasks["one.md"][1],
				eventTasks["one.md"][2],
				eventTasks["one.md"][3],
				eventTasks["zero.md"][0],
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantTasks, c.ListTasks(tasks.FetchAllTasks(), tasks.WithSorters(tt.sorter)))
		})
	}
}

func TestSortersMultiple(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name      string
		sorters   []tasks.TaskSorter
		wantTasks []tasks.Task
	}{
		{
			name: "Sort by status -> agenda order, then by priority",
			sorters: []tasks.TaskSorter{
				tasks.SortByStatus(tasks.AgendaOrder()),
				tasks.SortByPriority(),
			},
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][2],
				eventTasks["one.md"][0],
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
				eventTasks["one.md"][1],
				eventTasks["one.md"][2],
				eventTasks["one.md"][3],
				eventTasks["zero.md"][1],
				eventTasks["zero.md"][0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantTasks, c.ListTasks(tasks.FetchAllTasks(), tasks.WithSorters(tt.sorters...)))
		})
	}
}
