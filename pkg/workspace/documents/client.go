package documents

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liamawhite/nl/internal/fsnotify"
	"github.com/liamawhite/nl/pkg/ast"
	"golang.org/x/sync/semaphore"
)

type document struct {
	lastUpdated int64
	document    ast.Document
}

// The document client is responsible for maintaining a cache of parsed files from the workspace.
// These documents can either be updated directly by the client or other instances of the client.
type Client struct {
	root     string
	clientId string

	// Documents indexed by their relative path
	documents map[string]document

	mutex      sync.RWMutex
	watcher    *fsnotify.RecursiveWatcher
	processors sync.WaitGroup

    // Everytime a goroutine makes a blocking syscall (in our case file i/o) it uses a new thread so to avoid
    // large workspaces exhausting the thread limit we use a semaphore to limit the number of concurrent goroutines
	threadLimit *semaphore.Weighted
}

func NewClient(root string, clientId string) (*Client, error) {
	watcher, err := fsnotify.NewRecursiveWatcher(root)
	if err != nil {
		return nil, err
	}

	client := &Client{
		root:        root,
		clientId:    clientId,
		documents:   make(map[string]document),
		mutex:       sync.RWMutex{},
		watcher:     watcher,
		threadLimit: semaphore.NewWeighted(1000), // Avoid exhausting golang max threads
	}

    go client.fileWatcher()

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
	client.WaitForProcessingCompletion()

	return client, nil
}

func (c *Client) absolute(relative string) string {
	return filepath.Join(c.root, relative)
}

func (c *Client) relative(absolute string) (string, error) {
	return filepath.Rel(c.root, absolute)
}
