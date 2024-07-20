package workspace

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liamawhite/nl/internal/cache"
	"github.com/liamawhite/nl/internal/fsnotify"
	"github.com/liamawhite/nl/internal/persistor"
	"github.com/liamawhite/nl/pkg/ast"
)

// TODO: Make these configurable
const (
	projectsDir = "projects"
	dailyDir    = "daily"
)

func New(root string) (*Workspace, error) {
	if !filepath.IsAbs(root) {
		return nil, fmt.Errorf("%v is not an absolute path", root)
	}

	watcher, err := fsnotify.NewRecursiveWatcher(root)
	if err != nil {
		return nil, err
	}

	// Ensure the projects and daily directories exist
	for _, dir := range []string{projectsDir, dailyDir} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			return nil, err
		}
	}

	ws := &Workspace{
		root:    root,
		watcher: watcher,
		directories: directories{
			DailyNotes: filepath.Join(root, dailyDir),
			Projects:   filepath.Join(root, projectsDir),
		},
		tasks:             make(map[string]map[int]*Task),
		documents:         make(map[string]*document),
		mutex:             &sync.Mutex{},
		cache:             cache.NewCache(root),
		persistor:         persistor.NewPersistor(),
		processingCounter: NewAtomicCounter(),
		files:             make(chan string, 1000),
		docs:              make(chan docChan, 1000),
	}
	go ws.runProcessor() // Handles updating the cache
	go ws.runEventLoop() // Handles watching for file changes

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
			ws.processingCounter.Increment()
			ws.files <- path
		}
		return nil
	})

	// Wait for the initial state to be built
	ws.synchronize()
	slog.Debug("initial state built", slog.Int("files", len(ws.documents)))

	return ws, nil
}

type directories struct {
	DailyNotes string `json:"daily_notes"`
	Projects   string `json:"projects"`
}

type document struct {
	markers ast.Markers
}

type Workspace struct {
	root              string
	cache             cache.Cache
	persistor         *persistor.Persistor
	processingCounter *AtomicCounter

	directories directories

	// Map of tasks by the file path and the line number
	// This is so we can quickly update the task when a file changes
	// This is also the only truly unique identifier for a given task
	tasks map[string]map[int]*Task
	// Map of documents by the file path
	documents map[string]*document

	files chan string
	docs  chan docChan

	watcher *fsnotify.RecursiveWatcher
	mutex   *sync.Mutex
}
