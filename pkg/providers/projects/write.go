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

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
)

func (c *ProjectClient) Create(path string, name string, status Status, options ...ProjectOption) error {
	options = append(options, WithStatus(status))
	project := NewProject(NewIdentifier(path, ""), options...)
	slog.Debug("creating project", "identifier", project.Identifier().String(), "project", project.String())

	metadata := reader.Metadata{reader.MetadataTypeKey: MetadataKey, StatusKey: status}
	contents := []byte(fmt.Sprintf("# %s\n\n", name))

	return c.writer.Add(path, metadata, contents)
}
