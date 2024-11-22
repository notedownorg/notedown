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

	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/stretchr/testify/assert"
)

func TestFetchAllProjects(t *testing.T) {
	events := loadEvents()
	c, _ := buildClient(events, test.Validators{})
	notes := c.ListProjects(projects.FetchAllProjects())
	assert.ElementsMatch(t, eventNotes, notes)
}
