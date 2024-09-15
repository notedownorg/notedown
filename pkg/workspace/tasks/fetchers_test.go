package tasks_test

import (
	"testing"

	"github.com/liamawhite/nl/pkg/workspace/tasks"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllTasks(t *testing.T) {
	events := defaultEvents()
	c, _ := buildClient(events...)
	tasks, err := c.ListTasks(tasks.FetchAllTasks())
	wantTasks := append(tasksBuilder(events[0].Document), tasksBuilder(events[1].Document)...)

	assert.NoError(t, err)
	assert.ElementsMatch(t, wantTasks, tasks)
}

func TestFetchTasksForDocument(t *testing.T) {
	events := defaultEvents()
	c, _ := buildClient(events...)
	tasks, err := c.ListTasks(tasks.FetchTasksForDocument("two.md"))
	wantTasks := tasksBuilder(events[1].Document)

	assert.NoError(t, err)
	assert.ElementsMatch(t, wantTasks, tasks)
}
