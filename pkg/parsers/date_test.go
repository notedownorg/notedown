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

package parsers_test

import (
	"testing"
	"time"

	"github.com/a-h/parse"
	"github.com/stretchr/testify/assert"

	"github.com/notedownorg/notedown/pkg/parsers"
)

func TestYearMonthDay(t *testing.T) {
	tests := []struct {
		input    string
		want     time.Time
		notFound bool
	}{
		{
			input: "2021-01-01",
			want:  time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "2021-01-32",
			notFound: true,
		},
		{
			input:    "2021-13-01",
			notFound: true,
		},
		{
			input:    "2021-00-01",
			notFound: true,
		},
		{
			input:    "2021-01-00",
			notFound: true,
		},
		{
			input:    "2021-02-29", // not a leap year
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

func TestDayOfWeek(t *testing.T) {
	tests := []struct {
		input    string
		want     time.Weekday
		notFound bool
	}{
		{
			input: "monday",
			want:  time.Monday,
		},
		{
			input: "mon",
			want:  time.Monday,
		},
		{
			input: "tuesday",
			want:  time.Tuesday,
		},
		{
			input: "tues",
			want:  time.Tuesday,
		},
		{
			input: "tue",
			want:  time.Tuesday,
		},
		{
			input: "wednesday",
			want:  time.Wednesday,
		},
		{
			input: "weds",
			want:  time.Wednesday,
		},
		{
			input: "wed",
			want:  time.Wednesday,
		},
		{
			input: "thursday",
			want:  time.Thursday,
		},
		{
			input: "thurs",
			want:  time.Thursday,
		},
		{
			input: "thur",
			want:  time.Thursday,
		},
		{
			input: "thu",
			want:  time.Thursday,
		},
		{
			input: "friday",
			want:  time.Friday,
		},
		{
			input: "fri",
			want:  time.Friday,
		},
		{
			input: "saturday",
			want:  time.Saturday,
		},
		{
			input: "sat",
			want:  time.Saturday,
		},
		{
			input: "sunday",
			want:  time.Sunday,
		},
		{
			input: "sun",
			want:  time.Sunday,
		},
		{
			input:    "foo",
			notFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			got, found, _ := parsers.DayOfWeek.Parse(input)
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
			assert.Equal(t, len(tt.input), input.Index(), "expected input to be consumed")
		})
	}
}

func TestDaysOfWeek(t *testing.T) {
	tests := []struct {
		input    string
		want     []time.Weekday
		notFound bool
	}{
		{
			input: "monday",
			want:  []time.Weekday{time.Monday},
		},
		{
			input: "mon tues wed thu fri sat sunday",
			want:  []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
		},
		{
			input:    "not a day",
			notFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			got, found, _ := parsers.DaysOfWeek.Parse(input)
			if tt.notFound {
				if found {
					t.Fatalf("expected not found, got %v", got)
				}
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, len(tt.input), input.Index(), "expected input to be consumed")
		})
	}
}

func TestMonthDay(t *testing.T) {
	tests := []struct {
		input    string
		want     int
		notFound bool
	}{
		{
			input: "1st",
			want:  1,
		},
		{
			input: "1",
			want:  1,
		},
		{
			input: "2nd",
			want:  2,
		},
		{
			input: "3rd",
			want:  3,
		},
		{
			input: "4th",
			want:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			got, found, _ := parsers.MonthDay.Parse(input)
			if tt.notFound {
				if found {
					t.Fatalf("expected not found, got %v", got)
				}
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, len(tt.input), input.Index(), "expected input to be consumed")
		})
	}
}

func TestMonthOfYear(t *testing.T) {
	tests := []struct {
		input    string
		want     time.Month
		notFound bool
	}{
		{
			input: "january",
			want:  time.January,
		},
		{
			input: "jan",
			want:  time.January,
		},
		{
			input: "february",
			want:  time.February,
		},
		{
			input: "feb",
			want:  time.February,
		},
		{
			input: "march",
			want:  time.March,
		},
		{
			input: "mar",
			want:  time.March,
		},
		{
			input: "april",
			want:  time.April,
		},
		{
			input: "apr",
			want:  time.April,
		},
		{
			input: "may",
			want:  time.May,
		},
		{
			input: "june",
			want:  time.June,
		},
		{
			input: "jun",
			want:  time.June,
		},
		{
			input: "july",
			want:  time.July,
		},
		{
			input: "jul",
			want:  time.July,
		},
		{
			input: "august",
			want:  time.August,
		},
		{
			input: "aug",
			want:  time.August,
		},
		{
			input: "september",
			want:  time.September,
		},
		{
			input: "sep",
			want:  time.September,
		},
		{
			input: "october",
			want:  time.October,
		},
		{
			input: "oct",
			want:  time.October,
		},
		{
			input: "november",
			want:  time.November,
		},
		{
			input: "nov",
			want:  time.November,
		},
		{
			input: "december",
			want:  time.December,
		},
		{
			input: "dec",
			want:  time.December,
		},
		{
			input:    "foo",
			notFound: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			got, found, _ := parsers.MonthOfYear.Parse(input)
			if tt.notFound {
				if found {
					t.Fatalf("expected not found, got %v", got)
				}
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, len(tt.input), input.Index(), "expected input to be consumed")
		})
	}
}
