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
	"fmt"
	"path/filepath"
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
)

// Set wait to 0 to not wait for the file to appear in the cache
// Otherwise ensure will error if the file does not appear in the cache within the wait duration
func (c *Client) Ensure(date time.Time, wait time.Duration) (Daily, bool, error) {
	// O(n) but probably fine
	// Unless humans achieve immortality or pre-emptively generate daily notes assuming they will live forever...
	matches := c.ListDailyNotes(FetchAllNotes(), WithFilters(FilterByDate(&date, &date)))
	if len(matches) > 0 {
		return matches[0], true, nil
	}
	if err := c.Create(date); err != nil {
		return Daily{}, false, fmt.Errorf("failed to create daily note: %w", err)
	}

	if wait == 0 {
		return Daily{}, false, nil
	}

	// If wait is set we wait for the file to appear in our cache before returning
	start := time.Now()
	for {
		if time.Since(start) > wait {
			return Daily{}, false, fmt.Errorf("timed out waiting for daily note to appear")
		}
		matches = c.ListDailyNotes(FetchAllNotes(), WithFilters(FilterByDate(&date, &date)))
		if len(matches) > 0 {
			return matches[0], true, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (c *Client) Create(date time.Time) error {
	name := date.Format("2006-01-02")
	path := filepath.Join("daily", fmt.Sprintf("%s.md", name))
	return c.writer.AddDocument(path, reader.Metadata{reader.MetadataTypeKey: MetadataKey}, []byte{})
}
