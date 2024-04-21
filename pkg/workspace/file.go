package workspace

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/liamawhite/nl/internal/fsnotify"
	"github.com/liamawhite/nl/pkg/ast"
	"github.com/liamawhite/nl/pkg/parsers"
)

type docChan struct {
	doc          *ast.Document
	file         string
	lastModified time.Time
}

func (w *Workspace) runProcessor() {
	for {
		select {
		case file := <-w.files:
			// If we have already parsed this file since the last time it was modified, we can just pass that along
			if d, lastModified, ok := w.cache.Get(file); ok {
				slog.Debug("using cached document", slog.String("file", file))
				w.docs <- docChan{doc: d, file: file, lastModified: lastModified}
				continue
			}

			slog.Debug("processing file", slog.String("file", file))
			fileInfo, err := os.Stat(file)
			if err != nil {
				slog.Error("error getting file info", slog.Any("error", err))
			}
			contents, err := os.ReadFile(file)
			if err != nil {
				slog.Error("error reading file", slog.Any("error", err), slog.String("file", file))
			}
			d, err := parsers.Document(time.Now())(string(contents))
			if err != nil {
				slog.Error("error parsing file", slog.Any("error", err), slog.String("file", file))
			}
			w.docs <- docChan{doc: &d, file: file, lastModified: fileInfo.ModTime()}

		case d := <-w.docs:
			slog.Debug("processing document", slog.String("file", d.file))
			tasks := map[int]*Task{}
			for i := range d.doc.Tasks {
				task := d.doc.Tasks[i]
				project := ""
				if typ, ok := d.doc.Metadata["type"].(string); ok {
					if typ == "project" {
						project = strings.ReplaceAll(filepath.Base(d.file), filepath.Ext(d.file), "")
					}
				}
				// Use paths relative to the workspace root in Ids to maintain cache portability
				rel, err := filepath.Rel(w.root, d.file)
				if err != nil {
					slog.Error("error getting relative path", slog.Any("error", err), slog.String("file", d.file))
				}
				tasks[task.Line] = &Task{
					id:        fmt.Sprintf("%s:%d", rel, task.Line),
					Name:      task.Name,
					Status:    Status(task.Status),
					Due:       task.Due,
					Scheduled: task.Scheduled,
					Completed: task.Completed,
					Priority:  task.Priority,
					Every:     task.Every,
					Project:   project,
				}
			}
			w.mutex.Lock()
			w.tasks[d.file] = tasks
			w.mutex.Unlock()
			w.cache.Set(d.file, d.lastModified, d.doc)
		}
	}

}

func (w *Workspace) runEventLoop() {
	defer w.watcher.Close()
	for {
		select {
		case event := <-w.watcher.Events():
			switch event.Op {
			case fsnotify.Create:
				w.handleCreateEvent(event)
			case fsnotify.Remove:
				w.handleRemoveEvent(event)
			case fsnotify.Rename:
				w.handleRenameEvent(event)
			case fsnotify.Write:
				w.handleWriteEvent(event)
			}
		case err := <-w.watcher.Errors():
			log.Printf("error: %s", err)
		}
	}
}

func (w *Workspace) handleCreateEvent(event fsnotify.Event) {
	slog.Debug("handling file create event", slog.String("file", event.Name))
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.tasks[event.Name] = make(map[int]*Task)
}

func (w *Workspace) handleRemoveEvent(event fsnotify.Event) {
	slog.Debug("handling file remove event", slog.String("file", event.Name))
	w.mutex.Lock()
	defer w.mutex.Unlock()
	delete(w.tasks, event.Name)
}

func (w *Workspace) handleRenameEvent(event fsnotify.Event) {
	slog.Debug("handling file rename event", slog.String("file", event.Name))
	// TODO implement when I have a better understanding of how to handle this
}

func (w *Workspace) handleWriteEvent(event fsnotify.Event) {
	slog.Debug("handling file write event", slog.String("file", event.Name))
	w.files <- event.Name
}
