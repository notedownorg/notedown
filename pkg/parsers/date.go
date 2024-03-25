package parsers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/a-h/parse"
)

var YearMonthDay = parse.Func(func(in *parse.Input) (match time.Time, ok bool, err error) {
	// Create parsers for year, month and day.
	year := parse.StringFrom(parse.Times(4, parse.ZeroToNine))
	month := parse.StringFrom(parse.RuneIn("01"), parse.ZeroToNine)
	day := parse.StringFrom(parse.RuneIn("0123"), parse.ZeroToNine)

	// Create string parser for yyyy-MM-dd.
	date := parse.StringFrom(parse.All(year, parse.Rune('-'), month, parse.Rune('-'), day))

	s, ok, err := date.Parse(in)
	if err != nil || !ok {
		return time.Time{}, false, err
	}

	// Parse the date.
	match, err = time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("failed to parse date: %w", err)
	}

	return match, true, nil
})

// mon, monday, tue, tues, tuesday, wed, weds, wednesday, thu, thur, thurs, thursday, fri, friday, sat, saturday, sun, sunday
var DayOfWeek = parse.Func(func(in *parse.Input) (match time.Weekday, ok bool, err error) {
	m := map[time.Weekday]parse.Parser[string]{
		// Ordering is important! Use more specific parsers first.
		time.Monday:    parse.Any(parse.String("monday"), parse.String("mon")),
		time.Tuesday:   parse.Any(parse.String("tuesday"), parse.String("tues"), parse.String("tue")),
		time.Wednesday: parse.Any(parse.String("wednesday"), parse.String("weds"), parse.String("wed")),
		time.Thursday:  parse.Any(parse.String("thursday"), parse.String("thurs"), parse.String("thur"), parse.String("thu")),
		time.Friday:    parse.Any(parse.String("friday"), parse.String("fri")),
		time.Saturday:  parse.Any(parse.String("saturday"), parse.String("sat")),
		time.Sunday:    parse.Any(parse.String("sunday"), parse.String("sun")),
	}

	for day, parser := range m {
		_, ok, err := parser.Parse(in)
		if err != nil {
			return time.Weekday(-1), false, err
		}
		if ok {
			return day, true, nil
		}
	}

	return time.Weekday(-1), false, nil
})

// Space separated list of days of the week.
var DaysOfWeek = parse.Func(func(in *parse.Input) (match []time.Weekday, ok bool, err error) {
	var days []time.Weekday

	delimiter := parse.RuneIn(" ")

	for {
		day, ok, err := DayOfWeek.Parse(in)
		if err != nil {
			return nil, false, err
		}
		if !ok {
			break
		}
		_, ok, err = delimiter.Parse(in)
		if err != nil {
			return nil, false, err
		}
		days = append(days, day)
	}

	if len(days) == 0 {
		return nil, false, nil
	}

	return days, true, nil
})

// Parse a number followed by an optional ordinal (st, nd, rd, th).
var MonthDay = parse.Func(func(in *parse.Input) (match int, ok bool, err error) {
	// Parse a number.
	n, ok, err := parse.StringFrom(parse.AtLeast(1, parse.ZeroToNine)).Parse(in)
	if err != nil {
		return 0, false, err
	}
	if !ok {
		return 0, false, nil
	}

	// Parse an optional ordinal.
	ordinal := parse.Any(parse.String("st"), parse.String("nd"), parse.String("rd"), parse.String("th"))
	_, _, err = ordinal.Parse(in)
	if err != nil {
		return 0, false, err
	}

	// Convert the number to an integer.
	number, err := strconv.Atoi(n)
	if err != nil {
		return 0, false, fmt.Errorf("failed to parse number: %w", err)
	}

	return number, true, nil
})

var MonthOfYear = parse.Func(func(in *parse.Input) (match time.Month, ok bool, err error) {
	m := map[time.Month]parse.Parser[string]{
		// Ordering is important! Use more specific parsers first.
		time.January:   parse.Any(parse.String("january"), parse.String("jan")),
		time.February:  parse.Any(parse.String("february"), parse.String("feb")),
		time.March:     parse.Any(parse.String("march"), parse.String("mar")),
		time.April:     parse.Any(parse.String("april"), parse.String("apr")),
		time.May:       parse.String("may"),
		time.June:      parse.Any(parse.String("june"), parse.String("jun")),
		time.July:      parse.Any(parse.String("july"), parse.String("jul")),
		time.August:    parse.Any(parse.String("august"), parse.String("aug")),
		time.September: parse.Any(parse.String("september"), parse.String("sept"), parse.String("sep")),
		time.October:   parse.Any(parse.String("october"), parse.String("oct")),
		time.November:  parse.Any(parse.String("november"), parse.String("nov")),
		time.December:  parse.Any(parse.String("december"), parse.String("dec")),
	}

	for month, parser := range m {
		_, ok, err := parser.Parse(in)
		if err != nil {
			return time.Month(-1), false, err
		}
		if ok {
			return month, true, nil
		}
	}

	return 0, false, nil
})
