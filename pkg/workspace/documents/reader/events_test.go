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
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDocuments_Client_Events_SubscribeWithInitialDocuments_Sync(t *testing.T) {
	// Do the setup and ensure its correct
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(dir, "testclient")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, client.documents, 1)
	go ensureNoErrors(t, client.Errors())

	// Create a subscriber and ensure it receives the load complete event
	sub := make(chan Event)
	done := false
	loaded := 0
	go func() {
		for ev := range sub {
			if ev.Op == SubscriberLoadComplete {
				done = true
			}
			if ev.Op == Load {
				loaded++
			}
		}
	}()

	client.Subscribe(sub, WithInitialDocuments())

	// Ensure we eventually receive the load complete event and that an event was received for each document
	waiter := func(d bool) func() bool { return func() bool { return done } }(done)
	assert.Eventually(t, waiter, 3*time.Second, time.Millisecond*200, "wg didn't finish in time")
	assert.Len(t, client.documents, loaded)
}

func TestDocuments_Client_Events_SubscribeWithInitialDocuments_Async(t *testing.T) {

	// Do the setup and ensure its correct
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(dir, "testclient")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, client.documents, 1)
	go ensureNoErrors(t, client.Errors())

	sub := make(chan Event)
	got := map[string]bool{}

	go func() {
		for {
			select {
			case ev := <-sub:
				if ev.Op == Load {
					got[ev.Key] = true
				}
			}
		}
	}()

	// Hook them up to the client and ensure we eventually receive all the initial documents
	client.Subscribe(sub, WithInitialDocuments())
	assert.Eventually(t, func() bool { return len(client.documents) == len(got) }, 3*time.Second, time.Millisecond*200, "sub finished with %v documents, expected %v", len(got), len(client.documents))
}

func TestDocuments_Client_Events_Fuzz(t *testing.T) {

	// Do the setup and ensure its correct
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(dir, "testclient")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, client.documents, 1)
	go ensureNoErrors(t, client.Errors())

	prexistingDocs := map[string]bool{}
	for k := range client.documents {
		prexistingDocs[k] = true
	}

	// Create two subscribers
	sub1 := make(chan Event)
	sub2 := make(chan Event)

	got1, got2 := map[string]bool{}, map[string]bool{}
	go func() {
		for {
			select {
			case ev := <-sub1:
				switch ev.Op {
				case Change:
					got1[ev.Key] = true
				case Delete:
					delete(got1, ev.Key)
				}
			case ev := <-sub2:
				switch ev.Op {
				case Change:
					got2[ev.Key] = true
				case Delete:
					delete(got2, ev.Key)
				}
			}
		}
	}()

	// Hook them up to the client
	client.Subscribe(sub1)
	client.Subscribe(sub2)

	// Throw a bunch of events at the client and ensure the subscribers are notified correctly
	errChan := make(chan error)
	go ensureNoErrors(t, errChan)
	wantAbs := map[string]bool{}
	wantRel := map[string]bool{}

	for i := 0; i < 1000; i++ {
		switch rand.Intn(4) {
		case 0:
			wantAbs[createFile(dir, "# Test Document", errChan)] = true
		case 1:
			wantAbs[createThenUpdateFile(dir, "# Test Document Updated", errChan)] = true
		case 2:
			createThenDeleteFile(dir, errChan)
		case 3:
			wantAbs[createThenRenameFile(dir, "# Test Document", errChan)] = true
		}
	}

	// We have to make the keys relative...
	for k := range wantAbs {
		rel, _ := client.relative(k)
		if err == nil {
			wantRel[rel] = true
		}
	}

	// To remove non-determinism we need to remove any pre-existing documents from the gots
	// Because of the way go schedules goroutines, we can't guarantee that the subscribers won't receive these events
	// but they dont actually matter for real use cases as the events are idempotent
	for k := range prexistingDocs {
		delete(got1, k)
		delete(got2, k)
	}

	// Wait until we have handled all the events
	assert.Eventually(t, func() bool { return len(wantRel) == len(got1) }, 3*time.Second, time.Millisecond*200, "sub1 finished with %v documents, expected %v", len(got1), len(wantAbs))
	assert.Eventually(t, func() bool { return len(wantRel) == len(got2) }, 3*time.Second, time.Millisecond*200, "sub2 finished with %v documents, expected %v", len(got2), len(wantAbs))

	// Check the subscribers got the expected events
	assert.Equal(t, wantRel, got1)
	assert.Equal(t, wantRel, got2)

}
