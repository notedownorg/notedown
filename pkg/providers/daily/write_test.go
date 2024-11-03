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

package daily_test

import (
	"testing"
	"time"

	"github.com/notedownorg/notedown/pkg/fileserver/reader"
	"github.com/notedownorg/notedown/pkg/fileserver/writer"
	"github.com/notedownorg/notedown/pkg/providers/daily"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	date := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	client, _ := buildClient(loadEvents(),
		// Ensure doesn't exist
		func(method string, doc writer.Document, metadata reader.Metadata, content []byte, feed chan reader.Event) error {
			path := "daily/2023-12-31.md"
			assert.Equal(t, "add", method)
			assert.Equal(t, writer.Document{Path: path}, doc)
			feed <- reader.Event{
				Op:       reader.Change,
				Key:      path,
				Document: reader.Document{Metadata: reader.Metadata{reader.MetadataTypeKey: daily.MetadataKey}},
			}
			return nil
		},

		// Create
		func(method string, doc writer.Document, metadata reader.Metadata, content []byte, feed chan reader.Event) error {
			assert.Equal(t, "add", method)
			assert.Equal(t, writer.Document{Path: "daily/2023-12-31.md"}, doc)
			return nil
		},
	)

	got, ok, err := client.Ensure(date.AddDate(0, 0, 1), time.Second) // 2024-01-01 already exists
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, daily.NewDaily(daily.NewIdentifier("daily/2024-01-01.md", "version")), got)

	got, ok, err = client.Ensure(date, time.Second)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, daily.NewDaily(daily.NewIdentifier("daily/2023-12-31.md", "")), got)

	assert.NoError(t, client.Create(date))
}
