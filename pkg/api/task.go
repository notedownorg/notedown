package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
)

type Status string

const (
	Todo      Status = "todo"
	Blocked   Status = "blocked"
	Doing     Status = "doing"
	Done      Status = "done"
	Abandoned Status = "abandoned"
)

type Task struct {
	Line      int
	Name      string
	Status    Status
	Due       *time.Time
	Scheduled *time.Time
	Completed *time.Time
	Priority  *int
	Every     *rrule.RRule
    Indent    int
    SubTasks  []*Task
}

func (t *Task) AddChild(task *Task) bool {
    fmt.Printf("trying to add %q to %q\n", task.Name, t.Name)
    if task.Indent <= t.Indent {
        fmt.Printf("ERR: %q indent %v is below parent %q %v\n", task.Name, task.Indent, t.Name, t.Indent)
        return false
    }
    if len(t.SubTasks) == 0 {
        fmt.Printf("%q has no subtasks so adding %q\n", t.Name, task.Name)
        t.SubTasks = append(t.SubTasks, task)
        return true
    }
    tail := t.SubTasks[len(t.SubTasks)-1]
    if task.Indent == tail.Indent {
        fmt.Printf("%q is a sibling of %q whos parent is %q\n", task.Name, tail.Name, t.Name)
        t.SubTasks = append(t.SubTasks, task)
        return true
    }
    if task.Indent > tail.Indent {
        fmt.Printf("%q indent %v is greater than subtask %q (%v) of %q, so recursing\n", task.Name, task.Indent, tail.Name, tail.Indent, t.Name)
        return tail.AddChild(task)
    }
    // indent is somehow in between, still add to my subtasks
    fmt.Printf("%q indent %v is somehow in between parent and child, adding anyway\n", task.Name, task.Indent)
    t.SubTasks = append(t.SubTasks, task)
    return true
}

func (t *Task) PrintChildren() {

    fmt.Printf("%s- [%s] %s\n", strings.Repeat(" ", t.Indent), t.Status, t.Name)
    for _, st := range t.SubTasks {
        st.PrintChildren()
    }

}
