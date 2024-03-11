package parsers

import "github.com/a-h/parse"

var (
    // More specific first.
    day = parse.Any(parse.String("days"), parse.String("day"))
    week = parse.Any(parse.String("weeks"), parse.String("week"))
    month = parse.Any(parse.String("months"), parse.String("month"))
    year = parse.Any(parse.String("years"), parse.String("year"))
)
