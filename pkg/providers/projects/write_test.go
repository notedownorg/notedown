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

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/providers/projects"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	client, _ := buildClient(loadEvents(),
		test.Validators{
			Create: []test.CreateValidator{
				func(doc writer.Document, metadata reader.Metadata, content []byte, feed chan reader.Event) error {
					assert.Equal(t, writer.Document{Path: "projects/project.md"}, doc)
					assert.Equal(t, reader.Metadata{reader.MetadataTypeKey: projects.MetadataKey, projects.StatusKey: projects.Backlog, projects.NameKey: "project"}, metadata)
					assert.Equal(t, []byte("# project\n\n"), content)
					return nil
				},
			},
			MetadataUpdate: []test.MetadataUpdateValidator{
				func(doc writer.Document, metadata reader.Metadata) error {
					assert.Equal(t, writer.Document{Path: "projects/project.md"}, doc)
					assert.Equal(t, reader.Metadata{reader.MetadataTypeKey: projects.MetadataKey, projects.StatusKey: projects.Active, projects.NameKey: "project"}, metadata)
					return nil
				},
				// Used by rename
				func(doc writer.Document, metadata reader.Metadata) error {
					assert.Equal(t, writer.Document{Path: "projects/project.md"}, doc)
					assert.Equal(t, reader.Metadata{reader.MetadataTypeKey: projects.MetadataKey, projects.StatusKey: projects.Active, projects.NameKey: "new-project"}, metadata)
					return nil
				},
			},
			Rename: []test.RenameValidator{
				func(oldPath, newPath string) error {
					assert.Equal(t, "projects/project.md", oldPath)
					assert.Equal(t, "projects/new-project.md", newPath)
					return nil
				},
			},
			Delete: []test.DeleteValidator{
				func(doc writer.Document) error {
					assert.Equal(t, writer.Document{Path: "projects/project.md"}, doc)
					return nil
				},
			},
		},
	)
	project := projects.NewProject(projects.NewIdentifier("projects/project.md", ""), projects.WithStatus(projects.Active))
	assert.NoError(t, client.CreateProject("projects/project.md", "project", projects.Backlog))
	assert.NoError(t, client.UpdateProject(project))
	assert.NoError(t, client.RenameProject(project, "new-project"))
	assert.NoError(t, client.DeleteProject(projects.NewProject(projects.NewIdentifier("projects/project.md", ""))))
}
