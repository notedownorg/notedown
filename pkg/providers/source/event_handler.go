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

	"github.com/go-viper/mapstructure/v2"
	"github.com/notedownorg/notedown/pkg/providers/pkg/traits"
	"github.com/notedownorg/notedown/pkg/workspace/reader"
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
		c.sourcesMutex.Lock()
		delete(c.sources, event.Key)
		c.sourcesMutex.Unlock()
		c.docsMutex.Lock()
		delete(c.docs, event.Key)
		c.docsMutex.Unlock()
		c.publisher.Events <- Event{Op: Delete}
		slog.Debug("removed source", "path", event.Key)
	}
}

func (c *SourceClient) handleChanges(event reader.Event) bool {
	if event.Document.Metadata == nil || len(event.Document.Metadata) == 0 {
		return false
	}

	if event.Document.Metadata.Type() != MetadataKey {
		return false
	}

	var source Source
	if err := mapstructure.Decode(event.Document.Metadata, &source); err != nil {
		slog.Error("failed to decode frontmatter", "error", err)
		return false
	}
	source.path = event.Key

	c.sourcesMutex.Lock()
	c.sources[event.Key] = source
	c.sourcesMutex.Unlock()
	c.docsMutex.Lock()
	c.docs[event.Key] = event.Document
	c.docsMutex.Unlock()
	slog.Debug("added source", "path", event.Key)
	return true
}
