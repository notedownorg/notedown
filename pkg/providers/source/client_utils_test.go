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

package source_test

import (
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/providers/source"
)

func buildClient(events []reader.Event, validators test.Validators) (*source.SourceClient, chan reader.Event) {
	feed := make(chan reader.Event)
	go func() {
		for _, event := range events {
			feed <- event
		}
	}()

	client := source.NewClient(
		&test.MockDocumentWriter{Validators: validators, Feed: feed},
		feed,
		source.WithInitialLoadWaiter(100*time.Millisecond),
	)
	return client, feed
}

func sourceCount(events []reader.Event) int {
	count := 0
	for _, event := range events {
		if event.Op == reader.Load && event.Document.Metadata[reader.MetadataTypeKey] == "source" {
			count++
		}
	}
	return count
}

var eventNotes = []source.Source{
	source.NewArticle(source.NewIdentifier("library/one.md", "version"), "example.com"),
	source.NewVideo(source.NewIdentifier("library/two.md", "version"), "example.com"),
	source.NewArticle(source.NewIdentifier("library/three.md", "version"), "example.com"),
	source.NewVideo(source.NewIdentifier("library/four.md", "version"), "example.com"),
	source.NewArticle(source.NewIdentifier("library/five.md", "version"), "example.com"),
	source.NewSource(source.NewIdentifier("library/six.md", "version"), source.Unknown),
}

func loadEvents() []reader.Event {
	return []reader.Event{
		// Sources format as string to reflect how it actually receives data in practice
		{
			Op:  reader.Load,
			Key: "library/one.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: source.MetadataKey,
					source.FormatKey:       string(source.Article),
					source.UrlKey:          "example.com",
				},
				Contents: []byte(`# One`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "library/two.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: source.MetadataKey,
					source.FormatKey:       string(source.Video),
					source.UrlKey:          "example.com",
				},
				Contents: []byte(`# Two`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "library/three.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: source.MetadataKey,
					source.FormatKey:       string(source.Article),
					source.UrlKey:          "example.com",
				},
				Contents: []byte(`# Three`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "library/four.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: source.MetadataKey,
					source.FormatKey:       string(source.Video),
					source.UrlKey:          "example.com",
				},
				Contents: []byte(`# Four`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "library/five.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: source.MetadataKey,
					source.FormatKey:       string(source.Article),
					source.UrlKey:          "example.com",
				},
				Contents: []byte(`# Five`),
				Checksum: "version",
			},
		},
		// No format set, we should prevent this where possible!
		// But theres nothing stopping someone from hand editing a file...
		{
			Op:  reader.Load,
			Key: "library/six.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: source.MetadataKey,
				},
				Contents: []byte(`# Six`),
				Checksum: "version",
			},
		},
		// Non-source note
		{
			Op:  reader.Load,
			Key: "someothernote.md",
			Document: reader.Document{
				Contents: []byte(`# Some other note`),
				Checksum: "version",
			},
		},
		// Load complete
		{
			Op: reader.SubscriberLoadComplete,
		},
	}
}
