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

package source

import (
	"testing"
	"time"

	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/workspace/reader"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	ch := make(chan reader.Event)
	events := loadEvents()
	go func() {
		for _, event := range events {
			ch <- event
		}
	}()

	client := NewClient(&workspaceConfig.Sources, &test.MockDocumentWriter{}, ch)

	// Assert that we eventually get the correct number of notes
	waitFor, tick := 3*time.Second, 200*time.Millisecond
	assert.Eventually(t, func() bool { return len(client.ListSources(FetchAllSources())) == sourceCount(events) }, waitFor, tick)
}

func TestClient_InitialLoadWaiter(t *testing.T) {
	ch := make(chan reader.Event)
	events := loadEvents()
	go func() {
		for _, event := range events {
			ch <- event
		}
		ch <- reader.Event{Op: reader.SubscriberLoadComplete}
	}()

	client := NewClient(&workspaceConfig.Sources, &test.MockDocumentWriter{}, ch, WithInitialLoadWaiter(100*time.Millisecond))

	// Assert that the client has the correct number of notes
	assert.Equal(t, sourceCount(events), len(client.ListSources((FetchAllSources()))))
}
