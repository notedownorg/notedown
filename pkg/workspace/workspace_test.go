package workspace_test

import (
	"os"
	"testing"
	"time"

	"github.com/liamawhite/nl/pkg/workspace"
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

func TestWorkspace(t *testing.T) {
	// Copy the testdata into a temporary directory so we don't modify the original
	tmp := copyTestData(t)

	ws, err := workspace.New(tmp)
	if err != nil {
		t.Fatal(err)
	}

	// Check the tasks were loaded correctly
	time.Sleep(1 * time.Second) // Remove once we have a way to wait for the initial state to be built
	tasks := ws.ListTasks()
	assert.Equal(t, []workspace.Task{
		{Name: "Project One, Task One", Id: "project-one.md:4", Project: "project-one", Status: workspace.Todo, Due: date(2024, 1, 1)},
		{Name: "Project One, Task Two", Id: "project-one.md:5", Project: "project-one", Status: workspace.Done, Due: date(2024, 1, 1)},
	}, tasks)
}
