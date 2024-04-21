package workspace

import (
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
	id        string
	Name      string
	Status    Status
	Due       *time.Time
	Scheduled *time.Time
	Completed *time.Time
	Priority  *int
	Every     *rrule.RRule
	Project   string
}

func (t Task) Id() string {
    return t.id
}

func (w Workspace) ListTasks() []Task {
	res := []Task{}
	for _, tasks := range w.tasks {
		for _, task := range tasks {
			res = append(res, *task)
		}
	}
	return res
}

func (w *Workspace) AddTask(task Task) error {
	panic("not implemented")
}

func (w *Workspace) UpdateTask(task Task) error {
	panic("not implemented")
}

func (w *Workspace) DeleteTask(id string) error {
	panic("not implemented")
}
