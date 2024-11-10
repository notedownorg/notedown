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
	"log"
	"log/slog"
	"os"

	"github.com/notedownorg/notedown/internal/fsnotify"
)

func (c *Client) fileWatcher() {
	defer c.watcher.Close()
	for {
		select {
		case event := <-c.watcher.Events():
			switch event.Op {
			case fsnotify.Create:
				c.handleCreateEvent(event)
			case fsnotify.Remove:
				c.handleRemoveEvent(event)
			case fsnotify.Rename:
				c.handleRenameEvent(event)
			case fsnotify.Write:
				c.handleWriteEvent(event)
			}
		case err := <-c.watcher.Errors():
			log.Printf("error: %s", err)
		}
	}
}

func isDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func (c *Client) handleCreateEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		slog.Debug("ignoring directory create event", slog.String("dir", event.Name))
		return
	}
	slog.Debug("handling file create event", slog.String("file", event.Name))
	c.processFile(event.Name, false)
}

func (c *Client) handleRemoveEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		slog.Debug("ignoring directory remove event", slog.String("dir", event.Name))
		return
	}
	slog.Debug("handling file remove event", slog.String("file", event.Name))
	rel, err := c.relative(event.Name)
	if err != nil {
		slog.Error("failed to get relative path", slog.String("file", event.Name), slog.String("error", err.Error()))
		c.errors <- fmt.Errorf("failed to get relative path: %w", err)
		return
	}
	c.docMutex.Lock()
	defer c.docMutex.Unlock()
	delete(c.documents, rel)
	c.events <- Event{Op: Delete, Document: Document{}, Key: rel}
}

func (c *Client) handleRenameEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		slog.Debug("ignoring directory rename event", slog.String("dir", event.Name))
		return
	}
	slog.Debug("handling file rename event", slog.String("file", event.Name))
	c.handleRemoveEvent(event) // rename sends the name of the old file, presumably it sends a create event for the new file
}

func (c *Client) handleWriteEvent(event fsnotify.Event) {
	if isDir(event.Name) {
		slog.Debug("ignoring directory write event", slog.String("dir", event.Name))
		return
	}
	slog.Debug("handling file write event", slog.String("file", event.Name))
	c.processFile(event.Name, false)
}
