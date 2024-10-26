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
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/notedownorg/notedown/pkg/parsers"
)

func (c *Client) processFile(path string, load bool) {
	slog.Debug("processing file", slog.String("file", path))
	// If we have already processed this file and it is up to date, we can skip it
	if c.isUpToDate(path) {
		slog.Debug("file is up to date, stopping short", slog.String("file", path))
		return
	}

	// Do the rest in a goroutine so we can continue doing other things
	c.threadLimit.Acquire(context.Background(), 1) // acquire semaphore as we will be making a blocking syscall
	go func() {
		slog.Debug("parsing file", slog.String("file", path))
		defer c.threadLimit.Release(1)
		contents, err := os.ReadFile(path)
		if err != nil {
			slog.Error("failed to read file", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to read file: %w", err)
			return
		}
		rel, err := c.relative(path)
		if err != nil {
			slog.Error("failed to get relative path", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to get relative path: %w", err)
			return
		}

		// Calculate and set the hash of the contents to use as the version
		hash := sha256.New()
		hash.Write(contents)
		version := fmt.Sprintf("%x", hash.Sum(nil))

		d, err := parsers.Document(rel, version, time.Now())(string(contents))
		if err != nil {
			slog.Error("failed to parse document", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to parse document: %w", err)
			return
		}

		slog.Debug("updating document in cache", slog.String("file", path), slog.String("relative", rel))
		doc := Document{Document: d, lastUpdated: time.Now().Unix()}

		c.docMutex.Lock()
		c.documents[rel] = doc
		c.docMutex.Unlock()
		op := func() Operation {
			if load {
				return Load
			} else {
				return Change
			}
		}()
		c.events <- Event{Op: op, Document: doc, Key: rel}
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
	return ok && doc.lastUpdated >= info.ModTime().Unix()
}
