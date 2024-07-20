package documents

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/liamawhite/nl/pkg/parsers"
)

func (c *Client) processFile(path string) {
	// If we have already processed this file and it is up to date, we can skip it
	if c.isUpToDate(path) {
		return
	}

	// Do the rest in a goroutine so we can continue doing other things
	c.processors.Add(1)
	c.threadLimit.Acquire(context.Background(), 1) // acquire semaphore as we will be making a blocking syscall
	go func() {
		defer c.processors.Done()
		defer c.threadLimit.Release(1)
		contents, err := os.ReadFile(path)
		if err != nil {
			slog.Error("failed to read file", slog.String("file", path), slog.String("error", err.Error()))
			return
		}
		d, err := parsers.Document(time.Now())(string(contents))
		if err != nil {
			slog.Error("failed to parse document", slog.String("file", path), slog.String("error", err.Error()))
			return
		}
		rel, err := c.relative(path)
		if err != nil {
			slog.Error("failed to get relative path", slog.String("file", path), slog.String("error", err.Error()))
			return
		}

		c.mutex.Lock()
		slog.Debug("updating document in cache", slog.String("file", path), slog.String("relative", rel))
		c.documents[rel] = document{document: d, lastUpdated: time.Now().Unix()}
		c.mutex.Unlock()
	}()
}

// Wait for all files to finish processing
func (c *Client) WaitForProcessingCompletion() {
	c.processors.Wait()
}

func (c *Client) isUpToDate(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		slog.Error("Failed to get file info", slog.String("file", file), slog.String("error", err.Error()))
		return false
	}
	rel, err := c.relative(file)
	if err != nil {
		slog.Error("Failed to get relative path", slog.String("file", file), slog.String("error", err.Error()))
		return false
	}
	c.mutex.RLock()
	doc, ok := c.documents[rel]
	c.mutex.RUnlock()
	return ok && doc.lastUpdated >= info.ModTime().Unix()
}
