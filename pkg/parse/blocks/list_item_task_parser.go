// Copyright 2025 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blocks

import (
	"strconv"
	"time"

	. "github.com/liamawhite/parse/core"
	. "github.com/liamawhite/parse/time"
	"github.com/teambition/rrule-go"
)

type taskBuilder struct {
	text   string
	status TaskStatus
	opts   []ListItemTaskOption
}

// Assumes the list bullet has already been parsed and we're effectively being passed a paragraph
func listItemTaskParser(ctx Context) Parser[*taskBuilder] {
	return func(in Input) (*taskBuilder, bool, error) {
		start := in.Checkpoint()

		// Parse the task status
		status, found, err := taskStatusParser(in)
		if err != nil || !found {
			return nil, false, err
		}

		// Read up until we hit a field, a new line, or the end of the input
		text, found, err := StringWhileNotEOFOr(Any(StringFrom(SequenceOf2(RuneIn(" "), anyFieldKey)), NewLine))(in)
		if err != nil || !found {
			in.Restore(start)
			return nil, false, err
		}

		// If we've hit newline or EOF, we're done
		next, _ := in.Peek(1)
		if next == "\n" || next == "" {
			NewLine(in) // consume the newline if we're done for roundtrip consistency
			return &taskBuilder{text: text, status: status}, true, nil
		}

		// Otherwise, we have fields to parse
		// Because there is no defined order of fields we scan for the presence of each field before
		// the end of the line and parse them if they are present, resetting the input each time.
		taskOpts := make([]ListItemTaskOption, 0)

		due, ok, err := fieldParser(dueKey, dueParser)(in)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}
		if ok {
			taskOpts = append(taskOpts, TaskWithDue(due))
		}

		scheduled, ok, err := fieldParser(scheduledKey, scheduledParser)(in)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}
		if ok {
			taskOpts = append(taskOpts, TaskWithScheduled(scheduled))
		}

		completed, ok, err := fieldParser(completedKey, completedParser)(in)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}
		if ok {
			taskOpts = append(taskOpts, TaskWithCompleted(completed))
		}

		priority, ok, err := fieldParser(priorityKey, priorityParser)(in)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}
		if ok {
			taskOpts = append(taskOpts, TaskWithPriority(priority))
		}

		every, ok, err := fieldParser(everyKey, everyParser(ctx))(in)
		if err != nil {
			in.Restore(start)
			return nil, false, err
		}
		if ok {
			taskOpts = append(taskOpts, TaskWithEvery(every))
		}

		// Consume the rest of the line and then the newline if we're not at EOF
		StringWhileNotEOFOr(NewLine)(in)
		NewLine(in)

		return &taskBuilder{text: text, status: status, opts: taskOpts}, true, nil
	}
}

var taskStatusLookup = map[string]TaskStatus{
	" ": Todo,
	"x": Done,
	"X": Done,
	"/": Doing,
	"b": Blocked,
	"B": Blocked,
	"a": Abandoned,
	"A": Abandoned,
}

func taskStatusParser(in Input) (TaskStatus, bool, error) {
	// Read the open bracket
	_, ok, err := Rune('[')(in)
	if err != nil || !ok {
		return "", false, err
	}

	// Read the status rune
	s, ok, err := RuneIn(" xX/bBaA")(in)
	if err != nil || !ok {
		return "", false, err
	}

	// Read the close bracket
	_, ok, err = Rune(']')(in)
	if err != nil || !ok {
		return "", false, err
	}

	// Eat the trailing space
	_, ok, err = Rune(' ')(in)
	if err != nil || !ok {
		return "", false, err
	}

	return taskStatusLookup[s], true, nil
}

var (
	dueKeyLong  = String("due:")
	dueKeyShort = String("d:")
	dueKey      = Any(dueKeyLong, dueKeyShort)

	scheduledKeyLong  = String("scheduled:")
	scheduledKeyShort = String("s:")
	scheduledKey      = Any(scheduledKeyLong, scheduledKeyShort)

	completedKey = Any(String("completed:"))

	everyKeyLong  = String("every:")
	everyKeyShort = String("e:")
	everyKey      = Any(everyKeyLong, everyKeyShort)

	priorityKeyLong  = String("priority:")
	priorityKeyShort = String("p:")
	priorityKey      = Any(priorityKeyLong, priorityKeyShort)

	anyFieldKey = Any(dueKey, scheduledKey, everyKey, priorityKey, completedKey)
)

func dueParser(in Input) (time.Time, bool, error) {
	start := in.Checkpoint()
	_, ok, err := dueKey(in)
	if err != nil || !ok {
		return time.Time{}, false, err
	}
	t, found, err := YearMonthDay(in)
	if err != nil || !found {
		in.Restore(start)
		return t, false, err
	}
	return t, true, nil
}

func scheduledParser(in Input) (time.Time, bool, error) {
	start := in.Checkpoint()
	_, ok, err := scheduledKey(in)
	if err != nil || !ok {
		return time.Time{}, false, err
	}
	t, found, err := YearMonthDay(in)
	if err != nil || !found {
		in.Restore(start)
		return t, false, err
	}
	return t, true, nil
}

func completedParser(in Input) (time.Time, bool, error) {
	start := in.Checkpoint()
	_, ok, err := completedKey(in)
	if err != nil || !ok {
		return time.Time{}, false, err
	}
	t, found, err := YearMonthDay(in)
	if err != nil || !found {
		in.Restore(start)
		return t, false, err
	}
	return t, true, nil
}

func priorityParser(in Input) (int, bool, error) {
	start := in.Checkpoint()
	_, ok, err := priorityKey(in)
	if err != nil || !ok {
		return 0, false, err
	}
	priority, ok, err := StringFrom(AtLeast(1, Digit))(in)
	if err != nil || !ok {
		in.Restore(start)
		return 0, false, err
	}
	p, err := strconv.Atoi(priority)
	if err != nil {
		in.Restore(start)
		return 0, false, nil
	}
	return p, true, nil
}

func everyParser(ctx Context) Parser[TaskEvery] {
	return func(in Input) (TaskEvery, bool, error) {
		start := in.Checkpoint()
		_, ok, err := everyKey(in)
		if err != nil || !ok {
			return TaskEvery{}, false, err
		}

		// Helper function for reducing the boilerplate of returning the final result and handling errors
		handleResult := func(opts rrule.ROption, err error) (TaskEvery, bool, error) {
			if err != nil {
				in.Restore(start)
				return TaskEvery{}, false, err
			}
			rr, err := rrule.NewRRule(opts)
			if err != nil {
				in.Restore(start)
				return TaskEvery{}, false, err
			}

			// Extract the original text for roundtripping
			end := in.Checkpoint()
			in.Restore(start)
			text, ok := in.Take(end - start)
			if !ok {
				return TaskEvery{}, false, nil
			}
			return TaskEvery{rrule: rr, text: text}, true, nil
		}

		rruleOpts := rrule.ROption{Dtstart: ctx.RelativeTo()}

		// There are a limited number of single words that can be used to describe the frequency.
		// So lets get those out of the way first. (day, week, month, year, weekday, weekend)
		// Note that the order of these is important, as "week" is a prefix of "weekday" and "weekend".
		single, ok, err := Any(Day, String("weekend"), String("weekday"), Month, Year, Week)(in)
		if err != nil {
			return handleResult(rruleOpts, err)
		}
		if ok {
			switch single {
			case "day":
				rruleOpts.Freq = rrule.DAILY
			case "week":
				rruleOpts.Freq = rrule.WEEKLY
			case "month":
				rruleOpts.Freq = rrule.MONTHLY
			case "year":
				rruleOpts.Freq = rrule.YEARLY
			case "weekday":
				rruleOpts.Byweekday = []rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR}
				rruleOpts.Freq = rrule.WEEKLY
			case "weekend":
				rruleOpts.Byweekday = []rrule.Weekday{rrule.SA}
				rruleOpts.Freq = rrule.WEEKLY
			}
			return handleResult(rruleOpts, nil)
		}

		// Every <day of week> or list of <day of week>
		daysOfWeek, ok, err := DaysOfWeek(in)
		if err != nil {
			return handleResult(rruleOpts, err)
		}
		if ok {
			for _, d := range daysOfWeek {
				rruleOpts.Byweekday = append(rruleOpts.Byweekday, rruleDayOfWeek(d))
			}
			rruleOpts.Freq = rrule.WEEKLY
			return handleResult(rruleOpts, nil)
		}

		// Every <number> <day/week/month/year>
		tuple, ok, err := SequenceOf3(
			StringFrom(AtLeast(1, RuneIn("0123456789"))),
			String(" "),
			Any(Day, Week, Month, Year),
		)(in)
		if err != nil {
			return handleResult(rruleOpts, err)
		}
		if ok {
			n, _, unit := tuple.Values()
			switch unit {
			case "day", "days":
				rruleOpts.Freq = rrule.DAILY
			case "week", "weeks":
				rruleOpts.Freq = rrule.WEEKLY
			case "month", "months":
				rruleOpts.Freq = rrule.MONTHLY
			case "year", "years":
				rruleOpts.Freq = rrule.YEARLY
			}
			rruleOpts.Interval, _ = strconv.Atoi(n)
			return handleResult(rruleOpts, nil)
		}

		// Some combination of month days and/or months
		// Keep reading input until we can't parse a month day or month.
		optionalDelimiter := Optional(Rune(' '))
		for {
			optionalDelimiter(in)
			monthDay, monthDayOk, err := MonthDay(in)
			if err != nil {
				return handleResult(rruleOpts, err)
			}
			if monthDayOk {
				rruleOpts.Bymonthday = append(rruleOpts.Bymonthday, monthDay)
			}

			optionalDelimiter(in)
			month, monthOk, err := MonthOfYear(in)
			if err != nil {
				return handleResult(rruleOpts, err)
			}
			if monthOk {
				rruleOpts.Bymonth = append(rruleOpts.Bymonth, int(month))
			}

			if !monthDayOk && !monthOk {
				break
			}
		}
		if len(rruleOpts.Bymonthday) > 0 || len(rruleOpts.Bymonth) > 0 {
			// If there are no days set, default to the first of the month.
			if len(rruleOpts.Bymonthday) == 0 {
				rruleOpts.Bymonthday = append(rruleOpts.Bymonthday, 1)
			}
			return handleResult(rruleOpts, nil)
		}
		return TaskEvery{}, false, nil
	}
}

func rruleDayOfWeek(d time.Weekday) rrule.Weekday {
	switch d {
	case time.Sunday:
		return rrule.SU
	case time.Monday:
		return rrule.MO
	case time.Tuesday:
		return rrule.TU
	case time.Wednesday:
		return rrule.WE
	case time.Thursday:
		return rrule.TH
	case time.Friday:
		return rrule.FR
	case time.Saturday:
		return rrule.SA
	}
	return rrule.MO
}

func fieldParser[T any](key Parser[string], full Parser[T]) Parser[T] {
	return func(in Input) (T, bool, error) {
		var empty T
		start := in.Checkpoint()

		// Read until we hit the field key or a newline
		_, ok, err := StringWhileNot(Any(StringFrom(SequenceOf2(RuneIn(" \t"), key)), NewLine))(in)
		if err != nil {
			return empty, false, err
		}

		// If we hit a newline, there is no field to parse
		if ch, _ := in.Peek(1); ch == "\n" {
			in.Restore(start)
			return empty, false, nil
		}

		// Otherwise parse the field
		if ok {
			//
			v, ok, err := SequenceOf2(RuneIn(" \t"), full)(in)
			if err != nil {
				return empty, false, err
			}
			if ok {
				_, value := v.Values()
				in.Restore(start) // restore to the start of the fields so we can parse the next field
				return value, true, nil
			}
		}

		// Gettting here means we hit the field key but couldn't parse the field
		// So reset the input and return false
		in.Restore(start)
		return empty, false, nil
	}
}
