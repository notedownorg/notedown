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

	"github.com/notedownorg/notedown/pkg/ast"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/stretchr/testify/assert"
)

func TestSorters(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name      string
		sorter    tasks.TaskSorter
		wantTasks []ast.Task
	}{
		{
			name:   "Sort by status -> kanban order",
			sorter: tasks.SortByStatus(tasks.KanbanOrder()),
			wantTasks: []ast.Task{
				events[1].Document.Tasks[1],
				events[1].Document.Tasks[2],
				events[1].Document.Tasks[0],
				events[0].Document.Tasks[4],
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[2],
				events[0].Document.Tasks[1],
				events[0].Document.Tasks[0],
			},
		},
		{
			name:   "Sort by status -> agenda order",
			sorter: tasks.SortByStatus(tasks.AgendaOrder()),
			wantTasks: []ast.Task{
				events[1].Document.Tasks[0],
				events[0].Document.Tasks[2],
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[4],
				events[1].Document.Tasks[1],
				events[1].Document.Tasks[2],
				events[0].Document.Tasks[1],
				events[0].Document.Tasks[0],
			},
		},
		{
			name:   "Sort by priority",
			sorter: tasks.SortByPriority(),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[1],
				events[0].Document.Tasks[2],
				events[1].Document.Tasks[0],
				events[1].Document.Tasks[1],
				events[1].Document.Tasks[2],
				events[0].Document.Tasks[0],
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[4],
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
		wantTasks []ast.Task
	}{
		{
			name: "Sort by status -> agenda order, then by priority",
			sorters: []tasks.TaskSorter{
				tasks.SortByStatus(tasks.AgendaOrder()),
				tasks.SortByPriority(),
			},
			wantTasks: []ast.Task{
				events[0].Document.Tasks[2],
				events[1].Document.Tasks[0],
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[4],
				events[1].Document.Tasks[1],
				events[1].Document.Tasks[2],
				events[0].Document.Tasks[1],
				events[0].Document.Tasks[0],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantTasks, c.ListTasks(tasks.FetchAllTasks(), tasks.WithSorters(tt.sorters...)))
		})
	}
}
