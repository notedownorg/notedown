package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/ast"
	"github.com/teambition/rrule-go"
)

var statusLookup = map[string]ast.Status{
	" ": ast.Todo,
	"x": ast.Done,
	"X": ast.Done,
	"/": ast.Doing,
	"b": ast.Blocked,
	"B": ast.Blocked,
	"a": ast.Abandoned,
	"A": ast.Abandoned,
}

var StatusRuneLookup = map[ast.Status]rune{
	ast.Todo:      ' ',
	ast.Blocked:   'b',
	ast.Doing:     '/',
	ast.Done:      'x',
	ast.Abandoned: 'a',
}

var statusParser = parse.Func(func(in *parse.Input) (ast.Status, bool, error) {
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

var (
	dueKeyLong  = parse.String("due:")
	dueKeyShort = parse.String("d:")
	dueKey      = parse.Any(dueKeyLong, dueKeyShort)

	scheduledKeyLong  = parse.String("scheduled:")
	scheduledKeyShort = parse.String("s:")
	scheduledKey      = parse.Any(scheduledKeyLong, scheduledKeyShort)

	completedKey = parse.Any(parse.String("completed:"))

	everyKeyLong  = parse.String("every:")
	everyKeyShort = parse.String("e:")
	everyKey      = parse.Any(everyKeyLong, everyKeyShort)

	priorityKeyLong  = parse.String("priority:")
	priorityKeyShort = parse.String("p:")
	priorityKey      = parse.Any(priorityKeyLong, priorityKeyShort)

	anyFieldKey = parse.Any(dueKey, scheduledKey, everyKey, priorityKey, completedKey)
)

var dueParser = parse.Func(func(in *parse.Input) (time.Time, bool, error) {
	inlineWhitespaceRunes.Parse(in) // dump any leading whitespace
	_, longOk, err := dueKeyLong.Parse(in)
	if err != nil {
		return time.Time{}, false, err
	}
	_, shortOk, err := dueKeyShort.Parse(in)
	if err != nil {
		return time.Time{}, false, err
	}

	// Ensure short key is not the end of a longer key.
	// Basically it either must be the start of the input or preceeded by a space.
	if shortOk {
		curr := in.Index()
		isStart := !in.Seek(curr - 3)
		start, _ := in.Peek(1)
		if !isStart && start != " " {
			return time.Time{}, false, nil
		}
		in.Seek(curr)
	}

	if !longOk && !shortOk {
		return time.Time{}, false, nil
	}
	return YearMonthDay.Parse(in)
})

var scheduledParser = parse.Func(func(in *parse.Input) (time.Time, bool, error) {
	inlineWhitespaceRunes.Parse(in) // dump any leading whitespace
	_, longOk, err := scheduledKeyLong.Parse(in)
	if err != nil {
		return time.Time{}, false, err
	}
	_, shortOk, err := scheduledKeyShort.Parse(in)
	if err != nil {
		return time.Time{}, false, err
	}

	// Ensure short key is not the end of a longer key.
	// Basically it either must be the start of the input or preceeded by a space.
	if shortOk {
		curr := in.Index()
		isStart := !in.Seek(curr - 3)
		start, _ := in.Peek(1)
		if !isStart && start != " " {
			return time.Time{}, false, nil
		}
		in.Seek(curr)
	}

	if !longOk && !shortOk {
		return time.Time{}, false, nil
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
	inlineWhitespaceRunes.Parse(in) // dump any leading whitespace
	_, longOk, err := priorityKeyLong.Parse(in)
	if err != nil {
		return 0, false, err
	}
	_, shortOk, err := priorityKeyShort.Parse(in)
	if err != nil {
		return 0, false, err
	}

	// Ensure short key is not the end of a longer key.
	// Basically it either must be the start of the input or preceeded by a space.
	if shortOk {
		curr := in.Index()
		isStart := !in.Seek(curr - 3)
		start, _ := in.Peek(1)
		if !isStart && start != " " {
			return 0, false, nil
		}
		in.Seek(curr)
	}

	if !longOk && !shortOk {
		return 0, false, nil
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

var everyParser = func(relativeTo time.Time) parse.Parser[ast.Every] {
	return parse.Func(func(in *parse.Input) (ast.Every, bool, error) {
		inlineWhitespaceRunes.Parse(in) // dump any leading inlineWhitespace
		_, longOk, err := everyKeyLong.Parse(in)
		if err != nil {
			return ast.Every{}, false, err
		}
		_, shortOk, err := everyKeyShort.Parse(in)
		if err != nil {
			return ast.Every{}, false, err
		}

		// Ensure short key is not the end of a longer key.
		// Basically it either must be the start of the input or preceeded by a space.
		if shortOk {
			curr := in.Index()
			isStart := !in.Seek(curr - 3)
			start, _ := in.Peek(1)
			if !isStart && start != " " {
				return ast.Every{}, false, nil
			}
			in.Seek(curr)
		}

		if !longOk && !shortOk {
			return ast.Every{}, false, nil
		}

		rruleOpts := rrule.ROption{Dtstart: relativeTo}

		// This closure keeps track of where we started so we can store the original text.
		buildResult := func() func(rrule.ROption, error) (ast.Every, bool, error) {
			start := in.Index()
			return func(opts rrule.ROption, err error) (ast.Every, bool, error) {
				if err != nil {
					return ast.Every{}, false, err
				}
				rr, err := rrule.NewRRule(opts)
				if err != nil {
					return ast.Every{}, false, err
				}

				// Get the text
				end := in.Index()
				in.Seek(start)
				text, ok := in.Take(end - start)
				if !ok {
					return ast.Every{}, false, fmt.Errorf("failed to store original every text start: %d end: %d", start, end)
				}

				return ast.Every{RRule: rr, Text: strings.TrimSpace(text)}, true, nil
			}
		}()

		// There are a limited number of single words that can be used to describe the frequency.
		// So lets get those out of the way first. (day, week, month, year, weekday, weekend)
		// Note that the order of these is important, as "week" is a prefix of "weekday" and "weekend".
		single, ok, err := parse.Any(day, parse.String("weekend"), parse.String("weekday"), month, year, week).Parse(in)
		if err != nil {
			return buildResult(rruleOpts, err)
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
			return buildResult(rruleOpts, nil)
		}

		// Every <day of week> or list of <day of week>
		daysOfWeek, ok, err := DaysOfWeek.Parse(in)
		if err != nil {
			return buildResult(rruleOpts, err)
		}
		if ok {
			for _, d := range daysOfWeek {
				rruleOpts.Byweekday = append(rruleOpts.Byweekday, rruleDayOfWeek(d))
			}
			rruleOpts.Freq = rrule.WEEKLY
			return buildResult(rruleOpts, nil)
		}

		// Every <number> <day/week/month/year>
		tuple, ok, err := parse.SequenceOf3(
			parse.StringFrom(parse.AtLeast(1, parse.ZeroToNine)),
			parse.String(" "),
			parse.Any(day, week, month, year),
		).Parse(in)
		if err != nil {
			return buildResult(rruleOpts, err)
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
			return buildResult(rruleOpts, nil)
		}

		// Some combination of month days and/or months
		// Keep reading input until we can't parse a month day or month.
		optionalDelimiter := parse.Optional(parse.Rune(' '))
		for {
			optionalDelimiter.Parse(in)
			monthDay, monthDayOk, err := MonthDay.Parse(in)
			if err != nil {
				return buildResult(rruleOpts, err)
			}
			if monthDayOk {
				rruleOpts.Bymonthday = append(rruleOpts.Bymonthday, monthDay)
			}

			optionalDelimiter.Parse(in)
			month, monthOk, err := MonthOfYear.Parse(in)
			if err != nil {
				return buildResult(rruleOpts, err)
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
			return buildResult(rruleOpts, nil)
		}

		return ast.Every{}, false, nil
	})
}

var Task = func(relativeTo time.Time) parse.Parser[ast.Task] {
	return parse.Func(func(in *parse.Input) (ast.Task, bool, error) {
		res := ast.Task{Line: in.Position().Line}

		// Read and dump the list item open
		_, ok, err := listItemOpen.Parse(in)
		if err != nil || !ok {
			return ast.Task{}, false, err
		}

		// Read the task status
		status, ok, err := statusParser.Parse(in)
		if err != nil || !ok {
			return ast.Task{}, false, err
		}
		res.Status = status

		// Read until we hit a key, newline or eof to get the name.
		name, ok, err := parse.StringUntil(parse.Any[string](anyFieldKey, newLineOrEOF)).Parse(in)
		if err != nil || !ok {
			return ast.Task{}, false, err
		}
		res.Name = strings.TrimSpace(name)

		// Parse the fields
		start := in.Index()

		// Due
		_, ok, err = parse.StringUntil(parse.Any[string](dueKey, newLineOrEOF)).Parse(in)
		if err != nil {
			return ast.Task{}, false, err
		}
		if ok {
			due, ok, err := dueParser.Parse(in)
			if err != nil {
				return ast.Task{}, false, err
			}
			if ok {
				res.Due = &due
			}
		}
		in.Seek(start)

		// Scheduled
		_, ok, err = parse.StringUntil(parse.Any[string](scheduledKey, newLineOrEOF)).Parse(in)
		if err != nil {
			return ast.Task{}, false, err
		}
		if ok {
			scheduled, ok, err := scheduledParser.Parse(in)
			if err != nil {
				return ast.Task{}, false, err
			}
			if ok {
				res.Scheduled = &scheduled
			}
		}
		in.Seek(start)

		// Completed
		_, ok, err = parse.StringUntil(parse.Any[string](completedKey, newLineOrEOF)).Parse(in)
		if err != nil {
			return ast.Task{}, false, err
		}
		if ok {
			completed, ok, err := completedParser.Parse(in)
			if err != nil {
				return ast.Task{}, false, err
			}
			if ok {
				res.Completed = &completed
			}
		}
		in.Seek(start)

		// Priority
		_, ok, err = parse.StringUntil(parse.Any[string](priorityKey, newLineOrEOF)).Parse(in)
		if err != nil {
			return ast.Task{}, false, err
		}
		if ok {
			priority, ok, err := priorityParser.Parse(in)
			if err != nil {
				return ast.Task{}, false, err
			}
			if ok {
				res.Priority = &priority
			}
		}
		in.Seek(start)

		// Every
		_, ok, err = parse.StringUntil(parse.Any[string](everyKey, newLineOrEOF)).Parse(in)
		if err != nil {
			return ast.Task{}, false, err
		}
		if ok {
			every, ok, err := everyParser(relativeTo).Parse(in)
			if err != nil {
				return ast.Task{}, false, err
			}
			if ok {
				res.Every = &every
			}
		}
		in.Seek(start)

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
