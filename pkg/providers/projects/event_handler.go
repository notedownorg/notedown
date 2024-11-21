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
	"log/slog"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/pkg/traits"
)

func onLoad(c *ProjectClient) traits.EventHandler {
	return func(event reader.Event) {
		c.handleChanges(event)
		c.publisher.Events <- Event{Op: Load}
	}
}

func onChange(c *ProjectClient) traits.EventHandler {
	return func(event reader.Event) {
		c.handleChanges(event)
		c.publisher.Events <- Event{Op: Change}
	}
}

func onDelete(c *ProjectClient) traits.EventHandler {
	return func(event reader.Event) {
		c.notesMutex.Lock()
		delete(c.notes, event.Key)
		c.notesMutex.Unlock()
		c.publisher.Events <- Event{Op: Delete}
		slog.Debug("removed project", "path", event.Key)
	}
}

func (c *ProjectClient) handleChanges(event reader.Event) {
	if event.Document.Metadata.Type() != MetadataKey {
		return
	}

	// Handle metadata
	opts := make([]ProjectOption, 0)
	opts = append(opts, WithStatus(extractStatus(event.Key, event.Document.Metadata)))

	p := NewProject(NewIdentifier(event.Key, event.Document.Checksum))
	for _, opt := range opts {
		opt(&p)
	}

	c.notesMutex.Lock()
	c.notes[event.Key] = p
	c.notesMutex.Unlock()
	slog.Debug("added project", "path", event.Key)
}

func extractStatus(path string, metadata reader.Metadata) Status {
	if metadata[StatusKey] == nil {
		slog.Error("status key not found, defaulting to backlog", "path", path)
		return Backlog
	}

	str, ok := metadata[StatusKey].(string)
	if !ok {
		slog.Error("invalid status value, defaulting to backlog", "status", metadata[StatusKey], "path", path)
		return Backlog
	}

	status, ok := statusMap[str]
	if !ok {
		slog.Error("invalid status value, defaulting to backlog", "status", metadata[StatusKey], "path", path)
		return Backlog
	}

	return status
}
