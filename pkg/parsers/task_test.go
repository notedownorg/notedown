package parsers

import (
	"testing"
	"time"

	"github.com/a-h/parse"
	"github.com/liamawhite/nl/pkg/api"
	"github.com/stretchr/testify/assert"
)

func date(y, m, d int) *time.Time {
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	return &t
}

func intPtr(i int) *int {
	return &i
}

func TestTask(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected api.Task
	}{
		// Test each status
		{
			name:     "Todo",
			input:    "- [ ] Task\n",
			expected: api.Task{Status: api.Todo, Name: "Task"},
		},
		{
			name:     "Done (lowercase)",
			input:    "- [x] Task\n",
			expected: api.Task{Status: api.Done, Name: "Task"},
		},
		{
			name:     "Done (uppercase)",
			input:    "- [X] Task\n",
			expected: api.Task{Status: api.Done, Name: "Task"},
		},
		{
			name:     "Doing",
			input:    "- [/] Task\n",
			expected: api.Task{Status: api.Doing, Name: "Task"},
		},
		{
			name:     "Blocked (lowercase)",
			input:    "- [b] Task\n",
			expected: api.Task{Status: api.Blocked, Name: "Task"},
		},
		{
			name:     "Blocked (uppercase)",
			input:    "- [B] Task\n",
			expected: api.Task{Status: api.Blocked, Name: "Task"},
		},
		{
			name:     "Abandoned (lowercase)",
			input:    "- [a] Task\n",
			expected: api.Task{Status: api.Abandoned, Name: "Task"},
		},
		{
			name:     "Abandoned (uppercase)",
			input:    "- [A] Task\n",
			expected: api.Task{Status: api.Abandoned, Name: "Task"},
		},
		// Whitespace tests
		{
			name:     "Leading space",
			input:    " - [ ] Task\n",
			expected: api.Task{Status: api.Todo, Name: "Task"},
		},
		{
			name:     "Trailing space",
			input:    "- [ ] Task \n",
			expected: api.Task{Status: api.Todo, Name: "Task"},
		},
		{
			name:     "Task name with spaces",
			input:    "- [ ] Task Name\n",
			expected: api.Task{Status: api.Todo, Name: "Task Name"},
		},
		{
			name:     "Task name with lots of random spaces",
			input:    "          - [ ]   Task   Name   \n",
			expected: api.Task{Status: api.Todo, Name: "Task   Name"},
		},
		// Fields
		{
			name:     "Due date",
			input:    "- [ ] Task due:2021-01-01\n",
			expected: api.Task{Status: api.Todo, Name: "Task", Due: date(2021, 1, 1)},
		},
		{
			name:     "Scheduled date",
			input:    "- [ ] Task scheduled:2021-01-01\n",
			expected: api.Task{Status: api.Todo, Name: "Task", Scheduled: date(2021, 1, 1)},
		},
		{
			name:     "Priority",
			input:    "- [ ] Task priority:1\n",
			expected: api.Task{Status: api.Todo, Name: "Task", Priority: intPtr(1)},
		},
        {
            name:     "All fields long",
            input:    "- [ ] Task due:2021-01-01 scheduled:2021-01-01 priority:1\n",
            expected: api.Task{Status: api.Todo, Name: "Task", Due: date(2021, 1, 1), Scheduled: date(2021, 1, 1), Priority: intPtr(1)},
        },
        {
            name:     "All fields short",
            input:    "- [ ] Task d:2021-01-01 s:2021-01-01 p:1\n",
            expected: api.Task{Status: api.Todo, Name: "Task", Due: date(2021, 1, 1), Scheduled: date(2021, 1, 1), Priority: intPtr(1)},
        },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := parse.NewInput(test.input)
			result, found, _ := Task.Parse(in)
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

func TestPriority(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected int
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
    }
    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            in := parse.NewInput(test.input)
            result, found, _ := priorityParser.Parse(in)
            if !found {
                t.Fatal("expected found")
            }
            assert.Equal(t, test.expected, result)
        })
    }
}
