package workspace

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
)

func copyTestData(t *testing.T, name string) string {
	// If we're running in a CI environment, we dont want to create temp directories
	// This ensures we can store the artifacts for debugging
	dir := "/testdata"
	if os.Getenv("CI") != "true" {
        var err error
		dir, err = os.MkdirTemp("", fmt.Sprintf("nl-%v-", name))
		if err != nil {
			t.Fatal(err)
		}
	}

    if err := cp.Copy("testdata/workspace", dir); err != nil {
		t.Fatal(err)
	}

	return dir
}

func date(year, month, day int) *time.Time {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return &date
}

var deterministicTasks = func(tasks []Task) []Task {
	slices.SortFunc(tasks, func(a, b Task) int { return strings.Compare(a.Id(), b.Id()) })
	return tasks
}

var originalTasks = []Task{
	{Name: "Project One, Task One", id: "projects/project-one.md:4", Project: "project-one", Status: Todo, Due: date(2024, 1, 1)},
	{Name: "Project One, Task Two", id: "projects/project-one.md:5", Project: "project-one", Status: Done, Due: date(2024, 1, 1)},
}

func taskWithData(id string, project string, task Task) Task {
	task.Project = project
	task.id = id
	return task
}

func TestWorkspace_Tasks(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// Copy the testdata into a temporary directory so we don't modify the original
	tmp := copyTestData(t, "tasks")
	fmt.Println("Created temp dir: ", tmp)

	ws, err := New(tmp)
	assert.NoError(t, err)

	t.Run("Workspace Load", func(t *testing.T) {
		time.Sleep(1 * time.Second) // remove once we have a way to wait for the initial state to be built
		tasks := deterministicTasks(ws.ListTasks())
		assert.Equal(t, originalTasks, tasks)
	})

	t.Run("Add Tasks", func(t *testing.T) {
		// Add tasks to existing projects
		inFrontmatter := Task{Name: "Frontmatter Task", Status: Abandoned, Due: date(2024, 1, 2)}
		assert.NoError(t, ws.AddTask("projects/project-one.md", 1, inFrontmatter)) // Add a task to the frontmatter of a project
		beginningOfProject := Task{Name: "Beginning Task", Status: Doing}
		assert.NoError(t, ws.AddTask("projects/project-one.md", 0, beginningOfProject)) // Add a task to the beginning of a project
		endOfProject := Task{Name: "End Task", Status: Todo}
		assert.NoError(t, ws.AddTask("projects/project-one.md", -1, endOfProject)) // Add a task to the end of a project

		// Add tasks without a project (should be added to the daily note)
		// For the first one the note wont exist yet so make we ensure it's created
		dailyNotePath, _ := ws.DailyNotePath(*date(2021, 1, 1))
		dailyTask1 := Task{Name: "Daily Task 1", Status: Done}
		assert.NoError(t, ws.AddTask(dailyNotePath, 0, dailyTask1)) // Add a task to a non-existent daily note (will be created)
		dailyTask2 := Task{Name: "Daily Task 2", Status: Todo}
		assert.NoError(t, ws.AddTask(dailyNotePath, -1, dailyTask2)) // Add a task to an existing note

		time.Sleep(1 * time.Second) // remove once we have a way to wait for updates to be processed
		tasks := deterministicTasks(ws.ListTasks())

		assert.Equal(t, deterministicTasks([]Task{
			taskWithData("projects/project-one.md:3", "project-one", beginningOfProject),
			taskWithData("projects/project-one.md:4", "project-one", inFrontmatter),
			taskWithData("projects/project-one.md:6", "project-one", originalTasks[0]),
			taskWithData("projects/project-one.md:7", "project-one", originalTasks[1]),
			taskWithData("projects/project-one.md:10", "project-one", endOfProject),
			taskWithData("daily/2021-01-01.md:0", "", dailyTask1),
			taskWithData("daily/2021-01-01.md:2", "", dailyTask2),
		}), tasks)

	})
}
