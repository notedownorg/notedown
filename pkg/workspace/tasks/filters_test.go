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
	"github.com/notedownorg/notedown/pkg/workspace/documents/reader"
	"github.com/notedownorg/notedown/pkg/workspace/tasks"
	"github.com/stretchr/testify/assert"
)

func TestTaskFilters(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name      string
		filter    tasks.TaskFilter
		wantTasks []ast.Task
	}{
		{
			name:   "Filter by single priority",
			filter: tasks.FilterByPriority(1),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[1],
				events[0].Document.Tasks[2],
			},
		},
		{
			name:   "Filter by multiple priorities",
			filter: tasks.FilterByPriority(1, 2),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[1],
				events[0].Document.Tasks[2],
				events[1].Document.Tasks[0],
			},
		},
		{
			name:      "Filter by status",
			filter:    tasks.FilterByStatus(ast.Done),
			wantTasks: []ast.Task{events[0].Document.Tasks[1]},
		},
		{
			name:      "Filter by multiple statuses",
			filter:    tasks.FilterByStatus(ast.Todo, ast.Done),
			wantTasks: []ast.Task{events[1].Document.Tasks[1], events[0].Document.Tasks[1]},
		},
		{
			name:      "Filter by due date when after and before are set",
			filter:    tasks.FilterByDueDate(date(1, 1, 2, 0), date(1, 1, 3, -1)),
			wantTasks: []ast.Task{events[0].Document.Tasks[3]},
		},
		{
			name:   "Filter by due date is set using nil-nil",
			filter: tasks.FilterByDueDate(nil, nil),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[2],
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[4],
			},
		},
		{
			name:   "Filter by due date before",
			filter: tasks.FilterByDueDate(nil, date(1, 1, 2, 0)),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[2],
				events[0].Document.Tasks[3],
			},
		},
		{
			name:   "Filter by due date after",
			filter: tasks.FilterByDueDate(date(1, 1, 2, 0), nil),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[4],
			},
		},
		{
			name:      "Filter by completed date before and after",
			filter:    tasks.FilterByCompletedDate(date(1, 1, 2, 0), date(1, 1, 3, -1)),
			wantTasks: []ast.Task{events[0].Document.Tasks[3]},
		},
		{
			name:   "Filter by completed date is set using nil-nil",
			filter: tasks.FilterByCompletedDate(nil, nil),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[2],
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[4],
			},
		},
		{
			name:   "Filter by completed date before",
			filter: tasks.FilterByCompletedDate(nil, date(1, 1, 2, 0)),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[2],
				events[0].Document.Tasks[3],
			},
		},
		{
			name:   "Filter by completed date after",
			filter: tasks.FilterByCompletedDate(date(1, 1, 2, 0), nil),
			wantTasks: []ast.Task{
				events[0].Document.Tasks[3],
				events[0].Document.Tasks[4],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.wantTasks, c.ListTasks(tasks.FetchAllTasks(), tasks.WithFilters(tt.filter)))
		})
	}
}

func TestDocumentFilters(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name          string
		filter        tasks.DocumentFilter
		wantDocuments map[string]reader.Document
	}{
		{
			name:   "Filter by document type",
			filter: tasks.FilterByDocumentType("project"),
			wantDocuments: map[string]reader.Document{
				"zero.md":  events[0].Document,
				"three.md": events[3].Document,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantDocuments, c.ListDocuments(tasks.FetchAllDocuments(), tt.filter))
		})
	}
}
