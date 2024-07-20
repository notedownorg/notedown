package documents

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDocuments_Client_FileEvents(t *testing.T) {
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
	assert.Eventually(t, func() bool { return len(client.documents) == 3 }, time.Second, time.Millisecond*100)
}

func TestDocuments_Client_FileEvents_Fuzz(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

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

	// Throw a bunch of events at the client and ensure the documents are updated correctly
	var wg sync.WaitGroup
	c, u, d := 0, 0, 0
	for i := 0; i < 1000; i++ {
		switch rand.Intn(3) {
		case 0:
			c++
			go createFile(&wg, dir, "# Test Document")
		case 1:
			u++
			go createThenUpdateFile(&wg, dir, "# Test Document Updated")
		case 2:
			d++
			go createThenDeleteFile(&wg, dir)
		}
	}
	wg.Wait()
	slog.Info("all file events have been sent", slog.Int("create", c), slog.Int("update", u), slog.Int("delete", d))
	assert.Eventually(t, func() bool { return len(client.documents) == 1+c+u }, time.Second, time.Millisecond*100)
}

func createFile(wg *sync.WaitGroup, dir string, content string) (string, error) {
	wg.Add(1)
	defer wg.Done()
	filename := uuid.New().String()
	path := fmt.Sprintf("%v/%v.md", dir, filename)
	return path, os.WriteFile(path, []byte(content), 0644)
}

func createThenDeleteFile(wg *sync.WaitGroup, dir string) error {
	wg.Add(1)
	defer wg.Done()
	content := "some random text"
	path, err := createFile(wg, dir, content)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func createThenUpdateFile(wg *sync.WaitGroup, dir string, content string) error {
	wg.Add(1)
	defer wg.Done()
	path, err := createFile(wg, dir, content)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte("some random updated text"), 0644)
}
