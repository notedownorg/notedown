package parsers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/api"
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

var dueKey = parse.Any(parse.String("due:"), parse.String("d:"))
var scheduledKey = parse.Any(parse.String("scheduled:"), parse.String("s:"))
var everyKey = parse.Any(parse.String("every:"), parse.String("e:"))
var priorityKey = parse.Any(parse.String("priority:"), parse.String("p:"))
var anyFieldKey = parse.Any(dueKey, scheduledKey, everyKey, priorityKey)

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

var priorityParser = parse.Func(func(in *parse.Input) (int, bool, error) {
	_, ok, err := priorityKey.Parse(in)
	if err != nil || !ok {
		return 0, false, err
	}

	priority, ok, err := parse.StringFrom(parse.ZeroToNine).Parse(in)
	if err != nil || !ok {
		return 0, false, err
	}
	p, err := strconv.Atoi(priority)
	if err != nil {
		return 0, false, fmt.Errorf("invalid priority: %w", err)
	}
	return p, true, nil
})

var Task = parse.Func(func(in *parse.Input) (api.Task, bool, error) {
	// Read and dump the list item open
	_, ok, err := listItemOpen.Parse(in)
	if err != nil || !ok {
		return api.Task{}, false, err
	}

	// Read the task status
	status, ok, err := statusParser.Parse(in)
	if err != nil || !ok {
		return api.Task{}, false, err
	}

	// Attempt to parse each of the fields resetting the input index each time.
	// Keep track of the shortest until string as that will be our name.
	res := api.Task{Status: status}
	start := in.Index()

	// Start name with the rest of the line. If we find a field (i.e. theres a shorter name) we'll use that.
	name, ok, err := parse.StringUntil(parse.NewLine).Parse(in)
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

	// Name
	res.Name = strings.TrimSpace(name)

    // Consume to the rest of the line.
    parse.StringUntil(parse.NewLine).Parse(in)
    parse.NewLine.Parse(in)

	return res, true, nil
})

func evaluateCandidate(ok bool, candidate, name string) string {
	if !ok {
		return name
	}
	if len(candidate) < len(name) {
		return candidate
	}
	return name
}
