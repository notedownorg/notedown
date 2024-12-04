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
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/providers/projects"
)

func buildClient(events []reader.Event, validators test.Validators) (*projects.ProjectClient, chan reader.Event) {
	feed := make(chan reader.Event)
	go func() {
		for _, event := range events {
			feed <- event
		}
	}()

	client := projects.NewClient(
		&test.MockDocumentWriter{Validators: validators, Feed: feed},
		feed,
		projects.WithInitialLoadWaiter(100*time.Millisecond),
	)
	return client, feed
}

func projectsCount(events []reader.Event) int {
	count := 0
	for _, event := range events {
		if event.Op == reader.Load && event.Document.Metadata[reader.MetadataTypeKey] == "project" {
			count++
		}
	}
	return count
}

var eventNotes = []projects.Project{
	projects.NewProject(projects.NewIdentifier("projects/one.md", "version"), projects.WithStatus(projects.Active)),
	projects.NewProject(projects.NewIdentifier("projects/two.md", "version"), projects.WithStatus(projects.Backlog)),
	projects.NewProject(projects.NewIdentifier("projects/three.md", "version"), projects.WithStatus(projects.Abandoned)),
	projects.NewProject(projects.NewIdentifier("projects/four.md", "version"), projects.WithStatus(projects.Archived)),
	projects.NewProject(projects.NewIdentifier("projects/five.md", "version"), projects.WithStatus(projects.Blocked)),
	projects.NewProject(projects.NewIdentifier("projects/six.md", "version"), projects.WithStatus(projects.Backlog)),
}

func loadEvents() []reader.Event {
	return []reader.Event{
		// Projects status as string to reflect how it actually receives data in practice
		{
			Op:  reader.Load,
			Key: "projects/one.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: projects.MetadataKey,
					projects.StatusKey:     string(projects.Active),
                    projects.NameKey:       "one",
				},
				Contents: []byte(`# One`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "projects/two.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: projects.MetadataKey,
					projects.StatusKey:     string(projects.Backlog),
                    projects.NameKey:       "two",
				},
				Contents: []byte(`# Two`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "projects/three.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: projects.MetadataKey,
					projects.StatusKey:     string(projects.Abandoned),
                    projects.NameKey:       "three",
				},
				Contents: []byte(`# Three`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "projects/four.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: projects.MetadataKey,
					projects.StatusKey:     string(projects.Archived),
                    projects.NameKey:       "four",
				},
				Contents: []byte(`# Four`),
				Checksum: "version",
			},
		},
		{
			Op:  reader.Load,
			Key: "projects/five.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: projects.MetadataKey,
					projects.StatusKey:     string(projects.Blocked),
                    projects.NameKey:       "five",
				},
				Contents: []byte(`# Five`),
				Checksum: "version",
			},
		},
		// No status set, we should prevent this where possible!
		// But theres nothing stopping someone from hand editing a file...
		{
			Op:  reader.Load,
			Key: "projects/six.md",
			Document: reader.Document{
				Metadata: reader.Metadata{
					reader.MetadataTypeKey: projects.MetadataKey,
                    projects.NameKey:       "six",
				},
				Contents: []byte(`# Six`),
				Checksum: "version",
			},
		},
		// Non-project note
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
