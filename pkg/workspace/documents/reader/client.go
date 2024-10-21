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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/notedownorg/notedown/internal/fsnotify"
	"golang.org/x/sync/semaphore"
)

type identifier struct {
	application string
	hostname    string
	pid         int
}

func (i identifier) String() string {
	return fmt.Sprintf("%s-%s-%d", i.application, i.hostname, i.pid)
}

func newIdentifier(application string) identifier {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = fmt.Sprintf("unknown_%d", time.Now().Unix())
	}
	pid := os.Getpid()
	return identifier{
		application: application,
		hostname:    hostname,
		pid:         pid,
	}
}

// The document client is responsible for maintaining a cache of parsed files from the workspace.
// These documents can either be updated directly by the client or other instances of the client.
type Client struct {
	root     string
	clientId identifier

	// Documents indexed by their relative path
	documents map[string]Document
	docMutex  sync.RWMutex

	watcher    *fsnotify.RecursiveWatcher
	processors sync.WaitGroup

	subscribers []chan Event

	// Everytime a goroutine makes a blocking syscall (in our case usually file i/o) it uses a new thread so to avoid
	// large workspaces exhausting the thread limit we use a semaphore to limit the number of concurrent goroutines
	threadLimit *semaphore.Weighted

	errors chan error
	events chan Event
}

func NewClient(root string, application string) (*Client, error) {
	watcher, err := fsnotify.NewRecursiveWatcher(root)
	if err != nil {
		return nil, err
	}

	client := &Client{
		root:        root,
		clientId:    newIdentifier(application),
		documents:   make(map[string]Document),
		docMutex:    sync.RWMutex{},
		watcher:     watcher,
		subscribers: make([]chan Event, 0),
		threadLimit: semaphore.NewWeighted(1000), // Avoid exhausting golang max threads
		errors:      make(chan error),
		events:      make(chan Event),
	}

	go client.fileWatcher()
	go client.eventDispatcher()

	// Recurse through the root directory and process all the files to build the initial state
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, ".git") || strings.Contains(path, ".stversions") {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			client.processFile(path)
		}
		return nil
	})

	// Wait for all the processors to finish
	client.Wait()

	return client, nil
}

func (c *Client) absolute(relative string) string {
	return filepath.Join(c.root, relative)
}

func (c *Client) relative(absolute string) (string, error) {
	return filepath.Rel(c.root, absolute)
}
