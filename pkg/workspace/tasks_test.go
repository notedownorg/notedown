package workspace

import (
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

func copyTestData(t *testing.T) string {
	tmp, err := os.MkdirTemp("", "nl-test-")
	if err != nil {
		t.Fatal(err)
	}

	err = cp.Copy("testdata/workspace", tmp)
	if err != nil {
		t.Fatal(err)
	}

	return tmp
}

func date(year, month, day int) *time.Time {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return &date
}

var deterministicTasks = func(tasks []Task) []Task {
    slices.SortFunc(tasks, func(a, b Task) int { return strings.Compare(a.Id(), b.Id()) })
    return tasks
}

func TestWorkspace_Tasks(t *testing.T) {
	// Copy the testdata into a temporary directory so we don't modify the original
	tmp := copyTestData(t)

	ws, err :=New(tmp)
	if err != nil {
		t.Fatal(err)
	}

	// Check the tasks were loaded correctly
	time.Sleep(1 * time.Second) // remove once we have a way to wait for the initial state to be built
	tasks := deterministicTasks(ws.ListTasks())
	assert.Equal(t, []Task{
		{Name: "Project One, Task One", id: "project-one.md:4", Project: "project-one", Status:Todo, Due: date(2024, 1, 1)},
		{Name: "Project One, Task Two", id: "project-one.md:5", Project: "project-one", Status:Done, Due: date(2024, 1, 1)},
	}, tasks)

    // Check adding tasks
    ws.AddTask(Task{Name: "New Task", Status:Todo, Due: date(2024, 1, 1)})
}
