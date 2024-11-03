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

package daily

import (
	"sync"
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
	"github.com/notedownorg/notedown/pkg/providers/pkg/traits"
)

// Use a type alias to hide the implementation details of the traits
type watcher = traits.Watcher
type publisher = traits.Publisher[Event]

type Client struct {
	*watcher
	*publisher
	writer writer.DocumentWriter
	dir    string

	// notes maps between file paths to notes it should ONLY be updated in response
	// to events from the docuuments client and should otherwise be read-only.
	notes      map[string]Daily
	notesMutex sync.RWMutex
}

type clientOptions func(*Client)

// Inform NewClient to wait for the initial load to complete before returning
func WithInitialLoadWaiter(tick time.Duration) clientOptions {
	return func(client *Client) {
		traits.WithInitialLoadWaiter(client.watcher)(tick)
	}
}

func NewClient(writer writer.DocumentWriter, feed <-chan reader.Event, opts ...clientOptions) *Client {
	client := &Client{
		notes:  make(map[string]Daily),
		writer: writer,
		dir:    "daily",
	}

	client.publisher = traits.NewPublisher[Event]()
	client.watcher = traits.NewWatcher(feed, onLoad(client), onChange(client), onDelete(client))

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *Client) Summary() int {
	c.notesMutex.RLock()
	defer c.notesMutex.RUnlock()
	return len(c.notes)
}

// Opts are applied in order so filters should be applied before sorters
func (c *Client) ListDailyNotes(fetcher collections.Fetcher[Client, Daily], opts ...collections.ListOption[Daily]) []Daily {
	return collections.List(c, fetcher, opts...)
}
