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

package projects

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
)

func (c *ProjectClient) CreateProject(path string, name string, status Status, options ...ProjectOption) error {
	options = append(options, WithStatus(status))
	project := NewProject(NewIdentifier(path, ""), options...)
	slog.Debug("creating project", "identifier", project.Identifier().String(), "project", project.String())

	metadata := reader.Metadata{reader.MetadataTypeKey: MetadataKey, StatusKey: status, NameKey: name}
	contents := []byte(fmt.Sprintf("# %s\n\n", name))

	return c.writer.Create(path, metadata, contents)
}

func (c *ProjectClient) UpdateProject(project Project) error {
	slog.Debug("updating project", "identifier", project.Identifier().String(), "project", project.String())

	metadata := reader.Metadata{reader.MetadataTypeKey: MetadataKey, StatusKey: project.Status(), NameKey: project.Name()}
	return c.writer.UpdateMetadata(writer.Document{Path: project.Path()}, metadata)
}

func (c *ProjectClient) RenameProject(project Project, newName string) error {
	slog.Debug("renaming project", "oldName", project.Name(), "newName", newName, "path", project.Path())

	// Update the project first
	project = NewProjectFromProject(project, withName(newName))
	if err := c.UpdateProject(project); err != nil {
		return fmt.Errorf("failed to update project prior to file rename: %w", err)
	}

	// Then rename the file itself, we don't change directories just the file name
	currDir := filepath.Dir(project.Path())
	newPath := filepath.Join(currDir, fmt.Sprintf("%s.md", newName))
	return c.writer.Rename(project.Path(), newPath)
}

func (c *ProjectClient) DeleteProject(project Project) error {
	slog.Debug("deleting project", "name", project.Name(), "path", project.Path())
	return c.writer.Delete(writer.Document{Path: project.Path()})
}
