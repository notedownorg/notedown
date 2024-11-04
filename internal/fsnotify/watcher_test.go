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

package fsnotify_test

import (
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/notedownorg/notedown/internal/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/tjarratt/babble"
)

var babbler = babble.NewBabbler()

func randomFile(root string) string {
	var b strings.Builder
	b.WriteString(root)
	b.WriteString("/")
	for i := 0; i < rand.Intn(3); i++ {
		b.WriteString(babbler.Babble())
		b.WriteString("/")
	}
	b.WriteString(babbler.Babble())
    b.WriteString(".file")
	return b.String()
}

func TestRecursiveWatcher(t *testing.T) {
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug.Level()})))

	dir, _ := os.MkdirTemp("", "testrecursivewatcher")
	w, _ := fsnotify.NewRecursiveWatcher(dir)

	// What we want to test is that we get an accurate view of the filesystem based on the events we receive
	// This because events are non-deteministic even if you dont take ordering into account
	got := make(fileview)
	go tracker(t, w, got)

	// Do a bunch of things
	want := make(fileview)
	for i := 0; i < 10000; i++ {
		path := randomFile(dir)
		if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
			t.Fatal(err)
		}

		// Create
		content := babbler.Babble()
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		want.add(path)

		// Randomly update, rename or remove
		switch rand.Intn(3) {
		case 0: // Update
			slog.Debug("updating file", "path", path)
			content = babbler.Babble()
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
			want.add(path)
		case 1: // Rename
			newpath := randomFile(dir)
			slog.Debug("renaming file", "path", path, "newpath", newpath)
			if err := os.MkdirAll(filepath.Dir(newpath), 0777); err != nil {
				t.Fatal(err)
			}
			if err := os.Rename(path, newpath); err != nil {
				t.Fatal(err)
			}
			want.add(newpath)
			delete(want, path)
		case 2: // Remove
			slog.Debug("removing file", "path", path)
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			delete(want, path)
		}
	}

	// Wait for the tracker to catch up then compare the views
	time.Sleep(3 * time.Second)
	assert.Equal(t, want, got)
}

// Map of file paths to their content
type fileview map[string]string

func (f *fileview) add(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	(*f)[path] = string(data)
}

func tracker(t *testing.T, w *fsnotify.RecursiveWatcher, view fileview) {
	for {
		select {
		case event := <-w.Events():
			if event.Op.Has(fsnotify.Create) {
				view.add(event.Name)
			}
			if event.Op.Has(fsnotify.Remove) || event.Op.Has(fsnotify.Rename) {
				delete(view, event.Name)
			}
			if event.Op.Has(fsnotify.Write) {
				view.add(event.Name)
			}
		case err := <-w.Errors():
			t.Log(err)
		}
	}
}
