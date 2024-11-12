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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/parse"
	. "github.com/notedownorg/notedown/pkg/parsers"
	"github.com/teambition/rrule-go"
)

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
	_, longOk, err := dueKeyLong.Parse(in)
	if err != nil {
		return time.Time{}, false, err
	}
	_, shortOk, err := dueKeyShort.Parse(in)
	if err != nil {
		return time.Time{}, false, err
	}

	if !longOk && !shortOk {
		return time.Time{}, false, nil
	}
	return YearMonthDay.Parse(in)
})

var scheduledParser = parse.Func(func(in *parse.Input) (time.Time, bool, error) {
	_, longOk, err := scheduledKeyLong.Parse(in)
	if err != nil {
		return time.Time{}, false, err
	}
	_, shortOk, err := scheduledKeyShort.Parse(in)
	if err != nil {
		return time.Time{}, false, err
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
	_, longOk, err := priorityKeyLong.Parse(in)
	if err != nil {
		return 0, false, err
	}
	_, shortOk, err := priorityKeyShort.Parse(in)
	if err != nil {
		return 0, false, err
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

var everyParser = func(relativeTo time.Time) parse.Parser[Every] {
	return parse.Func(func(in *parse.Input) (Every, bool, error) {
		_, longOk, err := everyKeyLong.Parse(in)
		if err != nil {
			return Every{}, false, err
		}
		_, shortOk, err := everyKeyShort.Parse(in)
		if err != nil {
			return Every{}, false, err
		}

		if !longOk && !shortOk {
			return Every{}, false, nil
		}

		rruleOpts := rrule.ROption{Dtstart: relativeTo}

		// This closure keeps track of where we started so we can store the original text.
		buildResult := func() func(rrule.ROption, error) (Every, bool, error) {
			start := in.Index()
			return func(opts rrule.ROption, err error) (Every, bool, error) {
				if err != nil {
					return Every{}, false, err
				}
				rr, err := rrule.NewRRule(opts)
				if err != nil {
					return Every{}, false, err
				}

				// Get the text
				end := in.Index()
				in.Seek(start)
				text, ok := in.Take(end - start)
				if !ok {
					return Every{}, false, fmt.Errorf("failed to store original every text start: %d end: %d", start, end)
				}

				return Every{rrule: rr, text: strings.TrimSpace(text)}, true, nil
			}
		}()

		// There are a limited number of single words that can be used to describe the frequency.
		// So lets get those out of the way first. (day, week, month, year, weekday, weekend)
		// Note that the order of these is important, as "week" is a prefix of "weekday" and "weekend".
		single, ok, err := parse.Any(Day, parse.String("weekend"), parse.String("weekday"), Month, Year, Week).Parse(in)
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
			parse.Any(Day, Week, Month, Year),
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

		return Every{}, false, nil
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
