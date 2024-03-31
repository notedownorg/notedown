package cache

import (
	"log/slog"
	"os"
	"time"
)

// Gargabe collection is required because we always load the entire existing cache regardless
// of whether or not the file still exists.
func (c *cache) runGarbageCollector() {
	ticker := time.NewTicker(60 * time.Second)
	c.done = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				c.gc()
			case <-c.done:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *cache) stopGarbageCollector() {
	close(c.done)
}

func (c *cache) gc() {
	slog.Debug("running garbage collector")
	for path := range c.Docs {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			delete(c.Docs, path)
			continue
		}
		if err != nil {
			slog.Error("error getting file info", slog.Any("error", err), slog.String("file", path))
			continue
		}
	}
	slog.Debug("finished garbage collection")
}
