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
	"math/rand"
	"testing"
	"time"

	"github.com/notedownorg/notedown/pkg/providers/pkg/test"
	"github.com/notedownorg/notedown/pkg/workspace"
	"github.com/notedownorg/notedown/pkg/workspace/reader"
	"github.com/stretchr/testify/assert"
)

func TestEventBroadcast_Fuzz(t *testing.T) {
	c, feed := buildClient(
		loadEvents(),
		test.Validators{},
	)

	// Create two subscribers
	sub1 := make(chan Event)
	sub2 := make(chan Event)

	// Listen for events from the daily client
	got1, got2 := make([]Operation, 0), make([]Operation, 0)
	go func() {
		for {
			select {
			case event := <-sub1:
				got1 = append(got1, event.Op)
			case event := <-sub2:
				got2 = append(got2, event.Op)
			}
		}
	}()

	// Subscribe the listeners
	c.Subscribe(sub1)
	c.Subscribe(sub2)

	// Throw some events at the sources client and ensure we are notified correctly
	want := make([]Operation, 0)
	count := 1000
	d := workspace.Document{
		Metadata: workspace.Metadata{
			workspace.MetadataTagsKey: SourceTag("test"),
			TitleKey:                  "test",
			FormatKey:                 string(Article),
			UrlKey:                    "example.com"},
	}
	for i := 0; i < count; i++ {
		switch rand.Intn(3) {
		case 0:
			feed <- reader.Event{Op: reader.Load, Key: "test.md", Document: d}
			want = append(want, Load)
		case 1:
			feed <- reader.Event{Op: reader.Change, Key: "test.md", Document: d}
			want = append(want, Change)
		case 2:
			feed <- reader.Event{Op: reader.Delete, Key: "test.md", Document: d}
			want = append(want, Delete)
		}
	}

	// Ensure we received all events
	assert.Eventually(t, func() bool { return len(got1) == count }, 3*time.Second, 200*time.Millisecond)
	assert.Eventually(t, func() bool { return len(got2) == count }, 3*time.Second, 200*time.Millisecond)

	// Ensure the events match, note we need to check elements because the order is not guaranteed.
	assert.ElementsMatch(t, want, got1)
	assert.ElementsMatch(t, want, got2)
}
