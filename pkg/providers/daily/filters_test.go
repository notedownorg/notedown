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

package daily_test

import (
	"testing"

	"github.com/notedownorg/notedown/pkg/providers/daily"
	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
	"github.com/stretchr/testify/assert"
)

func TestTaskFilters(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name      string
		filter    collections.Filter[daily.Daily]
		wantNotes []daily.Daily
	}{
		{
			name:      "Filter by date",
			filter:    daily.FilterByDate(date(2024, 1, 2, 0), date(2024, 1, 3, -1)),
			wantNotes: []daily.Daily{eventNotes[1]},
		},
		{
			name:      "Filter by same before and after date",
			filter:    daily.FilterByDate(date(2024, 1, 2, 0), date(2024, 1, 2, 0)),
			wantNotes: []daily.Daily{eventNotes[1]},
		},
		{
			name:      "Filter by date is set using nil-nil",
			filter:    daily.FilterByDate(nil, nil),
			wantNotes: eventNotes,
		},
		{
			name:      "Filter by date before",
			filter:    daily.FilterByDate(nil, date(2024, 1, 2, 0)),
			wantNotes: []daily.Daily{eventNotes[0], eventNotes[1]},
		},
		{
			name:      "Filter by date after",
			filter:    daily.FilterByDate(date(2024, 1, 2, 0), nil),
			wantNotes: []daily.Daily{eventNotes[1], eventNotes[2]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.wantNotes, c.ListDailyNotes(daily.FetchAllNotes(), daily.WithFilter(tt.filter)))
		})
	}
}
