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
	"testing"
	"time"

	"github.com/a-h/parse"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/rrule-go"
)

func date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

func intPtr(i int) *int {
	return &i
}

var relativeTo, _ = time.Parse(time.RFC3339, "2020-01-02T00:00:00Z") // thurs
var dailyRule, _ = rrule.NewRRule(rrule.ROption{Freq: rrule.DAILY, Dtstart: relativeTo})
var spacesRule, _ = rrule.NewRRule(rrule.ROption{Freq: rrule.WEEKLY, Dtstart: relativeTo, Byweekday: []rrule.Weekday{rrule.MO, rrule.WE, rrule.FR}})

func TestTask(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      Task
		leftOverInput bool
	}{
		// Test each status
		{
			name:     "Todo",
			input:    "- [ ] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Todo),
		},
		{
			name:     "Done (lowercase)",
			input:    "- [x] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Done),
		},
		{
			name:     "Done (uppercase)",
			input:    "- [X] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Done),
		},
		{
			name:     "Doing",
			input:    "- [/] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Doing),
		},
		{
			name:     "Blocked (lowercase)",
			input:    "- [b] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Blocked),
		},
		{
			name:     "Blocked (uppercase)",
			input:    "- [B] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Blocked),
		},
		{
			name:     "Abandoned (lowercase)",
			input:    "- [a] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Abandoned),
		},
		{
			name:     "Abandoned (uppercase)",
			input:    "- [A] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Abandoned),
		},
		// Whitespace tests
		{
			name:     "Leading space",
			input:    " - [ ] Task",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Todo),
		},
		{
			name:     "Trailing space",
			input:    "- [ ] Task ",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Todo),
		},
		{
			name:     "Task name with spaces",
			input:    "- [ ] Task Name",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task Name", Todo),
		},
		{
			name:     "Task name with lots of random spaces",
			input:    "          - [ ]   Task   Name   ",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task   Name", Todo),
		},
		// Fields
		{
			name:     "Due date",
			input:    "- [ ] Task due:2021-01-01",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Todo, WithDue(date(2021, 1, 1))),
		},
		{
			name:          "Due date on different task",
			input:         "- [ ] Task 1\n- [ ] Task 2 due:2021-01-01",
			expected:      NewTask(NewIdentifier("path", "version", 1), "Task 1", Todo),
			leftOverInput: true,
		},
		{
			name:     "Scheduled date",
			input:    "- [ ] Task scheduled:2021-01-01",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Todo, WithScheduled(date(2021, 1, 1))),
		},
		{
			name:          "Scheduled date on different task",
			input:         "- [ ] Task 1\n- [ ] Task 2 scheduled:2021-01-01",
			expected:      NewTask(NewIdentifier("path", "version", 1), "Task 1", Todo),
			leftOverInput: true,
		},
		{
			name:     "Completed date",
			input:    "- [ ] Task completed:2021-01-01",
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Todo, WithCompleted(date(2021, 1, 1))),
		},
		{
			name:          "Completed date on different task",
			input:         "- [ ] Task 1\n- [ ] Task 2 completed:2021-01-01",
			expected:      NewTask(NewIdentifier("path", "version", 1), "Task 1", Todo),
			leftOverInput: true,
		},
		{
			name:     "Conflicting short and long fields",
			input:    "- [ ] Task scheduled:2021-01-01 completed:2021-01-02", // both end in d: so make sure theres no due date
			expected: NewTask(NewIdentifier("path", "version", 1), "Task", Todo, WithScheduled(date(2021, 1, 1)), WithCompleted(date(2021, 1, 2))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := ParseTask("path", "version", relativeTo).Parse(in)
			if !found {
				t.Fatal("expected found")
			}
			assert.Equal(t, test.expected, result)
			if test.leftOverInput {
				assert.NotEqual(t, len(test.input), in.Index(), "expected there to be leftover input")
			} else {
				assert.Equal(t, len(test.input), in.Index(), "expected to consume the entire input")
			}
		})
	}
}

func TestDueDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		notFound bool
	}{
		{
			name:     "Long",
			input:    "due:2021-01-01",
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Short",
			input:    "d:2021-01-01",
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Short with leading space",
			input:    " d:2021-01-01",
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := dueParser.Parse(in)
			if !found {
				t.Fatal("expected found")
			}
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestScheduledDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "Long",
			input:    "scheduled:2021-01-01",
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Short",
			input:    "s:2021-01-01",
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Short with leading space",
			input:    " s:2021-01-01",
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := scheduledParser.Parse(in)
			if !found {
				t.Fatal("expected found")
			}
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestCompletedDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "Long",
			input:    "completed:2021-01-01",
			expected: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := completedParser.Parse(in)
			if !found {
				t.Fatal("expected found")
			}
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestPriority(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		notFound bool
	}{
		{
			name:     "Long",
			input:    "priority:1",
			expected: 1,
		},
		{
			name:     "Short",
			input:    "p:1",
			expected: 1,
		},
		{
			name:     "Short with leading space",
			input:    " p:1",
			expected: 1,
		},
		{
			name:     "Zero",
			input:    "p:0",
			expected: 0,
		},
		{
			name:     "Nine",
			input:    "p:9",
			expected: 9,
		},
		{
			name:     "Ten",
			input:    "p:10",
			expected: 10,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := priorityParser.Parse(in)
			if test.notFound {
				if found {
					t.Fatal("expected not found")
				}
				return
			}
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestEvery(t *testing.T) {
	tests := []struct {
		input        string
		expected     []time.Time
		expectedText string
		end          time.Time
		notFound     bool
	}{
		{
			input:        "every:day",
			expectedText: "day",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "e:day",
			expectedText: "day",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        " e:day",
			expectedText: "day",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:week",
			expectedText: "week",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 16, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 23, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 20, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 2, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:month",
			expectedText: "month",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 4, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 5, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 6, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 7, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 8, 2, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 8, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:year",
			expectedText: "year",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2027, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:weekday",
			expectedText: "weekday",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 13, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 1, 13, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:weekend",
			expectedText: "weekend",
			expected: []time.Time{
				time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 11, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 18, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 25, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 22, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 2, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:mon",
			expectedText: "mon",
			expected: []time.Time{
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 3, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 17, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 24, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 2, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:mon tues wed",
			expectedText: "mon tues wed",
			expected: []time.Time{
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 14, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 20, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 21, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 1, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:tues thurs",
			expectedText: "tues thurs",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 9, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 14, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 16, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 21, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 23, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 28, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:2 days",
			expectedText: "2 days",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 12, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 14, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 16, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 1, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:2 weeks",
			expectedText: "2 weeks",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 16, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 13, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 27, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 12, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 26, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 4, 9, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 4, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:18 months",
			expectedText: "18 months",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 7, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 7, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2027, 7, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2029, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2030, 7, 2, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2030, 7, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:3 years",
			expectedText: "3 years",
			expected: []time.Time{
				time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2029, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2032, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2035, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2038, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2041, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2041, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:jan mar sept",
			expectedText: "jan mar sept",
			expected: []time.Time{
				time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2022, 9, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:1st 15th jan sept",
			expectedText: "1st 15th jan sept",
			expected: []time.Time{
				time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 9, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 9, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 9, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 9, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:1 15",
			expectedText: "1 15",
			expected: []time.Time{
				time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 3, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 4, 15, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2020, 5, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:        "every:1 30 jan february",
			expectedText: "1 30 jan february",
			expected: []time.Time{
				time.Date(2020, 1, 30, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 1, 30, 0, 0, 0, 0, time.UTC),
				time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 1, 30, 0, 0, 0, 0, time.UTC),
				time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			end: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, err := everyParser(relativeTo).Parse(in)
			assert.NoError(t, err)
			if test.notFound {
				if found {
					t.Fatal("expected not found")
				}
				return
			}
			if !found {
				t.Fatal("expected found")
			}
			assert.Equal(t, test.expected, result.rrule.Between(relativeTo, test.end, true))
			assert.Equal(t, test.expectedText, result.text)
		})
	}
}
