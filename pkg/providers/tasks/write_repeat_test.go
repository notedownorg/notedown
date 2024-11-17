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

	"github.com/stretchr/testify/assert"
)

func tsk(status Status, opts ...TaskOption) Task {
	return NewTask(
		NewIdentifier("", "", 1),
		"task",
		status,
		opts...,
	)
}

func ev(s string) Every {
	e, _ := NewEvery(s)
	return e
}

func TestNewForRepeat(t *testing.T) {
	type want struct {
		scheduled *time.Time
		due       *time.Time
	}

	ptr := func(d time.Time) *time.Time {
		return &d
	}

	now := date(2021, 1, 1)

	testCases := []struct {
		name     string
		input    Task
		want     want
		noRepeat bool
	}{
		{
			name:     "no repeat",
			input:    tsk(Todo),
			noRepeat: true,
		},
		{
			name:  "every day",
			input: tsk(Todo, WithEvery(ev("day"))),
			want: want{
				due: ptr(now.AddDate(0, 0, 1)),
			},
		},
		{
			name:  "every day with scheduled",
			input: tsk(Blocked, WithEvery(ev("day")), WithScheduled(now)),
			want: want{
				scheduled: ptr(now.AddDate(0, 0, 1)),
			},
		},
		{
			name:  "every three days",
			input: tsk(Doing, WithEvery(ev("3 days"))),
			want: want{
				due: ptr(now.AddDate(0, 0, 3)),
			},
		},
		{
			name:  "every week",
			input: tsk(Abandoned, WithEvery(ev("week"))),
			want: want{
				due: ptr(now.AddDate(0, 0, 7)),
			},
		},
		{
			name:  "every month",
			input: tsk(Todo, WithEvery(ev("month"))),
			want: want{
				due: ptr(now.AddDate(0, 1, 0)),
			},
		},
		{
			name:  "every year",
			input: tsk(Todo, WithEvery(ev("year"))),
			want: want{
				due: ptr(now.AddDate(1, 0, 0)),
			},
		},
		{
			name:  "every weekday", // 1st jan 2021 was a friday
			input: tsk(Todo, WithEvery(ev("weekday"))),
			want: want{
				due: ptr(now.AddDate(0, 0, 3)),
			},
		},
		{
			name:  "every weekend", // 1st jan 2021 was a friday
			input: tsk(Todo, WithEvery(ev("weekend"))),
			want: want{
				due: ptr(now.AddDate(0, 0, 1)),
			},
		},
		{
			name:  "every monday", // 1st jan 2021 was a friday
			input: tsk(Todo, WithEvery(ev("monday"))),
			want: want{
				due: ptr(now.AddDate(0, 0, 3)),
			},
		},
		{
			name:  "every friday", // 1st jan 2021 was a friday
			input: tsk(Todo, WithEvery(ev("friday"))),
			want: want{
				due: ptr(now.AddDate(0, 0, 7)),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			completed := NewTaskFromTask(tc.input, WithStatus(Done, now))
			got, repeat := newForRepeat(completed)

			assert.Equal(t, tc.noRepeat, !repeat)
			if tc.noRepeat {
				return
			}

			assert.NotNil(t, got.Every())
			assert.Equal(t, Todo, got.Status())
			assert.Equal(t, tc.want.scheduled, got.Scheduled())
			assert.Equal(t, tc.want.due, got.Due())
			assert.Nil(t, got.Completed())
		})
	}
}
