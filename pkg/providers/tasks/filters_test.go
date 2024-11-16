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

	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
	"github.com/notedownorg/notedown/pkg/providers/tasks"
	"github.com/stretchr/testify/assert"
)

func TestTaskFilters(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name      string
		filter    collections.Filter[tasks.Task]
		wantTasks []tasks.Task
	}{
		{
			name:   "Filter by single priority",
			filter: tasks.FilterByPriority(1),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][1],
				eventTasks["zero.md"][2],
			},
		},
		{
			name:   "Filter by multiple priorities",
			filter: tasks.FilterByPriority(1, 2),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][1],
				eventTasks["zero.md"][2],
				eventTasks["one.md"][0],
			},
		},
		{
			name:      "Filter by status",
			filter:    tasks.FilterByStatus(tasks.Done),
			wantTasks: []tasks.Task{eventTasks["zero.md"][1]},
		},
		{
			name:      "Filter by multiple statuses",
			filter:    tasks.FilterByStatus(tasks.Todo, tasks.Done),
			wantTasks: []tasks.Task{eventTasks["one.md"][1], eventTasks["zero.md"][1]},
		},
		{
			name:      "Filter by due date when after and before are set",
			filter:    tasks.FilterByDueDate(date(1, 1, 2, 0), date(1, 1, 3, -1)),
			wantTasks: []tasks.Task{eventTasks["zero.md"][3]},
		},
		{
			name:   "Filter by due date is set using nil-nil",
			filter: tasks.FilterByDueDate(nil, nil),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
			},
		},
		{
			name:   "Filter by due date before",
			filter: tasks.FilterByDueDate(nil, date(1, 1, 2, 0)),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
			},
		},
		{
			name:   "Filter by due date after",
			filter: tasks.FilterByDueDate(date(1, 1, 2, 0), nil),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
			},
		},
		{
			name:      "Filter by completed date before and after",
			filter:    tasks.FilterByCompletedDate(date(1, 1, 2, 0), date(1, 1, 3, -1)),
			wantTasks: []tasks.Task{eventTasks["zero.md"][3]},
		},
		{
			name:   "Filter by completed date is set using nil-nil",
			filter: tasks.FilterByCompletedDate(nil, nil),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
			},
		},
		{
			name:   "Filter by completed date before",
			filter: tasks.FilterByCompletedDate(nil, date(1, 1, 2, 0)),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
			},
		},
		{
			name:   "Filter by completed date after",
			filter: tasks.FilterByCompletedDate(date(1, 1, 2, 0), nil),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
			},
		},
		{
			name:      "Filter by scheduled date before and after",
			filter:    tasks.FilterByScheduledDate(date(1, 1, 2, 0), date(1, 1, 3, -1)),
			wantTasks: []tasks.Task{eventTasks["zero.md"][3]},
		},
		{
			name:   "Filter by scheduled date is set using nil-nil",
			filter: tasks.FilterByScheduledDate(nil, nil),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
			},
		},
		{
			name:   "Filter by scheduled date before",
			filter: tasks.FilterByScheduledDate(nil, date(1, 1, 2, 0)),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][2],
				eventTasks["zero.md"][3],
			},
		},
		{
			name:   "Filter by scheduled date after",
			filter: tasks.FilterByScheduledDate(date(1, 1, 2, 0), nil),
			wantTasks: []tasks.Task{
				eventTasks["zero.md"][3],
				eventTasks["zero.md"][4],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.wantTasks, c.ListTasks(tasks.FetchAllTasks(), tasks.WithFilter(tt.filter)))
		})
	}
}
