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

package reader

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/notedownorg/notedown/internal/fsnotify"
	"github.com/notedownorg/notedown/pkg/configuration"
	"github.com/notedownorg/notedown/pkg/workspace"
	"golang.org/x/sync/semaphore"
)

// The document client is responsible for maintaining a cache of parsed files from the workspace.
// These documents can either be updated directly by the client or other instances of the client.
type Client struct {
	root string

	// Documents indexed by their relative path
	documents map[string]workspace.Document
	docMutex  sync.RWMutex

	watcher *fsnotify.RecursiveWatcher

	subscribers []chan Event

	// Everytime a goroutine makes a blocking syscall (in our case usually file i/o) it uses a new thread so to avoid
	// large workspaces exhausting the thread limit we use a semaphore to limit the number of concurrent goroutines
	threadLimit *semaphore.Weighted

	errors chan error
	events chan Event
}

func NewClient(ws *configuration.Workspace, application string) (*Client, error) {
	ignoredDirs := []string{".git", ".vscode", ".debug", ".stversions", ".stfolder"}
	watcher, err := fsnotify.NewRecursiveWatcher(ws.Location, fsnotify.WithIgnoredDirs(ignoredDirs))
	if err != nil {
		return nil, err
	}

	client := &Client{
		root:        ws.Location,
		documents:   make(map[string]workspace.Document),
		docMutex:    sync.RWMutex{},
		watcher:     watcher,
		subscribers: make([]chan Event, 0),
		threadLimit: semaphore.NewWeighted(1000), // Avoid exhausting golang max threads
		errors:      make(chan error),
		events:      make(chan Event),
	}

	// Create a subscription so we can listen for the initial load events
	sub := make(chan Event)
	subscriberIndex := client.Subscribe(sub)

	// For each file we process on intial load, a load event is emitted
	// Therefore if our subscriber has received a load event for each file we have finished the initial load
	var wg sync.WaitGroup
	go func() {
		for ev := range sub {
			if ev.Op == Load {
				wg.Done()
			}
		}
	}()

	go client.fileWatcher()
	go client.eventDispatcher()

	// Recurse through the root directory and process all the files to build the initial state
	slog.Debug("walking workspace to build initial state")
	err = filepath.Walk(client.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		for _, ignoredDir := range ignoredDirs {
			if strings.Contains(path, ignoredDir) {
				return nil
			}
		}
		if strings.HasSuffix(path, ".md") {
			wg.Add(1) // Increment the wait group for each file we process
			client.processFile(path, true)
		}
		return nil
	})

	// Wait for all initial loads to finish, unsubscribe and close the channel
	slog.Debug("waiting for initial load to complete")
	wg.Wait()
	client.Unsubscribe(subscriberIndex)
	close(sub)

	return client, nil
}

func (c *Client) absolute(relative string) string {
	return filepath.Join(c.root, relative)
}

func (c *Client) relative(absolute string) (string, error) {
	return filepath.Rel(c.root, absolute)
}
