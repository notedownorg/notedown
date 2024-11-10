// Copyright 2024 Notedown Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fsnotify

import (
	"log/slog"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	root        string
	ignoredDirs []string

	w      *fsnotify.Watcher
	events chan fsnotify.Event
	errors chan error

	watchers map[string]struct{}
}

type Option func(*RecursiveWatcher)

func WithIgnoredDirs(dirs []string) Option {
	return func(rw *RecursiveWatcher) {
		rw.ignoredDirs = dirs
	}
}

func NewRecursiveWatcher(root string, opts ...Option) (*RecursiveWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	rw := &RecursiveWatcher{
		root:        root,
		ignoredDirs: []string{},
		w:           w,
		events:      make(chan fsnotify.Event),
		errors:      make(chan error),
		watchers:    make(map[string]struct{}),
	}

	for _, opt := range opts {
		opt(rw)
	}

	go rw.eventLoop()
	go rw.add(root) // run in a goroutine so we can return before events are being read

	return rw, nil
}

func (rw RecursiveWatcher) Events() <-chan fsnotify.Event {
	return rw.events
}

func (rw RecursiveWatcher) Errors() <-chan error {
	return rw.errors
}

func (rw *RecursiveWatcher) Close() error {
	return rw.w.Close()
}

func (rw *RecursiveWatcher) eventLoop() {
	for {
		select {
		case event := <-rw.w.Events:
			if event.Op.Has(fsnotify.Create) {
				slog.Debug("received create event", "path", event.Name)
				rw.handleCreate(event)
			}
			if event.Op.Has(fsnotify.Remove) {
				slog.Debug("received remove event", "path", event.Name)
				rw.handleRemove(event)
			}
			if event.Op.Has(fsnotify.Rename) {
				slog.Debug("received rename event", "path", event.Name)
				rw.handleRemove(event)
			}
			if event.Op.Has(fsnotify.Write) {
				slog.Debug("received write event", "path", event.Name)
				rw.events <- event
			}
		case err := <-rw.w.Errors:
			rw.errors <- err
		}
	}
}

func (rw *RecursiveWatcher) add(path string) {
	slog.Debug("processing path", "path", path)

	// Check if the path is already being watched
	if _, ok := rw.watchers[path]; ok {
		slog.Debug("path already being watched", "path", path)
		return
	}

	// Check the path does not match any ignored directories
	for _, ignoredDir := range rw.ignoredDirs {
		if strings.Contains(path, ignoredDir) {
			slog.Debug("path matches ignored directory", "path", path, "ignoredDir", ignoredDir)
			return
		}
	}

	// If the path is not a directory, return
	// See fsnotify docs as to why we dont bother watching individual files
	if !isDir(path) {
		slog.Debug("skipping path, not a directory", "path", path)
		return
	}

	// Add the path to the watcher
	if err := rw.w.Add(path); err != nil {
		rw.errors <- err
		return
	}

	// Add the path to the watchers map
	rw.watchers[path] = struct{}{}
	slog.Debug("watching path", "path", path)

	// Iterate over all the entries in the directory (subdirectories) and recurse
	entries, err := os.ReadDir(path)
	if err != nil {
		rw.errors <- err
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			rw.add(path + "/" + entry.Name())
		}
	}

	// This is a bit of a hack to ensure we get events when the following race condition occurs:
	// 1. A new directory is created
	// 2. A file is created in the new directory
	// 3. The directory is added to the watcher
	// 4. No event is sent for the file == sad face

	// Reload the entries
	entries, err = os.ReadDir(path)
	if err != nil {
		rw.errors <- err
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			slog.Debug("sending create event for file in new directory", "path", path+"/"+entry.Name())
			rw.events <- fsnotify.Event{Name: path + "/" + entry.Name(), Op: fsnotify.Create}
		}
	}
}

func isDir(path string) bool {
	fi, error := os.Stat(path)
	if error != nil {
		return false
	}
	return fi.IsDir()
}
