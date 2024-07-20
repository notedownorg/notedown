package documents

import (
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	rel, err := c.relative(event.Name)
	if err != nil {
		slog.Error("failed to get relative path", slog.String("file", event.Name), slog.String("error", err.Error()))
		return
	}
	delete(c.documents, rel)
}

func (c *Client) handleRenameEvent(event fsnotify.Event) {
	slog.Debug("handling file rename event", slog.String("file", event.Name))
	// This should probably be a transaction?
	c.handleRemoveEvent(event)
	c.handleCreateEvent(event)
}

func (c *Client) handleWriteEvent(event fsnotify.Event) {
	slog.Debug("handling file write event", slog.String("file", event.Name))
	c.processFile(event.Name)
}
