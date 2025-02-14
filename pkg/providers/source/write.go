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

package source

import (
	"log/slog"
	"path/filepath"

	"github.com/notedownorg/notedown/pkg/workspace"
)

func (c *SourceClient) CreateSource(title string, format Format, url string, options ...SourceOption) error {
	options = append(options, WithUrl(url))
	src := NewSource(title, format, options...)
	slog.Debug("creating source", "path", src.path)

	metadata := workspace.NewMetadata()
	metadata[workspace.MetadataTagsKey] = SourceTag(title)
	metadata[TitleKey] = title
	metadata[FormatKey] = format
	metadata[UrlKey] = url

	return c.writer.Create(workspace.NewDocument(filepath.Join(c.dir, title+".md"), metadata))
}

func (c *SourceClient) DeleteSource(source Source) error {
	slog.Debug("deleting source", "path", source.path)
	return c.writer.Delete(source.path)
}
