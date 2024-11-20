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
	"sync"
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/providers/pkg/collections"
	"github.com/notedownorg/notedown/pkg/providers/pkg/traits"
)

// Use a type alias to hide the implementation details of the traits
type watcher = traits.Watcher
type publisher = traits.Publisher[Event]

type DocumentWriter interface {
	Add(path string, metadata reader.Metadata, content []byte) error
}

type ProjectClient struct {
	*watcher
	*publisher
	writer DocumentWriter
	dir    string

	// notes maps between file paths to notes it should ONLY be updated in response
	// to events from the docuuments client and should otherwise be read-only.
	notes      map[string]Project
	notesMutex sync.RWMutex
}

type clientOptions func(*ProjectClient)

// Inform NewClient to wait for the initial load to complete before returning
func WithInitialLoadWaiter(tick time.Duration) clientOptions {
	return func(client *ProjectClient) {
		traits.WithInitialLoadWaiter(client.watcher)(tick)
	}
}

func NewClient(writer DocumentWriter, feed <-chan reader.Event, opts ...clientOptions) *ProjectClient {
	client := &ProjectClient{
		notes:  make(map[string]Project),
		writer: writer,
		dir:    "projects",
	}

	client.publisher = traits.NewPublisher[Event]()
	client.watcher = traits.NewWatcher(feed, onLoad(client), onChange(client), onDelete(client))

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c *ProjectClient) Summary() int {
	c.notesMutex.RLock()
	defer c.notesMutex.RUnlock()
	return len(c.notes)
}

// Opts are applied in order so filters should be applied before sorters
func (c *ProjectClient) ListProjects(fetcher collections.Fetcher[ProjectClient, Project], opts ...collections.ListOption[Project]) []Project {
	return collections.List(c, fetcher, opts...)
}
