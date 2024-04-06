package cache

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	    "github.com/liamawhite/nl/pkg/ast"
)

// We keep track of when documents were last modified so we can only re-parse them if they've changed
// This is much faster than hashing the file and comparing it to a hash because we don't have to read the file
type doc struct {
	LastModified time.Time     `json:"lastModified"`
	Data         *ast.Document `json:"data"`
}

// We dont have any custom logic for multiple client caches because the cache contents are deterministic
type cache struct {
	root string

	// Flusher debouncing... We debounce to ensure we are fully initialized before we start flushing
	lastUpdate time.Time
	done       chan struct{}

	Docs map[string]doc `json:"docs"`
}

type Cache interface {
	Get(path string) (*ast.Document, time.Time, bool)
	Set(path string, lastModified time.Time, d *ast.Document)
}

func cacheFile(root string) string {
	return filepath.Join(root, ".nl", "cache.json")
}

func NewCache(root string) Cache {
	emptyCache := &cache{Docs: make(map[string]doc), root: root}

	// Ensure the parent directory exists
	err := os.MkdirAll(filepath.Join(root, ".nl"), 0755)
	if err != nil {
		slog.Error("error creating cache directory", slog.Any("error", err))
		return emptyCache.start()
	}

	// Check if a cache exists
	if _, err := os.Stat(cacheFile(root)); err != nil {
		if !os.IsNotExist(err) {
			slog.Error("error checking for cache file", slog.Any("error", err))
		}
		return emptyCache.start()
	}

	// Load the cache
	data, err := os.ReadFile(cacheFile(root))
	if err != nil {
		slog.Error("error reading cache file", slog.Any("error", err))
		return emptyCache.start()
	}
	c := &cache{}
	err = json.Unmarshal(data, c)
	if err != nil {
		slog.Error("error unmarshalling cache file", slog.Any("error", err))
		return emptyCache.start()
	}
	c.root = root

	return c.start()
}

func (c *cache) start() *cache {
	c.runFlusher()
	c.runGarbageCollector()
	return c
}

func (c *cache) Get(path string) (*ast.Document, time.Time, bool) {
	d, ok := c.Docs[path]
	if !ok {
		return nil, time.Time{}, false
	}
	f, err := os.Stat(path)
	if err != nil {
		slog.Error("error getting file info", slog.Any("error", err))
		return nil, time.Time{}, false
	}
	if f.ModTime().Equal(d.LastModified) {
		return d.Data, d.LastModified, true
	}
	// If the file has been modified since we last parsed it, we need to re-parse it
	if f.ModTime().After(d.LastModified) {
		return nil, time.Time{}, false
	}
	// If we somehow got to the point where the file was last modified before we parsed it, we should error
	slog.Error("file was last modified before we parsed it, this shouldn't be possible", slog.Any("file", path), slog.Any("lastModified", d.LastModified), slog.Any("fileLastModified", f.ModTime()))
	return nil, time.Time{}, false
}

func (c *cache) Set(path string, lastModified time.Time, d *ast.Document) {
	c.Docs[path] = doc{LastModified: lastModified, Data: d}
	c.lastUpdate = time.Now()
}
