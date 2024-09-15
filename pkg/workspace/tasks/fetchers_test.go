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
	wantTasks := append(events[0].Document.Tasks, events[1].Document.Tasks...)

	assert.NoError(t, err)
	assert.ElementsMatch(t, wantTasks, tasks)
}

func TestFetchTasksForDocument(t *testing.T) {
	events := defaultEvents()
	c, _ := buildClient(events...)
	tasks, err := c.ListTasks(tasks.FetchTasksForDocument("two.md"))
	wantTasks := events[1].Document.Tasks

	assert.NoError(t, err)
	assert.ElementsMatch(t, wantTasks, tasks)
}
