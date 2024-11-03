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

package daily_test

import (
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/daily"
	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
)

func buildClient(events []reader.Event, validators ...test.LineWriterValidator) (*daily.Client, chan reader.Event) {
	feed := make(chan reader.Event)
	go func() {
		for _, event := range events {
			feed <- event
		}
	}()

	client := daily.NewClient(
		&test.MockLineWriter{Validators: validators},
		feed,
		daily.WithInitialLoadWaiter(100*time.Millisecond),
	)
	return client, feed
}

func dailyCount(events []reader.Event) int {
	count := 0
	for _, event := range events {
		if event.Op == reader.Load && event.Document.Metadata[reader.MetadataType] == "daily" {
			count++
		}
	}
	return count
}

var eventNotes = []daily.Daily{
	daily.NewDaily(daily.NewIdentifier("2024-01-01.md", "version")),
}

func loadEvents() []reader.Event {
	return []reader.Event{
		// Daily note
		{
			Op:  reader.Load,
			Key: "2024-01-01.md",
			Document: reader.Document{
				Metadata: reader.Metadata{reader.MetadataType: "daily"},
				Contents: []byte(`# 2024-01-01`),
				Checksum: "version",
			},
		},
		// Non-daily note
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
