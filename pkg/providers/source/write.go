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
	"fmt"
	"log/slog"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
)

func (c *SourceClient) CreateSource(path string, title string, format Format, url string, options ...SourceOption) error {
	options = append(options, WithUrl(url))

	source := NewSource(NewIdentifier(path, ""), title, format, options...)
	slog.Debug("creating source", "identifier", source.Identifier(), "source", source.String())

	metadata := reader.Metadata{
		reader.MetadataTypeKey: MetadataKey,
		TitleKey:               title,
		FormatKey:              format,
		UrlKey:                 url,
	}
	content := []byte(fmt.Sprintf("\n# %s\n\n", title))

	return c.writer.Create(path, metadata, content)
}

func (c *SourceClient) DeleteSource(source Source) error {
	slog.Debug("deleting source", "identifier", source.Identifier(), "source", source.String())
	return c.writer.Delete(writer.Document{Path: source.Path()})
}
