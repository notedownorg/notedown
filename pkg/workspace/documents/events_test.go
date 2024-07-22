package documents

import (
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDocuments_Client_Events_Fuzz(t *testing.T) {
	// change to debug if you want to see the events, too noisy to leave on permanently though
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

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
	sub1 := client.Subscribe()
	sub2 := client.Subscribe()

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
