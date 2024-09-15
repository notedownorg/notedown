package documents

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/liamawhite/nl/internal/fsnotify"
)

func (c *Client) fileWatcher() {
	defer c.watcher.Close()
	for {
		select {
		case event := <-c.watcher.Events():
			switch event.Op {
			case fsnotify.Create:
				c.handleCreateEvent(event)
			case fsnotify.Remove:
				c.handleRemoveEvent(event)
			case fsnotify.Rename:
				c.handleRenameEvent(event)
			case fsnotify.Write:
				c.handleWriteEvent(event)
			}
		case err := <-c.watcher.Errors():
			log.Printf("error: %s", err)
		}
	}
}

func (c *Client) handleCreateEvent(event fsnotify.Event) {
	slog.Debug("handling file create event", slog.String("file", event.Name))
	c.processFile(event.Name)
}

func (c *Client) handleRemoveEvent(event fsnotify.Event) {
	slog.Debug("handling file remove event", slog.String("file", event.Name))
	rel, err := c.relative(event.Name)
	if err != nil {
		slog.Error("failed to get relative path", slog.String("file", event.Name), slog.String("error", err.Error()))
		c.errors <- fmt.Errorf("failed to get relative path: %w", err)
		return
	}
	c.docMutex.Lock()
	defer c.docMutex.Unlock()
	delete(c.documents, rel)
	c.events <- Event{Op: Delete, Document: Document{}, Key: rel}
}

func (c *Client) handleRenameEvent(event fsnotify.Event) {
	slog.Debug("handling file rename event", slog.String("file", event.Name))
	c.handleRemoveEvent(event) // rename sends the name of the old file, presumably it sends a create event for the new file
}

func (c *Client) handleWriteEvent(event fsnotify.Event) {
	slog.Debug("handling file write event", slog.String("file", event.Name))
	c.processFile(event.Name)
}
