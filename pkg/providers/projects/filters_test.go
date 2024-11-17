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

package projects_test

import (
	"testing"

	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/stretchr/testify/assert"
)

func TestTaskFilters(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events)

	tests := []struct {
		name      string
		filter    collections.Filter[projects.Project]
		wantNotes []projects.Project
	}{
		{
			name:      "Filter by status",
			filter:    projects.FilterByStatus(projects.Active),
			wantNotes: []projects.Project{eventNotes[0]},
		},
		{
			name:      "Filter by multiple statuses",
			filter:    projects.FilterByStatus(projects.Active, projects.Archived),
			wantNotes: []projects.Project{eventNotes[0], eventNotes[3]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.wantNotes, c.ListProjects(projects.FetchAllProjects(), projects.WithFilter(tt.filter)))
		})
	}
}
