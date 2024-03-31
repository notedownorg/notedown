package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/fsnotify"
	"github.com/liamawhite/nl/pkg/workspace/cache"
)

func New(root string) (*Workspace, error) {
	if !filepath.IsAbs(root) {
		return nil, fmt.Errorf("%v is not an absolute path", root)
	}

	watcher, err := fsnotify.NewRecursiveWatcher(root)
	if err != nil {
		return nil, err
	}

	ws := &Workspace{
		root:    root,
		watcher: watcher,
		tasks:   make(map[string]map[int]*ast.Task),
		mutex:   &sync.Mutex{},
		cache:   cache.NewCache(root),
		files:   make(chan string, 1000),
		docs:    make(chan docChan, 1000),
	}
	go ws.runProcessor()

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
			ws.files <- path
		}
		return nil
	})
	return ws, nil
}

type Workspace struct {
	root  string
	cache cache.Cache

	// Map of tasks by the file path and the line number
	// This is so we can quickly update the task when a file changes
	// This is also the only truly unique identifier for a given task
	tasks map[string]map[int]*ast.Task

	files chan string
	docs  chan docChan

	watcher *fsnotify.RecursiveWatcher
	mutex   *sync.Mutex
}

func (w Workspace) Tasks() map[string]map[int]*ast.Task {
	return w.tasks
}
