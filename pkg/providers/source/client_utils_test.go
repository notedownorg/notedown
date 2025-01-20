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
	"time"

	"github.com/notedownorg/notedown/pkg/configuration"
	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/workspace"
	"github.com/notedownorg/notedown/pkg/workspace/reader"
)

var workspaceConfig = &configuration.WorkspaceConfiguration{
	Sources: configuration.Sources{
		DefaultDirectory: "sources",
	},
}

func buildClient(events []reader.Event, validators test.Validators) (*SourceClient, chan reader.Event) {
	feed := make(chan reader.Event)
	go func() {
		for _, event := range events {
			feed <- event
		}
	}()

	client := NewClient(
		workspaceConfig,
		&test.MockDocumentWriter{Validators: validators, Feed: feed},
		feed,
		WithInitialLoadWaiter(100*time.Millisecond),
	)
	return client, feed
}

func sourceCount(events []reader.Event) int {
	count := 0
	for _, event := range events {
		if event.Op == reader.Load && event.Document.Metadata[workspace.MetadataTypeKey] == "source" {
			count++
		}
	}
	return count
}

var eventNotes = []Source{
	{Title: "one", Format: Article, Url: "example.com", path: "library/one.md"},
	{Title: "two", Format: Video, Url: "example.com", path: "library/two.md"},
	{Title: "three", Format: Article, Url: "example.com", path: "library/three.md"},
	{Title: "four", Format: Video, Url: "example.com", path: "library/four.md"},
	{Title: "five", Format: Article, Url: "example.com", path: "library/five.md"},
	{Format: Unknown, Url: "", path: "library/six.md"},
}

func loadEvents() []reader.Event {
	return []reader.Event{
		// Sources format as string to reflect how it actually receives data in practice
		{
			Op:  reader.Load,
			Key: "library/one.md",
			Document: workspace.NewDocument(
				"library/one.md",
				workspace.Metadata{
					workspace.MetadataTypeKey: MetadataKey,
					TitleKey:                  "one",
					FormatKey:                 string(Article),
					UrlKey:                    "example.com",
				},
			),
		},
		{
			Op:  reader.Load,
			Key: "library/two.md",
			Document: workspace.NewDocument(
				"library/two.md",
				workspace.Metadata{
					workspace.MetadataTypeKey: MetadataKey,
					TitleKey:                  "two",
					FormatKey:                 string(Video),
					UrlKey:                    "example.com",
				},
			),
		},
		{
			Op:  reader.Load,
			Key: "library/three.md",
			Document: workspace.NewDocument(
				"library/three.md",
				workspace.Metadata{
					workspace.MetadataTypeKey: MetadataKey,
					TitleKey:                  "three",
					FormatKey:                 string(Article),
					UrlKey:                    "example.com",
				},
			),
		},
		{
			Op:  reader.Load,
			Key: "library/four.md",
			Document: workspace.NewDocument(
				"library/four.md",
				workspace.Metadata{
					workspace.MetadataTypeKey: MetadataKey,
					TitleKey:                  "four",
					FormatKey:                 string(Video),
					UrlKey:                    "example.com",
				},
			),
		},
		{
			Op:  reader.Load,
			Key: "library/five.md",
			Document: workspace.NewDocument(
				"library/five.md",
				workspace.Metadata{
					workspace.MetadataTypeKey: MetadataKey,
					TitleKey:                  "five",
					FormatKey:                 string(Article),
					UrlKey:                    "example.com",
				},
			),
		},
		// No format, url or title set, we should prevent this where possible!
		// But theres nothing stopping someone from hand editing a file...
		{
			Op:  reader.Load,
			Key: "library/six.md",
			Document: workspace.NewDocument(
				"library/six.md",
				workspace.Metadata{
					workspace.MetadataTypeKey: MetadataKey,
				},
			),
		},
		// Non-source note
		{
			Op:  reader.Load,
			Key: "someothernote.md",
			Document: workspace.NewDocument(
				"someothernote.md",
				nil,
			),
		},
		// Load complete
		{
			Op: reader.SubscriberLoadComplete,
		},
	}
}
