// Copyright 2024 Notedown Authors
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

package tasks

import (
	"strings"
	"time"

	"github.com/a-h/parse"
	. "github.com/notedownorg/notedown/pkg/parsers"
)

var statusLookup = map[string]Status{
	" ": Todo,
	"x": Done,
	"X": Done,
	"/": Doing,
	"b": Blocked,
	"B": Blocked,
	"a": Abandoned,
	"A": Abandoned,
}

var StatusRuneLookup = map[Status]rune{
	Todo:      ' ',
	Blocked:   'b',
	Doing:     '/',
	Done:      'x',
	Abandoned: 'a',
}

var statusParser = parse.Func(func(in *parse.Input) (Status, bool, error) {
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

var listItemOpen = parse.StringFrom(RemainingInlineWhitespace, parse.Rune('-'), RemainingInlineWhitespace)

var ParseTask = func(path string, checksum string, relativeTo time.Time) parse.Parser[Task] {
	return parse.Func(func(in *parse.Input) (Task, bool, error) {
		// Line is 1-indexed not 0-indexed, this is so it's a bit more user friendly and also to allow for 0 to represent the beginning of the file.
		line, taskOpts := in.Position().Line+1, []TaskOption{}

		// Read and dump the list item open
		_, ok, err := listItemOpen.Parse(in)
		if err != nil || !ok {
			return Task{}, false, err
		}

		// Read the task status
		status, ok, err := statusParser.Parse(in)
		if err != nil || !ok {
			return Task{}, false, err
		}

		// Read until we hit a key, newline or eof to get the name.
		name, ok, err := parse.StringUntil(parse.Any(LeadingWhitespace(anyFieldKey), NewLineOrEOF)).Parse(in)
		if err != nil || !ok {
			return Task{}, false, err
		}
		name = strings.TrimSpace(name)

		// Parse the fields
		start := in.Index()

		// Due
		_, ok, err = parse.StringUntil(parse.Any(LeadingWhitespace(dueKey), NewLineOrEOF)).Parse(in)
		if err != nil {
			return Task{}, false, err
		}
		if ok {
			due, ok, err := LeadingWhitespace(dueParser).Parse(in)
			if err != nil {
				return Task{}, false, err
			}
			if ok {
				taskOpts = append(taskOpts, WithDue(due))
			}
		}
		in.Seek(start)

		// Scheduled
		_, ok, err = parse.StringUntil(parse.Any(LeadingWhitespace(scheduledKey), NewLineOrEOF)).Parse(in)
		if err != nil {
			return Task{}, false, err
		}
		if ok {
			scheduled, ok, err := LeadingWhitespace(scheduledParser).Parse(in)
			if err != nil {
				return Task{}, false, err
			}
			if ok {
				taskOpts = append(taskOpts, WithScheduled(scheduled))
			}
		}
		in.Seek(start)

		// Completed
		_, ok, err = parse.StringUntil(parse.Any(LeadingWhitespace(completedKey), NewLineOrEOF)).Parse(in)
		if err != nil {
			return Task{}, false, err
		}
		if ok {
			completed, ok, err := LeadingWhitespace(completedParser).Parse(in)
			if err != nil {
				return Task{}, false, err
			}
			if ok {
				taskOpts = append(taskOpts, WithCompleted(completed))
			}
		}
		in.Seek(start)

		// Priority
		_, ok, err = parse.StringUntil(parse.Any(LeadingWhitespace(priorityKey), NewLineOrEOF)).Parse(in)
		if err != nil {
			return Task{}, false, err
		}
		if ok {
			priority, ok, err := LeadingWhitespace(priorityParser).Parse(in)
			if err != nil {
				return Task{}, false, err
			}
			if ok {
				taskOpts = append(taskOpts, WithPriority(priority))
			}
		}
		in.Seek(start)

		// Every
		_, ok, err = parse.StringUntil(parse.Any(LeadingWhitespace(everyKey), NewLineOrEOF)).Parse(in)
		if err != nil {
			return Task{}, false, err
		}
		if ok {
			every, ok, err := LeadingWhitespace(everyParser(relativeTo)).Parse(in)
			if err != nil {
				return Task{}, false, err
			}
			if ok {
				taskOpts = append(taskOpts, WithEvery(every))
			}
		}
		in.Seek(start)

		// Consume to the next line or eof.
		parse.StringUntil(NewLineOrEOF).Parse(in)
		NewLineOrEOF.Parse(in)

		return NewTask(NewIdentifier(path, checksum, line), name, status, taskOpts...), true, nil
	})
}
