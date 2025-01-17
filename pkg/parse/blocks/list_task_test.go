// Copyright 2025 Notedown Authors
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

package blocks_test

import (
	"testing"
	"time"

	. "github.com/liamawhite/parse/test"
	"github.com/notedownorg/notedown/pkg/parse/ast"
	. "github.com/notedownorg/notedown/pkg/parse/blocks"
	. "github.com/notedownorg/notedown/pkg/parse/test"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/rrule-go"
)

func TestListTaskLists(t *testing.T) {
	tests := []ParserTest[ast.Block]{
		{
			Name:          "Status Todo",
			Input:         "- [ ] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status Done (lowercase)",
			Input:         "- [x] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Done, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status Done (uppercase)",
			Input:         "- [X] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Done, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status Doing",
			Input:         "- [/] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Doing, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status blocked (lowercase)",
			Input:         "- [b] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Blocked, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status blocked (uppercase)",
			Input:         "- [B] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Blocked, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status abandoned (lowercase)",
			Input:         "- [a] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Abandoned, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status abandoned (uppercase)",
			Input:         "- [A] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Abandoned, "Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Status invalid",
			Input:         "- [z] Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Ul(Uli("", "-", " ", P("[z] Task"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task name with lots of random spaces",
			Input:         "- [ ]  Task  with  lots  of  spaces  ",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, " Task  with  lots  of  spaces  "))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with paragraph",
			Input:         "- [ ] Task\n  with paragraph",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task"), P("with paragraph"))),
			ExpectedOK:    true,
		},
		{
			Name:          "Nested tasklist",
			Input:         "- [ ] Task\n  - [ ] Nested Task",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task"), Tl(Tli("", "-", " ", T(Todo, "Nested Task"))))),
			ExpectedOK:    true,
		},

		// Fields
		{
			Name:          "Task with due date",
			Input:         "- [ ] Task due:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Due(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with short due date",
			Input:         "- [ ] Task d:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Due(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Tasks with due date on second task",
			Input:         "- [ ] Task 1\n- [ ] Task 2 due:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task 1")), Tli("", "-", " ", T(Todo, "Task 2", Due(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with scheduled date",
			Input:         "- [ ] Task scheduled:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Scheduled(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with short scheduled date",
			Input:         "- [ ] Task s:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Scheduled(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Tasks with scheduled date on second task",
			Input:         "- [ ] Task 1\n- [ ] Task 2 scheduled:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task 1")), Tli("", "-", " ", T(Todo, "Task 2", Scheduled(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with completed date",
			Input:         "- [ ] Task completed:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Completed(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Tasks with completed date on second task",
			Input:         "- [ ] Task 1\n- [ ] Task 2 completed:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task 1")), Tli("", "-", " ", T(Todo, "Task 2", Completed(2021, 1, 1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with priority",
			Input:         "- [ ] Task priority:1",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Priority(1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with short priority",
			Input:         "- [ ] Task p:1",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Priority(1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Tasks with priority 999",
			Input:         "- [ ] Task p:999",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Priority(999)))),
			ExpectedOK:    true,
		},

		{
			Name:          "Tasks with priority on second task",
			Input:         "- [ ] Task 1\n- [ ] Task 2 priority:1",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task 1")), Tli("", "-", " ", T(Todo, "Task 2", Priority(1)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every",
			Input:         "- [ ] Task every:day",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(daily, "every:day")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with short every",
			Input:         "- [ ] Task e:day",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(daily, "e:day")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every on second task",
			Input:         "- [ ] Task 1\n- [ ] Task 2 every:day",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task 1")), Tli("", "-", " ", T(Todo, "Task 2", Every(daily, "every:day")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every weekly",
			Input:         "- [ ] Task every:week",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(weekly, "every:week")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every monthly",
			Input:         "- [ ] Task every:month",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(monthly, "every:month")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every yearly",
			Input:         "- [ ] Task every:year",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(yearly, "every:year")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every weekday",
			Input:         "- [ ] Task every:weekday",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(weekday, "every:weekday")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every weekend",
			Input:         "- [ ] Task every:weekend",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(weekend, "every:weekend")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every monday",
			Input:         "- [ ] Task every:mon",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(monday, "every:mon")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every mon tue wed",
			Input:         "- [ ] Task every:mon tue wed",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(monTuesWed, "every:mon tue wed")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every tue thurs",
			Input:         "- [ ] Task every:tue thurs",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(tuesThurs, "every:tue thurs")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every two days",
			Input:         "- [ ] Task every:2 days",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(twoDays, "every:2 days")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every two weeks",
			Input:         "- [ ] Task every:2 weeks",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(twoWeeks, "every:2 weeks")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every 18 months",
			Input:         "- [ ] Task every:18 months",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(eighteenMonths, "every:18 months")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every 3 years",
			Input:         "- [ ] Task every:3 years",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(threeYears, "every:3 years")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every jan mar sept",
			Input:         "- [ ] Task every:jan mar sept",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(janMarSept, "every:jan mar sept")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every first and fifteenth of jan and sept",
			Input:         "- [ ] Task every:1st 15th jan sept",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(firstFifteenthJanSept, "every:1st 15th jan sept")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every first and fifteenth",
			Input:         "- [ ] Task every:1 15",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(firstFifteenth, "every:1 15")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Task with every first and thirtieth jan and feb",
			Input:         "- [ ] Task every:1 30 jan feb",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Every(firstThirtiethJanFeb, "every:1 30 jan feb")))),
			ExpectedOK:    true,
		},

		// Some edge cases
		{
			Name:          "conflicting short and long fields", // both end in d: so make sure theres no due date
			Input:         "- [ ] Task scheduled:2021-01-01 completed:2021-01-02",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Scheduled(2021, 1, 1), Completed(2021, 1, 2)))),
			ExpectedOK:    true,
		},
		{
			Name:          "Fields parse order",
			Input:         "- [ ] Task due:2021-01-01 scheduled:2021-01-02 completed:2021-01-03 priority:1 every:day",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Due(2021, 1, 1), Scheduled(2021, 1, 2), Completed(2021, 1, 3), Priority(1), Every(daily, "every:day")))),
			ExpectedOK:    true,
		},
		{
			Name:          "Fields reverse parse order",
			Input:         "- [ ] Task every:day priority:1 completed:2021-01-03 scheduled:2021-01-02 due:2021-01-01",
			Parser:        ListParser(ctx),
			ExpectedMatch: Tl(Tli("", "-", " ", T(Todo, "Task", Due(2021, 1, 1), Scheduled(2021, 1, 2), Completed(2021, 1, 3), Priority(1), Every(daily, "every:day")))),
			ExpectedOK:    true,
		},
	}

	RunTests(t, tests)
}

var (
	daily   = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.DAILY}
	weekly  = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.WEEKLY}
	monthly = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.MONTHLY}
	yearly  = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.YEARLY}
	weekday = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.WEEKLY, Byweekday: []rrule.Weekday{rrule.MO, rrule.TU, rrule.WE, rrule.TH, rrule.FR}}
	// we only need sat because next weekend from either day is still the saturday
	weekend               = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.WEEKLY, Byweekday: []rrule.Weekday{rrule.SA}}
	monday                = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.WEEKLY, Byweekday: []rrule.Weekday{rrule.MO}}
	monTuesWed            = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.WEEKLY, Byweekday: []rrule.Weekday{rrule.MO, rrule.TU, rrule.WE}}
	tuesThurs             = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.WEEKLY, Byweekday: []rrule.Weekday{rrule.TU, rrule.TH}}
	twoDays               = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.DAILY, Interval: 2}
	twoWeeks              = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.WEEKLY, Interval: 2}
	eighteenMonths        = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.MONTHLY, Interval: 18}
	threeYears            = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.YEARLY, Interval: 3}
	janMarSept            = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.YEARLY, Bymonth: []int{1, 3, 9}, Bymonthday: []int{1}}
	firstFifteenthJanSept = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.YEARLY, Bymonthday: []int{1, 15}, Bymonth: []int{1, 9}}
	firstFifteenth        = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.YEARLY, Bymonthday: []int{1, 15}}
	firstThirtiethJanFeb  = rrule.ROption{Dtstart: ctx.RelativeTo(), Freq: rrule.YEARLY, Bymonthday: []int{1, 30}, Bymonth: []int{1, 2}}
)

// Test that the options above are actually correct.
// It's not always easy to tell that the options do what you think they do.
func TestROptions(t *testing.T) {
	tests := []struct {
		name     string
		input    rrule.ROption
		expected []time.Time
		end      time.Time
	}{
		{
			name:  "daily",
			input: daily,
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
			name:  "weekly",
			input: weekly,
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
			name:  "monthly",
			input: monthly,
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
			name:  "yearly",
			input: yearly,
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
			name:  "weekday",
			input: weekday,
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
			name:  "weekend",
			input: weekend,
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
			name:  "monday",
			input: monday,
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
			name:  "monTuesWed",
			input: monTuesWed,
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
			name:  "tuesThurs",
			input: tuesThurs,
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
			name:  "twoDays",
			input: twoDays,
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
			name:  "twoWeeks",
			input: twoWeeks,
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
			name:  "eighteenMonths",
			input: eighteenMonths,
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
			name:  "threeYears",
			input: threeYears,
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
			name:  "janMarSept",
			input: janMarSept,
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
			name:  "firstFifteenthJanSept",
			input: firstFifteenthJanSept,
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
			name:  "firstFifteenth",
			input: firstFifteenth,
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
			name:  "firstThirtiethJanFeb",
			input: firstThirtiethJanFeb,
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
		t.Run(test.name, func(t *testing.T) {
			rr, err := rrule.NewRRule(test.input)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, rr.Between(relativeTo, test.end, true))
		})
	}

}
