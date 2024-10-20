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
	"encoding/json"
	"log/slog"
	"os"
	"time"
)

func (c *cache) runFlusher() {
	ticker := time.NewTicker(5 * time.Second)
	c.done = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if time.Since(c.lastUpdate) > 5*time.Second {
					c.flush()
				}
			case <-c.done:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *cache) stopFlusher() {
	close(c.done)
}

func (c cache) flush() {
	slog.Debug("flushing cache")
	data, err := json.Marshal(c)
	if err != nil {
		slog.Error("error marshalling cache", slog.Any("error", err))
		return
	}
	err = os.WriteFile(cacheFile(c.root), data, 0644)
	if err != nil {
		slog.Error("error writing cache file", slog.Any("error", err))
	}
	slog.Debug("flushed cache")
}
