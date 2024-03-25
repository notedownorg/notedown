package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/api"
	"github.com/teambition/rrule-go"
)

var statusLookup = map[string]api.Status{
	" ": api.Todo,
	"x": api.Done,
	"X": api.Done,
	"/": api.Doing,
	"b": api.Blocked,
	"B": api.Blocked,
	"a": api.Abandoned,
	"A": api.Abandoned,
}

var statusRuneLookup = map[api.Status]rune{
	api.Todo:      ' ',
	api.Blocked:   'b',
	api.Doing:     '/',
	api.Done:      'x',
	api.Abandoned: 'a',
}

var statusParser = parse.Func(func(in *parse.Input) (api.Status, bool, error) {
	// Read the open bracket
	_, ok, err := parse.Rune('[').Parse(in)
	if err != nil || !ok {
		return "", false, err
	}

	// Read the status rune
	s, ok, err := parse.RuneIn(" xX/bBaA").Parse(in)
	if err != nil || !ok {
		return "", false, err
	}

	// Read the close bracket
	_, ok, err = parse.Rune(']').Parse(in)
	if err != nil || !ok {
		return "", false, err
	}

	// Eat the trailing space
	_, ok, err = parse.Rune(' ').Parse(in)
	if err != nil || !ok {
		return "", false, err
	}

	return statusLookup[s], true, nil
})

var listItemOpen = parse.StringFrom(remainingInlineWhitespace, parse.Rune('-'), remainingInlineWhitespace)
var listItemOpen2 = parse.Func(func(in *parse.Input) (int, bool, error) {
    indent, ok, err := remainingInlineWhitespace.Parse(in)
    if err != nil || !ok {
        fmt.Println("at the indent")
        return 0, false, err
    }
    _, ok, err = parse.Rune('-').Parse(in)
    if err != nil || !ok {
        fmt.Println("at searching for -")
        return 0, false, err
    }
    _, ok, err = remainingInlineWhitespace.Parse(in)
    if err != nil || !ok {
        fmt.Println("at more whitespace")
        return 0, false, err
    }

    return len(indent), true, nil
})

var dueKey = parse.Any(parse.String("due:"), parse.String("d:"))
var scheduledKey = parse.Any(parse.String("scheduled:"), parse.String("s:"))
var completedKey = parse.Any(parse.String("completed:"))
var everyKey = parse.Any(parse.String("every:"), parse.String("e:"))
var priorityKey = parse.Any(parse.String("priority:"), parse.String("p:"))
var anyFieldKey = parse.Any(dueKey, scheduledKey, everyKey, priorityKey, completedKey)

var dueParser = parse.Func(func(in *parse.Input) (time.Time, bool, error) {
	_, ok, err := dueKey.Parse(in)
	if err != nil || !ok {
		return time.Time{}, false, err
	}
	return YearMonthDay.Parse(in)
})

var scheduledParser = parse.Func(func(in *parse.Input) (time.Time, bool, error) {
	_, ok, err := scheduledKey.Parse(in)
	if err != nil || !ok {
		return time.Time{}, false, err
	}
	return YearMonthDay.Parse(in)
})

var completedParser = parse.Func(func(in *parse.Input) (time.Time, bool, error) {
	_, ok, err := completedKey.Parse(in)
	if err != nil || !ok {
		return time.Time{}, false, err
	}
	return YearMonthDay.Parse(in)
})

var priorityParser = parse.Func(func(in *parse.Input) (int, bool, error) {
	_, ok, err := priorityKey.Parse(in)
	if err != nil || !ok {
		return 0, false, err
	}

	priority, ok, err := parse.StringFrom(parse.AtLeast(1, parse.ZeroToNine)).Parse(in)
	if err != nil || !ok {
		return 0, false, err
	}

	p, err := strconv.Atoi(priority)
	if err != nil {
		return 0, false, fmt.Errorf("invalid priority: %w", err)
	}
	return p, true, nil
})

var everyParser = func(relativeTo time.Time) parse.Parser[rrule.RRule] {
	return parse.Func(func(in *parse.Input) (rrule.RRule, bool, error) {
		_, ok, err := everyKey.Parse(in)
		if err != nil || !ok {
			return rrule.RRule{}, false, err
		}
		rruleOpts := rrule.ROption{Dtstart: relativeTo}

		// There are a limited number of single words that can be used to describe the frequency.
		// So lets get those out of the way first. (day, week, month, year, weekday, weekend)
		// Note that the order of these is important, as "week" is a prefix of "weekday" and "weekend".
		single, ok, err := parse.Any(day, parse.String("weekend"), parse.String("weekday"), month, year, week).Parse(in)
		if err != nil {
			return rrule.RRule{}, false, err
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
			rr, _ := rrule.NewRRule(rruleOpts)
			return *rr, true, nil
		}

		// Every <day of week> or list of <day of week>
		daysOfWeek, ok, err := DaysOfWeek.Parse(in)
		if err != nil {
			return rrule.RRule{}, false, err
		}
		if ok {
			for _, d := range daysOfWeek {
				rruleOpts.Byweekday = append(rruleOpts.Byweekday, rruleDayOfWeek(d))
			}
			rruleOpts.Freq = rrule.WEEKLY
			rr, _ := rrule.NewRRule(rruleOpts)
			return *rr, true, nil
		}

		// Every <number> <day/week/month/year>
		tuple, ok, err := parse.SequenceOf3(
			parse.StringFrom(parse.AtLeast(1, parse.ZeroToNine)),
			parse.String(" "),
			parse.Any(day, week, month, year),
		).Parse(in)
		if err != nil {
			return rrule.RRule{}, false, err
		}
		if ok {
			n, unit := tuple.A, tuple.C
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
			rr, _ := rrule.NewRRule(rruleOpts)
			return *rr, true, nil
		}

		// Some combination of month days and/or months
		// Keep reading input until we can't parse a month day or month.
		optionalDelimiter := parse.Optional(parse.Rune(' '))
		for {
			optionalDelimiter.Parse(in)
			monthDay, monthDayOk, err := MonthDay.Parse(in)
			if err != nil {
				return rrule.RRule{}, false, err
			}
			if monthDayOk {
				rruleOpts.Bymonthday = append(rruleOpts.Bymonthday, monthDay)
			}

			optionalDelimiter.Parse(in)
			month, monthOk, err := MonthOfYear.Parse(in)
			if err != nil {
				return rrule.RRule{}, false, err
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
			rr, err := rrule.NewRRule(rruleOpts)
			if err != nil {
				return rrule.RRule{}, false, err
			}
			return *rr, true, nil
		}

		return rrule.RRule{}, false, nil

	})
}

var Task = func(relativeTo time.Time) parse.Parser[api.Task] {
	return parse.Func(func(in *parse.Input) (api.Task, bool, error) {
		res := api.Task{Line: in.Position().Line}

		// Read and dump the list item open
		indent, ok, err := listItemOpen2.Parse(in)
		if err != nil || !ok {
			return api.Task{}, false, err
		}

		// Read the task status
		status, ok, err := statusParser.Parse(in)
		if err != nil || !ok {
			return api.Task{}, false, err
		}
		res.Status = status

		// Attempt to parse each of the fields resetting the input index each time.
		// Keep track of the shortest until string as that will be our name.
		start := in.Index()

		// Start name with the rest of the line. If we find a field (i.e. theres a shorter name) we'll use that.
		name, ok, err := parse.StringUntil(newLineOrEOF).Parse(in)
		if err != nil || !ok {
			return api.Task{}, false, err
		}
		in.Seek(start)

		// Due
		// We need to make sure the space is there to avoid matching on the single chars that match the end of a longer one.
		candidate, ok, err := parse.StringUntil(parse.StringFrom(parse.Rune(' '), dueKey)).Parse(in)
		if err != nil {
			return api.Task{}, false, err
		}
		name = evaluateCandidate(ok, candidate, name)
		if ok {
			parse.Rune(' ').Parse(in) // pop the space
			due, ok, err := dueParser.Parse(in)
			if err != nil || !ok {
				return api.Task{}, false, err
			}
			res.Due = &due
			in.Seek(start)
		}

		// Scheduled
		// We need to make sure the space is there to avoid matching on the single chars that match the end of a longer one.
		candidate, ok, err = parse.StringUntil(parse.StringFrom(parse.Rune(' '), scheduledKey)).Parse(in)
		if err != nil {
			return api.Task{}, false, err
		}
		name = evaluateCandidate(ok, candidate, name)
		if ok {
			parse.Rune(' ').Parse(in) // pop the space
			scheduled, ok, err := scheduledParser.Parse(in)
			if err != nil || !ok {
				return api.Task{}, false, err
			}
			res.Scheduled = &scheduled
			in.Seek(start)
		}

		// Completed
		// We need to make sure the space is there to avoid matching on the single chars that match the end of a longer one.
		candidate, ok, err = parse.StringUntil(parse.StringFrom(parse.Rune(' '), completedKey)).Parse(in)
		if err != nil {
			return api.Task{}, false, err
		}
		name = evaluateCandidate(ok, candidate, name)
		if ok {
			parse.Rune(' ').Parse(in) // pop the space
			completed, ok, err := completedParser.Parse(in)
			if err != nil || !ok {
				return api.Task{}, false, err
			}
			res.Completed = &completed
			in.Seek(start)
		}

		// Priority
		// We need to make sure the space is there to avoid matching on the single chars that match the end of a longer one.
		candidate, ok, err = parse.StringUntil(parse.StringFrom(parse.Rune(' '), priorityKey)).Parse(in)
		if err != nil {
			return api.Task{}, false, err
		}
		name = evaluateCandidate(ok, candidate, name)
		if ok {
			parse.Rune(' ').Parse(in) // pop the space
			priority, ok, err := priorityParser.Parse(in)
			if err != nil || !ok {
				return api.Task{}, false, err
			}
			res.Priority = &priority
			in.Seek(start)
		}

		// Every
		// We need to make sure the space is there to avoid matching on the single chars that match the end of a longer one.
		candidate, ok, err = parse.StringUntil(parse.StringFrom(parse.Rune(' '), everyKey)).Parse(in)
		if err != nil {
			return api.Task{}, false, err
		}
		name = evaluateCandidate(ok, candidate, name)
		if ok {
			parse.Rune(' ').Parse(in) // pop the space
			every, ok, err := everyParser(relativeTo).Parse(in)
			if err != nil || !ok {
				return api.Task{}, false, err
			}
			res.Every = &every
			in.Seek(start)
		}

		// Name
		res.Name = strings.TrimSpace(name)

        res.Indent = indent

		// Consume to the next line or eof.
		parse.StringUntil(newLineOrEOF).Parse(in)
		newLineOrEOF.Parse(in)

		return res, true, nil
	})
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

func evaluateCandidate(ok bool, candidate, name string) string {
	if !ok {
		return name
	}
	if len(candidate) < len(name) {
		return candidate
	}
	return name
}
