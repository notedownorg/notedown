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
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/notedownorg/notedown/pkg/providers/pkg/traits"
	"github.com/notedownorg/notedown/pkg/workspace"
	"github.com/notedownorg/notedown/pkg/workspace/reader"
	"github.com/notedownorg/notedown/pkg/workspace/writer"
)

// Use a type alias to hide the implementation details of the traits
type watcher = traits.Watcher
type publisher = traits.Publisher[Event]

type DocumentWriter interface {
	Create(doc workspace.Document) error
	Delete(path string) error
}

var _ DocumentWriter = writer.Client{}

type SourceClient struct {
	*watcher
	*publisher

	writer DocumentWriter
	dir    string

	// sources maps file paths to source. It should ONLY be updated in response
	// to events from the docuuments client and should otherwise be read-only.
	sources      map[string]Source
	sourcesMutex sync.RWMutex

	// documents maps file paths to documents. It should ONLY be updated in response
	// to events from the documents client and should otherwise be read-only.
	docs      map[string]workspace.Document
	docsMutex sync.RWMutex
}

type clientOptions func(*SourceClient)

// Inform NewClient to wait for the initial load to complete before returning
func WithInitialLoadWaiter(tick time.Duration) clientOptions {
	return func(client *SourceClient) {
		traits.WithInitialLoadWaiter(client.watcher)(tick)
	}
}

func NewClient(writer DocumentWriter, feed <-chan reader.Event, opts ...clientOptions) *SourceClient {
	client := &SourceClient{
		sources: make(map[string]Source),
		docs:    make(map[string]workspace.Document),
		writer:  writer,
		dir:     "library",
	}

	client.publisher = traits.NewPublisher[Event]()
	client.watcher = traits.NewWatcher(feed, onLoad(client), onChange(client), onDelete(client))

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Where a new source should be created based on the name
func (c *SourceClient) NewSourceLocation(name string) string {
	if name == "" {
		return ""
	}
	if strings.ContainsAny(name, "./\\") {
		return ""
	}
	return filepath.Join(c.dir, fmt.Sprintf("%s.md", name))
}
