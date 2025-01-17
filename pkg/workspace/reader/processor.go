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
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/notedownorg/notedown/pkg/workspace"
)

func (c *Client) processFile(path string, load bool) {
	slog.Debug("processing file", slog.String("file", path))
	// If we have already processed this file and it is up to date, we can skip it
	if c.isUpToDate(path) {
		slog.Debug("file is up to date, stopping short", slog.String("file", path))

		// To enable us to know when all files are loaded we must always report load events
		if load {
			rel, err := c.relative(path)
			if err != nil {
				c.errors <- fmt.Errorf("failed to get relative path: %w", err)
			}
			c.events <- Event{Op: Load, Document: c.documents[rel], Key: rel}
		}
		return
	}
	// Do the rest in a goroutine so we can continue doing other things
	c.threadLimit.Acquire(context.Background(), 1) // acquire semaphore as we will be making a blocking syscall
	go func() {
		slog.Debug("parsing file", slog.String("file", path))
		defer c.threadLimit.Release(1)

		rel, err := c.relative(path)
		if err != nil {
			slog.Error("failed to get relative path", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to get relative path: %w", err)
			return
		}

		d, err := workspace.LoadDocument(c.root, rel, time.Now())
		if err != nil {
			slog.Error("failed to parse document", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to parse document: %w", err)
			return
		}

		slog.Debug("updating document in cache", slog.String("file", path), slog.String("relative", rel))

		c.docMutex.Lock()
		c.documents[rel] = d
		c.docMutex.Unlock()
		op := func() Operation {
			if load {
				return Load
			} else {
				return Change
			}
		}()
		c.events <- Event{Op: op, Document: d, Key: rel}
	}()
}

func (c *Client) isUpToDate(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		slog.Error("Failed to get file info", slog.String("file", file), slog.String("error", err.Error()))
		c.errors <- fmt.Errorf("failed to get file info: %w", err)
		return false
	}
	rel, err := c.relative(file)
	if err != nil {
		slog.Error("Failed to get relative path", slog.String("file", file), slog.String("error", err.Error()))
		c.errors <- fmt.Errorf("failed to get relative path: %w", err)
		return false
	}
	c.docMutex.RLock()
	doc, ok := c.documents[rel]
	c.docMutex.RUnlock()
	return ok && doc.Modified(info.ModTime())
}
