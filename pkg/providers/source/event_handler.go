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
	"strings"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/pkg/traits"
)

func onLoad(c *SourceClient) traits.EventHandler {
	return func(event reader.Event) {
		if c.handleChanges(event) {
			c.publisher.Events <- Event{Op: Load}
		}
	}
}

func onChange(c *SourceClient) traits.EventHandler {
	return func(event reader.Event) {
		if c.handleChanges(event) {
			c.publisher.Events <- Event{Op: Change}
		}
	}
}

func onDelete(c *SourceClient) traits.EventHandler {
	return func(event reader.Event) {
		c.notesMutex.Lock()
		delete(c.notes, event.Key)
		c.notesMutex.Unlock()
		c.publisher.Events <- Event{Op: Delete}
		slog.Debug("removed source", "path", event.Key)
	}
}

func (c *SourceClient) handleChanges(event reader.Event) bool {
	if event.Document.Metadata.Type() != MetadataKey {
		return false
	}

	title := extractTitle(event.Key, event.Document.Metadata)
	format := extractFormat(event.Key, event.Document.Metadata)
	url := extractUrl(event.Document.Metadata)
	p := NewSource(NewIdentifier(event.Key, event.Document.Checksum), title, format, WithUrl(url))

	c.notesMutex.Lock()
	c.notes[event.Key] = p
	c.notesMutex.Unlock()
	slog.Debug("added source", "path", event.Key)
	return true
}

func extractFormat(path string, metadata reader.Metadata) Format {
	if metadata[FormatKey] == nil {
		slog.Error("format key not found", "path", path)
		return Unknown
	}

	str, ok := metadata[FormatKey].(string)
	if !ok {
		slog.Error("invalid format type", "format", metadata[FormatKey], "path", path)
		return Unknown
	}

	format, ok := formatMap[str]
	if !ok {
		slog.Error("invalid format value", "format", metadata[FormatKey], "path", path)
		return Unknown
	}

	return format
}

func extractUrl(metadata reader.Metadata) string {
	if metadata[UrlKey] == nil {
		return ""
	}

	str, ok := metadata[UrlKey].(string)
	if !ok {
		slog.Error("invalid url type", "url", metadata[UrlKey])
		return ""
	}

	return str
}

func extractTitle(path string, title reader.Metadata) string {
	if title[TitleKey] == nil {
		slog.Error("title key not found, defaulting to filename", "path", path)
		return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	str, ok := title[TitleKey].(string)
	if !ok {
		slog.Error("invalid title type", "title", title[TitleKey])
		return ""
	}

	return str
}
