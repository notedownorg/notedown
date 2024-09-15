package documents

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/liamawhite/nl/pkg/parsers"
)

func (c *Client) processFile(path string) {
	slog.Debug("processing file", slog.String("file", path))
	// If we have already processed this file and it is up to date, we can skip it
	if c.isUpToDate(path) {
		slog.Debug("file is up to date, stopping short", slog.String("file", path))
		return
	}

	// Do the rest in a goroutine so we can continue doing other things
	c.processors.Add(1)
	c.threadLimit.Acquire(context.Background(), 1) // acquire semaphore as we will be making a blocking syscall
	go func() {
		slog.Debug("parsing file", slog.String("file", path))
		defer c.processors.Done()
		defer c.threadLimit.Release(1)
		contents, err := os.ReadFile(path)
		if err != nil {
			slog.Error("failed to read file", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to read file: %w", err)
			return
		}
		d, err := parsers.Document(time.Now())(string(contents))
		if err != nil {
			slog.Error("failed to parse document", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to parse document: %w", err)
			return
		}
		rel, err := c.relative(path)
		if err != nil {
			slog.Error("failed to get relative path", slog.String("file", path), slog.String("error", err.Error()))
			c.errors <- fmt.Errorf("failed to get relative path: %w", err)
			return
		}

		// Work out the hash of the contents
		hash := sha256.New()
		hash.Write(contents)
		hashSum := hash.Sum(nil)
		hashStr := hex.EncodeToString(hashSum)

		slog.Debug("updating document in cache", slog.String("file", path), slog.String("relative", rel))
		doc := Document{Document: d, lastUpdated: time.Now().Unix(), Hash: hashStr}

		c.docMutex.Lock()
		c.documents[rel] = doc
		c.docMutex.Unlock()
		c.events <- Event{Op: Change, Document: doc, Key: rel}
	}()
}

// Wait for all files to finish processing
func (c *Client) Wait() {
	c.processors.Wait()
}

func (c *Client) isUpToDate(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		slog.Error("Failed to get file info", slog.String("file", file), slog.String("error", err.Error()))
		c.errors <- fmt.Errorf("failed to get file info: %w", err)
		return false
	}
	rel, err := c.relative(file)
	if err != nil {
		slog.Error("Failed to get relative path", slog.String("file", file), slog.String("error", err.Error()))
		c.errors <- fmt.Errorf("failed to get relative path: %w", err)
		return false
	}
	c.docMutex.RLock()
	doc, ok := c.documents[rel]
	c.docMutex.RUnlock()
	return ok && doc.lastUpdated >= info.ModTime().Unix()
}
