package api

import (
	"time"

	"github.com/teambition/rrule-go"
)

type Status string

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
}
