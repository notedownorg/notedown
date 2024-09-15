package tasks_test

import (
	"testing"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/workspace/tasks"
	"github.com/stretchr/testify/assert"
)

func TestFilters(t *testing.T) {
	events := defaultEvents()
	c, _ := buildClient(events...)

	tests := []struct {
		name      string
		filter    tasks.TaskFilter
		wantTasks []ast.Task
	}{
		{
			name:      "Filter by single priority",
			filter:    tasks.FilterByPriority(1),
			wantTasks: []ast.Task{events[0].Document.Tasks[1]},
		},
		{
			name:      "Filter by multiple priorities",
			filter:    tasks.FilterByPriority(1, 2),
			wantTasks: []ast.Task{events[0].Document.Tasks[1], events[1].Document.Tasks[0]},
		},
		{
			name:      "Filter by status",
			filter:    tasks.FilterByStatus(ast.Done),
			wantTasks: []ast.Task{events[1].Document.Tasks[0]},
		},
		{
			name:      "Filter by multiple statuses",
			filter:    tasks.FilterByStatus(ast.Todo, ast.Done),
			wantTasks: []ast.Task{events[0].Document.Tasks[0], events[1].Document.Tasks[0]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := c.ListTasks(tasks.FetchAllTasks(), tt.filter)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.wantTasks, tasks)
		})
	}

}
