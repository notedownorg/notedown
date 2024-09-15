package reader

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liamawhite/nl/internal/fsnotify"
	"golang.org/x/sync/semaphore"
)

// The document client is responsible for maintaining a cache of parsed files from the workspace.
// These documents can either be updated directly by the client or other instances of the client.
type Client struct {
	root     string
	clientId string

	// Documents indexed by their relative path
	documents map[string]Document
	docMutex  sync.RWMutex

	watcher    *fsnotify.RecursiveWatcher
	processors sync.WaitGroup

	subscribers []chan Event

	// Everytime a goroutine makes a blocking syscall (in our case usually file i/o) it uses a new thread so to avoid
	// large workspaces exhausting the thread limit we use a semaphore to limit the number of concurrent goroutines
	threadLimit *semaphore.Weighted

	errors chan error
	events chan Event
}

func NewClient(root string, clientId string) (*Client, error) {
	watcher, err := fsnotify.NewRecursiveWatcher(root)
	if err != nil {
		return nil, err
	}

	client := &Client{
		root:        root,
		clientId:    clientId,
		documents:   make(map[string]Document),
		docMutex:    sync.RWMutex{},
		watcher:     watcher,
		subscribers: make([]chan Event, 0),
		threadLimit: semaphore.NewWeighted(1000), // Avoid exhausting golang max threads
		errors:      make(chan error),
		events:      make(chan Event),
	}

	go client.fileWatcher()
	go client.eventDispatcher()

	// Recurse through the root directory and process all the files to build the initial state
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, ".git") || strings.Contains(path, ".stversions") {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			client.processFile(path)
		}
		return nil
	})

	// Wait for all the processors to finish
	client.Wait()

	return client, nil
}

func (c *Client) absolute(relative string) string {
	return filepath.Join(c.root, relative)
}

func (c *Client) relative(absolute string) (string, error) {
	return filepath.Rel(c.root, absolute)
}
