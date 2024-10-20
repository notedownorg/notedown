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

package cache

import (
	"log/slog"
	"os"
	"time"
)

// Gargabe collection is required because we always load the entire existing cache regardless
// of whether or not the file still exists.
func (c *cache) runGarbageCollector() {
	ticker := time.NewTicker(60 * time.Second)
	c.done = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				c.gc()
			case <-c.done:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *cache) stopGarbageCollector() {
	close(c.done)
}

func (c *cache) gc() {
	slog.Debug("running garbage collector")
	for path := range c.Docs {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			delete(c.Docs, path)
			continue
		}
		if err != nil {
			slog.Error("error getting file info", slog.Any("error", err), slog.String("file", path))
			continue
		}
	}
	slog.Debug("finished garbage collection")
}
