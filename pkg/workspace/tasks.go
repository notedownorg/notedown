package workspace

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/liamawhite/nl/pkg/ast"
)

type Status ast.Status

const (
	Todo      Status = Status(ast.Todo)
	Blocked   Status = Status(ast.Blocked)
	Doing     Status = Status(ast.Doing)
	Done      Status = Status(ast.Done)
	Abandoned Status = Status(ast.Abandoned)
)

func fromAst(relativePath string, project string, task ast.Task) *Task {
	return &Task{
		id:        fmt.Sprintf("%s:%d", relativePath, task.Line),
		Name:      task.Name,
		Status:    Status(task.Status),
		Due:       task.Due,
		Scheduled: task.Scheduled,
		Completed: task.Completed,
		Priority:  task.Priority,
		Every:     task.Every,
		Project:   project,
	}
}

func toAst(line int, task Task) ast.Task {
	return ast.Task{
		Name:      task.Name,
		Line:      line,
		Status:    ast.Status(task.Status),
		Due:       task.Due,
		Scheduled: task.Scheduled,
		Completed: task.Completed,
		Priority:  task.Priority,
		Every:     task.Every,
	}
}

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
	Every     *ast.Every
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

// AddTask adds a task to the workspace in the passed document at the passed line number.
// If the line number is -1, the task will be added to the end of the document.
// If the line number is 0 or would be placed before the end of the front matter, the task will be added just after the front matter.
func (w *Workspace) AddTask(documentPath string, line int, task Task) error {
	slog.Debug("adding task", slog.String("file", documentPath), slog.String("task", task.Name))
	abs := w.absolutePath(documentPath)
	if _, ok := w.documents[abs]; !ok {
		return fmt.Errorf("document %s not found", documentPath)
	}

	if line >= 0 && line < w.documents[abs].markers.ContentStart {
		line = w.documents[abs].markers.ContentStart
	}

	if err := w.persistor.AddLine(abs, line, toAst(line, task).String()); err != nil {
		return err
	}

	// TODO: Are there forced cache invalidations that need to happen here?
	// TODO: Maybe we should just wait until the cache is updated before returning?
	// TODO: Or maybe enventual consistency is fine?
	return nil
}

func (w *Workspace) UpdateTask(task Task) error {
	panic("not implemented")
}

func (w *Workspace) DeleteTask(id string) error {
	panic("not implemented")
}
