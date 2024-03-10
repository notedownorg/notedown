package parsers

import (
	"fmt"
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
