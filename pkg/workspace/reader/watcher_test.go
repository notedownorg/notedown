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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDocuments_Client_Watcher(t *testing.T) {
	// Do the setup and ensure its correct
	dir, err := copyTestData(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(dir, "testclient")
	if err != nil {
		t.Fatal(err)
	}
	go ensureNoErrors(t, client.Errors())
	assert.Len(t, client.documents, 1)

	// Throw a bunch of events at the client and ensure the documents are updated correctly
	writeFile(dir, "1.md", "# Test Document 1") // doc count: 2
	writeFile(dir, "2.md", "# Test Document 2") // doc count: 3
	writeFile(dir, "3.md", "# Test Document 3") // doc count: 4

	// Do some updates
	writeFile(dir, "1.md", "# Test Document 1 Updated") // doc count: 4
	writeFile(dir, "2.md", "# Test Document 2 Updated") // doc count: 4

	// Do some deletes
	os.Remove(dir + "/3.md") // doc count: 3

	// As file watching has to be done async theres no way to deterministically wait for the events to be processed
	assert.Eventually(t, func() bool { return len(client.documents) == 3 }, time.Second, time.Millisecond*100, "expected %v documents got %v", 3, len(client.documents))
}

func TestDocuments_Client_Watcher_Fuzz(t *testing.T) {
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

	// Throw a bunch of events at the client and ensure the documents are updated correctly
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
		rel, err := client.relative(k)
		if err != nil {
			t.Fatal(err)
		}
		wantRel[rel] = true
	}

	// Add prexisting docs to the wantRel map
	for k := range prexistingDocs {
		wantRel[k] = true
	}

	// Wait for all files to finish processing
	assert.Eventually(t, func() bool { return len(wantRel) == len(client.documents) }, 5*time.Second, time.Millisecond*100, "expected %v documents got %v", len(wantRel), len(client.documents))

	// Ensure the documents paths are correct
	got := map[string]bool{}
	client.docMutex.Lock()
	for k := range client.documents {
		got[k] = true
	}
	client.docMutex.Unlock()
	assert.Equal(t, wantRel, got)
}
