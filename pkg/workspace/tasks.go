package workspace

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/teambition/rrule-go"
)

type Status ast.Status

const (
    Todo      Status = Status(ast.Todo)
    Blocked   Status = Status(ast.Blocked)
    Doing     Status = Status(ast.Doing)
    Done      Status = Status(ast.Done)
    Abandoned Status = Status(ast.Abandoned)
)

func OrderStatus(status Status) int {
	switch status {
	case Todo:
		return 1
	case Blocked:
		return 2
	case Doing:
		return 3
	case Done:
		return 4
	case Abandoned:
		return 5
	}
	return 0
}

type Task struct {
	Id        string
	Name      string
	Status    Status
	Due       *time.Time
	Scheduled *time.Time
	Completed *time.Time
	Priority  *int
	Every     *rrule.RRule
	Project   string
}

func (w Workspace) ListTasks() []Task {
	res := []Task{}
	for path, tasks := range w.tasks {
		for line, task := range tasks {
			res = append(res, Task{
				Id:        fmt.Sprintf("%s:%d", path, line),
				Name:      task.Name,
				Status:    Status(task.Status),
				Due:       task.Due,
				Scheduled: task.Scheduled,
				Completed: task.Completed,
				Priority:  task.Priority,
				Every:     task.Every,
				Project:   strings.ReplaceAll(path, filepath.Ext(path), ""),
			})
		}
	}
	return res
}

func (w *Workspace) AddTask(task *ast.Task) error {
	panic("not implemented")
}

func (w *Workspace) UpdateTask(task *ast.Task) error {
	panic("not implemented")
}

func (w *Workspace) DeleteTask(id string) error {
	panic("not implemented")
}
