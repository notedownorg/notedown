package cache

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"
)

func (c *cache) runFlusher() {
	ticker := time.NewTicker(5 * time.Second)
	c.done = make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if time.Since(c.lastUpdate) > 5*time.Second {
					c.flush()
				}
			case <-c.done:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *cache) stopFlusher() {
	close(c.done)
}

func (c cache) flush() {
	slog.Debug("flushing cache")
	data, err := json.Marshal(c)
	if err != nil {
		slog.Error("error marshalling cache", slog.Any("error", err))
		return
	}
	err = os.WriteFile(cacheFile(c.root), data, 0644)
	if err != nil {
		slog.Error("error writing cache file", slog.Any("error", err))
	}
	slog.Debug("flushed cache")
}
