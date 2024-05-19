package ast

import (
	"fmt"
	"time"

	"github.com/teambition/rrule-go"
)

type Status string

const (
	Todo      Status = " "
	Blocked   Status = "b"
	Doing     Status = "/"
	Done      Status = "x"
	Abandoned Status = "a"
)

type Task struct {
	Line      int
	Name      string
	Status    Status
	Due       *time.Time
	Scheduled *time.Time
	Completed *time.Time
	Priority  *int
	Every     *Every
}

type Every struct {
	RRule *rrule.RRule
	Text  string // maintain the original text for every so we can write it back out
}

func (t Task) String() string {
	res := fmt.Sprintf("- [%v] %v", t.Status, t.Name)
	if t.Due != nil {
		res = fmt.Sprintf("%v due:%v", res, t.Due.Format("2006-01-02"))
	}
	if t.Scheduled != nil {
		res = fmt.Sprintf("%v scheduled:%v", res, t.Scheduled.Format("2006-01-02"))
	}
	if t.Priority != nil {
		res = fmt.Sprintf("%v priority:%v", res, *t.Priority)
	}
	if t.Every != nil {
		res = fmt.Sprintf("%v every:%v", res, t.Every.Text)
	}
	if t.Completed != nil {
		res = fmt.Sprintf("%v completed:%v", res, t.Completed.Format("2006-01-02"))
	}
	return res
}
