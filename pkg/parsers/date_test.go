package parsers_test

import (
	"testing"
	"time"

	"github.com/a-h/parse"
	"github.com/stretchr/testify/assert"

	"github.com/liamawhite/nl/pkg/parsers"
)

func TestYearMonthDay(t *testing.T) {
    tests := []struct {
        input string
        want  time.Time
        notFound bool
    }{
        {
            input: "2021-01-01",
            want:  time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
        },
        {
            input: "2021-01-32",
            notFound: true,
        },
        {
            input: "2021-13-01",
            notFound: true,
        },
        {
            input: "2021-00-01",
            notFound: true,
        },
        {
            input: "2021-01-00",
            notFound: true,
        },
        {
            input: "2021-02-29", // not a leap year
            notFound: true,
        },
        {
            input: "2024-02-29", // valid leap year
            want:  time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
        },
    }
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            input := parse.NewInput(tt.input)
            got, found, _ := parsers.YearMonthDay.Parse(input)
            if tt.notFound {
                if found {
                    t.Fatalf("expected not found, got %v", got)
                }
                return
            }
            if !found {
                t.Fatalf("expected found, got not found")
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
