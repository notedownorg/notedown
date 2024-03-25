package workspace

import (
	"fmt"
	"time"

	"github.com/liamawhite/nl/pkg/api"
	"github.com/teambition/rrule-go"
)

type Task struct {
	Id        string
	Name      string
	Status    api.Status
	Due       *time.Time
	Scheduled *time.Time
	Completed *time.Time
	Priority  *int
	Every     *rrule.RRule
}

func (w Workspace) ListTasks() []Task {
	res := []Task{}
	for path, tasks := range w.tasks {
		for line, task := range tasks {
			res = append(res, Task{
				Id:        fmt.Sprintf("%s:%d", path, line),
				Name:      task.Name,
				Status:    task.Status,
				Due:       task.Due,
				Scheduled: task.Scheduled,
				Completed: task.Completed,
				Priority:  task.Priority,
				Every:     task.Every,
			})
		}
	}
	return res
}

func (w *Workspace) AddTask(task *api.Task) error {
	panic("not implemented")
}

func (w *Workspace) UpdateTask(task *api.Task) error {
	panic("not implemented")
}

func (w *Workspace) DeleteTask(id string) error {
	panic("not implemented")
}
