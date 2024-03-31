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

func OrderStatus(status Status) int {
	switch ast.Status(status) {
	case ast.Todo:
		return 1
	case ast.Blocked:
		return 2
	case ast.Doing:
		return 3
	case ast.Done:
		return 4
	case ast.Abandoned:
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
