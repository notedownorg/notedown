package api

import (
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
    Line        int
	Name      string
	Status    Status
	Due       *time.Time
	Scheduled *time.Time
    Completed *time.Time
	Priority  *int
	Every     *rrule.RRule
}
