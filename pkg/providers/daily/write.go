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

func (c *Client) Ensure(date time.Time) error {
	// O(n) but probably fine
	// Unless humans achieve immortality or pre-emptively generate daily notes assuming they will live forever...
	matches := c.ListDailyNotes(FetchAllNotes(), WithFilters(FilterByDate(&date, &date)))
	if len(matches) > 0 {
		return nil
	}
	return c.Create(date)
}

func (c *Client) Create(date time.Time) error {
	name := date.Format("2006-01-02")
	path := filepath.Join("daily", fmt.Sprintf("%s.md", name))
	return c.writer.AddDocument(path, reader.Metadata{reader.MetadataTypeKey: MetadataKey}, []byte{})
}
