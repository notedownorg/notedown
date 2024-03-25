package parsers

import (
	"testing"
	"time"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/rrule-go"
)

func date(y, m, d int) *time.Time {
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return &t
}

func intPtr(i int) *int {
	return &i
}

var relativeTo, _ = time.Parse(time.RFC3339, "2020-01-02T00:00:00Z") // thurs
var dailyRule, _ = rrule.NewRRule(rrule.ROption{Freq: rrule.DAILY, Dtstart: relativeTo})
var spacesRule, _ = rrule.NewRRule(rrule.ROption{Freq: rrule.WEEKLY, Dtstart: relativeTo, Byweekday: []rrule.Weekday{rrule.MO, rrule.WE, rrule.FR}})

func TestTask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected api.Task
	}{
		// Test each status
		{
			name:     "Todo",
			input:    "- [ ] Task",
			expected: api.Task{Status: api.Todo, Name: "Task"},
		},
		{
			name:     "Done (lowercase)",
			input:    "- [x] Task",
			expected: api.Task{Status: api.Done, Name: "Task"},
		},
		{
			name:     "Done (uppercase)",
			input:    "- [X] Task",
			expected: api.Task{Status: api.Done, Name: "Task"},
		},
		{
			name:     "Doing",
			input:    "- [/] Task",
			expected: api.Task{Status: api.Doing, Name: "Task"},
		},
		{
			name:     "Blocked (lowercase)",
			input:    "- [b] Task",
			expected: api.Task{Status: api.Blocked, Name: "Task"},
		},
		{
			name:     "Blocked (uppercase)",
			input:    "- [B] Task",
			expected: api.Task{Status: api.Blocked, Name: "Task"},
		},
		{
			name:     "Abandoned (lowercase)",
			input:    "- [a] Task",
			expected: api.Task{Status: api.Abandoned, Name: "Task"},
		},
		{
			name:     "Abandoned (uppercase)",
			input:    "- [A] Task",
			expected: api.Task{Status: api.Abandoned, Name: "Task"},
		},
		// Whitespace tests
		{
			name:     "Leading space",
			input:    " - [ ] Task",
			expected: api.Task{Status: api.Todo, Name: "Task", Indent: 1},
		},
		{
			name:     "Trailing space",
			input:    "- [ ] Task ",
			expected: api.Task{Status: api.Todo, Name: "Task"},
		},
		{
			name:     "Task name with spaces",
			input:    "- [ ] Task Name",
			expected: api.Task{Status: api.Todo, Name: "Task Name"},
		},
		{
			name:     "Task name with lots of random spaces",
			input:    "          - [ ]   Task   Name   ",
			expected: api.Task{Status: api.Todo, Name: "Task   Name", Indent: 10},
		},
		// Fields
		{
			name:     "Due date",
			input:    "- [ ] Task due:2021-01-01",
			expected: api.Task{Status: api.Todo, Name: "Task", Due: date(2021, 1, 1)},
		},
		{
			name:     "Scheduled date",
			input:    "- [ ] Task scheduled:2021-01-01",
			expected: api.Task{Status: api.Todo, Name: "Task", Scheduled: date(2021, 1, 1)},
		},
		{
			name:     "Completed date",
			input:    "- [ ] Task completed:2021-01-01",
			expected: api.Task{Status: api.Todo, Name: "Task", Completed: date(2021, 1, 1)},
		},
		{
			name:     "Priority",
			input:    "- [ ] Task priority:1",
			expected: api.Task{Status: api.Todo, Name: "Task", Priority: intPtr(1)},
		},
		{
			name:     "Every",
			input:    "- [ ] Task every:day",
			expected: api.Task{Status: api.Todo, Name: "Task", Every: dailyRule},
		},
		{
			name:     "Every with spaces",
			input:    "- [ ] Task every:mon wed fri",
			expected: api.Task{Status: api.Todo, Name: "Task", Every: spacesRule},
		},
		{
			name:     "All fields long",
			input:    "- [ ] Task due:2021-01-01 every:mon wed fri scheduled:2021-01-01 completed:2021-01-01 priority:1",
			expected: api.Task{Status: api.Todo, Name: "Task", Due: date(2021, 1, 1), Scheduled: date(2021, 1, 1), Completed: date(2021, 1, 1), Priority: intPtr(1), Every: spacesRule},
		},
		{
			name:     "All fields short",
			input:    "- [ ] Task d:2021-01-01 e:mon wed fri s:2021-01-01 p:1 completed:2021-01-01",
			expected: api.Task{Status: api.Todo, Name: "Task", Due: date(2021, 1, 1), Scheduled: date(2021, 1, 1), Completed: date(2021, 1, 1), Priority: intPtr(1), Every: spacesRule},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := Task(relativeTo).Parse(in)
			if !found {
				t.Fatal("expected found")
			}
			assert.Equal(t, test.expected, result)
			assert.Equal(t, len(test.input), in.Index(), "expected to consume the entire input")
		})
	}
}

func TestDueDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
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
		input    string
		expected []time.Time
		end      time.Time
		notFound bool
	}{
		{
			input: "every:day",
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
			input: "e:day",
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
			input: "every:week",
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
			input: "every:month",
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
			input: "every:year",
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
			input: "every:weekday",
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
			input: "every:weekend",
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
			input: "every:mon",
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
			input: "every:mon tues wed",
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
			input: "every:tues thurs",
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
			input: "every:2 days",
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
			input: "every:2 weeks",
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
			input: "every:18 months",
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
			input: "every:3 years",
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
			input: "every:jan mar sept",
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
			input: "every:1st 15th jan sept",
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
			input: "every:1 15",
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
			input: "every:1 30 jan february",
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
			result, found, _ := everyParser(relativeTo).Parse(in)
			if test.notFound {
				if found {
					t.Fatal("expected not found")
				}
				return
			}
			if !found {
				t.Fatal("expected found")
			}
			assert.Equal(t, test.expected, result.Between(relativeTo, test.end, true))
		})
	}
}
